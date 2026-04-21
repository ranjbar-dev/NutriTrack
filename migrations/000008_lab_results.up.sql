CREATE TABLE lab_results (
    id              UUID        NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id       UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    nutritionist_id UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    file_path       TEXT        NOT NULL,
    original_name   TEXT        NOT NULL,
    file_type       TEXT        NOT NULL,
    file_size       BIGINT      NOT NULL,
    notes           TEXT        NOT NULL DEFAULT '',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_lab_results_client_id ON lab_results(client_id);
