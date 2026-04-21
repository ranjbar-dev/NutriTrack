package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranjbar-dev/nutritrack/internal/domain/foodrequest/entity"
)

// FoodRequestRepository defines persistence operations for the FoodRequest aggregate.
type FoodRequestRepository interface {
	Create(ctx context.Context, req *entity.FoodRequest) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.FoodRequest, error)
	ListPending(ctx context.Context, nutritionistID uuid.UUID, limit, offset int32) ([]*entity.FoodRequest, error)
	CountPending(ctx context.Context, nutritionistID uuid.UUID) (int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.FoodRequestStatus, rejectionReason *string, createdFoodID *uuid.UUID) (*entity.FoodRequest, error)
}
