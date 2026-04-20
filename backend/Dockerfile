# ── Stage 1: Build ──
FROM golang:1.25.9-alpine AS builder

RUN apk add --no-cache ca-certificates git

WORKDIR /build

# Cache dependency downloads
COPY go.mod go.sum ./
RUN go mod download

# Build static binaries
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /seed ./cmd/seed

# ── Stage 2: Run (~10MB final image) ──
FROM alpine:3.20

# TLS certificates for outbound HTTPS (e.g. Kavenegar API)
RUN apk add --no-cache ca-certificates wget

# Application binary
COPY --from=builder /app /app

# Seed binary (run once to create super admin)
COPY --from=builder /seed /seed

# Database migrations (golang-migrate runs on startup)
COPY --from=builder /build/db/migrations /db/migrations

EXPOSE 8080

ENTRYPOINT ["/app"]
