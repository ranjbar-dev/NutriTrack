-- name: CreateFoodCategory :one
INSERT INTO food_categories (name) VALUES ($1) RETURNING *;

-- name: GetFoodCategoryByID :one
SELECT * FROM food_categories WHERE id = $1 LIMIT 1;

-- name: GetFoodCategoryByName :one
SELECT * FROM food_categories WHERE name = $1 LIMIT 1;

-- name: ListFoodCategories :many
SELECT * FROM food_categories ORDER BY name;

-- name: DeleteFoodCategory :exec
DELETE FROM food_categories WHERE id = $1;
