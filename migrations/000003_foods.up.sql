CREATE TABLE food_categories (
    id         uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    name       text        NOT NULL UNIQUE,
    created_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE foods (
    id              uuid         PRIMARY KEY DEFAULT gen_random_uuid(),
    name            text         NOT NULL,
    name_normalized text         NOT NULL,
    unit            text         NOT NULL DEFAULT 'گرم',
    calories        numeric(8,2) NOT NULL DEFAULT 0,
    protein         numeric(8,2) NOT NULL DEFAULT 0,
    carbohydrate    numeric(8,2) NOT NULL DEFAULT 0,
    fat             numeric(8,2) NOT NULL DEFAULT 0,
    fiber           numeric(8,2) NOT NULL DEFAULT 0,
    created_by      uuid         REFERENCES users(id) ON DELETE SET NULL,
    is_active       boolean      NOT NULL DEFAULT true,
    created_at      timestamptz  NOT NULL DEFAULT NOW(),
    updated_at      timestamptz  NOT NULL DEFAULT NOW()
);

CREATE TABLE food_category_mappings (
    food_id     uuid NOT NULL REFERENCES foods(id) ON DELETE CASCADE,
    category_id uuid NOT NULL REFERENCES food_categories(id) ON DELETE CASCADE,
    PRIMARY KEY (food_id, category_id)
);

CREATE INDEX idx_foods_name_trgm  ON foods USING gin(name_normalized gin_trgm_ops);
CREATE INDEX idx_foods_created_by ON foods(created_by);
CREATE INDEX idx_foods_is_active  ON foods(is_active);
