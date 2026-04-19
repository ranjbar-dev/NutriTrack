package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"

	"github.com/ranjbar-dev/nutritrack/backend/internal/model"
	"github.com/ranjbar-dev/nutritrack/backend/internal/model/dto"
	"github.com/ranjbar-dev/nutritrack/backend/internal/repository"
	"github.com/ranjbar-dev/nutritrack/backend/internal/repository/sqlc"
	jwtpkg "github.com/ranjbar-dev/nutritrack/backend/pkg/jwt"
	"github.com/ranjbar-dev/nutritrack/backend/pkg/sms"
)

const (
	otpLength     = 6
	otpExpiry     = 2 * time.Minute
	otpMaxAttempt = 3
	bcryptCost    = 12
)

// ErrInvalidCredentials is returned for any authentication failure.
// Using a single error message for all auth failures prevents user enumeration (T-04-01, T-04-02).
var (
	ErrInvalidCredentials = errors.New("ایمیل یا رمز عبور اشتباه است")
	ErrInvalidOTP         = errors.New("کد تایید نامعتبر است")
	ErrInvalidToken       = errors.New("توکن نامعتبر است")
	ErrUserNotFound       = errors.New("کاربر یافت نشد")
)

// AuthService handles authentication business logic.
type AuthService struct {
	userRepo  repository.UserRepository
	otpRepo   repository.OTPRepository
	tokenRepo repository.TokenRepository
	smsSender sms.Sender
	jwtSecret []byte
	logger    zerolog.Logger
}

// NewAuthService creates a new AuthService with all required dependencies.
func NewAuthService(
	userRepo repository.UserRepository,
	otpRepo repository.OTPRepository,
	tokenRepo repository.TokenRepository,
	smsSender sms.Sender,
	jwtSecret []byte,
	logger zerolog.Logger,
) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		otpRepo:   otpRepo,
		tokenRepo: tokenRepo,
		smsSender: smsSender,
		jwtSecret: jwtSecret,
		logger:    logger,
	}
}

// LoginWithPassword authenticates a user with email and password (AUTH-01, AUTH-03).
// Returns the same error for both "user not found" and "wrong password" to prevent enumeration (T-04-01).
func (s *AuthService) LoginWithPassword(ctx context.Context, email, password string) (*model.TokenPair, *dto.UserResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		// User not found — return same error as wrong password (anti-enumeration)
		s.logger.Debug().Str("email", email).Msg("login failed: user not found")
		return nil, nil, ErrInvalidCredentials
	}

	// Check password hash
	if !user.PasswordHash.Valid {
		s.logger.Warn().Str("email", email).Msg("login failed: user has no password hash")
		return nil, nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.String), []byte(password)); err != nil {
		s.logger.Debug().Str("email", email).Msg("login failed: wrong password")
		return nil, nil, ErrInvalidCredentials
	}

	// Create token pair
	userID := uuid.UUID(user.ID.Bytes)
	tokens, err := s.createAndStoreTokens(ctx, userID, string(user.Role))
	if err != nil {
		return nil, nil, fmt.Errorf("create tokens: %w", err)
	}

	userResp := sqlcUserToResponse(user)
	s.logger.Info().Str("user_id", userID.String()).Str("role", string(user.Role)).Msg("user logged in via password")

	return tokens, userResp, nil
}

// RequestOTP generates and sends an OTP code to the given mobile number (AUTH-05).
// Always returns nil regardless of whether the user exists to prevent enumeration (T-04-02).
func (s *AuthService) RequestOTP(ctx context.Context, mobile string) error {
	mobile = normalizePhone(mobile)

	// Generate a cryptographically secure 6-digit OTP (crypto/rand, not math/rand)
	code, err := generateOTP(otpLength)
	if err != nil {
		return fmt.Errorf("generate OTP: %w", err)
	}

	// Hash OTP with SHA-256 before storage — never store plaintext OTP (T-04-06)
	codeHash := hashSHA256(code)

	// Store OTP record with 2-minute expiry and max 3 attempts (AUTH-06)
	_, err = s.otpRepo.Create(ctx, sqlc.CreateOTPParams{
		Mobile:      mobile,
		CodeHash:    codeHash,
		ExpiresAt:   pgtype.Timestamptz{Time: time.Now().Add(otpExpiry), Valid: true},
		MaxAttempts: otpMaxAttempt,
	})
	if err != nil {
		return fmt.Errorf("store OTP: %w", err)
	}

	// Send OTP via SMS (MockSender in dev, KavenegarSender in prod per D-04)
	if err := s.smsSender.SendOTP(mobile, code); err != nil {
		s.logger.Error().Err(err).Str("mobile", mobile).Msg("failed to send OTP SMS")
		// Don't return error to the caller — still return success for anti-enumeration
	}

	s.logger.Info().Str("mobile", mobile).Msg("OTP requested")
	return nil
}

