---
phase: 01-foundation-infrastructure
plan: 01
title: "Monorepo Scaffold & Go Backend Bootstrap"
subsystem: backend
tags: [go, gin, postgresql, sqlc, migrations, scaffold]
dependency_graph:
  requires: []
  provides: [go-module, db-schema, sqlc-generated-repo, config-module, health-endpoint, domain-models, auth-dtos]
  affects: [01-02, 01-03, 01-04]
tech_stack:
  added: [gin@1.12.0, pgx/v5@5.9.2, golang-migrate/v4@4.19.1, zerolog@1.35.0, golang-jwt/v5@5.3.1, validator/v10@10.30.2, uuid@latest]
  patterns: [layered-architecture, env-config, graceful-shutdown, parameterized-sql, pgx-v5-types]
key_files:
  created:
    - backend/go.mod
    - backend/cmd/api/main.go
    - backend/cmd/seed/main.go
    - backend/internal/config/config.go
    - backend/internal/handler/health_handler.go
    - backend/internal/model/user.go
    - backend/internal/model/dto/auth_dto.go
    - backend/internal/validator/validator.go
    - backend/sqlc.yaml
    - backend/db/migrations/000001_create_users.up.sql
    - backend/db/migrations/000002_create_otp_codes.up.sql
    - backend/db/migrations/000003_create_refresh_tokens.up.sql
    - backend/db/queries/users.sql
    - backend/db/queries/otp.sql
    - backend/db/queries/refresh_tokens.sql
    - backend/internal/repository/sqlc/db.go
    - backend/internal/repository/sqlc/models.go
    - backend/internal/repository/sqlc/querier.go
    - backend/.env.example
    - .gitignore
  modified: []
decisions:
  - "Used golang-migrate pgx/v5 driver (not postgres driver) to match pgx pool driver"
  - "Go auto-upgraded from 1.24.2 to 1.25.0 as Gin v1.12.0 requires Go 1.25+"
  - "Placed custom validator in internal/validator package (not inside handler or model)"
metrics:
  duration: "8 minutes"
  completed: "2026-04-19T18:58:00Z"
---

# Phase 01 Plan 01: Monorepo Scaffold & Go Backend Bootstrap Summary

Go backend bootstrapped with Gin HTTP framework, pgx/v5 PostgreSQL driver, 3 migration pairs (users/otp_codes/refresh_tokens), sqlc-generated type-safe repository code, domain models with validation, and health check endpoint at GET /api/health.

## What Was Done

### Task 1: Monorepo directory structure and Go module initialization
- Created full directory structure per D-07: `cmd/api`, `cmd/seed`, `internal/{handler,service,repository/sqlc,model/dto,middleware,config}`, `db/{migrations,queries}`, `pkg/{jwt,sms}`
- Initialized Go module `github.com/ranjbar-dev/nutritrack/backend`
- Go auto-upgraded from 1.24.2 → 1.25.0 (Gin v1.12.0 requires Go 1.25+)
- Installed all dependencies: Gin, pgx/v5, pgxpool, golang-jwt/v5, golang-migrate/v4, validator/v10, zerolog, uuid, x/crypto
- Created `config.go` with environment variable loading and validation (PORT, DATABASE_URL, JWT_SECRET, ENVIRONMENT, FRONTEND_URL, SMS_API_KEY, SMS_TEMPLATE)
- Created `main.go` with: config loading → zerolog setup → pgxpool connection → migrations → Gin server with graceful shutdown
- Used `gin.New()` (not `gin.Default()`) per D-07 — custom middleware added in Plan 03
- Created `.env.example` and root `.gitignore`
- **Commit:** `00c7084`

### Task 2: Database migrations and sqlc code generation
- Created 3 migration pairs (up/down):
  - `000001_create_users`: users table with `user_role`/`gender_type` enums, UUID PKs, partial indexes on mobile/email/nutritionist_id
  - `000002_create_otp_codes`: otp_codes table with code_hash, attempt tracking, expiry
  - `000003_create_refresh_tokens`: refresh_tokens table with family_id for rotation theft detection, ON DELETE CASCADE
- All timestamps use `TIMESTAMPTZ` (never `TIMESTAMP`) per Pitfall 6
- Created `sqlc.yaml` targeting pgx/v5 with `emit_interface`, `emit_empty_slices`, `emit_json_tags`
- Created 3 query files with fully parameterized SQL (SEC-03):
  - `users.sql`: 6 queries (GetByEmail, GetByMobile, GetByID, Create, GetClientsByNutritionist, UpdateActive)
  - `otp.sql`: 5 queries (Create, GetActive, IncrementAttempts, MarkVerified, DeleteExpired)
  - `refresh_tokens.sql`: 5 queries (Create, GetByHash, Revoke, RevokeFamily, RevokeUser)
- Generated type-safe Go code via `sqlc generate` → 6 files in `internal/repository/sqlc/`
- Generated models use pgx/v5 types: `pgtype.UUID`, `pgtype.Timestamptz`, `pgtype.Text`, etc.
- **Commit:** `460dabe`

### Task 3: Domain models, DTOs, and health check handler
- Created `model/user.go` with `UserRole`, `GenderType` enums and `TokenPair` struct
- Created `model/dto/auth_dto.go` with all request/response DTOs:
  - `LoginRequest`, `OTPRequestDTO`, `OTPVerifyDTO`, `CreateNutritionistRequest`, `RegisterClientRequest`
  - `AuthResponse`, `UserResponse`, `ErrorResponse`, `HealthResponse`
  - All DTOs have proper `binding` validation tags for Gin
- Created `internal/validator/validator.go` with `iranian_mobile` custom validator (`^09[0-9]{9}$` per D-25)
- Registered custom validators at startup in `main.go`
- Updated health handler to use `dto.HealthResponse`
- **Commit:** `3184fcf`

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 3 - Blocking] golang-migrate pgx/v5 driver URL scheme**
- **Found during:** Task 1
- **Issue:** The `golang-migrate` pgx/v5 driver uses `pgx5://` URL scheme, but DATABASE_URL uses standard `postgres://`
- **Fix:** Added URL scheme conversion in `runMigrations()` — converts `postgres://` to `pgx5://` for the migrate library while keeping pgxpool using the original URL
- **Files modified:** `backend/cmd/api/main.go`
- **Commit:** `00c7084`

**2. [Rule 3 - Blocking] Go version auto-upgrade**
- **Found during:** Task 1
- **Issue:** Local Go was 1.24.2 but Gin v1.12.0 requires Go 1.25+
- **Fix:** `go get` auto-upgraded `go.mod` to Go 1.25.0 and downloaded the toolchain
- **Files modified:** `backend/go.mod`
- **Commit:** `00c7084`

**3. [Rule 2 - Missing functionality] Custom validator package**
- **Found during:** Task 3
- **Issue:** Plan specified registering Iranian mobile validator in main.go but didn't specify a dedicated package for it
- **Fix:** Created `internal/validator/validator.go` as a clean, reusable package for all custom validators, imported with alias `customvalidator` in main.go
- **Files modified:** `backend/internal/validator/validator.go`, `backend/cmd/api/main.go`
- **Commit:** `3184fcf`

## Verification Results

```
cd backend && go build ./cmd/api/   → PASS
cd backend && go vet ./...          → PASS
cd backend && sqlc generate         → PASS (regenerates cleanly)
cd backend && go build ./...        → PASS (compiles after regen)
```

## Self-Check: PASSED

All 23 key files verified present. All 3 task commits verified in git log.
