-- Enable pg_trgm extension for Persian fuzzy search
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Create normalize_persian function for Persian character normalization
-- Converts ی (U+06CC FARSI YEH) → ي (U+064A YEH)
-- Converts ک (U+06A9 KEHEH) → ك (U+0643 KAF)
CREATE OR REPLACE FUNCTION normalize_persian(text) RETURNS text AS $$
BEGIN
    RETURN REPLACE(REPLACE($1, 'ی', 'ي'), 'ک', 'ك');
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- Create food_category enum
CREATE TYPE food_category AS ENUM (
    'breakfast',
    'lunch',
    'dinner',
    'snack',
    'fruit',
    'beverage',
    'supplement',
    'other'
);

-- Create measurement_unit enum
CREATE TYPE measurement_unit AS ENUM (
    'gram',
    'kg',
    'tablespoon',
    'teaspoon',
    'cup',
    'piece',
    'slice',
    'palm',
    'matchbox',
    'bowl',
    'ml',
    'liter'
);

-- Create foods table
CREATE TABLE foods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    name_normalized VARCHAR(200) NOT NULL,
    description TEXT,
    calories DECIMAL(8,2) NOT NULL DEFAULT 0,
    protein_g DECIMAL(8,2) NOT NULL DEFAULT 0,
    carbs_g DECIMAL(8,2) NOT NULL DEFAULT 0,
    fat_g DECIMAL(8,2) NOT NULL DEFAULT 0,
    fiber_g DECIMAL(8,2) NOT NULL DEFAULT 0,
    sugar_g DECIMAL(8,2) NOT NULL DEFAULT 0,
    sodium_mg DECIMAL(8,2) NOT NULL DEFAULT 0,
    measurement_unit measurement_unit NOT NULL DEFAULT 'gram',
    measurement_amount DECIMAL(8,2) NOT NULL DEFAULT 100,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create GIN trigram index on normalized name for fuzzy search
CREATE INDEX idx_foods_name_trgm ON foods USING GIN (name_normalized gin_trgm_ops);

-- Create additional indexes
CREATE INDEX idx_foods_is_active ON foods (is_active);
CREATE INDEX idx_foods_created_by ON foods (created_by);

-- Create food_categories junction table
CREATE TABLE food_categories (
    food_id UUID NOT NULL REFERENCES foods(id) ON DELETE CASCADE,
    category food_category NOT NULL,
    PRIMARY KEY (food_id, category)
);
