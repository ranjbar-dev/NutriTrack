-- name: CreateFoodRequest :one
INSERT INTO food_requests (food_name, description, requested_by)
VALUES ($1, $2, $3)
RETURNING id, food_name, description, status, rejection_reason, requested_by, reviewed_by, created_at, updated_at;

-- name: ListFoodRequestsByClient :many
SELECT id, food_name, description, status, rejection_reason, requested_by, reviewed_by, created_at, updated_at
FROM food_requests
WHERE requested_by = $1
ORDER BY created_at DESC;

-- name: ListPendingFoodRequestsForNutritionist :many
SELECT fr.id, fr.food_name, fr.description, fr.status, fr.rejection_reason, fr.requested_by, fr.reviewed_by, fr.created_at, fr.updated_at
FROM food_requests fr
JOIN users u ON u.id = fr.requested_by AND u.nutritionist_id = $1
WHERE fr.status = 'pending'
ORDER BY fr.created_at ASC;

-- name: GetFoodRequestByID :one
SELECT id, food_name, description, status, rejection_reason, requested_by, reviewed_by, created_at, updated_at
FROM food_requests
WHERE id = $1;

-- name: ApproveFoodRequest :one
UPDATE food_requests
SET status = 'approved', reviewed_by = $2, updated_at = NOW()
WHERE id = $1 AND status = 'pending'
RETURNING id, food_name, description, status, rejection_reason, requested_by, reviewed_by, created_at, updated_at;

-- name: RejectFoodRequest :one
UPDATE food_requests
SET status = 'rejected', rejection_reason = $2, reviewed_by = $3, updated_at = NOW()
WHERE id = $1 AND status = 'pending'
RETURNING id, food_name, description, status, rejection_reason, requested_by, reviewed_by, created_at, updated_at;
