---
phase: 01-foundation-infrastructure
plan: 04
title: "Go Auth Services, Handlers & Router Wiring"
subsystem: backend
tags: [go, gin, auth, jwt, otp, bcrypt, services, handlers, router, seeder]
dependency_graph:
  requires: [01-01, 01-03]
  provides: [auth-service, user-service, auth-endpoints, admin-endpoints, nutritionist-endpoints, router-wiring, super-admin-seeder]
  affects: [01-05, 02-01, 02-02]
tech_stack:
  added: []
  patterns: [layered-architecture, repository-pattern, dependency-injection, anti-enumeration, token-rotation, sha256-otp-hash, pgtype-conversion]
key_files:
  created:
    - backend/internal/repository/user_repo.go
    - backend/internal/repository/otp_repo.go
    - backend/internal/repository/token_repo.go
    - backend/internal/service/auth_service.go
    - backend/internal/service/user_service.go
    - backend/internal/handler/auth_handler.go
    - backend/internal/handler/admin_handler.go
    - backend/internal/handler/client_handler.go
    - backend/internal/middleware/recovery.go
  modified:
    - backend/cmd/api/main.go
    - backend/cmd/seed/main.go
    - backend/db/queries/refresh_tokens.sql
    - backend/internal/repository/sqlc/querier.go
    - backend/internal/repository/sqlc/refresh_tokens.sql.go
decisions:
  - "Added GetRefreshTokenByHashAny sqlc query for replay detection — existing query filtered revoked=false, blocking theft detection"
  - "Recovery middleware placed in dedicated file (recovery.go) rather than inline in main.go for consistency"
  - "AuthService.GetUserByID added to maintain layered arch — handlers never call repositories directly"
  - "setupLogger returns zerolog.Logger instance for dependency injection into services and middleware"
metrics:
  duration: "9 minutes"
  completed: "2026-04-19T19:43:00Z"
---

# Phase 01 Plan 04: Go Auth Services, Handlers & Router Wiring Summary

Complete auth backend with repository wrappers over sqlc, AuthService (login/OTP/refresh/logout), UserService (nutritionist/client creation), HTTP handlers, full Gin router with middleware groups, and bcrypt-12 Super Admin seeder CLI.

## What Was Done

### Task 1: Repository wrappers and auth/user services
- Created `UserRepository`, `OTPRepository`, `TokenRepository` interface/implementation pairs wrapping sqlc-generated queries
- Repository wrappers handle `pgtype.UUID` ↔ `uuid.UUID` and `pgtype.Text` ↔ `string` conversions
- Created `AuthService` with full dependency injection (repos, SMS sender, JWT secret, logger)
  - `LoginWithPassword`: email lookup → bcrypt compare → token pair creation
  - `RequestOTP`: phone normalization → crypto/rand 6-digit OTP → SHA-256 hash → store → send SMS
  - `VerifyOTP`: active OTP lookup → attempt check/increment → hash compare → user lookup → token pair
  - `RefreshTokens`: SHA-256 hash lookup → replay detection (revoked token → family revocation) → rotation
  - `Logout`: revoke refresh token by hash
  - `GetUserByID`: user lookup for GET /api/auth/me
- Created `UserService` with `CreateNutritionist` (bcrypt cost 12) and `RegisterClient` (nutritionist-initiated)
- Anti-enumeration: same error message for user not found / wrong password / invalid OTP
- OTP generation uses `crypto/rand` (not `math/rand`), hashed with SHA-256 before storage
- **Commit:** `227633c`

### Task 2: Auth, admin, and client HTTP handlers
- Created `AuthHandler` with Login, RequestOTP, VerifyOTP, Refresh, Logout, GetMe methods
- Created `AdminHandler` with CreateNutritionist (POST /api/admin/nutritionists)
- Created `ClientHandler` with RegisterClient (POST /api/nutritionist/clients)
- Cookie helpers: `setAuthCookies` (access: /api, 15min; refresh: /api/auth/refresh, 30d)
- `clearAuthCookies` sets maxAge=-1 on both cookies
- No public client self-registration endpoint exists (AUTH-12, D-05)
- Handlers contain zero business logic — pure HTTP concern layer per D-07
- **Commit:** `714ccf3`

### Task 3: Router wiring in main.go and Super Admin seeder CLI
- Rewrote `main.go` with full dependency injection chain: config → pool → repos → services → handlers → router
- pgxpool configured with explicit pool settings: MaxConns=20, MinConns=5, 1h lifetime, 30min idle, 1min health check
- Router groups: public (/api), authed (/api with Auth), admin (/api/admin with Auth+RoleGuard), nutritionist (/api/nutritionist), client (/api/client)
- Global middleware stack: Recovery → SecurityHeaders → RequestID → Logger → CORS
- Rate limiter on OTP endpoints: 3 requests per 10-minute window
- SMS sender selection: MockSender in development, KavenegarSender in production
- Created `recovery.go` middleware with Persian error message for panic recovery
- Implemented Super Admin seeder CLI in `cmd/seed/main.go`: reads env vars, bcrypt cost 12, idempotent duplicate handling
- Graceful shutdown: SIGINT/SIGTERM with 5-second timeout, pool close
- **Commit:** `3c4468a`

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 3 - Blocking] Added GetRefreshTokenByHashAny sqlc query for replay detection**
- **Found during:** Task 1
- **Issue:** Existing `GetRefreshTokenByHash` query filters `revoked = false`, making it impossible to detect replayed (already-revoked) tokens for family-wide revocation (T-04-05)
- **Fix:** Added `GetRefreshTokenByHashAny` query to `refresh_tokens.sql` (no revoked filter) and regenerated sqlc
- **Files modified:** `backend/db/queries/refresh_tokens.sql`, `backend/internal/repository/sqlc/querier.go`, `backend/internal/repository/sqlc/refresh_tokens.sql.go`
- **Commit:** `227633c`

**2. [Rule 2 - Missing functionality] Recovery middleware**
- **Found during:** Task 3
- **Issue:** Plan specified adding Recovery middleware but it didn't exist in the codebase (Plan 03 created other middleware but not recovery)
- **Fix:** Created `backend/internal/middleware/recovery.go` with `gin.CustomRecovery` returning Persian 500 error
- **Files modified:** `backend/internal/middleware/recovery.go`
- **Commit:** `3c4468a`

**3. [Rule 2 - Missing functionality] AuthService.GetUserByID for /api/auth/me**
- **Found during:** Task 1
- **Issue:** Plan's GetMe handler calls `userRepo.GetByID` directly from handler, violating layered architecture (handler → service → repository). Handlers should not access repositories.
- **Fix:** Added `GetUserByID` method to `AuthService` so the handler calls the service layer, maintaining proper layering per D-07
- **Files modified:** `backend/internal/service/auth_service.go`
- **Commit:** `227633c`

## Verification Results

```
cd backend && go build ./cmd/api/    → PASS
cd backend && go build ./cmd/seed/   → PASS
cd backend && go vet ./...           → PASS
```

## Known Stubs

None — all endpoints are fully wired with business logic.

## Self-Check: PASSED

All 11 key files verified present. All 3 task commits verified in git log (`227633c`, `714ccf3`, `3c4468a`).
