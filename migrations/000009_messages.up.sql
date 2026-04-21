CREATE TABLE messages (
    id              UUID        NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    sender_id       UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    receiver_id     UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content         TEXT        NOT NULL DEFAULT '',
    attachment_path TEXT,
    attachment_type TEXT,
    attachment_size BIGINT,
    attachment_name TEXT,
    read_at         TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_messages_receiver_id ON messages(receiver_id);
CREATE INDEX idx_messages_conversation ON messages(
    LEAST(sender_id, receiver_id),
    GREATEST(sender_id, receiver_id),
    created_at
);
