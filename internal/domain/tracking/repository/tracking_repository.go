package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ranjbar-dev/nutritrack/internal/domain/tracking/entity"
)

// TrackingRepository defines persistence operations for all 6 tracking types.
type TrackingRepository interface {
	// Food logs
	UpsertFoodLog(ctx context.Context, log *entity.FoodLog) (inserted bool, err error)
	ListFoodLogsByDate(ctx context.Context, clientID uuid.UUID, date time.Time) ([]*entity.FoodLog, error)
	ListFoodLogs(ctx context.Context, clientID uuid.UUID, limit, offset int32) ([]*entity.FoodLog, int64, error)

	// Water logs
	UpsertWaterLog(ctx context.Context, log *entity.WaterLog) (inserted bool, err error)
	ListWaterLogsByDate(ctx context.Context, clientID uuid.UUID, date time.Time) ([]*entity.WaterLog, error)
	ListWaterLogs(ctx context.Context, clientID uuid.UUID, limit, offset int32) ([]*entity.WaterLog, int64, error)

	// Sleep logs
	UpsertSleepLog(ctx context.Context, log *entity.SleepLog) (inserted bool, err error)
	ListSleepLogsByDate(ctx context.Context, clientID uuid.UUID, date time.Time) ([]*entity.SleepLog, error)
	ListSleepLogs(ctx context.Context, clientID uuid.UUID, limit, offset int32) ([]*entity.SleepLog, int64, error)

	// Exercise logs
	UpsertExerciseLog(ctx context.Context, log *entity.ExerciseLog) (inserted bool, err error)
	ListExerciseLogsByDate(ctx context.Context, clientID uuid.UUID, date time.Time) ([]*entity.ExerciseLog, error)
	ListExerciseLogs(ctx context.Context, clientID uuid.UUID, limit, offset int32) ([]*entity.ExerciseLog, int64, error)

	// Medication logs
	UpsertMedicationLog(ctx context.Context, log *entity.MedicationLog) (inserted bool, err error)
	ListMedicationLogsByDate(ctx context.Context, clientID uuid.UUID, date time.Time) ([]*entity.MedicationLog, error)
	ListMedicationLogs(ctx context.Context, clientID uuid.UUID, limit, offset int32) ([]*entity.MedicationLog, int64, error)

	// Body measurements
	UpsertBodyMeasurement(ctx context.Context, m *entity.BodyMeasurement) (inserted bool, err error)
	ListBodyMeasurementsByDate(ctx context.Context, clientID uuid.UUID, date time.Time) ([]*entity.BodyMeasurement, error)
	ListBodyMeasurements(ctx context.Context, clientID uuid.UUID, limit, offset int32) ([]*entity.BodyMeasurement, int64, error)
}
