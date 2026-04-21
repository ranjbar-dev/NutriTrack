package db

import (
	"context"

	"github.com/google/uuid"
)

const createMessage = `-- name: CreateMessage :one
INSERT INTO messages (sender_id, receiver_id, content, attachment_path, attachment_type, attachment_size, attachment_name)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, sender_id, receiver_id, content, attachment_path, attachment_type, attachment_size, attachment_name, read_at, created_at`

// CreateMessageParams holds parameters for creating a message.
type CreateMessageParams struct {
	SenderID       uuid.UUID `db:"sender_id"`
	ReceiverID     uuid.UUID `db:"receiver_id"`
	Content        string    `db:"content"`
	AttachmentPath *string   `db:"attachment_path"`
	AttachmentType *string   `db:"attachment_type"`
	AttachmentSize *int64    `db:"attachment_size"`
	AttachmentName *string   `db:"attachment_name"`
}

// CreateMessage inserts a new message and returns the created row.
func (q *Queries) CreateMessage(ctx context.Context, arg CreateMessageParams) (Message, error) {
	row := q.db.QueryRow(ctx, createMessage,
		arg.SenderID, arg.ReceiverID, arg.Content,
		arg.AttachmentPath, arg.AttachmentType, arg.AttachmentSize, arg.AttachmentName,
	)
	var i Message
	err := row.Scan(
		&i.ID, &i.SenderID, &i.ReceiverID, &i.Content,
		&i.AttachmentPath, &i.AttachmentType, &i.AttachmentSize, &i.AttachmentName,
		&i.ReadAt, &i.CreatedAt,
	)
	return i, err
}

// ListConversationMessagesParams holds parameters for listing conversation messages.
type ListConversationMessagesParams struct {
	UserA  uuid.UUID `db:"user_a"`
	UserB  uuid.UUID `db:"user_b"`
	Limit  int32     `db:"limit"`
	Offset int32     `db:"offset"`
}

const listConversationMessages = `-- name: ListConversationMessages :many
SELECT id, sender_id, receiver_id, content, attachment_path, attachment_type, attachment_size, attachment_name, read_at, created_at
FROM messages
WHERE (sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1)
ORDER BY created_at ASC
LIMIT $3 OFFSET $4`

// ListConversationMessages returns paginated messages between two users.
func (q *Queries) ListConversationMessages(ctx context.Context, arg ListConversationMessagesParams) ([]Message, error) {
	rows, err := q.db.Query(ctx, listConversationMessages, arg.UserA, arg.UserB, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Message
	for rows.Next() {
		var i Message
		if err := rows.Scan(
			&i.ID, &i.SenderID, &i.ReceiverID, &i.Content,
			&i.AttachmentPath, &i.AttachmentType, &i.AttachmentSize, &i.AttachmentName,
			&i.ReadAt, &i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

const countConversationMessages = `-- name: CountConversationMessages :one
SELECT COUNT(*) FROM messages
WHERE (sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1)`

// CountConversationMessages counts total messages between two users.
func (q *Queries) CountConversationMessages(ctx context.Context, userA, userB uuid.UUID) (int64, error) {
	var count int64
	err := q.db.QueryRow(ctx, countConversationMessages, userA, userB).Scan(&count)
	return count, err
}

const markConversationRead = `-- name: MarkConversationRead :exec
UPDATE messages SET read_at = NOW()
WHERE receiver_id = $1 AND sender_id = $2 AND read_at IS NULL`

// MarkConversationRead marks all messages from senderID to receiverID as read.
func (q *Queries) MarkConversationRead(ctx context.Context, receiverID, senderID uuid.UUID) error {
	_, err := q.db.Exec(ctx, markConversationRead, receiverID, senderID)
	return err
}

const countUnreadMessages = `-- name: CountUnreadMessages :one
SELECT COUNT(*) FROM messages WHERE receiver_id = $1 AND read_at IS NULL`

// CountUnreadMessages returns the total number of unread messages for a user.
func (q *Queries) CountUnreadMessages(ctx context.Context, receiverID uuid.UUID) (int64, error) {
	var count int64
	err := q.db.QueryRow(ctx, countUnreadMessages, receiverID).Scan(&count)
	return count, err
}
