-- name: UpsertBodyMeasurement :one
INSERT INTO body_measurements (id, client_id, local_id, date, weight_kg, waist_cm, hip_cm, abdomen_cm, thigh_cm, chest_cm, wrist_cm, recorded_by)
VALUES (gen_random_uuid(), @client_id, @local_id, @date, @weight_kg, @waist_cm, @hip_cm, @abdomen_cm, @thigh_cm, @chest_cm, @wrist_cm, @recorded_by)
ON CONFLICT (client_id, date) DO UPDATE SET
    weight_kg   = EXCLUDED.weight_kg,
    waist_cm    = EXCLUDED.waist_cm,
    hip_cm      = EXCLUDED.hip_cm,
    abdomen_cm  = EXCLUDED.abdomen_cm,
    thigh_cm    = EXCLUDED.thigh_cm,
    chest_cm    = EXCLUDED.chest_cm,
    wrist_cm    = EXCLUDED.wrist_cm,
    recorded_by = EXCLUDED.recorded_by,
    updated_at  = NOW()
RETURNING *;

-- name: GetBodyMeasurementByDate :one
SELECT * FROM body_measurements WHERE client_id = @client_id AND date = @date;

-- name: ListBodyMeasurementsForNutritionist :many
SELECT bm.* FROM body_measurements bm
JOIN users u ON u.id = bm.client_id AND u.nutritionist_id = @nutritionist_id
WHERE bm.client_id = @client_id AND bm.date BETWEEN @from_date AND @to_date
ORDER BY bm.date DESC;

-- name: GetWeightHistory :many
SELECT date, weight_kg FROM body_measurements
WHERE client_id = @client_id AND weight_kg IS NOT NULL
ORDER BY date ASC;
