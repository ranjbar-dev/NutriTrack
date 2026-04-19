-- Reverse migration 000007: Drop Diet Plan Engine tables in reverse FK order
DROP TABLE IF EXISTS plan_medications;
DROP TABLE IF EXISTS plan_exercises;
DROP TABLE IF EXISTS meal_option_items;
DROP TABLE IF EXISTS meal_options;
DROP TABLE IF EXISTS meals;
DROP TABLE IF EXISTS plan_days;
DROP TABLE IF EXISTS diet_plans;
DROP TYPE IF EXISTS diet_plan_status;
