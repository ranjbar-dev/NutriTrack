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

-- name: ListMealOptionItemsWithFood :many
SELECT
    moi.id,
    moi.option_id,
    moi.food_id,
    moi.quantity,
    moi.unit,
    moi.notes,
    moi.created_at,
    f.name        AS food_name,
    f.unit        AS food_unit,
    f.calories    AS food_calories,
    f.protein     AS food_protein,
    f.carbohydrate AS food_carbohydrate,
    f.fat         AS food_fat,
    f.fiber       AS food_fiber
FROM meal_option_items moi
JOIN foods f ON moi.food_id = f.id
WHERE moi.option_id = $1
ORDER BY moi.created_at;
