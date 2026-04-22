package notification

import (
	"github.com/ranjbar-dev/nutritrack/internal/domain/notification/entity"
	db "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/sqlc"
)

func toDomain(row db.NotificationPreference) entity.NotificationPreference {
	pref, _ := entity.NewNotificationPreference(row.ID, row.UserID)
	pref.SetMealReminders(row.MealReminders)
	pref.SetWaterReminders(row.WaterReminders)
	pref.SetMessageAlerts(row.MessageAlerts)
	pref.SetDietUpdates(row.DietUpdates)
	return *pref
}
