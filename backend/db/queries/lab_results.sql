-- name: CreateLabResult :one
INSERT INTO lab_results (
    id, client_id, local_id, uploaded_by, title, lab_type, test_date, file_path, external_link, original_filename, mime_type, file_size_bytes
) VALUES (
    gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
ON CONFLICT (local_id) DO NOTHING
RETURNING *;

-- name: GetLabResultByLocalID :one
SELECT * FROM lab_results
WHERE local_id = $1
  AND client_id = $2;

-- name: GetLabResultByID :one
SELECT * FROM lab_results
WHERE id = $1
  AND client_id = $2;

-- name: ListLabResultsByClient :many
SELECT * FROM lab_results
WHERE client_id = $1
ORDER BY test_date DESC, created_at DESC;

-- name: ListLabResultsForNutritionist :many
SELECT lr.*
FROM lab_results lr
JOIN users u
  ON u.id = lr.client_id
 AND u.nutritionist_id = $1
WHERE lr.client_id = $2
ORDER BY lr.test_date DESC, lr.created_at DESC;


-- name: GetLabResultForNutritionist :one
SELECT lr.*
FROM lab_results lr
JOIN users u
  ON u.id = lr.client_id
 AND u.nutritionist_id = $1
WHERE lr.id = $2
  AND lr.client_id = $3;
