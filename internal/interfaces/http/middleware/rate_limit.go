package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	"github.com/ranjbar-dev/nutritrack/internal/interfaces/http/dto"
)

// ErrRateLimitExceeded is the Persian error returned when a client exceeds the request rate limit.
var ErrRateLimitExceeded = &shared.AppError{
	Code:       "RATE_LIMIT_EXCEEDED",
	Message:    "تعداد درخواست‌های شما بیش از حد مجاز است",
	HTTPStatus: 429,
}

const (
	ipRateLimitWindow    = time.Minute
	ipRateLimitKeyPrefix = "rate:ip:"
)

// RateLimitByIP returns a per-IP sliding-window rate limiting middleware backed by Redis.
// max specifies the maximum number of requests allowed per ipRateLimitWindow (1 minute).
// On Redis failure the middleware fails open (request is allowed) to avoid false positives.
func RateLimitByIP(rdb *redis.Client, max int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := ipRateLimitKeyPrefix + ip

		ctx := c.Request.Context()
		count, err := rdb.Incr(ctx, key).Result()
		if err != nil {
			// Fail open: do not block requests when Redis is unavailable.
			c.Next()
			return
		}

		// Set the expiry only on the first increment to start the window.
		if count == 1 {
			rdb.Expire(ctx, key, ipRateLimitWindow)
		}

		if count > max {
			dto.Abort(c, ErrRateLimitExceeded)
			return
		}

		c.Next()
	}
}
