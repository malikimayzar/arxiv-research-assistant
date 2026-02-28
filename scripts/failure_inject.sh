#!/bin/bash
SCENARIO=$1
case $SCENARIO in
  "kill-qdrant")
    echo "🔴 Scenario 1: Killing Qdrant..."
    docker stop arxiv-qdrant
    echo "Qdrant stopped."
    ;;
  "restore-qdrant")
    echo "🟢 Restoring Qdrant..."
    docker start arxiv-qdrant
    echo "Qdrant restored."
    ;;
  "kill-postgres")
    echo "🔴 Scenario 2: Killing PostgreSQL..."
    docker stop arxiv-postgres
    echo "PostgreSQL stopped."
    ;;
  "restore-postgres")
    echo "🟢 Restoring PostgreSQL..."
    docker start arxiv-postgres
    echo "PostgreSQL restored."
    ;;
  "kill-ml")
    echo "🔴 Scenario 3: Killing ML service..."
    fuser -k 8001/tcp 2>/dev/null || true
    echo "ML service stopped."
    ;;
  "kill-go")
    echo "🔴 Scenario 4: Killing Go backend..."
    fuser -k 8080/tcp 2>/dev/null || true
    echo "Go backend stopped."
    ;;
  "inject-garbage")
    echo "🔴 Scenario 4: Injecting garbage chunks to Qdrant..."
    python3 -c "
from qdrant_client import QdrantClient
from qdrant_client.models import PointStruct
import uuid, random
client = QdrantClient(host='localhost', port=6333)
points = [PointStruct(id=str(uuid.uuid4()), vector=[random.uniform(-1,1) for _ in range(384)], payload={'text':'asjdhaksjdh garbage nonsense xyzxyz 1234567890','arxiv_id':'0000.00000','categories':['garbage']}) for _ in range(20)]
client.upsert(collection_name='arxiv_chunks', points=points)
print(f'Injected {len(points)} garbage chunks')
"
    ;;
  "remove-garbage")
    echo "🟢 Removing garbage chunks from Qdrant..."
    python3 -c "
from qdrant_client import QdrantClient
from qdrant_client.models import Filter, FieldCondition, MatchValue
client = QdrantClient(host='localhost', port=6333)
client.delete(collection_name='arxiv_chunks', points_selector=Filter(must=[FieldCondition(key='arxiv_id', match=MatchValue(value='0000.00000'))]))
print('Garbage chunks removed')
"
    ;;
  "fill-disk")
    echo "🔴 Scenario 5: Filling PostgreSQL container disk..."
    docker exec arxiv-postgres bash -c "dd if=/dev/urandom of=/tmp/diskfill bs=1M count=2048 2>&1 | tail -1"
    docker exec arxiv-postgres df -h /tmp
    ;;
  "free-disk")
    echo "🟢 Freeing disk space..."
    docker exec arxiv-postgres rm -f /tmp/diskfill
    docker exec arxiv-postgres df -h /tmp
    echo "✅ Disk freed"
    ;;
  "status")
    echo "📊 System Status:"
    docker compose -f infra/docker-compose.yml ps
    echo ""
    curl -s http://localhost:8080/health | python3 -m json.tool
    ;;
  *)
    echo "Usage: $0 [scenario]"
    echo "Scenarios:"
    echo "  kill-qdrant / restore-qdrant"
    echo "  kill-postgres / restore-postgres"
    echo "  kill-ml / kill-go"
    echo "  inject-garbage / remove-garbage"
    echo "  fill-disk / free-disk"
    echo "  status"
    ;;
esac
