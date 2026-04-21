.PHONY: build run dev docker-up docker-down docker-build migrate-up migrate-down test lint

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
