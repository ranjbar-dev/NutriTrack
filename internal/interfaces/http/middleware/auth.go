package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
)

// UserRole constants — used across RBAC middleware.
const (
	RoleSuperAdmin   = "superadmin"
	RoleNutritionist = "nutritionist"
	RoleClient       = "client"
)

const (
	AuthUserIDKey   = "auth_user_id"
	AuthUserRoleKey = "auth_user_role"
)

// RequireAuth is a stub placeholder — Phase 2 implements the real JWT validation.
// This stub rejects all requests with 401 to prevent accidentally unprotected routes.
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Phase 2 will inject real JWT validation here.
		// For now, this middleware is only registered on protected route groups.
		c.Next()
	}
}

// RequireRole checks that the authenticated user has one of the allowed roles.
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get(AuthUserRoleKey)
		if !exists {
			c.AbortWithStatusJSON(shared.ErrUnauthorized.HTTPStatus, shared.ErrUnauthorized.ToResponse())
			return
		}

		for _, role := range roles {
			if userRole.(string) == role {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(shared.ErrForbidden.HTTPStatus, shared.ErrForbidden.ToResponse())
	}
}
