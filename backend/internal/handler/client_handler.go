package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ranjbar-dev/nutritrack/backend/internal/model/dto"
	"github.com/ranjbar-dev/nutritrack/backend/internal/service"
)

// ClientHandler handles client-related HTTP endpoints.
// NOTE: There is NO public /api/client/register endpoint (AUTH-12, D-05).
// Client registration is exclusively through the nutritionist-protected route.
type ClientHandler struct {
	userService *service.UserService
}

// NewClientHandler creates a new ClientHandler with the given user service.
func NewClientHandler(userService *service.UserService) *ClientHandler {
	return &ClientHandler{userService: userService}
}

// RegisterClient handles POST /api/nutritionist/clients (AUTH-12, CLNT-01).
// Nutritionist-initiated client registration only.
// The nutritionist_id is extracted from the JWT (set by Auth middleware).
func (h *ClientHandler) RegisterClient(c *gin.Context) {
	var req dto.RegisterClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "اطلاعات ورودی نامعتبر است"})
		return
	}

	// Extract nutritionist ID from JWT context (set by Auth middleware)
	nutritionistIDStr := c.GetString("user_id")
	if nutritionistIDStr == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "احراز هویت الزامی است"})
		return
	}

	nutritionistID, err := uuid.Parse(nutritionistIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "توکن نامعتبر است"})
		return
	}

	userResp, err := h.userService.RegisterClient(c.Request.Context(), nutritionistID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "خطا در ثبت‌نام مراجع"})
		return
	}

	c.JSON(http.StatusCreated, dto.AuthResponse{User: *userResp})
}
