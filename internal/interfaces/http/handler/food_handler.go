package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	appFood "github.com/ranjbar-dev/nutritrack/internal/application/food"
	"github.com/ranjbar-dev/nutritrack/internal/domain/food/entity"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	"github.com/ranjbar-dev/nutritrack/internal/interfaces/http/dto"
	"github.com/ranjbar-dev/nutritrack/internal/interfaces/http/middleware"
)

// FoodHandler handles food CRUD and search endpoints.
type FoodHandler struct {
	svc *appFood.FoodService
}

// NewFoodHandler constructs a FoodHandler.
func NewFoodHandler(svc *appFood.FoodService) *FoodHandler {
	return &FoodHandler{svc: svc}
}

// foodRequest is the shared request body for POST and PATCH.
type foodRequest struct {
	Name         string      `json:"name"          binding:"required"`
	Unit         string      `json:"unit"          binding:"required"`
	Calories     float64     `json:"calories"      binding:"gte=0"`
	Protein      float64     `json:"protein"       binding:"gte=0"`
	Carbohydrate float64     `json:"carbohydrate"  binding:"gte=0"`
	Fat          float64     `json:"fat"           binding:"gte=0"`
	Fiber        float64     `json:"fiber"         binding:"gte=0"`
	CategoryIDs  []uuid.UUID `json:"category_ids"`
}

// Create handles POST /api/v1/foods.
func (h *FoodHandler) Create(c *gin.Context) {
	var req foodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	callerID, callerRole := callerContext(c)

	food, err := h.svc.CreateFood(c.Request.Context(), appFood.CreateFoodRequest{
		Name:         req.Name,
		Unit:         req.Unit,
		Calories:     req.Calories,
		Protein:      req.Protein,
		Carbohydrate: req.Carbohydrate,
		Fat:          req.Fat,
		Fiber:        req.Fiber,
		CategoryIDs:  req.CategoryIDs,
		CallerID:     callerID,
		CallerRole:   callerRole,
	})
	if err != nil {
		dto.Abort(c, toAppError(err))
		return
	}

	dto.Created(c, toFoodResponse(food))
}

// GetOne handles GET /api/v1/foods/:id.
func (h *FoodHandler) GetOne(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	food, err := h.svc.GetFood(c.Request.Context(), id)
	if err != nil {
		dto.Abort(c, toAppError(err))
		return
	}

	dto.OK(c, toFoodResponse(food))
}

// Search handles GET /api/v1/foods.
func (h *FoodHandler) Search(c *gin.Context) {
	query := c.DefaultQuery("q", "")
	pg := dto.ParsePagination(c)

	var categoryID *uuid.UUID
	if catStr := c.Query("category_id"); catStr != "" {
		parsed, err := uuid.Parse(catStr)
		if err != nil {
			dto.Abort(c, shared.ErrValidation)
			return
		}
		categoryID = &parsed
	}

	foods, total, err := h.svc.SearchFoods(c.Request.Context(), query, categoryID, pg.Limit(), pg.Offset())
	if err != nil {
		dto.Abort(c, toAppError(err))
		return
	}

	resp := make([]gin.H, len(foods))
	for i, f := range foods {
		resp[i] = toFoodResponse(f)
	}
	dto.Paginated(c, resp, total, pg.Page, pg.PageSize)
}

// Update handles PATCH /api/v1/foods/:id.
func (h *FoodHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	var req foodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	callerID, callerRole := callerContext(c)

	food, err := h.svc.UpdateFood(c.Request.Context(), appFood.UpdateFoodRequest{
		ID:           id,
		Name:         req.Name,
		Unit:         req.Unit,
		Calories:     req.Calories,
		Protein:      req.Protein,
		Carbohydrate: req.Carbohydrate,
		Fat:          req.Fat,
		Fiber:        req.Fiber,
		CategoryIDs:  req.CategoryIDs,
		CallerID:     callerID,
		CallerRole:   callerRole,
	})
	if err != nil {
		dto.Abort(c, toAppError(err))
		return
	}

	dto.OK(c, toFoodResponse(food))
}

// Delete handles DELETE /api/v1/foods/:id.
func (h *FoodHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	callerID, callerRole := callerContext(c)

	if err := h.svc.DeleteFood(c.Request.Context(), id, callerID, callerRole); err != nil {
		dto.Abort(c, toAppError(err))
		return
	}

	c.Status(http.StatusNoContent)
}

// --- helpers ---

// callerContext extracts the authenticated user ID and role from gin context.
func callerContext(c *gin.Context) (uuid.UUID, string) {
	callerIDRaw, _ := c.Get(middleware.AuthUserIDKey)
	callerID, _ := callerIDRaw.(uuid.UUID)

	callerRoleRaw, _ := c.Get(middleware.AuthUserRoleKey)
	callerRole, _ := callerRoleRaw.(string)

	return callerID, callerRole
}

// toAppError converts any error to *shared.AppError, defaulting to ErrInternal.
func toAppError(err error) *shared.AppError {
	appErr, ok := err.(*shared.AppError)
	if !ok {
		return shared.ErrInternal
	}
	return appErr
}

// toFoodResponse converts a domain Food to a JSON-serialisable map.
func toFoodResponse(f *entity.Food) gin.H {
	cats := make([]gin.H, len(f.Categories))
	for i, c := range f.Categories {
		cats[i] = gin.H{
			"id":   c.ID,
			"name": c.Name,
		}
	}

	resp := gin.H{
		"id":           f.ID,
		"name":         f.Name,
		"unit":         f.Unit,
		"calories":     f.Calories,
		"protein":      f.Protein,
		"carbohydrate": f.Carbohydrate,
		"fat":          f.Fat,
		"fiber":        f.Fiber,
		"is_active":    f.IsActive,
		"categories":   cats,
		"created_at":   f.CreatedAt,
		"updated_at":   f.UpdatedAt,
	}

	if f.CreatedBy != nil {
		resp["created_by"] = *f.CreatedBy
	} else {
		resp["created_by"] = nil
	}

	return resp
}
