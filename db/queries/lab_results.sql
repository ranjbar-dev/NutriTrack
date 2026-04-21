-- name: CreateLabResult :one
INSERT INTO lab_results (client_id, nutritionist_id, file_path, original_name, file_type, file_size, notes)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, client_id, nutritionist_id, file_path, original_name, file_type, file_size, notes, created_at;

-- name: GetLabResultByID :one
SELECT id, client_id, nutritionist_id, file_path, original_name, file_type, file_size, notes, created_at
FROM lab_results
WHERE id = $1;

-- name: ListLabResultsByClientID :many
SELECT id, client_id, nutritionist_id, file_path, original_name, file_type, file_size, notes, created_at
FROM lab_results
WHERE client_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountLabResultsByClientID :one
SELECT COUNT(*) FROM lab_results WHERE client_id = $1;