// VerifyOTP verifies an OTP code for the given mobile number (AUTH-05, AUTH-06).
// Returns the same error for all failure modes to prevent enumeration (T-04-02).
func (s *AuthService) VerifyOTP(ctx context.Context, mobile, code string) (*model.TokenPair, *dto.UserResponse, error) {
	mobile = normalizePhone(mobile)

	// Retrieve active (non-expired, non-verified) OTP for this mobile
	otp, err := s.otpRepo.GetActiveByMobile(ctx, mobile)
	if err != nil {
		s.logger.Debug().Str("mobile", mobile).Msg("OTP verify failed: no active OTP")
		return nil, nil, ErrInvalidOTP
	}

	otpID := uuid.UUID(otp.ID.Bytes)

	// Check max attempts (AUTH-06: max 3 verification attempts)
	if otp.Attempts >= otp.MaxAttempts {
		s.logger.Debug().Str("mobile", mobile).Msg("OTP verify failed: max attempts exceeded")
		return nil, nil, ErrInvalidOTP
	}

	// Increment attempt counter before checking hash (prevents race conditions)
	if err := s.otpRepo.IncrementAttempts(ctx, otpID); err != nil {
		return nil, nil, fmt.Errorf("increment OTP attempts: %w", err)
	}

	// Compare SHA-256 hash of submitted code with stored code_hash (T-04-06)
	if hashSHA256(code) != otp.CodeHash {
		s.logger.Debug().Str("mobile", mobile).Msg("OTP verify failed: wrong code")
		return nil, nil, ErrInvalidOTP
	}

	// Mark OTP as verified
	if err := s.otpRepo.MarkVerified(ctx, otpID); err != nil {
		return nil, nil, fmt.Errorf("mark OTP verified: %w", err)
	}

	// Look up the user by mobile
	user, err := s.userRepo.GetByMobile(ctx, mobile)
	if err != nil {
		s.logger.Debug().Str("mobile", mobile).Msg("OTP verify failed: user not found for mobile")
		return nil, nil, ErrInvalidOTP
	}

	// Create token pair
	userID := uuid.UUID(user.ID.Bytes)
	tokens, err := s.createAndStoreTokens(ctx, userID, string(user.Role))
	if err != nil {
		return nil, nil, fmt.Errorf("create tokens: %w", err)
	}

	userResp := sqlcUserToResponse(user)
	s.logger.Info().Str("user_id", userID.String()).Msg("user logged in via OTP")

	return tokens, userResp, nil
}

// RefreshTokens rotates the refresh token and issues a new token pair (AUTH-08).
// Implements refresh token rotation with theft detection via family_id (T-04-05).
func (s *AuthService) RefreshTokens(ctx context.Context, refreshTokenStr string) (*model.TokenPair, error) {
	tokenHash := hashSHA256(refreshTokenStr)

	// Try to find the token (non-revoked)
	token, err := s.tokenRepo.GetByHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Token not found among non-revoked — check if it was already revoked (replay detection)
			revokedToken, replayErr := s.tokenRepo.GetByHashAny(ctx, tokenHash)
			if replayErr == nil && revokedToken.Revoked {
				// REPLAY DETECTED: Token was already used and revoked.
				// Revoke entire family to protect against token theft (T-04-05).
				familyID := uuid.UUID(revokedToken.FamilyID.Bytes)
				_ = s.tokenRepo.RevokeFamily(ctx, familyID)
				s.logger.Warn().
					Str("family_id", familyID.String()).
					Msg("refresh token replay detected — entire family revoked")
			}
		}
		return nil, ErrInvalidToken
	}

	// Check expiry
	if token.ExpiresAt.Valid && token.ExpiresAt.Time.Before(time.Now()) {
		return nil, ErrInvalidToken
	}

	// Revoke the old token
	tokenID := uuid.UUID(token.ID.Bytes)
	if err := s.tokenRepo.Revoke(ctx, tokenID); err != nil {
		return nil, fmt.Errorf("revoke old token: %w", err)
	}

	// Parse the refresh token to get user info for new pair
	claims, err := jwtpkg.ParseToken(refreshTokenStr, s.jwtSecret)
	if err != nil {
		return nil, ErrInvalidToken
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// Create new token pair with same family_id (rotation within the same family)
	familyID := uuid.UUID(token.FamilyID.Bytes)
	tokens, err := s.createAndStoreTokensWithFamily(ctx, userID, claims.Role, familyID)
	if err != nil {
		return nil, fmt.Errorf("create new tokens: %w", err)
	}

	s.logger.Debug().Str("user_id", claims.UserID).Msg("tokens refreshed")
	return tokens, nil
}

