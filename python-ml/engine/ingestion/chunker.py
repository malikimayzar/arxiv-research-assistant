import logging

logger = logging.getLogger(__name__)

def chunk_text(text: str, chunk_size: int = 512, overlap: int = 50) -> list[str]:
    chunks = []
    start = 0
    
    while start < len(text):
        end = start + chunk_size
        chunk = text[start:end].strip()
        if chunk:
            chunks.append(chunk)
        start += chunk_size - overlap
    
    logger.info(f"Chunked text into {len(chunks)} chunks")
    return chunks
