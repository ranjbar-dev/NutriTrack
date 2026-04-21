package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
)

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		requestID, _ := c.Get(RequestIDKey)
		log.Error().
			Str("request_id", requestID.(string)).
			Interface("panic", recovered).
			Msg("panic recovered")

		c.JSON(http.StatusInternalServerError, shared.ErrInternal.ToResponse())
		c.Abort()
	})
}
