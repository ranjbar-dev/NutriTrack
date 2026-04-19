package dto

// ─── Request structs ─────────────────────────────────────────────────────────

// CreateDietPlanRequest payload for POST /api/diet-plans
type CreateDietPlanRequest struct {
	ClientID           string  `json:"client_id"            binding:"required,uuid"`
	StartDate          string  `json:"start_date"           binding:"required"`
	EndDate            string  `json:"end_date"             binding:"required"`
	Notes              *string `json:"notes"                binding:"omitempty,max=2000"`
	DailyWaterTargetMl *int    `json:"daily_water_target_ml" binding:"omitempty,gte=0"`
}

// UpdateDietPlanRequest payload for PATCH /api/diet-plans/:id
type UpdateDietPlanRequest struct {
	StartDate          string  `json:"start_date"           binding:"required"`
	EndDate            string  `json:"end_date"             binding:"required"`
	Notes              *string `json:"notes"                binding:"omitempty,max=2000"`
	DailyWaterTargetMl *int    `json:"daily_water_target_ml" binding:"omitempty,gte=0"`
}

// CreateDayRequest payload for POST /api/diet-plans/:id/days
type CreateDayRequest struct {
	DayNumber int     `json:"day_number" binding:"required,min=1"`
	Label     *string `json:"label"      binding:"omitempty,max=100"`
}

// UpdateDayRequest payload for PUT /api/diet-plans/:id/days/:dayId
type UpdateDayRequest struct {
	Label *string `json:"label" binding:"omitempty,max=100"`
}

// CreateMealRequest payload for POST /api/diet-plans/:id/days/:dayId/meals
type CreateMealRequest struct {
	Title         string  `json:"title"          binding:"required,max=200"`
	ScheduledTime *string `json:"scheduled_time" binding:"omitempty"`
	DisplayOrder  int     `json:"display_order"  binding:"gte=0"`
}

// UpdateMealRequest payload for PUT /api/diet-plans/:id/days/:dayId/meals/:mealId
type UpdateMealRequest struct {
	Title         string  `json:"title"          binding:"required,max=200"`
	ScheduledTime *string `json:"scheduled_time" binding:"omitempty"`
	DisplayOrder  int     `json:"display_order"  binding:"gte=0"`
}

// CreateMealOptionRequest payload for POST .../meals/:mealId/options
type CreateMealOptionRequest struct {
	Label *string `json:"label" binding:"omitempty,max=100"`
}

// CreateMealOptionItemRequest payload for POST .../options/:optionId/items
type CreateMealOptionItemRequest struct {
	FoodID          string  `json:"food_id"          binding:"required,uuid"`
	Quantity        float64 `json:"quantity"         binding:"required,gt=0"`
	MeasurementUnit string  `json:"measurement_unit" binding:"required,oneof=gram kg tablespoon teaspoon cup piece slice palm matchbox bowl ml liter"`
	Notes           *string `json:"notes"            binding:"omitempty,max=500"`
}

// UpdateMealOptionItemRequest payload for PUT .../items/:itemId
type UpdateMealOptionItemRequest struct {
	Quantity        float64 `json:"quantity"         binding:"required,gt=0"`
	MeasurementUnit string  `json:"measurement_unit" binding:"required,oneof=gram kg tablespoon teaspoon cup piece slice palm matchbox bowl ml liter"`
	Notes           *string `json:"notes"            binding:"omitempty,max=500"`
}

// CreateExerciseRequest payload for POST .../days/:dayId/exercises
type CreateExerciseRequest struct {
	ExerciseName         string  `json:"exercise_name"          binding:"required,max=200"`
	DurationMinutes      int     `json:"duration_minutes"       binding:"required,gte=1"`
	Description          *string `json:"description"            binding:"omitempty,max=1000"`
	CaloriesBurnEstimate *int    `json:"calories_burn_estimate" binding:"omitempty,gte=0"`
	DisplayOrder         int     `json:"display_order"          binding:"gte=0"`
}

// UpdateExerciseRequest payload for PUT .../exercises/:exerciseId
type UpdateExerciseRequest struct {
	ExerciseName         string  `json:"exercise_name"          binding:"required,max=200"`
	DurationMinutes      int     `json:"duration_minutes"       binding:"required,gte=1"`
	Description          *string `json:"description"            binding:"omitempty,max=1000"`
	CaloriesBurnEstimate *int    `json:"calories_burn_estimate" binding:"omitempty,gte=0"`
}

// CreateMedicationPrescriptionRequest payload for POST .../medications
type CreateMedicationPrescriptionRequest struct {
	MedicationID string   `json:"medication_id" binding:"required,uuid"`
	Dosage       string   `json:"dosage"        binding:"required,max=100"`
	Frequency    string   `json:"frequency"     binding:"required,max=200"`
	Times        []string `json:"times"         binding:"required"`
	Instructions *string  `json:"instructions"  binding:"omitempty,max=1000"`
	StartDate    *string  `json:"start_date"    binding:"omitempty"`
	EndDate      *string  `json:"end_date"      binding:"omitempty"`
}

