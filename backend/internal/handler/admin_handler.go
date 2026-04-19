package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ranjbar-dev/nutritrack/backend/internal/model/dto"
	"github.com/ranjbar-dev/nutritrack/backend/internal/service"
)

// AdminHandler handles admin-only HTTP endpoints.
type AdminHandler struct {
	userService *service.UserService
}

// NewAdminHandler creates a new AdminHandler with the given user service.
func NewAdminHandler(userService *service.UserService) *AdminHandler {
	return &AdminHandler{userService: userService}
}

// CreateNutritionist handles POST /api/admin/nutritionists (AUTH-04).
// Only accessible by super_admin role (enforced by RoleGuard middleware).
func (h *AdminHandler) CreateNutritionist(c *gin.Context) {
	var req dto.CreateNutritionistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "اطلاعات ورودی نامعتبر است"})
		return
	}

	userResp, err := h.userService.CreateNutritionist(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "خطا در ایجاد حساب کارشناس تغذیه"})
		return
	}

	c.JSON(http.StatusCreated, dto.AuthResponse{User: *userResp})
}
