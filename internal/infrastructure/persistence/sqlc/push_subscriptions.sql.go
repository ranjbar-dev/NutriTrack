package db

import (
	"context"

	"github.com/google/uuid"
)

const upsertPushSubscription = `-- name: UpsertPushSubscription :one
INSERT INTO push_subscriptions (user_id, endpoint, p256dh, auth)
VALUES ($1, $2, $3, $4)
ON CONFLICT (user_id, endpoint) DO UPDATE
    SET p256dh = EXCLUDED.p256dh, auth = EXCLUDED.auth
RETURNING id, user_id, endpoint, p256dh, auth, created_at`

// UpsertPushSubscriptionParams holds parameters for upserting a push subscription.
type UpsertPushSubscriptionParams struct {
	UserID   uuid.UUID `db:"user_id"`
	Endpoint string    `db:"endpoint"`
	P256dh   string    `db:"p256dh"`
	Auth     string    `db:"auth"`
}

// UpsertPushSubscription upserts a push subscription for a user.
func (q *Queries) UpsertPushSubscription(ctx context.Context, arg UpsertPushSubscriptionParams) (PushSubscription, error) {
	row := q.db.QueryRow(ctx, upsertPushSubscription, arg.UserID, arg.Endpoint, arg.P256dh, arg.Auth)
	var i PushSubscription
	err := row.Scan(&i.ID, &i.UserID, &i.Endpoint, &i.P256dh, &i.Auth, &i.CreatedAt)
	return i, err
}

const deletePushSubscription = `-- name: DeletePushSubscription :exec
DELETE FROM push_subscriptions WHERE user_id = $1 AND endpoint = $2`

// DeletePushSubscription removes a push subscription.
func (q *Queries) DeletePushSubscription(ctx context.Context, userID uuid.UUID, endpoint string) error {
	_, err := q.db.Exec(ctx, deletePushSubscription, userID, endpoint)
	return err
}

const listPushSubscriptionsByUser = `-- name: ListPushSubscriptionsByUser :many
SELECT id, user_id, endpoint, p256dh, auth, created_at
FROM push_subscriptions
WHERE user_id = $1`

// ListPushSubscriptionsByUser returns all push subscriptions for a user.
func (q *Queries) ListPushSubscriptionsByUser(ctx context.Context, userID uuid.UUID) ([]PushSubscription, error) {
	rows, err := q.db.Query(ctx, listPushSubscriptionsByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PushSubscription
	for rows.Next() {
		var i PushSubscription
		if err := rows.Scan(&i.ID, &i.UserID, &i.Endpoint, &i.P256dh, &i.Auth, &i.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}
