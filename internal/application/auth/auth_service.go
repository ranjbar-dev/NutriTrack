package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	userRepo "github.com/ranjbar-dev/nutritrack/internal/domain/user/repository"
	"github.com/ranjbar-dev/nutritrack/internal/domain/user/valueobject"
	"github.com/ranjbar-dev/nutritrack/internal/infrastructure/redis"
)

const otpLength = 6

// AuthService orchestrates all authentication flows.
type AuthService struct {
	userRepo    userRepo.UserRepository
	otpStore    userRepo.OTPRepository
	blacklist   *redis.TokenBlacklist
	jwtService  *JWTService
	smsProvider shared.SMSProvider
}

func NewAuthService(
	userRepo userRepo.UserRepository,
	otpStore userRepo.OTPRepository,
	blacklist *redis.TokenBlacklist,
	jwtService *JWTService,
	smsProvider shared.SMSProvider,
) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		otpStore:    otpStore,
		blacklist:   blacklist,
		jwtService:  jwtService,
		smsProvider: smsProvider,
	}
}

// Login authenticates superadmin/nutritionist with email + password.
func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		log.Error().Err(err).Msg("login: find user by email failed")
		return nil, shared.ErrInternal
	}
	if user == nil {
		return nil, shared.ErrInvalidCredentials
	}
	if !user.IsActive {
		return nil, shared.ErrForbidden
	}
	if user.IsClient() {
		// Clients use OTP, not password
		return nil, shared.ErrInvalidCredentials
	}
	if !CheckPassword(req.Password, user.PasswordHash) {
		return nil, shared.ErrInvalidCredentials
	}

	tokens, err := s.jwtService.GenerateTokenPair(user.ID, user.Role)
	if err != nil {
		log.Error().Err(err).Msg("login: token generation failed")
		return nil, shared.ErrInternal
	}

	return &AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TokenType:    "Bearer",
		UserID:       user.ID,
		Role:         user.Role,
	}, nil
}

// SendOTP sends a 6-digit OTP to the client's mobile number.
// Rate limited: max 3 sends per 10-minute window (atomic Redis INCR).
func (s *AuthService) SendOTP(ctx context.Context, req OTPSendRequest) error {
	mob, err := valueobject.NewMobile(req.Mobile)
	if err != nil {
		return err // already an *AppError
	}

	// Atomic rate limit check
	count, err := s.otpStore.IncrRateLimit(ctx, mob.String())
	if err != nil {
		log.Error().Err(err).Msg("sendOTP: incr rate limit failed")
		return shared.ErrInternal
	}
	if count > redis.MaxOTPRateLimit() {
		return shared.ErrOTPRateLimit
	}

	otp, err := shared.GenerateOTP(otpLength)
	if err != nil {
		log.Error().Err(err).Msg("sendOTP: generate OTP failed")
		return shared.ErrInternal
	}

	if err := s.otpStore.StoreOTP(ctx, mob.String(), otp); err != nil {
		log.Error().Err(err).Msg("sendOTP: store OTP failed")
		return shared.ErrInternal
	}

	if err := s.smsProvider.SendOTP(ctx, mob.Local(), otp); err != nil {
		log.Error().Err(err).Str("mobile", mob.String()).Msg("sendOTP: SMS send failed")
		// Don't leak SMS errors to caller — OTP is stored, client can retry verify
		return shared.ErrInternal
	}

	log.Info().Str("mobile", mob.String()).Msg("OTP sent")
	return nil
}

// VerifyOTP validates the OTP and returns JWT tokens.
// Locked after 3 failed attempts.
func (s *AuthService) VerifyOTP(ctx context.Context, req OTPVerifyRequest) (*AuthResponse, error) {
	mob, err := valueobject.NewMobile(req.Mobile)
	if err != nil {
		return nil, err
	}

	// Check attempt count before validating
	attempts, err := s.otpStore.GetAttempts(ctx, mob.String())
	if err != nil {
		return nil, shared.ErrInternal
	}
	if attempts >= redis.MaxOTPAttempts() {
		return nil, shared.ErrOTPMaxAttempts
	}

	storedOTP, err := s.otpStore.GetOTP(ctx, mob.String())
	if err != nil {
		return nil, shared.ErrInternal
	}
	if storedOTP == "" {
		return nil, shared.ErrOTPExpired
	}

	if storedOTP != req.Code {
		if _, incrErr := s.otpStore.IncrAttempts(ctx, mob.String()); incrErr != nil {
			log.Error().Err(incrErr).Msg("verifyOTP: incr attempts failed")
		}
		return nil, shared.ErrOTPInvalid
	}

	// OTP valid — clean up
	_ = s.otpStore.DeleteOTP(ctx, mob.String())
	_ = s.otpStore.DeleteAttempts(ctx, mob.String())

	// Find or auto-create client user
	user, err := s.userRepo.FindByMobile(ctx, mob.String())
	if err != nil {
		return nil, shared.ErrInternal
	}
	if user == nil {
		return nil, shared.ErrUserNotFound
	}
	if !user.IsActive {
		return nil, shared.ErrForbidden
	}

	tokens, err := s.jwtService.GenerateTokenPair(user.ID, user.Role)
	if err != nil {
		return nil, shared.ErrInternal
	}

	return &AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TokenType:    "Bearer",
		UserID:       user.ID,
		Role:         user.Role,
	}, nil
}

// RefreshToken validates a refresh token and returns a new access token.
func (s *AuthService) RefreshToken(ctx context.Context, req RefreshRequest) (*AuthResponse, error) {
	claims, err := s.jwtService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, shared.ErrInvalidToken
	}

	// Check blacklist
	revoked, err := s.blacklist.IsRevoked(ctx, claims.ID)
	if err != nil {
		return nil, shared.ErrInternal
	}
	if revoked {
		return nil, shared.ErrTokenRevoked
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, shared.ErrInvalidToken
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, shared.ErrInternal
	}
	if user == nil || !user.IsActive {
		return nil, shared.ErrUnauthorized
	}

	// Rotate: revoke old refresh token
	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl > 0 {
		_ = s.blacklist.Revoke(ctx, claims.ID, ttl)
	}

	tokens, err := s.jwtService.GenerateTokenPair(user.ID, user.Role)
	if err != nil {
		return nil, shared.ErrInternal
	}

	return &AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TokenType:    "Bearer",
		UserID:       user.ID,
		Role:         user.Role,
	}, nil
}

// RevokeToken adds a token JTI to the blacklist with the given TTL.
// Called on logout to enable immediate access token invalidation.
func (s *AuthService) RevokeToken(ctx context.Context, jti string, ttl time.Duration) error {
	if jti == "" || ttl <= 0 {
		return nil
	}
	return s.blacklist.Revoke(ctx, jti, ttl)
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	claims, err := s.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		// Already invalid — treat as logged out
		return nil
	}

	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl > 0 {
		if err := s.blacklist.Revoke(ctx, claims.ID, ttl); err != nil {
			log.Error().Err(err).Msg("logout: revoke token failed")
			return shared.ErrInternal
		}
	}

	return nil
}
