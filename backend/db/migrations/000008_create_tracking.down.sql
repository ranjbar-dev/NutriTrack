-- Migration 000008 rollback: drop tracking tables and enums in reverse dependency order
DROP TABLE IF EXISTS lab_results;
DROP TABLE IF EXISTS body_measurements;
DROP TABLE IF EXISTS medication_logs;
DROP TABLE IF EXISTS exercise_logs;
DROP TABLE IF EXISTS sleep_logs;
DROP TABLE IF EXISTS water_logs;
DROP TABLE IF EXISTS food_logs;
DROP TYPE IF EXISTS lab_result_type;
DROP TYPE IF EXISTS sleep_quality;
