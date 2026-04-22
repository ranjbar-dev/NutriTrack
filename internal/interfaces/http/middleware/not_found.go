package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	"github.com/ranjbar-dev/nutritrack/internal/interfaces/http/dto"
)

func NotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		dto.Error(c, shared.ErrNotFound)
	}
}
