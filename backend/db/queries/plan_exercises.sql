-- name: CreatePlanExercise :one
INSERT INTO plan_exercises (day_id, exercise_name, duration_minutes, description, calories_burn_estimate, display_order)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetPlanExerciseByID :one
SELECT * FROM plan_exercises
WHERE id = $1 AND day_id = $2;

-- name: UpdatePlanExercise :one
UPDATE plan_exercises
SET exercise_name          = $2,
    duration_minutes       = $3,
    description            = $4,
    calories_burn_estimate = $5
WHERE id = $1
RETURNING *;

-- name: DeletePlanExercise :exec
DELETE FROM plan_exercises
WHERE id = $1 AND day_id = $2;
