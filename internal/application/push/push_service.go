package push

import (
	"context"
	"encoding/json"

	webpush "github.com/SherClockHolmes/webpush-go"
	"github.com/google/uuid"
	"github.com/ranjbar-dev/nutritrack/internal/domain/push/entity"
	pushRepo "github.com/ranjbar-dev/nutritrack/internal/domain/push/repository"
)

// PushService manages Web Push subscriptions and delivers notifications.
type PushService struct {
	repo            pushRepo.PushSubscriptionRepository
	vapidPublicKey  string
	vapidPrivateKey string
}

// NewPushService creates a new PushService.
func NewPushService(repo pushRepo.PushSubscriptionRepository, vapidPublicKey, vapidPrivateKey string) *PushService {
	return &PushService{
		repo:            repo,
		vapidPublicKey:  vapidPublicKey,
		vapidPrivateKey: vapidPrivateKey,
	}
}

// Subscribe registers or updates a push subscription for a user.
func (s *PushService) Subscribe(ctx context.Context, userID uuid.UUID, endpoint, p256dh, auth string) (*entity.PushSubscription, error) {
	sub := &entity.PushSubscription{
		UserID:   userID,
		Endpoint: endpoint,
		P256dh:   p256dh,
		Auth:     auth,
	}
	return s.repo.Upsert(ctx, sub)
}

// Unsubscribe removes a push subscription.
func (s *PushService) Unsubscribe(ctx context.Context, userID uuid.UUID, endpoint string) error {
	return s.repo.Delete(ctx, userID, endpoint)
}

// Send delivers a push notification to all active subscriptions of a user.
// Returns nil immediately if VAPID keys are not configured.
// Individual subscription failures are silently ignored (best-effort delivery).
func (s *PushService) Send(ctx context.Context, userID uuid.UUID, title, body string) error {
	if s.vapidPublicKey == "" || s.vapidPrivateKey == "" {
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
		resp, sendErr := webpush.SendNotification(payload, &webpush.Subscription{
			Endpoint: sub.Endpoint,
			Keys: webpush.Keys{
				Auth:   sub.Auth,
				P256dh: sub.P256dh,
			},
		}, &webpush.Options{
			Subscriber:      "mailto:info@nutritrack.ir",
			VAPIDPublicKey:  s.vapidPublicKey,
			VAPIDPrivateKey: s.vapidPrivateKey,
			TTL:             30,
		})
		if sendErr == nil && resp != nil {
			resp.Body.Close()
		}
	}

	return nil
}
