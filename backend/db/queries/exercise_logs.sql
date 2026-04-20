-- name: CreateExerciseLog :one
INSERT INTO exercise_logs (id, client_id, local_id, date, exercise_name, duration_minutes, calories_burned, notes)
VALUES (gen_random_uuid(), @client_id, @local_id, @date, @exercise_name, @duration_minutes, @calories_burned, @notes)
ON CONFLICT (local_id) DO NOTHING
RETURNING *;

-- name: GetExerciseLogByLocalID :one
SELECT * FROM exercise_logs WHERE local_id = @local_id AND client_id = @client_id;

-- name: ListExerciseLogsByDate :many
SELECT * FROM exercise_logs WHERE client_id = @client_id AND date = @date ORDER BY created_at ASC;

-- name: ListExerciseLogsForNutritionist :many
SELECT el.* FROM exercise_logs el
JOIN users u ON u.id = el.client_id AND u.nutritionist_id = @nutritionist_id
WHERE el.client_id = @client_id AND el.date BETWEEN @from_date AND @to_date
ORDER BY el.date DESC, el.created_at DESC;
