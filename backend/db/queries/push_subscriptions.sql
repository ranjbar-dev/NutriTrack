-- name: UpsertPushSubscription :one
INSERT INTO push_subscriptions (client_id, endpoint, p256dh_key, auth_key, user_agent)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (client_id, endpoint) DO UPDATE
  SET p256dh_key = EXCLUDED.p256dh_key,
      auth_key   = EXCLUDED.auth_key,
      user_agent = EXCLUDED.user_agent,
      updated_at = NOW()
RETURNING *;

-- name: DeletePushSubscriptionByEndpoint :exec
DELETE FROM push_subscriptions
WHERE client_id = $1 AND endpoint = $2;

-- name: GetPushSubscriptionsByClient :many
SELECT * FROM push_subscriptions WHERE client_id = $1;

-- name: InsertSentReminder :exec
INSERT INTO sent_reminders (client_id, dedup_key)
VALUES ($1, $2)
ON CONFLICT (client_id, dedup_key) DO NOTHING;

-- name: ReminderAlreadySent :one
SELECT EXISTS (
  SELECT 1 FROM sent_reminders WHERE client_id = $1 AND dedup_key = $2
) AS "exists";

-- name: PurgeSentRemindersOlderThan :exec
DELETE FROM sent_reminders WHERE sent_at < $1;
