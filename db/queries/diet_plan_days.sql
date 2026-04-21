-- name: CreateDietPlanDay :one
INSERT INTO diet_plan_days (plan_id, day_number) VALUES ($1, $2) RETURNING *;

-- name: GetDietPlanDay :one
SELECT * FROM diet_plan_days WHERE id = $1 LIMIT 1;

-- name: ListDietPlanDays :many
SELECT * FROM diet_plan_days WHERE plan_id = $1 ORDER BY day_number;

-- name: DeleteDietPlanDay :exec
DELETE FROM diet_plan_days WHERE id = $1;
