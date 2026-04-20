package dto

// ─── Request structs ─────────────────────────────────────────────────────────

type DailyTrackingQueryParams struct {
	Date string `form:"date" binding:"required"`
}

type TrackingHistoryQueryParams struct {
	From string `form:"from" binding:"required"`
	To   string `form:"to" binding:"required"`
}

type LogFoodRequest struct {
	LocalID          string  `json:"local_id"           binding:"required,uuid"`
	Date             string  `json:"date"               binding:"required"`
	MealID           string  `json:"meal_id"            binding:"required,uuid"`
	SelectedOptionID *string `json:"selected_option_id" binding:"omitempty,uuid"`
	IsSkipped        bool    `json:"is_skipped"`
	Notes            *string `json:"notes"              binding:"omitempty,max=1000"`
}

type LogWaterRequest struct {
	LocalID    string  `json:"local_id"    binding:"required,uuid"`
	Date       string  `json:"date"        binding:"required"`
	AmountMl   int     `json:"amount_ml"   binding:"required,gt=0"`
	LoggedTime *string `json:"logged_time" binding:"omitempty"`
}

type UpsertSleepRequest struct {
	LocalID   string  `json:"local_id"   binding:"required,uuid"`
	Date      string  `json:"date"       binding:"required"`
	SleepTime string  `json:"sleep_time" binding:"required"`
	WakeTime  string  `json:"wake_time"  binding:"required"`
	Quality   string  `json:"quality"    binding:"required,oneof=good fair poor"`
	Notes     *string `json:"notes"      binding:"omitempty,max=1000"`
}

type LogExerciseRequest struct {
	LocalID         string  `json:"local_id"          binding:"required,uuid"`
	Date            string  `json:"date"              binding:"required"`
	ExerciseName    string  `json:"exercise_name"     binding:"required,max=200"`
	DurationMinutes int     `json:"duration_minutes"  binding:"required,gte=1"`
	CaloriesBurned  *int    `json:"calories_burned"   binding:"omitempty,gte=0"`
	Notes           *string `json:"notes"             binding:"omitempty,max=1000"`
}

type LogMedicationRequest struct {
	LocalID                string  `json:"local_id"                 binding:"required,uuid"`
	Date                   string  `json:"date"                     binding:"required"`
	PrescribedMedicationID *string `json:"prescribed_medication_id" binding:"omitempty,uuid"`
	MedicationName         string  `json:"medication_name"          binding:"required,max=200"`
	Dosage                 *string `json:"dosage"                   binding:"omitempty,max=100"`
	TakenAt                string  `json:"taken_at"                 binding:"required"`
	Notes                  *string `json:"notes"                    binding:"omitempty,max=1000"`
	IsSelfReported         bool    `json:"is_self_reported"`
}

type UpsertBodyMeasurementRequest struct {
	LocalID   string   `json:"local_id"   binding:"required,uuid"`
	Date      string   `json:"date"       binding:"required"`
	WeightKg  *float64 `json:"weight_kg"  binding:"omitempty,gt=0"`
	WaistCm   *float64 `json:"waist_cm"   binding:"omitempty,gt=0"`
	HipCm     *float64 `json:"hip_cm"     binding:"omitempty,gt=0"`
	AbdomenCm *float64 `json:"abdomen_cm" binding:"omitempty,gt=0"`
	ThighCm   *float64 `json:"thigh_cm"   binding:"omitempty,gt=0"`
	ChestCm   *float64 `json:"chest_cm"   binding:"omitempty,gt=0"`
	WristCm   *float64 `json:"wrist_cm"   binding:"omitempty,gt=0"`
}

type CreateLabResultRequest struct {
	LocalID          string  `json:"local_id"          form:"local_id"          binding:"required,uuid"`
	Title            string  `json:"title"             form:"title"             binding:"required,max=200"`
	LabType          string  `json:"lab_type"          form:"lab_type"          binding:"required,oneof=blood_test urine_test thyroid hormone allergy other"`
	TestDate         string  `json:"test_date"         form:"test_date"         binding:"required"`
	ExternalLink     *string `json:"external_link"     form:"external_link"     binding:"omitempty,max=2000"`
	FilePath         *string `json:"file_path,omitempty"`
	OriginalFilename *string `json:"original_filename,omitempty"`
	MimeType         *string `json:"mime_type,omitempty"`
	FileSizeBytes    *int64  `json:"file_size_bytes,omitempty"`
}

// ─── Response structs ────────────────────────────────────────────────────────

