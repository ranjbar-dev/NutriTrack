-- name: GetFoodLogByLocalID :one
SELECT * FROM food_logs
WHERE local_id = $1
  AND client_id = $2;

-- name: UpsertFoodLog :one
INSERT INTO food_logs (
    id, client_id, local_id, date, meal_id, selected_option_id, is_skipped, notes
) VALUES (
    gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7
)
ON CONFLICT (client_id, date, meal_id)
DO UPDATE SET
    selected_option_id = EXCLUDED.selected_option_id,
    is_skipped = EXCLUDED.is_skipped,
    notes = EXCLUDED.notes,
    updated_at = NOW()
RETURNING *;

-- name: ListFoodLogsByDate :many
SELECT * FROM food_logs
WHERE client_id = $1
  AND date = $2
ORDER BY created_at ASC;

-- name: ListFoodLogsByDateRange :many
SELECT * FROM food_logs
WHERE client_id = $1
  AND date BETWEEN $2 AND $3
ORDER BY date DESC, created_at DESC;

-- name: ListFoodLogsForNutritionist :many
SELECT fl.*
FROM food_logs fl
JOIN users u
  ON u.id = fl.client_id
 AND u.nutritionist_id = $1
WHERE fl.client_id = $2
  AND fl.date BETWEEN $3 AND $4
ORDER BY fl.date DESC, fl.created_at DESC;

