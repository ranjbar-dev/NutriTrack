package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ranjbar-dev/nutritrack/backend/internal/model/dto"
)

// HealthCheck handles GET /api/health and returns the server status.
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, dto.HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}
