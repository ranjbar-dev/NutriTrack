package foodrequest

import (
	"context"

	"github.com/google/uuid"
	appFood "github.com/ranjbar-dev/nutritrack/internal/application/food"
	foodEntity "github.com/ranjbar-dev/nutritrack/internal/domain/food/entity"
	frEntity "github.com/ranjbar-dev/nutritrack/internal/domain/foodrequest/entity"
	frRepo "github.com/ranjbar-dev/nutritrack/internal/domain/foodrequest/repository"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	userRepo "github.com/ranjbar-dev/nutritrack/internal/domain/user/repository"
)

// FoodCreator is the port through which FoodRequestService creates foods.
// Defined here to avoid a concrete dependency on *appFood.FoodService.
type FoodCreator interface {
	CreateFood(ctx context.Context, req appFood.CreateFoodRequest) (*foodEntity.Food, error)
}

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
	foodSvc  FoodCreator
}

// NewFoodRequestService constructs a FoodRequestService.
func NewFoodRequestService(repo frRepo.FoodRequestRepository, userRepo userRepo.UserRepository, foodSvc FoodCreator) *FoodRequestService {
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
	if client.GetNutritionistID() == nil {
		return nil, shared.ErrForbidden
	}

	req, err := frEntity.NewFoodRequest(clientID, *client.GetNutritionistID(), foodName)
	if err != nil {
		return nil, err
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
	if user == nil || !user.IsNutritionist() {
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
	if foodReq.GetNutritionistID() != nutritionistID {
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

	updated, err := s.repo.UpdateStatus(ctx, requestID, frEntity.FoodRequestStatusApproved, nil, func() *uuid.UUID { id := food.ID(); return &id }())
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
	if foodReq.GetNutritionistID() != nutritionistID {
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