// UpdateMedicationPrescriptionRequest payload for PUT .../medications/:medId
type UpdateMedicationPrescriptionRequest struct {
	Dosage       string   `json:"dosage"       binding:"required,max=100"`
	Frequency    string   `json:"frequency"    binding:"required,max=200"`
	Times        []string `json:"times"        binding:"required"`
	Instructions *string  `json:"instructions" binding:"omitempty,max=1000"`
	StartDate    *string  `json:"start_date"   binding:"omitempty"`
	EndDate      *string  `json:"end_date"     binding:"omitempty"`
}

// DietPlanListQueryParams query params for GET /api/clients/:clientId/plans
type DietPlanListQueryParams struct {
	Page  int `form:"page,default=1"   binding:"gte=1"`
	Limit int `form:"limit,default=20" binding:"gte=1,lte=100"`
}

// ─── Response structs ─────────────────────────────────────────────────────────

// FoodEmbedded — full nutritional data embedded in every MealOptionItemResponse (D-15)
// Required for client-side nutrition computation (D-14).
type FoodEmbedded struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	Calories          float64 `json:"calories"`
	ProteinG          float64 `json:"protein_g"`
	CarbsG            float64 `json:"carbs_g"`
	FatG              float64 `json:"fat_g"`
	FiberG            float64 `json:"fiber_g"`
	MeasurementUnit   string  `json:"measurement_unit"`
	MeasurementAmount float64 `json:"measurement_amount"`
}

// MealOptionItemResponse single food item within a meal option
type MealOptionItemResponse struct {
	ID              string       `json:"id"`
	Quantity        float64      `json:"quantity"`
	MeasurementUnit string       `json:"measurement_unit"`
	Notes           *string      `json:"notes,omitempty"`
	Food            FoodEmbedded `json:"food"`
}

// MealOptionResponse one option within a meal (client picks one)
type MealOptionResponse struct {
	ID           string                   `json:"id"`
	OptionNumber int                      `json:"option_number"`
	Label        *string                  `json:"label,omitempty"`
	Items        []MealOptionItemResponse `json:"items"`
}

// MealResponse a meal within a plan day
type MealResponse struct {
	ID            string               `json:"id"`
	Title         string               `json:"title"`
	ScheduledTime *string              `json:"scheduled_time,omitempty"`
	DisplayOrder  int                  `json:"display_order"`
	Options       []MealOptionResponse `json:"options"`
}

// PlanExerciseResponse exercise recommendation for a plan day
type PlanExerciseResponse struct {
	ID                   string  `json:"id"`
	ExerciseName         string  `json:"exercise_name"`
	DurationMinutes      int     `json:"duration_minutes"`
	Description          *string `json:"description,omitempty"`
	CaloriesBurnEstimate *int    `json:"calories_burn_estimate,omitempty"`
	DisplayOrder         int     `json:"display_order"`
}

// PlanMedicationResponse prescribed medication within a diet plan
type PlanMedicationResponse struct {
	ID             string   `json:"id"`
	MedicationID   string   `json:"medication_id"`
	MedicationName string   `json:"medication_name"`
	MedicationForm string   `json:"medication_form"`
	Dosage         string   `json:"dosage"`
	Frequency      string   `json:"frequency"`
	Times          []string `json:"times"`
	Instructions   *string  `json:"instructions,omitempty"`
	StartDate      *string  `json:"start_date,omitempty"`
	EndDate        *string  `json:"end_date,omitempty"`
}

// PlanDayResponse a single day within a diet plan aggregate
type PlanDayResponse struct {
	ID        string                 `json:"id"`
	DayNumber int                    `json:"day_number"`
	Label     *string                `json:"label,omitempty"`
	Meals     []MealResponse         `json:"meals"`
	Exercises []PlanExerciseResponse `json:"exercises"`
}

// DietPlanResponse full aggregate response (used for plan detail / builder view)
type DietPlanResponse struct {
	ID                 string                   `json:"id"`
	ClientID           string                   `json:"client_id"`
	NutritionistID     string                   `json:"nutritionist_id"`
	StartDate          string                   `json:"start_date"`
	EndDate            string                   `json:"end_date"`
	Notes              *string                  `json:"notes,omitempty"`
	DailyWaterTargetMl *int                     `json:"daily_water_target_ml,omitempty"`
	Status             string                   `json:"status"`
	Days               []PlanDayResponse        `json:"days"`
	Medications        []PlanMedicationResponse `json:"medications"`
	CreatedAt          string                   `json:"created_at"`
	UpdatedAt          string                   `json:"updated_at"`
}

// DietPlanSummaryResponse flat plan summary used in plan list (D-32) — no nested Days slice
type DietPlanSummaryResponse struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	DayCount  int64  `json:"day_count"`
	CreatedAt string `json:"created_at"`
}

// DietPlanListResponse paginated list of diet plan summaries
type DietPlanListResponse struct {
	Data    []DietPlanSummaryResponse `json:"data"`
	Total   int64                     `json:"total"`
	Page    int                       `json:"page"`
	Limit   int                       `json:"limit"`
	HasMore bool                      `json:"has_more"`
}
