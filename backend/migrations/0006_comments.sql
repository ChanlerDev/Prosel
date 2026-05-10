CREATE TABLE IF NOT EXISTS comments (
    id VARCHAR(36) PRIMARY KEY,
    ref_type VARCHAR(20) NOT NULL,
    ref_id VARCHAR(36) NOT NULL,
    parent_id VARCHAR(36) REFERENCES comments(id) ON DELETE CASCADE,
    root_id VARCHAR(36) REFERENCES comments(id) ON DELETE CASCADE,
    author_name VARCHAR(80) NOT NULL,
    author_email VARCHAR(255) NOT NULL,
    author_website VARCHAR(500),
    author_ip VARCHAR(64),
    user_agent TEXT,
    content TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    is_admin_reply BOOLEAN NOT NULL DEFAULT FALSE,
    is_pinned BOOLEAN NOT NULL DEFAULT FALSE,
    reply_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_comments_ref_type CHECK (ref_type IN ('post', 'note', 'page')),
    CONSTRAINT chk_comments_status CHECK (status IN ('pending', 'approved', 'rejected', 'spam'))
);

CREATE INDEX IF NOT EXISTS idx_comments_ref ON comments(ref_type, ref_id, status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_comments_parent ON comments(parent_id);
CREATE INDEX IF NOT EXISTS idx_comments_root ON comments(root_id, created_at ASC);
CREATE INDEX IF NOT EXISTS idx_comments_status ON comments(status, created_at DESC);
