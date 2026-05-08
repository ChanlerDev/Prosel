CREATE TABLE IF NOT EXISTS schema_migrations (
    version VARCHAR(64) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS site_settings (
    id VARCHAR(36) PRIMARY KEY,
    setting_key VARCHAR(100) NOT NULL UNIQUE,
    setting_value TEXT,
    value_type VARCHAR(20) NOT NULL DEFAULT 'string',
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_site_settings_key ON site_settings(setting_key);

INSERT INTO site_settings (id, setting_key, setting_value, value_type, description)
VALUES
    ('00000000-0000-0000-0000-000000000001', 'site_name', 'Prosel', 'string', 'Public site name'),
    ('00000000-0000-0000-0000-000000000002', 'site_description', 'A personal blog powered by Prosel', 'string', 'Public site description'),
    ('00000000-0000-0000-0000-000000000003', 'site_url', 'http://localhost:3000', 'string', 'Canonical site URL'),
    ('00000000-0000-0000-0000-000000000004', 'posts_per_page', '10', 'number', 'Default public posts per page'),
    ('00000000-0000-0000-0000-000000000005', 'comment_moderation', 'true', 'boolean', 'Whether comments require moderation'),
    ('00000000-0000-0000-0000-000000000006', 'analytics_enabled', 'false', 'boolean', 'Whether lightweight analytics is enabled')
ON CONFLICT (setting_key) DO NOTHING;