type FoodLogResponse struct {
	ID               string  `json:"id"`
	ClientID         string  `json:"client_id"`
	LocalID          string  `json:"local_id"`
	Date             string  `json:"date"`
	MealID           string  `json:"meal_id"`
	SelectedOptionID *string `json:"selected_option_id,omitempty"`
	IsSkipped        bool    `json:"is_skipped"`
	Notes            *string `json:"notes,omitempty"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}

type WaterLogResponse struct {
	ID         string  `json:"id"`
	ClientID   string  `json:"client_id"`
	LocalID    string  `json:"local_id"`
	Date       string  `json:"date"`
	AmountMl   int     `json:"amount_ml"`
	LoggedTime *string `json:"logged_time,omitempty"`
	CreatedAt  string  `json:"created_at"`
}

type SleepLogResponse struct {
	ID        string  `json:"id"`
	ClientID  string  `json:"client_id"`
	LocalID   string  `json:"local_id"`
	Date      string  `json:"date"`
	SleepTime string  `json:"sleep_time"`
	WakeTime  string  `json:"wake_time"`
	Quality   string  `json:"quality"`
	Notes     *string `json:"notes,omitempty"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type ExerciseLogResponse struct {
	ID              string  `json:"id"`
	ClientID        string  `json:"client_id"`
	LocalID         string  `json:"local_id"`
	Date            string  `json:"date"`
	ExerciseName    string  `json:"exercise_name"`
	DurationMinutes int     `json:"duration_minutes"`
	CaloriesBurned  *int    `json:"calories_burned,omitempty"`
	Notes           *string `json:"notes,omitempty"`
	CreatedAt       string  `json:"created_at"`
}

type MedicationLogResponse struct {
	ID                     string  `json:"id"`
	ClientID               string  `json:"client_id"`
	LocalID                string  `json:"local_id"`
	Date                   string  `json:"date"`
	PrescribedMedicationID *string `json:"prescribed_medication_id,omitempty"`
	MedicationName         string  `json:"medication_name"`
	Dosage                 *string `json:"dosage,omitempty"`
	TakenAt                string  `json:"taken_at"`
	Notes                  *string `json:"notes,omitempty"`
	IsSelfReported         bool    `json:"is_self_reported"`
	CreatedAt              string  `json:"created_at"`
}

type BodyMeasurementResponse struct {
	ID         string   `json:"id"`
	ClientID   string   `json:"client_id"`
	LocalID    string   `json:"local_id"`
	Date       string   `json:"date"`
	WeightKg   *float64 `json:"weight_kg,omitempty"`
	WaistCm    *float64 `json:"waist_cm,omitempty"`
	HipCm      *float64 `json:"hip_cm,omitempty"`
	AbdomenCm  *float64 `json:"abdomen_cm,omitempty"`
	ThighCm    *float64 `json:"thigh_cm,omitempty"`
	ChestCm    *float64 `json:"chest_cm,omitempty"`
	WristCm    *float64 `json:"wrist_cm,omitempty"`
	RecordedBy string   `json:"recorded_by"`
	CreatedAt  string   `json:"created_at"`
	UpdatedAt  string   `json:"updated_at"`
}

type WeightHistoryPointResponse struct {
	Date     string  `json:"date"`
	WeightKg float64 `json:"weight_kg"`
}

type LabResultResponse struct {
	ID               string  `json:"id"`
	ClientID         string  `json:"client_id"`
	LocalID          string  `json:"local_id"`
	UploadedBy       string  `json:"uploaded_by"`
	Title            string  `json:"title"`
	LabType          string  `json:"lab_type"`
	TestDate         string  `json:"test_date"`
	HasFile          bool    `json:"has_file"`
	FilePath         *string `json:"-"`
	ExternalLink     *string `json:"external_link,omitempty"`
	OriginalFilename *string `json:"original_filename,omitempty"`
	MimeType         *string `json:"mime_type,omitempty"`
	FileSizeBytes    *int64  `json:"file_size_bytes,omitempty"`
	CreatedAt        string  `json:"created_at"`
}

type DailyDashboardResponse struct {
	Date                 string                   `json:"date"`
	WaterTotalMl         int                      `json:"water_total_ml"`
	WaterTargetMl        *int                     `json:"water_target_ml,omitempty"`
	MealsLogged          int                      `json:"meals_logged"`
	MealsTotal           int                      `json:"meals_total"`
	SleepLog             *SleepLogResponse        `json:"sleep_log,omitempty"`
	ExerciseCount        int                      `json:"exercise_count"`
	MedicationTakenCount int                      `json:"medication_taken_count"`
	BodyLoggedToday      bool                     `json:"body_logged_today"`
	TodayBodyMeasurement *BodyMeasurementResponse `json:"today_body_measurement,omitempty"`
	RecentLabResults     []LabResultResponse      `json:"recent_lab_results"`
}
