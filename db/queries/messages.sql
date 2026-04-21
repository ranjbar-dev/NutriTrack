-- name: CreateMessage :one
INSERT INTO messages (sender_id, receiver_id, content, attachment_path, attachment_type, attachment_size, attachment_name)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, sender_id, receiver_id, content, attachment_path, attachment_type, attachment_size, attachment_name, read_at, created_at;

-- name: ListConversationMessages :many
SELECT id, sender_id, receiver_id, content, attachment_path, attachment_type, attachment_size, attachment_name, read_at, created_at
FROM messages
WHERE (sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1)
ORDER BY created_at ASC
LIMIT $3 OFFSET $4;

-- name: CountConversationMessages :one
SELECT COUNT(*) FROM messages
WHERE (sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1);

-- name: MarkConversationRead :exec
UPDATE messages
SET read_at = NOW()
WHERE receiver_id = $1 AND sender_id = $2 AND read_at IS NULL;

-- name: CountUnreadMessages :one
SELECT COUNT(*) FROM messages WHERE receiver_id = $1 AND read_at IS NULL;
