package push

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ranjbar-dev/nutritrack/internal/domain/push/entity"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	db "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/sqlc"
)

// PgPushSubscriptionRepository is the PostgreSQL implementation of PushSubscriptionRepository.
type PgPushSubscriptionRepository struct {
	q *db.Queries
}

// NewPgPushSubscriptionRepository creates a new PgPushSubscriptionRepository.
func NewPgPushSubscriptionRepository(pool *pgxpool.Pool) *PgPushSubscriptionRepository {
	return &PgPushSubscriptionRepository{q: db.New(pool)}
}

func (r *PgPushSubscriptionRepository) Upsert(ctx context.Context, sub *entity.PushSubscription) (*entity.PushSubscription, error) {
	row, err := r.q.UpsertPushSubscription(ctx, db.UpsertPushSubscriptionParams{
		UserID:   sub.GetUserID(),
		Endpoint: sub.GetEndpoint(),
		P256dh:   sub.GetP256dh(),
		Auth:     sub.GetAuth(),
	})
	if err != nil {
		return nil, shared.ErrInternal
	}
	return toDomain(row), nil
}

func (r *PgPushSubscriptionRepository) Delete(ctx context.Context, userID uuid.UUID, endpoint string) error {
	if err := r.q.DeletePushSubscription(ctx, userID, endpoint); err != nil {
		return shared.ErrInternal
	}
	return nil
}

func (r *PgPushSubscriptionRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]*entity.PushSubscription, error) {
	rows, err := r.q.ListPushSubscriptionsByUser(ctx, userID)
	if err != nil {
		return nil, shared.ErrInternal
	}
	result := make([]*entity.PushSubscription, len(rows))
	for i, row := range rows {
		result[i] = toDomain(row)
	}
	return result, nil
}
