# ArXiv Research Assistant

Production-grade RAG system for querying academic papers from ArXiv.

## Stack
- **Go + Fiber** — Backend API
- **Python + FastAPI** — ML Service (embedding, chunking, generation)
- **PostgreSQL** — Metadata storage
- **Qdrant** — Vector database
- **Redis** — Cache
- **Ollama** — Local LLM (Mistral / phi3:mini)
- **Prometheus + Grafana** — Monitoring

## Quick Start

### Prerequisites
- Go 1.22+
- Python 3.12+
- Docker + Docker Compose
- Ollama with `phi3:mini` or `mistral`

### Setup
```bash
# 1. Clone repo
git clone <repo-url>
cd arxiv-research-assistant

# 2. Start infrastructure
make up

# 3. Start Python ML service
cd python-ml
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
uvicorn main:app --port 8001

# 4. Start Go backend (new terminal)
make run
```

### Verify
```bash
curl http://localhost:8080/health
# Expected: {"status":"ok","services":{"ml":"ok","postgres":"ok",...}}
```

## Usage

### Ingest Papers
```bash
curl -X POST http://localhost:8001/ingest \
  -H "Content-Type: application/json" \
  -d '{"arxiv_ids": ["2312.10997"], "chunk_size": 512, "overlap": 50}'
```

### Query
```bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"query": "What is RAG?", "top_k": 2, "model": "phi3:mini"}'
```

### Monitor
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000 (admin/admin)

### Chaos Testing
```bash
./scripts/failure_inject.sh kill-qdrant
./scripts/failure_inject.sh status
./scripts/failure_inject.sh restore-qdrant
```

## API Reference

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /health | System health check |
| POST | /query | Query papers with natural language |
| GET | /metrics | Prometheus metrics |

## Architecture
```
Client
  └── Go Backend :8080
        ├── PostgreSQL :5432 (metadata)
        └── Python ML :8001
              ├── Qdrant :6333 (vectors)
              └── Ollama :11434 (LLM)
```

## Development
```bash
make test    # Run tests
make lint    # Go vet
make up      # Start Docker services
make down    # Stop Docker services
make ps      # Service status
```
