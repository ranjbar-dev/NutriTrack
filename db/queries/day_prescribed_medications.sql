-- name: CreateDayPrescribedMedication :one
INSERT INTO day_prescribed_medications (day_id, medication_id, dosage, frequency, times, instructions, start_date, end_date)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: GetDayPrescribedMedication :one
SELECT * FROM day_prescribed_medications WHERE id = $1 LIMIT 1;

-- name: ListDayPrescribedMedicationsWithMedication :many
SELECT
    dpm.id,
    dpm.day_id,
    dpm.medication_id,
    dpm.dosage,
    dpm.frequency,
    dpm.times,
    dpm.instructions,
    dpm.start_date,
    dpm.end_date,
    dpm.created_at,
    m.name AS medication_name,
    m.unit AS medication_unit
FROM day_prescribed_medications dpm
JOIN medications m ON dpm.medication_id = m.id
WHERE dpm.day_id = $1
ORDER BY dpm.created_at;

-- name: DeleteDayPrescribedMedication :exec
DELETE FROM day_prescribed_medications WHERE id = $1;
