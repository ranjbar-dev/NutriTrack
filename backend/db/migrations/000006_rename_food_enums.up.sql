-- Rename food_category enum to avoid sqlc naming conflict with food_categories table
ALTER TYPE food_category RENAME TO food_category_type;
ALTER TYPE measurement_unit RENAME TO measurement_unit_type;
