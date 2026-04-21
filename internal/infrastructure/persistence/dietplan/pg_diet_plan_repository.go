package dietplan

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ranjbar-dev/nutritrack/internal/domain/dietplan/entity"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	db "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/sqlc"
)

// PgDietPlanRepository implements domain/dietplan/repository.DietPlanRepository using PostgreSQL + sqlc.
type PgDietPlanRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

// NewPgDietPlanRepository constructs a PgDietPlanRepository.
func NewPgDietPlanRepository(pool *pgxpool.Pool) *PgDietPlanRepository {
	return &PgDietPlanRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

// CreateWithArchive atomically archives any existing active plan for the client,
// then inserts the new plan in a single DB transaction.
func (r *PgDietPlanRepository) CreateWithArchive(ctx context.Context, plan *entity.DietPlan) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return shared.ErrInternal
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	// 1. Archive existing active plan (exec, no error if no rows).
	if err := qtx.ArchiveActivePlanForClient(ctx, plan.ClientID); err != nil {
		return shared.ErrInternal
	}

	// 2. Insert new plan.
	created, err := qtx.CreateDietPlan(ctx, db.CreateDietPlanParams{
		ClientID:           plan.ClientID,
		NutritionistID:     plan.NutritionistID,
		Title:              plan.Title,
		StartDate:          plan.StartDate,
		EndDate:            plan.EndDate,
		Notes:              plan.Notes,
		DailyWaterTargetMl: int32(plan.DailyWaterTargetML),
	})
	if err != nil {
		return shared.ErrInternal
	}

	if err := tx.Commit(ctx); err != nil {
		return shared.ErrInternal
	}

	plan.ID = created.ID
	plan.Status = entity.PlanStatusActive
	plan.CreatedAt = created.CreatedAt
	plan.UpdatedAt = created.UpdatedAt
	return nil
}

// FindByID retrieves a diet plan by its ID. Returns nil, nil when not found.
func (r *PgDietPlanRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.DietPlan, error) {
	row, err := r.queries.GetDietPlanByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, shared.ErrInternal
	}
	return dietPlanToDomain(row), nil
}

// FindActiveByClientID retrieves the active plan for a client. Returns nil, nil when not found.
func (r *PgDietPlanRepository) FindActiveByClientID(ctx context.Context, clientID uuid.UUID) (*entity.DietPlan, error) {
	row, err := r.queries.GetActivePlanByClientID(ctx, clientID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, shared.ErrInternal
	}
	return dietPlanToDomain(row), nil
}

