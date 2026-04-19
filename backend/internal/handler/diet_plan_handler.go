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

// DietPlanHandler handles all diet plan HTTP endpoints.
type DietPlanHandler struct {
	planService *service.DietPlanService
}

// NewDietPlanHandler creates a new DietPlanHandler.
func NewDietPlanHandler(planService *service.DietPlanService) *DietPlanHandler {
	return &DietPlanHandler{planService: planService}
}

// ─── Plan-level handlers ──────────────────────────────────────────────────────

// CreatePlan handles POST /api/diet-plans.
func (h *DietPlanHandler) CreatePlan(c *gin.Context) {
	var req dto.CreateDietPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "اطلاعات ورودی نامعتبر است"})
		return
	}

	nutritionistID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "توکن نامعتبر است"})
		return
	}

	resp, err := h.planService.CreatePlan(c.Request.Context(), nutritionistID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrPlanInvalidDateRange):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		default:
			h.handlePlanError(c, err)
		}
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetPlanAggregate handles GET /api/diet-plans/:id.
func (h *DietPlanHandler) GetPlanAggregate(c *gin.Context) {
	planID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه نامعتبر است"})
		return
	}

	nutritionistID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "توکن نامعتبر است"})
		return
	}

	resp, err := h.planService.GetPlanAggregate(c.Request.Context(), planID, nutritionistID)
	if err != nil {
		h.handlePlanError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdatePlanHeader handles PATCH /api/diet-plans/:id.
func (h *DietPlanHandler) UpdatePlanHeader(c *gin.Context) {
	planID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه نامعتبر است"})
		return
	}

	var req dto.UpdateDietPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "اطلاعات ورودی نامعتبر است"})
		return
	}

	nutritionistID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "توکن نامعتبر است"})
		return
	}

	resp, err := h.planService.UpdatePlanHeader(c.Request.Context(), planID, nutritionistID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrPlanInvalidDateRange):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		default:
			h.handlePlanError(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ActivatePlan handles PATCH /api/diet-plans/:id/activate.
func (h *DietPlanHandler) ActivatePlan(c *gin.Context) {
	planID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه نامعتبر است"})
		return
	}

	nutritionistID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "توکن نامعتبر است"})
		return
	}

	if err := h.planService.ActivatePlan(c.Request.Context(), planID, nutritionistID); err != nil {
		h.handlePlanError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "برنامه با موفقیت فعال شد"})
}

