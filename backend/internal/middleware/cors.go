package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS creates a CORS middleware that restricts access to the specified frontend URL (D-24, SEC-06).
// Access-Control-Allow-Credentials is true to allow httpOnly cookies (T-03-03).
// Origin MUST be an explicit URL (not "*") when credentials are allowed.
func CORS(frontendURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" {
			c.Header("Vary", "Origin")
		}

		if origin == frontendURL {
			c.Header("Access-Control-Allow-Origin", frontendURL)
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Max-Age", "86400")
		} else if origin != "" && c.Request.Method == http.MethodOptions {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "مبدا درخواست مجاز نیست"})
			return
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
