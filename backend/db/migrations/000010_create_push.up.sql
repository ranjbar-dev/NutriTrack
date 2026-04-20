-- Push subscriptions (one per client device)
CREATE TABLE push_subscriptions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    endpoint    TEXT NOT NULL,
    p256dh_key  TEXT NOT NULL,
    auth_key    TEXT NOT NULL,
    user_agent  TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (client_id, endpoint)
);
CREATE INDEX idx_push_subscriptions_client ON push_subscriptions(client_id);

-- Per-client notification type preferences
CREATE TABLE notification_preferences (
    client_id               UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    new_message             BOOLEAN NOT NULL DEFAULT TRUE,
    plan_activated          BOOLEAN NOT NULL DEFAULT TRUE,
    food_request_decision   BOOLEAN NOT NULL DEFAULT TRUE,
    meal_reminder           BOOLEAN NOT NULL DEFAULT TRUE,
    medication_reminder     BOOLEAN NOT NULL DEFAULT TRUE,
    water_reminder          BOOLEAN NOT NULL DEFAULT FALSE,
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Reminder dedup (D-19): prevents duplicate reminder pushes
CREATE TABLE sent_reminders (
    id          BIGSERIAL PRIMARY KEY,
    client_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    dedup_key   TEXT NOT NULL,
    sent_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (client_id, dedup_key)
);
CREATE INDEX idx_sent_reminders_client ON sent_reminders(client_id);
