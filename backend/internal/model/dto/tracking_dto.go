package dto

import "time"

// ─── Food Log ───────────────────────────────────────────────────────────────

// LogFoodRequest is the payload for POST /api/client/tracking/food
type LogFoodRequest struct {
	LocalID          string  `json:"local_id"           binding:"required,uuid"`
	Date             string  `json:"date"               binding:"required"`          // YYYY-MM-DD
	MealID           string  `json:"meal_id"            binding:"required,uuid"`
	SelectedOptionID *string `json:"selected_option_id" binding:"omitempty,uuid"`
	IsSkipped        bool    `json:"is_skipped"`
	Notes            *string `json:"notes"`
}

// FoodLogResponse is returned for food log reads and writes.
type FoodLogResponse struct {
	ID               string    `json:"id"`
	LocalID          string    `json:"local_id"`
	Date             string    `json:"date"`
	MealID           string    `json:"meal_id"`
	SelectedOptionID *string   `json:"selected_option_id"`
	IsSkipped        bool      `json:"is_skipped"`
	Notes            *string   `json:"notes"`
	CreatedAt        time.Time `json:"created_at"`
}

// ─── Water Log ──────────────────────────────────────────────────────────────

// LogWaterRequest is the payload for POST /api/client/tracking/water
type LogWaterRequest struct {
	LocalID    string  `json:"local_id"    binding:"required,uuid"`
	Date       string  `json:"date"        binding:"required"`
	AmountMl   int     `json:"amount_ml"   binding:"required,min=1"`
	LoggedTime *string `json:"logged_time"` // HH:MM or HH:MM:SS, optional
}

// WaterLogResponse is returned for water log reads and writes.
type WaterLogResponse struct {
	ID         string    `json:"id"`
	LocalID    string    `json:"local_id"`
	Date       string    `json:"date"`
	AmountMl   int       `json:"amount_ml"`
	LoggedTime *string   `json:"logged_time"`
	CreatedAt  time.Time `json:"created_at"`
}

// ─── Sleep Log ──────────────────────────────────────────────────────────────

// UpsertSleepRequest is the payload for POST /api/client/tracking/sleep
type UpsertSleepRequest struct {
	LocalID   string  `json:"local_id"   binding:"required,uuid"`
	Date      string  `json:"date"       binding:"required"`
	SleepTime string  `json:"sleep_time" binding:"required"` // HH:MM
	WakeTime  string  `json:"wake_time"  binding:"required"` // HH:MM
	Quality   string  `json:"quality"    binding:"required,oneof=good fair poor"`
	Notes     *string `json:"notes"`
}

