package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	appTracking "github.com/ranjbar-dev/nutritrack/internal/application/tracking"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	"github.com/ranjbar-dev/nutritrack/internal/interfaces/http/dto"
	"github.com/ranjbar-dev/nutritrack/internal/interfaces/http/middleware"
)

// TrackingHandler handles HTTP requests for all tracking types.
type TrackingHandler struct {
	service *appTracking.TrackingService
}

// NewTrackingHandler creates a new TrackingHandler.
func NewTrackingHandler(service *appTracking.TrackingService) *TrackingHandler {
	return &TrackingHandler{service: service}
}

// --- Request structs ---

type logFoodBody struct {
	LocalID  string    `json:"local_id"`
	LoggedAt time.Time `json:"logged_at"`
	FoodID   *string   `json:"food_id"`
	FoodName string    `json:"food_name"`
	Quantity float64   `json:"quantity"`
	Unit     string    `json:"unit"`
	Calories float64   `json:"calories"`
	Protein  float64   `json:"protein"`
	Carbs    float64   `json:"carbs"`
	Fat      float64   `json:"fat"`
	Notes    string    `json:"notes"`
}

type logWaterBody struct {
	LocalID  string    `json:"local_id"`
	LoggedAt time.Time `json:"logged_at"`
	AmountMl int       `json:"amount_ml"`
	Notes    string    `json:"notes"`
}

type logSleepBody struct {
	LocalID         string    `json:"local_id"`
	SleepStart      time.Time `json:"sleep_start"`
	SleepEnd        time.Time `json:"sleep_end"`
	DurationMinutes int       `json:"duration_minutes"`
	Quality         int       `json:"quality"`
	Notes           string    `json:"notes"`
}

type logExerciseBody struct {
	LocalID         string    `json:"local_id"`
	LoggedAt        time.Time `json:"logged_at"`
	ExerciseName    string    `json:"exercise_name"`
	DurationMinutes int       `json:"duration_minutes"`
	CaloriesBurned  int       `json:"calories_burned"`
	Notes           string    `json:"notes"`
}

type logMedicationBody struct {
	LocalID        string    `json:"local_id"`
	LoggedAt       time.Time `json:"logged_at"`
	MedicationID   *string   `json:"medication_id"`
	MedicationName string    `json:"medication_name"`
	Dosage         string    `json:"dosage"`
	Notes          string    `json:"notes"`
}

type logBodyBody struct {
	LocalID    string    `json:"local_id"`
	MeasuredAt time.Time `json:"measured_at"`
	WeightKg   *float64  `json:"weight_kg"`
	HeightCm   *float64  `json:"height_cm"`
	WaistCm    *float64  `json:"waist_cm"`
	HipCm      *float64  `json:"hip_cm"`
	ChestCm    *float64  `json:"chest_cm"`
	ArmCm      *float64  `json:"arm_cm"`
	Notes      string    `json:"notes"`
}

type bulkSyncBody struct {
	Entries []appTracking.SyncEntry `json:"entries"`
}

// --- Handlers ---

// LogFood handles POST /tracking/food
func (h *TrackingHandler) LogFood(c *gin.Context) {
	callerIDVal, _ := c.Get(middleware.AuthUserIDKey)
	clientID := callerIDVal.(uuid.UUID)

	var body logFoodBody
	if err := c.ShouldBindJSON(&body); err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}
	if body.LocalID == "" {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	var foodID *uuid.UUID
	if body.FoodID != nil && *body.FoodID != "" {
		parsed, err := uuid.Parse(*body.FoodID)
		if err != nil {
			dto.Abort(c, shared.ErrValidation)
			return
		}
		foodID = &parsed
	}

	loggedAt := body.LoggedAt
	if loggedAt.IsZero() {
		loggedAt = shared.NowTehran()
	}

	result, svcErr := h.service.LogFood(c.Request.Context(), clientID, appTracking.LogFoodRequest{
		LocalID:  body.LocalID,
		LoggedAt: loggedAt,
		FoodID:   foodID,
		FoodName: body.FoodName,
		Quantity: body.Quantity,
		Unit:     body.Unit,
		Calories: body.Calories,
		Protein:  body.Protein,
		Carbs:    body.Carbs,
		Fat:      body.Fat,
		Notes:    body.Notes,
	})
	if svcErr != nil {
		if appErr, ok := svcErr.(*shared.AppError); ok {
			dto.Abort(c, appErr)
			return
		}
		dto.Abort(c, shared.ErrInternal)
		return
	}

	dto.OK(c, appTracking.MapFoodLog(result))
}

