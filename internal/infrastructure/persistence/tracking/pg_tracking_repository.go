package tracking

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	"github.com/ranjbar-dev/nutritrack/internal/domain/tracking/entity"
	db "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/sqlc"
)

// PgTrackingRepository implements TrackingRepository using PostgreSQL.
type PgTrackingRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

// NewPgTrackingRepository creates a new PgTrackingRepository.
func NewPgTrackingRepository(pool *pgxpool.Pool) *PgTrackingRepository {
	return &PgTrackingRepository{pool: pool, queries: db.New(pool)}
}

// --- Food Logs ---

// UpsertFoodLog inserts or updates a food log entry (idempotent via local_id).
// The entity is populated with DB-assigned fields on success.
func (r *PgTrackingRepository) UpsertFoodLog(ctx context.Context, log *entity.FoodLog) (bool, error) {
	row, inserted, err := r.queries.UpsertFoodLog(ctx, db.UpsertFoodLogParams{
		ClientID:   log.ClientID,
		LocalID:    log.LocalID,
		LoggedAt:   log.LoggedAt,
		LoggedDate: log.LoggedDate,
		FoodID:     log.FoodID,
		FoodName:   log.FoodName,
		Quantity:   float64ToNumeric(log.Quantity),
		Unit:       log.Unit,
		Calories:   float64ToNumeric(log.Calories),
		Protein:    float64ToNumeric(log.Protein),
		Carbs:      float64ToNumeric(log.Carbs),
		Fat:        float64ToNumeric(log.Fat),
		Notes:      log.Notes,
	})
	if err != nil {
		return false, shared.ErrInternal
	}
	*log = *foodLogToDomain(row)
	return inserted, nil
}

// ListFoodLogsByDate returns all food logs for a client on a given date.
func (r *PgTrackingRepository) ListFoodLogsByDate(ctx context.Context, clientID uuid.UUID, date time.Time) ([]*entity.FoodLog, error) {
	rows, err := r.queries.ListFoodLogsByClientAndDate(ctx, db.ListFoodLogsByClientAndDateParams{
		ClientID:   clientID,
		LoggedDate: date,
	})
	if err != nil {
		return nil, shared.ErrInternal
	}
	result := make([]*entity.FoodLog, len(rows))
	for i, row := range rows {
		result[i] = foodLogToDomain(row)
	}
	return result, nil
}

