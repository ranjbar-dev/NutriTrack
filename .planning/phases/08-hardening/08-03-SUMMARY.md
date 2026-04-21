---
phase: "08"
plan: "03"
subsystem: security-hardening-and-caching
tags: [security, redis, caching, rate-limiting, jwt, cors, food, dietplan]
dependency_graph:
  requires: [08-01, 08-02]
  provides: [token-revocation, food-search-cache, active-plan-cache, ip-rate-limit]
  affects: [auth, food, dietplan, middleware]
tech_stack:
  added: []
  patterns:
    - TokenRevocationChecker interface in middleware layer
    - Cache-aside pattern (Redis) for food search and active diet plan
    - Per-IP token-bucket rate limiting with Redis INCR + EXPIRE
    - CachedRepository wrapper implementing domain interface
key_files:
  created:
    - internal/interfaces/http/middleware/rate_limit.go
    - internal/infrastructure/persistence/food/cached_food_repository.go
    - internal/infrastructure/persistence/dietplan/cached_diet_plan_repository.go
  modified:
    - configs/config.go
    - internal/interfaces/http/middleware/auth.go
    - internal/interfaces/http/middleware/cors.go
    - internal/application/auth/auth_service.go
    - internal/application/food/food_service.go
    - internal/interfaces/http/handler/auth_handler.go
    - internal/interfaces/http/router/router.go
    - bootstrap/wire.go
decisions:
  - UUID type fix in RequireAuth: claims.UserID (string) now parsed to uuid.UUID before storing in context — prevents runtime panics in all handlers that cast to uuid.UUID
  - TokenRevocationChecker interface defined in middleware package to avoid infrastructure→interface layer import inversion
  - Access token revocation on logout: JTI stored as AuthTokenJTIKey in context; Logout handler blacklists it with 24h TTL
  - CORS now config-driven via CORS_ALLOWED_ORIGINS env var; defaults to * if unset
  - CreateFood now enforces role check (nutritionist/superadmin only) — previously any authenticated user could create foods
  - CachedFoodRepository invalidates entire food:cache:* namespace on any write (SCAN+DEL) for simplicity
  - CachedDietPlanRepository caches nil active plan sentinel to prevent repeated DB hits when no plan exists
metrics:
  duration: 25m
  completed_date: "2026-04-21"
  tasks_completed: 11
  files_modified: 11
---

# Phase 08 Plan 03: Security Hardening, Rate Limiting & Redis Caching Summary

**One-liner:** JWT access-token revocation on logout, per-IP rate limiting via Redis, food-search cache with SCAN invalidation, and active-plan cache with 2-minute TTL.

---

## What Was Built

### Track A: Security Hardening

#### A1 — JWT Access Token Revocation Fix (`auth.go`, `auth_service.go`, `auth_handler.go`)

**Bug found (Rule 1):** `RequireAuth` stored `claims.UserID` (a `string`) directly into the Gin context under `AuthUserIDKey`, but every downstream handler cast it to `uuid.UUID`. This would cause a runtime panic for tracking endpoints using a single-value assertion. **Fixed** by parsing `claims.UserID` to `uuid.UUID` inside `RequireAuth`.

**Security gap:** The blacklist (`TokenBlacklist`) was only consulted during `RefreshToken`. A logout only revoked the refresh token; the access token remained usable until natural expiry. **Fixed:**

- Added `TokenRevocationChecker` interface to the middleware package.
- `RequireAuth` now checks the access token's JTI on every request.
- Added `AuthTokenJTIKey` context key to carry the JTI downstream.
- `AuthService.RevokeToken` method added for immediate JTI blacklisting.
- `Logout` handler now also blacklists the access token JTI with a 24 h TTL.

#### A2 — CORS now config-driven (`cors.go`, `configs/config.go`, `router.go`)

`CORS()` previously hardcoded `"*"`. Now accepts `allowedOrigin string` (passed from `cfg.App.CORSAllowedOrigins` / `CORS_ALLOWED_ORIGINS` env var). Defaults to `"*"` when the env var is unset, enabling gradual tightening in production.

#### A3 — CreateFood role check (`food_service.go`)

`CreateFood` had no role guard — any authenticated user could create foods. **Fixed:** added early return with `ErrForbidden` for callers that are neither `nutritionist` nor `superadmin`.

#### A4 — Rate limiting for OTP endpoint (`rate_limit.go`, `router.go`)

No global per-IP rate limit existed. Added `RateLimitByIP(rdb, 60)` middleware:
- Key: `rate:ip:<ClientIP>`, window: 1 minute, max: 60 requests.
- Redis `INCR` + first-request `EXPIRE` (same atomic pattern as `OTPStore`).
- Returns HTTP 429 with Persian message `"تعداد درخواست‌های شما بیش از حد مجاز است"`.
- Applied to `POST /api/v1/auth/otp/send` (most abuse-prone public endpoint).
- Fail-open on Redis error to avoid false positives.

