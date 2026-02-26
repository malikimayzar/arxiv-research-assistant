# Postmortem: ML Service Down

## Scenario
Python ML service dihentikan via pkill.

## Observed Behavior
- /health masih return ml: "ok" — FALSE POSITIVE
- Query masih berhasil — ML service ternyata tidak benar-benar mati di WSL
- pkill tidak efektif untuk mematikan uvicorn di WSL environment

## Bug Found
Health check ML service hanya cek saat startup, tidak real-time per request.
Harusnya setiap /health request melakukan ping ke ML service.

## Improvement
1. Fix pkill script untuk WSL — gunakan port-based kill
2. Health check harus aktif ping ML service setiap request
