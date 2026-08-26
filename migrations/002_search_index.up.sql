CREATE TABLE IF NOT EXISTS search_tenants(id TEXT PRIMARY KEY,name TEXT NOT NULL);
CREATE TABLE IF NOT EXISTS search_collections(id TEXT NOT NULL,tenant_id TEXT NOT NULL,name TEXT NOT NULL,mappings JSONB NOT NULL,version INT NOT NULL,status TEXT NOT NULL,PRIMARY KEY(tenant_id,id));
CREATE TABLE IF NOT EXISTS search_document_changes(id BIGSERIAL PRIMARY KEY,tenant_id TEXT NOT NULL,collection_id TEXT NOT NULL,document_id TEXT NOT NULL,version BIGINT NOT NULL,payload JSONB NOT NULL,idempotency_key TEXT UNIQUE,created_at TIMESTAMPTZ NOT NULL);
CREATE TABLE IF NOT EXISTS search_index_versions(id TEXT PRIMARY KEY,tenant_id TEXT NOT NULL,collection_id TEXT NOT NULL,status TEXT NOT NULL,created_at TIMESTAMPTZ NOT NULL);
CREATE TABLE IF NOT EXISTS search_index_tasks(id TEXT PRIMARY KEY,tenant_id TEXT NOT NULL,collection_id TEXT NOT NULL,status TEXT NOT NULL,attempts INT NOT NULL DEFAULT 0,updated_at TIMESTAMPTZ NOT NULL);
CREATE INDEX IF NOT EXISTS idx_search_changes_replay ON search_document_changes(tenant_id,collection_id,id);