// SleepLogResponse is returned for sleep log reads and writes.
type SleepLogResponse struct {
	ID        string    `json:"id"`
	LocalID   string    `json:"local_id"`
	Date      string    `json:"date"`
	SleepTime string    `json:"sleep_time"`
	WakeTime  string    `json:"wake_time"`
	Quality   string    `json:"quality"`
	Notes     *string   `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
}

// ─── Exercise Log ───────────────────────────────────────────────────────────

// LogExerciseRequest is the payload for POST /api/client/tracking/exercise
type LogExerciseRequest struct {
	LocalID         string  `json:"local_id"          binding:"required,uuid"`
	Date            string  `json:"date"              binding:"required"`
	ExerciseName    string  `json:"exercise_name"     binding:"required,max=200"`
	DurationMinutes int     `json:"duration_minutes"  binding:"required,min=1"`
	CaloriesBurned  *int    `json:"calories_burned"   binding:"omitempty,min=0"`
	Notes           *string `json:"notes"`
}

// ExerciseLogResponse is returned for exercise log reads and writes.
type ExerciseLogResponse struct {
	ID              string    `json:"id"`
	LocalID         string    `json:"local_id"`
	Date            string    `json:"date"`
	ExerciseName    string    `json:"exercise_name"`
	DurationMinutes int       `json:"duration_minutes"`
	CaloriesBurned  *int      `json:"calories_burned"`
	Notes           *string   `json:"notes"`
	CreatedAt       time.Time `json:"created_at"`
}

// ─── Medication Log ─────────────────────────────────────────────────────────

// LogMedicationRequest is the payload for POST /api/client/tracking/medication
type LogMedicationRequest struct {
	LocalID                string  `json:"local_id"                binding:"required,uuid"`
	Date                   string  `json:"date"                    binding:"required"`
	PrescribedMedicationID *string `json:"prescribed_medication_id" binding:"omitempty,uuid"`
	MedicationName         string  `json:"medication_name"         binding:"required,max=200"`
	Dosage                 *string `json:"dosage"                  binding:"omitempty,max=100"`
	TakenAt                string  `json:"taken_at"                binding:"required"` // HH:MM
	Notes                  *string `json:"notes"`
	IsSelfReported         bool    `json:"is_self_reported"`
}

// MedicationLogResponse is returned for medication log reads and writes.
type MedicationLogResponse struct {
	ID                     string    `json:"id"`
	LocalID                string    `json:"local_id"`
	Date                   string    `json:"date"`
	PrescribedMedicationID *string   `json:"prescribed_medication_id"`
	MedicationName         string    `json:"medication_name"`
	Dosage                 *string   `json:"dosage"`
	TakenAt                string    `json:"taken_at"`
	Notes                  *string   `json:"notes"`
	IsSelfReported         bool      `json:"is_self_reported"`
	CreatedAt              time.Time `json:"created_at"`
}

// ─── Body Measurement ───────────────────────────────────────────────────────

// UpsertBodyMeasurementRequest is the payload for POST /api/client/tracking/body
type UpsertBodyMeasurementRequest struct {
	LocalID   string   `json:"local_id" binding:"required,uuid"`
	Date      string   `json:"date"     binding:"required"`
	WeightKg  *float64 `json:"weight_kg"  binding:"omitempty,gt=0"`
	WaistCm   *float64 `json:"waist_cm"   binding:"omitempty,gt=0"`
	HipCm     *float64 `json:"hip_cm"     binding:"omitempty,gt=0"`
	AbdomenCm *float64 `json:"abdomen_cm" binding:"omitempty,gt=0"`
	ThighCm   *float64 `json:"thigh_cm"   binding:"omitempty,gt=0"`
	ChestCm   *float64 `json:"chest_cm"   binding:"omitempty,gt=0"`
	WristCm   *float64 `json:"wrist_cm"   binding:"omitempty,gt=0"`
}

// BodyMeasurementResponse is returned for body measurement reads and writes.
type BodyMeasurementResponse struct {
	ID         string    `json:"id"`
	LocalID    string    `json:"local_id"`
	Date       string    `json:"date"`
	WeightKg   *float64  `json:"weight_kg"`
	WaistCm    *float64  `json:"waist_cm"`
	HipCm      *float64  `json:"hip_cm"`
	AbdomenCm  *float64  `json:"abdomen_cm"`
	ThighCm    *float64  `json:"thigh_cm"`
	ChestCm    *float64  `json:"chest_cm"`
	WristCm    *float64  `json:"wrist_cm"`
	RecordedBy string    `json:"recorded_by"`
	CreatedAt  time.Time `json:"created_at"`
}

// WeightPoint is a single data point for the weight history chart (oldest-first).
type WeightPoint struct {
	Date     string  `json:"date"`      // YYYY-MM-DD
	WeightKg float64 `json:"weight_kg"`
}

// ─── Lab Result ─────────────────────────────────────────────────────────────

// CreateLabResultRequest is the multipart form payload for POST /api/client/tracking/lab-results
// File is handled via c.FormFile("file"); service validates at least one of file or ExternalLink (LAB-02).
type CreateLabResultRequest struct {
	LocalID      string  `form:"local_id"      binding:"required,uuid"`
	Title        string  `form:"title"         binding:"required,max=200"`
	LabType      string  `form:"lab_type"      binding:"required,oneof=blood_test urine_test thyroid hormone allergy other"`
	TestDate     string  `form:"test_date"     binding:"required"`
	ExternalLink *string `form:"external_link" binding:"omitempty,url"`
	// File field is extracted from multipart via c.FormFile("file") — not a binding field.
}

// LabResultResponse is returned for lab result reads and writes.
type LabResultResponse struct {
	ID               string    `json:"id"`
	Title            string    `json:"title"`
	LabType          string    `json:"lab_type"`
	TestDate         string    `json:"test_date"`
	HasFile          bool      `json:"has_file"`
	ExternalLink     *string   `json:"external_link"`
	OriginalFilename *string   `json:"original_filename"`
	FileSizeBytes    *int64    `json:"file_size_bytes"`
	MimeType         *string   `json:"mime_type"`
	UploadedBy       string    `json:"uploaded_by"`
	CreatedAt        time.Time `json:"created_at"`
}

// ─── Daily Dashboard ────────────────────────────────────────────────────────

// DailyDashboardResponse is the aggregate summary for GET /api/client/tracking/dashboard
type DailyDashboardResponse struct {
	Date                 string  `json:"date"`
	WaterTotalMl         int64   `json:"water_total_ml"`
	WaterTargetMl        int     `json:"water_target_ml"`       // from active plan
	SleepTime            *string `json:"sleep_time"`            // HH:MM or nil
	WakeTime             *string `json:"wake_time"`             // HH:MM or nil
	SleepQuality         *string `json:"sleep_quality"`         // "good"/"fair"/"poor" or nil
	ExerciseCount        int     `json:"exercise_count"`
	MedicationTakenCount int     `json:"medication_taken_count"`
	BodyLoggedToday      bool    `json:"body_logged_today"`
	MealsLoggedCount     int     `json:"meals_logged_count"`
	TotalMealsCount      int     `json:"total_meals_count"`
}
