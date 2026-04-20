package repository

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ranjbar-dev/nutritrack/backend/internal/model/dto"
	"github.com/ranjbar-dev/nutritrack/backend/internal/repository/sqlc"
)

const sqlGetDailyDashboard = `
WITH target_plan AS (
    SELECT id, start_date, COALESCE(daily_water_target_ml, 0) AS water_target_ml
    FROM diet_plans
    WHERE client_id = $1
      AND status = 'active'
    ORDER BY start_date DESC
    LIMIT 1
),
day_ctx AS (
    SELECT
        tp.id AS plan_id,
        GREATEST(1, ($2::date - tp.start_date::date + 1))::integer AS day_number,
        tp.water_target_ml
    FROM target_plan tp
)
SELECT
    COALESCE((
        SELECT SUM(amount_ml)::bigint
        FROM water_logs
        WHERE client_id = $1
          AND date = $2
    ), 0) AS water_total_ml,
    COALESCE((SELECT water_target_ml FROM day_ctx), 0)::integer AS water_target_ml,
    COALESCE((
        SELECT COUNT(*)::integer
        FROM food_logs
        WHERE client_id = $1
          AND date = $2
    ), 0) AS meals_logged,
    COALESCE((
        SELECT COUNT(*)::integer
        FROM meals m
        JOIN plan_days pd ON pd.id = m.day_id
        JOIN day_ctx dc ON dc.plan_id = pd.plan_id AND pd.day_number = dc.day_number
    ), 0) AS meals_total,
    (SELECT sleep_time FROM sleep_logs WHERE client_id = $1 AND date = $2 ORDER BY updated_at DESC LIMIT 1) AS sleep_time,
    (SELECT wake_time FROM sleep_logs WHERE client_id = $1 AND date = $2 ORDER BY updated_at DESC LIMIT 1) AS wake_time,
    (SELECT quality FROM sleep_logs WHERE client_id = $1 AND date = $2 ORDER BY updated_at DESC LIMIT 1) AS sleep_quality,
    COALESCE((
        SELECT COUNT(*)::integer
        FROM exercise_logs
        WHERE client_id = $1
          AND date = $2
    ), 0) AS exercise_count,
    COALESCE((
        SELECT COUNT(*)::integer
        FROM medication_logs
        WHERE client_id = $1
          AND date = $2
    ), 0) AS medication_taken_count,
    EXISTS(
        SELECT 1
        FROM body_measurements
        WHERE client_id = $1
          AND date = $2
    ) AS body_logged_today
`

type TrackingRepository interface {
	LogFood(ctx context.Context, clientID uuid.UUID, req dto.LogFoodRequest) (*dto.FoodLogResponse, error)
	ListFoodLogs(ctx context.Context, clientID uuid.UUID, date string) ([]dto.FoodLogResponse, error)
	ListFoodLogsForNutritionist(ctx context.Context, clientID, nutritionistID uuid.UUID, from, to string) ([]dto.FoodLogResponse, error)

	LogWater(ctx context.Context, clientID uuid.UUID, req dto.LogWaterRequest) (*dto.WaterLogResponse, error)
	ListWaterLogs(ctx context.Context, clientID uuid.UUID, date string) ([]dto.WaterLogResponse, error)
	ListWaterLogsForNutritionist(ctx context.Context, clientID, nutritionistID uuid.UUID, from, to string) ([]dto.WaterLogResponse, error)

	UpsertSleep(ctx context.Context, clientID uuid.UUID, req dto.UpsertSleepRequest) (*dto.SleepLogResponse, error)
	GetSleepLog(ctx context.Context, clientID uuid.UUID, date string) (*dto.SleepLogResponse, error)
	ListSleepLogsForNutritionist(ctx context.Context, clientID, nutritionistID uuid.UUID, from, to string) ([]dto.SleepLogResponse, error)

	LogExercise(ctx context.Context, clientID uuid.UUID, req dto.LogExerciseRequest) (*dto.ExerciseLogResponse, error)
	ListExerciseLogs(ctx context.Context, clientID uuid.UUID, date string) ([]dto.ExerciseLogResponse, error)
	ListExerciseLogsForNutritionist(ctx context.Context, clientID, nutritionistID uuid.UUID, from, to string) ([]dto.ExerciseLogResponse, error)

	LogMedication(ctx context.Context, clientID uuid.UUID, req dto.LogMedicationRequest) (*dto.MedicationLogResponse, error)
	ListMedicationLogs(ctx context.Context, clientID uuid.UUID, date string) ([]dto.MedicationLogResponse, error)
	ListMedicationLogsForNutritionist(ctx context.Context, clientID, nutritionistID uuid.UUID, from, to string) ([]dto.MedicationLogResponse, error)

	UpsertBodyMeasurement(ctx context.Context, clientID, recordedBy uuid.UUID, req dto.UpsertBodyMeasurementRequest) (*dto.BodyMeasurementResponse, error)
	GetBodyMeasurement(ctx context.Context, clientID uuid.UUID, date string) (*dto.BodyMeasurementResponse, error)
	ListBodyMeasurements(ctx context.Context, clientID uuid.UUID, from, to string) ([]dto.BodyMeasurementResponse, error)
	ListBodyMeasurementsForNutritionist(ctx context.Context, clientID, nutritionistID uuid.UUID, from, to string) ([]dto.BodyMeasurementResponse, error)
	GetWeightHistory(ctx context.Context, clientID uuid.UUID, from, to string) ([]dto.WeightHistoryPointResponse, error)
	GetWeightHistoryForNutritionist(ctx context.Context, clientID, nutritionistID uuid.UUID, from, to string) ([]dto.WeightHistoryPointResponse, error)

	CreateLabResult(ctx context.Context, p CreateLabResultParams) (*dto.LabResultResponse, error)
	ListLabResults(ctx context.Context, clientID uuid.UUID) ([]dto.LabResultResponse, error)
	ListLabResultsForNutritionist(ctx context.Context, clientID, nutritionistID uuid.UUID) ([]dto.LabResultResponse, error)
	GetLabResultForNutritionist(ctx context.Context, labID, clientID, nutritionistID uuid.UUID) (*dto.LabResultResponse, error)

	GetDailyDashboard(ctx context.Context, clientID uuid.UUID, date string) (*dto.DailyDashboardResponse, error)
}