// Logout revokes the refresh token for the given token string.
func (s *AuthService) Logout(ctx context.Context, refreshTokenStr string) error {
	if refreshTokenStr == "" {
		return nil
	}

	tokenHash := hashSHA256(refreshTokenStr)
	token, err := s.tokenRepo.GetByHash(ctx, tokenHash)
	if err != nil {
		// Token not found or already revoked — nothing to do
		return nil
	}

	tokenID := uuid.UUID(token.ID.Bytes)
	return s.tokenRepo.Revoke(ctx, tokenID)
}

// GetUserByID retrieves a user by their UUID (used by GET /api/auth/me).
func (s *AuthService) GetUserByID(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	resp := sqlcUserToResponse(user)
	return resp, nil
}

// createAndStoreTokens creates a new JWT token pair and stores the refresh token hash in DB.
// Generates a new family_id for fresh logins.
func (s *AuthService) createAndStoreTokens(ctx context.Context, userID uuid.UUID, role string) (*model.TokenPair, error) {
	familyID := uuid.New()
	return s.createAndStoreTokensWithFamily(ctx, userID, role, familyID)
}

// createAndStoreTokensWithFamily creates a JWT token pair and stores the refresh token hash
// with the given family_id (used for token rotation within the same family).
func (s *AuthService) createAndStoreTokensWithFamily(ctx context.Context, userID uuid.UUID, role string, familyID uuid.UUID) (*model.TokenPair, error) {
	tokens, err := jwtpkg.CreateTokenPair(userID.String(), role, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("create token pair: %w", err)
	}

	// Hash refresh token with SHA-256 before storage
	refreshHash := hashSHA256(tokens.RefreshToken)

	_, err = s.tokenRepo.Create(ctx, sqlc.CreateRefreshTokenParams{
		UserID:    pgtype.UUID{Bytes: userID, Valid: true},
		TokenHash: refreshHash,
		FamilyID:  pgtype.UUID{Bytes: familyID, Valid: true},
		ExpiresAt: pgtype.Timestamptz{Time: time.Now().Add(30 * 24 * time.Hour), Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("store refresh token: %w", err)
	}

	return tokens, nil
}

// normalizePhone normalizes an Iranian phone number to canonical 10-digit format (9XXXXXXXXX).
// Strips leading +98, 0098, or 0 prefix.
func normalizePhone(phone string) string {
	phone = strings.TrimSpace(phone)

	switch {
	case strings.HasPrefix(phone, "+98"):
		phone = phone[3:]
	case strings.HasPrefix(phone, "0098"):
		phone = phone[4:]
	case strings.HasPrefix(phone, "0"):
		phone = phone[1:]
	}

	return phone
}

// generateOTP generates a cryptographically secure n-digit OTP code using crypto/rand.
func generateOTP(length int) (string, error) {
	max := new(big.Int)
	max.SetString(strings.Repeat("9", length), 10)
	max.Add(max, big.NewInt(1))

	// Generate a random number in [0, 10^length)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", fmt.Errorf("generate random number: %w", err)
	}

	// Pad with leading zeros to ensure exactly `length` digits
	format := fmt.Sprintf("%%0%dd", length)
	return fmt.Sprintf(format, n), nil
}

// hashSHA256 returns the hex-encoded SHA-256 hash of the input string.
func hashSHA256(input string) string {
	h := sha256.Sum256([]byte(input))
	return hex.EncodeToString(h[:])
}

// sqlcUserToResponse converts a sqlc.User to a dto.UserResponse.
func sqlcUserToResponse(user *sqlc.User) *dto.UserResponse {
	resp := &dto.UserResponse{
		ID:       uuid.UUID(user.ID.Bytes).String(),
		Role:     string(user.Role),
		FullName: user.FullName,
	}
	if user.Email.Valid {
		resp.Email = user.Email.String
	}
	if user.Mobile.Valid {
		resp.Mobile = user.Mobile.String
	}
	return resp
}
