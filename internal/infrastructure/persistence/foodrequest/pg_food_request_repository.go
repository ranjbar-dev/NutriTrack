package foodrequest

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ranjbar-dev/nutritrack/internal/domain/foodrequest/entity"
	db "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/sqlc"
)

// PgFoodRequestRepository is the PostgreSQL implementation of FoodRequestRepository.
type PgFoodRequestRepository struct {
	queries *db.Queries
}

// NewPgFoodRequestRepository creates a new PgFoodRequestRepository.
func NewPgFoodRequestRepository(pool *pgxpool.Pool) *PgFoodRequestRepository {
	return &PgFoodRequestRepository{queries: db.New(pool)}
}

// Create inserts a new food request.
func (r *PgFoodRequestRepository) Create(ctx context.Context, req *entity.FoodRequest) error {
	row, err := r.queries.CreateFoodRequest(ctx, db.CreateFoodRequestParams{
		ClientID:       req.ClientID,
		NutritionistID: req.NutritionistID,
		FoodName:       req.FoodName,
	})
	if err != nil {
		return err
	}
	// Populate generated fields back onto the entity.
	req.ID = row.ID
	req.Status = entity.FoodRequestStatus(row.Status)
	req.CreatedAt = row.CreatedAt
	req.UpdatedAt = row.UpdatedAt
	return nil
}

// FindByID retrieves a food request by its ID. Returns nil if not found.
func (r *PgFoodRequestRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.FoodRequest, error) {
	row, err := r.queries.GetFoodRequest(ctx, id)
	if err != nil {
		if isNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return toDomain(row), nil
}

// ListPending returns pending food requests for a nutritionist with pagination.
func (r *PgFoodRequestRepository) ListPending(ctx context.Context, nutritionistID uuid.UUID, limit, offset int32) ([]*entity.FoodRequest, error) {
	rows, err := r.queries.ListPendingFoodRequests(ctx, db.ListPendingFoodRequestsParams{
		NutritionistID: nutritionistID,
		Limit:          limit,
		Offset:         offset,
	})
	if err != nil {
		return nil, err
	}
	result := make([]*entity.FoodRequest, len(rows))
	for i, row := range rows {
		result[i] = toDomain(row)
	}
	return result, nil
}

// CountPending counts pending food requests for a nutritionist.
func (r *PgFoodRequestRepository) CountPending(ctx context.Context, nutritionistID uuid.UUID) (int64, error) {
	return r.queries.CountPendingFoodRequests(ctx, nutritionistID)
}

// UpdateStatus updates the status, rejection reason, and created food ID of a food request.
func (r *PgFoodRequestRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.FoodRequestStatus, rejectionReason *string, createdFoodID *uuid.UUID) (*entity.FoodRequest, error) {
	row, err := r.queries.UpdateFoodRequestStatus(ctx, db.UpdateFoodRequestStatusParams{
		ID:              id,
		Status:          string(status),
		RejectionReason: rejectionReason,
		CreatedFoodID:   createdFoodID,
	})
	if err != nil {
		return nil, err
	}
	return toDomain(row), nil
}
