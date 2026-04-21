CREATE TABLE notification_preferences (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    meal_reminders     BOOLEAN NOT NULL DEFAULT true,
    water_reminders    BOOLEAN NOT NULL DEFAULT true,
    message_alerts     BOOLEAN NOT NULL DEFAULT true,
    diet_updates       BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id)
);
