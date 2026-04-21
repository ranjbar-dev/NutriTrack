# Technology Stack

**Project:** NutriTrack — Go Backend REST API  
**Researched:** 2025  
**Confidence:** HIGH (all versions verified against proxy.golang.org)

---

## Recommended Stack

### Core Framework

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| `github.com/gin-gonic/gin` | v1.12.0 | HTTP router & middleware | Project requirement; most popular Go HTTP framework, mature ecosystem, excellent middleware support, best-in-class performance for REST APIs |
| `github.com/gin-contrib/cors` | v1.7.7 | CORS middleware | Official Gin contrib, handles preflight for Nuxt.js PWA frontend |
| `github.com/gin-contrib/requestid` | v1.0.6 | Request ID injection | Traces requests through logs end-to-end |

### Database — Primary

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| PostgreSQL | 17-alpine (Docker) | Primary data store | Project requirement; pg_trgm extension for Persian full-text search |
| `github.com/jackc/pgx/v5` | v5.9.2 | PostgreSQL driver | Faster than `lib/pq`, native pgx protocol, required by sqlc v2 config |
| `github.com/sqlc-dev/sqlc` | v1.31.0 | Type-safe SQL code gen (CLI tool) | Project requirement; generates type-safe Go from raw SQL, zero reflection, compile-time safety |
| `github.com/golang-migrate/migrate/v4` | v4.19.1 | DB schema migrations | Project requirement; versioned SQL files, CLI + embedded usage, works perfectly alongside sqlc |

### Database — Cache / Session

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| Redis | 7-alpine (Docker) | OTP, rate limiting, JWT tokens, caching | Project requirement; O(1) TTL-based ops, atomic INCR for rate limiting |
| `github.com/redis/go-redis/v9` | v9.18.0 | Redis client | Official Redis client for Go, context-aware, supports Redis 7 features, pipeline/transactions |

### Authentication & Security

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| `github.com/golang-jwt/jwt/v5` | v5.3.1 | JWT access/refresh tokens | Community standard successor to `dgrijalva/jwt-go`, v5 adds RFC-compliant parsing |
| `golang.org/x/crypto` | v0.50.0 | bcrypt password hashing | Standard library extension; use cost=12 as specified |
| `github.com/google/uuid` | v1.6.0 | JWT `jti` claim + entity IDs | RFC 4122 compliant UUIDs for token blacklisting, idempotency keys |

### Logging

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| `github.com/rs/zerolog` | v1.35.1 | Structured JSON logging | **Use zerolog over zap**: zero-allocation JSON logger, simpler API, smaller binary size, same performance class. Outputs to stdout for Loki collection. |

> **Decision: zerolog over zap** — Both are excellent. zerolog has simpler chaining API (`log.Info().Str("key","val").Msg("...")`), lower allocation overhead at P99, and the `zerolog.SetGlobalLevel()` pattern fits 12-factor app config. zap's sugared logger is ergonomic but adds indirection. For a new project, zerolog is the cleaner choice.

### Validation

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| `github.com/go-playground/validator/v10` | v10.30.2 | Request struct validation | Gin's built-in validator uses this; tag-based validation, custom validators for Persian phone format |

### Configuration

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| `github.com/spf13/viper` | v1.21.0 | Config management | Reads from env vars + `.env` file; 12-factor compliant; `mapstructure` tags for typed config struct |
| `github.com/joho/godotenv` | v1.5.1 | `.env` file loading | Load `.env` in development before viper reads env vars |

### Notifications & SMS

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| `github.com/SherClockHolmes/webpush-go` | v1.4.0 | Web Push (VAPID) notifications | Project requirement; handles VAPID key management and payload encryption automatically |
| SMS (Kavenegar) | HTTP client only | OTP delivery | Use standard `net/http` with adapter interface — do NOT add an SDK dependency. Kavenegar's Go SDK is unmaintained; raw HTTP is more reliable |

