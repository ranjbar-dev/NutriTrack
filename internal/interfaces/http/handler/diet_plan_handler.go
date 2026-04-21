package handler

import (
"time"

"github.com/gin-gonic/gin"
"github.com/google/uuid"
appDietPlan "github.com/ranjbar-dev/nutritrack/internal/application/dietplan"
"github.com/ranjbar-dev/nutritrack/internal/domain/dietplan/entity"
"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
"github.com/ranjbar-dev/nutritrack/internal/interfaces/http/dto"
"github.com/ranjbar-dev/nutritrack/internal/interfaces/http/middleware"
)

// DietPlanHandler handles HTTP requests for diet plan management.
type DietPlanHandler struct {
service *appDietPlan.DietPlanService
}

// NewDietPlanHandler creates a new DietPlanHandler.
func NewDietPlanHandler(service *appDietPlan.DietPlanService) *DietPlanHandler {
return &DietPlanHandler{service: service}
}

// createPlanRequest is the request body for creating a diet plan.
type createPlanRequest struct {
Title              string `json:"title"`
StartDate          string `json:"start_date"` // "2006-01-02"
EndDate            string `json:"end_date"`
Notes              string `json:"notes"`
DailyWaterTargetML int    `json:"daily_water_target_ml"`
}

// addDayRequest is the request body for adding a day to a plan.
type addDayRequest struct {
DayNumber int `json:"day_number"`
}

// addMealRequest is the request body for adding a meal to a plan day.
type addMealRequest struct {
Title         string `json:"title"`
ScheduledTime string `json:"scheduled_time"` // "HH:MM"
DisplayOrder  int    `json:"display_order"`
}

// addOptionRequest is the request body for adding an option to a meal.
type addOptionRequest struct {
OptionNumber int `json:"option_number"`
}

// CreatePlan handles POST /clients/:id/plans
func (h *DietPlanHandler) CreatePlan(c *gin.Context) {
clientIDStr := c.Param("id")
clientID, err := uuid.Parse(clientIDStr)
if err != nil {
dto.Abort(c, shared.ErrValidation)
return
}

callerIDVal, _ := c.Get(middleware.AuthUserIDKey)
callerRoleVal, _ := c.Get(middleware.AuthUserRoleKey)

var req createPlanRequest
if err := c.ShouldBindJSON(&req); err != nil {
dto.Abort(c, shared.ErrValidation)
return
}

startDate, err := time.Parse("2006-01-02", req.StartDate)
if err != nil {
dto.Abort(c, shared.ErrValidation)
return
}
endDate, err := time.Parse("2006-01-02", req.EndDate)
if err != nil {
dto.Abort(c, shared.ErrValidation)
return
}

_ = callerRoleVal // role not needed for create – nutritionist is the caller
plan, svcErr := h.service.CreatePlan(c.Request.Context(), appDietPlan.CreatePlanRequest{
ClientID:           clientID,
NutritionistID:     callerIDVal.(uuid.UUID),
Title:              req.Title,
StartDate:          startDate,
EndDate:            endDate,
Notes:              req.Notes,
DailyWaterTargetML: req.DailyWaterTargetML,
})
if svcErr != nil {
if appErr, ok := svcErr.(*shared.AppError); ok {
dto.Abort(c, appErr)
return
}
dto.Abort(c, shared.ErrInternal)
return
}

dto.Created(c, planToMap(plan))
}

// ListClientPlans handles GET /clients/:id/plans
func (h *DietPlanHandler) ListClientPlans(c *gin.Context) {
clientIDStr := c.Param("id")
clientID, err := uuid.Parse(clientIDStr)
if err != nil {
dto.Abort(c, shared.ErrValidation)
return
}

callerIDVal, _ := c.Get(middleware.AuthUserIDKey)
callerRoleVal, _ := c.Get(middleware.AuthUserRoleKey)

pg := dto.ParsePagination(c)

plans, total, svcErr := h.service.ListClientPlans(
c.Request.Context(),
clientID,
callerIDVal.(uuid.UUID),
callerRoleVal.(string),
pg.Limit(),
pg.Offset(),
)
if svcErr != nil {
if appErr, ok := svcErr.(*shared.AppError); ok {
dto.Abort(c, appErr)
return
}
dto.Abort(c, shared.ErrInternal)
return
}

result := make([]any, len(plans))
for i, p := range plans {
result[i] = planToMap(p)
}

dto.Paginated(c, result, total, pg.Page, pg.PageSize)
}

