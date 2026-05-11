CREATE TABLE IF NOT EXISTS search_documents (
    id VARCHAR(36) PRIMARY KEY,
    ref_type VARCHAR(20) NOT NULL,
    ref_id VARCHAR(36) NOT NULL,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255),
    excerpt VARCHAR(500),
    search_text TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'published',
    published_at TIMESTAMPTZ,
    search_vector TSVECTOR GENERATED ALWAYS AS (
        setweight(to_tsvector('simple', coalesce(title, '')), 'A') ||
        setweight(to_tsvector('simple', coalesce(excerpt, '')), 'B') ||
        setweight(to_tsvector('simple', coalesce(search_text, '')), 'C')
    ) STORED,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_search_ref_type CHECK (ref_type IN ('post', 'note', 'page')),
    UNIQUE (ref_type, ref_id)
);

CREATE INDEX IF NOT EXISTS idx_search_documents_vector ON search_documents USING GIN(search_vector);
CREATE INDEX IF NOT EXISTS idx_search_documents_ref ON search_documents(ref_type, ref_id);
CREATE INDEX IF NOT EXISTS idx_search_documents_published ON search_documents(status, published_at DESC);
