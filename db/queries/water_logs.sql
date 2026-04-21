-- name: UpsertWaterLog :one
INSERT INTO water_logs (client_id, local_id, logged_at, logged_date, amount_ml, notes)
VALUES (@client_id, @local_id, @logged_at, @logged_date, @amount_ml, @notes)
ON CONFLICT (client_id, local_id) DO UPDATE SET client_id = EXCLUDED.client_id
RETURNING *;

-- name: ListWaterLogsByClientAndDate :many
SELECT * FROM water_logs
WHERE client_id = @client_id AND logged_date = @logged_date
ORDER BY logged_at ASC;

-- name: ListWaterLogsByClient :many
SELECT * FROM water_logs
WHERE client_id = @client_id
ORDER BY logged_date DESC, logged_at DESC
LIMIT @limit OFFSET @offset;

-- name: CountWaterLogsByClient :one
SELECT COUNT(*) FROM water_logs WHERE client_id = @client_id;
