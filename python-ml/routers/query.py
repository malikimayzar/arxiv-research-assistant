from fastapi import APIRouter, HTTPException
from pydantic import BaseModel
import logging
import time
import main as app_state
from engine.retrieval.qdrant_client import get_client, search
from engine.generation.ollama_client import generate

logger = logging.getLogger(__name__)
router = APIRouter()

class QueryRequest(BaseModel):
    query: str
    top_k: int = 5
    model: str = "mistral"
    arxiv_id: str | None = None

class QueryResponse(BaseModel):
    answer: str
    sources: list[dict]
    retrieval_ms: int
    generation_ms: int
    model: str

@router.post("/query", response_model=QueryResponse)
def query(request: QueryRequest):
    if not request.query.strip():
        raise HTTPException(status_code=400, detail="query cannot be empty")
    
    if app_state.model is None:
        raise HTTPException(status_code=503, detail="model not loaded")

    # Embed query
    t0 = time.time()
    query_vector = app_state.model.encode([request.query])[0].tolist()
    
    # Search Qdrant
    qdrant = get_client()
    results = search(qdrant, query_vector, request.top_k, request.arxiv_id)
    retrieval_ms = int((time.time() - t0) * 1000)

    if not results:
        raise HTTPException(status_code=404, detail="No relevant chunks found")

    # Generate answer
    context_chunks = [r["text"] for r in results]
    generation = generate(request.query, context_chunks, request.model)

    return QueryResponse(
        answer=generation["answer"],
        sources=results,
        retrieval_ms=retrieval_ms,
        generation_ms=generation["latency_ms"],
        model=request.model,
    )