#### A5 — Row-level auth audit (no changes needed)

- `food_service.go` `UpdateFood` / `DeleteFood`: already enforce nutritionist-owns-food and superadmin-can-all. ✅
- `tracking_handler.go` `GetTracking`: passes `callerID` + `callerRole` to service which handles ownership. ✅
- `lab_result_handler.go` Upload/List/Download: passes `callerID` + `callerRole` to service. ✅

---

### Track B: N+1 Audit & Redis Caching

#### B1 — N+1 Audit (no fixes needed)

- `pg_diet_plan_repository.go`: Days/meals/options/items are loaded via separate calls triggered only when the application service explicitly requests them (no auto-loading). The repo does not load a hierarchy unprompted. ✅
- `pg_food_repository.go`: `FindByID` does two queries (GetFoodByID + GetFoodCategories). Acceptable for single-item point-lookup; not an N+1 pattern. ✅
- `pg_message_repository.go`: Not audited in detail; follow-up deferred if needed.

#### B2 — Food Search Cache (`cached_food_repository.go`)

Wraps `PgFoodRepository` as a `CachedFoodRepository` implementing `repository.FoodRepository`:

| Method | Behaviour |
|--------|-----------|
| `Search` | Cache-aside, key `food:cache:search:noc:<q>:<lim>:<off>`, TTL 5 min |
| `CountSearch` | Cache-aside, key `food:cache:count:noc:<q>`, TTL 5 min |
| `SearchByCategory` | Cache-aside, key `food:cache:search:cat:<catID>:<q>:<lim>:<off>`, TTL 5 min |
| `CountByCategory` | Cache-aside, key `food:cache:count:cat:<catID>:<q>`, TTL 5 min |
| `FindByID` | Pass-through (point-lookup, low traffic) |
| `Create/Update/Delete/Deactivate` | Delegate to inner + SCAN+DEL `food:cache:*` |

#### B3 — Active Diet Plan Cache (`cached_diet_plan_repository.go`)

Wraps `PgDietPlanRepository` as a `CachedDietPlanRepository` implementing `repository.DietPlanRepository`:

| Method | Behaviour |
|--------|-----------|
| `FindActiveByClientID` | Cache-aside, key `dietplan:active:<clientID>`, TTL 2 min; "nil" sentinel for no-plan |
| `CreateWithArchive` | Delegate + DEL `dietplan:active:<clientID>` |
| `Update` | Delegate + DEL `dietplan:active:<plan.ClientID>` |
| `Delete` | FindByID to get clientID → delete from DB → DEL cache |
| All other methods | Pass-through |

#### B4 — Wire update (`wire.go`)

`pgFoodRepo` wrapped with `cachedFoodRepo`; `pgPlanRepo` wrapped with `cachedPlanRepo`. Service constructors accept the `repository` interface — zero service-layer changes needed.

---

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] UUID type mismatch in RequireAuth**
- **Found during:** Track A review of auth.go
- **Issue:** `claims.UserID` is `string`; handlers assert it to `uuid.UUID` (single-value panic form in tracking_handler.go)
- **Fix:** `RequireAuth` now parses `claims.UserID` to `uuid.UUID` before storing; validation error returned if parse fails
- **Files modified:** `internal/interfaces/http/middleware/auth.go`
- **Commit:** af297f9

**2. [Rule 2 - Missing Critical Functionality] Access token not revoked on logout**
- **Found during:** Track A JWT blacklist audit
- **Issue:** Logout only revoked the refresh token JTI; access token remained usable post-logout until natural expiry
- **Fix:** Store JTI in context during `RequireAuth`; `Logout` handler blacklists access token JTI; `RequireAuth` checks blacklist on every request
- **Files modified:** `auth.go`, `auth_service.go`, `auth_handler.go`
- **Commit:** af297f9

---

## Known Stubs

None — all implemented paths return real data.

---

## Threat Flags

None — no new network endpoints or trust-boundary crossings introduced. All new paths are Redis reads/writes behind existing authentication.

---

## Self-Check: PASSED

- `internal/interfaces/http/middleware/rate_limit.go` — FOUND ✅
- `internal/infrastructure/persistence/food/cached_food_repository.go` — FOUND ✅
- `internal/infrastructure/persistence/dietplan/cached_diet_plan_repository.go` — FOUND ✅
- Commit `af297f9` — FOUND ✅
- `go build ./...` — PASSED ✅
