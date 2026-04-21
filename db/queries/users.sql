-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: GetUserByMobile :one
SELECT * FROM users WHERE mobile = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
    role, mobile, email, password_hash,
    first_name, last_name, gender, birth_date,
    height, weight, avatar_url, is_active, nutritionist_id
) VALUES (
    $1, $2, $3, $4,
    $5, $6, $7, $8,
    $9, $10, $11, $12, $13
) RETURNING *;

-- name: UpdateUser :one
UPDATE users SET
    first_name      = $2,
    last_name       = $3,
    gender          = $4,
    birth_date      = $5,
    height          = $6,
    weight          = $7,
    avatar_url      = $8,
    is_active       = $9,
    password_hash   = $10,
    updated_at      = NOW()
WHERE id = $1
RETURNING *;

-- name: ListClientsByNutritionist :many
SELECT * FROM users
WHERE nutritionist_id = $1 AND role = 'client'
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountClientsByNutritionist :one
SELECT COUNT(*) FROM users
WHERE nutritionist_id = $1 AND role = 'client';

-- name: ListNutritionists :many
SELECT * FROM users
WHERE role = 'nutritionist'
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountNutritionists :one
SELECT COUNT(*) FROM users WHERE role = 'nutritionist';

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: ExistsByMobile :one
SELECT EXISTS(SELECT 1 FROM users WHERE mobile = $1)::boolean;

-- name: ExistsByEmail :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)::boolean;
