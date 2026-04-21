package entity

import "github.com/google/uuid"

// NotificationPreference holds a user's notification preference settings.
type NotificationPreference struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	MealReminders  bool
	WaterReminders bool
	MessageAlerts  bool
	DietUpdates    bool
}
