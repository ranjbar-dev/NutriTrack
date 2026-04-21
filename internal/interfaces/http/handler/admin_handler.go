package handler

import (
	"github.com/gin-gonic/gin"
	appAdmin "github.com/ranjbar-dev/nutritrack/internal/application/admin"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	"github.com/ranjbar-dev/nutritrack/internal/interfaces/http/dto"
)

// AdminHandler handles super-admin dashboard endpoints.
type AdminHandler struct {
	service *appAdmin.AdminService
}

// NewAdminHandler constructs an AdminHandler.
func NewAdminHandler(svc *appAdmin.AdminService) *AdminHandler {
	return &AdminHandler{service: svc}
}

// GetStats handles GET /api/v1/admin/stats.
// Returns aggregated platform counts for the super-admin dashboard.
func (h *AdminHandler) GetStats(c *gin.Context) {
	stats, err := h.service.GetStats(c.Request.Context())
	if err != nil {
		dto.Abort(c, shared.ErrInternal)
		return
	}
	dto.OK(c, gin.H{
		"total_nutritionists":    stats.TotalNutritionists,
		"active_nutritionists":   stats.ActiveNutritionists,
		"inactive_nutritionists": stats.InactiveNutritionists,
		"total_clients":          stats.TotalClients,
		"total_foods":            stats.TotalFoods,
		"active_diet_plans":      stats.ActiveDietPlans,
	})
}
