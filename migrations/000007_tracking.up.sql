-- Daily tracking: 6 types with offline-sync idempotency via UNIQUE(client_id, local_id)

-- Food logs: what the client ate
CREATE TABLE food_logs (
    id              uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    client_id       uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    local_id        text NOT NULL,
    logged_at       timestamptz NOT NULL DEFAULT NOW(),
    logged_date     date NOT NULL,
    food_id         uuid REFERENCES foods(id) ON DELETE SET NULL,
    food_name       text NOT NULL DEFAULT '',
    quantity        numeric(10,3) NOT NULL DEFAULT 0,
    unit            text NOT NULL DEFAULT '',
    calories        numeric(10,2) NOT NULL DEFAULT 0,
    protein         numeric(10,2) NOT NULL DEFAULT 0,
    carbs           numeric(10,2) NOT NULL DEFAULT 0,
    fat             numeric(10,2) NOT NULL DEFAULT 0,
    notes           text NOT NULL DEFAULT '',
    created_at      timestamptz NOT NULL DEFAULT NOW(),
    UNIQUE(client_id, local_id)
);
CREATE INDEX idx_food_logs_client_date ON food_logs(client_id, logged_date);

-- Water logs
CREATE TABLE water_logs (
    id              uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    client_id       uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    local_id        text NOT NULL,
    logged_at       timestamptz NOT NULL DEFAULT NOW(),
    logged_date     date NOT NULL,
    amount_ml       int NOT NULL DEFAULT 0,
    notes           text NOT NULL DEFAULT '',
    created_at      timestamptz NOT NULL DEFAULT NOW(),
    UNIQUE(client_id, local_id)
);
CREATE INDEX idx_water_logs_client_date ON water_logs(client_id, logged_date);

-- Sleep logs
CREATE TABLE sleep_logs (
    id              uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    client_id       uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    local_id        text NOT NULL,
    logged_date     date NOT NULL,
    sleep_start     timestamptz NOT NULL,
    sleep_end       timestamptz NOT NULL,
    duration_minutes int NOT NULL DEFAULT 0,
    quality         int NOT NULL DEFAULT 0,
    notes           text NOT NULL DEFAULT '',
    created_at      timestamptz NOT NULL DEFAULT NOW(),
    UNIQUE(client_id, local_id)
);
CREATE INDEX idx_sleep_logs_client_date ON sleep_logs(client_id, logged_date);

-- Exercise logs
CREATE TABLE exercise_logs (
    id              uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    client_id       uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    local_id        text NOT NULL,
    logged_at       timestamptz NOT NULL DEFAULT NOW(),
    logged_date     date NOT NULL,
    exercise_name   text NOT NULL DEFAULT '',
    duration_minutes int NOT NULL DEFAULT 0,
    calories_burned  int NOT NULL DEFAULT 0,
    notes           text NOT NULL DEFAULT '',
    created_at      timestamptz NOT NULL DEFAULT NOW(),
    UNIQUE(client_id, local_id)
);
CREATE INDEX idx_exercise_logs_client_date ON exercise_logs(client_id, logged_date);

-- Medication logs
CREATE TABLE medication_logs (
    id              uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    client_id       uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    local_id        text NOT NULL,
    logged_at       timestamptz NOT NULL DEFAULT NOW(),
    logged_date     date NOT NULL,
    medication_id   uuid REFERENCES medications(id) ON DELETE SET NULL,
    medication_name text NOT NULL DEFAULT '',
    dosage          text NOT NULL DEFAULT '',
    notes           text NOT NULL DEFAULT '',
    created_at      timestamptz NOT NULL DEFAULT NOW(),
    UNIQUE(client_id, local_id)
);
CREATE INDEX idx_medication_logs_client_date ON medication_logs(client_id, logged_date);

-- Body measurements
CREATE TABLE body_measurements (
    id              uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    client_id       uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    local_id        text NOT NULL,
    measured_at     timestamptz NOT NULL DEFAULT NOW(),
    measured_date   date NOT NULL,
    weight_kg       numeric(6,2),
    height_cm       numeric(5,2),
    waist_cm        numeric(5,2),
    hip_cm          numeric(5,2),
    chest_cm        numeric(5,2),
    arm_cm          numeric(5,2),
    notes           text NOT NULL DEFAULT '',
    created_at      timestamptz NOT NULL DEFAULT NOW(),
    UNIQUE(client_id, local_id)
);
CREATE INDEX idx_body_measurements_client_date ON body_measurements(client_id, measured_date);
