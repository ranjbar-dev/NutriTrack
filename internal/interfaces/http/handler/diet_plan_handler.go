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

// addExerciseRequest is the request body for adding an exercise recommendation to a plan day.
type addExerciseRequest struct {
	ExerciseName         string `json:"exercise_name"`
	DurationMinutes      int    `json:"duration_minutes"`
	Description          string `json:"description"`
	CaloriesBurnEstimate int    `json:"calories_burn_estimate"`
}

// addPrescriptionRequest is the request body for adding a medication prescription to a plan day.
type addPrescriptionRequest struct {
	MedicationID string   `json:"medication_id"`
	Dosage       string   `json:"dosage"`
	Frequency    string   `json:"frequency"`
	Times        []string `json:"times"`
	Instructions string   `json:"instructions"`
	StartDate    string   `json:"start_date"` // "2006-01-02" or ""
	EndDate      string   `json:"end_date"`
}

// addItemRequest is the request body for adding a food item to a meal option.
type addItemRequest struct {
FoodID   string  `json:"food_id"`
Quantity float64 `json:"quantity"`
Unit     string  `json:"unit"`
Notes    string  `json:"notes"`
}

// AddExercise handles POST /plans/:id/days/:day_id/exercises
func (h *DietPlanHandler) AddExercise(c *gin.Context) {
dayIDStr := c.Param("day_id")
dayID, err := uuid.Parse(dayIDStr)
if err != nil {
dto.Abort(c, shared.ErrValidation)
return
}

callerIDVal, _ := c.Get(middleware.AuthUserIDKey)
callerRoleVal, _ := c.Get(middleware.AuthUserRoleKey)

var req addExerciseRequest
if err := c.ShouldBindJSON(&req); err != nil {
dto.Abort(c, shared.ErrValidation)
return
}

ex, svcErr := h.service.AddExercise(c.Request.Context(), appDietPlan.AddExerciseRequest{
DayID:                dayID,
ExerciseName:         req.ExerciseName,
DurationMinutes:      req.DurationMinutes,
Description:          req.Description,
CaloriesBurnEstimate: req.CaloriesBurnEstimate,
CallerID:             callerIDVal.(uuid.UUID),
CallerRole:           callerRoleVal.(string),
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
"id":                     ex.ID,
"day_id":                 ex.DayID,
"exercise_name":          ex.ExerciseName,
"duration_minutes":       ex.DurationMinutes,
"description":            ex.Description,
"calories_burn_estimate": ex.CaloriesBurnEstimate,
"created_at":             ex.CreatedAt,
})
}

