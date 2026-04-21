package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORS sets CORS headers.
// allowedOrigin is the value for Access-Control-Allow-Origin.
// Pass "*" for local development; in production, pass the exact frontend origin
// (e.g. "https://app.nutritrack.ir") via the CORS_ALLOWED_ORIGINS env variable.
func CORS(allowedOrigin string) gin.HandlerFunc {
	if allowedOrigin == "" {
		allowedOrigin = "*"
	}
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", allowedOrigin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, X-Request-ID")
		c.Header("Access-Control-Expose-Headers", "X-Request-ID")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
