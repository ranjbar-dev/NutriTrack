package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ranjbar-dev/nutritrack/backend/pkg/jwt"
)

// Auth extracts the JWT from the "access_token" httpOnly cookie (D-01)
// and sets "user_id" and "role" in the Gin context for downstream handlers.
func Auth(jwtSecret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("access_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "احراز هویت الزامی است"})
			return
		}

		claims, err := jwt.ParseToken(tokenStr, jwtSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "توکن نامعتبر است"})
			return
		}

		// Set user context for downstream handlers
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
