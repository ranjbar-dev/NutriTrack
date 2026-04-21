---
phase: "01"
plan: "01"
subsystem: "foundation"
tags: ["go-module", "ddd", "config", "logging", "middleware"]
dependency_graph:
  requires: []
  provides: ["go-module", "ddd-skeleton", "viper-config", "zerolog-logger", "gin-router", "apperror-catalog", "timezone-utils", "persian-normalization"]
  affects: ["all subsequent plans"]
tech_stack:
  added:
    - "github.com/gin-gonic/gin v1.12.0"
    - "github.com/rs/zerolog v1.35.1"
    - "github.com/spf13/viper v1.21.0"
    - "github.com/golang-jwt/jwt/v5 v5.2.2"
    - "github.com/redis/go-redis/v9 v9.18.0"
    - "github.com/jackc/pgx/v5 v5.9.2"
    - "github.com/sqlc-dev/sqlc v1.28.0"
    - "github.com/golang-migrate/migrate/v4 v4.18.3"
    - "github.com/google/uuid v1.6.0"
    - "golang.org/x/crypto v0.38.0"
  patterns:
    - "DDD folder structure (domain/application/infrastructure/interfaces)"
    - "Viper config loading from .env file + env vars"
    - "zerolog structured logging (JSON prod, console dev)"
    - "Persian AppError catalog with ToResponse() method"
    - "Tehran timezone helpers (NowTehran/TodayTehran)"
key_files:
  created:
    - go.mod
    - go.sum
    - .env.example
    - configs/config.go
    - bootstrap/logger.go
    - bootstrap/database.go
    - bootstrap/redis.go
    - cmd/server/main.go
    - internal/domain/shared/apperror.go
    - internal/domain/shared/timeutil.go
    - internal/domain/shared/normalize.go
    - internal/interfaces/http/router/router.go
    - internal/interfaces/http/middleware/request_id.go
    - internal/interfaces/http/middleware/logger.go
    - internal/interfaces/http/middleware/recovery.go
    - internal/interfaces/http/middleware/not_found.go
  modified:
    - .gitignore
    - .env.example
decisions:
  - "Used import _ time/tzdata in cmd/server/main.go to embed timezone data for Alpine containers (Asia/Tehran support)"
  - "go mod tidy upgraded dependencies to latest compatible versions (pgx v5.9.2, gin v1.12.0, zerolog v1.35.1)"
  - "Added .env to .gitignore to prevent local secrets from being committed (Rule 2 - security)"
metrics:
  duration: "~10 minutes"
  completed: "2025-01-01"
  tasks_completed: 15
  files_created: 68
---

# Phase 01 Plan 01: Go Module Init, DDD Folder Skeleton, Viper Config, zerolog Setup Summary

**One-liner:** Go module initialized with full DDD skeleton, Viper+zerolog config/logging infrastructure, Persian AppError catalog, and Tehran timezone utilities for the NutriTrack backend.

## What Was Built

This plan established the complete foundational structure for the NutriTrack Go backend:

1. **Go module** `github.com/ranjbar-dev/nutritrack` with all required dependencies
2. **DDD folder skeleton** — 50+ directories across `domain/`, `application/`, `infrastructure/`, `interfaces/`, `bootstrap/`, `cmd/server/`, `migrations/`, `db/`
3. **Viper configuration** — typed `Config` struct with App/Database/Redis/JWT/SMS sub-configs, `.env` file loading with env var override
4. **zerolog logger** — console writer in development, JSON structured logs in production
5. **pgx/v5 pool bootstrap** — `NewPostgresPool()` with ping verification
6. **Redis client bootstrap** — `NewRedisClient()` with ping verification
7. **cmd/server/main.go** — graceful shutdown with SIGINT/SIGTERM, critical `import _ "time/tzdata"` for Alpine containers
8. **Persian AppError catalog** — 25 error entries with Persian messages, `ToResponse()` method
9. **Tehran timezone utilities** — `NowTehran()`, `TodayTehran()`, `TehranLoc()`, `ToTehran()`
10. **NormalizePersian()** — Arabic Kaf/Yeh to Persian equivalents for consistent pg_trgm search
11. **Gin router skeleton** — `/health` endpoint, `/api/v1` group placeholder, production mode support
12. **Middleware** — RequestID (UUID), structured Logger, Recovery (panic handler), NotFound

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 2 - Security] Added .env to .gitignore**
- **Found during:** Task 4 (creating .env.example)
- **Issue:** `.gitignore` only had `backend/.env` but not `.env` at repository root — the `.env` file we created for local dev would have been committed
- **Fix:** Added `.env` and Go build artifact patterns to `.gitignore`
- **Files modified:** `.gitignore`
- **Commit:** `3e88626`

**2. [Rule 3 - Blocking] Dependency version upgrades during go mod tidy**
- **Found during:** Task 14 (go mod tidy)
- **Issue:** `go mod tidy` upgraded several packages to minimum compatible versions required by the dependency graph (pgx v5.7.5 → v5.9.2 due to Go 1.25 toolchain requirement, gin v1.10.0 → v1.12.0, zerolog v1.33.0 → v1.35.1, go-redis v9.7.3 → v9.18.0, viper v1.20.1 → v1.21.0)
- **Fix:** Allowed go mod tidy to resolve; all APIs are backward compatible
- **Files modified:** `go.mod`, `go.sum`
- **Commit:** `3e88626`

## Success Criteria Verification

- ✅ `go build ./...` succeeds with no errors
- ✅ All DDD folders exist with .gitkeep files (50+ directories)
- ✅ `import _ "time/tzdata"` present in `cmd/server/main.go`
- ✅ Persian AppError catalog has 25 error entries (exceeds 20 requirement)
- ✅ `TodayTehran()` and `NowTehran()` functions exist in `internal/domain/shared/timeutil.go`
- ✅ `NormalizePersian()` exists in `internal/domain/shared/normalize.go`
- ✅ `.env.example` exists with all config keys documented

## Known Stubs

None — this plan is infrastructure/skeleton only; no data flows or UI rendering.

## Self-Check: PASSED

- `configs/config.go` — FOUND
- `bootstrap/logger.go` — FOUND
- `bootstrap/database.go` — FOUND
- `bootstrap/redis.go` — FOUND
- `cmd/server/main.go` — FOUND
- `internal/domain/shared/apperror.go` — FOUND
- `internal/domain/shared/timeutil.go` — FOUND
- `internal/domain/shared/normalize.go` — FOUND
- `internal/interfaces/http/router/router.go` — FOUND
- `internal/interfaces/http/middleware/` — FOUND (4 files)
- Commit `3e88626` — FOUND
