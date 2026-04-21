package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranjbar-dev/nutritrack/internal/domain/food/entity"
)

// FoodRepository defines persistence operations for the Food aggregate.
type FoodRepository interface {
	Create(ctx context.Context, food *entity.Food) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Food, error)
	Update(ctx context.Context, food *entity.Food) error
	Delete(ctx context.Context, id uuid.UUID) error
	Deactivate(ctx context.Context, id uuid.UUID) error
	Search(ctx context.Context, query string, limit, offset int32) ([]*entity.Food, error)
	CountSearch(ctx context.Context, query string) (int64, error)
	SearchByCategory(ctx context.Context, categoryID uuid.UUID, query string, limit, offset int32) ([]*entity.Food, error)
	CountByCategory(ctx context.Context, categoryID uuid.UUID, query string) (int64, error)
}
