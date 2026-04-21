package db

import (
	"context"

	"github.com/google/uuid"
)

const createFoodRequest = `-- name: CreateFoodRequest :one
INSERT INTO food_requests (client_id, nutritionist_id, food_name)
VALUES ($1, $2, $3)
RETURNING id, client_id, nutritionist_id, food_name, status, rejection_reason, created_food_id, created_at, updated_at`

// CreateFoodRequestParams holds parameters for creating a food request.
type CreateFoodRequestParams struct {
	ClientID       uuid.UUID `db:"client_id"`
	NutritionistID uuid.UUID `db:"nutritionist_id"`
	FoodName       string    `db:"food_name"`
}

// CreateFoodRequest inserts a new food request and returns the created row.
func (q *Queries) CreateFoodRequest(ctx context.Context, arg CreateFoodRequestParams) (FoodRequest, error) {
	row := q.db.QueryRow(ctx, createFoodRequest,
		arg.ClientID, arg.NutritionistID, arg.FoodName,
	)
	var i FoodRequest
	err := row.Scan(
		&i.ID, &i.ClientID, &i.NutritionistID, &i.FoodName,
		&i.Status, &i.RejectionReason, &i.CreatedFoodID,
		&i.CreatedAt, &i.UpdatedAt,
	)
	return i, err
}

const getFoodRequest = `-- name: GetFoodRequest :one
SELECT id, client_id, nutritionist_id, food_name, status, rejection_reason, created_food_id, created_at, updated_at
FROM food_requests
WHERE id = $1`

// GetFoodRequest retrieves a food request by ID.
func (q *Queries) GetFoodRequest(ctx context.Context, id uuid.UUID) (FoodRequest, error) {
	row := q.db.QueryRow(ctx, getFoodRequest, id)
	var i FoodRequest
	err := row.Scan(
		&i.ID, &i.ClientID, &i.NutritionistID, &i.FoodName,
		&i.Status, &i.RejectionReason, &i.CreatedFoodID,
		&i.CreatedAt, &i.UpdatedAt,
	)
	return i, err
}

// ListPendingFoodRequestsParams holds parameters for listing pending food requests.
type ListPendingFoodRequestsParams struct {
	NutritionistID uuid.UUID `db:"nutritionist_id"`
	Limit          int32     `db:"limit"`
	Offset         int32     `db:"offset"`
}

const listPendingFoodRequests = `-- name: ListPendingFoodRequests :many
SELECT id, client_id, nutritionist_id, food_name, status, rejection_reason, created_food_id, created_at, updated_at
FROM food_requests
WHERE nutritionist_id = $1 AND status = 'pending'
ORDER BY created_at ASC
LIMIT $2 OFFSET $3`

// ListPendingFoodRequests returns paginated pending food requests for a nutritionist.
func (q *Queries) ListPendingFoodRequests(ctx context.Context, arg ListPendingFoodRequestsParams) ([]FoodRequest, error) {
	rows, err := q.db.Query(ctx, listPendingFoodRequests, arg.NutritionistID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FoodRequest
	for rows.Next() {
		var i FoodRequest
		if err := rows.Scan(
			&i.ID, &i.ClientID, &i.NutritionistID, &i.FoodName,
			&i.Status, &i.RejectionReason, &i.CreatedFoodID,
			&i.CreatedAt, &i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

const countPendingFoodRequests = `-- name: CountPendingFoodRequests :one
SELECT COUNT(*) FROM food_requests
WHERE nutritionist_id = $1 AND status = 'pending'`

// CountPendingFoodRequests counts pending food requests for a nutritionist.
func (q *Queries) CountPendingFoodRequests(ctx context.Context, nutritionistID uuid.UUID) (int64, error) {
	var count int64
	err := q.db.QueryRow(ctx, countPendingFoodRequests, nutritionistID).Scan(&count)
	return count, err
}

// UpdateFoodRequestStatusParams holds parameters for updating a food request status.
type UpdateFoodRequestStatusParams struct {
	ID              uuid.UUID  `db:"id"`
	Status          string     `db:"status"`
	RejectionReason *string    `db:"rejection_reason"`
	CreatedFoodID   *uuid.UUID `db:"created_food_id"`
}

const updateFoodRequestStatus = `-- name: UpdateFoodRequestStatus :one
UPDATE food_requests
SET status = $2, rejection_reason = $3, created_food_id = $4, updated_at = NOW()
WHERE id = $1
RETURNING id, client_id, nutritionist_id, food_name, status, rejection_reason, created_food_id, created_at, updated_at`

// UpdateFoodRequestStatus updates the status (and related fields) of a food request.
func (q *Queries) UpdateFoodRequestStatus(ctx context.Context, arg UpdateFoodRequestStatusParams) (FoodRequest, error) {
	row := q.db.QueryRow(ctx, updateFoodRequestStatus,
		arg.ID, arg.Status, arg.RejectionReason, arg.CreatedFoodID,
	)
	var i FoodRequest
	err := row.Scan(
		&i.ID, &i.ClientID, &i.NutritionistID, &i.FoodName,
		&i.Status, &i.RejectionReason, &i.CreatedFoodID,
		&i.CreatedAt, &i.UpdatedAt,
	)
	return i, err
}
