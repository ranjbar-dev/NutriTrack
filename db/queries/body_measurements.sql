-- name: UpsertBodyMeasurement :one
INSERT INTO body_measurements (client_id, local_id, measured_at, measured_date, weight_kg, height_cm, waist_cm, hip_cm, chest_cm, arm_cm, notes)
VALUES (@client_id, @local_id, @measured_at, @measured_date, @weight_kg, @height_cm, @waist_cm, @hip_cm, @chest_cm, @arm_cm, @notes)
ON CONFLICT (client_id, local_id) DO UPDATE SET client_id = EXCLUDED.client_id
RETURNING *;

-- name: ListBodyMeasurementsByClientAndDate :many
SELECT * FROM body_measurements
WHERE client_id = @client_id AND measured_date = @measured_date
ORDER BY measured_at ASC;

-- name: ListBodyMeasurementsByClient :many
SELECT * FROM body_measurements
WHERE client_id = @client_id
ORDER BY measured_date DESC, measured_at DESC
LIMIT @limit OFFSET @offset;

-- name: CountBodyMeasurementsByClient :one
SELECT COUNT(*) FROM body_measurements WHERE client_id = @client_id;
