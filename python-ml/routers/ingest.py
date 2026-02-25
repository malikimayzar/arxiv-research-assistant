from fastapi import APIRouter, HTTPException
from pydantic import BaseModel
import logging
import main as app_state
from engine.ingestion.arxiv_fetcher import fetch_metadata
from engine.ingestion.pdf_parser import download_and_parse_pdf
from engine.ingestion.chunker import chunk_text
from engine.retrieval.qdrant_client import get_client, ensure_collection, upsert_chunks

logger = logging.getLogger(__name__)
router = APIRouter()

class IngestRequest(BaseModel):
    arxiv_ids: list[str]
    chunk_size: int = 512
    overlap: int = 50

class IngestResponse(BaseModel):
    ingested: list[str]
    failed: list[str]
    total_chunks: int

@router.post("/ingest", response_model=IngestResponse)
def ingest(request: IngestRequest):
    if not request.arxiv_ids:
        raise HTTPException(status_code=400, detail="arxiv_ids cannot be empty")
    
    if app_state.model is None:
        raise HTTPException(status_code=503, detail="model not loaded")

    qdrant = get_client()
    ensure_collection(qdrant)

    ingested = []
    failed = []
    total_chunks = 0

    papers = fetch_metadata(request.arxiv_ids)
    if not papers:
        raise HTTPException(status_code=404, detail="No papers found")

    for paper in papers:
        try:
            # Parse PDF
            text = download_and_parse_pdf(paper.pdf_url)
            
            # Chunk
            chunks = chunk_text(text, request.chunk_size, request.overlap)
            
            # Embed all chunks
            embeddings = app_state.model.encode(chunks).tolist()
            
            # Store in Qdrant
            chunk_dicts = [
                {
                    "text": chunk,
                    "embedding": embeddings[i],
                    "arxiv_id": paper.arxiv_id,
                    "categories": paper.categories,
                }
                for i, chunk in enumerate(chunks)
            ]
            upsert_chunks(qdrant, chunk_dicts)
            
            ingested.append(paper.arxiv_id)
            total_chunks += len(chunks)
            logger.info(f"✅ Ingested {paper.arxiv_id}: {len(chunks)} chunks")

        except Exception as e:
            logger.error(f"❌ Failed to ingest {paper.arxiv_id}: {e}")
            failed.append(paper.arxiv_id)

    return IngestResponse(
        ingested=ingested,
        failed=failed,
        total_chunks=total_chunks
    )
