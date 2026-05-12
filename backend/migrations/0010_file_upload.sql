CREATE TABLE IF NOT EXISTS files (
    id VARCHAR(36) PRIMARY KEY,
    uploader_id VARCHAR(36) REFERENCES users(id) ON DELETE SET NULL,
    original_name VARCHAR(255) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    storage_type VARCHAR(20) NOT NULL DEFAULT 'local',
    object_key VARCHAR(500) NOT NULL,
    public_url VARCHAR(500) NOT NULL,
    mime_type VARCHAR(120) NOT NULL,
    byte_size BIGINT NOT NULL,
    width INTEGER,
    height INTEGER,
    ref_type VARCHAR(20),
    ref_id VARCHAR(36),
    status VARCHAR(20) NOT NULL DEFAULT 'attached',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_files_storage CHECK (storage_type IN ('local', 's3')),
    CONSTRAINT chk_files_status CHECK (status IN ('attached', 'orphan', 'deleted'))
);

CREATE INDEX IF NOT EXISTS idx_files_ref ON files(ref_type, ref_id);
CREATE INDEX IF NOT EXISTS idx_files_status_created ON files(status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_files_uploader ON files(uploader_id, created_at DESC);
