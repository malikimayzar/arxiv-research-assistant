-- Prompt versions table
CREATE TABLE IF NOT EXISTS prompt_versions (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  version     VARCHAR(20) NOT NULL UNIQUE,  -- e.g. 'v1.0', 'v1.1'
  name        VARCHAR(100) NOT NULL,
  template    TEXT NOT NULL,
  description TEXT,
  is_active   BOOLEAN DEFAULT FALSE,
  created_at  TIMESTAMPTZ DEFAULT NOW()
);

-- Seed default prompt
INSERT INTO prompt_versions (version, name, template, description, is_active)
VALUES (
  'v1.0',
  'strict_context',
  'Answer strictly based on the provided context.\nIf the answer is not in the context, say "Not found in context.".\n\nContext:\n{context}\n\nQuestion: {query}\n\nAnswer:',
  'Default prompt — strict context adherence',
  TRUE
);

-- Embedding model versions table
CREATE TABLE IF NOT EXISTS embedding_versions (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  model_name  VARCHAR(100) NOT NULL,
  model_dim   INT NOT NULL,
  batch_id    UUID DEFAULT gen_random_uuid(),
  paper_count INT DEFAULT 0,
  chunk_count INT DEFAULT 0,
  created_at  TIMESTAMPTZ DEFAULT NOW()
);

-- Seed current embedding model
INSERT INTO embedding_versions (model_name, model_dim, paper_count, chunk_count)
VALUES ('all-MiniLM-L6-v2', 384, 1, 238);

-- Add version tracking to query_logs
ALTER TABLE query_logs ADD COLUMN IF NOT EXISTS prompt_version VARCHAR(20) DEFAULT 'v1.0';
ALTER TABLE query_logs ADD COLUMN IF NOT EXISTS embedding_version VARCHAR(100) DEFAULT 'all-MiniLM-L6-v2';
