package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ranjbar-dev/nutritrack/configs"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
)

// TokenType distinguishes access tokens from refresh tokens.
type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

// Claims is the JWT payload for NutriTrack tokens.
type Claims struct {
	UserID string    `json:"uid"`
	Role   string    `json:"role"`
	Type   TokenType `json:"type"`
	jwt.RegisteredClaims
}

// TokenPair holds an access/refresh token pair together with their expiry times.
type TokenPair struct {
	AccessToken   string
	RefreshToken  string
	AccessExpiry  time.Time
	RefreshExpiry time.Time
}

// JWTService generates and validates JWT access/refresh token pairs.
type JWTService struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

// NewJWTService creates a JWTService wired to the application JWT configuration.
func NewJWTService(cfg *configs.JWTConfig) *JWTService {
	return &JWTService{
		accessSecret:  []byte(cfg.AccessSecret),
		refreshSecret: []byte(cfg.RefreshSecret),
		accessTTL:     time.Duration(cfg.AccessTTLMin) * time.Minute,
		refreshTTL:    time.Duration(cfg.RefreshTTLDay) * 24 * time.Hour,
	}
}

// GenerateTokenPair creates a new access + refresh token pair for the given user.
func (s *JWTService) GenerateTokenPair(userID uuid.UUID, role string) (*TokenPair, error) {
	now := time.Now()

	accessExpiry := now.Add(s.accessTTL)
	accessToken, err := s.generateToken(userID, role, AccessToken, accessExpiry, s.accessSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshExpiry := now.Add(s.refreshTTL)
	refreshToken, err := s.generateToken(userID, role, RefreshToken, refreshExpiry, s.refreshSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		AccessExpiry:  accessExpiry,
		RefreshExpiry: refreshExpiry,
	}, nil
}

// ValidateAccessToken parses and validates an access token string.
func (s *JWTService) ValidateAccessToken(tokenStr string) (*Claims, error) {
	return s.validateToken(tokenStr, s.accessSecret, AccessToken)
}

// ValidateRefreshToken parses and validates a refresh token string.
func (s *JWTService) ValidateRefreshToken(tokenStr string) (*Claims, error) {
	return s.validateToken(tokenStr, s.refreshSecret, RefreshToken)
}

func (s *JWTService) generateToken(userID uuid.UUID, role string, tokenType TokenType, expiry time.Time, secret []byte) (string, error) {
	claims := Claims{
		UserID: userID.String(),
		Role:   role,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expiry),
			Issuer:    "nutritrack",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func (s *JWTService) validateToken(tokenStr string, secret []byte, expectedType TokenType) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, shared.ErrInvalidToken
		}
		return nil, shared.ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, shared.ErrInvalidToken
	}

	if claims.Type != expectedType {
		return nil, shared.ErrInvalidToken
	}

	return claims, nil
}
