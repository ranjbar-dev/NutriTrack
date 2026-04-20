-- name: CreateLabResult :one
INSERT INTO lab_results (id, client_id, uploaded_by, title, lab_type, test_date, file_path, external_link, original_filename, mime_type, file_size_bytes)
VALUES (gen_random_uuid(), @client_id, @uploaded_by, @title, @lab_type, @test_date, @file_path, @external_link, @original_filename, @mime_type, @file_size_bytes)
RETURNING *;

-- name: ListLabResultsByClient :many
SELECT * FROM lab_results WHERE client_id = @client_id ORDER BY test_date DESC, created_at DESC;

-- name: GetLabResultForNutritionist :one
SELECT lr.* FROM lab_results lr
JOIN users u ON u.id = lr.client_id AND u.nutritionist_id = @nutritionist_id
WHERE lr.id = @id AND lr.client_id = @client_id;

-- name: ListLabResultsForNutritionist :many
SELECT lr.* FROM lab_results lr
JOIN users u ON u.id = lr.client_id AND u.nutritionist_id = @nutritionist_id
WHERE lr.client_id = @client_id
ORDER BY lr.test_date DESC, lr.created_at DESC;
