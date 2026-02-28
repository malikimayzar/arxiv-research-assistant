-- Prompt versions table
CREATE TABLE IF NOT EXISTS prompt_versions (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  version     VARCHAR(20) NOT NULL UNIQUE,
  name        VARCHAR(100) NOT NULL,
  template    TEXT NOT NULL,
  description TEXT,
  is_active   BOOLEAN DEFAULT FALSE,
  created_at  TIMESTAMPTZ DEFAULT NOW()
);

INSERT INTO prompt_versions (version, name, template, description, is_active)
VALUES (
  'v1.0',
  'strict_context',
  'Answer strictly based on the provided context.
If the answer is not in the context, say "Not found in context.".

Context:
{context}

Question: {query}

Answer:',
  'Default prompt — strict context adherence',
  TRUE
) ON CONFLICT (version) DO NOTHING;

-- Embedding versions table
CREATE TABLE IF NOT EXISTS embedding_versions (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  model_name  VARCHAR(100) NOT NULL,
  model_dim   INT NOT NULL,
  batch_id    UUID DEFAULT gen_random_uuid(),
  paper_count INT DEFAULT 0,
  chunk_count INT DEFAULT 0,
  created_at  TIMESTAMPTZ DEFAULT NOW()
);

INSERT INTO embedding_versions (model_name, model_dim, paper_count, chunk_count)
VALUES ('all-MiniLM-L6-v2', 384, 1, 238);

-- Version tracking columns on query_logs
ALTER TABLE query_logs ADD COLUMN IF NOT EXISTS prompt_version VARCHAR(20) DEFAULT 'v1.0';
ALTER TABLE query_logs ADD COLUMN IF NOT EXISTS embedding_version VARCHAR(100) DEFAULT 'all-MiniLM-L6-v2';
