package db

import (
	"context"

	"github.com/google/uuid"
)

const upsertNotificationPreferences = `-- name: UpsertNotificationPreferences :one
INSERT INTO notification_preferences (user_id, meal_reminders, water_reminders, message_alerts, diet_updates)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (user_id) DO UPDATE SET
    meal_reminders  = EXCLUDED.meal_reminders,
    water_reminders = EXCLUDED.water_reminders,
    message_alerts  = EXCLUDED.message_alerts,
    diet_updates    = EXCLUDED.diet_updates,
    updated_at      = NOW()
RETURNING id, user_id, meal_reminders, water_reminders, message_alerts, diet_updates, created_at, updated_at`

// UpsertNotificationPreferencesParams holds parameters for upserting notification preferences.
type UpsertNotificationPreferencesParams struct {
	UserID         uuid.UUID `db:"user_id"`
	MealReminders  bool      `db:"meal_reminders"`
	WaterReminders bool      `db:"water_reminders"`
	MessageAlerts  bool      `db:"message_alerts"`
	DietUpdates    bool      `db:"diet_updates"`
}

// UpsertNotificationPreferences inserts or updates notification preferences for a user.
func (q *Queries) UpsertNotificationPreferences(ctx context.Context, arg UpsertNotificationPreferencesParams) (NotificationPreference, error) {
	row := q.db.QueryRow(ctx, upsertNotificationPreferences,
		arg.UserID,
		arg.MealReminders,
		arg.WaterReminders,
		arg.MessageAlerts,
		arg.DietUpdates,
	)
	var i NotificationPreference
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.MealReminders,
		&i.WaterReminders,
		&i.MessageAlerts,
		&i.DietUpdates,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getNotificationPreferences = `-- name: GetNotificationPreferences :one
SELECT id, user_id, meal_reminders, water_reminders, message_alerts, diet_updates, created_at, updated_at
FROM notification_preferences WHERE user_id = $1`

// GetNotificationPreferences returns the notification preferences for a user.
func (q *Queries) GetNotificationPreferences(ctx context.Context, userID uuid.UUID) (NotificationPreference, error) {
	row := q.db.QueryRow(ctx, getNotificationPreferences, userID)
	var i NotificationPreference
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.MealReminders,
		&i.WaterReminders,
		&i.MessageAlerts,
		&i.DietUpdates,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
