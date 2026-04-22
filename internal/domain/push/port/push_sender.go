// Package port defines the outbound ports (driven adapters) for the push domain.
package port

import "context"

// PushSender is the domain port for delivering Web Push notifications.
// Implementations live in internal/infrastructure/push/.
type PushSender interface {
	// SendToSubscription delivers a raw JSON payload to a single push subscription.
	// Returns nil on success. The implementation is responsible for transport details
	// (VAPID signing, HTTP delivery, TTL, etc.).
	SendToSubscription(ctx context.Context, endpoint, p256dh, auth string, payload []byte) error
}
