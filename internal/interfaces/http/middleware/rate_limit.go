package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	"github.com/ranjbar-dev/nutritrack/internal/interfaces/http/dto"
)

// RateLimiter abstracts the rate-limiting store (e.g. Redis) from the middleware.
// Modelled after the TokenRevocationChecker pattern: the concrete implementation
// lives in internal/infrastructure/redis and is injected by the router.
type RateLimiter interface {
	Allow(ctx context.Context, key string, max int64, window time.Duration) (bool, error)
}

// ErrRateLimitExceeded is the Persian error returned when a client exceeds the request rate limit.
var ErrRateLimitExceeded = &shared.AppError{
	Code:    "RATE_LIMIT_EXCEEDED",
	Message: "تعداد درخواست‌های شما بیش از حد مجاز است",
}

const (
	ipRateLimitWindow    = time.Minute
	ipRateLimitKeyPrefix = "rate:ip:"
)

// RateLimitByIP returns a per-IP sliding-window rate limiting middleware.
// max specifies the maximum number of requests allowed per ipRateLimitWindow (1 minute).
// On limiter failure the middleware fails open (request is allowed) to avoid false positives.
func RateLimitByIP(limiter RateLimiter, max int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := ipRateLimitKeyPrefix + ip

		ctx := c.Request.Context()
		allowed, err := limiter.Allow(ctx, key, max, ipRateLimitWindow)
		if err != nil {
			// Fail open: do not block requests when the limiter is unavailable.
			c.Next()
			return
		}

		if !allowed {
			dto.Abort(c, ErrRateLimitExceeded)
			return
		}

		c.Next()
	}
}
