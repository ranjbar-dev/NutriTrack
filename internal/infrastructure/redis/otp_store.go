package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	otpTTL          = 2 * time.Minute
	attemptsTTL     = 10 * time.Minute
	rateLimitTTL    = 10 * time.Minute
	maxOTPAttempts  = 3
	maxOTPRateLimit = 3 // max OTP sends per rateLimitTTL window
)

// OTPStore implements OTPRepository using Redis.
type OTPStore struct {
	client *redis.Client
}

func NewOTPStore(client *redis.Client) *OTPStore {
	return &OTPStore{client: client}
}

func (s *OTPStore) otpKey(mobile string) string {
	return fmt.Sprintf("otp:%s", mobile)
}

func (s *OTPStore) attemptsKey(mobile string) string {
	return fmt.Sprintf("otp_attempts:%s", mobile)
}

func (s *OTPStore) rateLimitKey(mobile string) string {
	return fmt.Sprintf("otp_rate:%s", mobile)
}

func (s *OTPStore) StoreOTP(ctx context.Context, mobile, otp string) error {
	return s.client.Set(ctx, s.otpKey(mobile), otp, otpTTL).Err()
}

func (s *OTPStore) GetOTP(ctx context.Context, mobile string) (string, error) {
	val, err := s.client.Get(ctx, s.otpKey(mobile)).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (s *OTPStore) DeleteOTP(ctx context.Context, mobile string) error {
	return s.client.Del(ctx, s.otpKey(mobile)).Err()
}

// IncrAttempts uses atomic INCR and sets TTL only on first increment.
// This prevents the TOCTOU race condition that GET+SET would have.
func (s *OTPStore) IncrAttempts(ctx context.Context, mobile string) (int64, error) {
	key := s.attemptsKey(mobile)
	count, err := s.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	// Set TTL only on first increment (count == 1)
	if count == 1 {
		s.client.Expire(ctx, key, attemptsTTL)
	}
	return count, nil
}

func (s *OTPStore) GetAttempts(ctx context.Context, mobile string) (int64, error) {
	val, err := s.client.Get(ctx, s.attemptsKey(mobile)).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return val, err
}

func (s *OTPStore) DeleteAttempts(ctx context.Context, mobile string) error {
	return s.client.Del(ctx, s.attemptsKey(mobile)).Err()
}

// IncrRateLimit uses atomic INCR — same pattern as IncrAttempts.
// Returns new count so caller can check against maxOTPRateLimit.
func (s *OTPStore) IncrRateLimit(ctx context.Context, mobile string) (int64, error) {
	key := s.rateLimitKey(mobile)
	count, err := s.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if count == 1 {
		s.client.Expire(ctx, key, rateLimitTTL)
	}
	return count, nil
}

func (s *OTPStore) GetRateLimit(ctx context.Context, mobile string) (int64, error) {
	val, err := s.client.Get(ctx, s.rateLimitKey(mobile)).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return val, err
}

// MaxOTPAttempts returns the maximum allowed failed attempts.
func MaxOTPAttempts() int64 { return maxOTPAttempts }

// MaxOTPRateLimit returns the maximum OTP sends per rate limit window.
func MaxOTPRateLimit() int64 { return maxOTPRateLimit }
