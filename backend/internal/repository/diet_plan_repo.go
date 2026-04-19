package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ranjbar-dev/nutritrack/backend/internal/model/dto"
	"github.com/ranjbar-dev/nutritrack/backend/internal/repository/sqlc"
)

// ─── Batch SQL constants (raw queries — NOT sqlc-generated) ──────────────────

const (
	// Phase 1: plan header + days in a single query (LEFT JOIN so empty plan still returns)
	sqlGetPlanAndDays = `
		SELECT dp.id, dp.client_id, dp.nutritionist_id, dp.start_date, dp.end_date,
		       dp.notes, dp.daily_water_target_ml, dp.status, dp.created_at, dp.updated_at,
		       pd.id as day_id, pd.day_number, pd.label, pd.created_at as day_created_at
		FROM diet_plans dp
		LEFT JOIN plan_days pd ON pd.plan_id = dp.id
		WHERE dp.id = $1
		ORDER BY pd.day_number ASC`

	// Phase 2: sub-queries keyed on collected IDs
	sqlGetMealsByDayIDs = `
		SELECT id, day_id, title, scheduled_time, display_order
		FROM meals WHERE day_id = ANY($1)
		ORDER BY display_order ASC, scheduled_time ASC NULLS LAST`

	sqlGetOptionsByMealIDs = `
		SELECT id, meal_id, option_number, label
		FROM meal_options WHERE meal_id = ANY($1)
		ORDER BY option_number ASC`

	sqlGetItemsWithFoodByOptionIDs = `
		SELECT moi.id, moi.option_id, moi.food_id, moi.quantity, moi.measurement_unit, moi.notes,
		       f.name as food_name, f.calories, f.protein_g, f.carbs_g, f.fat_g, f.fiber_g,
		       f.measurement_unit as food_unit, f.measurement_amount
		FROM meal_option_items moi
		JOIN foods f ON moi.food_id = f.id
		WHERE moi.option_id = ANY($1)`

	sqlGetExercisesByDayIDs = `
		SELECT id, day_id, exercise_name, duration_minutes, description, calories_burn_estimate, display_order
		FROM plan_exercises WHERE day_id = ANY($1)
		ORDER BY display_order ASC`

	sqlGetMedicationsByPlanID = `
		SELECT pm.id, pm.plan_id, pm.medication_id, pm.dosage, pm.frequency, pm.times,
		       pm.instructions, pm.start_date, pm.end_date,
		       m.name as medication_name, m.form as medication_form
		FROM plan_medications pm
		JOIN medications m ON pm.medication_id = m.id
		WHERE pm.plan_id = $1`
)

// ─── Intermediate row structs for batch scanning ─────────────────────────────

type aggDayRow struct {
	ID        pgtype.UUID
	DayNumber int32
	Label     pgtype.Text
	CreatedAt pgtype.Timestamptz
}

type aggMealRow struct {
	ID            pgtype.UUID
	DayID         pgtype.UUID
	Title         string
	ScheduledTime pgtype.Time
	DisplayOrder  int32
}

type aggOptionRow struct {
	ID           pgtype.UUID
	MealID       pgtype.UUID
	OptionNumber int16
	Label        pgtype.Text
}

type aggItemRow struct {
	ID                pgtype.UUID
	OptionID          pgtype.UUID
	FoodID            pgtype.UUID
	Quantity          pgtype.Numeric
	MeasurementUnit   string
	Notes             pgtype.Text
	FoodName          string
	FoodCalories      pgtype.Numeric
	FoodProteinG      pgtype.Numeric
	FoodCarbsG        pgtype.Numeric
	FoodFatG          pgtype.Numeric
	FoodFiberG        pgtype.Numeric
	FoodUnit          string
	FoodMeasureAmount pgtype.Numeric
}

type aggExerciseRow struct {
	ID                   pgtype.UUID
	DayID                pgtype.UUID
	ExerciseName         string
	DurationMinutes      int32
	Description          pgtype.Text
	CaloriesBurnEstimate pgtype.Int4
	DisplayOrder         int32
}

type aggMedRow struct {
	ID             pgtype.UUID
	PlanID         pgtype.UUID
	MedicationID   pgtype.UUID
	Dosage         string
	Frequency      string
	Times          []byte
	Instructions   pgtype.Text
	StartDate      pgtype.Date
	EndDate        pgtype.Date
	MedicationName string
	MedicationForm string
}

// ─── Interface ────────────────────────────────────────────────────────────────

