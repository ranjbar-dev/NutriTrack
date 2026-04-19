-- Create medication_form enum
CREATE TYPE medication_form AS ENUM (
    'tablet',
    'capsule',
    'syrup',
    'injection',
    'drop',
    'powder',
    'other'
);

-- Create medications table
CREATE TABLE medications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    name_normalized VARCHAR(200) NOT NULL,
    generic_name VARCHAR(200),
    generic_name_normalized VARCHAR(200),
    form medication_form NOT NULL DEFAULT 'tablet',
    dosage_unit VARCHAR(50),
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create GIN trigram indexes for fuzzy search
CREATE INDEX idx_medications_name_trgm ON medications USING GIN (name_normalized gin_trgm_ops);
CREATE INDEX idx_medications_generic_trgm ON medications USING GIN (generic_name_normalized gin_trgm_ops) WHERE generic_name_normalized IS NOT NULL;

-- Create additional indexes
CREATE INDEX idx_medications_is_active ON medications (is_active);
CREATE INDEX idx_medications_created_by ON medications (created_by);
