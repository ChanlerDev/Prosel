CREATE TABLE IF NOT EXISTS pages (
    id VARCHAR(36) PRIMARY KEY,
    author_id VARCHAR(36) REFERENCES users(id) ON DELETE SET NULL,
    title VARCHAR(200) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    subtitle VARCHAR(300),
    content_markdown TEXT NOT NULL DEFAULT '',
    content_text TEXT NOT NULL DEFAULT '',
    template VARCHAR(40) NOT NULL DEFAULT 'default',
    status VARCHAR(20) NOT NULL DEFAULT 'published',
    sort_order INTEGER NOT NULL DEFAULT 0,
    seo_title VARCHAR(200),
    seo_description VARCHAR(500),
    view_count BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_pages_template CHECK (template IN ('default', 'about', 'friends', 'projects')),
    CONSTRAINT chk_pages_status CHECK (status IN ('draft', 'published', 'archived'))
);

CREATE INDEX IF NOT EXISTS idx_pages_status_order ON pages(status, sort_order);

CREATE TABLE IF NOT EXISTS friends (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    url VARCHAR(500) NOT NULL UNIQUE,
    avatar_url VARCHAR(500),
    description VARCHAR(500),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_friends_status CHECK (status IN ('active', 'pending', 'hidden'))
);

CREATE INDEX IF NOT EXISTS idx_friends_status_order ON friends(status, sort_order);
