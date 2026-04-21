-- name: CreateFood :one
INSERT INTO foods (name, name_normalized, unit, calories, protein, carbohydrate, fat, fiber, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetFoodByID :one
SELECT * FROM foods WHERE id = $1 LIMIT 1;

-- name: SearchFoods :many
SELECT * FROM foods
WHERE is_active = true
  AND (
    sqlc.arg(query)::text = ''
    OR similarity(name_normalized, sqlc.arg(query)::text) > 0.15
    OR name_normalized ILIKE '%' || sqlc.arg(query)::text || '%'
  )
ORDER BY
  CASE WHEN sqlc.arg(query)::text = '' THEN 0.0 ELSE -similarity(name_normalized, sqlc.arg(query)::text) END,
  created_at DESC
LIMIT sqlc.arg(lim)::int OFFSET sqlc.arg(off)::int;

-- name: CountSearchFoods :one
SELECT COUNT(*) FROM foods
WHERE is_active = true
  AND (
    sqlc.arg(query)::text = ''
    OR similarity(name_normalized, sqlc.arg(query)::text) > 0.15
    OR name_normalized ILIKE '%' || sqlc.arg(query)::text || '%'
  );

-- name: UpdateFood :one
UPDATE foods
SET name            = $2,
    name_normalized = $3,
    unit            = $4,
    calories        = $5,
    protein         = $6,
    carbohydrate    = $7,
    fat             = $8,
    fiber           = $9,
    updated_at      = NOW()
WHERE id = $1
RETURNING *;

-- name: DeactivateFood :exec
UPDATE foods SET is_active = false, updated_at = NOW() WHERE id = $1;

-- name: DeleteFood :exec
DELETE FROM foods WHERE id = $1;

-- name: GetFoodCategories :many
SELECT fc.id, fc.name, fc.created_at
FROM food_categories fc
JOIN food_category_mappings fcm ON fc.id = fcm.category_id
WHERE fcm.food_id = $1;

-- name: AddFoodCategory :exec
INSERT INTO food_category_mappings (food_id, category_id) VALUES ($1, $2) ON CONFLICT DO NOTHING;

-- name: RemoveFoodCategories :exec
DELETE FROM food_category_mappings WHERE food_id = $1;
