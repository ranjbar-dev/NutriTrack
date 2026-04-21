-- name: CreateDietMeal :one
INSERT INTO diet_meals (day_id, title, scheduled_time, display_order)
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetDietMeal :one
SELECT * FROM diet_meals WHERE id = $1 LIMIT 1;

-- name: ListDietMeals :many
SELECT * FROM diet_meals WHERE day_id = $1 ORDER BY display_order, scheduled_time;

-- name: DeleteDietMeal :exec
DELETE FROM diet_meals WHERE id = $1;
