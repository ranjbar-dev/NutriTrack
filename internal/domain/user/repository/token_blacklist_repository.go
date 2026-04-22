package repository

import (
	"context"
	"time"
)

// TokenBlacklistRepository defines the contract for revoking and checking JWT tokens.
// The Redis adapter implements this interface in internal/infrastructure/redis/.
type TokenBlacklistRepository interface {
	// Revoke adds a token to the blacklist with the given TTL.
	Revoke(ctx context.Context, tokenID string, ttl time.Duration) error

	// IsRevoked returns true if the token has been revoked.
	IsRevoked(ctx context.Context, tokenID string) (bool, error)
}
