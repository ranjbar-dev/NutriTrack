-- name: UpsertSleepLog :one
INSERT INTO sleep_logs (client_id, local_id, logged_date, sleep_start, sleep_end, duration_minutes, quality, notes)
VALUES (@client_id, @local_id, @logged_date, @sleep_start, @sleep_end, @duration_minutes, @quality, @notes)
ON CONFLICT (client_id, local_id) DO UPDATE SET client_id = EXCLUDED.client_id
RETURNING *;

-- name: ListSleepLogsByClientAndDate :many
SELECT * FROM sleep_logs
WHERE client_id = @client_id AND logged_date = @logged_date
ORDER BY sleep_start ASC;

-- name: ListSleepLogsByClient :many
SELECT * FROM sleep_logs
WHERE client_id = @client_id
ORDER BY logged_date DESC, sleep_start DESC
LIMIT @limit OFFSET @offset;

-- name: CountSleepLogsByClient :one
SELECT COUNT(*) FROM sleep_logs WHERE client_id = @client_id;
