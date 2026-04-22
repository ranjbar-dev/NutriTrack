package food

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranjbar-dev/nutritrack/internal/domain/food/entity"
	"github.com/ranjbar-dev/nutritrack/internal/domain/food/repository"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
)

// CreateFoodRequest carries input for the CreateFood use case.
type CreateFoodRequest struct {
	Name         string
	Unit         string
	Calories     float64
	Protein      float64
	Carbohydrate float64
	Fat          float64
	Fiber        float64
	Sugar        float64
	Sodium       float64
	Amount       float64
	CategoryIDs  []uuid.UUID
	CallerID     uuid.UUID
	CallerRole   string
}

// UpdateFoodRequest carries input for the UpdateFood use case.
type UpdateFoodRequest struct {
	ID           uuid.UUID
	Name         string
	Unit         string
	Calories     float64
	Protein      float64
	Carbohydrate float64
	Fat          float64
	Fiber        float64
	Sugar        float64
	Sodium       float64
	Amount       float64
	CategoryIDs  []uuid.UUID
	CallerID     uuid.UUID
	CallerRole   string
}

// FoodService implements food management use cases.
type FoodService struct {
	foodRepo     repository.FoodRepository
	categoryRepo repository.FoodCategoryRepository
}

// NewFoodService constructs a FoodService.
func NewFoodService(foodRepo repository.FoodRepository, categoryRepo repository.FoodCategoryRepository) *FoodService {
	return &FoodService{
		foodRepo:     foodRepo,
		categoryRepo: categoryRepo,
	}
}

// CreateFood creates a new food item.
// Only nutritionists and super admins may create foods.
func (s *FoodService) CreateFood(ctx context.Context, req CreateFoodRequest) (*entity.Food, error) {
	if req.CallerRole != "nutritionist" && req.CallerRole != "superadmin" {
		return nil, shared.ErrForbidden
	}

	normalized := shared.NormalizePersian(req.Name)

	categories := make([]entity.FoodCategory, 0, len(req.CategoryIDs))
	for _, catID := range req.CategoryIDs {
		cat, err := s.categoryRepo.FindByID(ctx, catID)
		if err != nil {
			return nil, err
		}
		if cat != nil {
			categories = append(categories, *cat)
		}
	}

	food, err := entity.NewFood(
		req.Name, normalized, req.Unit,
		req.Calories, req.Protein, req.Carbohydrate, req.Fat, req.Fiber, req.Sugar, req.Sodium, req.Amount,
		&req.CallerID, categories,
	)
	if err != nil {
		return nil, err
	}

	if err := s.foodRepo.Create(ctx, food); err != nil {
		return nil, err
	}

	return food, nil
}

// GetFood retrieves an active food by ID.
func (s *FoodService) GetFood(ctx context.Context, id uuid.UUID) (*entity.Food, error) {
	food, err := s.foodRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if food == nil || !food.IsActive() {
		return nil, shared.ErrFoodNotFound
	}
	return food, nil
}

// SearchFoods performs a pg_trgm similarity search on active foods, optionally filtered by category.
func (s *FoodService) SearchFoods(ctx context.Context, query string, categoryID *uuid.UUID, limit, offset int32) ([]*entity.Food, int64, error) {
	normalized := shared.NormalizePersian(query)

	if categoryID != nil {
		foods, err := s.foodRepo.SearchByCategory(ctx, *categoryID, normalized, limit, offset)
		if err != nil {
			return nil, 0, err
		}
		total, err := s.foodRepo.CountByCategory(ctx, *categoryID, normalized)
		if err != nil {
			return nil, 0, err
		}
		return foods, total, nil
	}

	foods, err := s.foodRepo.Search(ctx, normalized, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.foodRepo.CountSearch(ctx, normalized)
	if err != nil {
		return nil, 0, err
	}

	return foods, total, nil
}

// UpdateFood updates an existing food item.
// Nutritionists may only update foods they created; superadmin can update any.
func (s *FoodService) UpdateFood(ctx context.Context, req UpdateFoodRequest) (*entity.Food, error) {
	food, err := s.foodRepo.FindByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if food == nil || !food.IsActive() {
		return nil, shared.ErrFoodNotFound
	}

	// Row-level ownership check for nutritionists.
	if req.CallerRole == "nutritionist" {
		if food.CreatedBy() == nil || *food.CreatedBy() != req.CallerID {
			return nil, shared.ErrForbidden
		}
	} else if req.CallerRole != "superadmin" {
		return nil, shared.ErrForbidden
	}

	// Resolve category entities.
	categories := make([]entity.FoodCategory, 0, len(req.CategoryIDs))
	for _, catID := range req.CategoryIDs {
		cat, err := s.categoryRepo.FindByID(ctx, catID)
		if err != nil {
			return nil, err
		}
		if cat != nil {
			categories = append(categories, *cat)
		}
	}

	if err := food.Update(
		req.Name, shared.NormalizePersian(req.Name), req.Unit,
		req.Calories, req.Protein, req.Carbohydrate, req.Fat, req.Fiber, req.Sugar, req.Sodium, req.Amount,
		categories,
	); err != nil {
		return nil, err
	}

	if err := s.foodRepo.Update(ctx, food); err != nil {
		return nil, err
	}

	return food, nil
}

// DeleteFood removes a food item.
// Superadmin performs a hard delete; nutritionist performs a soft delete (deactivate).
func (s *FoodService) DeleteFood(ctx context.Context, id uuid.UUID, callerID uuid.UUID, callerRole string) error {
	food, err := s.foodRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if food == nil {
		return shared.ErrFoodNotFound
	}

	switch callerRole {
	case "superadmin":
		return s.foodRepo.Delete(ctx, id)
	case "nutritionist":
		if food.CreatedBy() == nil || *food.CreatedBy() != callerID {
			return shared.ErrForbidden
		}
		return s.foodRepo.Deactivate(ctx, id)
	default:
		return shared.ErrForbidden
	}
}
