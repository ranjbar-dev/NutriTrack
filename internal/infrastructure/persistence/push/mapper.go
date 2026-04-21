package push

import (
	"github.com/ranjbar-dev/nutritrack/internal/domain/push/entity"
	db "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/sqlc"
)

func toDomain(row db.PushSubscription) *entity.PushSubscription {
	return &entity.PushSubscription{
		ID:        row.ID,
		UserID:    row.UserID,
		Endpoint:  row.Endpoint,
		P256dh:    row.P256dh,
		Auth:      row.Auth,
		CreatedAt: row.CreatedAt,
	}
}
