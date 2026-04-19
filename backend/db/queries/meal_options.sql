-- name: CreateMealOption :one
INSERT INTO meal_options (meal_id, option_number, label)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetMealOptionByID :one
SELECT * FROM meal_options
WHERE id = $1 AND meal_id = $2;

-- name: ListMealOptions :many
SELECT * FROM meal_options
WHERE meal_id = $1
ORDER BY option_number ASC;

-- name: DeleteMealOption :exec
DELETE FROM meal_options
WHERE id = $1 AND meal_id = $2;

-- name: GetNextOptionNumber :one
SELECT COALESCE(MAX(option_number), 0) + 1 AS next_option_number
FROM meal_options
WHERE meal_id = $1;
