CREATE TABLE IF NOT EXISTS subscribers (
    id VARCHAR(36) PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(100),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    verify_token VARCHAR(100) NOT NULL UNIQUE,
    unsubscribe_token VARCHAR(100) NOT NULL UNIQUE,
    verified_at TIMESTAMPTZ,
    unsubscribed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_subscribers_status CHECK (status IN ('pending', 'active', 'unsubscribed', 'bounced'))
);

CREATE TABLE IF NOT EXISTS email_deliveries (
    id VARCHAR(36) PRIMARY KEY,
    subscriber_id VARCHAR(36) REFERENCES subscribers(id) ON DELETE SET NULL,
    subject VARCHAR(255) NOT NULL,
    ref_type VARCHAR(20),
    ref_id VARCHAR(36),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    error_message TEXT,
    sent_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_email_deliveries_status CHECK (status IN ('pending', 'sent', 'failed'))
);

CREATE INDEX IF NOT EXISTS idx_subscribers_status ON subscribers(status);
CREATE INDEX IF NOT EXISTS idx_email_deliveries_ref ON email_deliveries(ref_type, ref_id);
