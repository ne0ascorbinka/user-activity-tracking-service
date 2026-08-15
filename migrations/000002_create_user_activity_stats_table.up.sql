CREATE TABLE IF NOT EXISTS user_activity_stats (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    event_count INTEGER NOT NULL,
    period_start TIMESTAMPTZ NOT NULL,
    period_end TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_user_activity_stats_user_period UNIQUE (user_id, period_start, period_end)
);

-- Index for querying aggregate history by user_id
CREATE INDEX IF NOT EXISTS idx_user_activity_stats_user_id ON user_activity_stats (user_id);

-- Index for querying aggregate history by period
CREATE INDEX IF NOT EXISTS idx_user_activity_stats_period ON user_activity_stats (period_start, period_end);