### Testing

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| `github.com/stretchr/testify` | v1.11.1 | Test assertions & mocking | Standard Go testing companion; `assert`, `require`, `mock` packages |

### Utilities

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| `github.com/go-playground/validator/v10` | v10.30.2 | Input validation | Already listed above |

---

## Go Version

**Use Go 1.24+** — required for `range` over integers, improved type inference, and `slices`/`maps` packages that simplify domain code.

---

## Alternatives Considered and Rejected

| Category | Recommended | Rejected | Why Rejected |
|----------|-------------|----------|--------------|
| HTTP Framework | Gin v1.12.0 | Fiber, Echo | Project constraint: Gin required |
| ORM/Query | sqlc v1.31.0 | GORM, ent, sqlx | Project constraint: sqlc required; GORM's reflection overhead + magic behavior conflicts with DDD clarity |
| DB Driver | pgx/v5 | lib/pq | `lib/pq` is in maintenance mode; pgx/v5 is faster and required for sqlc v2 `pgx/v5` sql_package |
| Logging | zerolog | zap | Both valid; zerolog wins on API simplicity and zero-alloc chain style |
| Config | viper | envconfig, godotenv alone | Viper handles `.env` + env vars + nested structs in one pass |
| Rate Limiting | Redis INCR+EXPIRE (manual) | ulule/limiter | For this scale, manual Redis INCR is simpler, explicit, and avoids dependency for 4 lines of logic |
| JWT | golang-jwt/v5 | lestrrat-go/jwx | golang-jwt is lighter; jwx is for complex JWKS scenarios not needed here |
| UUID | google/uuid | sonyflake, ULID | UUIDs are standard; sonyflake is for distributed ID generation at scale (overkill) |

---

## DDD Folder Structure