// DeletePlan handles DELETE /api/diet-plans/:id.
func (h *DietPlanHandler) DeletePlan(c *gin.Context) {
	planID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه نامعتبر است"})
		return
	}

	nutritionistID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "توکن نامعتبر است"})
		return
	}

	if err := h.planService.DeletePlan(c.Request.Context(), planID, nutritionistID); err != nil {
		h.handlePlanError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ListClientPlans handles GET /api/clients/:clientId/plans.
func (h *DietPlanHandler) ListClientPlans(c *gin.Context) {
	clientID, err := uuid.Parse(c.Param("clientId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه نامعتبر است"})
		return
	}

	nutritionistID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "توکن نامعتبر است"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	resp, err := h.planService.ListClientPlans(c.Request.Context(), clientID, nutritionistID, page, limit)
	if err != nil {
		h.handlePlanError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ─── Day handlers ─────────────────────────────────────────────────────────────

// AddDay handles POST /api/diet-plans/:id/days.
func (h *DietPlanHandler) AddDay(c *gin.Context) {
	planID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه نامعتبر است"})
		return
	}

	var req dto.CreateDayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "اطلاعات ورودی نامعتبر است"})
		return
	}

	nutritionistID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "توکن نامعتبر است"})
		return
	}

	resp, err := h.planService.AddDay(c.Request.Context(), planID, nutritionistID, req)
	if err != nil {
		h.handlePlanError(c, err)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// UpdateDay handles PUT /api/diet-plans/:id/days/:dayId.
func (h *DietPlanHandler) UpdateDay(c *gin.Context) {
	planID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه نامعتبر است"})
		return
	}
	dayID, err := uuid.Parse(c.Param("dayId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه روز نامعتبر است"})
		return
	}

	var req dto.UpdateDayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "اطلاعات ورودی نامعتبر است"})
		return
	}

	nutritionistID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "توکن نامعتبر است"})
		return
	}

	resp, err := h.planService.UpdateDay(c.Request.Context(), planID, dayID, nutritionistID, req)
	if err != nil {
		h.handlePlanError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteDay handles DELETE /api/diet-plans/:id/days/:dayId.
func (h *DietPlanHandler) DeleteDay(c *gin.Context) {
	planID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه نامعتبر است"})
		return
	}
	dayID, err := uuid.Parse(c.Param("dayId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه روز نامعتبر است"})
		return
	}

	nutritionistID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "توکن نامعتبر است"})
		return
	}

	if err := h.planService.DeleteDay(c.Request.Context(), planID, dayID, nutritionistID); err != nil {
		h.handlePlanError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ─── Meal handlers ────────────────────────────────────────────────────────────

// AddMeal handles POST /api/diet-plans/:id/days/:dayId/meals.
func (h *DietPlanHandler) AddMeal(c *gin.Context) {
	planID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه نامعتبر است"})
		return
	}
	dayID, err := uuid.Parse(c.Param("dayId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه روز نامعتبر است"})
		return
	}

	var req dto.CreateMealRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "اطلاعات ورودی نامعتبر است"})
		return
	}

	nutritionistID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "توکن نامعتبر است"})
		return
	}

	resp, err := h.planService.AddMeal(c.Request.Context(), planID, dayID, nutritionistID, req)
	if err != nil {
		h.handlePlanError(c, err)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// UpdateMeal handles PUT /api/diet-plans/:id/days/:dayId/meals/:mealId.
func (h *DietPlanHandler) UpdateMeal(c *gin.Context) {
	planID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه نامعتبر است"})
		return
	}
	dayID, err := uuid.Parse(c.Param("dayId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه روز نامعتبر است"})
		return
	}
	mealID, err := uuid.Parse(c.Param("mealId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه وعده نامعتبر است"})
		return
	}

	var req dto.UpdateMealRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "اطلاعات ورودی نامعتبر است"})
		return
	}

	nutritionistID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "توکن نامعتبر است"})
		return
	}

	resp, err := h.planService.UpdateMeal(c.Request.Context(), planID, dayID, mealID, nutritionistID, req)
	if err != nil {
		h.handlePlanError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteMeal handles DELETE /api/diet-plans/:id/days/:dayId/meals/:mealId.
func (h *DietPlanHandler) DeleteMeal(c *gin.Context) {
	planID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه نامعتبر است"})
		return
	}
	dayID, err := uuid.Parse(c.Param("dayId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه روز نامعتبر است"})
		return
	}
	mealID, err := uuid.Parse(c.Param("mealId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه وعده نامعتبر است"})
		return
	}

	nutritionistID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "توکن نامعتبر است"})
		return
	}

	if err := h.planService.DeleteMeal(c.Request.Context(), planID, dayID, mealID, nutritionistID); err != nil {
		h.handlePlanError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ReorderMeal handles PATCH /api/diet-plans/:id/days/:dayId/meals/:mealId/order.
func (h *DietPlanHandler) ReorderMeal(c *gin.Context) {
	planID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه نامعتبر است"})
		return
	}
	dayID, err := uuid.Parse(c.Param("dayId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه روز نامعتبر است"})
		return
	}
	mealID, err := uuid.Parse(c.Param("mealId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه وعده نامعتبر است"})
		return
	}

	var body struct {
		DisplayOrder int `json:"display_order"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "اطلاعات ورودی نامعتبر است"})
		return
	}

	nutritionistID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "توکن نامعتبر است"})
		return
	}

	if err := h.planService.ReorderMeal(c.Request.Context(), planID, dayID, mealID, nutritionistID, int32(body.DisplayOrder)); err != nil {
		h.handlePlanError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ترتیب وعده با موفقیت به‌روز شد"})
}

// ─── Option handlers ──────────────────────────────────────────────────────────

// AddOption handles POST /api/diet-plans/:id/days/:dayId/meals/:mealId/options.
func (h *DietPlanHandler) AddOption(c *gin.Context) {
	planID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه نامعتبر است"})
		return
	}
	dayID, err := uuid.Parse(c.Param("dayId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه روز نامعتبر است"})
		return
	}
	mealID, err := uuid.Parse(c.Param("mealId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه وعده نامعتبر است"})
		return
	}

	nutritionistID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "توکن نامعتبر است"})
		return
	}

	resp, err := h.planService.AddOption(c.Request.Context(), planID, dayID, mealID, nutritionistID)
	if err != nil {
		h.handlePlanError(c, err)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// DeleteOption handles DELETE /api/diet-plans/:id/days/:dayId/meals/:mealId/options/:optId.
func (h *DietPlanHandler) DeleteOption(c *gin.Context) {
	planID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه نامعتبر است"})
		return
	}
	dayID, err := uuid.Parse(c.Param("dayId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه روز نامعتبر است"})
		return
	}
	mealID, err := uuid.Parse(c.Param("mealId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه وعده نامعتبر است"})
		return
	}
	optID, err := uuid.Parse(c.Param("optId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه گزینه نامعتبر است"})
		return
	}

	nutritionistID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "توکن نامعتبر است"})
		return
	}

	if err := h.planService.DeleteOption(c.Request.Context(), planID, dayID, mealID, optID, nutritionistID); err != nil {
		h.handlePlanError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ─── Item handlers ────────────────────────────────────────────────────────────

// AddItem handles POST .../options/:optId/items.
func (h *DietPlanHandler) AddItem(c *gin.Context) {
	planID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه نامعتبر است"})
		return
	}
	dayID, err := uuid.Parse(c.Param("dayId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه روز نامعتبر است"})
		return
	}
	mealID, err := uuid.Parse(c.Param("mealId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه وعده نامعتبر است"})
		return
	}
	optID, err := uuid.Parse(c.Param("optId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه گزینه نامعتبر است"})
		return
	}

	var req dto.CreateMealOptionItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "اطلاعات ورودی نامعتبر است"})
		return
	}

	nutritionistID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "توکن نامعتبر است"})
		return
	}

	resp, err := h.planService.AddItem(c.Request.Context(), planID, dayID, mealID, optID, nutritionistID, req)
	if err != nil {
		h.handlePlanError(c, err)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// UpdateItem handles PUT .../options/:optId/items/:itemId.
func (h *DietPlanHandler) UpdateItem(c *gin.Context) {
	planID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه نامعتبر است"})
		return
	}
	dayID, err := uuid.Parse(c.Param("dayId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه روز نامعتبر است"})
		return
	}
	mealID, err := uuid.Parse(c.Param("mealId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه وعده نامعتبر است"})
		return
	}
	optID, err := uuid.Parse(c.Param("optId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه گزینه نامعتبر است"})
		return
	}
	itemID, err := uuid.Parse(c.Param("itemId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه آیتم نامعتبر است"})
		return
	}

	var req dto.UpdateMealOptionItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "اطلاعات ورودی نامعتبر است"})
		return
	}

	nutritionistID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "توکن نامعتبر است"})
		return
	}

	resp, err := h.planService.UpdateItem(c.Request.Context(), planID, dayID, mealID, optID, itemID, nutritionistID, req)
	if err != nil {
		h.handlePlanError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteItem handles DELETE .../options/:optId/items/:itemId.
func (h *DietPlanHandler) DeleteItem(c *gin.Context) {
	planID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه نامعتبر است"})
		return
	}
	dayID, err := uuid.Parse(c.Param("dayId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه روز نامعتبر است"})
		return
	}
	mealID, err := uuid.Parse(c.Param("mealId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه وعده نامعتبر است"})
		return
	}
	optID, err := uuid.Parse(c.Param("optId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه گزینه نامعتبر است"})
		return
	}
	itemID, err := uuid.Parse(c.Param("itemId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه آیتم نامعتبر است"})
		return
	}

	nutritionistID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "توکن نامعتبر است"})
		return
	}

	if err := h.planService.DeleteItem(c.Request.Context(), planID, dayID, mealID, optID, itemID, nutritionistID); err != nil {
		h.handlePlanError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ─── Exercise handlers ────────────────────────────────────────────────────────

// AddExercise handles POST /api/diet-plans/:id/days/:dayId/exercises.
func (h *DietPlanHandler) AddExercise(c *gin.Context) {
	planID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه نامعتبر است"})
		return
	}
	dayID, err := uuid.Parse(c.Param("dayId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه روز نامعتبر است"})
		return
	}

	var req dto.CreateExerciseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "اطلاعات ورودی نامعتبر است"})
		return
	}

	nutritionistID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "توکن نامعتبر است"})
		return
	}

	resp, err := h.planService.AddExercise(c.Request.Context(), planID, dayID, nutritionistID, req)
	if err != nil {
		h.handlePlanError(c, err)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// UpdateExercise handles PUT /api/diet-plans/:id/days/:dayId/exercises/:exId.
func (h *DietPlanHandler) UpdateExercise(c *gin.Context) {
	planID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه نامعتبر است"})
		return
	}
	dayID, err := uuid.Parse(c.Param("dayId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه روز نامعتبر است"})
		return
	}
	exID, err := uuid.Parse(c.Param("exId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه تمرین نامعتبر است"})
		return
	}

	var req dto.UpdateExerciseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "اطلاعات ورودی نامعتبر است"})
		return
	}

	nutritionistID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "توکن نامعتبر است"})
		return
	}

	resp, err := h.planService.UpdateExercise(c.Request.Context(), planID, dayID, exID, nutritionistID, req)
	if err != nil {
		h.handlePlanError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteExercise handles DELETE /api/diet-plans/:id/days/:dayId/exercises/:exId.
func (h *DietPlanHandler) DeleteExercise(c *gin.Context) {
	planID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه نامعتبر است"})
		return
	}
	dayID, err := uuid.Parse(c.Param("dayId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه روز نامعتبر است"})
		return
	}
	exID, err := uuid.Parse(c.Param("exId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه تمرین نامعتبر است"})
		return
	}

	nutritionistID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "توکن نامعتبر است"})
		return
	}

	if err := h.planService.DeleteExercise(c.Request.Context(), planID, dayID, exID, nutritionistID); err != nil {
		h.handlePlanError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ─── Medication handlers ──────────────────────────────────────────────────────

// AddMedication handles POST /api/diet-plans/:id/medications.
func (h *DietPlanHandler) AddMedication(c *gin.Context) {
	planID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه نامعتبر است"})
		return
	}

	var req dto.CreateMedicationPrescriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "اطلاعات ورودی نامعتبر است"})
		return
	}

	nutritionistID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "توکن نامعتبر است"})
		return
	}

	resp, err := h.planService.AddMedication(c.Request.Context(), planID, nutritionistID, req)
	if err != nil {
		h.handlePlanError(c, err)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// UpdateMedication handles PUT /api/diet-plans/:id/medications/:medId.
func (h *DietPlanHandler) UpdateMedication(c *gin.Context) {
	planID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه نامعتبر است"})
		return
	}
	medID, err := uuid.Parse(c.Param("medId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه دارو نامعتبر است"})
		return
	}

	var req dto.UpdateMedicationPrescriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "اطلاعات ورودی نامعتبر است"})
		return
	}

	nutritionistID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "توکن نامعتبر است"})
		return
	}

	resp, err := h.planService.UpdateMedication(c.Request.Context(), planID, medID, nutritionistID, req)
	if err != nil {
		h.handlePlanError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteMedication handles DELETE /api/diet-plans/:id/medications/:medId.
func (h *DietPlanHandler) DeleteMedication(c *gin.Context) {
	planID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه نامعتبر است"})
		return
	}
	medID, err := uuid.Parse(c.Param("medId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه دارو نامعتبر است"})
		return
	}

	nutritionistID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "توکن نامعتبر است"})
		return
	}

	if err := h.planService.DeleteMedication(c.Request.Context(), planID, medID, nutritionistID); err != nil {
		h.handlePlanError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ─── Client handler ───────────────────────────────────────────────────────────

// GetActivePlan handles GET /api/clients/me/active-plan (client role only).
func (h *DietPlanHandler) GetActivePlan(c *gin.Context) {
	clientID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "توکن نامعتبر است"})
		return
	}

	resp, err := h.planService.GetActivePlanForClient(c.Request.Context(), clientID)
	if err != nil {
		if errors.Is(err, service.ErrPlanNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "برنامه‌ای فعال ندارید"})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "خطا در دریافت برنامه"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ─── Error helper ─────────────────────────────────────────────────────────────

// handlePlanError maps service sentinel errors to appropriate HTTP responses.
func (h *DietPlanHandler) handlePlanError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrPlanNotFound),
		errors.Is(err, service.ErrDayNotFound),
		errors.Is(err, service.ErrMealNotFound),
		errors.Is(err, service.ErrOptionNotFound),
		errors.Is(err, service.ErrItemNotFound),
		errors.Is(err, service.ErrExerciseNotFound),
		errors.Is(err, service.ErrMedicationPrescNotFound):
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
	case errors.Is(err, service.ErrPlanUnauthorized):
		c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: err.Error()})
	case errors.Is(err, service.ErrPlanNotDraft):
		c.JSON(http.StatusConflict, dto.ErrorResponse{Error: err.Error()})
	case errors.Is(err, service.ErrPlanIncomplete):
		c.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{Error: err.Error()})
	case errors.Is(err, service.ErrPlanAlreadyActive):
		c.JSON(http.StatusConflict, dto.ErrorResponse{Error: err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "خطا در پردازش درخواست"})
	}
}