// ListFoodLogs returns paginated food logs for a client along with total count.
func (r *PgTrackingRepository) ListFoodLogs(ctx context.Context, clientID uuid.UUID, limit, offset int32) ([]*entity.FoodLog, int64, error) {
	rows, err := r.queries.ListFoodLogsByClient(ctx, db.ListFoodLogsByClientParams{
		ClientID: clientID,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, 0, shared.ErrInternal
	}
	count, err := r.queries.CountFoodLogsByClient(ctx, clientID)
	if err != nil {
		return nil, 0, shared.ErrInternal
	}
	result := make([]*entity.FoodLog, len(rows))
	for i, row := range rows {
		result[i] = foodLogToDomain(row)
	}
	return result, count, nil
}

// --- Water Logs ---

// UpsertWaterLog inserts or updates a water log entry (idempotent via local_id).
func (r *PgTrackingRepository) UpsertWaterLog(ctx context.Context, log *entity.WaterLog) (bool, error) {
	row, inserted, err := r.queries.UpsertWaterLog(ctx, db.UpsertWaterLogParams{
		ClientID:   log.ClientID,
		LocalID:    log.LocalID,
		LoggedAt:   log.LoggedAt,
		LoggedDate: log.LoggedDate,
		AmountMl:   int32(log.AmountMl),
		Notes:      log.Notes,
	})
	if err != nil {
		return false, shared.ErrInternal
	}
	*log = *waterLogToDomain(row)
	return inserted, nil
}

// ListWaterLogsByDate returns all water logs for a client on a given date.
func (r *PgTrackingRepository) ListWaterLogsByDate(ctx context.Context, clientID uuid.UUID, date time.Time) ([]*entity.WaterLog, error) {
	rows, err := r.queries.ListWaterLogsByClientAndDate(ctx, db.ListWaterLogsByClientAndDateParams{
		ClientID:   clientID,
		LoggedDate: date,
	})
	if err != nil {
		return nil, shared.ErrInternal
	}
	result := make([]*entity.WaterLog, len(rows))
	for i, row := range rows {
		result[i] = waterLogToDomain(row)
	}
	return result, nil
}

// ListWaterLogs returns paginated water logs for a client along with total count.
func (r *PgTrackingRepository) ListWaterLogs(ctx context.Context, clientID uuid.UUID, limit, offset int32) ([]*entity.WaterLog, int64, error) {
	rows, err := r.queries.ListWaterLogsByClient(ctx, db.ListWaterLogsByClientParams{
		ClientID: clientID,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, 0, shared.ErrInternal
	}
	count, err := r.queries.CountWaterLogsByClient(ctx, clientID)
	if err != nil {
		return nil, 0, shared.ErrInternal
	}
	result := make([]*entity.WaterLog, len(rows))
	for i, row := range rows {
		result[i] = waterLogToDomain(row)
	}
	return result, count, nil
}

// --- Sleep Logs ---

// UpsertSleepLog inserts or updates a sleep log entry (idempotent via local_id).
func (r *PgTrackingRepository) UpsertSleepLog(ctx context.Context, log *entity.SleepLog) (bool, error) {
	row, inserted, err := r.queries.UpsertSleepLog(ctx, db.UpsertSleepLogParams{
		ClientID:        log.ClientID,
		LocalID:         log.LocalID,
		LoggedDate:      log.LoggedDate,
		SleepStart:      log.SleepStart,
		SleepEnd:        log.SleepEnd,
		DurationMinutes: int32(log.DurationMinutes),
		Quality:         int32(log.Quality),
		Notes:           log.Notes,
	})
	if err != nil {
		return false, shared.ErrInternal
	}
	*log = *sleepLogToDomain(row)
	return inserted, nil
}

// ListSleepLogsByDate returns all sleep logs for a client on a given date.
func (r *PgTrackingRepository) ListSleepLogsByDate(ctx context.Context, clientID uuid.UUID, date time.Time) ([]*entity.SleepLog, error) {
	rows, err := r.queries.ListSleepLogsByClientAndDate(ctx, db.ListSleepLogsByClientAndDateParams{
		ClientID:   clientID,
		LoggedDate: date,
	})
	if err != nil {
		return nil, shared.ErrInternal
	}
	result := make([]*entity.SleepLog, len(rows))
	for i, row := range rows {
		result[i] = sleepLogToDomain(row)
	}
	return result, nil
}

// ListSleepLogs returns paginated sleep logs for a client along with total count.
func (r *PgTrackingRepository) ListSleepLogs(ctx context.Context, clientID uuid.UUID, limit, offset int32) ([]*entity.SleepLog, int64, error) {
	rows, err := r.queries.ListSleepLogsByClient(ctx, db.ListSleepLogsByClientParams{
		ClientID: clientID,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, 0, shared.ErrInternal
	}
	count, err := r.queries.CountSleepLogsByClient(ctx, clientID)
	if err != nil {
		return nil, 0, shared.ErrInternal
	}
	result := make([]*entity.SleepLog, len(rows))
	for i, row := range rows {
		result[i] = sleepLogToDomain(row)
	}
	return result, count, nil
}

// --- Exercise Logs ---

// UpsertExerciseLog inserts or updates an exercise log entry (idempotent via local_id).
func (r *PgTrackingRepository) UpsertExerciseLog(ctx context.Context, log *entity.ExerciseLog) (bool, error) {
	row, inserted, err := r.queries.UpsertExerciseLog(ctx, db.UpsertExerciseLogParams{
		ClientID:        log.ClientID,
		LocalID:         log.LocalID,
		LoggedAt:        log.LoggedAt,
		LoggedDate:      log.LoggedDate,
		ExerciseName:    log.ExerciseName,
		DurationMinutes: int32(log.DurationMinutes),
		CaloriesBurned:  int32(log.CaloriesBurned),
		Notes:           log.Notes,
	})
	if err != nil {
		return false, shared.ErrInternal
	}
	*log = *exerciseLogToDomain(row)
	return inserted, nil
}

// ListExerciseLogsByDate returns all exercise logs for a client on a given date.
func (r *PgTrackingRepository) ListExerciseLogsByDate(ctx context.Context, clientID uuid.UUID, date time.Time) ([]*entity.ExerciseLog, error) {
	rows, err := r.queries.ListExerciseLogsByClientAndDate(ctx, db.ListExerciseLogsByClientAndDateParams{
		ClientID:   clientID,
		LoggedDate: date,
	})
	if err != nil {
		return nil, shared.ErrInternal
	}
	result := make([]*entity.ExerciseLog, len(rows))
	for i, row := range rows {
		result[i] = exerciseLogToDomain(row)
	}
	return result, nil
}

// ListExerciseLogs returns paginated exercise logs for a client along with total count.
func (r *PgTrackingRepository) ListExerciseLogs(ctx context.Context, clientID uuid.UUID, limit, offset int32) ([]*entity.ExerciseLog, int64, error) {
	rows, err := r.queries.ListExerciseLogsByClient(ctx, db.ListExerciseLogsByClientParams{
		ClientID: clientID,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, 0, shared.ErrInternal
	}
	count, err := r.queries.CountExerciseLogsByClient(ctx, clientID)
	if err != nil {
		return nil, 0, shared.ErrInternal
	}
	result := make([]*entity.ExerciseLog, len(rows))
	for i, row := range rows {
		result[i] = exerciseLogToDomain(row)
	}
	return result, count, nil
}

// --- Medication Logs ---

// UpsertMedicationLog inserts or updates a medication log entry (idempotent via local_id).
func (r *PgTrackingRepository) UpsertMedicationLog(ctx context.Context, log *entity.MedicationLog) (bool, error) {
	row, inserted, err := r.queries.UpsertMedicationLog(ctx, db.UpsertMedicationLogParams{
		ClientID:       log.ClientID,
		LocalID:        log.LocalID,
		LoggedAt:       log.LoggedAt,
		LoggedDate:     log.LoggedDate,
		MedicationID:   log.MedicationID,
		MedicationName: log.MedicationName,
		Dosage:         log.Dosage,
		Notes:          log.Notes,
	})
	if err != nil {
		return false, shared.ErrInternal
	}
	*log = *medicationLogToDomain(row)
	return inserted, nil
}

// ListMedicationLogsByDate returns all medication logs for a client on a given date.
func (r *PgTrackingRepository) ListMedicationLogsByDate(ctx context.Context, clientID uuid.UUID, date time.Time) ([]*entity.MedicationLog, error) {
	rows, err := r.queries.ListMedicationLogsByClientAndDate(ctx, db.ListMedicationLogsByClientAndDateParams{
		ClientID:   clientID,
		LoggedDate: date,
	})
	if err != nil {
		return nil, shared.ErrInternal
	}
	result := make([]*entity.MedicationLog, len(rows))
	for i, row := range rows {
		result[i] = medicationLogToDomain(row)
	}
	return result, nil
}

// ListMedicationLogs returns paginated medication logs for a client along with total count.
func (r *PgTrackingRepository) ListMedicationLogs(ctx context.Context, clientID uuid.UUID, limit, offset int32) ([]*entity.MedicationLog, int64, error) {
	rows, err := r.queries.ListMedicationLogsByClient(ctx, db.ListMedicationLogsByClientParams{
		ClientID: clientID,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, 0, shared.ErrInternal
	}
	count, err := r.queries.CountMedicationLogsByClient(ctx, clientID)
	if err != nil {
		return nil, 0, shared.ErrInternal
	}
	result := make([]*entity.MedicationLog, len(rows))
	for i, row := range rows {
		result[i] = medicationLogToDomain(row)
	}
	return result, count, nil
}

// --- Body Measurements ---

// UpsertBodyMeasurement inserts or updates a body measurement entry (idempotent via local_id).
func (r *PgTrackingRepository) UpsertBodyMeasurement(ctx context.Context, m *entity.BodyMeasurement) (bool, error) {
	row, inserted, err := r.queries.UpsertBodyMeasurement(ctx, db.UpsertBodyMeasurementParams{
		ClientID:     m.ClientID,
		LocalID:      m.LocalID,
		MeasuredAt:   m.MeasuredAt,
		MeasuredDate: m.MeasuredDate,
		WeightKg:     optFloat64ToNumeric(m.WeightKg),
		HeightCm:     optFloat64ToNumeric(m.HeightCm),
		WaistCm:      optFloat64ToNumeric(m.WaistCm),
		HipCm:        optFloat64ToNumeric(m.HipCm),
		ChestCm:      optFloat64ToNumeric(m.ChestCm),
		ArmCm:        optFloat64ToNumeric(m.ArmCm),
		Notes:        m.Notes,
	})
	if err != nil {
		return false, shared.ErrInternal
	}
	*m = *bodyMeasurementToDomain(row)
	return inserted, nil
}

// ListBodyMeasurementsByDate returns all body measurements for a client on a given date.
func (r *PgTrackingRepository) ListBodyMeasurementsByDate(ctx context.Context, clientID uuid.UUID, date time.Time) ([]*entity.BodyMeasurement, error) {
	rows, err := r.queries.ListBodyMeasurementsByClientAndDate(ctx, db.ListBodyMeasurementsByClientAndDateParams{
		ClientID:     clientID,
		MeasuredDate: date,
	})
	if err != nil {
		return nil, shared.ErrInternal
	}
	result := make([]*entity.BodyMeasurement, len(rows))
	for i, row := range rows {
		result[i] = bodyMeasurementToDomain(row)
	}
	return result, nil
}

// ListBodyMeasurements returns paginated body measurements for a client along with total count.
func (r *PgTrackingRepository) ListBodyMeasurements(ctx context.Context, clientID uuid.UUID, limit, offset int32) ([]*entity.BodyMeasurement, int64, error) {
	rows, err := r.queries.ListBodyMeasurementsByClient(ctx, db.ListBodyMeasurementsByClientParams{
		ClientID: clientID,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, 0, shared.ErrInternal
	}
	count, err := r.queries.CountBodyMeasurementsByClient(ctx, clientID)
	if err != nil {
		return nil, 0, shared.ErrInternal
	}
	result := make([]*entity.BodyMeasurement, len(rows))
	for i, row := range rows {
		result[i] = bodyMeasurementToDomain(row)
	}
	return result, count, nil
}
