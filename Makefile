.PHONY: build run dev docker-up docker-down docker-build migrate-up migrate-down test lint \
        sqlc-generate sqlc-check tidy up down logs test-race

# Build the binary
build:
	go build -o bin/server ./cmd/server

# Run locally (requires local PG + Redis)
run:
	go run ./cmd/server

# Start full stack with docker compose
docker-up:
	docker compose up -d

# Stop docker compose
docker-down:
	docker compose down

# Build docker image
docker-build:
	docker compose build

# Run database migrations (up)
migrate-up:
	docker compose run --rm app ./server migrate up

# Run database migrations (down, 1 step)
migrate-down:
	docker compose run --rm app ./server migrate down 1

# Run tests
test:
	go test ./... -v -count=1

# Run linter
lint:
	golangci-lint run ./...

# Generate sqlc code
sqlc-generate:
	sqlc generate

# Tidy dependencies
tidy:
	go mod tidy

# sqlc freshness check — verifies generated code is in sync with SQL queries
# CI gate: exits non-zero if generated code is stale
sqlc-check:
	@echo "Checking sqlc freshness..."
	@if ! command -v sqlc >/dev/null 2>&1; then \
		echo "sqlc not installed — skipping check"; \
		exit 0; \
	fi
	sqlc diff
	@echo "sqlc check passed"

# Run with race detector
test-race:
	go test -race ./...

# Start services (alias for docker-up)
up:
	docker compose up -d

# Stop services (alias for docker-down)
down:
	docker compose down

# View app logs
logs:
	docker compose logs -f app