type CreateLabResultParams struct {
	ClientID         uuid.UUID
	LocalID          uuid.UUID
	UploadedBy       uuid.UUID
	Title            string
	LabType          string
	TestDate         string
	FilePath         *string
	ExternalLink     *string
	OriginalFilename *string
	MimeType         *string
	FileSizeBytes    *int64
}

type trackingRepository struct {
	pool *pgxpool.Pool
	q    *sqlc.Queries
}

func NewTrackingRepository(pool *pgxpool.Pool) TrackingRepository {
	return &trackingRepository{pool: pool, q: sqlc.New(pool)}
}

func (r *trackingRepository) LogFood(ctx context.Context, clientID uuid.UUID, req dto.LogFoodRequest) (*dto.FoodLogResponse, error) {
	date, err := parseDateValue(req.Date)
	if err != nil {
		return nil, err
	}
	localID, err := uuid.Parse(req.LocalID)
	if err != nil {
		return nil, err
	}
	mealID, err := uuid.Parse(req.MealID)
	if err != nil {
		return nil, err
	}
	row, err := r.q.UpsertFoodLog(ctx, sqlc.UpsertFoodLogParams{
		ClientID:         uuidParam(clientID),
		LocalID:          uuidParam(localID),
		Date:             date,
		MealID:           uuidParam(mealID),
		SelectedOptionID: optionalUUIDParam(req.SelectedOptionID),
		IsSkipped:        req.IsSkipped,
		Notes:            optionalTextParam(req.Notes),
	})
	if err != nil {
		return nil, err
	}
	resp := toFoodLogResponse(row)
	return &resp, nil
}

func (r *trackingRepository) ListFoodLogs(ctx context.Context, clientID uuid.UUID, date string) ([]dto.FoodLogResponse, error) {
	parsedDate, err := parseDateValue(date)
	if err != nil {
		return nil, err
	}
	rows, err := r.q.ListFoodLogsByDate(ctx, sqlc.ListFoodLogsByDateParams{ClientID: uuidParam(clientID), Date: parsedDate})
	if err != nil {
		return nil, err
	}
	return mapFoodLogs(rows), nil
}

func (r *trackingRepository) ListFoodLogsForNutritionist(ctx context.Context, clientID, nutritionistID uuid.UUID, from, to string) ([]dto.FoodLogResponse, error) {
	fromDate, toDate, err := parseDateRange(from, to)
	if err != nil {
		return nil, err
	}
	rows, err := r.q.ListFoodLogsForNutritionist(ctx, sqlc.ListFoodLogsForNutritionistParams{
		NutritionistID: uuidParam(nutritionistID),
		ClientID:       uuidParam(clientID),
		Date:           fromDate,
		Date_2:         toDate,
	})
	if err != nil {
		return nil, err
	}
	return mapFoodLogs(rows), nil
}

