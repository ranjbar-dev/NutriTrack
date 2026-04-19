package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/ranjbar-dev/nutritrack/backend/internal/model"
)

const (
	issuer             = "nutritrack"
	accessTokenExpiry  = 15 * time.Minute
	refreshTokenExpiry = 30 * 24 * time.Hour
)

// Claims represents the JWT claims payload with user identification.
type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// CreateAccessToken creates a signed JWT access token with 15-minute expiry.
func CreateAccessToken(userID, role string, secret []byte) (string, error) {
	return createToken(userID, role, secret, accessTokenExpiry)
}

// CreateRefreshToken creates a signed JWT refresh token with 30-day expiry.
func CreateRefreshToken(userID, role string, secret []byte) (string, error) {
	return createToken(userID, role, secret, refreshTokenExpiry)
}

// CreateTokenPair creates both access and refresh tokens.
func CreateTokenPair(userID, role string, secret []byte) (*model.TokenPair, error) {
	accessToken, err := CreateAccessToken(userID, role, secret)
	if err != nil {
		return nil, fmt.Errorf("create access token: %w", err)
	}

	refreshToken, err := CreateRefreshToken(userID, role, secret)
	if err != nil {
		return nil, fmt.Errorf("create refresh token: %w", err)
	}

	return &model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// ParseToken parses and validates a JWT token string.
// It validates the signing method is HMAC (prevents algorithm confusion attacks)
// and returns the claims if the token is valid.
func ParseToken(tokenString string, secret []byte) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method is HMAC to prevent algorithm confusion attack (T-03-01)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

// createToken is the internal helper that creates a signed JWT with the given expiry.
func createToken(userID, role string, secret []byte, expiry time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
