package entity

import (
	"time"

	"github.com/google/uuid"
)

// PushSubscription represents a Web Push subscription for a user.
type PushSubscription struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Endpoint  string
	P256dh    string
	Auth      string
	CreatedAt time.Time
}
