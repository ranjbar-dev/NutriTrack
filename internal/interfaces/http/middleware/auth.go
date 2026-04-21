package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ranjbar-dev/nutritrack/internal/application/auth"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	"github.com/ranjbar-dev/nutritrack/internal/interfaces/http/dto"
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

// RequireAuth validates the JWT Bearer token and injects user ID + role into context.
func RequireAuth(jwtService *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			dto.Abort(c, shared.ErrUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := jwtService.ValidateAccessToken(tokenStr)
		if err != nil {
			dto.Abort(c, shared.ErrInvalidToken)
			return
		}

		c.Set(AuthUserIDKey, claims.UserID)
		c.Set(AuthUserRoleKey, claims.Role)
		c.Next()
	}
}

// RequireRole checks that the authenticated user has one of the allowed roles.
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get(AuthUserRoleKey)
		if !exists {
			dto.Abort(c, shared.ErrUnauthorized)
			return
		}

		for _, role := range roles {
			if userRole.(string) == role {
				c.Next()
				return
			}
		}

		dto.Abort(c, shared.ErrForbidden)
	}
}
