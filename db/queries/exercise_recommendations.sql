-- name: CreateExerciseRecommendation :one
INSERT INTO exercise_recommendations (day_id, exercise_name, duration_minutes, description, calories_burn_estimate)
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetExerciseRecommendation :one
SELECT * FROM exercise_recommendations WHERE id = $1 LIMIT 1;

-- name: ListExerciseRecommendations :many
SELECT * FROM exercise_recommendations WHERE day_id = $1 ORDER BY created_at;

-- name: DeleteExerciseRecommendation :exec
DELETE FROM exercise_recommendations WHERE id = $1;
