-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 AND is_active = true;

-- name: GetUserByMobile :one
SELECT * FROM users
WHERE mobile = $1 AND is_active = true;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (
    role, full_name, email, password_hash, mobile,
    date_of_birth, height_cm, gender, nutritionist_id, notes
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8, $9, $10
) RETURNING *;

-- name: GetClientsByNutritionistID :many
SELECT * FROM users
WHERE nutritionist_id = $1 AND role = 'client'
ORDER BY created_at DESC;

-- name: UpdateUserActive :exec
UPDATE users
SET is_active = $2, updated_at = NOW()
WHERE id = $1;
