-- name: CreatePlanDay :one
INSERT INTO plan_days (plan_id, day_number, label)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetPlanDayByID :one
SELECT * FROM plan_days
WHERE id = $1 AND plan_id = $2;

-- name: ListPlanDays :many
SELECT * FROM plan_days
WHERE plan_id = $1
ORDER BY day_number ASC;

-- name: UpdatePlanDay :one
UPDATE plan_days
SET label = $2
WHERE id = $1
RETURNING *;

-- name: DeletePlanDay :exec
DELETE FROM plan_days
WHERE id = $1 AND plan_id = $2;

-- name: CountPlanDays :one
SELECT COUNT(*) FROM plan_days WHERE plan_id = $1;
