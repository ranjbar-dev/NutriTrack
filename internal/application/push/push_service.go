package push

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/ranjbar-dev/nutritrack/internal/domain/push/entity"
	pushPort "github.com/ranjbar-dev/nutritrack/internal/domain/push/port"
	pushRepo "github.com/ranjbar-dev/nutritrack/internal/domain/push/repository"
)

// PushService manages Web Push subscriptions and delivers notifications.
type PushService struct {
	repo   pushRepo.PushSubscriptionRepository
	sender pushPort.PushSender
}

// NewPushService creates a new PushService.
func NewPushService(repo pushRepo.PushSubscriptionRepository, sender pushPort.PushSender) *PushService {
	return &PushService{
		repo:   repo,
		sender: sender,
	}
}

// Subscribe registers or updates a push subscription for a user.
func (s *PushService) Subscribe(ctx context.Context, userID uuid.UUID, endpoint, p256dh, auth string) (*entity.PushSubscription, error) {
	sub, err := entity.NewPushSubscription(userID, endpoint, p256dh, auth)
	if err != nil {
		return nil, err
	}
	return s.repo.Upsert(ctx, sub)
}

// Unsubscribe removes a push subscription.
func (s *PushService) Unsubscribe(ctx context.Context, userID uuid.UUID, endpoint string) error {
	return s.repo.Delete(ctx, userID, endpoint)
}

// Send delivers a push notification to all active subscriptions of a user.
// Individual subscription failures are silently ignored (best-effort delivery).
func (s *PushService) Send(ctx context.Context, userID uuid.UUID, title, body string) error {
	if s.sender == nil {
		return nil
	}

	subs, err := s.repo.ListByUser(ctx, userID)
	if err != nil {
		return err
	}

	payload, _ := json.Marshal(map[string]any{
		"title": title,
		"body":  body,
	})

	for _, sub := range subs {
		_ = s.sender.SendToSubscription(ctx, sub.GetEndpoint(), sub.GetP256dh(), sub.GetAuth(), payload)
	}

	return nil
}
