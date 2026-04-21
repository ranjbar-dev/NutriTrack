-- name: UpsertMedicationLog :one
INSERT INTO medication_logs (client_id, local_id, logged_at, logged_date, medication_id, medication_name, dosage, notes)
VALUES (@client_id, @local_id, @logged_at, @logged_date, @medication_id, @medication_name, @dosage, @notes)
ON CONFLICT (client_id, local_id) DO UPDATE SET client_id = EXCLUDED.client_id
RETURNING *;

-- name: ListMedicationLogsByClientAndDate :many
SELECT * FROM medication_logs
WHERE client_id = @client_id AND logged_date = @logged_date
ORDER BY logged_at ASC;

-- name: ListMedicationLogsByClient :many
SELECT * FROM medication_logs
WHERE client_id = @client_id
ORDER BY logged_date DESC, logged_at DESC
LIMIT @limit OFFSET @offset;

-- name: CountMedicationLogsByClient :one
SELECT COUNT(*) FROM medication_logs WHERE client_id = @client_id;
