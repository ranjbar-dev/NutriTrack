CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE user_role AS ENUM ('super_admin', 'nutritionist', 'client');
CREATE TYPE gender_type AS ENUM ('male', 'female');

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role user_role NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE,
    password_hash VARCHAR(255),
    mobile VARCHAR(15) UNIQUE,
    date_of_birth DATE,
    height_cm REAL,
    gender gender_type,
    nutritionist_id UUID REFERENCES users(id),
    is_active BOOLEAN NOT NULL DEFAULT true,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_mobile ON users (mobile) WHERE mobile IS NOT NULL;
CREATE INDEX idx_users_nutritionist ON users (nutritionist_id) WHERE role = 'client';
CREATE INDEX idx_users_email ON users (email) WHERE email IS NOT NULL;
