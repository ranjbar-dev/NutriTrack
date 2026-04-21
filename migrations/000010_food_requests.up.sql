CREATE TABLE food_requests (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id        UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    nutritionist_id  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    food_name        TEXT NOT NULL,
    status           TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','approved','rejected')),
    rejection_reason TEXT,
    created_food_id  UUID REFERENCES foods(id) ON DELETE SET NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_food_requests_client ON food_requests(client_id);
CREATE INDEX idx_food_requests_nutritionist_pending ON food_requests(nutritionist_id, status) WHERE status = 'pending';
