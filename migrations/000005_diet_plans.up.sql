CREATE TYPE diet_plan_status AS ENUM ('active', 'archived', 'draft');

CREATE TABLE diet_plans (
    id                    uuid             PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id             uuid             NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    nutritionist_id       uuid             NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title                 text             NOT NULL DEFAULT '',
    start_date            date             NOT NULL,
    end_date              date             NOT NULL,
    notes                 text             NOT NULL DEFAULT '',
    daily_water_target_ml integer          NOT NULL DEFAULT 0,
    status                diet_plan_status NOT NULL DEFAULT 'active',
    created_at            timestamptz      NOT NULL DEFAULT NOW(),
    updated_at            timestamptz      NOT NULL DEFAULT NOW()
);

CREATE TABLE diet_plan_days (
    id         uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id    uuid        NOT NULL REFERENCES diet_plans(id) ON DELETE CASCADE,
    day_number integer     NOT NULL CHECK (day_number >= 1),
    created_at timestamptz NOT NULL DEFAULT NOW(),
    UNIQUE(plan_id, day_number)
);

CREATE TABLE diet_meals (
    id             uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    day_id         uuid        NOT NULL REFERENCES diet_plan_days(id) ON DELETE CASCADE,
    title          text        NOT NULL,
    scheduled_time time        NOT NULL DEFAULT '08:00:00',
    display_order  integer     NOT NULL DEFAULT 0,
    created_at     timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE meal_options (
    id            uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    meal_id       uuid        NOT NULL REFERENCES diet_meals(id) ON DELETE CASCADE,
    option_number integer     NOT NULL CHECK (option_number >= 1),
    created_at    timestamptz NOT NULL DEFAULT NOW(),
    UNIQUE(meal_id, option_number)
);

CREATE TABLE meal_option_items (
    id          uuid         PRIMARY KEY DEFAULT gen_random_uuid(),
    option_id   uuid         NOT NULL REFERENCES meal_options(id) ON DELETE CASCADE,
    food_id     uuid         NOT NULL REFERENCES foods(id) ON DELETE RESTRICT,
    quantity    numeric(8,2) NOT NULL DEFAULT 1,
    unit        text         NOT NULL DEFAULT 'گرم',
    notes       text         NOT NULL DEFAULT '',
    created_at  timestamptz  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_diet_plans_client_id        ON diet_plans(client_id);
CREATE INDEX idx_diet_plans_nutritionist_id  ON diet_plans(nutritionist_id);
CREATE INDEX idx_diet_plans_status           ON diet_plans(status);
CREATE INDEX idx_diet_plan_days_plan_id      ON diet_plan_days(plan_id);
CREATE INDEX idx_diet_meals_day_id           ON diet_meals(day_id);
CREATE INDEX idx_meal_options_meal_id        ON meal_options(meal_id);
CREATE INDEX idx_meal_option_items_option_id ON meal_option_items(option_id);
