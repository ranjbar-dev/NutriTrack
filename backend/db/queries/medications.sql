-- name: CreateMedication :one
INSERT INTO medications (
    name, name_normalized, generic_name, generic_name_normalized, form, dosage_unit, description, created_by
) VALUES (
    sqlc.arg('name'),
    LOWER(normalize_persian(sqlc.arg('name'))),
    sqlc.narg('generic_name'),
    CASE WHEN sqlc.narg('generic_name')::text IS NOT NULL THEN LOWER(normalize_persian(sqlc.narg('generic_name')::text)) ELSE NULL END,
    sqlc.arg('form')::medication_form,
    sqlc.narg('dosage_unit'),
    sqlc.narg('description'),
    sqlc.arg('created_by')
) RETURNING *;

-- name: GetMedicationByID :one
SELECT m.*, u.full_name AS creator_name
FROM medications m
JOIN users u ON m.created_by = u.id
WHERE m.id = $1;

-- name: ListMedications :many
SELECT m.*, u.full_name AS creator_name
FROM medications m
JOIN users u ON m.created_by = u.id
WHERE m.is_active = COALESCE(sqlc.narg('is_active')::boolean, true)
  AND (sqlc.narg('search')::text IS NULL
    OR m.name_normalized ILIKE '%' || LOWER(normalize_persian(sqlc.narg('search')::text)) || '%'
    OR (m.generic_name_normalized IS NOT NULL
        AND m.generic_name_normalized ILIKE '%' || LOWER(normalize_persian(sqlc.narg('search')::text)) || '%'))
ORDER BY
  CASE WHEN sqlc.narg('search')::text IS NOT NULL
    THEN similarity(m.name_normalized, LOWER(normalize_persian(sqlc.narg('search')::text)))
    ELSE 0 END DESC,
  m.created_at DESC
LIMIT sqlc.arg('limit_val')::bigint
OFFSET sqlc.arg('offset_val')::bigint;

-- name: CountMedications :one
SELECT COUNT(*)
FROM medications m
WHERE m.is_active = COALESCE(sqlc.narg('is_active')::boolean, true)
  AND (sqlc.narg('search')::text IS NULL
    OR m.name_normalized ILIKE '%' || LOWER(normalize_persian(sqlc.narg('search')::text)) || '%'
    OR (m.generic_name_normalized IS NOT NULL
        AND m.generic_name_normalized ILIKE '%' || LOWER(normalize_persian(sqlc.narg('search')::text)) || '%'));

-- name: UpdateMedication :one
UPDATE medications
SET name = sqlc.arg('name'),
    name_normalized = LOWER(normalize_persian(sqlc.arg('name'))),
    generic_name = sqlc.narg('generic_name'),
    generic_name_normalized = CASE WHEN sqlc.narg('generic_name')::text IS NOT NULL THEN LOWER(normalize_persian(sqlc.narg('generic_name')::text)) ELSE NULL END,
    form = sqlc.arg('form')::medication_form,
    dosage_unit = sqlc.narg('dosage_unit'),
    description = sqlc.narg('description'),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: SoftDeleteMedication :exec
UPDATE medications SET is_active = false, updated_at = NOW() WHERE id = $1;

-- name: SoftDeleteMedicationByOwner :exec
UPDATE medications SET is_active = false, updated_at = NOW() WHERE id = $1 AND created_by = $2;

-- name: CheckDuplicateMedicationName :one
SELECT EXISTS(
  SELECT 1 FROM medications
  WHERE LOWER(normalize_persian(name)) = LOWER(normalize_persian(sqlc.arg('name')))
    AND is_active = true
    AND id != COALESCE(sqlc.narg('exclude_id')::uuid, '00000000-0000-0000-0000-000000000000'::uuid)
) AS is_duplicate;
