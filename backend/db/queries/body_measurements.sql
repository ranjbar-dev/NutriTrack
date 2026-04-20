-- name: GetBodyMeasurementByLocalID :one
SELECT * FROM body_measurements
WHERE local_id = $1
  AND client_id = $2;

-- name: UpsertBodyMeasurement :one
INSERT INTO body_measurements (
    id, client_id, local_id, date, weight_kg, waist_cm, hip_cm, abdomen_cm, thigh_cm, chest_cm, wrist_cm, recorded_by
) VALUES (
    gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
ON CONFLICT (client_id, date)
DO UPDATE SET
    weight_kg = EXCLUDED.weight_kg,
    waist_cm = EXCLUDED.waist_cm,
    hip_cm = EXCLUDED.hip_cm,
    abdomen_cm = EXCLUDED.abdomen_cm,
    thigh_cm = EXCLUDED.thigh_cm,
    chest_cm = EXCLUDED.chest_cm,
    wrist_cm = EXCLUDED.wrist_cm,
    recorded_by = EXCLUDED.recorded_by,
    updated_at = NOW()
RETURNING *;

-- name: GetBodyMeasurementByDate :one
SELECT * FROM body_measurements
WHERE client_id = $1
  AND date = $2;

-- name: ListBodyMeasurementsByDateRange :many
SELECT * FROM body_measurements
WHERE client_id = $1
  AND date BETWEEN $2 AND $3
ORDER BY date DESC;

-- name: ListBodyMeasurementsForNutritionist :many
SELECT bm.*
FROM body_measurements bm
JOIN users u
  ON u.id = bm.client_id
 AND u.nutritionist_id = $1
WHERE bm.client_id = $2
  AND bm.date BETWEEN $3 AND $4
ORDER BY bm.date DESC;

-- name: ListWeightHistory :many
SELECT date, weight_kg
FROM body_measurements
WHERE client_id = $1
  AND weight_kg IS NOT NULL
  AND date BETWEEN $2 AND $3
ORDER BY date ASC;

-- name: ListWeightHistoryForNutritionist :many
SELECT bm.date, bm.weight_kg
FROM body_measurements bm
JOIN users u
  ON u.id = bm.client_id
 AND u.nutritionist_id = $1
WHERE bm.client_id = $2
  AND bm.weight_kg IS NOT NULL
  AND bm.date BETWEEN $3 AND $4
ORDER BY bm.date ASC;

