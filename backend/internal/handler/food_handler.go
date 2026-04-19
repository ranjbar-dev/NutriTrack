package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ranjbar-dev/nutritrack/backend/internal/model/dto"
	"github.com/ranjbar-dev/nutritrack/backend/internal/service"
)

// FoodHandler handles food CRUD HTTP endpoints.
type FoodHandler struct {
	foodService *service.FoodService
}

// NewFoodHandler creates a new FoodHandler with the given food service.
func NewFoodHandler(foodService *service.FoodService) *FoodHandler {
	return &FoodHandler{foodService: foodService}
}

// Create handles POST /api/foods.
func (h *FoodHandler) Create(c *gin.Context) {
	var req dto.CreateFoodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "اطلاعات ورودی نامعتبر است"})
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "توکن نامعتبر است"})
		return
	}

	resp, err := h.foodService.CreateFood(c.Request.Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrFoodDuplicate):
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "خطا در ثبت غذا"})
		}
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// Get handles GET /api/foods/:id.
func (h *FoodHandler) Get(c *gin.Context) {
	foodID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "اطلاعات ورودی نامعتبر است"})
		return
	}

	resp, err := h.foodService.GetFood(c.Request.Context(), foodID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrFoodNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "خطا در دریافت غذا"})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

// List handles GET /api/foods.
func (h *FoodHandler) List(c *gin.Context) {
	var query dto.FoodListQueryParams
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "اطلاعات ورودی نامعتبر است"})
		return
	}

	resp, err := h.foodService.ListFoods(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "خطا در دریافت لیست غذاها"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Update handles PUT /api/foods/:id.
func (h *FoodHandler) Update(c *gin.Context) {
	foodID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "اطلاعات ورودی نامعتبر است"})
		return
	}

	var req dto.UpdateFoodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "اطلاعات ورودی نامعتبر است"})
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "توکن نامعتبر است"})
		return
	}

	resp, err := h.foodService.UpdateFood(c.Request.Context(), foodID, userID, c.GetString("role"), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrFoodNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		case errors.Is(err, service.ErrFoodUnauthorizedEdit):
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: err.Error()})
		case errors.Is(err, service.ErrFoodDuplicate):
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "خطا در ویرایش غذا"})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Delete handles DELETE /api/foods/:id.
func (h *FoodHandler) Delete(c *gin.Context) {
	foodID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "اطلاعات ورودی نامعتبر است"})
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "توکن نامعتبر است"})
		return
	}

	err = h.foodService.DeleteFood(c.Request.Context(), foodID, userID, c.GetString("role"))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrFoodNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		case errors.Is(err, service.ErrFoodUnauthorizedDelete):
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "خطا در حذف غذا"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "غذا با موفقیت حذف شد"})
}
