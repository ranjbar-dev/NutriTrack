package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	appUser "github.com/ranjbar-dev/nutritrack/internal/application/user"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	"github.com/ranjbar-dev/nutritrack/internal/domain/user/entity"
	"github.com/ranjbar-dev/nutritrack/internal/interfaces/http/dto"
)

// NutritionistHandler handles super-admin nutritionist management endpoints.
type NutritionistHandler struct {
	svc *appUser.NutritionistService
}

// NewNutritionistHandler constructs a NutritionistHandler.
func NewNutritionistHandler(svc *appUser.NutritionistService) *NutritionistHandler {
	return &NutritionistHandler{svc: svc}
}

// Create handles POST /api/v1/admin/nutritionists.
func (h *NutritionistHandler) Create(c *gin.Context) {
	var req struct {
		Email     string `json:"email"      binding:"required,email"`
		Password  string `json:"password"   binding:"required,min=8"`
		FirstName string `json:"first_name" binding:"required"`
		LastName  string `json:"last_name"  binding:"required"`
		Mobile    string `json:"mobile"     binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	user, err := h.svc.Create(c.Request.Context(), appUser.CreateNutritionistRequest{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Mobile:    req.Mobile,
	})
	if err != nil {
		appErr, ok := err.(*shared.AppError)
		if !ok {
			appErr = shared.ErrInternal
		}
		dto.Abort(c, appErr)
		return
	}
	dto.Created(c, toNutritionistResponse(user))
}

// Get handles GET /api/v1/admin/nutritionists/:id.
func (h *NutritionistHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	user, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		appErr, ok := err.(*shared.AppError)
		if !ok {
			appErr = shared.ErrInternal
		}
		dto.Abort(c, appErr)
		return
	}
	dto.OK(c, toNutritionistResponse(user))
}

// List handles GET /api/v1/admin/nutritionists.
func (h *NutritionistHandler) List(c *gin.Context) {
	pg := dto.ParsePagination(c)

	users, total, err := h.svc.List(c.Request.Context(), pg.Limit(), pg.Offset())
	if err != nil {
		dto.Abort(c, shared.ErrInternal)
		return
	}

	resp := make([]gin.H, len(users))
	for i, u := range users {
		resp[i] = toNutritionistResponse(u)
	}
	dto.Paginated(c, resp, total, pg.Page, pg.PageSize)
}

// Update handles PATCH /api/v1/admin/nutritionists/:id.
func (h *NutritionistHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	user, err := h.svc.Update(c.Request.Context(), appUser.UpdateNutritionistRequest{
		ID:        id,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})
	if err != nil {
		appErr, ok := err.(*shared.AppError)
		if !ok {
			appErr = shared.ErrInternal
		}
		dto.Abort(c, appErr)
		return
	}
	dto.OK(c, toNutritionistResponse(user))
}

// SetStatus handles PATCH /api/v1/admin/nutritionists/:id/status.
func (h *NutritionistHandler) SetStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	if err := h.svc.SetStatus(c.Request.Context(), id, req.IsActive); err != nil {
		appErr, ok := err.(*shared.AppError)
		if !ok {
			appErr = shared.ErrInternal
		}
		dto.Abort(c, appErr)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "وضعیت متخصص تغذیه با موفقیت به‌روز شد"})
}

// toNutritionistResponse converts a domain User to a JSON-serialisable map.
func toNutritionistResponse(u *entity.User) gin.H {
	return gin.H{
		"id":         u.ID,
		"email":      u.Email,
		"mobile":     u.Mobile,
		"first_name": u.FirstName,
		"last_name":  u.LastName,
		"is_active":  u.IsActive,
		"created_at": u.CreatedAt,
	}
}
