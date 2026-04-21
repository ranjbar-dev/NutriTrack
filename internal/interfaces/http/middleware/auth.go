package middleware

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	// AuthTokenJTIKey holds the JWT ID (JTI) of the current access token, used for blacklisting.
	AuthTokenJTIKey = "auth_token_jti"
)

// TokenRevocationChecker is satisfied by any type that can check whether a token JTI has been
// revoked (e.g. infrastructure/redis.TokenBlacklist). Defined here to avoid a direct dependency
// on the infrastructure layer from the interface layer.
type TokenRevocationChecker interface {
	IsRevoked(ctx context.Context, tokenID string) (bool, error)
}

// RequireAuth validates the JWT Bearer token and injects user ID + role into context.
// If revocationChecker is non-nil, it also verifies the token JTI has not been revoked —
// enabling immediate invalidation when a user logs out.
func RequireAuth(jwtService *auth.JWTService, revocationChecker TokenRevocationChecker) gin.HandlerFunc {
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

		// Check token blacklist on every authenticated request.
		if revocationChecker != nil {
			revoked, rErr := revocationChecker.IsRevoked(c.Request.Context(), claims.ID)
			if rErr == nil && revoked {
				dto.Abort(c, shared.ErrTokenRevoked)
				return
			}
		}

		// Parse UserID string into uuid.UUID — required by all downstream handlers.
		userID, parseErr := uuid.Parse(claims.UserID)
		if parseErr != nil {
			dto.Abort(c, shared.ErrInvalidToken)
			return
		}

		c.Set(AuthUserIDKey, userID)
		c.Set(AuthUserRoleKey, claims.Role)
		c.Set(AuthTokenJTIKey, claims.ID)
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
