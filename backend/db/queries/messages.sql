-- name: CreateMessage :one
INSERT INTO messages (sender_id, receiver_id, content, attachment_type, attachment_path, attachment_name)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, sender_id, receiver_id, content, attachment_type, attachment_path, attachment_name, sent_at, read_at;

-- name: ListMessages :many
SELECT id, sender_id, receiver_id, content, attachment_type, attachment_path, attachment_name, sent_at, read_at
FROM messages
WHERE (sender_id = $1 AND receiver_id = $2)
   OR (sender_id = $2 AND receiver_id = $1)
ORDER BY sent_at ASC
LIMIT $3 OFFSET $4;

-- name: ListMessagesSince :many
SELECT id, sender_id, receiver_id, content, attachment_type, attachment_path, attachment_name, sent_at, read_at
FROM messages
WHERE ((sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1))
  AND sent_at > $3
ORDER BY sent_at ASC;

-- name: MarkMessagesRead :exec
UPDATE messages
SET read_at = NOW()
WHERE receiver_id = $1
  AND sender_id = $2
  AND read_at IS NULL;

-- name: CountUnreadMessages :one
SELECT COUNT(*)::integer
FROM messages
WHERE receiver_id = $1
  AND read_at IS NULL;

-- name: GetMessageByID :one
SELECT id, sender_id, receiver_id, content, attachment_type, attachment_path, attachment_name, sent_at, read_at
FROM messages
WHERE id = $1;
