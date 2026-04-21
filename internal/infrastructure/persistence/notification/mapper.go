package notification

import (
	"github.com/ranjbar-dev/nutritrack/internal/domain/notification/entity"
	db "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/sqlc"
)

func toDomain(row db.NotificationPreference) entity.NotificationPreference {
	return entity.NotificationPreference{
		ID:             row.ID,
		UserID:         row.UserID,
		MealReminders:  row.MealReminders,
		WaterReminders: row.WaterReminders,
		MessageAlerts:  row.MessageAlerts,
		DietUpdates:    row.DietUpdates,
	}
}
