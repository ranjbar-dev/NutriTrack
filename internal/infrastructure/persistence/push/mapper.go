package push

import (
	"github.com/ranjbar-dev/nutritrack/internal/domain/push/entity"
	db "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/sqlc"
)

func toDomain(row db.PushSubscription) *entity.PushSubscription {
	return entity.NewPushSubscriptionFromDB(row.ID, row.UserID, row.Endpoint, row.P256dh, row.Auth, row.CreatedAt)
}
