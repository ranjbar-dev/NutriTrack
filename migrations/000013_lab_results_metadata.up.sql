ALTER TABLE lab_results
    ADD COLUMN title       TEXT NOT NULL DEFAULT '',
    ADD COLUMN result_type TEXT NOT NULL DEFAULT 'other',
    ADD COLUMN test_date   DATE,
    ADD COLUMN link        TEXT;