// GetPlan handles GET /plans/:id (returns full tree)
func (h *DietPlanHandler) GetPlan(c *gin.Context) {
planIDStr := c.Param("id")
planID, err := uuid.Parse(planIDStr)
if err != nil {
dto.Abort(c, shared.ErrValidation)
return
}

callerIDVal, _ := c.Get(middleware.AuthUserIDKey)
callerRoleVal, _ := c.Get(middleware.AuthUserRoleKey)

plan, svcErr := h.service.GetFullPlan(
c.Request.Context(),
planID,
callerIDVal.(uuid.UUID),
callerRoleVal.(string),
)
if svcErr != nil {
if appErr, ok := svcErr.(*shared.AppError); ok {
dto.Abort(c, appErr)
return
}
dto.Abort(c, shared.ErrInternal)
return
}

dto.OK(c, planFullToMap(plan))
}

// AddDay handles POST /plans/:id/days
func (h *DietPlanHandler) AddDay(c *gin.Context) {
planIDStr := c.Param("id")
planID, err := uuid.Parse(planIDStr)
if err != nil {
dto.Abort(c, shared.ErrValidation)
return
}

callerIDVal, _ := c.Get(middleware.AuthUserIDKey)
callerRoleVal, _ := c.Get(middleware.AuthUserRoleKey)

var req addDayRequest
if err := c.ShouldBindJSON(&req); err != nil {
dto.Abort(c, shared.ErrValidation)
return
}

day, svcErr := h.service.AddDay(c.Request.Context(), appDietPlan.AddDayRequest{
PlanID:     planID,
DayNumber:  req.DayNumber,
CallerID:   callerIDVal.(uuid.UUID),
CallerRole: callerRoleVal.(string),
})
if svcErr != nil {
if appErr, ok := svcErr.(*shared.AppError); ok {
dto.Abort(c, appErr)
return
}
dto.Abort(c, shared.ErrInternal)
return
}

dto.Created(c, gin.H{
"id":         day.ID,
"plan_id":    day.PlanID,
"day_number": day.DayNumber,
"created_at": day.CreatedAt,
})
}

// AddMeal handles POST /plans/:id/days/:day_id/meals
func (h *DietPlanHandler) AddMeal(c *gin.Context) {
dayIDStr := c.Param("day_id")
dayID, err := uuid.Parse(dayIDStr)
if err != nil {
dto.Abort(c, shared.ErrValidation)
return
}

callerIDVal, _ := c.Get(middleware.AuthUserIDKey)
callerRoleVal, _ := c.Get(middleware.AuthUserRoleKey)

var req addMealRequest
if err := c.ShouldBindJSON(&req); err != nil {
dto.Abort(c, shared.ErrValidation)
return
}

meal, svcErr := h.service.AddMeal(c.Request.Context(), appDietPlan.AddMealRequest{
DayID:         dayID,
Title:         req.Title,
ScheduledTime: req.ScheduledTime,
DisplayOrder:  req.DisplayOrder,
CallerID:      callerIDVal.(uuid.UUID),
CallerRole:    callerRoleVal.(string),
})
if svcErr != nil {
if appErr, ok := svcErr.(*shared.AppError); ok {
dto.Abort(c, appErr)
return
}
dto.Abort(c, shared.ErrInternal)
return
}

dto.Created(c, gin.H{
"id":             meal.ID,
"day_id":         meal.DayID,
"title":          meal.Title,
"scheduled_time": meal.ScheduledTime,
"display_order":  meal.DisplayOrder,
"created_at":     meal.CreatedAt,
})
}

