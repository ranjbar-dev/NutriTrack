package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranjbar-dev/nutritrack/internal/domain/message/entity"
)

// MessageRepository defines persistence operations for messages.
type MessageRepository interface {
	Create(ctx context.Context, msg *entity.Message) (*entity.Message, error)
	ListConversation(ctx context.Context, userA, userB uuid.UUID, limit, offset int32) ([]*entity.Message, int64, error)
	MarkRead(ctx context.Context, receiverID, senderID uuid.UUID) error
	CountUnread(ctx context.Context, userID uuid.UUID) (int64, error)
}
