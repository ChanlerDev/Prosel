CREATE TABLE IF NOT EXISTS posts (
    id VARCHAR(36) PRIMARY KEY,
    author_id VARCHAR(36) REFERENCES users(id) ON DELETE SET NULL,
    category_id VARCHAR(36),
    title VARCHAR(200) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    excerpt VARCHAR(500),
    content_markdown TEXT NOT NULL DEFAULT '',
    content_text TEXT NOT NULL DEFAULT '',
    cover_image VARCHAR(500),
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    featured BOOLEAN NOT NULL DEFAULT FALSE,
    pinned_at TIMESTAMPTZ,
    published_at TIMESTAMPTZ,
    seo_title VARCHAR(200),
    seo_description VARCHAR(500),
    view_count BIGINT NOT NULL DEFAULT 0,
    like_count BIGINT NOT NULL DEFAULT 0,
    comment_count BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_posts_status CHECK (status IN ('draft', 'published', 'archived'))
);

CREATE INDEX IF NOT EXISTS idx_posts_status_published ON posts(status, published_at DESC);
CREATE INDEX IF NOT EXISTS idx_posts_featured ON posts(featured, published_at DESC);
CREATE INDEX IF NOT EXISTS idx_posts_category ON posts(category_id, published_at DESC);
CREATE INDEX IF NOT EXISTS idx_posts_author ON posts(author_id);
