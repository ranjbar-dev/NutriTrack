-- name: UpsertSleepLog :one
INSERT INTO sleep_logs (id, client_id, local_id, date, sleep_time, wake_time, quality, notes)
VALUES (gen_random_uuid(), @client_id, @local_id, @date, @sleep_time, @wake_time, @quality, @notes)
ON CONFLICT (client_id, date) DO UPDATE SET
    sleep_time = EXCLUDED.sleep_time,
    wake_time  = EXCLUDED.wake_time,
    quality    = EXCLUDED.quality,
    notes      = EXCLUDED.notes,
    updated_at = NOW()
RETURNING *;

-- name: GetSleepLogByDate :one
SELECT * FROM sleep_logs WHERE client_id = @client_id AND date = @date;

-- name: ListSleepLogsForNutritionist :many
SELECT sl.* FROM sleep_logs sl
JOIN users u ON u.id = sl.client_id AND u.nutritionist_id = @nutritionist_id
WHERE sl.client_id = @client_id AND sl.date BETWEEN @from_date AND @to_date
ORDER BY sl.date DESC;
