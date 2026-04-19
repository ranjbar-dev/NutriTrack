package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// Logger creates a structured JSON logging middleware using zerolog (D-23, INFRA-04).
// Logs: request_id, method, path, status, duration_ms, client_ip.
// Timestamp is added automatically by zerolog.
// Output: structured JSON to stdout for Loki collection (T-03-06).
func Logger(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)

		logger.Info().
			Str("request_id", c.GetString("request_id")).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", c.Writer.Status()).
			Dur("duration_ms", duration).
			Str("client_ip", c.ClientIP()).
			Msg("request")
	}
}
