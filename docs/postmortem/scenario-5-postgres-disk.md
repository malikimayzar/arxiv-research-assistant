# Scenario 5: PostgreSQL Disk Pressure

## Date
2026-02-28

## Cause
Simulasi disk pressure pada PostgreSQL container dengan mengisi /tmp menggunakan dd.

## Injection Command
bash scripts/failure_inject.sh fill-disk

## Limitation of Simulation
WSL2 disk size (~1TB) terlalu besar untuk di-fill secara praktis.
Scenario ini lebih relevan di environment dengan disk terbatas (production VPS 20-50GB).

## Realistic Failure Mode
Pada disk yang benar-benar penuh:
- PostgreSQL INSERT gagal: could not extend file: No space left on device
- query_logs write gagal — query tetap diproses tapi tidak ter-log
- Paper ingest gagal di tengah jalan — partial data
- PostgreSQL bisa crash dan require manual recovery

## Observed Behavior (Simulated)
- Read operations (GET /papers, GET /query/history) tetap jalan
- Write operations gagal graceful — error di-log tapi response tetap 200
- Go backend tidak crash karena query log failure di-handle goroutine terpisah

## Impact
- Severity: High — data loss, query logs tidak tersimpan
- Detection: Prometheus error rate naik, postgres connections drop
- User experience: Query tetap jalan, tapi history hilang

## Mitigation
1. Alert Prometheus kalau disk usage lebih dari 80 persen
2. WAL monitoring
3. Separate disk untuk PostgreSQL data vs system
4. Query log failure tidak boleh gagalin response user — sudah implemented

## Recovery
bash scripts/failure_inject.sh free-disk
