-- name: CreateMealOptionItem :one
INSERT INTO meal_option_items (option_id, food_id, quantity, unit, notes)
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetMealOptionItem :one
SELECT * FROM meal_option_items WHERE id = $1 LIMIT 1;

-- name: ListMealOptionItems :many
SELECT * FROM meal_option_items WHERE option_id = $1 ORDER BY created_at;

-- name: DeleteMealOptionItem :exec
DELETE FROM meal_option_items WHERE id = $1;

-- name: DeleteMealOptionItemsByOption :exec
DELETE FROM meal_option_items WHERE option_id = $1;