// RemoveExercise handles DELETE /plans/:id/days/:day_id/exercises/:exercise_id
func (h *DietPlanHandler) RemoveExercise(c *gin.Context) {
exerciseIDStr := c.Param("exercise_id")
exerciseID, err := uuid.Parse(exerciseIDStr)
if err != nil {
dto.Abort(c, shared.ErrValidation)
return
}

callerIDVal, _ := c.Get(middleware.AuthUserIDKey)
callerRoleVal, _ := c.Get(middleware.AuthUserRoleKey)

svcErr := h.service.RemoveExercise(c.Request.Context(),
exerciseID,
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

// AddPrescription handles POST /plans/:id/days/:day_id/prescriptions
func (h *DietPlanHandler) AddPrescription(c *gin.Context) {
dayIDStr := c.Param("day_id")
dayID, err := uuid.Parse(dayIDStr)
if err != nil {
dto.Abort(c, shared.ErrValidation)
return
}

callerIDVal, _ := c.Get(middleware.AuthUserIDKey)
callerRoleVal, _ := c.Get(middleware.AuthUserRoleKey)

var req addPrescriptionRequest
if err := c.ShouldBindJSON(&req); err != nil {
dto.Abort(c, shared.ErrValidation)
return
}

medicationID, err := uuid.Parse(req.MedicationID)
if err != nil {
dto.Abort(c, shared.ErrValidation)
return
}

var startDate *time.Time
if req.StartDate != "" {
t, err := time.Parse("2006-01-02", req.StartDate)
if err != nil {
dto.Abort(c, shared.ErrValidation)
return
}
startDate = &t
}

var endDate *time.Time
if req.EndDate != "" {
t, err := time.Parse("2006-01-02", req.EndDate)
if err != nil {
dto.Abort(c, shared.ErrValidation)
return
}
endDate = &t
}

times := req.Times
if times == nil {
times = []string{}
}

rx, svcErr := h.service.AddPrescription(c.Request.Context(), appDietPlan.AddPrescriptionRequest{
DayID:        dayID,
MedicationID: medicationID,
Dosage:       req.Dosage,
Frequency:    req.Frequency,
Times:        times,
Instructions: req.Instructions,
StartDate:    startDate,
EndDate:      endDate,
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
"id":            rx.ID,
"day_id":        rx.DayID,
"medication_id": rx.MedicationID,
"dosage":        rx.Dosage,
"frequency":     rx.Frequency,
"times":         rx.Times,
"instructions":  rx.Instructions,
"start_date":    rx.StartDate,
"end_date":      rx.EndDate,
"created_at":    rx.CreatedAt,
})
}

// RemovePrescription handles DELETE /plans/:id/days/:day_id/prescriptions/:prescription_id
func (h *DietPlanHandler) RemovePrescription(c *gin.Context) {
prescriptionIDStr := c.Param("prescription_id")
prescriptionID, err := uuid.Parse(prescriptionIDStr)
if err != nil {
dto.Abort(c, shared.ErrValidation)
return
}

callerIDVal, _ := c.Get(middleware.AuthUserIDKey)
callerRoleVal, _ := c.Get(middleware.AuthUserRoleKey)

svcErr := h.service.RemovePrescription(c.Request.Context(),
prescriptionID,
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

// AddItem handles POST /plans/:id/days/:day_id/meals/:meal_id/options/:option_id/items
func (h *DietPlanHandler) AddItem(c *gin.Context) {
optionIDStr := c.Param("option_id")
optionID, err := uuid.Parse(optionIDStr)
if err != nil {
dto.Abort(c, shared.ErrValidation)
return
}

callerIDVal, _ := c.Get(middleware.AuthUserIDKey)
callerRoleVal, _ := c.Get(middleware.AuthUserRoleKey)

var req addItemRequest
if err := c.ShouldBindJSON(&req); err != nil {
dto.Abort(c, shared.ErrValidation)
return
}

foodID, err := uuid.Parse(req.FoodID)
if err != nil {
dto.Abort(c, shared.ErrValidation)
return
}

item, svcErr := h.service.AddItem(c.Request.Context(), appDietPlan.AddItemRequest{
OptionID:   optionID,
FoodID:     foodID,
Quantity:   req.Quantity,
Unit:       req.Unit,
Notes:      req.Notes,
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
"id":         item.ID,
"option_id":  item.OptionID,
"food_id":    item.FoodID,
"quantity":   item.Quantity,
"unit":       item.Unit,
"notes":      item.Notes,
"created_at": item.CreatedAt,
})
}

// RemoveItem handles DELETE /plans/:id/days/:day_id/meals/:meal_id/options/:option_id/items/:item_id
func (h *DietPlanHandler) RemoveItem(c *gin.Context) {
itemIDStr := c.Param("item_id")
itemID, err := uuid.Parse(itemIDStr)
if err != nil {
dto.Abort(c, shared.ErrValidation)
return
}

callerIDVal, _ := c.Get(middleware.AuthUserIDKey)
callerRoleVal, _ := c.Get(middleware.AuthUserRoleKey)

svcErr := h.service.RemoveItem(c.Request.Context(),
itemID,
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
"totals":        opt.Totals,
"items":         itemsToSlice(opt.Items),
}
}
meals[j] = map[string]any{
"id":             meal.ID,
"title":          meal.Title,
"scheduled_time": meal.ScheduledTime,
"display_order":  meal.DisplayOrder,
"total_range":    meal.TotalRange,
"options":        options,
}
}
days[i] = map[string]any{
"id":            day.ID,
"day_number":    day.DayNumber,
"total_range":   day.TotalRange,
"meals":         meals,
"exercises":     exercisesToSlice(day.Exercises),
"prescriptions": prescriptionsToSlice(day.Prescriptions),
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

// itemsToSlice converts a slice of MealOptionItem entities to a slice of maps for JSON response.
func itemsToSlice(items []*entity.MealOptionItem) []any {
result := make([]any, len(items))
for i, item := range items {
m := map[string]any{
"id":         item.ID,
"option_id":  item.OptionID,
"food_id":    item.FoodID,
"quantity":   item.Quantity,
"unit":       item.Unit,
"notes":      item.Notes,
"computed":   item.Computed,
"created_at": item.CreatedAt,
}
if item.Food != nil {
m["food"] = map[string]any{
"id":   item.Food.ID,
"name": item.Food.Name,
"unit": item.Food.Unit,
}
}
result[i] = m
}
return result
}

// exercisesToSlice converts a slice of ExerciseRecommendation entities to a slice of maps for JSON response.
func exercisesToSlice(exercises []*entity.ExerciseRecommendation) []any {
result := make([]any, len(exercises))
for i, e := range exercises {
result[i] = map[string]any{
"id":                     e.ID,
"exercise_name":          e.ExerciseName,
"duration_minutes":       e.DurationMinutes,
"description":            e.Description,
"calories_burn_estimate": e.CaloriesBurnEstimate,
"created_at":             e.CreatedAt,
}
}
return result
}

// prescriptionsToSlice converts a slice of PrescribedMedication entities to a slice of maps for JSON response.
func prescriptionsToSlice(prescriptions []*entity.PrescribedMedication) []any {
result := make([]any, len(prescriptions))
for i, rx := range prescriptions {
m := map[string]any{
"id":            rx.ID,
"medication_id": rx.MedicationID,
"dosage":        rx.Dosage,
"frequency":     rx.Frequency,
"times":         rx.Times,
"instructions":  rx.Instructions,
"start_date":    rx.StartDate,
"end_date":      rx.EndDate,
"created_at":    rx.CreatedAt,
}
if rx.Medication != nil {
m["medication"] = map[string]any{
"id":   rx.Medication.ID,
"name": rx.Medication.Name,
"unit": rx.Medication.Unit,
}
}
result[i] = m
}
return result
}