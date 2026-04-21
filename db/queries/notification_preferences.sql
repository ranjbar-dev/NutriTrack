-- name: UpsertNotificationPreferences :one
INSERT INTO notification_preferences (user_id, meal_reminders, water_reminders, message_alerts, diet_updates)
VALUES (@user_id, @meal_reminders, @water_reminders, @message_alerts, @diet_updates)
ON CONFLICT (user_id) DO UPDATE SET
    meal_reminders  = EXCLUDED.meal_reminders,
    water_reminders = EXCLUDED.water_reminders,
    message_alerts  = EXCLUDED.message_alerts,
    diet_updates    = EXCLUDED.diet_updates,
    updated_at      = NOW()
RETURNING *;

-- name: GetNotificationPreferences :one
SELECT * FROM notification_preferences WHERE user_id = @user_id;