```
nutritrack/
├── cmd/
│   └── server/
│       └── main.go                    # Wire everything together, start HTTP server
│
├── internal/
│   │
│   ├── domain/                        # PURE DOMAIN — no framework imports, no DB imports
│   │   ├── user/
│   │   │   ├── entity.go              # User, UserRole value object
│   │   │   ├── errors.go              # Domain errors (Persian messages)
│   │   │   └── repository.go          # UserRepository interface (port)
│   │   ├── diet/
│   │   │   ├── entity.go              # DietPlan, PlanDay, Meal, MealOption, MealOptionItem
│   │   │   ├── aggregate.go           # DietPlan aggregate root with computed totals
│   │   │   ├── value_object.go        # NutritionFacts, MealTime value objects
│   │   │   └── repository.go          # DietPlanRepository interface
│   │   ├── food/
│   │   │   ├── entity.go              # FoodItem entity
│   │   │   ├── value_object.go        # MeasurementUnit, FoodCategory
│   │   │   └── repository.go          # FoodRepository interface
│   │   ├── tracking/
│   │   │   ├── entity.go              # FoodLog, WaterLog, SleepLog, ExerciseLog, Measurement
│   │   │   └── repository.go
│   │   ├── medication/
│   │   │   ├── entity.go
│   │   │   └── repository.go
│   │   ├── notification/
│   │   │   ├── entity.go              # PushSubscription entity
│   │   │   └── service.go             # NotificationService interface
│   │   └── auth/
│   │       ├── value_object.go        # OTP, Token value objects
│   │       └── service.go             # AuthDomainService interface
│   │
│   ├── application/                   # USE CASES — orchestrates domain, calls infrastructure ports
│   │   ├── auth/
│   │   │   ├── service.go             # SendOTP, VerifyOTP, RefreshToken, Logout
│   │   │   └── dto.go                 # SendOTPRequest, VerifyOTPResponse
│   │   ├── diet/
│   │   │   ├── service.go             # CreateDietPlan, AssignPlan, GetClientPlan
│   │   │   └── dto.go
│   │   ├── food/
│   │   │   ├── service.go             # CreateFood, SearchFood, ApproveRequest
│   │   │   └── dto.go
│   │   ├── tracking/
│   │   │   ├── service.go             # LogFood, LogWater, LogSleep, LogMeasurement
│   │   │   └── dto.go
│   │   ├── user/
│   │   │   ├── service.go             # RegisterClient, UpdateProfile
│   │   │   └── dto.go
│   │   └── notification/
│   │       ├── service.go             # Subscribe, SendReminder
│   │       └── dto.go
│   │
│   ├── infrastructure/                # ADAPTERS — implements domain interfaces
│   │   ├── persistence/
│   │   │   ├── queries/               # Raw .sql files (input to sqlc)
│   │   │   │   ├── users.sql
│   │   │   │   ├── diet_plans.sql
│   │   │   │   ├── food_items.sql
│   │   │   │   └── ...
│   │   │   ├── db/                    # sqlc-generated Go code (DO NOT EDIT)
│   │   │   │   ├── db.go
│   │   │   │   ├── models.go
│   │   │   │   ├── users.sql.go
│   │   │   │   └── ...
│   │   │   ├── user_repository.go     # Implements domain/user.UserRepository
│   │   │   ├── diet_repository.go     # Implements domain/diet.DietPlanRepository
│   │   │   ├── food_repository.go
│   │   │   └── tracking_repository.go
│   │   ├── cache/
│   │   │   ├── redis.go               # Redis client initialization
│   │   │   ├── otp_store.go           # OTP CRUD with TTL
│   │   │   ├── token_store.go         # Refresh token + JWT blacklist
│   │   │   └── rate_limiter.go        # INCR+EXPIRE rate limiting
│   │   ├── sms/
│   │   │   ├── interface.go           # SMSProvider interface
│   │   │   ├── kavenegar.go           # Kavenegar HTTP adapter
│   │   │   └── mock.go                # Mock for testing
│   │   ├── webpush/
│   │   │   └── service.go             # Implements domain/notification.NotificationService
│   │   └── storage/
│   │       └── local.go               # Local filesystem upload (/data/uploads/)
│   │
│   └── interfaces/                    # HTTP LAYER — Gin handlers + middleware
│       ├── http/
│       │   ├── handler/
│       │   │   ├── auth_handler.go
│       │   │   ├── diet_handler.go
│       │   │   ├── food_handler.go
│       │   │   ├── tracking_handler.go
│       │   │   ├── user_handler.go
│       │   │   └── notification_handler.go
│       │   ├── middleware/
│       │   │   ├── auth.go            # JWT extraction + validation
│       │   │   ├── rate_limit.go      # Redis-backed rate limiting
│       │   │   ├── logger.go          # zerolog request logging
│       │   │   ├── recovery.go        # Panic recovery → Persian error
│       │   │   └── role.go            # Role-based access control
│       │   ├── request/               # Input DTOs (bind + validate)
│       │   └── response/              # Output DTOs + Persian error wrapper
│       └── router.go                  # Gin engine setup, route registration
│
├── pkg/                               # SHARED UTILITIES — no business logic
│   ├── apperrors/
│   │   └── errors.go                  # Typed errors with Persian messages
│   ├── pagination/
│   │   └── pagination.go              # Cursor/offset pagination helpers
│   ├── timezone/
│   │   └── tehran.go                  # time.LoadLocation("Asia/Tehran") singleton
│   └── logger/
│       └── logger.go                  # zerolog initialization from config
│
├── migrations/                        # golang-migrate versioned SQL files
│   ├── 000001_create_users.up.sql
│   ├── 000001_create_users.down.sql
│   ├── 000002_create_foods.up.sql
│   └── ...
│
├── sqlc.yaml                          # sqlc configuration
├── Dockerfile
├── docker-compose.yml
├── docker-compose.override.yml        # Local dev overrides
├── .env.example
├── Makefile
└── go.mod
```

---

## sqlc Configuration

