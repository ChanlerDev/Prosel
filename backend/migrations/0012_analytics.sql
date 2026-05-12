CREATE TABLE IF NOT EXISTS analytics_events (
    id VARCHAR(36) PRIMARY KEY,
    event_type VARCHAR(40) NOT NULL DEFAULT 'page_view',
    path VARCHAR(500) NOT NULL,
    ref_type VARCHAR(20),
    ref_id VARCHAR(36),
    referer VARCHAR(500),
    ip_hash VARCHAR(128),
    user_agent TEXT,
    country VARCHAR(80),
    device_type VARCHAR(40),
    browser VARCHAR(80),
    os VARCHAR(80),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_analytics_ref_type CHECK (ref_type IS NULL OR ref_type IN ('post', 'note', 'page'))
);

CREATE INDEX IF NOT EXISTS idx_analytics_created ON analytics_events(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_analytics_path_created ON analytics_events(path, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_analytics_ref_created ON analytics_events(ref_type, ref_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_analytics_referer_created ON analytics_events(referer, created_at DESC);
