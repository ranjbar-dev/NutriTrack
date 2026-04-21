package tracking

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	"github.com/ranjbar-dev/nutritrack/internal/domain/tracking/entity"
	trackRepo "github.com/ranjbar-dev/nutritrack/internal/domain/tracking/repository"
	userRepo "github.com/ranjbar-dev/nutritrack/internal/domain/user/repository"
)

// TrackingService provides business logic for all 6 tracking types.
type TrackingService struct {
	trackRepo trackRepo.TrackingRepository
	userRepo  userRepo.UserRepository
}

// NewTrackingService creates a new TrackingService.
func NewTrackingService(trackRepo trackRepo.TrackingRepository, userRepo userRepo.UserRepository) *TrackingService {
	return &TrackingService{trackRepo: trackRepo, userRepo: userRepo}
}

// --- Request types ---

type LogFoodRequest struct {
	LocalID  string
	LoggedAt time.Time
	FoodID   *uuid.UUID
	FoodName string
	Quantity float64
	Unit     string
	Calories float64
	Protein  float64
	Carbs    float64
	Fat      float64
	Notes    string
}

type LogWaterRequest struct {
	LocalID  string
	LoggedAt time.Time
	AmountMl int
	Notes    string
}

type LogSleepRequest struct {
	LocalID         string
	SleepStart      time.Time
	SleepEnd        time.Time
	DurationMinutes int
	Quality         int
	Notes           string
}

type LogExerciseRequest struct {
	LocalID         string
	LoggedAt        time.Time
	ExerciseName    string
	DurationMinutes int
	CaloriesBurned  int
	Notes           string
}

type LogMedicationRequest struct {
	LocalID        string
	LoggedAt       time.Time
	MedicationID   *uuid.UUID
	MedicationName string
	Dosage         string
	Notes          string
}

type LogBodyRequest struct {
	LocalID    string
	MeasuredAt time.Time
	WeightKg   *float64
	HeightCm   *float64
	WaistCm    *float64
	HipCm      *float64
	ChestCm    *float64
	ArmCm      *float64
	Notes      string
}

// SyncEntry is one entry in a bulk sync request (type-polymorphic).
type SyncEntry struct {
	Type    string          `json:"type"`
	LocalID string          `json:"local_id"`
	Data    json.RawMessage `json:"data"`
}

// BulkSyncResult reports the outcome of a bulk sync operation.
type BulkSyncResult struct {
	Received int `json:"received"`
	Inserted int `json:"inserted"`
	Skipped  int `json:"skipped"`
}

// --- Tehran date helper ---
// tehranDate extracts the date-only portion of a time.Time in Asia/Tehran timezone.
func tehranDate(t time.Time) time.Time {
	th := shared.ToTehran(t)
	return time.Date(th.Year(), th.Month(), th.Day(), 0, 0, 0, 0, time.UTC)
}

// --- Access control helper ---
// checkClientAccess verifies that the caller is allowed to access a client's data.
func (s *TrackingService) checkClientAccess(ctx context.Context, clientID, callerID uuid.UUID, callerRole string) error {
	if callerRole == "superadmin" {
		return nil
	}
	if callerRole == "client" {
		if callerID != clientID {
			return shared.ErrForbidden
		}
		return nil
	}
	// Nutritionist: verify the client belongs to them.
	client, err := s.userRepo.FindByID(ctx, clientID)
	if err != nil {
		return err
	}
	if client == nil {
		return shared.ErrUserNotFound
	}
	if !client.BelongsTo(callerID) {
		return shared.ErrForbidden
	}
	return nil
}

// --- Logging methods ---

// LogFood logs a food consumption entry for the authenticated client.
func (s *TrackingService) LogFood(ctx context.Context, clientID uuid.UUID, req LogFoodRequest) (*entity.FoodLog, error) {
	log := &entity.FoodLog{
		ClientID:   clientID,
		LocalID:    req.LocalID,
		LoggedAt:   req.LoggedAt,
		LoggedDate: tehranDate(req.LoggedAt),
		FoodID:     req.FoodID,
		FoodName:   req.FoodName,
		Quantity:   req.Quantity,
		Unit:       req.Unit,
		Calories:   req.Calories,
		Protein:    req.Protein,
		Carbs:      req.Carbs,
		Fat:        req.Fat,
		Notes:      req.Notes,
	}
	if _, err := s.trackRepo.UpsertFoodLog(ctx, log); err != nil {
		return nil, err
	}
	return log, nil
}