**`sqlc.yaml`** — place at project root:

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "internal/infrastructure/persistence/queries"
    schema: "migrations"
    gen:
      go:
        package: "db"
        out: "internal/infrastructure/persistence/db"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_db_tags: true
        emit_interface: true
        emit_exact_table_names: false
        emit_empty_slices: true
        emit_pointers_for_null_types: true
        overrides:
          - db_type: "uuid"
            go_type: "github.com/google/uuid.UUID"
          - db_type: "pg_catalog.numeric"
            go_type: "float64"
```

**Key sqlc decisions:**
- `sql_package: "pgx/v5"` — use pgx directly, not `database/sql` wrapper; needed for pgx-native types and named args
- `emit_pointers_for_null_types: true` — nullable columns become `*string`, `*float64` (cleaner than `sql.NullString`)
- `emit_interface: true` — generates a `Querier` interface usable for mocking in tests
- Schema source is `migrations/` — sqlc reads migration files directly to infer schema

---

## Redis Usage Patterns

### OTP Storage

```go
// Store OTP (2 min TTL)
func (s *OTPStore) Set(ctx context.Context, phone, code string) error {
    key := fmt.Sprintf("otp:%s", phone)
    return s.client.Set(ctx, key, code, 2*time.Minute).Err()
}

// Verify and delete (atomic: get then del)
func (s *OTPStore) Verify(ctx context.Context, phone, code string) (bool, error) {
    key := fmt.Sprintf("otp:%s", phone)
    val, err := s.client.Get(ctx, key).Result()
    if err == redis.Nil {
        return false, nil // expired or not found
    }
    if err != nil {
        return false, err
    }
    if val != code {
        return false, nil
    }
    s.client.Del(ctx, key) // consume OTP — single use
    return true, nil
}
```

### OTP Rate Limiting (max 3 requests per phone per 10 min)

```go
func (r *RateLimiter) CheckOTPLimit(ctx context.Context, phone string) error {
    key := fmt.Sprintf("otp_rl:%s", phone)
    count, err := r.client.Incr(ctx, key).Result()
    if err != nil {
        return err
    }
    if count == 1 {
        // First request in window — set expiry
        r.client.Expire(ctx, key, 10*time.Minute)
    }
    if count > 3 {
        return apperrors.ErrOTPRateLimitExceeded // "تعداد درخواست‌های شما بیش از حد مجاز است"
    }
    return nil
}
```

### OTP Attempt Limiting (max 3 wrong attempts)

```go
func (r *RateLimiter) IncrOTPAttempt(ctx context.Context, phone string) (int64, error) {
    key := fmt.Sprintf("otp_attempts:%s", phone)
    count, err := r.client.Incr(ctx, key).Result()
    if count == 1 {
        r.client.Expire(ctx, key, 2*time.Minute) // same TTL as OTP
    }
    return count, err
}
```

### JWT Blacklisting (logout / token invalidation)

```go
// On logout: blacklist the jti until the token's natural expiry
func (s *TokenStore) BlacklistToken(ctx context.Context, jti string, expiresAt time.Time) error {
    key := fmt.Sprintf("jwt_bl:%s", jti)
    ttl := time.Until(expiresAt)
    if ttl <= 0 {
        return nil // already expired, no need to blacklist
    }
    return s.client.Set(ctx, key, "1", ttl).Err()
}

func (s *TokenStore) IsBlacklisted(ctx context.Context, jti string) (bool, error) {
    key := fmt.Sprintf("jwt_bl:%s", jti)
    exists, err := s.client.Exists(ctx, key).Result()
    return exists > 0, err
}
```

### Refresh Token Storage (30-day TTL, per user+device)

```go
// Key pattern: refresh:{userID}:{deviceID}
func (s *TokenStore) StoreRefreshToken(ctx context.Context, userID, deviceID, token string) error {
    key := fmt.Sprintf("refresh:%s:%s", userID, deviceID)
    return s.client.Set(ctx, key, token, 30*24*time.Hour).Err()
}

