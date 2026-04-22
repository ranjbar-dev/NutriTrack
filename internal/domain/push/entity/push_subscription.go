package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrEmptyEndpoint        = errors.New("push subscription endpoint must not be empty")
	ErrEmptyP256dh          = errors.New("push subscription p256dh key must not be empty")
	ErrEmptyAuth            = errors.New("push subscription auth key must not be empty")
	ErrSubscriptionNotFound = errors.New("push subscription not found")
)

// PushSubscription represents a Web Push subscription for a user.
type PushSubscription struct {
	id        uuid.UUID
	userID    uuid.UUID
	endpoint  string
	p256dh    string
	auth      string
	createdAt time.Time
}

// NewPushSubscription creates a validated PushSubscription for a new subscription (before persistence).
func NewPushSubscription(userID uuid.UUID, endpoint, p256dh, auth string) (*PushSubscription, error) {
	if endpoint == "" {
		return nil, ErrEmptyEndpoint
	}
	if p256dh == "" {
		return nil, ErrEmptyP256dh
	}
	if auth == "" {
		return nil, ErrEmptyAuth
	}
	return &PushSubscription{
		userID:   userID,
		endpoint: endpoint,
		p256dh:   p256dh,
		auth:     auth,
	}, nil
}

// NewPushSubscriptionFromDB reconstructs a PushSubscription from persistent storage.
// This should only be called by infrastructure code.
func NewPushSubscriptionFromDB(id, userID uuid.UUID, endpoint, p256dh, auth string, createdAt time.Time) *PushSubscription {
	return &PushSubscription{
		id:        id,
		userID:    userID,
		endpoint:  endpoint,
		p256dh:    p256dh,
		auth:      auth,
		createdAt: createdAt,
	}
}

// Getters

func (ps PushSubscription) GetID() uuid.UUID        { return ps.id }
func (ps PushSubscription) GetUserID() uuid.UUID    { return ps.userID }
func (ps PushSubscription) GetEndpoint() string     { return ps.endpoint }
func (ps PushSubscription) GetP256dh() string       { return ps.p256dh }
func (ps PushSubscription) GetAuth() string         { return ps.auth }
func (ps PushSubscription) GetCreatedAt() time.Time { return ps.createdAt }
