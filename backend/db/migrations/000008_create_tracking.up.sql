-- Migration 000008: Create Client Tracking Suite tables
-- Creates: sleep_quality enum, lab_result_type enum, food_logs, water_logs,
--          sleep_logs, exercise_logs, medication_logs, body_measurements,
--          lab_results

CREATE TYPE sleep_quality AS ENUM ('good', 'fair', 'poor');

CREATE TYPE lab_result_type AS ENUM (
    'blood_test',
    'urine_test',
    'thyroid',
    'hormone',
    'allergy',
    'other'
);

CREATE TABLE food_logs (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id          UUID NOT NULL REFERENCES users(id),
    local_id           UUID NOT NULL,
    date               DATE NOT NULL,
    meal_id            UUID NOT NULL REFERENCES meals(id),
    selected_option_id UUID REFERENCES meal_options(id),
    is_skipped         BOOLEAN NOT NULL DEFAULT false,
    notes              TEXT,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_food_logs_local_id UNIQUE (local_id),
    CONSTRAINT uq_food_logs_client_date_meal UNIQUE (client_id, date, meal_id),
    CONSTRAINT chk_food_logs_selection_or_skip CHECK (selected_option_id IS NOT NULL OR is_skipped = true)
);
CREATE INDEX idx_food_logs_client_date ON food_logs (client_id, date);

CREATE TABLE water_logs (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id  UUID NOT NULL REFERENCES users(id),
    local_id   UUID NOT NULL UNIQUE,
    date       DATE NOT NULL,
    amount_ml  INTEGER NOT NULL CHECK (amount_ml > 0),
    logged_time TIME,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_water_logs_client_date ON water_logs (client_id, date);

CREATE TABLE sleep_logs (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id  UUID NOT NULL REFERENCES users(id),
    local_id   UUID NOT NULL UNIQUE,
    date       DATE NOT NULL,
    sleep_time TIME NOT NULL,
    wake_time  TIME NOT NULL,
    quality    sleep_quality NOT NULL,
    notes      TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_sleep_logs_client_date UNIQUE (client_id, date)
);
CREATE INDEX idx_sleep_logs_client_date ON sleep_logs (client_id, date);

CREATE TABLE exercise_logs (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id        UUID NOT NULL REFERENCES users(id),
    local_id         UUID NOT NULL UNIQUE,
    date             DATE NOT NULL,
    exercise_name    VARCHAR(200) NOT NULL,
    duration_minutes INTEGER NOT NULL CHECK (duration_minutes > 0),
    calories_burned  INTEGER CHECK (calories_burned >= 0),
    notes            TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_exercise_logs_client_date ON exercise_logs (client_id, date);

CREATE TABLE medication_logs (
    id                       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id                UUID NOT NULL REFERENCES users(id),
    local_id                 UUID NOT NULL UNIQUE,
    date                     DATE NOT NULL,
    prescribed_medication_id UUID REFERENCES plan_medications(id),
    medication_name          VARCHAR(200) NOT NULL,
    dosage                   VARCHAR(100),
    taken_at                 TIME NOT NULL,
    notes                    TEXT,
    is_self_reported         BOOLEAN NOT NULL DEFAULT false,
    created_at               TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_medication_logs_client_date ON medication_logs (client_id, date);

CREATE TABLE body_measurements (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id  UUID NOT NULL REFERENCES users(id),
    local_id   UUID NOT NULL UNIQUE,
    date       DATE NOT NULL,
    weight_kg  NUMERIC(5,2) CHECK (weight_kg > 0),
    waist_cm   NUMERIC(5,2) CHECK (waist_cm > 0),
    hip_cm     NUMERIC(5,2) CHECK (hip_cm > 0),
    abdomen_cm NUMERIC(5,2) CHECK (abdomen_cm > 0),
    thigh_cm   NUMERIC(5,2) CHECK (thigh_cm > 0),
    chest_cm   NUMERIC(5,2) CHECK (chest_cm > 0),
    wrist_cm   NUMERIC(5,2) CHECK (wrist_cm > 0),
    recorded_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_body_measurements_client_date UNIQUE (client_id, date)
);
CREATE INDEX idx_body_measurements_client_date ON body_measurements (client_id, date);

CREATE TABLE lab_results (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id         UUID NOT NULL REFERENCES users(id),
    local_id          UUID NOT NULL UNIQUE,
    uploaded_by       UUID NOT NULL REFERENCES users(id),
    title             VARCHAR(200) NOT NULL,
    lab_type          lab_result_type NOT NULL,
    test_date         DATE NOT NULL,
    file_path         TEXT,
    external_link     TEXT,
    original_filename VARCHAR(255),
    mime_type         VARCHAR(100),
    file_size_bytes   BIGINT,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_lab_results_file_or_link CHECK (file_path IS NOT NULL OR external_link IS NOT NULL)
);
CREATE INDEX idx_lab_results_client_id ON lab_results (client_id);

