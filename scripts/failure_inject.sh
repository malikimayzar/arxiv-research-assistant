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

  "status")
    echo "📊 System Status:"
    docker compose -f infra/docker-compose.yml ps
    echo ""
    curl -s http://localhost:8080/health | python3 -m json.tool
    ;;

  *)
    echo "Usage: $0 [scenario]"
    echo "Scenarios: kill-qdrant, restore-qdrant, kill-postgres, restore-postgres, kill-ml, kill-go, status"
    ;;
esac