// LogWater logs a water intake entry.
func (s *TrackingService) LogWater(ctx context.Context, clientID uuid.UUID, req LogWaterRequest) (*entity.WaterLog, error) {
	log := &entity.WaterLog{
		ClientID:   clientID,
		LocalID:    req.LocalID,
		LoggedAt:   req.LoggedAt,
		LoggedDate: tehranDate(req.LoggedAt),
		AmountMl:   req.AmountMl,
		Notes:      req.Notes,
	}
	if _, err := s.trackRepo.UpsertWaterLog(ctx, log); err != nil {
		return nil, err
	}
	return log, nil
}

// LogSleep logs a sleep record. LoggedDate derived from SleepStart.
func (s *TrackingService) LogSleep(ctx context.Context, clientID uuid.UUID, req LogSleepRequest) (*entity.SleepLog, error) {
	log := &entity.SleepLog{
		ClientID:        clientID,
		LocalID:         req.LocalID,
		LoggedDate:      tehranDate(req.SleepStart),
		SleepStart:      req.SleepStart,
		SleepEnd:        req.SleepEnd,
		DurationMinutes: req.DurationMinutes,
		Quality:         req.Quality,
		Notes:           req.Notes,
	}
	if _, err := s.trackRepo.UpsertSleepLog(ctx, log); err != nil {
		return nil, err
	}
	return log, nil
}

// LogExercise logs an exercise session.
func (s *TrackingService) LogExercise(ctx context.Context, clientID uuid.UUID, req LogExerciseRequest) (*entity.ExerciseLog, error) {
	log := &entity.ExerciseLog{
		ClientID:        clientID,
		LocalID:         req.LocalID,
		LoggedAt:        req.LoggedAt,
		LoggedDate:      tehranDate(req.LoggedAt),
		ExerciseName:    req.ExerciseName,
		DurationMinutes: req.DurationMinutes,
		CaloriesBurned:  req.CaloriesBurned,
		Notes:           req.Notes,
	}
	if _, err := s.trackRepo.UpsertExerciseLog(ctx, log); err != nil {
		return nil, err
	}
	return log, nil
}

// LogMedication logs a medication intake.
func (s *TrackingService) LogMedication(ctx context.Context, clientID uuid.UUID, req LogMedicationRequest) (*entity.MedicationLog, error) {
	log := &entity.MedicationLog{
		ClientID:       clientID,
		LocalID:        req.LocalID,
		LoggedAt:       req.LoggedAt,
		LoggedDate:     tehranDate(req.LoggedAt),
		MedicationID:   req.MedicationID,
		MedicationName: req.MedicationName,
		Dosage:         req.Dosage,
		Notes:          req.Notes,
	}
	if _, err := s.trackRepo.UpsertMedicationLog(ctx, log); err != nil {
		return nil, err
	}
	return log, nil
}

