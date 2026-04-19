-- name: CreateDietPlan :one
INSERT INTO diet_plans (client_id, nutritionist_id, start_date, end_date, notes, daily_water_target_ml)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetDietPlanByID :one
SELECT * FROM diet_plans
WHERE id = $1 AND nutritionist_id = $2;

-- name: GetDietPlanByIDForClient :one
SELECT * FROM diet_plans
WHERE id = $1 AND client_id = $2 AND status != 'draft';

-- name: GetActivePlanIDForClient :one
SELECT id FROM diet_plans
WHERE client_id = $1 AND status = 'active';

-- name: ListClientPlans :many
SELECT * FROM diet_plans
WHERE client_id = $1 AND nutritionist_id = $2
ORDER BY created_at DESC
LIMIT sqlc.arg(limit_val) OFFSET sqlc.arg(offset_val);

-- name: CountClientPlans :one
SELECT COUNT(*) FROM diet_plans
WHERE client_id = $1 AND nutritionist_id = $2;

-- name: UpdateDietPlanHeader :one
UPDATE diet_plans
SET start_date            = $3,
    end_date              = $4,
    notes                 = $5,
    daily_water_target_ml = $6,
    updated_at            = NOW()
WHERE id = $1 AND nutritionist_id = $2
RETURNING *;

-- name: ActivateDietPlan :exec
UPDATE diet_plans
SET status = 'active', updated_at = NOW()
WHERE id = $1 AND nutritionist_id = $2;

-- name: ArchivePreviousActivePlan :exec
UPDATE diet_plans
SET status = 'archived', updated_at = NOW()
WHERE client_id = $1 AND status = 'active' AND id != $2;

-- name: DeleteDietPlan :exec
DELETE FROM diet_plans
WHERE id = $1 AND nutritionist_id = $2 AND status = 'draft';

-- name: GetDietPlanStatus :one
SELECT status FROM diet_plans WHERE id = $1;

-- name: ListMyPlans :many
SELECT * FROM diet_plans
WHERE client_id = $1 AND status != 'draft'
ORDER BY created_at DESC
LIMIT sqlc.arg(limit_val) OFFSET sqlc.arg(offset_val);

-- name: CountMyPlans :one
SELECT COUNT(*) FROM diet_plans
WHERE client_id = $1 AND status != 'draft';
