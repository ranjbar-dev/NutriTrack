-- name: CreateDietPlan :one
INSERT INTO diet_plans (client_id, nutritionist_id, title, start_date, end_date, notes, daily_water_target_ml, status)
VALUES ($1, $2, $3, $4, $5, $6, $7, 'active')
RETURNING *;

-- name: ArchiveActivePlanForClient :exec
UPDATE diet_plans
SET status = 'archived', updated_at = NOW()
WHERE client_id = $1 AND status = 'active';

-- name: GetDietPlanByID :one
SELECT * FROM diet_plans WHERE id = $1 LIMIT 1;

-- name: GetActivePlanByClientID :one
SELECT * FROM diet_plans WHERE client_id = $1 AND status = 'active' LIMIT 1;

-- name: ListPlansByClientID :many
SELECT * FROM diet_plans WHERE client_id = $1
ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: CountPlansByClientID :one
SELECT COUNT(*) FROM diet_plans WHERE client_id = $1;

-- name: UpdateDietPlan :one
UPDATE diet_plans
SET title = $2, notes = $3, daily_water_target_ml = $4, updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: DeleteDietPlan :exec
DELETE FROM diet_plans WHERE id = $1;