func (r *trackingRepository) LogWater(ctx context.Context, clientID uuid.UUID, req dto.LogWaterRequest) (*dto.WaterLogResponse, error) {
	localID, err := uuid.Parse(req.LocalID)
	if err != nil {
		return nil, err
	}
	date, err := parseDateValue(req.Date)
	if err != nil {
		return nil, err
	}
	loggedTime, err := parseTimeValue(req.LoggedTime)
	if err != nil {
		return nil, err
	}
	row, err := r.q.CreateWaterLog(ctx, sqlc.CreateWaterLogParams{
		ClientID:   uuidParam(clientID),
		LocalID:    uuidParam(localID),
		Date:       date,
		AmountMl:   int32(req.AmountMl),
		LoggedTime: loggedTime,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		existing, getErr := r.q.GetWaterLogByLocalID(ctx, sqlc.GetWaterLogByLocalIDParams{LocalID: uuidParam(localID), ClientID: uuidParam(clientID)})
		if getErr != nil {
			return nil, getErr
		}
		resp := toWaterLogResponse(existing)
		return &resp, nil
	}
	if err != nil {
		return nil, err
	}
	resp := toWaterLogResponse(row)
	return &resp, nil
}

func (r *trackingRepository) ListWaterLogs(ctx context.Context, clientID uuid.UUID, date string) ([]dto.WaterLogResponse, error) {
	parsedDate, err := parseDateValue(date)
	if err != nil {
		return nil, err
	}
	rows, err := r.q.ListWaterLogsByDate(ctx, sqlc.ListWaterLogsByDateParams{ClientID: uuidParam(clientID), Date: parsedDate})
	if err != nil {
		return nil, err
	}
	return mapWaterLogs(rows), nil
}

func (r *trackingRepository) ListWaterLogsForNutritionist(ctx context.Context, clientID, nutritionistID uuid.UUID, from, to string) ([]dto.WaterLogResponse, error) {
	fromDate, toDate, err := parseDateRange(from, to)
	if err != nil {
		return nil, err
	}
	rows, err := r.q.ListWaterLogsForNutritionist(ctx, sqlc.ListWaterLogsForNutritionistParams{
		NutritionistID: uuidParam(nutritionistID),
		ClientID:       uuidParam(clientID),
		Date:           fromDate,
		Date_2:         toDate,
	})
	if err != nil {
		return nil, err
	}
	return mapWaterLogs(rows), nil
}

func (r *trackingRepository) UpsertSleep(ctx context.Context, clientID uuid.UUID, req dto.UpsertSleepRequest) (*dto.SleepLogResponse, error) {
	localID, err := uuid.Parse(req.LocalID)
	if err != nil {
		return nil, err
	}
	date, err := parseDateValue(req.Date)
	if err != nil {
		return nil, err
	}
	sleepTime, err := parseTimeValue(&req.SleepTime)
	if err != nil {
		return nil, err
	}
	wakeTime, err := parseTimeValue(&req.WakeTime)
	if err != nil {
		return nil, err
	}
	row, err := r.q.UpsertSleepLog(ctx, sqlc.UpsertSleepLogParams{
		ClientID:  uuidParam(clientID),
		LocalID:   uuidParam(localID),
		Date:      date,
		SleepTime: sleepTime,
		WakeTime:  wakeTime,
		Quality:   sqlc.SleepQuality(req.Quality),
		Notes:     optionalTextParam(req.Notes),
	})
	if err != nil {
		return nil, err
	}
	resp := toSleepLogResponse(row)
	return &resp, nil
}

func (r *trackingRepository) GetSleepLog(ctx context.Context, clientID uuid.UUID, date string) (*dto.SleepLogResponse, error) {
	parsedDate, err := parseDateValue(date)
	if err != nil {
		return nil, err
	}
	row, err := r.q.GetSleepLogByDate(ctx, sqlc.GetSleepLogByDateParams{ClientID: uuidParam(clientID), Date: parsedDate})
	if err != nil {
		return nil, err
	}
	resp := toSleepLogResponse(row)
	return &resp, nil
}

func (r *trackingRepository) ListSleepLogsForNutritionist(ctx context.Context, clientID, nutritionistID uuid.UUID, from, to string) ([]dto.SleepLogResponse, error) {
	fromDate, toDate, err := parseDateRange(from, to)
	if err != nil {
		return nil, err
	}
	rows, err := r.q.ListSleepLogsForNutritionist(ctx, sqlc.ListSleepLogsForNutritionistParams{
		NutritionistID: uuidParam(nutritionistID),
		ClientID:       uuidParam(clientID),
		Date:           fromDate,
		Date_2:         toDate,
	})
	if err != nil {
		return nil, err
	}
	items := make([]dto.SleepLogResponse, 0, len(rows))
	for _, row := range rows {
		items = append(items, toSleepLogResponse(row))
	}
	return items, nil
}

func (r *trackingRepository) LogExercise(ctx context.Context, clientID uuid.UUID, req dto.LogExerciseRequest) (*dto.ExerciseLogResponse, error) {
	localID, err := uuid.Parse(req.LocalID)
	if err != nil {
		return nil, err
	}
	date, err := parseDateValue(req.Date)
	if err != nil {
		return nil, err
	}
	row, err := r.q.CreateExerciseLog(ctx, sqlc.CreateExerciseLogParams{
		ClientID:        uuidParam(clientID),
		LocalID:         uuidParam(localID),
		Date:            date,
		ExerciseName:    req.ExerciseName,
		DurationMinutes: int32(req.DurationMinutes),
		CaloriesBurned:  optionalInt4Param(req.CaloriesBurned),
		Notes:           optionalTextParam(req.Notes),
	})
	if errors.Is(err, pgx.ErrNoRows) {
		existing, getErr := r.q.GetExerciseLogByLocalID(ctx, sqlc.GetExerciseLogByLocalIDParams{LocalID: uuidParam(localID), ClientID: uuidParam(clientID)})
		if getErr != nil {
			return nil, getErr
		}
		resp := toExerciseLogResponse(existing)
		return &resp, nil
	}
	if err != nil {
		return nil, err
	}
	resp := toExerciseLogResponse(row)
	return &resp, nil
}

func (r *trackingRepository) ListExerciseLogs(ctx context.Context, clientID uuid.UUID, date string) ([]dto.ExerciseLogResponse, error) {
	parsedDate, err := parseDateValue(date)
	if err != nil {
		return nil, err
	}
	rows, err := r.q.ListExerciseLogsByDate(ctx, sqlc.ListExerciseLogsByDateParams{ClientID: uuidParam(clientID), Date: parsedDate})
	if err != nil {
		return nil, err
	}
	items := make([]dto.ExerciseLogResponse, 0, len(rows))
	for _, row := range rows {
		items = append(items, toExerciseLogResponse(row))
	}
	return items, nil
}

func (r *trackingRepository) ListExerciseLogsForNutritionist(ctx context.Context, clientID, nutritionistID uuid.UUID, from, to string) ([]dto.ExerciseLogResponse, error) {
	fromDate, toDate, err := parseDateRange(from, to)
	if err != nil {
		return nil, err
	}
	rows, err := r.q.ListExerciseLogsForNutritionist(ctx, sqlc.ListExerciseLogsForNutritionistParams{
		NutritionistID: uuidParam(nutritionistID),
		ClientID:       uuidParam(clientID),
		Date:           fromDate,
		Date_2:         toDate,
	})
	if err != nil {
		return nil, err
	}
	items := make([]dto.ExerciseLogResponse, 0, len(rows))
	for _, row := range rows {
		items = append(items, toExerciseLogResponse(row))
	}
	return items, nil
}

func (r *trackingRepository) LogMedication(ctx context.Context, clientID uuid.UUID, req dto.LogMedicationRequest) (*dto.MedicationLogResponse, error) {
	localID, err := uuid.Parse(req.LocalID)
	if err != nil {
		return nil, err
	}
	date, err := parseDateValue(req.Date)
	if err != nil {
		return nil, err
	}
	takenAt, err := parseTimeValue(&req.TakenAt)
	if err != nil {
		return nil, err
	}
	row, err := r.q.CreateMedicationLog(ctx, sqlc.CreateMedicationLogParams{
		ClientID:               uuidParam(clientID),
		LocalID:                uuidParam(localID),
		Date:                   date,
		PrescribedMedicationID: optionalUUIDParam(req.PrescribedMedicationID),
		MedicationName:         req.MedicationName,
		Dosage:                 optionalTextParam(req.Dosage),
		TakenAt:                takenAt,
		Notes:                  optionalTextParam(req.Notes),
		IsSelfReported:         req.IsSelfReported,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		existing, getErr := r.q.GetMedicationLogByLocalID(ctx, sqlc.GetMedicationLogByLocalIDParams{LocalID: uuidParam(localID), ClientID: uuidParam(clientID)})
		if getErr != nil {
			return nil, getErr
		}
		resp := toMedicationLogResponse(existing)
		return &resp, nil
	}
	if err != nil {
		return nil, err
	}
	resp := toMedicationLogResponse(row)
	return &resp, nil
}

func (r *trackingRepository) ListMedicationLogs(ctx context.Context, clientID uuid.UUID, date string) ([]dto.MedicationLogResponse, error) {
	parsedDate, err := parseDateValue(date)
	if err != nil {
		return nil, err
	}
	rows, err := r.q.ListMedicationLogsByDate(ctx, sqlc.ListMedicationLogsByDateParams{ClientID: uuidParam(clientID), Date: parsedDate})
	if err != nil {
		return nil, err
	}
	items := make([]dto.MedicationLogResponse, 0, len(rows))
	for _, row := range rows {
		items = append(items, toMedicationLogResponse(row))
	}
	return items, nil
}

func (r *trackingRepository) ListMedicationLogsForNutritionist(ctx context.Context, clientID, nutritionistID uuid.UUID, from, to string) ([]dto.MedicationLogResponse, error) {
	fromDate, toDate, err := parseDateRange(from, to)
	if err != nil {
		return nil, err
	}
	rows, err := r.q.ListMedicationLogsForNutritionist(ctx, sqlc.ListMedicationLogsForNutritionistParams{
		NutritionistID: uuidParam(nutritionistID),
		ClientID:       uuidParam(clientID),
		Date:           fromDate,
		Date_2:         toDate,
	})
	if err != nil {
		return nil, err
	}
	items := make([]dto.MedicationLogResponse, 0, len(rows))
	for _, row := range rows {
		items = append(items, toMedicationLogResponse(row))
	}
	return items, nil
}

func (r *trackingRepository) UpsertBodyMeasurement(ctx context.Context, clientID, recordedBy uuid.UUID, req dto.UpsertBodyMeasurementRequest) (*dto.BodyMeasurementResponse, error) {
	localID, err := uuid.Parse(req.LocalID)
	if err != nil {
		return nil, err
	}
	date, err := parseDateValue(req.Date)
	if err != nil {
		return nil, err
	}
	row, err := r.q.UpsertBodyMeasurement(ctx, sqlc.UpsertBodyMeasurementParams{
		ClientID:   uuidParam(clientID),
		LocalID:    uuidParam(localID),
		Date:       date,
		WeightKg:   optionalNumericParam(req.WeightKg),
		WaistCm:    optionalNumericParam(req.WaistCm),
		HipCm:      optionalNumericParam(req.HipCm),
		AbdomenCm:  optionalNumericParam(req.AbdomenCm),
		ThighCm:    optionalNumericParam(req.ThighCm),
		ChestCm:    optionalNumericParam(req.ChestCm),
		WristCm:    optionalNumericParam(req.WristCm),
		RecordedBy: uuidParam(recordedBy),
	})
	if err != nil {
		return nil, err
	}
	resp, err := toBodyMeasurementResponse(row)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (r *trackingRepository) GetBodyMeasurement(ctx context.Context, clientID uuid.UUID, date string) (*dto.BodyMeasurementResponse, error) {
	parsedDate, err := parseDateValue(date)
	if err != nil {
		return nil, err
	}
	row, err := r.q.GetBodyMeasurementByDate(ctx, sqlc.GetBodyMeasurementByDateParams{ClientID: uuidParam(clientID), Date: parsedDate})
	if err != nil {
		return nil, err
	}
	resp, err := toBodyMeasurementResponse(row)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (r *trackingRepository) ListBodyMeasurements(ctx context.Context, clientID uuid.UUID, from, to string) ([]dto.BodyMeasurementResponse, error) {
	fromDate, toDate, err := parseDateRange(from, to)
	if err != nil {
		return nil, err
	}
	rows, err := r.q.ListBodyMeasurementsByDateRange(ctx, sqlc.ListBodyMeasurementsByDateRangeParams{ClientID: uuidParam(clientID), Date: fromDate, Date_2: toDate})
	if err != nil {
		return nil, err
	}
	return mapBodyMeasurements(rows)
}

func (r *trackingRepository) ListBodyMeasurementsForNutritionist(ctx context.Context, clientID, nutritionistID uuid.UUID, from, to string) ([]dto.BodyMeasurementResponse, error) {
	fromDate, toDate, err := parseDateRange(from, to)
	if err != nil {
		return nil, err
	}
	rows, err := r.q.ListBodyMeasurementsForNutritionist(ctx, sqlc.ListBodyMeasurementsForNutritionistParams{
		NutritionistID: uuidParam(nutritionistID),
		ClientID:       uuidParam(clientID),
		Date:           fromDate,
		Date_2:         toDate,
	})
	if err != nil {
		return nil, err
	}
	return mapBodyMeasurements(rows)
}

func (r *trackingRepository) GetWeightHistory(ctx context.Context, clientID uuid.UUID, from, to string) ([]dto.WeightHistoryPointResponse, error) {
	fromDate, toDate, err := parseDateRange(from, to)
	if err != nil {
		return nil, err
	}
	rows, err := r.q.ListWeightHistory(ctx, sqlc.ListWeightHistoryParams{ClientID: uuidParam(clientID), Date: fromDate, Date_2: toDate})
	if err != nil {
		return nil, err
	}
	return mapWeightPoints(rows)
}

func (r *trackingRepository) GetWeightHistoryForNutritionist(ctx context.Context, clientID, nutritionistID uuid.UUID, from, to string) ([]dto.WeightHistoryPointResponse, error) {
	fromDate, toDate, err := parseDateRange(from, to)
	if err != nil {
		return nil, err
	}
	rows, err := r.q.ListWeightHistoryForNutritionist(ctx, sqlc.ListWeightHistoryForNutritionistParams{
		NutritionistID: uuidParam(nutritionistID),
		ClientID:       uuidParam(clientID),
		Date:           fromDate,
		Date_2:         toDate,
	})
	if err != nil {
		return nil, err
	}
	return mapWeightPointsForNutritionist(rows)
}

func (r *trackingRepository) CreateLabResult(ctx context.Context, p CreateLabResultParams) (*dto.LabResultResponse, error) {
	testDate, err := parseDateValue(p.TestDate)
	if err != nil {
		return nil, err
	}
	row, err := r.q.CreateLabResult(ctx, sqlc.CreateLabResultParams{
		ClientID:         uuidParam(p.ClientID),
		LocalID:          uuidParam(p.LocalID),
		UploadedBy:       uuidParam(p.UploadedBy),
		Title:            p.Title,
		LabType:          sqlc.LabResultType(p.LabType),
		TestDate:         testDate,
		FilePath:         optionalStringParam(p.FilePath),
		ExternalLink:     optionalStringParam(p.ExternalLink),
		OriginalFilename: optionalStringParam(p.OriginalFilename),
		MimeType:         optionalStringParam(p.MimeType),
		FileSizeBytes:    optionalInt8Param(p.FileSizeBytes),
	})
	if errors.Is(err, pgx.ErrNoRows) {
		existing, getErr := r.q.GetLabResultByLocalID(ctx, sqlc.GetLabResultByLocalIDParams{LocalID: uuidParam(p.LocalID), ClientID: uuidParam(p.ClientID)})
		if getErr != nil {
			return nil, getErr
		}
		resp := toLabResultResponse(existing)
		return &resp, nil
	}
	if err != nil {
		return nil, err
	}
	resp := toLabResultResponse(row)
	return &resp, nil
}

func (r *trackingRepository) ListLabResults(ctx context.Context, clientID uuid.UUID) ([]dto.LabResultResponse, error) {
	rows, err := r.q.ListLabResultsByClient(ctx, uuidParam(clientID))
	if err != nil {
		return nil, err
	}
	items := make([]dto.LabResultResponse, 0, len(rows))
	for _, row := range rows {
		items = append(items, toLabResultResponse(row))
	}
	return items, nil
}

func (r *trackingRepository) ListLabResultsForNutritionist(ctx context.Context, clientID, nutritionistID uuid.UUID) ([]dto.LabResultResponse, error) {
	rows, err := r.q.ListLabResultsForNutritionist(ctx, sqlc.ListLabResultsForNutritionistParams{NutritionistID: uuidParam(nutritionistID), ClientID: uuidParam(clientID)})
	if err != nil {
		return nil, err
	}
	items := make([]dto.LabResultResponse, 0, len(rows))
	for _, row := range rows {
		items = append(items, toLabResultResponse(row))
	}
	return items, nil
}

func (r *trackingRepository) GetLabResultForNutritionist(ctx context.Context, labID, clientID, nutritionistID uuid.UUID) (*dto.LabResultResponse, error) {
	row, err := r.q.GetLabResultForNutritionist(ctx, sqlc.GetLabResultForNutritionistParams{
		NutritionistID: uuidParam(nutritionistID),
		ID:             uuidParam(labID),
		ClientID:       uuidParam(clientID),
	})
	if err != nil {
		return nil, err
	}
	resp := toLabResultResponse(row)
	return &resp, nil
}

func (r *trackingRepository) GetDailyDashboard(ctx context.Context, clientID uuid.UUID, date string) (*dto.DailyDashboardResponse, error) {
	parsedDate, err := parseDateValue(date)
	if err != nil {
		return nil, err
	}

	var (
		waterTotal      int64
		waterTarget     int32
		mealsLogged     int32
		mealsTotal      int32
		sleepTime       pgtype.Time
		wakeTime        pgtype.Time
		sleepQuality    sqlc.NullSleepQuality
		exerciseCount   int32
		medicationCount int32
		bodyLogged      bool
	)

	err = r.pool.QueryRow(ctx, sqlGetDailyDashboard, uuidParam(clientID), parsedDate).Scan(
		&waterTotal,
		&waterTarget,
		&mealsLogged,
		&mealsTotal,
		&sleepTime,
		&wakeTime,
		&sleepQuality,
		&exerciseCount,
		&medicationCount,
		&bodyLogged,
	)
	if err != nil {
		return nil, err
	}

	resp := &dto.DailyDashboardResponse{
		Date:                 formatDate(parsedDate),
		WaterTotalMl:         int(waterTotal),
		WaterTargetMl:        intPtr(int(waterTarget)),
		MealsLogged:          int(mealsLogged),
		MealsTotal:           int(mealsTotal),
		ExerciseCount:        int(exerciseCount),
		MedicationTakenCount: int(medicationCount),
		BodyLoggedToday:      bodyLogged,
		RecentLabResults:     []dto.LabResultResponse{},
	}

	if sleepQuality.Valid {
		resp.SleepLog = &dto.SleepLogResponse{
			Date:      formatDate(parsedDate),
			SleepTime: formatTime(sleepTime),
			WakeTime:  formatTime(wakeTime),
			Quality:   string(sleepQuality.SleepQuality),
		}
	}

	if bodyLogged {
		body, bodyErr := r.GetBodyMeasurement(ctx, clientID, date)
		if bodyErr == nil {
			resp.TodayBodyMeasurement = body
		}
	}

	return resp, nil
}

func mapFoodLogs(rows []sqlc.FoodLog) []dto.FoodLogResponse {
	items := make([]dto.FoodLogResponse, 0, len(rows))
	for _, row := range rows {
		items = append(items, toFoodLogResponse(row))
	}
	return items
}

func mapWaterLogs(rows []sqlc.WaterLog) []dto.WaterLogResponse {
	items := make([]dto.WaterLogResponse, 0, len(rows))
	for _, row := range rows {
		items = append(items, toWaterLogResponse(row))
	}
	return items
}

func mapBodyMeasurements(rows []sqlc.BodyMeasurement) ([]dto.BodyMeasurementResponse, error) {
	items := make([]dto.BodyMeasurementResponse, 0, len(rows))
	for _, row := range rows {
		resp, err := toBodyMeasurementResponse(row)
		if err != nil {
			return nil, err
		}
		items = append(items, resp)
	}
	return items, nil
}

func mapWeightPoints(rows []sqlc.ListWeightHistoryRow) ([]dto.WeightHistoryPointResponse, error) {
	items := make([]dto.WeightHistoryPointResponse, 0, len(rows))
	for _, row := range rows {
		value, err := numericToFloat64(row.WeightKg)
		if err != nil {
			return nil, err
		}
		items = append(items, dto.WeightHistoryPointResponse{Date: formatDate(row.Date), WeightKg: value})
	}
	return items, nil
}

func mapWeightPointsForNutritionist(rows []sqlc.ListWeightHistoryForNutritionistRow) ([]dto.WeightHistoryPointResponse, error) {
	items := make([]dto.WeightHistoryPointResponse, 0, len(rows))
	for _, row := range rows {
		value, err := numericToFloat64(row.WeightKg)
		if err != nil {
			return nil, err
		}
		items = append(items, dto.WeightHistoryPointResponse{Date: formatDate(row.Date), WeightKg: value})
	}
	return items, nil
}

func toFoodLogResponse(row sqlc.FoodLog) dto.FoodLogResponse {
	resp := dto.FoodLogResponse{
		ID:        uuid.UUID(row.ID.Bytes).String(),
		ClientID:  uuid.UUID(row.ClientID.Bytes).String(),
		LocalID:   uuid.UUID(row.LocalID.Bytes).String(),
		Date:      formatDate(row.Date),
		MealID:    uuid.UUID(row.MealID.Bytes).String(),
		IsSkipped: row.IsSkipped,
		CreatedAt: formatTimestamptz(row.CreatedAt),
		UpdatedAt: formatTimestamptz(row.UpdatedAt),
	}
	if row.SelectedOptionID.Valid {
		value := uuid.UUID(row.SelectedOptionID.Bytes).String()
		resp.SelectedOptionID = &value
	}
	if row.Notes.Valid {
		value := row.Notes.String
		resp.Notes = &value
	}
	return resp
}

func toWaterLogResponse(row sqlc.WaterLog) dto.WaterLogResponse {
	resp := dto.WaterLogResponse{
		ID:        uuid.UUID(row.ID.Bytes).String(),
		ClientID:  uuid.UUID(row.ClientID.Bytes).String(),
		LocalID:   uuid.UUID(row.LocalID.Bytes).String(),
		Date:      formatDate(row.Date),
		AmountMl:  int(row.AmountMl),
		CreatedAt: formatTimestamptz(row.CreatedAt),
	}
	if row.LoggedTime.Valid {
		value := formatTime(row.LoggedTime)
		resp.LoggedTime = &value
	}
	return resp
}

func toSleepLogResponse(row sqlc.SleepLog) dto.SleepLogResponse {
	resp := dto.SleepLogResponse{
		ID:        uuid.UUID(row.ID.Bytes).String(),
		ClientID:  uuid.UUID(row.ClientID.Bytes).String(),
		LocalID:   uuid.UUID(row.LocalID.Bytes).String(),
		Date:      formatDate(row.Date),
		SleepTime: formatTime(row.SleepTime),
		WakeTime:  formatTime(row.WakeTime),
		Quality:   string(row.Quality),
		CreatedAt: formatTimestamptz(row.CreatedAt),
		UpdatedAt: formatTimestamptz(row.UpdatedAt),
	}
	if row.Notes.Valid {
		value := row.Notes.String
		resp.Notes = &value
	}
	return resp
}

func toExerciseLogResponse(row sqlc.ExerciseLog) dto.ExerciseLogResponse {
	resp := dto.ExerciseLogResponse{
		ID:              uuid.UUID(row.ID.Bytes).String(),
		ClientID:        uuid.UUID(row.ClientID.Bytes).String(),
		LocalID:         uuid.UUID(row.LocalID.Bytes).String(),
		Date:            formatDate(row.Date),
		ExerciseName:    row.ExerciseName,
		DurationMinutes: int(row.DurationMinutes),
		CreatedAt:       formatTimestamptz(row.CreatedAt),
	}
	if row.CaloriesBurned.Valid {
		value := int(row.CaloriesBurned.Int32)
		resp.CaloriesBurned = &value
	}
	if row.Notes.Valid {
		value := row.Notes.String
		resp.Notes = &value
	}
	return resp
}

func toMedicationLogResponse(row sqlc.MedicationLog) dto.MedicationLogResponse {
	resp := dto.MedicationLogResponse{
		ID:             uuid.UUID(row.ID.Bytes).String(),
		ClientID:       uuid.UUID(row.ClientID.Bytes).String(),
		LocalID:        uuid.UUID(row.LocalID.Bytes).String(),
		Date:           formatDate(row.Date),
		MedicationName: row.MedicationName,
		TakenAt:        formatTime(row.TakenAt),
		IsSelfReported: row.IsSelfReported,
		CreatedAt:      formatTimestamptz(row.CreatedAt),
	}
	if row.PrescribedMedicationID.Valid {
		value := uuid.UUID(row.PrescribedMedicationID.Bytes).String()
		resp.PrescribedMedicationID = &value
	}
	if row.Dosage.Valid {
		value := row.Dosage.String
		resp.Dosage = &value
	}
	if row.Notes.Valid {
		value := row.Notes.String
		resp.Notes = &value
	}
	return resp
}

func toBodyMeasurementResponse(row sqlc.BodyMeasurement) (dto.BodyMeasurementResponse, error) {
	resp := dto.BodyMeasurementResponse{
		ID:         uuid.UUID(row.ID.Bytes).String(),
		ClientID:   uuid.UUID(row.ClientID.Bytes).String(),
		LocalID:    uuid.UUID(row.LocalID.Bytes).String(),
		Date:       formatDate(row.Date),
		RecordedBy: uuid.UUID(row.RecordedBy.Bytes).String(),
		CreatedAt:  formatTimestamptz(row.CreatedAt),
		UpdatedAt:  formatTimestamptz(row.UpdatedAt),
	}
	var err error
	if resp.WeightKg, err = numericToFloat64Ptr(row.WeightKg); err != nil {
		return dto.BodyMeasurementResponse{}, err
	}
	if resp.WaistCm, err = numericToFloat64Ptr(row.WaistCm); err != nil {
		return dto.BodyMeasurementResponse{}, err
	}
	if resp.HipCm, err = numericToFloat64Ptr(row.HipCm); err != nil {
		return dto.BodyMeasurementResponse{}, err
	}
	if resp.AbdomenCm, err = numericToFloat64Ptr(row.AbdomenCm); err != nil {
		return dto.BodyMeasurementResponse{}, err
	}
	if resp.ThighCm, err = numericToFloat64Ptr(row.ThighCm); err != nil {
		return dto.BodyMeasurementResponse{}, err
	}
	if resp.ChestCm, err = numericToFloat64Ptr(row.ChestCm); err != nil {
		return dto.BodyMeasurementResponse{}, err
	}
	if resp.WristCm, err = numericToFloat64Ptr(row.WristCm); err != nil {
		return dto.BodyMeasurementResponse{}, err
	}
	return resp, nil
}

func toLabResultResponse(row sqlc.LabResult) dto.LabResultResponse {
	resp := dto.LabResultResponse{
		ID:         uuid.UUID(row.ID.Bytes).String(),
		ClientID:   uuid.UUID(row.ClientID.Bytes).String(),
		LocalID:    uuid.UUID(row.LocalID.Bytes).String(),
		UploadedBy: uuid.UUID(row.UploadedBy.Bytes).String(),
		Title:      row.Title,
		LabType:    string(row.LabType),
		TestDate:   formatDate(row.TestDate),
		CreatedAt:  formatTimestamptz(row.CreatedAt),
		HasFile:    row.FilePath.Valid,
	}
	if row.FilePath.Valid {
		value := row.FilePath.String
		resp.FilePath = &value
	}
	if row.ExternalLink.Valid {
		value := row.ExternalLink.String
		resp.ExternalLink = &value
	}
	if row.OriginalFilename.Valid {
		value := row.OriginalFilename.String
		resp.OriginalFilename = &value
	}
	if row.MimeType.Valid {
		value := row.MimeType.String
		resp.MimeType = &value
	}
	if row.FileSizeBytes.Valid {
		value := row.FileSizeBytes.Int64
		resp.FileSizeBytes = &value
	}
	return resp
}

func parseDateValue(value string) (pgtype.Date, error) {
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return pgtype.Date{}, fmt.Errorf("parse date %q: %w", value, err)
	}
	return pgtype.Date{Time: parsed, Valid: true}, nil
}

func parseDateRange(from, to string) (pgtype.Date, pgtype.Date, error) {
	fromDate, err := parseDateValue(from)
	if err != nil {
		return pgtype.Date{}, pgtype.Date{}, err
	}
	toDate, err := parseDateValue(to)
	if err != nil {
		return pgtype.Date{}, pgtype.Date{}, err
	}
	return fromDate, toDate, nil
}

func parseTimeValue(value *string) (pgtype.Time, error) {
	if value == nil || *value == "" {
		return pgtype.Time{Valid: false}, nil
	}
	parsed, err := time.Parse("15:04", *value)
	if err != nil {
		return pgtype.Time{}, fmt.Errorf("parse time %q: %w", *value, err)
	}
	micros := int64(parsed.Hour())*3_600_000_000 + int64(parsed.Minute())*60_000_000
	return pgtype.Time{Microseconds: micros, Valid: true}, nil
}

func uuidParam(value uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: value, Valid: true}
}

func optionalUUIDParam(value *string) pgtype.UUID {
	if value == nil || *value == "" {
		return pgtype.UUID{Valid: false}
	}
	parsed, err := uuid.Parse(*value)
	if err != nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: parsed, Valid: true}
}

func optionalTextParam(value *string) pgtype.Text {
	if value == nil || *value == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *value, Valid: true}
}