// On password change or deactivation: revoke all sessions
func (s *TokenStore) RevokeAllSessions(ctx context.Context, userID string) error {
    pattern := fmt.Sprintf("refresh:%s:*", userID)
    keys, _ := s.client.Keys(ctx, pattern).Result()
    if len(keys) > 0 {
        return s.client.Del(ctx, keys...).Err()
    }
    return nil
}
```

### General Rate Limiting (API endpoint)

```go
func (r *RateLimiter) CheckAPILimit(ctx context.Context, ip, endpoint string, max int, window time.Duration) error {
    key := fmt.Sprintf("api_rl:%s:%s", ip, endpoint)
    count, err := r.client.Incr(ctx, key).Result()
    if err != nil {
        return err
    }
    if count == 1 {
        r.client.Expire(ctx, key, window)
    }
    if count > int64(max) {
        return apperrors.ErrRateLimitExceeded
    }
    return nil
}
```

---

## Docker Compose Configuration

**`docker-compose.yml`** — production-oriented base:

```yaml
version: '3.9'

services:
  postgres:
    image: postgres:17-alpine
    restart: unless-stopped
    environment:
      POSTGRES_DB: ${POSTGRES_DB:-nutritrack}
      POSTGRES_USER: ${POSTGRES_USER:-nutritrack}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      TZ: Asia/Tehran
      PGTZ: Asia/Tehran
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - internal
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER:-nutritrack}"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    restart: unless-stopped
    command: >
      redis-server
      --requirepass ${REDIS_PASSWORD}
      --appendonly yes
      --appendfsync everysec
      --maxmemory 256mb
      --maxmemory-policy allkeys-lru
    environment:
      TZ: Asia/Tehran
    volumes:
      - redis_data:/data
    networks:
      - internal
    healthcheck:
      test: ["CMD", "redis-cli", "--no-auth-warning", "-a", "${REDIS_PASSWORD}", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  app:
    build:
      context: .
      dockerfile: Dockerfile
    restart: unless-stopped
    environment:
      TZ: Asia/Tehran
      APP_ENV: production
      SERVER_PORT: 8080
      DATABASE_URL: postgres://${POSTGRES_USER:-nutritrack}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB:-nutritrack}?sslmode=disable
      REDIS_URL: redis://:${REDIS_PASSWORD}@redis:6379/0
      JWT_SECRET: ${JWT_SECRET}
      JWT_ACCESS_TTL: 15m
      JWT_REFRESH_TTL: 720h
      OTP_TTL_SECONDS: 120
      OTP_MAX_ATTEMPTS: 3
      OTP_RATE_LIMIT_MAX: 3
      OTP_RATE_LIMIT_WINDOW_SECONDS: 600
      UPLOAD_PATH: /data/uploads
      LOG_LEVEL: info
      LOG_FORMAT: json
      KAVENEGAR_API_KEY: ${KAVENEGAR_API_KEY}
      KAVENEGAR_SENDER: ${KAVENEGAR_SENDER}
      VAPID_PUBLIC_KEY: ${VAPID_PUBLIC_KEY}
      VAPID_PRIVATE_KEY: ${VAPID_PRIVATE_KEY}
      VAPID_SUBSCRIBER: ${VAPID_SUBSCRIBER}
      BCRYPT_COST: 12
    volumes:
      - /data/uploads:/data/uploads
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - internal
      - traefik
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.nutritrack.rule=Host(`api.nutritrack.ir`)"
      - "traefik.http.routers.nutritrack.entrypoints=websecure"
      - "traefik.http.routers.nutritrack.tls.certresolver=letsencrypt"
      - "traefik.http.services.nutritrack.loadbalancer.server.port=8080"

volumes:
  postgres_data:
  redis_data:

networks:
  internal:
    driver: bridge
  traefik:
    external: true
```

**`docker-compose.override.yml`** — local development overrides:

```yaml
version: '3.9'

services:
  postgres:
    ports:
      - "5432:5432"

  redis:
    ports:
      - "6379:6379"

  app:
    build:
      target: development
    environment:
      APP_ENV: development
      LOG_LEVEL: debug
      LOG_FORMAT: console
    volumes:
      - .:/app
      - /app/tmp
    ports:
      - "8080:8080"
```

---

## Dockerfile

```dockerfile
# ---- Build Stage ----
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /app/bin/server ./cmd/server

# ---- Runtime Stage ----
FROM alpine:3.21 AS production

RUN apk add --no-cache ca-certificates tzdata
# Timezone data must be in image for Asia/Tehran to work
ENV TZ=Asia/Tehran

WORKDIR /app

COPY --from=builder /app/bin/server .
COPY --from=builder /app/migrations ./migrations

RUN mkdir -p /data/uploads

EXPOSE 8080

ENTRYPOINT ["./server"]
```

> **Critical:** Include `tzdata` in the Alpine image. Without it, `time.LoadLocation("Asia/Tehran")` panics at runtime. The `TZ=Asia/Tehran` env var alone is not enough without the timezone database.

---

## Persian Error Response Convention

All API error responses use a consistent JSON structure with Persian `message` fields:

```go
// pkg/apperrors/errors.go

type AppError struct {
    Code       string `json:"code"`       // Machine-readable: "OTP_EXPIRED"
    Message    string `json:"message"`    // Human-readable Persian: "کد تأیید منقضی شده است"
    StatusCode int    `json:"-"`
}

// Predefined domain errors
var (
    ErrOTPExpired           = &AppError{Code: "OTP_EXPIRED",            Message: "کد تأیید منقضی شده است",               StatusCode: 400}
    ErrOTPInvalid           = &AppError{Code: "OTP_INVALID",            Message: "کد تأیید اشتباه است",                  StatusCode: 400}
    ErrOTPRateLimitExceeded = &AppError{Code: "OTP_RATE_LIMIT",         Message: "تعداد درخواست‌های کد تأیید بیش از حد مجاز است", StatusCode: 429}
    ErrOTPMaxAttempts       = &AppError{Code: "OTP_MAX_ATTEMPTS",       Message: "تعداد تلاش‌های مجاز تمام شده است",      StatusCode: 400}
    ErrUnauthorized         = &AppError{Code: "UNAUTHORIZED",           Message: "دسترسی غیر مجاز",                      StatusCode: 401}
    ErrTokenExpired         = &AppError{Code: "TOKEN_EXPIRED",          Message: "توکن دسترسی منقضی شده است",            StatusCode: 401}
    ErrForbidden            = &AppError{Code: "FORBIDDEN",              Message: "شما مجاز به این عملیات نیستید",        StatusCode: 403}
    ErrNotFound             = &AppError{Code: "NOT_FOUND",              Message: "مورد مورد نظر یافت نشد",              StatusCode: 404}
    ErrValidation           = &AppError{Code: "VALIDATION_ERROR",       Message: "اطلاعات ورودی نامعتبر است",           StatusCode: 422}
    ErrRateLimitExceeded    = &AppError{Code: "RATE_LIMIT_EXCEEDED",    Message: "تعداد درخواست‌های شما بیش از حد مجاز است", StatusCode: 429}
    ErrInternal             = &AppError{Code: "INTERNAL_ERROR",         Message: "خطای داخلی سرور",                     StatusCode: 500}
    ErrDietPlanNotFound     = &AppError{Code: "DIET_PLAN_NOT_FOUND",    Message: "برنامه غذایی یافت نشد",               StatusCode: 404}
    ErrFoodNotFound         = &AppError{Code: "FOOD_NOT_FOUND",         Message: "ماده غذایی یافت نشد",                StatusCode: 404}
    ErrClientNotAssigned    = &AppError{Code: "CLIENT_NOT_ASSIGNED",    Message: "این مراجع به شما تعلق ندارد",         StatusCode: 403}
)
```

HTTP response wrapper in Gin handler:

```go
// interfaces/http/response/response.go

type SuccessResponse struct {
    Data    any    `json:"data"`
    Message string `json:"message,omitempty"`
}

type ErrorResponse struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details any    `json:"details,omitempty"`
}

