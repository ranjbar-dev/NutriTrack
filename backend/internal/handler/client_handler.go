package handler

import (
	"errors"
	"net/http"
	"strconv"

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

// NutriListClients handles GET /api/nutritionist/clients
func (h *ClientHandler) NutriListClients(c *gin.Context) {
	nutritionistID, ok := parseNutritionistID(c)
	if !ok {
		return
	}

	query := c.DefaultQuery("q", "")
	sortBy := c.DefaultQuery("sort", "created_at")

	var active *bool
	if activeStr := c.Query("active"); activeStr != "" {
		b, err := strconv.ParseBool(activeStr)
		if err == nil {
			active = &b
		}
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	resp, err := h.userService.ListClients(c.Request.Context(), nutritionistID, query, sortBy, active, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "خطا در دریافت لیست مراجعین"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// NutriGetClientProfile handles GET /api/nutritionist/clients/:clientId
func (h *ClientHandler) NutriGetClientProfile(c *gin.Context) {
	nutritionistID, ok := parseNutritionistID(c)
	if !ok {
		return
	}
	clientID, ok := parseClientParam(c)
	if !ok {
		return
	}

	profile, err := h.userService.GetClientProfile(c.Request.Context(), clientID, nutritionistID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "مراجع یافت نشد"})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "خطا در دریافت پروفایل"})
		return
	}
	c.JSON(http.StatusOK, profile)
}

// NutriActivateClient handles PATCH /api/nutritionist/clients/:clientId/activate
func (h *ClientHandler) NutriActivateClient(c *gin.Context) {
	nutritionistID, ok := parseNutritionistID(c)
	if !ok {
		return
	}
	clientID, ok := parseClientParam(c)
	if !ok {
		return
	}

	if err := h.userService.ActivateClient(c.Request.Context(), clientID, nutritionistID); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "مراجع یافت نشد"})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "خطا در فعال‌سازی حساب"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "حساب کاربری فعال شد"})
}

// NutriDeactivateClient handles PATCH /api/nutritionist/clients/:clientId/deactivate
func (h *ClientHandler) NutriDeactivateClient(c *gin.Context) {
	nutritionistID, ok := parseNutritionistID(c)
	if !ok {
		return
	}
	clientID, ok := parseClientParam(c)
	if !ok {
		return
	}

	if err := h.userService.DeactivateClient(c.Request.Context(), clientID, nutritionistID); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "مراجع یافت نشد"})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "خطا در غیرفعال‌سازی حساب"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "حساب کاربری غیرفعال شد"})
}

// NutriUpdateClientProfile handles PATCH /api/nutritionist/clients/:clientId/profile
func (h *ClientHandler) NutriUpdateClientProfile(c *gin.Context) {
	nutritionistID, ok := parseNutritionistID(c)
	if !ok {
		return
	}
	clientID, ok := parseClientParam(c)
	if !ok {
		return
	}

	var req dto.UpdateClientProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "اطلاعات ورودی نامعتبر است"})
		return
	}

	profile, err := h.userService.UpdateClientProfile(c.Request.Context(), clientID, nutritionistID, req)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "مراجع یافت نشد"})
			return
		}
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, profile)
}

// ─── helpers ──────────────────────────────────────────────────────────────────

func parseNutritionistID(c *gin.Context) (uuid.UUID, bool) {
	idStr := c.GetString("user_id")
	if idStr == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "احراز هویت الزامی است"})
		return uuid.UUID{}, false
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "توکن نامعتبر است"})
		return uuid.UUID{}, false
	}
	return id, true
}

func parseClientParam(c *gin.Context) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param("clientId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه مراجع نامعتبر است"})
		return uuid.UUID{}, false
	}
	return id, true
}

