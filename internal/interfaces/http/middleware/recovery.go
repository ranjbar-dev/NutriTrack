package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	"github.com/ranjbar-dev/nutritrack/internal/interfaces/http/dto"
)

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		requestID, _ := c.Get(RequestIDKey)
		log.Error().
			Str("request_id", requestID.(string)).
			Interface("panic", recovered).
			Msg("panic recovered")

		dto.Error(c, shared.ErrInternal)
		c.Abort()
	})
}