// DietPlanRepository defines all operations for the diet plan aggregate and sub-entities.
type DietPlanRepository interface {
	// Plan-level CRUD
	CreatePlan(ctx context.Context, params sqlc.CreateDietPlanParams) (*sqlc.DietPlan, error)
	GetPlanByID(ctx context.Context, planID, nutritionistID uuid.UUID) (*sqlc.DietPlan, error)
	ListClientPlans(ctx context.Context, params sqlc.ListClientPlansParams) ([]sqlc.DietPlan, error)
	CountClientPlans(ctx context.Context, params sqlc.CountClientPlansParams) (int64, error)
	ListMyPlans(ctx context.Context, params sqlc.ListMyPlansParams) ([]sqlc.DietPlan, error)
	CountMyPlans(ctx context.Context, clientID uuid.UUID) (int64, error)
	UpdatePlanHeader(ctx context.Context, params sqlc.UpdateDietPlanHeaderParams) (*sqlc.DietPlan, error)
	ActivatePlan(ctx context.Context, planID, nutritionistID uuid.UUID) error
	ArchivePreviousActivePlan(ctx context.Context, clientID, exceptPlanID uuid.UUID) error
	DeletePlan(ctx context.Context, planID, nutritionistID uuid.UUID) error
	GetPlanStatus(ctx context.Context, planID uuid.UUID) (sqlc.DietPlanStatus, error)
	CountPlanDays(ctx context.Context, planID uuid.UUID) (int64, error)

	// Day CRUD
	AddDay(ctx context.Context, params sqlc.CreatePlanDayParams) (*sqlc.PlanDay, error)
	GetDayByID(ctx context.Context, dayID, planID uuid.UUID) (*sqlc.PlanDay, error)
	UpdateDay(ctx context.Context, params sqlc.UpdatePlanDayParams) (*sqlc.PlanDay, error)
	DeleteDay(ctx context.Context, dayID, planID uuid.UUID) error

	// Meal CRUD
	AddMeal(ctx context.Context, params sqlc.CreateMealParams) (*sqlc.Meal, error)
	GetMealByID(ctx context.Context, mealID, dayID uuid.UUID) (*sqlc.Meal, error)
	UpdateMeal(ctx context.Context, params sqlc.UpdateMealParams) (*sqlc.Meal, error)
	DeleteMeal(ctx context.Context, mealID, dayID uuid.UUID) error
	ReorderMeal(ctx context.Context, mealID uuid.UUID, newOrder int32) error

	// Option CRUD
	AddOption(ctx context.Context, params sqlc.CreateMealOptionParams) (*sqlc.MealOption, error)
	GetOptionByID(ctx context.Context, optionID, mealID uuid.UUID) (*sqlc.MealOption, error)
	GetNextOptionNumber(ctx context.Context, mealID uuid.UUID) (int32, error)
	DeleteOption(ctx context.Context, optionID, mealID uuid.UUID) error

	// Item CRUD
	AddItem(ctx context.Context, params sqlc.CreateMealOptionItemParams) (*sqlc.MealOptionItem, error)
	GetItemByID(ctx context.Context, itemID, optionID uuid.UUID) (*sqlc.MealOptionItem, error)
	UpdateItem(ctx context.Context, params sqlc.UpdateMealOptionItemParams) (*sqlc.MealOptionItem, error)
	DeleteItem(ctx context.Context, itemID, optionID uuid.UUID) error

	// Exercise CRUD
	AddExercise(ctx context.Context, params sqlc.CreatePlanExerciseParams) (*sqlc.PlanExercise, error)
	GetExerciseByID(ctx context.Context, exerciseID, dayID uuid.UUID) (*sqlc.PlanExercise, error)
	UpdateExercise(ctx context.Context, params sqlc.UpdatePlanExerciseParams) (*sqlc.PlanExercise, error)
	DeleteExercise(ctx context.Context, exerciseID, dayID uuid.UUID) error

	// Medication prescription CRUD
	AddMedication(ctx context.Context, params sqlc.CreatePlanMedicationParams) (*sqlc.PlanMedication, error)
	GetMedicationByID(ctx context.Context, medID, planID uuid.UUID) (*sqlc.PlanMedication, error)
	UpdateMedication(ctx context.Context, params sqlc.UpdatePlanMedicationParams) (*sqlc.PlanMedication, error)
	DeleteMedication(ctx context.Context, medID, planID uuid.UUID) error

	// Batch aggregate (raw pgx — 2-phase, no N+1)
	GetFullPlanAggregate(ctx context.Context, planID uuid.UUID) (*dto.DietPlanResponse, error)
	GetActivePlanForClient(ctx context.Context, clientID uuid.UUID) (*dto.DietPlanResponse, error)
}

// ─── Implementation ───────────────────────────────────────────────────────────

type dietPlanRepository struct {
	q    *sqlc.Queries // standard sqlc CRUD
	pool *pgxpool.Pool // SendBatch aggregate queries
}

// NewDietPlanRepository creates a DietPlanRepository backed by the given pool.
// The pool is used both for sqlc.New (CRUD) and directly for SendBatch (aggregate).
func NewDietPlanRepository(pool *pgxpool.Pool) DietPlanRepository {
	return &dietPlanRepository{q: sqlc.New(pool), pool: pool}
}

// ─── Plan-level CRUD ──────────────────────────────────────────────────────────

