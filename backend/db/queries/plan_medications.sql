-- name: CreatePlanMedication :one
INSERT INTO plan_medications (plan_id, medication_id, dosage, frequency, times, instructions, start_date, end_date)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPlanMedicationByID :one
SELECT * FROM plan_medications
WHERE id = $1 AND plan_id = $2;

-- name: UpdatePlanMedication :one
UPDATE plan_medications
SET dosage        = $2,
    frequency     = $3,
    times         = $4,
    instructions  = $5,
    start_date    = $6,
    end_date      = $7,
    updated_at    = NOW()
WHERE id = $1
RETURNING *;

-- name: DeletePlanMedication :exec
DELETE FROM plan_medications
WHERE id = $1 AND plan_id = $2;
