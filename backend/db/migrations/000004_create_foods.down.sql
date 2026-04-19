-- Drop tables and types in reverse order
DROP TABLE IF EXISTS food_categories;
DROP TABLE IF EXISTS foods;
DROP TYPE IF EXISTS measurement_unit;
DROP TYPE IF EXISTS food_category;
DROP FUNCTION IF EXISTS normalize_persian;
DROP EXTENSION IF EXISTS pg_trgm;
