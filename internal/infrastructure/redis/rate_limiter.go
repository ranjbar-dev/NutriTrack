package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisRateLimiter implements the middleware.RateLimiter interface using Redis.
// It uses an increment-and-expire sliding window strategy.
type RedisRateLimiter struct {
	client *redis.Client
}

// NewRedisRateLimiter constructs a RedisRateLimiter backed by the given Redis client.
func NewRedisRateLimiter(client *redis.Client) *RedisRateLimiter {
	return &RedisRateLimiter{client: client}
}

// Allow increments the request counter for key and returns true if the count is within max.
// The first increment starts a sliding window of the given duration.
// On Redis error it returns (false, err) — callers should fail open.
func (r *RedisRateLimiter) Allow(ctx context.Context, key string, max int64, window time.Duration) (bool, error) {
	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		// Set TTL only on the first request to define the window boundary.
		r.client.Expire(ctx, key, window)
	}

	return count <= max, nil
}
