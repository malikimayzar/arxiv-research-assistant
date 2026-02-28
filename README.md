# ArXiv Research Assistant

> Production-grade RAG system for querying academic papers — built with Go, Python, and a full observability stack.

![CI](https://github.com/malikimayzar/arxiv-research-assistant/actions/workflows/test.yml/badge.svg)
![Lint](https://github.com/malikimayzar/arxiv-research-assistant/actions/workflows/lint.yml/badge.svg)

---

## What This Is

A system that lets you query academic papers in natural language. You ingest papers from ArXiv, they get chunked, embedded, and stored in a vector database. When you ask a question, the system retrieves relevant chunks and generates an answer using a local LLM.

Built to be **auditable** (every query logged), **observable** (Prometheus + Grafana), and **resilient** (circuit breaker, chaos tested).

---

## Architecture
```
Browser / CLI
     │
     ▼
Go Fiber API :8080
     │
     ├── PostgreSQL :5432   (paper metadata, query logs)
     │
     └── Python FastAPI :8001
               │
               ├── Qdrant :6333      (vector search)
               ├── Ollama :11434     (local LLM)
               └── Redis :6379       (cache)

Prometheus :9090  ←  scrapes Go + Python metrics
Grafana :3000     ←  dashboards
```

---

## Stack

| Layer | Tool | Why |
|-------|------|-----|
| Backend | Go + Fiber | Concurrency, explicit error handling, performance |
| ML Service | Python + FastAPI | Best ML ecosystem, clean service boundary |
| Vector DB | Qdrant | Production-grade, filterable, REST API |
| Relational DB | PostgreSQL | Metadata, query logs, versioning |
| Cache | Redis | Embedding cache, session |
| LLM | Ollama (phi3:mini) | Fully local, no API cost |
| Monitoring | Prometheus + Grafana | Full observability stack |
| CI/CD | GitHub Actions | Lint + Test on every push |

---

## Quick Start

**Prerequisites:** Go 1.22+, Python 3.12+, Docker, Ollama
```bash
# 1. Clone
git clone https://github.com/malikimayzar/arxiv-research-assistant.git
cd arxiv-research-assistant

# 2. Pull Ollama model
ollama pull phi3:mini

# 3. Setup everything
make setup

# 4. Terminal 1 — ML service
make run-ml

# 5. Terminal 2 — Go backend
make run

# 6. Ingest a paper
make ingest

# 7. Query
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"query": "What is RAG?", "top_k": 5}'
```

---

## API Reference

### Go Backend `:8080`

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | System health + service status |
| GET | `/papers` | List ingested papers |
| GET | `/papers/:id` | Get paper by ID |
| DELETE | `/papers/:id` | Delete paper |
| POST | `/query` | Natural language query |
| GET | `/query/history` | Query logs |
| GET | `/metrics` | Prometheus metrics |

### Python ML Service `:8001`

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/ingest` | Ingest ArXiv papers |
| POST | `/embed` | Generate embeddings |
| POST | `/chunk` | Chunk text |
| POST | `/query` | RAG pipeline |
| POST | `/evaluate` | Faithfulness scoring |
| GET | `/metrics` | Prometheus metrics |

---

## Observability

Prometheus scrapes both services. Grafana dashboards available at `http://localhost:3000`.

**Tracked metrics:**

| Metric | Type | What it detects |
|--------|------|-----------------|
| `http_request_duration_seconds` | Histogram | Slow endpoints |
| `http_requests_total` | Counter | Traffic + error rate |
| `query_faithfulness_score` | Gauge | Answer quality trend |
| `qdrant_search_duration_seconds` | Histogram | Vector search performance |
| `ollama_generation_duration_seconds` | Histogram | LLM latency |
| `failure_mode_total` | Counter | correct / insufficient_ctx / error |
| `active_goroutines` | Gauge | Goroutine leak detection |
| `postgres_connections_active` | Gauge | Connection pool exhaustion |

---

## Chaos Testing

5 failure scenarios tested and documented in `docs/postmortem/`:
```bash
bash scripts/failure_inject.sh kill-qdrant       # Scenario 1
bash scripts/failure_inject.sh kill-postgres     # Scenario 2
bash scripts/failure_inject.sh kill-ml           # Scenario 3
bash scripts/failure_inject.sh inject-garbage    # Scenario 4
bash scripts/failure_inject.sh fill-disk         # Scenario 5

bash scripts/failure_inject.sh status            # Check system
```

Circuit breaker in Go backend trips after 5 failures — ML service down returns graceful error instead of hanging.

---

## Ablation Study

Compare answer quality across different retrieval configs:
```bash
make ablation
```

---

## Development
```bash
make setup      # Full setup from scratch
make run        # Start Go backend
make run-ml     # Start Python ML service
make test       # Go tests
make lint       # Go vet
make build      # Go build
make health     # Check system health
make ingest     # Ingest sample paper (2312.10997)
make ablation   # Run ablation study
make up         # Docker services up
make down       # Docker services down
```

---

## Project Structure
```
arxiv-research-assistant/
├── go-backend/           # Go API service
│   ├── cmd/server/       # Entry point
│   ├── internal/
│   │   ├── api/          # HTTP handlers
│   │   ├── client/       # ML service client + circuit breaker
│   │   ├── middleware/   # Prometheus metrics, request ID
│   │   └── repository/   # PostgreSQL queries
│   └── migrations/       # SQL migrations
├── python-ml/            # Python ML service
│   ├── routers/          # FastAPI endpoints
│   └── engine/           # Ingestion, retrieval, generation, evaluation
├── infra/                # Docker Compose, Prometheus, Grafana
├── scripts/              # Chaos testing, ablation study
├── docs/                 # Architecture, threat model, postmortem
└── tests/                # Unit + integration tests
```

---

*Built by [Maliki Mayzar](https://github.com/malikimayzar) — 2026*