func OK(c *gin.Context, data any) {
    c.JSON(http.StatusOK, SuccessResponse{Data: data})
}

func Created(c *gin.Context, data any) {
    c.JSON(http.StatusCreated, SuccessResponse{Data: data})
}

func Fail(c *gin.Context, err *apperrors.AppError) {
    c.JSON(err.StatusCode, ErrorResponse{Code: err.Code, Message: err.Message})
    c.Abort()
}

func FailWithDetails(c *gin.Context, err *apperrors.AppError, details any) {
    c.JSON(err.StatusCode, ErrorResponse{Code: err.Code, Message: err.Message, Details: details})
    c.Abort()
}
```

---

## Asia/Tehran Timezone Pattern

```go
// pkg/timezone/tehran.go

package timezone

import (
    "sync"
    "time"
)

var (
    tehranOnce sync.Once
    tehranLoc  *time.Location
)

// Tehran returns the Asia/Tehran timezone location.
// Uses sync.Once so LoadLocation is called only once.
func Tehran() *time.Location {
    tehranOnce.Do(func() {
        var err error
        tehranLoc, err = time.LoadLocation("Asia/Tehran")
        if err != nil {
            // Should never happen if tzdata package is imported or tzdata Alpine package installed
            panic("failed to load Asia/Tehran timezone: " + err.Error())
        }
    })
    return tehranLoc
}

