-- name: GetNotificationPreferences :one
SELECT * FROM notification_preferences WHERE client_id = $1;

-- name: UpsertNotificationPreferences :one
INSERT INTO notification_preferences (client_id, new_message, plan_activated, food_request_decision,
  meal_reminder, medication_reminder, water_reminder)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (client_id) DO UPDATE
  SET new_message           = EXCLUDED.new_message,
      plan_activated        = EXCLUDED.plan_activated,
      food_request_decision = EXCLUDED.food_request_decision,
      meal_reminder         = EXCLUDED.meal_reminder,
      medication_reminder   = EXCLUDED.medication_reminder,
      water_reminder        = EXCLUDED.water_reminder,
      updated_at            = NOW()
RETURNING *;
