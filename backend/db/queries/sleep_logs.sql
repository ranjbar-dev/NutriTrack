-- name: GetSleepLogByLocalID :one
SELECT * FROM sleep_logs
WHERE local_id = $1
  AND client_id = $2;

-- name: UpsertSleepLog :one
INSERT INTO sleep_logs (
    id, client_id, local_id, date, sleep_time, wake_time, quality, notes
) VALUES (
    gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7
)
ON CONFLICT (client_id, date)
DO UPDATE SET
    sleep_time = EXCLUDED.sleep_time,
    wake_time = EXCLUDED.wake_time,
    quality = EXCLUDED.quality,
    notes = EXCLUDED.notes,
    updated_at = NOW()
RETURNING *;

-- name: GetSleepLogByDate :one
SELECT * FROM sleep_logs
WHERE client_id = $1
  AND date = $2;

-- name: ListSleepLogsByDateRange :many
SELECT * FROM sleep_logs
WHERE client_id = $1
  AND date BETWEEN $2 AND $3
ORDER BY date DESC;

-- name: ListSleepLogsForNutritionist :many
SELECT sl.*
FROM sleep_logs sl
JOIN users u
  ON u.id = sl.client_id
 AND u.nutritionist_id = $1
WHERE sl.client_id = $2
  AND sl.date BETWEEN $3 AND $4
ORDER BY sl.date DESC;

