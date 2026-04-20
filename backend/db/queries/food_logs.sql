-- name: UpsertFoodLog :one
INSERT INTO food_logs (id, client_id, local_id, date, meal_id, selected_option_id, is_skipped, notes)
VALUES (gen_random_uuid(), @client_id, @local_id, @date, @meal_id, @selected_option_id, @is_skipped, @notes)
ON CONFLICT (client_id, date, meal_id) DO UPDATE SET
    selected_option_id = EXCLUDED.selected_option_id,
    is_skipped         = EXCLUDED.is_skipped,
    notes              = EXCLUDED.notes,
    updated_at         = NOW()
RETURNING *;

-- name: GetFoodLogByLocalID :one
SELECT * FROM food_logs WHERE local_id = @local_id AND client_id = @client_id;

-- name: ListFoodLogsByDate :many
SELECT * FROM food_logs WHERE client_id = @client_id AND date = @date ORDER BY created_at ASC;

-- name: ListFoodLogsForNutritionist :many
SELECT fl.* FROM food_logs fl
JOIN users u ON u.id = fl.client_id AND u.nutritionist_id = @nutritionist_id
WHERE fl.client_id = @client_id AND fl.date BETWEEN @from_date AND @to_date
ORDER BY fl.date DESC, fl.created_at DESC;
