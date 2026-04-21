package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// TokenBlacklist stores revoked JWT refresh tokens in Redis.
type TokenBlacklist struct {
	client *redis.Client
}

func NewTokenBlacklist(client *redis.Client) *TokenBlacklist {
	return &TokenBlacklist{client: client}
}

func (b *TokenBlacklist) blacklistKey(tokenID string) string {
	return fmt.Sprintf("blacklist:%s", tokenID)
}

// Revoke adds a token to the blacklist with the given TTL.
func (b *TokenBlacklist) Revoke(ctx context.Context, tokenID string, ttl time.Duration) error {
	return b.client.Set(ctx, b.blacklistKey(tokenID), "1", ttl).Err()
}

// IsRevoked returns true if the token has been revoked.
func (b *TokenBlacklist) IsRevoked(ctx context.Context, tokenID string) (bool, error) {
	exists, err := b.client.Exists(ctx, b.blacklistKey(tokenID)).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}
