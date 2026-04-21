-- name: CreateMedication :one
INSERT INTO medications (name, name_normalized, description, unit, created_by)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetMedicationByID :one
SELECT * FROM medications WHERE id = $1 LIMIT 1;

-- name: SearchMedications :many
SELECT * FROM medications
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

-- name: CountSearchMedications :one
SELECT COUNT(*) FROM medications
WHERE is_active = true
  AND (
    sqlc.arg(query)::text = ''
    OR similarity(name_normalized, sqlc.arg(query)::text) > 0.15
    OR name_normalized ILIKE '%' || sqlc.arg(query)::text || '%'
  );

-- name: UpdateMedication :one
UPDATE medications
SET name            = $2,
    name_normalized = $3,
    description     = $4,
    unit            = $5,
    updated_at      = NOW()
WHERE id = $1
RETURNING *;

-- name: DeactivateMedication :exec
UPDATE medications SET is_active = false, updated_at = NOW() WHERE id = $1;

-- name: DeleteMedication :exec
DELETE FROM medications WHERE id = $1;
