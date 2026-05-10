CREATE TABLE IF NOT EXISTS notes (
    id VARCHAR(36) PRIMARY KEY,
    author_id VARCHAR(36) REFERENCES users(id) ON DELETE SET NULL,
    title VARCHAR(200),
    slug VARCHAR(255) UNIQUE,
    content_markdown TEXT NOT NULL DEFAULT '',
    content_text TEXT NOT NULL DEFAULT '',
    mood VARCHAR(80),
    weather VARCHAR(80),
    location VARCHAR(120),
    status VARCHAR(20) NOT NULL DEFAULT 'published',
    pinned_at TIMESTAMPTZ,
    published_at TIMESTAMPTZ,
    view_count BIGINT NOT NULL DEFAULT 0,
    like_count BIGINT NOT NULL DEFAULT 0,
    comment_count BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_notes_status CHECK (status IN ('draft', 'published', 'private', 'archived'))
);

CREATE INDEX IF NOT EXISTS idx_notes_status_published ON notes(status, published_at DESC);
CREATE INDEX IF NOT EXISTS idx_notes_pinned ON notes(pinned_at DESC);
