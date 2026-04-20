-- name: CreateMedicationLog :one
INSERT INTO medication_logs (
    id, client_id, local_id, date, prescribed_medication_id, medication_name, dosage, taken_at, notes, is_self_reported
) VALUES (
    gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, $9
)
ON CONFLICT (local_id) DO NOTHING
RETURNING *;

-- name: GetMedicationLogByLocalID :one
SELECT * FROM medication_logs
WHERE local_id = $1
  AND client_id = $2;

-- name: ListMedicationLogsByDate :many
SELECT * FROM medication_logs
WHERE client_id = $1
  AND date = $2
ORDER BY taken_at ASC, created_at ASC;

-- name: ListMedicationLogsByDateRange :many
SELECT * FROM medication_logs
WHERE client_id = $1
  AND date BETWEEN $2 AND $3
ORDER BY date DESC, taken_at DESC, created_at DESC;

-- name: ListMedicationLogsForNutritionist :many
SELECT ml.*
FROM medication_logs ml
JOIN users u
  ON u.id = ml.client_id
 AND u.nutritionist_id = $1
WHERE ml.client_id = $2
  AND ml.date BETWEEN $3 AND $4
ORDER BY ml.date DESC, ml.taken_at DESC, ml.created_at DESC;

