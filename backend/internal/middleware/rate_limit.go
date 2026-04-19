package middleware

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements an in-memory sliding window rate limiter.
// Safe for concurrent use.
type RateLimiter struct {
	mu      sync.Mutex
	entries map[string][]time.Time
	max     int
	window  time.Duration
}

// NewRateLimiter creates a new RateLimiter with the given max requests per window.
// Starts a background goroutine to clean up expired entries every minute.
func NewRateLimiter(max int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		entries: make(map[string][]time.Time),
		max:     max,
		window:  window,
	}

	// Background cleanup goroutine
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			rl.cleanup()
		}
	}()

	return rl
}

// Allow checks if a request from the given key is allowed under the rate limit.
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	// Filter out expired entries for this key
	timestamps := rl.entries[key]
	valid := timestamps[:0]
	for _, t := range timestamps {
		if t.After(windowStart) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= rl.max {
		rl.entries[key] = valid
		return false
	}

	rl.entries[key] = append(valid, now)
	return true
}

// cleanup removes all expired entries from all keys.
func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	for key, timestamps := range rl.entries {
		valid := timestamps[:0]
		for _, t := range timestamps {
			if t.After(windowStart) {
				valid = append(valid, t)
			}
		}
		if len(valid) == 0 {
			delete(rl.entries, key)
		} else {
			rl.entries[key] = valid
		}
	}
}

// RateLimit creates a Gin middleware that rate-limits requests based on phone number.
// It reads the "mobile" field from the JSON request body. If not available, falls back to client IP.
// Phone numbers are normalized to prevent bypass via format variation (T-03-02, Pitfall 3).
func RateLimit(limiter *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := extractRateLimitKey(c)

		if !limiter.Allow(key) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "تعداد درخواست‌ها بیش از حد مجاز است. لطفاً چند دقیقه صبر کنید.",
			})
			return
		}

		c.Next()
	}
}

// extractRateLimitKey extracts a rate-limit key from the request.
// Reads mobile from JSON body (peeking without consuming), normalizes it.
// Falls back to client IP if mobile is not available.
func extractRateLimitKey(c *gin.Context) string {
	// Try to read mobile from request body without consuming it
	if c.Request.Body != nil && c.ContentType() == "application/json" {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err == nil && len(bodyBytes) > 0 {
			// Restore body for downstream handlers
			c.Request.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))

			var body struct {
				Mobile string `json:"mobile"`
			}
			if err := json.Unmarshal(bodyBytes, &body); err == nil && body.Mobile != "" {
				normalized := normalizePhone(body.Mobile)
				if normalized != "" {
					return "phone:" + normalized
				}
			}
		}
	}

	// Fallback to client IP
	return "ip:" + c.ClientIP()
}

// normalizePhone normalizes an Iranian phone number to canonical 10-digit format (9XXXXXXXXX).
// Strips leading +98, 0098, or 0 prefix to prevent rate limit bypass via format variation.
func normalizePhone(phone string) string {
	phone = strings.TrimSpace(phone)

	// Strip common Iranian prefixes
	switch {
	case strings.HasPrefix(phone, "+98"):
		phone = phone[3:]
	case strings.HasPrefix(phone, "0098"):
		phone = phone[4:]
	case strings.HasPrefix(phone, "0"):
		phone = phone[1:]
	}

	// Validate length: should be 10 digits starting with 9
	if len(phone) == 10 && phone[0] == '9' {
		return phone
	}

	return phone // Return as-is if format is unexpected
}
