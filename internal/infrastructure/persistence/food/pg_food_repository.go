package food

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ranjbar-dev/nutritrack/internal/domain/food/entity"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	db "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/sqlc"
)

// PgFoodRepository implements domain/food/repository.FoodRepository using PostgreSQL + sqlc.
type PgFoodRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

// NewPgFoodRepository constructs a PgFoodRepository.
func NewPgFoodRepository(pool *pgxpool.Pool) *PgFoodRepository {
	return &PgFoodRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

// Create inserts a new food row and populates the entity with DB-generated fields.
// Category mappings are inserted after the food row.
func (r *PgFoodRepository) Create(ctx context.Context, food *entity.Food) error {
	created, err := r.queries.CreateFood(ctx, db.CreateFoodParams{
		Name:           food.Name(),
		NameNormalized: food.NameNormalized(),
		Unit:           food.Unit(),
		Calories:       float64ToNumeric(food.Calories()),
		Protein:        float64ToNumeric(food.Protein()),
		Carbohydrate:   float64ToNumeric(food.Carbohydrate()),
		Fat:            float64ToNumeric(food.Fat()),
		Fiber:          float64ToNumeric(food.Fiber()),
		Sugar:          float64ToNumeric(food.Sugar()),
		Sodium:         float64ToNumeric(food.Sodium()),
		Amount:         float64ToNumeric(food.Amount()),
		CreatedBy:      uuidToPgtypeUUID(food.CreatedBy()),
	})
	if err != nil {
		return shared.ErrInternal
	}

	// Populate entity with DB-generated values.
	food.SetPersistedState(created.ID, created.IsActive, created.CreatedAt, created.UpdatedAt)

	// Insert category mappings.
	for _, cat := range food.Categories() {
		if err := r.queries.AddFoodCategory(ctx, db.AddFoodCategoryParams{
			FoodID:     food.ID(),
			CategoryID: cat.ID(),
		}); err != nil {
			return shared.ErrInternal
		}
	}

	return nil
}

// FindByID retrieves a food by its ID including its categories.
// Returns nil, nil when not found.
func (r *PgFoodRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Food, error) {
	row, err := r.queries.GetFoodByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, shared.ErrInternal
	}

	food := foodToDomain(row)

	// Load categories.
	cats, err := r.queries.GetFoodCategories(ctx, id)
	if err != nil {
		return nil, shared.ErrInternal
	}
	food.SetCategories(categoriesToDomain(cats))

	return food, nil
}

// Update persists updated food fields. Category mappings are replaced entirely.
func (r *PgFoodRepository) Update(ctx context.Context, food *entity.Food) error {
	updated, err := r.queries.UpdateFood(ctx, db.UpdateFoodParams{
		ID:             food.ID(),
		Name:           food.Name(),
		NameNormalized: food.NameNormalized(),
		Unit:           food.Unit(),
		Calories:       float64ToNumeric(food.Calories()),
		Protein:        float64ToNumeric(food.Protein()),
		Carbohydrate:   float64ToNumeric(food.Carbohydrate()),
		Fat:            float64ToNumeric(food.Fat()),
		Fiber:          float64ToNumeric(food.Fiber()),
		Sugar:          float64ToNumeric(food.Sugar()),
		Sodium:         float64ToNumeric(food.Sodium()),
		Amount:         float64ToNumeric(food.Amount()),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return shared.ErrFoodNotFound
		}
		return shared.ErrInternal
	}

	food.SetUpdatedAt(updated.UpdatedAt)

	// Replace category mappings.
	if err := r.queries.RemoveFoodCategories(ctx, food.ID()); err != nil {
		return shared.ErrInternal
	}
	for _, cat := range food.Categories() {
		if err := r.queries.AddFoodCategory(ctx, db.AddFoodCategoryParams{
			FoodID:     food.ID(),
			CategoryID: cat.ID(),
		}); err != nil {
			return shared.ErrInternal
		}
	}

	return nil
}

// Delete hard-deletes a food record (superadmin only).
func (r *PgFoodRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.queries.DeleteFood(ctx, id); err != nil {
		return shared.ErrInternal
	}
	return nil
}

// Deactivate soft-deletes a food record (nutritionist).
func (r *PgFoodRepository) Deactivate(ctx context.Context, id uuid.UUID) error {
	if err := r.queries.DeactivateFood(ctx, id); err != nil {
		return shared.ErrInternal
	}
	return nil
}

// Search returns active foods matching the query with pg_trgm similarity.
func (r *PgFoodRepository) Search(ctx context.Context, query string, limit, offset int32) ([]*entity.Food, error) {
	rows, err := r.queries.SearchFoods(ctx, db.SearchFoodsParams{
		Query: query,
		Lim:   limit,
		Off:   offset,
	})
	if err != nil {
		return nil, shared.ErrInternal
	}

	foods := make([]*entity.Food, len(rows))
	for i, row := range rows {
		foods[i] = foodToDomain(row)
	}
	return foods, nil
}

// CountSearch returns the total count of active foods matching the query.
func (r *PgFoodRepository) CountSearch(ctx context.Context, query string) (int64, error) {
	count, err := r.queries.CountSearchFoods(ctx, query)
	if err != nil {
		return 0, shared.ErrInternal
	}
	return count, nil
}

// SearchByCategory returns active foods in a given category matching the query.
func (r *PgFoodRepository) SearchByCategory(ctx context.Context, categoryID uuid.UUID, query string, limit, offset int32) ([]*entity.Food, error) {
	rows, err := r.queries.SearchFoodsByCategory(ctx, db.SearchFoodsByCategoryParams{
		CategoryID: categoryID,
		Query:      query,
		Lim:        limit,
		Off:        offset,
	})
	if err != nil {
		return nil, shared.ErrInternal
	}

	foods := make([]*entity.Food, len(rows))
	for i, row := range rows {
		foods[i] = foodToDomain(row)
	}
	return foods, nil
}

// CountByCategory returns the total count of active foods in a given category matching the query.
func (r *PgFoodRepository) CountByCategory(ctx context.Context, categoryID uuid.UUID, query string) (int64, error) {
	count, err := r.queries.CountSearchFoodsByCategory(ctx, db.CountSearchFoodsByCategoryParams{
		CategoryID: categoryID,
		Query:      query,
	})
	if err != nil {
		return 0, shared.ErrInternal
	}
	return count, nil
}
