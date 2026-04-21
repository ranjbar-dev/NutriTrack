-- Migration 001: Foundation extensions and schema setup
-- pg_trgm: required for Persian trigram similarity search (food, medication names)
-- uuid-ossp: required for UUID primary key generation

CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Verify pg_trgm is active (will fail migration if extension is not available)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_extension WHERE extname = 'pg_trgm'
    ) THEN
        RAISE EXCEPTION 'pg_trgm extension is required but not available';
    END IF;
END $$;
