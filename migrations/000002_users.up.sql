-- Migration 002: Users table
CREATE TABLE IF NOT EXISTS users (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    role            VARCHAR(20)  NOT NULL CHECK (role IN ('superadmin', 'nutritionist', 'client')),
    mobile          VARCHAR(15)  NOT NULL UNIQUE,
    email           VARCHAR(255) UNIQUE,
    password_hash   TEXT,
    first_name      VARCHAR(100) NOT NULL DEFAULT '',
    last_name       VARCHAR(100) NOT NULL DEFAULT '',
    gender          VARCHAR(10)  CHECK (gender IN ('male', 'female')),
    birth_date      DATE,
    height          NUMERIC(5,2),       -- cm
    weight          NUMERIC(5,2),       -- kg
    avatar_url      TEXT,
    is_active       BOOLEAN      NOT NULL DEFAULT true,
    nutritionist_id UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- Indexes for common lookups
CREATE INDEX IF NOT EXISTS idx_users_mobile          ON users(mobile);
CREATE INDEX IF NOT EXISTS idx_users_email           ON users(email) WHERE email IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_users_role            ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_nutritionist_id ON users(nutritionist_id) WHERE nutritionist_id IS NOT NULL;

-- Seed the superadmin user (password: Admin@123456 — change in production)
INSERT INTO users (id, role, mobile, email, password_hash, first_name, last_name, is_active)
VALUES (
    uuid_generate_v4(),
    'superadmin',
    '+989000000000',
    'admin@nutritrack.ir',
    '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewYpfQN9nV.UmFZy', -- Admin@123456
    'مدیر',
    'سیستم',
    true
) ON CONFLICT (email) DO NOTHING;
