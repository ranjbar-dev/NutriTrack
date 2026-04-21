CREATE TABLE medications (
    id              uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    name            text        NOT NULL,
    name_normalized text        NOT NULL,
    description     text        NOT NULL DEFAULT '',
    unit            text        NOT NULL DEFAULT 'قرص',
    created_by      uuid        REFERENCES users(id) ON DELETE SET NULL,
    is_active       boolean     NOT NULL DEFAULT true,
    created_at      timestamptz NOT NULL DEFAULT NOW(),
    updated_at      timestamptz NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_medications_name_trgm  ON medications USING gin(name_normalized gin_trgm_ops);
CREATE INDEX idx_medications_created_by ON medications(created_by);
CREATE INDEX idx_medications_is_active  ON medications(is_active);
