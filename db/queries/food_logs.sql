-- name: UpsertFoodLog :one
INSERT INTO food_logs (client_id, local_id, logged_at, logged_date, food_id, food_name, quantity, unit, calories, protein, carbs, fat, notes)
VALUES (@client_id, @local_id, @logged_at, @logged_date, @food_id, @food_name, @quantity, @unit, @calories, @protein, @carbs, @fat, @notes)
ON CONFLICT (client_id, local_id) DO UPDATE SET client_id = EXCLUDED.client_id
RETURNING *;

-- name: ListFoodLogsByClientAndDate :many
SELECT * FROM food_logs
WHERE client_id = @client_id AND logged_date = @logged_date
ORDER BY logged_at ASC;

-- name: ListFoodLogsByClient :many
SELECT * FROM food_logs
WHERE client_id = @client_id
ORDER BY logged_date DESC, logged_at DESC
LIMIT @limit OFFSET @offset;

-- name: CountFoodLogsByClient :one
SELECT COUNT(*) FROM food_logs WHERE client_id = @client_id;