// ListByClientID returns a paginated list of diet plans for a client.
func (r *PgDietPlanRepository) ListByClientID(ctx context.Context, clientID uuid.UUID, limit, offset int32) ([]*entity.DietPlan, error) {
	rows, err := r.queries.ListPlansByClientID(ctx, db.ListPlansByClientIDParams{
		ClientID: clientID,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, shared.ErrInternal
	}
	plans := make([]*entity.DietPlan, len(rows))
	for i, row := range rows {
		plans[i] = dietPlanToDomain(row)
	}
	return plans, nil
}

// CountByClientID returns the total count of diet plans for a client.
func (r *PgDietPlanRepository) CountByClientID(ctx context.Context, clientID uuid.UUID) (int64, error) {
	count, err := r.queries.CountPlansByClientID(ctx, clientID)
	if err != nil {
		return 0, shared.ErrInternal
	}
	return count, nil
}

// Update persists updated diet plan fields.
func (r *PgDietPlanRepository) Update(ctx context.Context, plan *entity.DietPlan) error {
	updated, err := r.queries.UpdateDietPlan(ctx, db.UpdateDietPlanParams{
		ID:                 plan.ID,
		Title:              plan.Title,
		Notes:              plan.Notes,
		DailyWaterTargetMl: int32(plan.DailyWaterTargetML),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return shared.ErrPlanNotFound
		}
		return shared.ErrInternal
	}
	plan.UpdatedAt = updated.UpdatedAt
	return nil
}

// Delete removes a diet plan by ID.
func (r *PgDietPlanRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.queries.DeleteDietPlan(ctx, id); err != nil {
		return shared.ErrInternal
	}
	return nil
}

// AddDay inserts a new diet plan day and populates the entity with DB-generated fields.
func (r *PgDietPlanRepository) AddDay(ctx context.Context, day *entity.DietPlanDay) error {
	created, err := r.queries.CreateDietPlanDay(ctx, db.CreateDietPlanDayParams{
		PlanID:    day.PlanID,
		DayNumber: int32(day.DayNumber),
	})
	if err != nil {
		return shared.ErrInternal
	}
	day.ID = created.ID
	day.CreatedAt = created.CreatedAt
	return nil
}

// FindDayByID retrieves a diet plan day by its ID. Returns nil, nil when not found.
func (r *PgDietPlanRepository) FindDayByID(ctx context.Context, id uuid.UUID) (*entity.DietPlanDay, error) {
	row, err := r.queries.GetDietPlanDay(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, shared.ErrInternal
	}
	return dietPlanDayToDomain(row), nil
}

// ListDays returns all days for a diet plan ordered by day_number.
func (r *PgDietPlanRepository) ListDays(ctx context.Context, planID uuid.UUID) ([]*entity.DietPlanDay, error) {
	rows, err := r.queries.ListDietPlanDays(ctx, planID)
	if err != nil {
		return nil, shared.ErrInternal
	}
	days := make([]*entity.DietPlanDay, len(rows))
	for i, row := range rows {
		days[i] = dietPlanDayToDomain(row)
	}
	return days, nil
}

// DeleteDay removes a diet plan day by ID.
func (r *PgDietPlanRepository) DeleteDay(ctx context.Context, id uuid.UUID) error {
	if err := r.queries.DeleteDietPlanDay(ctx, id); err != nil {
		return shared.ErrInternal
	}
	return nil
}

// AddMeal inserts a new diet meal and populates the entity with DB-generated fields.
func (r *PgDietPlanRepository) AddMeal(ctx context.Context, meal *entity.DietMeal) error {
	created, err := r.queries.CreateDietMeal(ctx, db.CreateDietMealParams{
		DayID:         meal.DayID,
		Title:         meal.Title,
		ScheduledTime: stringToPgtime(meal.ScheduledTime),
		DisplayOrder:  int32(meal.DisplayOrder),
	})
	if err != nil {
		return shared.ErrInternal
	}
	meal.ID = created.ID
	meal.CreatedAt = created.CreatedAt
	return nil
}

// FindMealByID retrieves a diet meal by its ID. Returns nil, nil when not found.
func (r *PgDietPlanRepository) FindMealByID(ctx context.Context, id uuid.UUID) (*entity.DietMeal, error) {
	row, err := r.queries.GetDietMeal(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, shared.ErrInternal
	}
	return dietMealToDomain(row), nil
}

// ListMeals returns all meals for a day ordered by display_order, scheduled_time.
func (r *PgDietPlanRepository) ListMeals(ctx context.Context, dayID uuid.UUID) ([]*entity.DietMeal, error) {
	rows, err := r.queries.ListDietMeals(ctx, dayID)
	if err != nil {
		return nil, shared.ErrInternal
	}
	meals := make([]*entity.DietMeal, len(rows))
	for i, row := range rows {
		meals[i] = dietMealToDomain(row)
	}
	return meals, nil
}

// DeleteMeal removes a diet meal by ID.
func (r *PgDietPlanRepository) DeleteMeal(ctx context.Context, id uuid.UUID) error {
	if err := r.queries.DeleteDietMeal(ctx, id); err != nil {
		return shared.ErrInternal
	}
	return nil
}

// AddOption inserts a new meal option and populates the entity with DB-generated fields.
func (r *PgDietPlanRepository) AddOption(ctx context.Context, option *entity.MealOption) error {
	created, err := r.queries.CreateMealOption(ctx, db.CreateMealOptionParams{
		MealID:       option.MealID,
		OptionNumber: int32(option.OptionNumber),
	})
	if err != nil {
		return shared.ErrInternal
	}
	option.ID = created.ID
	option.CreatedAt = created.CreatedAt
	return nil
}

// FindOptionByID retrieves a meal option by its ID. Returns nil, nil when not found.
func (r *PgDietPlanRepository) FindOptionByID(ctx context.Context, id uuid.UUID) (*entity.MealOption, error) {
	row, err := r.queries.GetMealOption(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, shared.ErrInternal
	}
	return mealOptionToDomain(row), nil
}

// ListOptions returns all options for a meal ordered by option_number.
func (r *PgDietPlanRepository) ListOptions(ctx context.Context, mealID uuid.UUID) ([]*entity.MealOption, error) {
	rows, err := r.queries.ListMealOptions(ctx, mealID)
	if err != nil {
		return nil, shared.ErrInternal
	}
	options := make([]*entity.MealOption, len(rows))
	for i, row := range rows {
		options[i] = mealOptionToDomain(row)
	}
	return options, nil
}

// DeleteOption removes a meal option by ID.
func (r *PgDietPlanRepository) DeleteOption(ctx context.Context, id uuid.UUID) error {
	if err := r.queries.DeleteMealOption(ctx, id); err != nil {
		return shared.ErrInternal
	}
	return nil
}

// AddItem inserts a new meal option item and populates the entity with DB-generated fields.
func (r *PgDietPlanRepository) AddItem(ctx context.Context, item *entity.MealOptionItem) error {
	created, err := r.queries.CreateMealOptionItem(ctx, db.CreateMealOptionItemParams{
		OptionID: item.OptionID,
		FoodID:   item.FoodID,
		Quantity: float64ToNumeric(item.Quantity),
		Unit:     item.Unit,
		Notes:    item.Notes,
	})
	if err != nil {
		return shared.ErrInternal
	}
	item.ID = created.ID
	item.CreatedAt = created.CreatedAt
	return nil
}

// FindItemByID retrieves a meal option item by its ID. Returns nil, nil when not found.
func (r *PgDietPlanRepository) FindItemByID(ctx context.Context, id uuid.UUID) (*entity.MealOptionItem, error) {
	row, err := r.queries.GetMealOptionItem(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, shared.ErrInternal
	}
	return mealOptionItemToDomain(row), nil
}

// ListItems returns all items for a meal option ordered by created_at.
func (r *PgDietPlanRepository) ListItems(ctx context.Context, optionID uuid.UUID) ([]*entity.MealOptionItem, error) {
	rows, err := r.queries.ListMealOptionItems(ctx, optionID)
	if err != nil {
		return nil, shared.ErrInternal
	}
	items := make([]*entity.MealOptionItem, len(rows))
	for i, row := range rows {
		items[i] = mealOptionItemToDomain(row)
	}
	return items, nil
}

// DeleteItem removes a meal option item by ID.
func (r *PgDietPlanRepository) DeleteItem(ctx context.Context, id uuid.UUID) error {
	if err := r.queries.DeleteMealOptionItem(ctx, id); err != nil {
		return shared.ErrInternal
	}
	return nil
}

// DeleteItemsByOption removes all meal option items for a given option.
func (r *PgDietPlanRepository) DeleteItemsByOption(ctx context.Context, optionID uuid.UUID) error {
	if err := r.queries.DeleteMealOptionItemsByOption(ctx, optionID); err != nil {
		return shared.ErrInternal
	}
	return nil
}

// ListItemsWithFood returns items for an option joined with food data.
func (r *PgDietPlanRepository) ListItemsWithFood(ctx context.Context, optionID uuid.UUID) ([]*entity.MealOptionItem, error) {
	rows, err := r.queries.ListMealOptionItemsWithFood(ctx, optionID)
	if err != nil {
		return nil, shared.ErrInternal
	}
	items := make([]*entity.MealOptionItem, len(rows))
	for i, row := range rows {
		items[i] = mealOptionItemWithFoodToDomain(row)
	}
	return items, nil
}

// AddExercise inserts a new exercise recommendation and populates the entity with DB-generated fields.
func (r *PgDietPlanRepository) AddExercise(ctx context.Context, ex *entity.ExerciseRecommendation) error {
	created, err := r.queries.CreateExerciseRecommendation(ctx, db.CreateExerciseRecommendationParams{
		DayID:                ex.DayID,
		ExerciseName:         ex.ExerciseName,
		DurationMinutes:      int32(ex.DurationMinutes),
		Description:          ex.Description,
		CaloriesBurnEstimate: int32(ex.CaloriesBurnEstimate),
	})
	if err != nil {
		return shared.ErrInternal
	}
	ex.ID = created.ID
	ex.CreatedAt = created.CreatedAt
	return nil
}

// FindExerciseByID retrieves an exercise recommendation by its ID. Returns nil, nil when not found.
func (r *PgDietPlanRepository) FindExerciseByID(ctx context.Context, id uuid.UUID) (*entity.ExerciseRecommendation, error) {
	row, err := r.queries.GetExerciseRecommendation(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, shared.ErrInternal
	}
	return exerciseToDomain(row), nil
}

// ListExercises returns all exercise recommendations for a day.
func (r *PgDietPlanRepository) ListExercises(ctx context.Context, dayID uuid.UUID) ([]*entity.ExerciseRecommendation, error) {
	rows, err := r.queries.ListExerciseRecommendations(ctx, dayID)
	if err != nil {
		return nil, shared.ErrInternal
	}
	result := make([]*entity.ExerciseRecommendation, len(rows))
	for i, row := range rows {
		result[i] = exerciseToDomain(row)
	}
	return result, nil
}

// DeleteExercise removes an exercise recommendation by ID.
func (r *PgDietPlanRepository) DeleteExercise(ctx context.Context, id uuid.UUID) error {
	if err := r.queries.DeleteExerciseRecommendation(ctx, id); err != nil {
		return shared.ErrInternal
	}
	return nil
}

// AddPrescription inserts a new prescribed medication and populates the entity with DB-generated fields.
func (r *PgDietPlanRepository) AddPrescription(ctx context.Context, rx *entity.PrescribedMedication) error {
	created, err := r.queries.CreateDayPrescribedMedication(ctx, db.CreateDayPrescribedMedicationParams{
		DayID:        rx.DayID,
		MedicationID: rx.MedicationID,
		Dosage:       rx.Dosage,
		Frequency:    rx.Frequency,
		Times:        rx.Times,
		Instructions: rx.Instructions,
		StartDate:    rx.StartDate,
		EndDate:      rx.EndDate,
	})
	if err != nil {
		return shared.ErrInternal
	}
	rx.ID = created.ID
	rx.CreatedAt = created.CreatedAt
	return nil
}

// FindPrescriptionByID retrieves a prescribed medication by its ID. Returns nil, nil when not found.
func (r *PgDietPlanRepository) FindPrescriptionByID(ctx context.Context, id uuid.UUID) (*entity.PrescribedMedication, error) {
	row, err := r.queries.GetDayPrescribedMedication(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, shared.ErrInternal
	}
	return &entity.PrescribedMedication{
		ID:           row.ID,
		DayID:        row.DayID,
		MedicationID: row.MedicationID,
		Dosage:       row.Dosage,
		Frequency:    row.Frequency,
		Times:        row.Times,
		Instructions: row.Instructions,
		StartDate:    row.StartDate,
		EndDate:      row.EndDate,
		CreatedAt:    row.CreatedAt,
	}, nil
}

// ListPrescriptionsWithMedication returns all prescribed medications for a day joined with medication data.
func (r *PgDietPlanRepository) ListPrescriptionsWithMedication(ctx context.Context, dayID uuid.UUID) ([]*entity.PrescribedMedication, error) {
	rows, err := r.queries.ListDayPrescribedMedicationsWithMedication(ctx, dayID)
	if err != nil {
		return nil, shared.ErrInternal
	}
	result := make([]*entity.PrescribedMedication, len(rows))
	for i, row := range rows {
		result[i] = prescriptionWithMedToDomain(row)
	}
	return result, nil
}

// DeletePrescription removes a prescribed medication by ID.
func (r *PgDietPlanRepository) DeletePrescription(ctx context.Context, id uuid.UUID) error {
	if err := r.queries.DeleteDayPrescribedMedication(ctx, id); err != nil {
		return shared.ErrInternal
	}
	return nil
}