func optionalStringParam(value *string) pgtype.Text {
	if value == nil || *value == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *value, Valid: true}
}

func optionalInt4Param(value *int) pgtype.Int4 {
	if value == nil {
		return pgtype.Int4{Valid: false}
	}
	return pgtype.Int4{Int32: int32(*value), Valid: true}
}

func optionalInt8Param(value *int64) pgtype.Int8 {
	if value == nil {
		return pgtype.Int8{Valid: false}
	}
	return pgtype.Int8{Int64: *value, Valid: true}
}

func optionalNumericParam(value *float64) pgtype.Numeric {
	if value == nil {
		return pgtype.Numeric{Valid: false}
	}
	var numeric pgtype.Numeric
	if err := numeric.Scan(strconv.FormatFloat(*value, 'f', -1, 64)); err != nil {
		return pgtype.Numeric{Valid: false}
	}
	return numeric
}

func numericToFloat64(value pgtype.Numeric) (float64, error) {
	if !value.Valid {
		return 0, nil
	}
	floatValue, err := value.Float64Value()
	if err != nil {
		return 0, err
	}
	return floatValue.Float64, nil
}

func numericToFloat64Ptr(value pgtype.Numeric) (*float64, error) {
	if !value.Valid {
		return nil, nil
	}
	number, err := numericToFloat64(value)
	if err != nil {
		return nil, err
	}
	return &number, nil
}

func intPtr(value int) *int {
	return &value
}

func numericToBigFloat(value pgtype.Numeric) (*big.Float, error) {
	if !value.Valid || value.Int == nil {
		return nil, nil
	}
	rational := new(big.Rat).SetFrac(value.Int, big.NewInt(1))
	if value.Exp < 0 {
		scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(-value.Exp)), nil)
		rational = rational.Quo(rational, new(big.Rat).SetInt(scale))
	} else if value.Exp > 0 {
		scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(value.Exp)), nil)
		rational = rational.Mul(rational, new(big.Rat).SetInt(scale))
	}
	return new(big.Float).SetRat(rational), nil
}
