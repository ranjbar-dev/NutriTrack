-- name: CreateFoodRequest :one
INSERT INTO food_requests (client_id, nutritionist_id, food_name)
VALUES ($1, $2, $3)
RETURNING id, client_id, nutritionist_id, food_name, status, rejection_reason, created_food_id, created_at, updated_at;

-- name: GetFoodRequest :one
SELECT id, client_id, nutritionist_id, food_name, status, rejection_reason, created_food_id, created_at, updated_at
FROM food_requests
WHERE id = $1;

-- name: ListPendingFoodRequests :many
SELECT id, client_id, nutritionist_id, food_name, status, rejection_reason, created_food_id, created_at, updated_at
FROM food_requests
WHERE nutritionist_id = $1 AND status = 'pending'
ORDER BY created_at ASC
LIMIT $2 OFFSET $3;

-- name: CountPendingFoodRequests :one
SELECT COUNT(*) FROM food_requests
WHERE nutritionist_id = $1 AND status = 'pending';

-- name: UpdateFoodRequestStatus :one
UPDATE food_requests
SET status = $2, rejection_reason = $3, created_food_id = $4, updated_at = NOW()
WHERE id = $1
RETURNING id, client_id, nutritionist_id, food_name, status, rejection_reason, created_food_id, created_at, updated_at;
