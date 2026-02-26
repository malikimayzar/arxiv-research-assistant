import httpx
import logging
import time

logger = logging.getLogger(__name__)

OLLAMA_URL = "http://localhost:11434"

def generate(query: str, context_chunks: list[str], model: str = "phi3:mini") -> dict:
    context = "\n\n".join([chunk[:200] for chunk in context_chunks[:2]])
    
    prompt = f"""Answer briefly based on context only.

Context: {context}

Question: {query}

Short answer:"""

    start = time.time()
    
    with httpx.Client(timeout=600) as client:
        response = client.post(
            f"{OLLAMA_URL}/api/generate",
            json={
                "model": model,
                "prompt": prompt,
                "stream": False,
                "options": {
                    "num_predict": 150,
                    "temperature": 0.1,
                }
            }
        )
        response.raise_for_status()
    
    latency_ms = int((time.time() - start) * 1000)
    result = response.json()
    
    logger.info(f"Generated answer in {latency_ms}ms using {model}")
    
    return {
        "answer": result["response"].strip(),
        "latency_ms": latency_ms,
        "model": model,
    }