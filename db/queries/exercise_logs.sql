-- name: UpsertExerciseLog :one
INSERT INTO exercise_logs (client_id, local_id, logged_at, logged_date, exercise_name, duration_minutes, calories_burned, notes)
VALUES (@client_id, @local_id, @logged_at, @logged_date, @exercise_name, @duration_minutes, @calories_burned, @notes)
ON CONFLICT (client_id, local_id) DO UPDATE SET client_id = EXCLUDED.client_id
RETURNING *;

-- name: ListExerciseLogsByClientAndDate :many
SELECT * FROM exercise_logs
WHERE client_id = @client_id AND logged_date = @logged_date
ORDER BY logged_at ASC;

-- name: ListExerciseLogsByClient :many
SELECT * FROM exercise_logs
WHERE client_id = @client_id
ORDER BY logged_date DESC, logged_at DESC
LIMIT @limit OFFSET @offset;

-- name: CountExerciseLogsByClient :one
SELECT COUNT(*) FROM exercise_logs WHERE client_id = @client_id;
