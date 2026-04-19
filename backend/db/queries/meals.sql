-- name: CreateMeal :one
INSERT INTO meals (day_id, title, scheduled_time, display_order)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetMealByID :one
SELECT * FROM meals
WHERE id = $1 AND day_id = $2;

-- name: ListMeals :many
SELECT * FROM meals
WHERE day_id = $1
ORDER BY display_order ASC, scheduled_time ASC;

-- name: UpdateMeal :one
UPDATE meals
SET title          = $2,
    scheduled_time = $3,
    display_order  = $4,
    updated_at     = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteMeal :exec
DELETE FROM meals
WHERE id = $1 AND day_id = $2;

-- name: ReorderMeal :exec
UPDATE meals
SET display_order = $2
WHERE id = $1;
