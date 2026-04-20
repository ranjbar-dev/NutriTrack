-- name: CreateWaterLog :one
INSERT INTO water_logs (id, client_id, local_id, date, amount_ml, logged_time)
VALUES (gen_random_uuid(), @client_id, @local_id, @date, @amount_ml, @logged_time)
ON CONFLICT (local_id) DO NOTHING
RETURNING *;

-- name: GetWaterLogByLocalID :one
SELECT * FROM water_logs WHERE local_id = @local_id AND client_id = @client_id;

-- name: ListWaterLogsByDate :many
SELECT * FROM water_logs WHERE client_id = @client_id AND date = @date ORDER BY logged_time ASC NULLS LAST;

-- name: ListWaterLogsForNutritionist :many
SELECT wl.* FROM water_logs wl
JOIN users u ON u.id = wl.client_id AND u.nutritionist_id = @nutritionist_id
WHERE wl.client_id = @client_id AND wl.date BETWEEN @from_date AND @to_date
ORDER BY wl.date DESC, wl.logged_time DESC;

-- name: SumWaterByDate :one
SELECT COALESCE(SUM(amount_ml), 0)::bigint AS total_ml
FROM water_logs WHERE client_id = @client_id AND date = @date;