// NowTehran returns the current time in Asia/Tehran.
func NowTehran() time.Time {
    return time.Now().In(Tehran())
}

// ToTehran converts any time.Time to Asia/Tehran.
func ToTehran(t time.Time) time.Time {
    return t.In(Tehran())
}
```

**Timestamp rules:**
- Store all timestamps in **UTC** in PostgreSQL (`TIMESTAMPTZ` columns)
- Convert to Asia/Tehran only at the **API response boundary** when sending to clients
- All TTL calculations (OTP, JWT) use UTC durations — unaffected by timezone
- PostgreSQL session timezone: set `SET timezone = 'UTC'` in connection string or via pgx connect config

---

## zerolog Initialization Pattern

```go
// pkg/logger/logger.go

package logger

import (
    "os"
    "time"

    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)

func Init(level, format string) {
    // Set log level
    lvl, err := zerolog.ParseLevel(level)
    if err != nil {
        lvl = zerolog.InfoLevel
    }
    zerolog.SetGlobalLevel(lvl)

    // Set time format to RFC3339
    zerolog.TimeFieldFormat = time.RFC3339

    if format == "console" {
        // Development: human-readable colored output
        log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
    } else {
        // Production: JSON to stdout for Loki
        log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
    }
}
```

**Request logging middleware:**

```go
// interfaces/http/middleware/logger.go

func RequestLogger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path

        c.Next()

        log.Info().
            Str("method", c.Request.Method).
            Str("path", path).
            Int("status", c.Writer.Status()).
            Dur("latency", time.Since(start)).
            Str("ip", c.ClientIP()).
            Str("request_id", c.GetString("request_id")).
            Msg("request")
    }
}
```

---

## Makefile Commands

```makefile
.PHONY: run build migrate-up migrate-down sqlc lint test docker-up docker-down

# Development
run:
	go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

# Database migrations
migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down 1

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

# sqlc code generation
sqlc:
	sqlc generate

# Testing
test:
	go test ./... -v -race -cover

lint:
	golangci-lint run ./...

# Docker
docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f app

# VAPID key generation (run once)
vapid-keys:
	go run ./cmd/tools/vapid_keygen
```

---

## go.mod Structure

```
module github.com/ranjbar-dev/nutritrack

go 1.24

