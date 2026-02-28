# Scenario 4: Garbage Chunks Injection

## Date
2026-02-28

## Cause
20 garbage chunks dengan random vectors dan nonsense text di-inject ke Qdrant collection `arxiv_chunks`.

## Injection Command
```bash
bash scripts/failure_inject.sh inject-garbage
```

## Observed Behavior
- Query tetap jalan — sistem tidak crash
- Retrieval mengembalikan duplicate chunks karena garbage vectors mengacaukan similarity ranking
- Answer quality menurun — phi3:mini generate jawaban yang salah ("Representational Agreement Graphs") padahal RAG = Retrieval-Augmented Generation
- Faithfulness score diperkirakan turun karena context yang retrieved tidak relevan

## Impact
- **Severity**: Medium — sistem jalan tapi answer quality degraded
- **User experience**: Jawaban salah tanpa indikasi error
- **Detection**: Hanya bisa dideteksi via faithfulness score monitoring di Grafana

## Metrics
- retrieval_ms: 301ms (normal)
- generation_ms: 88220ms (normal)
- Answer correctness: ❌ salah

## Mitigation
1. Monitor faithfulness score di Grafana — alert kalau < 0.5
2. Implement input validation saat ingest — reject chunks dengan entropy terlalu rendah
3. Periodic Qdrant collection audit — cek distribusi similarity scores

## Recovery
```bash
bash scripts/failure_inject.sh remove-garbage
```
