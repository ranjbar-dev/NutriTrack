package foodrequest

import (
	"context"

	"github.com/google/uuid"
	appFood "github.com/ranjbar-dev/nutritrack/internal/application/food"
	frEntity "github.com/ranjbar-dev/nutritrack/internal/domain/foodrequest/entity"
	frRepo "github.com/ranjbar-dev/nutritrack/internal/domain/foodrequest/repository"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	userRepo "github.com/ranjbar-dev/nutritrack/internal/domain/user/repository"
)

// ApproveRequest carries input for approving a food request by creating the food.
type ApproveRequest struct {
	Name         string
	Unit         string
	Calories     float64
	Protein      float64
	Carbohydrate float64
	Fat          float64
	Fiber        float64
}

// FoodRequestService implements the food-request use cases.
type FoodRequestService struct {
	repo     frRepo.FoodRequestRepository
	userRepo userRepo.UserRepository
	foodSvc  *appFood.FoodService
}

// NewFoodRequestService constructs a FoodRequestService.
func NewFoodRequestService(repo frRepo.FoodRequestRepository, userRepo userRepo.UserRepository, foodSvc *appFood.FoodService) *FoodRequestService {
	return &FoodRequestService{repo: repo, userRepo: userRepo, foodSvc: foodSvc}
}

// Submit allows a client to request a custom food from their nutritionist.
func (s *FoodRequestService) Submit(ctx context.Context, clientID uuid.UUID, foodName string) (*frEntity.FoodRequest, error) {
	client, err := s.userRepo.FindByID(ctx, clientID)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, shared.ErrUserNotFound
	}
	if client.NutritionistID == nil {
		return nil, shared.ErrForbidden
	}

	req := &frEntity.FoodRequest{
		ClientID:       clientID,
		NutritionistID: *client.NutritionistID,
		FoodName:       foodName,
		Status:         frEntity.FoodRequestStatusPending,
	}

	if err := s.repo.Create(ctx, req); err != nil {
		return nil, err
	}

	return req, nil
}

// ListPending returns pending food requests for a nutritionist with pagination.
func (s *FoodRequestService) ListPending(ctx context.Context, nutritionistID uuid.UUID, limit, offset int32) ([]*frEntity.FoodRequest, int64, error) {
	user, err := s.userRepo.FindByID(ctx, nutritionistID)
	if err != nil {
		return nil, 0, err
	}
	if user == nil || user.Role != "nutritionist" {
		return nil, 0, shared.ErrForbidden
	}

	items, err := s.repo.ListPending(ctx, nutritionistID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountPending(ctx, nutritionistID)
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// Approve approves a food request by creating the food and updating the request status.
func (s *FoodRequestService) Approve(ctx context.Context, requestID uuid.UUID, nutritionistID uuid.UUID, req ApproveRequest) (*frEntity.FoodRequest, error) {
	foodReq, err := s.repo.FindByID(ctx, requestID)
	if err != nil {
		return nil, err
	}
	if foodReq == nil {
		return nil, shared.ErrFoodRequestNotFound
	}
	if foodReq.NutritionistID != nutritionistID {
		return nil, shared.ErrFoodRequestNotOwned
	}
	if !foodReq.IsPending() {
		return nil, shared.ErrFoodRequestAlreadyProcessed
	}

	food, err := s.foodSvc.CreateFood(ctx, appFood.CreateFoodRequest{
		Name:         req.Name,
		Unit:         req.Unit,
		Calories:     req.Calories,
		Protein:      req.Protein,
		Carbohydrate: req.Carbohydrate,
		Fat:          req.Fat,
		Fiber:        req.Fiber,
		CallerID:     nutritionistID,
		CallerRole:   "nutritionist",
	})
	if err != nil {
		return nil, err
	}

	updated, err := s.repo.UpdateStatus(ctx, requestID, frEntity.FoodRequestStatusApproved, nil, &food.ID)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

// Reject rejects a food request with a reason.
func (s *FoodRequestService) Reject(ctx context.Context, requestID uuid.UUID, nutritionistID uuid.UUID, reason string) (*frEntity.FoodRequest, error) {
	foodReq, err := s.repo.FindByID(ctx, requestID)
	if err != nil {
		return nil, err
	}
	if foodReq == nil {
		return nil, shared.ErrFoodRequestNotFound
	}
	if foodReq.NutritionistID != nutritionistID {
		return nil, shared.ErrFoodRequestNotOwned
	}
	if !foodReq.IsPending() {
		return nil, shared.ErrFoodRequestAlreadyProcessed
	}

	updated, err := s.repo.UpdateStatus(ctx, requestID, frEntity.FoodRequestStatusRejected, &reason, nil)
	if err != nil {
		return nil, err
	}

	return updated, nil
}