require (
    github.com/gin-contrib/cors v1.7.7
    github.com/gin-contrib/requestid v1.0.6
    github.com/gin-gonic/gin v1.12.0
    github.com/go-playground/validator/v10 v10.30.2
    github.com/golang-jwt/jwt/v5 v5.3.1
    github.com/golang-migrate/migrate/v4 v4.19.1
    github.com/google/uuid v1.6.0
    github.com/jackc/pgx/v5 v5.9.2
    github.com/joho/godotenv v1.5.1
    github.com/redis/go-redis/v9 v9.18.0
    github.com/rs/zerolog v1.35.1
    github.com/SherClockHolmes/webpush-go v1.4.0
    github.com/spf13/viper v1.21.0
    github.com/stretchr/testify v1.11.1
    golang.org/x/crypto v0.50.0
)
```

---

## Environment Variables Reference

**`.env.example`:**

```bash
# Application
APP_ENV=development
SERVER_PORT=8080
LOG_LEVEL=debug
LOG_FORMAT=console

# Database
POSTGRES_DB=nutritrack
POSTGRES_USER=nutritrack
POSTGRES_PASSWORD=change_me_in_production
DATABASE_URL=postgres://nutritrack:change_me_in_production@localhost:5432/nutritrack?sslmode=disable

# Redis
REDIS_PASSWORD=change_me_in_production
REDIS_URL=redis://:change_me_in_production@localhost:6379/0

# JWT
JWT_SECRET=change_me_minimum_32_chars_for_hmac_sha256
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=720h

# OTP
OTP_TTL_SECONDS=120
OTP_MAX_ATTEMPTS=3
OTP_RATE_LIMIT_MAX=3
OTP_RATE_LIMIT_WINDOW_SECONDS=600

# SMS (Kavenegar)
KAVENEGAR_API_KEY=
KAVENEGAR_SENDER=
KAVENEGAR_TEMPLATE=otp_template

# Web Push (VAPID) — generate with: make vapid-keys
VAPID_PUBLIC_KEY=
VAPID_PRIVATE_KEY=
VAPID_SUBSCRIBER=mailto:admin@nutritrack.ir

# Storage
UPLOAD_PATH=/data/uploads
UPLOAD_MAX_SIZE_MB=10

# Timezone (also set in Docker)
TZ=Asia/Tehran

# Security
BCRYPT_COST=12
```

---

## Sources

- Gin v1.12.0: https://proxy.golang.org/github.com/gin-gonic/gin/@latest — HIGH confidence
- sqlc v1.31.0: https://proxy.golang.org/github.com/sqlc-dev/sqlc/@latest — HIGH confidence
- go-redis v9.18.0: https://proxy.golang.org/github.com/redis/go-redis/v9/@latest — HIGH confidence
- golang-migrate v4.19.1: https://proxy.golang.org/github.com/golang-migrate/migrate/v4/@latest — HIGH confidence
- zerolog v1.35.1: https://proxy.golang.org/github.com/rs/zerolog/@latest — HIGH confidence
- golang-jwt/v5 v5.3.1: https://proxy.golang.org/github.com/golang-jwt/jwt/v5/@latest — HIGH confidence
- webpush-go v1.4.0: https://proxy.golang.org/github.com/!sher!clock!holmes/webpush-go/@latest — HIGH confidence
- pgx/v5 v5.9.2: https://proxy.golang.org/github.com/jackc/pgx/v5/@latest — HIGH confidence
- google/uuid v1.6.0: https://proxy.golang.org/github.com/google/uuid/@latest — HIGH confidence
- validator/v10 v10.30.2: https://proxy.golang.org/github.com/go-playground/validator/v10/@latest — HIGH confidence
- golang.org/x/crypto v0.50.0: https://proxy.golang.org/golang.org/x/crypto/@latest — HIGH confidence
- spf13/viper v1.21.0: https://proxy.golang.org/github.com/spf13/viper/@latest — HIGH confidence
- go.uber.org/zap v1.27.1: https://proxy.golang.org/go.uber.org/zap/@latest — HIGH confidence (rejected in favor of zerolog)
