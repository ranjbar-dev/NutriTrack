ALTER TABLE foods
    ADD COLUMN sugar  numeric(8,2) NOT NULL DEFAULT 0,
    ADD COLUMN sodium numeric(8,2) NOT NULL DEFAULT 0,
    ADD COLUMN amount numeric(8,2) NOT NULL DEFAULT 100;

COMMENT ON COLUMN foods.amount IS 'Base measurement amount for all nutritional values (e.g., 100g)';
COMMENT ON COLUMN foods.sodium IS 'Sodium content in milligrams per measurement unit';
