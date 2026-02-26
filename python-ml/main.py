from fastapi import FastAPI
from contextlib import asynccontextmanager
from sentence_transformers import SentenceTransformer
from routers import embed, chunk, ingest, query, evaluate
import logging

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

model: SentenceTransformer = None

@asynccontextmanager
async def lifespan(app: FastAPI):
    global model
    logger.info("Loading embedding model...")
    model = SentenceTransformer("all-MiniLM-L6-v2")
    logger.info("✅ Embedding model loaded")
    yield

app = FastAPI(title="ArXiv ML Service", version="0.1.0", lifespan=lifespan)

app.include_router(embed.router)
app.include_router(chunk.router)
app.include_router(ingest.router)
app.include_router(query.router)
app.include_router(evaluate.router)

@app.get("/health")
def health():
    return {
        "status": "ok",
        "model": "all-MiniLM-L6-v2",
        "model_loaded": model is not None
    }
