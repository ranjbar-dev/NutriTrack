-- name: CreateMedicationLog :one
INSERT INTO medication_logs (id, client_id, local_id, date, prescribed_medication_id, medication_name, dosage, taken_at, notes, is_self_reported)
VALUES (gen_random_uuid(), @client_id, @local_id, @date, @prescribed_medication_id, @medication_name, @dosage, @taken_at, @notes, @is_self_reported)
ON CONFLICT (local_id) DO NOTHING
RETURNING *;

-- name: GetMedicationLogByLocalID :one
SELECT * FROM medication_logs WHERE local_id = @local_id AND client_id = @client_id;

-- name: ListMedicationLogsByDate :many
SELECT * FROM medication_logs WHERE client_id = @client_id AND date = @date ORDER BY taken_at ASC;

-- name: ListMedicationLogsForNutritionist :many
SELECT ml.* FROM medication_logs ml
JOIN users u ON u.id = ml.client_id AND u.nutritionist_id = @nutritionist_id
WHERE ml.client_id = @client_id AND ml.date BETWEEN @from_date AND @to_date
ORDER BY ml.date DESC, ml.taken_at DESC;
