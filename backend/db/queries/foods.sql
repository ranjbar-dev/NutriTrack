-- name: CreateFood :one
INSERT INTO foods (
    name, name_normalized, description, calories, protein_g, carbs_g, fat_g, fiber_g, sugar_g, sodium_mg,
    measurement_unit, measurement_amount, created_by
) VALUES (
    $1, LOWER(normalize_persian($1)), $2, $3, $4, $5, $6, $7, $8, $9,
    $10, $11, $12
) RETURNING *;

-- name: AddFoodCategory :exec
INSERT INTO food_categories (food_id, category) VALUES ($1, $2::food_category_type) ON CONFLICT DO NOTHING;

-- name: GetFoodByID :one
SELECT f.*, u.full_name AS creator_name
FROM foods f
JOIN users u ON f.created_by = u.id
WHERE f.id = $1;

-- name: GetFoodCategories :many
SELECT category::text AS category FROM food_categories WHERE food_id = $1;

-- name: ListFoods :many
SELECT f.*, u.full_name AS creator_name
FROM foods f
JOIN users u ON f.created_by = u.id
WHERE f.is_active = COALESCE(sqlc.narg('is_active')::boolean, true)
  AND (sqlc.narg('search')::text IS NULL
    OR f.name_normalized ILIKE '%' || LOWER(normalize_persian(sqlc.narg('search')::text)) || '%')
  AND (sqlc.narg('category')::text IS NULL
    OR EXISTS (
      SELECT 1 FROM food_categories fc
      WHERE fc.food_id = f.id
        AND fc.category::text = sqlc.narg('category')::text
    ))
ORDER BY
  CASE WHEN sqlc.narg('search')::text IS NOT NULL
    THEN similarity(f.name_normalized, LOWER(normalize_persian(sqlc.narg('search')::text)))
    ELSE 0 END DESC,
  f.created_at DESC
LIMIT sqlc.arg('limit_val')::bigint
OFFSET sqlc.arg('offset_val')::bigint;

-- name: CountFoods :one
SELECT COUNT(*)
FROM foods f
WHERE f.is_active = COALESCE(sqlc.narg('is_active')::boolean, true)
  AND (sqlc.narg('search')::text IS NULL
    OR f.name_normalized ILIKE '%' || LOWER(normalize_persian(sqlc.narg('search')::text)) || '%')
  AND (sqlc.narg('category')::text IS NULL
    OR EXISTS (
      SELECT 1 FROM food_categories fc
      WHERE fc.food_id = f.id
        AND fc.category::text = sqlc.narg('category')::text
    ));

-- name: UpdateFood :one
UPDATE foods
SET name = $2,
    name_normalized = LOWER(normalize_persian($2)),
    description = $3,
    calories = $4,
    protein_g = $5,
    carbs_g = $6,
    fat_g = $7,
    fiber_g = $8,
    sugar_g = $9,
    sodium_mg = $10,
    measurement_unit = $11,
    measurement_amount = $12,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: SoftDeleteFood :exec
UPDATE foods SET is_active = false, updated_at = NOW() WHERE id = $1;

-- name: SoftDeleteFoodByOwner :exec
UPDATE foods SET is_active = false, updated_at = NOW() WHERE id = $1 AND created_by = $2;

-- name: DeleteFoodCategories :exec
DELETE FROM food_categories WHERE food_id = $1;

-- name: CheckDuplicateFoodName :one
SELECT EXISTS(
  SELECT 1 FROM foods
  WHERE LOWER(normalize_persian(name)) = LOWER(normalize_persian($1))
    AND is_active = true
    AND id != COALESCE(sqlc.narg('exclude_id')::uuid, '00000000-0000-0000-0000-000000000000'::uuid)
) AS is_duplicate;

-- name: CountActiveFoods :one
SELECT COUNT(*) FROM foods WHERE is_active = true;
