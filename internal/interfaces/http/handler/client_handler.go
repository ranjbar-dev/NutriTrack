package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	appUser "github.com/ranjbar-dev/nutritrack/internal/application/user"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	"github.com/ranjbar-dev/nutritrack/internal/domain/user/entity"
	"github.com/ranjbar-dev/nutritrack/internal/interfaces/http/dto"
	"github.com/ranjbar-dev/nutritrack/internal/interfaces/http/middleware"
)

// ClientHandler handles nutritionist-scoped client management endpoints.
type ClientHandler struct {
	svc *appUser.ClientService
}

// NewClientHandler constructs a ClientHandler.
func NewClientHandler(svc *appUser.ClientService) *ClientHandler {
	return &ClientHandler{svc: svc}
}

// RegisterClient handles POST /api/v1/clients.
func (h *ClientHandler) RegisterClient(c *gin.Context) {
	nutIDRaw, _ := c.Get(middleware.AuthUserIDKey)
	nutritionistID, ok := nutIDRaw.(uuid.UUID)
	if !ok {
		dto.Abort(c, shared.ErrUnauthorized)
		return
	}

	var req struct {
		Mobile    string  `json:"mobile"      binding:"required"`
		FirstName string  `json:"first_name"  binding:"required"`
		LastName  string  `json:"last_name"   binding:"required"`
		Gender    string  `json:"gender"`
		BirthDate string  `json:"birth_date"` // "2006-01-02"
		Height    *float64 `json:"height"`
		Weight    *float64 `json:"weight"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	var birthDate *time.Time
	if req.BirthDate != "" {
		t, err := time.Parse("2006-01-02", req.BirthDate)
		if err != nil {
			dto.Abort(c, shared.ErrValidation)
			return
		}
		birthDate = &t
	}

	user, err := h.svc.RegisterClient(c.Request.Context(), appUser.RegisterClientRequest{
		NutritionistID: nutritionistID,
		Mobile:         req.Mobile,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Gender:         req.Gender,
		BirthDate:      birthDate,
		Height:         req.Height,
		Weight:         req.Weight,
	})
	if err != nil {
		appErr, ok := err.(*shared.AppError)
		if !ok {
			appErr = shared.ErrInternal
		}
		dto.Abort(c, appErr)
		return
	}
	dto.Created(c, toClientResponse(user))
}

// ListClients handles GET /api/v1/clients.
func (h *ClientHandler) ListClients(c *gin.Context) {
	nutIDRaw, _ := c.Get(middleware.AuthUserIDKey)
	nutritionistID, ok := nutIDRaw.(uuid.UUID)
	if !ok {
		dto.Abort(c, shared.ErrUnauthorized)
		return
	}

	pg := dto.ParsePagination(c)

	users, total, err := h.svc.ListClients(c.Request.Context(), nutritionistID, pg.Limit(), pg.Offset())
	if err != nil {
		dto.Abort(c, shared.ErrInternal)
		return
	}

	resp := make([]gin.H, len(users))
	for i, u := range users {
		resp[i] = toClientResponse(u)
	}
	dto.Paginated(c, resp, total, pg.Page, pg.PageSize)
}

// GetClientProfile handles GET /api/v1/clients/:id.
func (h *ClientHandler) GetClientProfile(c *gin.Context) {
	nutIDRaw, _ := c.Get(middleware.AuthUserIDKey)
	nutritionistID, ok := nutIDRaw.(uuid.UUID)
	if !ok {
		dto.Abort(c, shared.ErrUnauthorized)
		return
	}

	clientID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	user, err := h.svc.GetClientProfile(c.Request.Context(), clientID, nutritionistID)
	if err != nil {
		appErr, ok := err.(*shared.AppError)
		if !ok {
			appErr = shared.ErrInternal
		}
		dto.Abort(c, appErr)
		return
	}
	dto.OK(c, toClientResponse(user))
}

// UpdateClient handles PATCH /api/v1/clients/:id.
func (h *ClientHandler) UpdateClient(c *gin.Context) {
	nutIDRaw, _ := c.Get(middleware.AuthUserIDKey)
	nutritionistID, ok := nutIDRaw.(uuid.UUID)
	if !ok {
		dto.Abort(c, shared.ErrUnauthorized)
		return
	}

	clientID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	var req struct {
		FirstName string   `json:"first_name"`
		LastName  string   `json:"last_name"`
		Gender    string   `json:"gender"`
		BirthDate string   `json:"birth_date"` // "2006-01-02"
		Height    *float64 `json:"height"`
		Weight    *float64 `json:"weight"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	var birthDate *time.Time
	if req.BirthDate != "" {
		t, err := time.Parse("2006-01-02", req.BirthDate)
		if err != nil {
			dto.Abort(c, shared.ErrValidation)
			return
		}
		birthDate = &t
	}

	user, err := h.svc.UpdateClient(c.Request.Context(), appUser.UpdateClientRequest{
		ClientID:       clientID,
		NutritionistID: nutritionistID,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Gender:         req.Gender,
		BirthDate:      birthDate,
		Height:         req.Height,
		Weight:         req.Weight,
	})
	if err != nil {
		appErr, ok := err.(*shared.AppError)
		if !ok {
			appErr = shared.ErrInternal
		}
		dto.Abort(c, appErr)
		return
	}
	dto.OK(c, toClientResponse(user))
}

// SetStatus handles PATCH /api/v1/clients/:id/status.
func (h *ClientHandler) SetStatus(c *gin.Context) {
	nutIDRaw, _ := c.Get(middleware.AuthUserIDKey)
	nutritionistID, ok := nutIDRaw.(uuid.UUID)
	if !ok {
		dto.Abort(c, shared.ErrUnauthorized)
		return
	}

	clientID, err := uuid.Parse(c.Param("id"))
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

	if err := h.svc.SetClientStatus(c.Request.Context(), clientID, nutritionistID, req.IsActive); err != nil {
		appErr, ok := err.(*shared.AppError)
		if !ok {
			appErr = shared.ErrInternal
		}
		dto.Abort(c, appErr)
		return
	}
	dto.OK(c, gin.H{"message": "وضعیت مراجع با موفقیت به‌روز شد"})
}

// toClientResponse converts a domain User (client) to a JSON-serialisable map.
func toClientResponse(u *entity.User) gin.H {
	var birthDate *string
	if u.BirthDate != nil {
		s := u.BirthDate.Format("2006-01-02")
		birthDate = &s
	}

	var nutID *string
	if u.NutritionistID != nil {
		s := u.NutritionistID.String()
		nutID = &s
	}

	return gin.H{
		"id":              u.ID,
		"mobile":          u.Mobile,
		"first_name":      u.FirstName,
		"last_name":       u.LastName,
		"full_name":       u.FullName(),
		"gender":          u.Gender,
		"birth_date":      birthDate,
		"height":          u.Height,
		"weight":          u.Weight,
		"bmi":             u.BMI(),
		"avatar_url":      u.AvatarURL,
		"is_active":       u.IsActive,
		"nutritionist_id": nutID,
		"created_at":      u.CreatedAt,
		"updated_at":      u.UpdatedAt,
	}
}
