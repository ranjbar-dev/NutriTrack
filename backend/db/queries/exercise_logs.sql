-- name: CreateExerciseLog :one
INSERT INTO exercise_logs (
    id, client_id, local_id, date, exercise_name, duration_minutes, calories_burned, notes
) VALUES (
    gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7
)
ON CONFLICT (local_id) DO NOTHING
RETURNING *;

-- name: GetExerciseLogByLocalID :one
SELECT * FROM exercise_logs
WHERE local_id = $1
  AND client_id = $2;

-- name: ListExerciseLogsByDate :many
SELECT * FROM exercise_logs
WHERE client_id = $1
  AND date = $2
ORDER BY created_at DESC;

-- name: ListExerciseLogsByDateRange :many
SELECT * FROM exercise_logs
WHERE client_id = $1
  AND date BETWEEN $2 AND $3
ORDER BY date DESC, created_at DESC;

-- name: ListExerciseLogsForNutritionist :many
SELECT el.*
FROM exercise_logs el
JOIN users u
  ON u.id = el.client_id
 AND u.nutritionist_id = $1
WHERE el.client_id = $2
  AND el.date BETWEEN $3 AND $4
ORDER BY el.date DESC, el.created_at DESC;

