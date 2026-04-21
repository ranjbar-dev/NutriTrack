package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranjbar-dev/nutritrack/internal/domain/push/entity"
)

// PushSubscriptionRepository handles persistence of push subscriptions.
type PushSubscriptionRepository interface {
	Upsert(ctx context.Context, sub *entity.PushSubscription) (*entity.PushSubscription, error)
	Delete(ctx context.Context, userID uuid.UUID, endpoint string) error
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*entity.PushSubscription, error)
}
