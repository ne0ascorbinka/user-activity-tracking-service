CREATE TABLE IF NOT EXISTS events (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    action VARCHAR(64) NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for filtering by user_id
CREATE INDEX IF NOT EXISTS idx_events_user_id ON events (user_id);

-- Index for filtering and range queries by created_at
CREATE INDEX IF NOT EXISTS idx_events_created_at ON events (created_at);

-- Composite index for efficient querying by user_id and created_at range
CREATE INDEX IF NOT EXISTS idx_events_user_id_created_at ON events (user_id, created_at);