// LogWater handles POST /tracking/water
func (h *TrackingHandler) LogWater(c *gin.Context) {
	callerIDVal, _ := c.Get(middleware.AuthUserIDKey)
	clientID := callerIDVal.(uuid.UUID)

	var body logWaterBody
	if err := c.ShouldBindJSON(&body); err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}
	if body.LocalID == "" {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	loggedAt := body.LoggedAt
	if loggedAt.IsZero() {
		loggedAt = shared.NowTehran()
	}

	result, svcErr := h.service.LogWater(c.Request.Context(), clientID, appTracking.LogWaterRequest{
		LocalID:  body.LocalID,
		LoggedAt: loggedAt,
		AmountMl: body.AmountMl,
		Notes:    body.Notes,
	})
	if svcErr != nil {
		if appErr, ok := svcErr.(*shared.AppError); ok {
			dto.Abort(c, appErr)
			return
		}
		dto.Abort(c, shared.ErrInternal)
		return
	}

	dto.OK(c, appTracking.MapWaterLog(result))
}

// LogSleep handles POST /tracking/sleep
func (h *TrackingHandler) LogSleep(c *gin.Context) {
	callerIDVal, _ := c.Get(middleware.AuthUserIDKey)
	clientID := callerIDVal.(uuid.UUID)

	var body logSleepBody
	if err := c.ShouldBindJSON(&body); err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}
	if body.LocalID == "" || body.SleepStart.IsZero() || body.SleepEnd.IsZero() {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	result, svcErr := h.service.LogSleep(c.Request.Context(), clientID, appTracking.LogSleepRequest{
		LocalID:         body.LocalID,
		SleepStart:      body.SleepStart,
		SleepEnd:        body.SleepEnd,
		DurationMinutes: body.DurationMinutes,
		Quality:         body.Quality,
		Notes:           body.Notes,
	})
	if svcErr != nil {
		if appErr, ok := svcErr.(*shared.AppError); ok {
			dto.Abort(c, appErr)
			return
		}
		dto.Abort(c, shared.ErrInternal)
		return
	}

	dto.OK(c, appTracking.MapSleepLog(result))
}

// LogExercise handles POST /tracking/exercise
func (h *TrackingHandler) LogExercise(c *gin.Context) {
	callerIDVal, _ := c.Get(middleware.AuthUserIDKey)
	clientID := callerIDVal.(uuid.UUID)

	var body logExerciseBody
	if err := c.ShouldBindJSON(&body); err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}
	if body.LocalID == "" {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	loggedAt := body.LoggedAt
	if loggedAt.IsZero() {
		loggedAt = shared.NowTehran()
	}

	result, svcErr := h.service.LogExercise(c.Request.Context(), clientID, appTracking.LogExerciseRequest{
		LocalID:         body.LocalID,
		LoggedAt:        loggedAt,
		ExerciseName:    body.ExerciseName,
		DurationMinutes: body.DurationMinutes,
		CaloriesBurned:  body.CaloriesBurned,
		Notes:           body.Notes,
	})
	if svcErr != nil {
		if appErr, ok := svcErr.(*shared.AppError); ok {
			dto.Abort(c, appErr)
			return
		}
		dto.Abort(c, shared.ErrInternal)
		return
	}

	dto.OK(c, appTracking.MapExerciseLog(result))
}

