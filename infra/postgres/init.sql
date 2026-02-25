-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Papers table
CREATE TABLE IF NOT EXISTS papers (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  arxiv_id    VARCHAR(20) UNIQUE NOT NULL,
  title       TEXT NOT NULL,
  authors     TEXT[],
  abstract    TEXT,
  categories  TEXT[],
  published   DATE,
  ingested_at TIMESTAMPTZ DEFAULT NOW()
);

-- Chunks table
CREATE TABLE IF NOT EXISTS chunks (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  paper_id    UUID REFERENCES papers(id) ON DELETE CASCADE,
  chunk_index INT NOT NULL,
  text        TEXT NOT NULL,
  qdrant_id   UUID,
  char_count  INT,
  created_at  TIMESTAMPTZ DEFAULT NOW()
);

-- Query logs table
CREATE TABLE IF NOT EXISTS query_logs (
  id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  query              TEXT NOT NULL,
  answer             TEXT,
  retrieval_ms       INT,
  generation_ms      INT,
  faithfulness       FLOAT,
  failure_mode       VARCHAR(50),
  model              VARCHAR(50),
  retrieval_strategy VARCHAR(50),
  chunk_size_used    INT,
  session_id         VARCHAR(100),
  created_at         TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_papers_arxiv_id ON papers(arxiv_id);
CREATE INDEX IF NOT EXISTS idx_chunks_paper_id ON chunks(paper_id);
CREATE INDEX IF NOT EXISTS idx_query_logs_created_at ON query_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_query_logs_failure_mode ON query_logs(failure_mode);
