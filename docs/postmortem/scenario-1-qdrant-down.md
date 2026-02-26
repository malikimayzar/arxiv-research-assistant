# Postmortem: Qdrant Service Down

## Scenario
Qdrant vector database container dihentikan secara tiba-tiba.

## How to Reproduce
```bash
./scripts/failure_inject.sh kill-qdrant
curl -X POST http://localhost:8080/query -H 'Content-Type: application/json' -d '{"query": "What is RAG?", "top_k": 2}'
```

## Observed Behavior
- Go backend return HTTP 500
- Response body: `{"status":500,"message":"ml service returned status 500","request_id":"..."}`
- Request ID tersedia untuk tracing
- /health endpoint masih jalan

## Impact
- Query endpoint: ❌ Down
- Health endpoint: ✅ Up
- Ingest endpoint: ❌ Down
- User experience: Error 500 tanpa context

## Root Cause
Python ML service gagal connect ke Qdrant saat vector search.
Error tidak di-handle secara graceful — langsung return 500.

## Mitigation
1. Tambah circuit breaker di Python ML service
2. Return 503 Service Unavailable (lebih tepat dari 500)
3. Tambah pesan error yang lebih informatif untuk user
4. Setup Prometheus alert kalau Qdrant down

## Timeline
- T+0: Qdrant stopped
- T+1s: Next query request gagal
- T+?: Alert Prometheus trigger (belum diimplementasi)