// LogMedication handles POST /tracking/medication
func (h *TrackingHandler) LogMedication(c *gin.Context) {
	callerIDVal, _ := c.Get(middleware.AuthUserIDKey)
	clientID := callerIDVal.(uuid.UUID)

	var body logMedicationBody
	if err := c.ShouldBindJSON(&body); err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}
	if body.LocalID == "" {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	loggedAt := body.LoggedAt
	if loggedAt.IsZero() {
		loggedAt = shared.NowTehran()
	}

	var medicationID *uuid.UUID
	if body.MedicationID != nil && *body.MedicationID != "" {
		parsed, err := uuid.Parse(*body.MedicationID)
		if err != nil {
			dto.Abort(c, shared.ErrValidation)
			return
		}
		medicationID = &parsed
	}

	result, svcErr := h.service.LogMedication(c.Request.Context(), clientID, appTracking.LogMedicationRequest{
		LocalID:        body.LocalID,
		LoggedAt:       loggedAt,
		MedicationID:   medicationID,
		MedicationName: body.MedicationName,
		Dosage:         body.Dosage,
		Notes:          body.Notes,
	})
	if svcErr != nil {
		if appErr, ok := svcErr.(*shared.AppError); ok {
			dto.Abort(c, appErr)
			return
		}
		dto.Abort(c, shared.ErrInternal)
		return
	}

	dto.OK(c, appTracking.MapMedicationLog(result))
}

// LogBody handles POST /tracking/body
func (h *TrackingHandler) LogBody(c *gin.Context) {
	callerIDVal, _ := c.Get(middleware.AuthUserIDKey)
	clientID := callerIDVal.(uuid.UUID)

	var body logBodyBody
	if err := c.ShouldBindJSON(&body); err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}
	if body.LocalID == "" {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	measuredAt := body.MeasuredAt
	if measuredAt.IsZero() {
		measuredAt = shared.NowTehran()
	}

	result, svcErr := h.service.LogBody(c.Request.Context(), clientID, appTracking.LogBodyRequest{
		LocalID:    body.LocalID,
		MeasuredAt: measuredAt,
		WeightKg:   body.WeightKg,
		HeightCm:   body.HeightCm,
		WaistCm:    body.WaistCm,
		HipCm:      body.HipCm,
		ChestCm:    body.ChestCm,
		ArmCm:      body.ArmCm,
		Notes:      body.Notes,
	})
	if svcErr != nil {
		if appErr, ok := svcErr.(*shared.AppError); ok {
			dto.Abort(c, appErr)
			return
		}
		dto.Abort(c, shared.ErrInternal)
		return
	}

	dto.OK(c, appTracking.MapBodyMeasurement(result))
}

// BulkSync handles POST /tracking/sync
func (h *TrackingHandler) BulkSync(c *gin.Context) {
	callerIDVal, _ := c.Get(middleware.AuthUserIDKey)
	clientID := callerIDVal.(uuid.UUID)

	var body bulkSyncBody
	if err := c.ShouldBindJSON(&body); err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	result, svcErr := h.service.BulkSync(c.Request.Context(), clientID, body.Entries)
	if svcErr != nil {
		if appErr, ok := svcErr.(*shared.AppError); ok {
			dto.Abort(c, appErr)
			return
		}
		dto.Abort(c, shared.ErrInternal)
		return
	}

	dto.OK(c, result)
}

// GetTracking handles GET /clients/:id/tracking?type=food&date=2025-06-01
func (h *TrackingHandler) GetTracking(c *gin.Context) {
	clientIDStr := c.Param("id")
	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	callerIDVal, _ := c.Get(middleware.AuthUserIDKey)
	callerRoleVal, _ := c.Get(middleware.AuthUserRoleKey)

	trackType := c.Query("type")
	if trackType == "" {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	var date time.Time
	dateStr := c.Query("date")
	if dateStr == "" {
		date = shared.TodayTehran()
	} else {
		parsed, parseErr := time.Parse("2006-01-02", dateStr)
		if parseErr != nil {
			dto.Abort(c, shared.ErrValidation)
			return
		}
		date = parsed
	}

	result, svcErr := h.service.GetTracking(
		c.Request.Context(),
		clientID,
		callerIDVal.(uuid.UUID),
		callerRoleVal.(string),
		trackType,
		date,
	)
	if svcErr != nil {
		if appErr, ok := svcErr.(*shared.AppError); ok {
			dto.Abort(c, appErr)
			return
		}
		dto.Abort(c, shared.ErrInternal)
		return
	}

	dto.OK(c, result)
}