func (r *dietPlanRepository) CreatePlan(ctx context.Context, params sqlc.CreateDietPlanParams) (*sqlc.DietPlan, error) {
	plan, err := r.q.CreateDietPlan(ctx, params)
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *dietPlanRepository) GetPlanByID(ctx context.Context, planID, nutritionistID uuid.UUID) (*sqlc.DietPlan, error) {
	plan, err := r.q.GetDietPlanByID(ctx, sqlc.GetDietPlanByIDParams{
		ID:             pgtype.UUID{Bytes: planID, Valid: true},
		NutritionistID: pgtype.UUID{Bytes: nutritionistID, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *dietPlanRepository) ListClientPlans(ctx context.Context, params sqlc.ListClientPlansParams) ([]sqlc.DietPlan, error) {
	return r.q.ListClientPlans(ctx, params)
}

func (r *dietPlanRepository) CountClientPlans(ctx context.Context, params sqlc.CountClientPlansParams) (int64, error) {
	return r.q.CountClientPlans(ctx, params)
}

func (r *dietPlanRepository) ListMyPlans(ctx context.Context, params sqlc.ListMyPlansParams) ([]sqlc.DietPlan, error) {
	return r.q.ListMyPlans(ctx, params)
}

func (r *dietPlanRepository) CountMyPlans(ctx context.Context, clientID uuid.UUID) (int64, error) {
	return r.q.CountMyPlans(ctx, pgtype.UUID{Bytes: clientID, Valid: true})
}

func (r *dietPlanRepository) UpdatePlanHeader(ctx context.Context, params sqlc.UpdateDietPlanHeaderParams) (*sqlc.DietPlan, error) {
	plan, err := r.q.UpdateDietPlanHeader(ctx, params)
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *dietPlanRepository) ActivatePlan(ctx context.Context, planID, nutritionistID uuid.UUID) error {
	return r.q.ActivateDietPlan(ctx, sqlc.ActivateDietPlanParams{
		ID:             pgtype.UUID{Bytes: planID, Valid: true},
		NutritionistID: pgtype.UUID{Bytes: nutritionistID, Valid: true},
	})
}

func (r *dietPlanRepository) ArchivePreviousActivePlan(ctx context.Context, clientID, exceptPlanID uuid.UUID) error {
	return r.q.ArchivePreviousActivePlan(ctx, sqlc.ArchivePreviousActivePlanParams{
		ClientID: pgtype.UUID{Bytes: clientID, Valid: true},
		ID:       pgtype.UUID{Bytes: exceptPlanID, Valid: true},
	})
}

func (r *dietPlanRepository) DeletePlan(ctx context.Context, planID, nutritionistID uuid.UUID) error {
	return r.q.DeleteDietPlan(ctx, sqlc.DeleteDietPlanParams{
		ID:             pgtype.UUID{Bytes: planID, Valid: true},
		NutritionistID: pgtype.UUID{Bytes: nutritionistID, Valid: true},
	})
}

func (r *dietPlanRepository) GetPlanStatus(ctx context.Context, planID uuid.UUID) (sqlc.DietPlanStatus, error) {
	return r.q.GetDietPlanStatus(ctx, pgtype.UUID{Bytes: planID, Valid: true})
}

func (r *dietPlanRepository) CountPlanDays(ctx context.Context, planID uuid.UUID) (int64, error) {
	return r.q.CountPlanDays(ctx, pgtype.UUID{Bytes: planID, Valid: true})
}

// ─── Day CRUD ─────────────────────────────────────────────────────────────────

func (r *dietPlanRepository) AddDay(ctx context.Context, params sqlc.CreatePlanDayParams) (*sqlc.PlanDay, error) {
	day, err := r.q.CreatePlanDay(ctx, params)
	if err != nil {
		return nil, err
	}
	return &day, nil
}

func (r *dietPlanRepository) GetDayByID(ctx context.Context, dayID, planID uuid.UUID) (*sqlc.PlanDay, error) {
	day, err := r.q.GetPlanDayByID(ctx, sqlc.GetPlanDayByIDParams{
		ID:     pgtype.UUID{Bytes: dayID, Valid: true},
		PlanID: pgtype.UUID{Bytes: planID, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return &day, nil
}

func (r *dietPlanRepository) UpdateDay(ctx context.Context, params sqlc.UpdatePlanDayParams) (*sqlc.PlanDay, error) {
	day, err := r.q.UpdatePlanDay(ctx, params)
	if err != nil {
		return nil, err
	}
	return &day, nil
}

func (r *dietPlanRepository) DeleteDay(ctx context.Context, dayID, planID uuid.UUID) error {
	return r.q.DeletePlanDay(ctx, sqlc.DeletePlanDayParams{
		ID:     pgtype.UUID{Bytes: dayID, Valid: true},
		PlanID: pgtype.UUID{Bytes: planID, Valid: true},
	})
}

// ─── Meal CRUD ────────────────────────────────────────────────────────────────

func (r *dietPlanRepository) AddMeal(ctx context.Context, params sqlc.CreateMealParams) (*sqlc.Meal, error) {
	meal, err := r.q.CreateMeal(ctx, params)
	if err != nil {
		return nil, err
	}
	return &meal, nil
}

func (r *dietPlanRepository) GetMealByID(ctx context.Context, mealID, dayID uuid.UUID) (*sqlc.Meal, error) {
	meal, err := r.q.GetMealByID(ctx, sqlc.GetMealByIDParams{
		ID:    pgtype.UUID{Bytes: mealID, Valid: true},
		DayID: pgtype.UUID{Bytes: dayID, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return &meal, nil
}

func (r *dietPlanRepository) UpdateMeal(ctx context.Context, params sqlc.UpdateMealParams) (*sqlc.Meal, error) {
	meal, err := r.q.UpdateMeal(ctx, params)
	if err != nil {
		return nil, err
	}
	return &meal, nil
}

func (r *dietPlanRepository) DeleteMeal(ctx context.Context, mealID, dayID uuid.UUID) error {
	return r.q.DeleteMeal(ctx, sqlc.DeleteMealParams{
		ID:    pgtype.UUID{Bytes: mealID, Valid: true},
		DayID: pgtype.UUID{Bytes: dayID, Valid: true},
	})
}

func (r *dietPlanRepository) ReorderMeal(ctx context.Context, mealID uuid.UUID, newOrder int32) error {
	return r.q.ReorderMeal(ctx, sqlc.ReorderMealParams{
		ID:           pgtype.UUID{Bytes: mealID, Valid: true},
		DisplayOrder: newOrder,
	})
}

// ─── Option CRUD ──────────────────────────────────────────────────────────────

func (r *dietPlanRepository) AddOption(ctx context.Context, params sqlc.CreateMealOptionParams) (*sqlc.MealOption, error) {
	opt, err := r.q.CreateMealOption(ctx, params)
	if err != nil {
		return nil, err
	}
	return &opt, nil
}

func (r *dietPlanRepository) GetOptionByID(ctx context.Context, optionID, mealID uuid.UUID) (*sqlc.MealOption, error) {
	opt, err := r.q.GetMealOptionByID(ctx, sqlc.GetMealOptionByIDParams{
		ID:     pgtype.UUID{Bytes: optionID, Valid: true},
		MealID: pgtype.UUID{Bytes: mealID, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return &opt, nil
}

func (r *dietPlanRepository) GetNextOptionNumber(ctx context.Context, mealID uuid.UUID) (int32, error) {
	return r.q.GetNextOptionNumber(ctx, pgtype.UUID{Bytes: mealID, Valid: true})
}

func (r *dietPlanRepository) DeleteOption(ctx context.Context, optionID, mealID uuid.UUID) error {
	return r.q.DeleteMealOption(ctx, sqlc.DeleteMealOptionParams{
		ID:     pgtype.UUID{Bytes: optionID, Valid: true},
		MealID: pgtype.UUID{Bytes: mealID, Valid: true},
	})
}

// ─── Item CRUD ────────────────────────────────────────────────────────────────

func (r *dietPlanRepository) AddItem(ctx context.Context, params sqlc.CreateMealOptionItemParams) (*sqlc.MealOptionItem, error) {
	item, err := r.q.CreateMealOptionItem(ctx, params)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *dietPlanRepository) GetItemByID(ctx context.Context, itemID, optionID uuid.UUID) (*sqlc.MealOptionItem, error) {
	item, err := r.q.GetMealOptionItemByID(ctx, sqlc.GetMealOptionItemByIDParams{
		ID:       pgtype.UUID{Bytes: itemID, Valid: true},
		OptionID: pgtype.UUID{Bytes: optionID, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *dietPlanRepository) UpdateItem(ctx context.Context, params sqlc.UpdateMealOptionItemParams) (*sqlc.MealOptionItem, error) {
	item, err := r.q.UpdateMealOptionItem(ctx, params)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *dietPlanRepository) DeleteItem(ctx context.Context, itemID, optionID uuid.UUID) error {
	return r.q.DeleteMealOptionItem(ctx, sqlc.DeleteMealOptionItemParams{
		ID:       pgtype.UUID{Bytes: itemID, Valid: true},
		OptionID: pgtype.UUID{Bytes: optionID, Valid: true},
	})
}

// ─── Exercise CRUD ────────────────────────────────────────────────────────────

func (r *dietPlanRepository) AddExercise(ctx context.Context, params sqlc.CreatePlanExerciseParams) (*sqlc.PlanExercise, error) {
	ex, err := r.q.CreatePlanExercise(ctx, params)
	if err != nil {
		return nil, err
	}
	return &ex, nil
}

func (r *dietPlanRepository) GetExerciseByID(ctx context.Context, exerciseID, dayID uuid.UUID) (*sqlc.PlanExercise, error) {
	ex, err := r.q.GetPlanExerciseByID(ctx, sqlc.GetPlanExerciseByIDParams{
		ID:    pgtype.UUID{Bytes: exerciseID, Valid: true},
		DayID: pgtype.UUID{Bytes: dayID, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return &ex, nil
}

func (r *dietPlanRepository) UpdateExercise(ctx context.Context, params sqlc.UpdatePlanExerciseParams) (*sqlc.PlanExercise, error) {
	ex, err := r.q.UpdatePlanExercise(ctx, params)
	if err != nil {
		return nil, err
	}
	return &ex, nil
}

func (r *dietPlanRepository) DeleteExercise(ctx context.Context, exerciseID, dayID uuid.UUID) error {
	return r.q.DeletePlanExercise(ctx, sqlc.DeletePlanExerciseParams{
		ID:    pgtype.UUID{Bytes: exerciseID, Valid: true},
		DayID: pgtype.UUID{Bytes: dayID, Valid: true},
	})
}

// ─── Medication CRUD ──────────────────────────────────────────────────────────

func (r *dietPlanRepository) AddMedication(ctx context.Context, params sqlc.CreatePlanMedicationParams) (*sqlc.PlanMedication, error) {
	med, err := r.q.CreatePlanMedication(ctx, params)
	if err != nil {
		return nil, err
	}
	return &med, nil
}

func (r *dietPlanRepository) GetMedicationByID(ctx context.Context, medID, planID uuid.UUID) (*sqlc.PlanMedication, error) {
	med, err := r.q.GetPlanMedicationByID(ctx, sqlc.GetPlanMedicationByIDParams{
		ID:     pgtype.UUID{Bytes: medID, Valid: true},
		PlanID: pgtype.UUID{Bytes: planID, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return &med, nil
}

func (r *dietPlanRepository) UpdateMedication(ctx context.Context, params sqlc.UpdatePlanMedicationParams) (*sqlc.PlanMedication, error) {
	med, err := r.q.UpdatePlanMedication(ctx, params)
	if err != nil {
		return nil, err
	}
	return &med, nil
}

func (r *dietPlanRepository) DeleteMedication(ctx context.Context, medID, planID uuid.UUID) error {
	return r.q.DeletePlanMedication(ctx, sqlc.DeletePlanMedicationParams{
		ID:     pgtype.UUID{Bytes: medID, Valid: true},
		PlanID: pgtype.UUID{Bytes: planID, Valid: true},
	})
}

// ─── Batch Aggregate: GetFullPlanAggregate ────────────────────────────────────

// GetFullPlanAggregate loads a complete plan tree (days→meals→options→items+food,
// exercises, medications) using a 2-phase approach to avoid N+1 queries.
//
// Phase 1 (1 round-trip): SELECT plan + days via LEFT JOIN.
// Phase 2 (3 sequential queries + 1 SendBatch of 3): meals, options, items+food, exercises, medications.
// Total: 2 round-trips for the full tree.
func (r *dietPlanRepository) GetFullPlanAggregate(ctx context.Context, planID uuid.UUID) (*dto.DietPlanResponse, error) {
	// ── Phase 1: plan header + days ──────────────────────────────────────────
	rows, err := r.pool.Query(ctx, sqlGetPlanAndDays, pgtype.UUID{Bytes: planID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("query plan and days: %w", err)
	}

	var plan *dto.DietPlanResponse
	var dayRows []aggDayRow
	var dayIDs []pgtype.UUID

	for rows.Next() {
		// Plan fields appear on every row (from dp.* columns)
		var (
			pID, pClientID, pNutrID        pgtype.UUID
			pStartDate, pEndDate           pgtype.Date
			pNotes                         pgtype.Text
			pWater                         pgtype.Int4
			pStatus                        sqlc.DietPlanStatus
			pCreatedAt, pUpdatedAt         pgtype.Timestamptz
			dayID                          pgtype.UUID
			dayNumber                      pgtype.Int4
			dayLabel                       pgtype.Text
			dayCreatedAt                   pgtype.Timestamptz
		)
		if err := rows.Scan(
			&pID, &pClientID, &pNutrID,
			&pStartDate, &pEndDate,
			&pNotes, &pWater, &pStatus,
			&pCreatedAt, &pUpdatedAt,
			&dayID, &dayNumber, &dayLabel, &dayCreatedAt,
		); err != nil {
			rows.Close()
			return nil, fmt.Errorf("scan plan+day row: %w", err)
		}

		// Initialize plan response on first row
		if plan == nil {
			plan = &dto.DietPlanResponse{
				ID:             uuid.UUID(pID.Bytes).String(),
				ClientID:       uuid.UUID(pClientID.Bytes).String(),
				NutritionistID: uuid.UUID(pNutrID.Bytes).String(),
				StartDate:      formatDate(pStartDate),
				EndDate:        formatDate(pEndDate),
				Status:         string(pStatus),
				CreatedAt:      formatTimestamptz(pCreatedAt),
				UpdatedAt:      formatTimestamptz(pUpdatedAt),
			}
			if pNotes.Valid {
				plan.Notes = &pNotes.String
			}
			if pWater.Valid {
				v := int(pWater.Int32)
				plan.DailyWaterTargetMl = &v
			}
		}

		// Collect day rows (dayID may be NULL if no days)
		if dayID.Valid {
			dayRows = append(dayRows, aggDayRow{
				ID:        dayID,
				DayNumber: dayNumber.Int32,
				Label:     dayLabel,
				CreatedAt: dayCreatedAt,
			})
			dayIDs = append(dayIDs, dayID)
		}
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate plan+day rows: %w", err)
	}
	if plan == nil {
		return nil, pgx.ErrNoRows
	}

	// No days → return plan with empty slices
	if len(dayIDs) == 0 {
		plan.Days = []dto.PlanDayResponse{}
		plan.Medications = []dto.PlanMedicationResponse{}
		return plan, nil
	}

	// ── Phase 2a: meals ───────────────────────────────────────────────────────
	mealRows, mealIDs, err := r.queryMeals(ctx, dayIDs)
	if err != nil {
		return nil, err
	}

	// ── Phase 2b: options ─────────────────────────────────────────────────────
	var optionRows []aggOptionRow
	var optionIDs []pgtype.UUID
	if len(mealIDs) > 0 {
		optionRows, optionIDs, err = r.queryOptions(ctx, mealIDs)
		if err != nil {
			return nil, err
		}
	}

	// ── Phase 2c: batch 3 remaining queries ───────────────────────────────────
	var itemRows []aggItemRow
	var exRows []aggExerciseRow
	var medRows []aggMedRow

	batch := &pgx.Batch{}
	if len(optionIDs) > 0 {
		batch.Queue(sqlGetItemsWithFoodByOptionIDs, optionIDs)
	} else {
		batch.Queue("SELECT 1 WHERE false", nil) // placeholder to keep cursor aligned
	}
	batch.Queue(sqlGetExercisesByDayIDs, dayIDs)
	batch.Queue(sqlGetMedicationsByPlanID, pgtype.UUID{Bytes: planID, Valid: true})

	br := r.pool.SendBatch(ctx, batch)
	defer br.Close() // MUST defer — releases connection even on error

	// Collect items+food
	itemRows2, err := r.scanItemRows(br)
	if err != nil {
		return nil, err
	}
	itemRows = append(itemRows, itemRows2...)

	// Collect exercises
	exRows, err = r.scanExerciseRows(br)
	if err != nil {
		return nil, err
	}

	// Collect medications
	medRows, err = r.scanMedRows(br)
	if err != nil {
		return nil, err
	}

	// ── Assemble tree ─────────────────────────────────────────────────────────
	return assemblePlanAggregate(plan, dayRows, mealRows, optionRows, itemRows, exRows, medRows)
}

// queryMeals fetches meals for given dayIDs and returns rows + collected mealIDs.
func (r *dietPlanRepository) queryMeals(ctx context.Context, dayIDs []pgtype.UUID) ([]aggMealRow, []pgtype.UUID, error) {
	rows, err := r.pool.Query(ctx, sqlGetMealsByDayIDs, dayIDs)
	if err != nil {
		return nil, nil, fmt.Errorf("query meals: %w", err)
	}
	defer rows.Close()

	var mealRows []aggMealRow
	var mealIDs []pgtype.UUID
	for rows.Next() {
		var mr aggMealRow
		if err := rows.Scan(&mr.ID, &mr.DayID, &mr.Title, &mr.ScheduledTime, &mr.DisplayOrder); err != nil {
			return nil, nil, fmt.Errorf("scan meal row: %w", err)
		}
		mealRows = append(mealRows, mr)
		mealIDs = append(mealIDs, mr.ID)
	}
	return mealRows, mealIDs, rows.Err()
}

// queryOptions fetches options for given mealIDs and returns rows + collected optionIDs.
func (r *dietPlanRepository) queryOptions(ctx context.Context, mealIDs []pgtype.UUID) ([]aggOptionRow, []pgtype.UUID, error) {
	rows, err := r.pool.Query(ctx, sqlGetOptionsByMealIDs, mealIDs)
	if err != nil {
		return nil, nil, fmt.Errorf("query options: %w", err)
	}
	defer rows.Close()

	var optionRows []aggOptionRow
	var optionIDs []pgtype.UUID
	for rows.Next() {
		var or aggOptionRow
		if err := rows.Scan(&or.ID, &or.MealID, &or.OptionNumber, &or.Label); err != nil {
			return nil, nil, fmt.Errorf("scan option row: %w", err)
		}
		optionRows = append(optionRows, or)
		optionIDs = append(optionIDs, or.ID)
	}
	return optionRows, optionIDs, rows.Err()
}

// scanItemRows collects the first BatchResults query (items+food).
func (r *dietPlanRepository) scanItemRows(br pgx.BatchResults) ([]aggItemRow, error) {
	rows, err := br.Query()
	if err != nil {
		return nil, fmt.Errorf("batch items query: %w", err)
	}
	defer rows.Close()

	var items []aggItemRow
	for rows.Next() {
		var ir aggItemRow
		if err := rows.Scan(
			&ir.ID, &ir.OptionID, &ir.FoodID,
			&ir.Quantity, &ir.MeasurementUnit, &ir.Notes,
			&ir.FoodName, &ir.FoodCalories, &ir.FoodProteinG, &ir.FoodCarbsG, &ir.FoodFatG, &ir.FoodFiberG,
			&ir.FoodUnit, &ir.FoodMeasureAmount,
		); err != nil {
			return nil, fmt.Errorf("scan item row: %w", err)
		}
		items = append(items, ir)
	}
	return items, rows.Err()
}

// scanExerciseRows collects the second BatchResults query (exercises).
func (r *dietPlanRepository) scanExerciseRows(br pgx.BatchResults) ([]aggExerciseRow, error) {
	rows, err := br.Query()
	if err != nil {
		return nil, fmt.Errorf("batch exercises query: %w", err)
	}
	defer rows.Close()

	var exercises []aggExerciseRow
	for rows.Next() {
		var er aggExerciseRow
		if err := rows.Scan(
			&er.ID, &er.DayID, &er.ExerciseName,
			&er.DurationMinutes, &er.Description, &er.CaloriesBurnEstimate, &er.DisplayOrder,
		); err != nil {
			return nil, fmt.Errorf("scan exercise row: %w", err)
		}
		exercises = append(exercises, er)
	}
	return exercises, rows.Err()
}

// scanMedRows collects the third BatchResults query (medications).
func (r *dietPlanRepository) scanMedRows(br pgx.BatchResults) ([]aggMedRow, error) {
	rows, err := br.Query()
	if err != nil {
		return nil, fmt.Errorf("batch medications query: %w", err)
	}
	defer rows.Close()

	var meds []aggMedRow
	for rows.Next() {
		var mr aggMedRow
		if err := rows.Scan(
			&mr.ID, &mr.PlanID, &mr.MedicationID,
			&mr.Dosage, &mr.Frequency, &mr.Times,
			&mr.Instructions, &mr.StartDate, &mr.EndDate,
			&mr.MedicationName, &mr.MedicationForm,
		); err != nil {
			return nil, fmt.Errorf("scan medication row: %w", err)
		}
		meds = append(meds, mr)
	}
	return meds, rows.Err()
}

// GetActivePlanForClient returns the active plan aggregate for a client, or pgx.ErrNoRows if none.
func (r *dietPlanRepository) GetActivePlanForClient(ctx context.Context, clientID uuid.UUID) (*dto.DietPlanResponse, error) {
	var planID pgtype.UUID
	row := r.pool.QueryRow(ctx, `SELECT id FROM diet_plans WHERE client_id = $1 AND status = 'active' LIMIT 1`,
		pgtype.UUID{Bytes: clientID, Valid: true})
	if err := row.Scan(&planID); err != nil {
		return nil, err // returns pgx.ErrNoRows naturally
	}
	return r.GetFullPlanAggregate(ctx, uuid.UUID(planID.Bytes))
}

// ─── Tree assembly (O(1) map-based) ──────────────────────────────────────────

func assemblePlanAggregate(
	plan *dto.DietPlanResponse,
	dayRows []aggDayRow,
	mealRows []aggMealRow,
	optionRows []aggOptionRow,
	itemRows []aggItemRow,
	exRows []aggExerciseRow,
	medRows []aggMedRow,
) (*dto.DietPlanResponse, error) {
	dayMap    := make(map[uuid.UUID]*dto.PlanDayResponse)
	mealMap   := make(map[uuid.UUID]*dto.MealResponse)
	optionMap := make(map[uuid.UUID]*dto.MealOptionResponse)

	plan.Days = make([]dto.PlanDayResponse, 0, len(dayRows))
	for _, dr := range dayRows {
		d := dto.PlanDayResponse{
			ID:        uuid.UUID(dr.ID.Bytes).String(),
			DayNumber: int(dr.DayNumber),
			Meals:     []dto.MealResponse{},
			Exercises: []dto.PlanExerciseResponse{},
		}
		if dr.Label.Valid {
			d.Label = &dr.Label.String
		}
		plan.Days = append(plan.Days, d)
		dayMap[uuid.UUID(dr.ID.Bytes)] = &plan.Days[len(plan.Days)-1]
	}

	for _, mr := range mealRows {
		m := dto.MealResponse{
			ID:           uuid.UUID(mr.ID.Bytes).String(),
			Title:        mr.Title,
			DisplayOrder: int(mr.DisplayOrder),
			Options:      []dto.MealOptionResponse{},
		}
		if mr.ScheduledTime.Valid {
			ts := formatTime(mr.ScheduledTime)
			m.ScheduledTime = &ts
		}
		dayKey := uuid.UUID(mr.DayID.Bytes)
		if day, ok := dayMap[dayKey]; ok {
			day.Meals = append(day.Meals, m)
			mealMap[uuid.UUID(mr.ID.Bytes)] = &day.Meals[len(day.Meals)-1]
		}
	}

	for _, or_ := range optionRows {
		o := dto.MealOptionResponse{
			ID:           uuid.UUID(or_.ID.Bytes).String(),
			OptionNumber: int(or_.OptionNumber),
			Items:        []dto.MealOptionItemResponse{},
		}
		if or_.Label.Valid {
			o.Label = &or_.Label.String
		}
		mealKey := uuid.UUID(or_.MealID.Bytes)
		if meal, ok := mealMap[mealKey]; ok {
			meal.Options = append(meal.Options, o)
			optionMap[uuid.UUID(or_.ID.Bytes)] = &meal.Options[len(meal.Options)-1]
		}
	}

	for _, ir := range itemRows {
		calories, _ := repoNumericToFloat64(ir.FoodCalories)
		proteinG, _ := repoNumericToFloat64(ir.FoodProteinG)
		carbsG, _   := repoNumericToFloat64(ir.FoodCarbsG)
		fatG, _     := repoNumericToFloat64(ir.FoodFatG)
		fiberG, _   := repoNumericToFloat64(ir.FoodFiberG)
		amount, _   := repoNumericToFloat64(ir.FoodMeasureAmount)
		qty, _      := repoNumericToFloat64(ir.Quantity)

		item := dto.MealOptionItemResponse{
			ID:              uuid.UUID(ir.ID.Bytes).String(),
			Quantity:        qty,
			MeasurementUnit: ir.MeasurementUnit,
			Food: dto.FoodEmbedded{
				ID:                uuid.UUID(ir.FoodID.Bytes).String(),
				Name:              ir.FoodName,
				Calories:          calories,
				ProteinG:          proteinG,
				CarbsG:            carbsG,
				FatG:              fatG,
				FiberG:            fiberG,
				MeasurementUnit:   ir.FoodUnit,
				MeasurementAmount: amount,
			},
		}
		if ir.Notes.Valid {
			item.Notes = &ir.Notes.String
		}
		optKey := uuid.UUID(ir.OptionID.Bytes)
		if opt, ok := optionMap[optKey]; ok {
			opt.Items = append(opt.Items, item)
		}
	}

	for _, er := range exRows {
		ex := dto.PlanExerciseResponse{
			ID:              uuid.UUID(er.ID.Bytes).String(),
			ExerciseName:    er.ExerciseName,
			DurationMinutes: int(er.DurationMinutes),
			DisplayOrder:    int(er.DisplayOrder),
		}
		if er.Description.Valid {
			ex.Description = &er.Description.String
		}
		if er.CaloriesBurnEstimate.Valid {
			v := int(er.CaloriesBurnEstimate.Int32)
			ex.CaloriesBurnEstimate = &v
		}
		dayKey := uuid.UUID(er.DayID.Bytes)
		if day, ok := dayMap[dayKey]; ok {
			day.Exercises = append(day.Exercises, ex)
		}
	}

	plan.Medications = make([]dto.PlanMedicationResponse, 0, len(medRows))
	for _, mr := range medRows {
		med := dto.PlanMedicationResponse{
			ID:             uuid.UUID(mr.ID.Bytes).String(),
			MedicationID:   uuid.UUID(mr.MedicationID.Bytes).String(),
			MedicationName: mr.MedicationName,
			MedicationForm: mr.MedicationForm,
			Dosage:         mr.Dosage,
			Frequency:      mr.Frequency,
			Times:          []string{},
		}
		if len(mr.Times) > 0 {
			_ = json.Unmarshal(mr.Times, &med.Times)
		}
		if mr.Instructions.Valid {
			med.Instructions = &mr.Instructions.String
		}
		if mr.StartDate.Valid {
			s := mr.StartDate.Time.Format("2006-01-02")
			med.StartDate = &s
		}
		if mr.EndDate.Valid {
			s := mr.EndDate.Time.Format("2006-01-02")
			med.EndDate = &s
		}
		plan.Medications = append(plan.Medications, med)
	}

	return plan, nil
}

// ─── Local helpers ────────────────────────────────────────────────────────────

func formatDate(d pgtype.Date) string {
	if !d.Valid {
		return ""
	}
	return d.Time.Format("2006-01-02")
}

func formatTimestamptz(ts pgtype.Timestamptz) string {
	if !ts.Valid {
		return ""
	}
	return ts.Time.Format(time.RFC3339)
}

func formatTime(t pgtype.Time) string {
	if !t.Valid {
		return ""
	}
	hours := t.Microseconds / 3_600_000_000
	minutes := (t.Microseconds % 3_600_000_000) / 60_000_000
	return fmt.Sprintf("%02d:%02d", hours, minutes)
}

func repoNumericToFloat64(n pgtype.Numeric) (float64, error) {
	if !n.Valid {
		return 0, nil
	}
	fv, err := n.Float64Value()
	if err != nil {
		return 0, err
	}
	return fv.Float64, nil
}

// repoNumericFromFloat64 converts a float64 to pgtype.Numeric (used by service tests).
func repoNumericFromFloat64(value float64) (pgtype.Numeric, error) {
	var numeric pgtype.Numeric
	if err := numeric.Scan(strconv.FormatFloat(value, 'f', -1, 64)); err != nil {
		return pgtype.Numeric{}, err
	}
	return numeric, nil
}
