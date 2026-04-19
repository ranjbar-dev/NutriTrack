package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/ranjbar-dev/nutritrack/backend/internal/repository/sqlc"
)

// FoodRepository defines operations on foods and food_categories.
type FoodRepository interface {
	Create(ctx context.Context, params sqlc.CreateFoodParams) (*sqlc.Food, error)
	AddCategory(ctx context.Context, foodID uuid.UUID, category string) error
	GetByID(ctx context.Context, id uuid.UUID) (*sqlc.GetFoodByIDRow, error)
	GetCategories(ctx context.Context, foodID uuid.UUID) ([]string, error)
	List(ctx context.Context, params sqlc.ListFoodsParams) ([]sqlc.ListFoodsRow, error)
	Count(ctx context.Context, params sqlc.CountFoodsParams) (int64, error)
	Update(ctx context.Context, params sqlc.UpdateFoodParams) (*sqlc.Food, error)
	SoftDelete(ctx context.Context, id uuid.UUID) error
	SoftDeleteByOwner(ctx context.Context, id, ownerID uuid.UUID) error
	DeleteCategories(ctx context.Context, foodID uuid.UUID) error
	CheckDuplicateName(ctx context.Context, name string, excludeID *uuid.UUID) (bool, error)
	CountActive(ctx context.Context) (int64, error)
}

type foodRepository struct {
	q *sqlc.Queries
}

// NewFoodRepository creates a new FoodRepository backed by the given sqlc.DBTX.
func NewFoodRepository(db sqlc.DBTX) FoodRepository {
	return &foodRepository{q: sqlc.New(db)}
}

func (r *foodRepository) Create(ctx context.Context, params sqlc.CreateFoodParams) (*sqlc.Food, error) {
	food, err := r.q.CreateFood(ctx, params)
	if err != nil {
		return nil, err
	}
	return &food, nil
}

func (r *foodRepository) AddCategory(ctx context.Context, foodID uuid.UUID, category string) error {
	return r.q.AddFoodCategory(ctx, sqlc.AddFoodCategoryParams{
		FoodID:  pgtype.UUID{Bytes: foodID, Valid: true},
		Column2: sqlc.FoodCategoryType(category),
	})
}

func (r *foodRepository) GetByID(ctx context.Context, id uuid.UUID) (*sqlc.GetFoodByIDRow, error) {
	food, err := r.q.GetFoodByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}
	return &food, nil
}

func (r *foodRepository) GetCategories(ctx context.Context, foodID uuid.UUID) ([]string, error) {
	return r.q.GetFoodCategories(ctx, pgtype.UUID{Bytes: foodID, Valid: true})
}

func (r *foodRepository) List(ctx context.Context, params sqlc.ListFoodsParams) ([]sqlc.ListFoodsRow, error) {
	return r.q.ListFoods(ctx, params)
}

func (r *foodRepository) Count(ctx context.Context, params sqlc.CountFoodsParams) (int64, error) {
	return r.q.CountFoods(ctx, params)
}

func (r *foodRepository) Update(ctx context.Context, params sqlc.UpdateFoodParams) (*sqlc.Food, error) {
	food, err := r.q.UpdateFood(ctx, params)
	if err != nil {
		return nil, err
	}
	return &food, nil
}

func (r *foodRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.q.SoftDeleteFood(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

func (r *foodRepository) SoftDeleteByOwner(ctx context.Context, id, ownerID uuid.UUID) error {
	return r.q.SoftDeleteFoodByOwner(ctx, sqlc.SoftDeleteFoodByOwnerParams{
		ID:        pgtype.UUID{Bytes: id, Valid: true},
		CreatedBy: pgtype.UUID{Bytes: ownerID, Valid: true},
	})
}

func (r *foodRepository) DeleteCategories(ctx context.Context, foodID uuid.UUID) error {
	return r.q.DeleteFoodCategories(ctx, pgtype.UUID{Bytes: foodID, Valid: true})
}

func (r *foodRepository) CheckDuplicateName(ctx context.Context, name string, excludeID *uuid.UUID) (bool, error) {
	arg := sqlc.CheckDuplicateFoodNameParams{
		NormalizePersian: name,
		ExcludeID:        pgtype.UUID{Valid: false},
	}
	if excludeID != nil {
		arg.ExcludeID = pgtype.UUID{Bytes: *excludeID, Valid: true}
	}
	return r.q.CheckDuplicateFoodName(ctx, arg)
}

func (r *foodRepository) CountActive(ctx context.Context) (int64, error) {
	return r.q.CountActiveFoods(ctx)
}
