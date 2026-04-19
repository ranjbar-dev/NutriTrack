-- Migration 000007: Create Diet Plan Engine tables
-- Creates: diet_plan_status enum, diet_plans, plan_days, meals,
--          meal_options, meal_option_items, plan_exercises, plan_medications

-- New enum for diet plan lifecycle status
CREATE TYPE diet_plan_status AS ENUM ('draft', 'active', 'archived');

-- Root aggregate: one plan per client-nutritionist pair, has a lifecycle status
CREATE TABLE diet_plans (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id             UUID NOT NULL REFERENCES users(id),
    nutritionist_id       UUID NOT NULL REFERENCES users(id),
    start_date            DATE NOT NULL,
    end_date              DATE NOT NULL,
    notes                 TEXT,
    daily_water_target_ml INTEGER,
    status                diet_plan_status NOT NULL DEFAULT 'draft',
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- D-02: Only one active plan per client is allowed at the DB level
CREATE UNIQUE INDEX idx_diet_plans_one_active_per_client
    ON diet_plans (client_id) WHERE status = 'active';

-- D-37: Performance indexes
CREATE INDEX idx_diet_plans_client_id_status   ON diet_plans (client_id, status);
CREATE INDEX idx_diet_plans_nutritionist_id    ON diet_plans (nutritionist_id);

-- Plan days: numbered 1-based, unique per plan
CREATE TABLE plan_days (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id    UUID NOT NULL REFERENCES diet_plans(id) ON DELETE CASCADE,
    day_number INTEGER NOT NULL CHECK (day_number >= 1),
    label      VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (plan_id, day_number)
);
CREATE INDEX idx_plan_days_plan_id ON plan_days (plan_id);

-- Meals per day, ordered by display_order then scheduled_time
CREATE TABLE meals (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    day_id        UUID NOT NULL REFERENCES plan_days(id) ON DELETE CASCADE,
    title         VARCHAR(200) NOT NULL,
    scheduled_time TIME,
    display_order INTEGER NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_meals_day_id ON meals (day_id);

-- Meal options: client picks one option per meal
CREATE TABLE meal_options (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    meal_id       UUID NOT NULL REFERENCES meals(id) ON DELETE CASCADE,
    option_number SMALLINT NOT NULL DEFAULT 1,
    label         VARCHAR(100),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_meal_options_meal_id ON meal_options (meal_id);

-- Meal option items: links food items to a meal option with quantity
-- Uses existing measurement_unit_type enum (renamed from measurement_unit in migration 000006)
CREATE TABLE meal_option_items (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    option_id        UUID NOT NULL REFERENCES meal_options(id) ON DELETE CASCADE,
    food_id          UUID NOT NULL REFERENCES foods(id),
    quantity         DECIMAL(8,2) NOT NULL CHECK (quantity > 0),
    measurement_unit measurement_unit_type NOT NULL DEFAULT 'gram',
    notes            TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_meal_option_items_option_id ON meal_option_items (option_id);

-- Exercise recommendations per plan day
CREATE TABLE plan_exercises (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    day_id                UUID NOT NULL REFERENCES plan_days(id) ON DELETE CASCADE,
    exercise_name         VARCHAR(200) NOT NULL,
    duration_minutes      INTEGER NOT NULL,
    description           TEXT,
    calories_burn_estimate INTEGER,
    display_order         INTEGER NOT NULL DEFAULT 0,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_plan_exercises_day_id ON plan_exercises (day_id);

-- Prescribed medications for a diet plan
CREATE TABLE plan_medications (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id       UUID NOT NULL REFERENCES diet_plans(id) ON DELETE CASCADE,
    medication_id UUID NOT NULL REFERENCES medications(id),
    dosage        VARCHAR(100) NOT NULL,
    frequency     VARCHAR(200) NOT NULL,
    times         JSONB NOT NULL DEFAULT '[]',
    instructions  TEXT,
    start_date    DATE,
    end_date      DATE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_plan_medications_plan_id ON plan_medications (plan_id);
