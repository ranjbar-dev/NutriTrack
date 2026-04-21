package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	appFood "github.com/ranjbar-dev/nutritrack/internal/application/food"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	"github.com/ranjbar-dev/nutritrack/internal/interfaces/http/dto"
)

// FoodCategoryHandler handles food category CRUD endpoints.
type FoodCategoryHandler struct {
	svc *appFood.FoodCategoryService
}

// NewFoodCategoryHandler constructs a FoodCategoryHandler.
func NewFoodCategoryHandler(svc *appFood.FoodCategoryService) *FoodCategoryHandler {
	return &FoodCategoryHandler{svc: svc}
}

// Create handles POST /admin/food-categories
func (h *FoodCategoryHandler) Create(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	cat, err := h.svc.Create(c.Request.Context(), req.Name)
	if err != nil {
		appErr, ok := err.(*shared.AppError)
		if !ok {
			appErr = shared.ErrInternal
		}
		dto.Abort(c, appErr)
		return
	}

	dto.Created(c, gin.H{"id": cat.ID, "name": cat.Name, "created_at": cat.CreatedAt})
}

// ListAll handles GET /food-categories
func (h *FoodCategoryHandler) ListAll(c *gin.Context) {
	cats, err := h.svc.ListAll(c.Request.Context())
	if err != nil {
		dto.Abort(c, shared.ErrInternal)
		return
	}

	result := make([]gin.H, len(cats))
	for i, cat := range cats {
		result[i] = gin.H{"id": cat.ID, "name": cat.Name}
	}
	dto.OK(c, result)
}

// Delete handles DELETE /admin/food-categories/:id
func (h *FoodCategoryHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		appErr, ok := err.(*shared.AppError)
		if !ok {
			appErr = shared.ErrInternal
		}
		dto.Abort(c, appErr)
		return
	}

	dto.OK(c, gin.H{"message": "دسته‌بندی با موفقیت حذف شد"})
}
