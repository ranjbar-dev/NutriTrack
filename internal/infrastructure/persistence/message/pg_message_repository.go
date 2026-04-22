package message

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ranjbar-dev/nutritrack/internal/domain/message/entity"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	db "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/sqlc"
)

// PgMessageRepository is the PostgreSQL implementation of MessageRepository.
type PgMessageRepository struct {
	queries *db.Queries
}

// NewPgMessageRepository creates a new PgMessageRepository.
func NewPgMessageRepository(pool *pgxpool.Pool) *PgMessageRepository {
	return &PgMessageRepository{queries: db.New(pool)}
}

// Create inserts a new message.
func (r *PgMessageRepository) Create(ctx context.Context, msg *entity.Message) (*entity.Message, error) {
	row, err := r.queries.CreateMessage(ctx, db.CreateMessageParams{
		SenderID:       msg.SenderID(),
		ReceiverID:     msg.ReceiverID(),
		Content:        msg.Content(),
		AttachmentPath: msg.AttachmentPath(),
		AttachmentType: msg.AttachmentType(),
		AttachmentSize: msg.AttachmentSize(),
		AttachmentName: msg.AttachmentName(),
	})
	if err != nil {
		return nil, shared.ErrInternal
	}
	return toDomain(row), nil
}

// ListConversation returns paginated messages between two users, plus total count.
func (r *PgMessageRepository) ListConversation(ctx context.Context, userA, userB uuid.UUID, limit, offset int32) ([]*entity.Message, int64, error) {
	rows, err := r.queries.ListConversationMessages(ctx, db.ListConversationMessagesParams{
		UserA:  userA,
		UserB:  userB,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, 0, shared.ErrInternal
	}
	total, err := r.queries.CountConversationMessages(ctx, userA, userB)
	if err != nil {
		return nil, 0, shared.ErrInternal
	}
	msgs := make([]*entity.Message, len(rows))
	for i, row := range rows {
		msgs[i] = toDomain(row)
	}
	return msgs, total, nil
}

// MarkRead marks all messages from senderID to receiverID as read.
func (r *PgMessageRepository) MarkRead(ctx context.Context, receiverID, senderID uuid.UUID) error {
	if err := r.queries.MarkConversationRead(ctx, receiverID, senderID); err != nil {
		return shared.ErrInternal
	}
	return nil
}

// CountUnread returns the number of unread messages for a user.
func (r *PgMessageRepository) CountUnread(ctx context.Context, userID uuid.UUID) (int64, error) {
	count, err := r.queries.CountUnreadMessages(ctx, userID)
	if err != nil {
		return 0, shared.ErrInternal
	}
	return count, nil
}
