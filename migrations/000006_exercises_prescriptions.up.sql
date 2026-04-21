CREATE TABLE exercise_recommendations (
    id                     uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    day_id                 uuid        NOT NULL REFERENCES diet_plan_days(id) ON DELETE CASCADE,
    exercise_name          text        NOT NULL,
    duration_minutes       integer     NOT NULL DEFAULT 0,
    description            text        NOT NULL DEFAULT '',
    calories_burn_estimate integer     NOT NULL DEFAULT 0,
    created_at             timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE day_prescribed_medications (
    id            uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    day_id        uuid        NOT NULL REFERENCES diet_plan_days(id) ON DELETE CASCADE,
    medication_id uuid        NOT NULL REFERENCES medications(id) ON DELETE RESTRICT,
    dosage        text        NOT NULL DEFAULT '',
    frequency     text        NOT NULL DEFAULT '',
    times         text[]      NOT NULL DEFAULT '{}',
    instructions  text        NOT NULL DEFAULT '',
    start_date    date,
    end_date      date,
    created_at    timestamptz NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_exercise_recommendations_day_id ON exercise_recommendations(day_id);
CREATE INDEX idx_day_prescribed_medications_day_id ON day_prescribed_medications(day_id);
