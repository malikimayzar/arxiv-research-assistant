.PHONY: run test lint tidy up down ps

run:
	cd go-backend && go run cmd/server/main.go

test:
	cd go-backend && go test ./... -v

tidy:
	cd go-backend && go mod tidy

lint:
	cd go-backend && go vet ./...

up:
	docker compose -f infra/docker-compose.yml up -d

down:
	docker compose -f infra/docker-compose.yml down

ps:
	docker compose -f infra/docker-compose.yml ps
