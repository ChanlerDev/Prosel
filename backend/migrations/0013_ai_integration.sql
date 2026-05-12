CREATE TABLE IF NOT EXISTS ai_summaries (
    id VARCHAR(36) PRIMARY KEY,
    ref_type VARCHAR(20) NOT NULL,
    ref_id VARCHAR(36) NOT NULL,
    language VARCHAR(20) NOT NULL DEFAULT 'zh',
    content_hash VARCHAR(64) NOT NULL,
    summary TEXT NOT NULL,
    keywords TEXT[] NOT NULL DEFAULT '{}',
    provider VARCHAR(50),
    model VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_ai_summaries_ref_type CHECK (ref_type IN ('post', 'note', 'page')),
    UNIQUE (ref_type, ref_id, language)
);

CREATE TABLE IF NOT EXISTS ai_translations (
    id VARCHAR(36) PRIMARY KEY,
    ref_type VARCHAR(20) NOT NULL,
    ref_id VARCHAR(36) NOT NULL,
    source_language VARCHAR(20) NOT NULL,
    target_language VARCHAR(20) NOT NULL,
    content_hash VARCHAR(64) NOT NULL,
    title VARCHAR(255),
    summary TEXT,
    content_markdown TEXT NOT NULL,
    provider VARCHAR(50),
    model VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_ai_translations_ref_type CHECK (ref_type IN ('post', 'note', 'page')),
    UNIQUE (ref_type, ref_id, target_language)
);

CREATE INDEX IF NOT EXISTS idx_ai_summaries_ref ON ai_summaries(ref_type, ref_id);
CREATE INDEX IF NOT EXISTS idx_ai_translations_ref ON ai_translations(ref_type, ref_id);
