-- Migration 000009: Create Communication & Collaboration tables

CREATE TYPE food_request_status AS ENUM ('pending', 'approved', 'rejected');

CREATE TABLE messages (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sender_id       UUID NOT NULL REFERENCES users(id),
    receiver_id     UUID NOT NULL REFERENCES users(id),
    content         TEXT,
    attachment_type TEXT CHECK (attachment_type IN ('image', 'pdf')),
    attachment_path TEXT,
    attachment_name TEXT,
    sent_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    read_at         TIMESTAMPTZ,
    CONSTRAINT chk_messages_has_content CHECK (content IS NOT NULL OR attachment_path IS NOT NULL)
);

CREATE INDEX idx_messages_sender_sent   ON messages (sender_id, sent_at);
CREATE INDEX idx_messages_receiver_sent ON messages (receiver_id, sent_at);
CREATE INDEX idx_messages_conversation  ON messages (LEAST(sender_id, receiver_id), GREATEST(sender_id, receiver_id), sent_at);
CREATE INDEX idx_messages_unread        ON messages (receiver_id, read_at) WHERE read_at IS NULL;

CREATE TABLE food_requests (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    food_name        TEXT NOT NULL,
    description      TEXT,
    status           food_request_status NOT NULL DEFAULT 'pending',
    rejection_reason TEXT,
    requested_by     UUID NOT NULL REFERENCES users(id),
    reviewed_by      UUID REFERENCES users(id),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_food_requests_requested_by ON food_requests (requested_by, created_at DESC);
CREATE INDEX idx_food_requests_status       ON food_requests (status) WHERE status = 'pending';
