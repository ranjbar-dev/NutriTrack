package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranjbar-dev/nutritrack/internal/domain/food/entity"
)

// FoodCategoryRepository defines persistence for food categories.
type FoodCategoryRepository interface {
	Create(ctx context.Context, name string) (*entity.FoodCategory, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.FoodCategory, error)
	FindByName(ctx context.Context, name string) (*entity.FoodCategory, error)
	ListAll(ctx context.Context) ([]entity.FoodCategory, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
