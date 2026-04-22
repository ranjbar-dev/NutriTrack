// Package push provides the infrastructure implementation of the PushSender domain port.
package push

import (
	"context"

	webpush "github.com/SherClockHolmes/webpush-go"
	pushPort "github.com/ranjbar-dev/nutritrack/internal/domain/push/port"
)

// Ensure WebpushSender satisfies the domain port at compile time.
var _ pushPort.PushSender = (*WebpushSender)(nil)

// WebpushSender delivers Web Push notifications using the webpush-go library.
type WebpushSender struct {
	vapidPublicKey  string
	vapidPrivateKey string
	subscriberEmail string
}

// NewWebpushSender creates a WebpushSender. VAPID keys and subscriber email must be non-empty.
func NewWebpushSender(vapidPublicKey, vapidPrivateKey, subscriberEmail string) *WebpushSender {
	return &WebpushSender{
		vapidPublicKey:  vapidPublicKey,
		vapidPrivateKey: vapidPrivateKey,
		subscriberEmail: subscriberEmail,
	}
}

// SendToSubscription delivers payload to the given Web Push subscription.
// A nil error is returned on success; delivery errors are returned to the caller.
func (s *WebpushSender) SendToSubscription(_ context.Context, endpoint, p256dh, auth string, payload []byte) error {
	resp, err := webpush.SendNotification(payload, &webpush.Subscription{
		Endpoint: endpoint,
		Keys: webpush.Keys{
			Auth:   auth,
			P256dh: p256dh,
		},
	}, &webpush.Options{
		Subscriber:      s.subscriberEmail,
		VAPIDPublicKey:  s.vapidPublicKey,
		VAPIDPrivateKey: s.vapidPrivateKey,
		TTL:             30,
	})
	if err != nil {
		return err
	}
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}
	return nil
}
