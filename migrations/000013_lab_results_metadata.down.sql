ALTER TABLE lab_results
    DROP COLUMN IF EXISTS title,
    DROP COLUMN IF EXISTS result_type,
    DROP COLUMN IF EXISTS test_date,
    DROP COLUMN IF EXISTS link;
