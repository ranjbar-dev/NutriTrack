# ── Stage 1: builder ──────────────────────────────────────────
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git (needed by go mod download for private modules if any)
RUN apk add --no-cache git

# Download dependencies first (layer cache optimization)
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/server \
    ./cmd/server

# ── Stage 2: final ────────────────────────────────────────────
FROM alpine:3.21

# tzdata is REQUIRED for time.LoadLocation("Asia/Tehran") in Alpine
RUN apk add --no-cache tzdata ca-certificates

# Set timezone
ENV TZ=Asia/Tehran

WORKDIR /app

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy binary from builder
COPY --from=builder /app/server .

# Ensure correct ownership
RUN chown appuser:appgroup /app/server

# Create uploads directory owned by appuser so uploads work at runtime
RUN mkdir -p /uploads && chown appuser:appgroup /uploads

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=15s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["/app/server"]
