-- name: CreateWaterLog :one
INSERT INTO water_logs (
    id, client_id, local_id, date, amount_ml, logged_time
) VALUES (
    gen_random_uuid(), $1, $2, $3, $4, $5
)
ON CONFLICT (local_id) DO NOTHING
RETURNING *;

-- name: GetWaterLogByLocalID :one
SELECT * FROM water_logs
WHERE local_id = $1
  AND client_id = $2;

-- name: ListWaterLogsByDate :many
SELECT * FROM water_logs
WHERE client_id = $1
  AND date = $2
ORDER BY logged_time ASC, created_at ASC;

-- name: ListWaterLogsByDateRange :many
SELECT * FROM water_logs
WHERE client_id = $1
  AND date BETWEEN $2 AND $3
ORDER BY date DESC, logged_time DESC, created_at DESC;

-- name: ListWaterLogsForNutritionist :many
SELECT wl.*
FROM water_logs wl
JOIN users u
  ON u.id = wl.client_id
 AND u.nutritionist_id = $1
WHERE wl.client_id = $2
  AND wl.date BETWEEN $3 AND $4
ORDER BY wl.date DESC, wl.logged_time DESC, wl.created_at DESC;

-- name: SumWaterByDate :one
SELECT COALESCE(SUM(amount_ml), 0)::bigint AS total_ml
FROM water_logs
WHERE client_id = $1
  AND date = $2;

