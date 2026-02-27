# Postmortem: PostgreSQL Service Down

## Scenario
PostgreSQL container dihentikan secara tiba-tiba.

## How to Reproduce
```bash
./scripts/failure_inject.sh kill-postgres
curl http://localhost:8080/health
```

## Observed Behavior
- /health return status "degraded"
- services.postgres = "error"
- Go backend masih jalan, tidak crash
- Query endpoint masih bisa serve (tidak butuh postgres langsung)

## Impact
- Health endpoint: Up (dengan status degraded)
- Query endpoint:  Masih jalan (postgres belum dipakai di query path)
- Ingest endpoint:  Akan gagal (butuh postgres untuk simpan metadata)
- Query logging:  Tidak bisa log ke query_logs table

## Root Cause
PostgreSQL container stopped. Connection pool exhausted setelah timeout.

## Improvement
1. Query logging harus graceful degrade — kalau postgres down, tetap serve query tapi skip logging
2. Ingest harus return 503 dengan pesan jelas kalau postgres tidak available
3. Setup Prometheus alert kalau postgres down lebih dari 30 detik
