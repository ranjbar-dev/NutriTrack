package repository

import "context"

// OTPRepository defines the contract for OTP storage and rate limiting.
// Implementation lives in internal/infrastructure/redis/otp_store.go.
type OTPRepository interface {
	// StoreOTP saves an OTP for the given mobile (2-minute TTL).
	StoreOTP(ctx context.Context, mobile, otp string) error

	// GetOTP retrieves the stored OTP for a mobile. Returns "" if not found/expired.
	GetOTP(ctx context.Context, mobile string) (string, error)

	// DeleteOTP removes the OTP (called after successful verification).
	DeleteOTP(ctx context.Context, mobile string) error

	// IncrAttempts atomically increments the failed attempt counter.
	// Returns the new count. TTL is set on first increment (10 minutes).
	IncrAttempts(ctx context.Context, mobile string) (int64, error)

	// GetAttempts returns the current failed attempt count.
	GetAttempts(ctx context.Context, mobile string) (int64, error)

	// DeleteAttempts removes the attempt counter (called after successful verification).
	DeleteAttempts(ctx context.Context, mobile string) error

	// IncrRateLimit atomically increments the OTP send rate limit counter.
	// Returns the new count. TTL is set on first increment (10 minutes).
	IncrRateLimit(ctx context.Context, mobile string) (int64, error)

	// GetRateLimit returns the current OTP send count for rate limiting.
	GetRateLimit(ctx context.Context, mobile string) (int64, error)
}
