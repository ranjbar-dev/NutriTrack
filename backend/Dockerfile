# ── Stage 1: Build ──
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache ca-certificates git

WORKDIR /build

# Cache dependency downloads
COPY go.mod go.sum ./
RUN go mod download

# Build static binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app ./cmd/api

# ── Stage 2: Run (~20MB final image) ──
FROM scratch

# TLS certificates for outbound HTTPS (e.g. Kavenegar API)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Application binary
COPY --from=builder /app /app

# Database migrations (golang-migrate runs on startup)
COPY --from=builder /build/db/migrations /db/migrations

EXPOSE 8080

ENTRYPOINT ["/app"]
