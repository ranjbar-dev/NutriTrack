-- name: CreateMealOption :one
INSERT INTO meal_options (meal_id, option_number) VALUES ($1, $2) RETURNING *;

-- name: GetMealOption :one
SELECT * FROM meal_options WHERE id = $1 LIMIT 1;

-- name: ListMealOptions :many
SELECT * FROM meal_options WHERE meal_id = $1 ORDER BY option_number;

-- name: DeleteMealOption :exec
DELETE FROM meal_options WHERE id = $1;
