-- name: CreateOTP :one
INSERT INTO otp_codes (
    mobile, code_hash, expires_at, max_attempts
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetActiveOTPByMobile :one
SELECT * FROM otp_codes
WHERE mobile = $1 AND expires_at > NOW() AND verified = false
ORDER BY created_at DESC
LIMIT 1;

-- name: IncrementOTPAttempts :exec
UPDATE otp_codes
SET attempts = attempts + 1
WHERE id = $1;

-- name: MarkOTPVerified :exec
UPDATE otp_codes
SET verified = true
WHERE id = $1;

-- name: DeleteExpiredOTPs :exec
DELETE FROM otp_codes
WHERE expires_at < NOW();