// LogBody logs a body measurement.
func (s *TrackingService) LogBody(ctx context.Context, clientID uuid.UUID, req LogBodyRequest) (*entity.BodyMeasurement, error) {
	m := &entity.BodyMeasurement{
		ClientID:     clientID,
		LocalID:      req.LocalID,
		MeasuredAt:   req.MeasuredAt,
		MeasuredDate: tehranDate(req.MeasuredAt),
		WeightKg:     req.WeightKg,
		HeightCm:     req.HeightCm,
		WaistCm:      req.WaistCm,
		HipCm:        req.HipCm,
		ChestCm:      req.ChestCm,
		ArmCm:        req.ArmCm,
		Notes:        req.Notes,
	}
	if _, err := s.trackRepo.UpsertBodyMeasurement(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

// BulkSync processes a batch of mixed tracking entries idempotently.
func (s *TrackingService) BulkSync(ctx context.Context, clientID uuid.UUID, entries []SyncEntry) (BulkSyncResult, error) {
	result := BulkSyncResult{Received: len(entries)}

	for _, entry := range entries {
		var inserted bool
		var err error

		switch entry.Type {
		case "food":
			var req LogFoodRequest
			if jsonErr := json.Unmarshal(entry.Data, &req); jsonErr != nil {
				continue // skip malformed entries
			}
			req.LocalID = entry.LocalID
			log := &entity.FoodLog{
				ClientID:   clientID,
				LocalID:    req.LocalID,
				LoggedAt:   req.LoggedAt,
				LoggedDate: tehranDate(req.LoggedAt),
				FoodID:     req.FoodID,
				FoodName:   req.FoodName,
				Quantity:   req.Quantity,
				Unit:       req.Unit,
				Calories:   req.Calories,
				Protein:    req.Protein,
				Carbs:      req.Carbs,
				Fat:        req.Fat,
				Notes:      req.Notes,
			}
			inserted, err = s.trackRepo.UpsertFoodLog(ctx, log)
		case "water":
			var req LogWaterRequest
			if jsonErr := json.Unmarshal(entry.Data, &req); jsonErr != nil {
				continue
			}
			req.LocalID = entry.LocalID
			log := &entity.WaterLog{
				ClientID:   clientID,
				LocalID:    req.LocalID,
				LoggedAt:   req.LoggedAt,
				LoggedDate: tehranDate(req.LoggedAt),
				AmountMl:   req.AmountMl,
				Notes:      req.Notes,
			}
			inserted, err = s.trackRepo.UpsertWaterLog(ctx, log)
		case "sleep":
			var req LogSleepRequest
			if jsonErr := json.Unmarshal(entry.Data, &req); jsonErr != nil {
				continue
			}
			req.LocalID = entry.LocalID
			log := &entity.SleepLog{
				ClientID:        clientID,
				LocalID:         req.LocalID,
				LoggedDate:      tehranDate(req.SleepStart),
				SleepStart:      req.SleepStart,
				SleepEnd:        req.SleepEnd,
				DurationMinutes: req.DurationMinutes,
				Quality:         req.Quality,
				Notes:           req.Notes,
			}
			inserted, err = s.trackRepo.UpsertSleepLog(ctx, log)
		case "exercise":
			var req LogExerciseRequest
			if jsonErr := json.Unmarshal(entry.Data, &req); jsonErr != nil {
				continue
			}
			req.LocalID = entry.LocalID
			log := &entity.ExerciseLog{
				ClientID:        clientID,
				LocalID:         req.LocalID,
				LoggedAt:        req.LoggedAt,
				LoggedDate:      tehranDate(req.LoggedAt),
				ExerciseName:    req.ExerciseName,
				DurationMinutes: req.DurationMinutes,
				CaloriesBurned:  req.CaloriesBurned,
				Notes:           req.Notes,
			}
			inserted, err = s.trackRepo.UpsertExerciseLog(ctx, log)
		case "medication":
			var req LogMedicationRequest
			if jsonErr := json.Unmarshal(entry.Data, &req); jsonErr != nil {
				continue
			}
			req.LocalID = entry.LocalID
			log := &entity.MedicationLog{
				ClientID:       clientID,
				LocalID:        req.LocalID,
				LoggedAt:       req.LoggedAt,
				LoggedDate:     tehranDate(req.LoggedAt),
				MedicationID:   req.MedicationID,
				MedicationName: req.MedicationName,
				Dosage:         req.Dosage,
				Notes:          req.Notes,
			}
			inserted, err = s.trackRepo.UpsertMedicationLog(ctx, log)
		case "body":
			var req LogBodyRequest
			if jsonErr := json.Unmarshal(entry.Data, &req); jsonErr != nil {
				continue
			}
			req.LocalID = entry.LocalID
			m := &entity.BodyMeasurement{
				ClientID:     clientID,
				LocalID:      req.LocalID,
				MeasuredAt:   req.MeasuredAt,
				MeasuredDate: tehranDate(req.MeasuredAt),
				WeightKg:     req.WeightKg,
				HeightCm:     req.HeightCm,
				WaistCm:      req.WaistCm,
				HipCm:        req.HipCm,
				ChestCm:      req.ChestCm,
				ArmCm:        req.ArmCm,
				Notes:        req.Notes,
			}
			inserted, err = s.trackRepo.UpsertBodyMeasurement(ctx, m)
		default:
			continue // unknown type, skip
		}

		if err != nil {
			continue // skip failed entries in bulk sync
		}
		if inserted {
			result.Inserted++
		} else {
			result.Skipped++
		}
	}

	return result, nil
}

// GetTracking retrieves tracking entries for a client by type and date.
// callerID and callerRole are used for access control.
func (s *TrackingService) GetTracking(ctx context.Context, clientID, callerID uuid.UUID, callerRole, trackType string, date time.Time) (any, error) {
	if err := s.checkClientAccess(ctx, clientID, callerID, callerRole); err != nil {
		return nil, err
	}

	switch trackType {
	case "food":
		return s.trackRepo.ListFoodLogsByDate(ctx, clientID, date)
	case "water":
		return s.trackRepo.ListWaterLogsByDate(ctx, clientID, date)
	case "sleep":
		return s.trackRepo.ListSleepLogsByDate(ctx, clientID, date)
	case "exercise":
		return s.trackRepo.ListExerciseLogsByDate(ctx, clientID, date)
	case "medication":
		return s.trackRepo.ListMedicationLogsByDate(ctx, clientID, date)
	case "body":
		return s.trackRepo.ListBodyMeasurementsByDate(ctx, clientID, date)
	default:
		return nil, shared.ErrValidation
	}
}
