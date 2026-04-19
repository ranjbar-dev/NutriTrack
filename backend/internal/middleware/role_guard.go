package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RoleGuard checks that the authenticated user's role (set by Auth middleware)
// matches one of the allowed roles. Returns 403 if unauthorized.
// MUST be applied after Auth middleware in the middleware chain.
func RoleGuard(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		for _, allowed := range allowedRoles {
			if role == allowed {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "دسترسی غیرمجاز"})
	}
}
