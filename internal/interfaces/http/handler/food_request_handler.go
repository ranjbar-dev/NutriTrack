package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	appFoodRequest "github.com/ranjbar-dev/nutritrack/internal/application/foodrequest"
	frEntity "github.com/ranjbar-dev/nutritrack/internal/domain/foodrequest/entity"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	"github.com/ranjbar-dev/nutritrack/internal/interfaces/http/dto"
	"github.com/ranjbar-dev/nutritrack/internal/interfaces/http/middleware"
)

// FoodRequestHandler handles HTTP requests for the food-request lifecycle.
type FoodRequestHandler struct {
	service *appFoodRequest.FoodRequestService
}

// NewFoodRequestHandler creates a new FoodRequestHandler.
func NewFoodRequestHandler(service *appFoodRequest.FoodRequestService) *FoodRequestHandler {
	return &FoodRequestHandler{service: service}
}

// Submit handles POST /food-requests — client submits a food request to their nutritionist.
func (h *FoodRequestHandler) Submit(c *gin.Context) {
	callerIDVal, _ := c.Get(middleware.AuthUserIDKey)
	callerRoleVal, _ := c.Get(middleware.AuthUserRoleKey)

	if callerRoleVal.(string) != "client" {
		dto.Abort(c, shared.ErrForbidden)
		return
	}

	var body struct {
		FoodName string `json:"food_name"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.FoodName == "" {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	result, svcErr := h.service.Submit(c.Request.Context(), callerIDVal.(uuid.UUID), body.FoodName)
	if svcErr != nil {
		if appErr, ok := svcErr.(*shared.AppError); ok {
			dto.Abort(c, appErr)
			return
		}
		dto.Abort(c, shared.ErrInternal)
		return
	}

	dto.Created(c, foodRequestToMap(result))
}

// ListPending handles GET /food-requests — nutritionist lists pending food requests.
func (h *FoodRequestHandler) ListPending(c *gin.Context) {
	callerIDVal, _ := c.Get(middleware.AuthUserIDKey)
	callerRoleVal, _ := c.Get(middleware.AuthUserRoleKey)

	if callerRoleVal.(string) != "nutritionist" {
		dto.Abort(c, shared.ErrForbidden)
		return
	}

	pg := dto.ParsePagination(c)
	items, total, svcErr := h.service.ListPending(
		c.Request.Context(),
		callerIDVal.(uuid.UUID),
		int32(pg.Limit()),
		int32(pg.Offset()),
	)
	if svcErr != nil {
		if appErr, ok := svcErr.(*shared.AppError); ok {
			dto.Abort(c, appErr)
			return
		}
		dto.Abort(c, shared.ErrInternal)
		return
	}

	result := make([]map[string]any, len(items))
	for i, item := range items {
		result[i] = foodRequestToMap(item)
	}
	dto.Paginated(c, result, total, pg.Page, pg.PageSize)
}

// Approve handles POST /food-requests/:id/approve — nutritionist approves a food request.
func (h *FoodRequestHandler) Approve(c *gin.Context) {
	requestIDStr := c.Param("id")
	requestID, err := uuid.Parse(requestIDStr)
	if err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	callerIDVal, _ := c.Get(middleware.AuthUserIDKey)
	callerRoleVal, _ := c.Get(middleware.AuthUserRoleKey)

	if callerRoleVal.(string) != "nutritionist" {
		dto.Abort(c, shared.ErrForbidden)
		return
	}

	var body struct {
		Name         string  `json:"name"`
		Unit         string  `json:"unit"`
		Calories     float64 `json:"calories"`
		Protein      float64 `json:"protein"`
		Carbohydrate float64 `json:"carbohydrate"`
		Fat          float64 `json:"fat"`
		Fiber        float64 `json:"fiber"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	result, svcErr := h.service.Approve(c.Request.Context(), requestID, callerIDVal.(uuid.UUID), appFoodRequest.ApproveRequest{
		Name:         body.Name,
		Unit:         body.Unit,
		Calories:     body.Calories,
		Protein:      body.Protein,
		Carbohydrate: body.Carbohydrate,
		Fat:          body.Fat,
		Fiber:        body.Fiber,
	})
	if svcErr != nil {
		if appErr, ok := svcErr.(*shared.AppError); ok {
			dto.Abort(c, appErr)
			return
		}
		dto.Abort(c, shared.ErrInternal)
		return
	}

	dto.OK(c, foodRequestToMap(result))
}

// Reject handles POST /food-requests/:id/reject — nutritionist rejects a food request.
func (h *FoodRequestHandler) Reject(c *gin.Context) {
	requestIDStr := c.Param("id")
	requestID, err := uuid.Parse(requestIDStr)
	if err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	callerIDVal, _ := c.Get(middleware.AuthUserIDKey)
	callerRoleVal, _ := c.Get(middleware.AuthUserRoleKey)

	if callerRoleVal.(string) != "nutritionist" {
		dto.Abort(c, shared.ErrForbidden)
		return
	}

	var body struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	result, svcErr := h.service.Reject(c.Request.Context(), requestID, callerIDVal.(uuid.UUID), body.Reason)
	if svcErr != nil {
		if appErr, ok := svcErr.(*shared.AppError); ok {
			dto.Abort(c, appErr)
			return
		}
		dto.Abort(c, shared.ErrInternal)
		return
	}

	dto.OK(c, foodRequestToMap(result))
}

// foodRequestToMap converts a FoodRequest entity to a JSON-friendly map.
func foodRequestToMap(r *frEntity.FoodRequest) map[string]any {
	return map[string]any{
		"id":               r.ID,
		"client_id":        r.ClientID,
		"nutritionist_id":  r.NutritionistID,
		"food_name":        r.FoodName,
		"status":           string(r.Status),
		"rejection_reason": r.RejectionReason,
		"created_food_id":  r.CreatedFoodID,
		"created_at":       r.CreatedAt,
		"updated_at":       r.UpdatedAt,
	}
}
