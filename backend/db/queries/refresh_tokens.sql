-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
    user_id, token_hash, family_id, expires_at
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetRefreshTokenByHash :one
SELECT * FROM refresh_tokens
WHERE token_hash = $1 AND revoked = false;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked = true
WHERE id = $1;

-- name: RevokeTokenFamily :exec
UPDATE refresh_tokens
SET revoked = true
WHERE family_id = $1;

-- name: RevokeUserTokens :exec
UPDATE refresh_tokens
SET revoked = true
WHERE user_id = $1;
