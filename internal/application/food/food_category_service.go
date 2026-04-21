package food

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranjbar-dev/nutritrack/internal/domain/food/entity"
	"github.com/ranjbar-dev/nutritrack/internal/domain/food/repository"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
)

// FoodCategoryService implements food category management use cases.
type FoodCategoryService struct {
	categoryRepo repository.FoodCategoryRepository
}

// NewFoodCategoryService constructs a FoodCategoryService.
func NewFoodCategoryService(repo repository.FoodCategoryRepository) *FoodCategoryService {
	return &FoodCategoryService{categoryRepo: repo}
}

// Create creates a new food category (superadmin only — enforced in handler).
func (s *FoodCategoryService) Create(ctx context.Context, name string) (*entity.FoodCategory, error) {
	normalized := shared.NormalizePersian(name)

	existing, err := s.categoryRepo.FindByName(ctx, normalized)
	if err != nil {
		return nil, shared.ErrInternal
	}
	if existing != nil {
		return nil, shared.ErrConflict
	}

	return s.categoryRepo.Create(ctx, normalized)
}

// ListAll returns all food categories.
func (s *FoodCategoryService) ListAll(ctx context.Context) ([]entity.FoodCategory, error) {
	return s.categoryRepo.ListAll(ctx)
}

// Delete removes a food category by ID (superadmin only).
func (s *FoodCategoryService) Delete(ctx context.Context, id uuid.UUID) error {
	cat, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return shared.ErrInternal
	}
	if cat == nil {
		return shared.ErrNotFound
	}
	return s.categoryRepo.Delete(ctx, id)
}