// AddOption handles POST /plans/:id/days/:day_id/meals/:meal_id/options
func (h *DietPlanHandler) AddOption(c *gin.Context) {
mealIDStr := c.Param("meal_id")
mealID, err := uuid.Parse(mealIDStr)
if err != nil {
dto.Abort(c, shared.ErrValidation)
return
}

callerIDVal, _ := c.Get(middleware.AuthUserIDKey)
callerRoleVal, _ := c.Get(middleware.AuthUserRoleKey)

var req addOptionRequest
if err := c.ShouldBindJSON(&req); err != nil {
dto.Abort(c, shared.ErrValidation)
return
}

option, svcErr := h.service.AddOption(c.Request.Context(), appDietPlan.AddOptionRequest{
MealID:       mealID,
OptionNumber: req.OptionNumber,
CallerID:     callerIDVal.(uuid.UUID),
CallerRole:   callerRoleVal.(string),
})
if svcErr != nil {
if appErr, ok := svcErr.(*shared.AppError); ok {
dto.Abort(c, appErr)
return
}
dto.Abort(c, shared.ErrInternal)
return
}

dto.Created(c, gin.H{
"id":            option.ID,
"meal_id":       option.MealID,
"option_number": option.OptionNumber,
"created_at":    option.CreatedAt,
})
}

// DeletePlan handles DELETE /plans/:id
func (h *DietPlanHandler) DeletePlan(c *gin.Context) {
planIDStr := c.Param("id")
planID, err := uuid.Parse(planIDStr)
if err != nil {
dto.Abort(c, shared.ErrValidation)
return
}

callerIDVal, _ := c.Get(middleware.AuthUserIDKey)
callerRoleVal, _ := c.Get(middleware.AuthUserRoleKey)

svcErr := h.service.DeletePlan(
c.Request.Context(),
planID,
callerIDVal.(uuid.UUID),
callerRoleVal.(string),
)
if svcErr != nil {
if appErr, ok := svcErr.(*shared.AppError); ok {
dto.Abort(c, appErr)
return
}
dto.Abort(c, shared.ErrInternal)
return
}

dto.NoContent(c)
}

// --- Response helpers ---

// planToMap builds a flat plan response map.
func planToMap(p *entity.DietPlan) map[string]any {
return map[string]any{
"id":                    p.ID,
"client_id":             p.ClientID,
"nutritionist_id":       p.NutritionistID,
"title":                 p.Title,
"start_date":            p.StartDate.Format("2006-01-02"),
"end_date":              p.EndDate.Format("2006-01-02"),
"notes":                 p.Notes,
"daily_water_target_ml": p.DailyWaterTargetML,
"status":                p.Status,
"created_at":            p.CreatedAt,
"updated_at":            p.UpdatedAt,
}
}

// planFullToMap builds a full plan response including the days tree.
func planFullToMap(plan *entity.DietPlan) map[string]any {
days := make([]any, len(plan.Days))
for i, day := range plan.Days {
meals := make([]any, len(day.Meals))
for j, meal := range day.Meals {
options := make([]any, len(meal.Options))
for k, opt := range meal.Options {
options[k] = map[string]any{
"id":            opt.ID,
"option_number": opt.OptionNumber,
"items":         opt.Items,
}
}
meals[j] = map[string]any{
"id":             meal.ID,
"title":          meal.Title,
"scheduled_time": meal.ScheduledTime,
"display_order":  meal.DisplayOrder,
"options":        options,
}
}
days[i] = map[string]any{
"id":         day.ID,
"day_number": day.DayNumber,
"meals":      meals,
}
}

return map[string]any{
"id":                    plan.ID,
"client_id":             plan.ClientID,
"nutritionist_id":       plan.NutritionistID,
"title":                 plan.Title,
"start_date":            plan.StartDate.Format("2006-01-02"),
"end_date":              plan.EndDate.Format("2006-01-02"),
"notes":                 plan.Notes,
"daily_water_target_ml": plan.DailyWaterTargetML,
"status":                plan.Status,
"days":                  days,
"created_at":            plan.CreatedAt,
"updated_at":            plan.UpdatedAt,
}
}