-- name: CreateMealOptionItem :one
INSERT INTO meal_option_items (option_id, food_id, quantity, measurement_unit, notes)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetMealOptionItemByID :one
SELECT * FROM meal_option_items
WHERE id = $1 AND option_id = $2;

-- name: UpdateMealOptionItem :one
UPDATE meal_option_items
SET quantity         = $2,
    measurement_unit = $3,
    notes            = $4
WHERE id = $1
RETURNING *;

-- name: DeleteMealOptionItem :exec
DELETE FROM meal_option_items
WHERE id = $1 AND option_id = $2;
