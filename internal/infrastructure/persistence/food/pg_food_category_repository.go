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

// PgFoodCategoryRepository implements domain/food/repository.FoodCategoryRepository.
type PgFoodCategoryRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

// NewPgFoodCategoryRepository constructs a PgFoodCategoryRepository.
func NewPgFoodCategoryRepository(pool *pgxpool.Pool) *PgFoodCategoryRepository {
	return &PgFoodCategoryRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

// Create inserts a new food category.
func (r *PgFoodCategoryRepository) Create(ctx context.Context, name string) (*entity.FoodCategory, error) {
	cat, err := r.queries.CreateFoodCategory(ctx, name)
	if err != nil {
		return nil, shared.ErrInternal
	}
	result := categoryToDomain(cat)
	return &result, nil
}

// FindByID retrieves a food category by ID. Returns nil, nil when not found.
func (r *PgFoodCategoryRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.FoodCategory, error) {
	cat, err := r.queries.GetFoodCategoryByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, shared.ErrInternal
	}
	result := categoryToDomain(cat)
	return &result, nil
}

// FindByName retrieves a food category by name. Returns nil, nil when not found.
func (r *PgFoodCategoryRepository) FindByName(ctx context.Context, name string) (*entity.FoodCategory, error) {
	cat, err := r.queries.GetFoodCategoryByName(ctx, name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, shared.ErrInternal
	}
	result := categoryToDomain(cat)
	return &result, nil
}

// ListAll returns all food categories ordered by name.
func (r *PgFoodCategoryRepository) ListAll(ctx context.Context) ([]entity.FoodCategory, error) {
	rows, err := r.queries.ListFoodCategories(ctx)
	if err != nil {
		return nil, shared.ErrInternal
	}
	return categoriesToDomain(rows), nil
}

// Delete removes a food category by ID.
func (r *PgFoodCategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.queries.DeleteFoodCategory(ctx, id); err != nil {
		return shared.ErrInternal
	}
	return nil
}
