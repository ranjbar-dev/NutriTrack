# DDD Audit: internal/interfaces/http
Layer: interfaces
Audited: 2026-04-22
Files reviewed: 26

## Summary
- CRITICAL: 0
- HIGH: 11
- MEDIUM: 4
- LOW: 4
- PASS: 10 files (handler/admin_handler.go, handler/notification_handler.go, handler/push_handler.go, middleware/cors.go, middleware/error_handler.go, middleware/logger.go, middleware/not_found.go, middleware/recovery.go, middleware/request_id.go, dto/response.go)

---

## Findings

### [HIGH-01] Nine handlers import domain entity packages directly for response mapping

**Files:** `handler/client_handler.go:9`, `handler/nutritionist_handler.go:8`, `handler/food_handler.go:7`, `handler/medication_handler.go:9`, `handler/food_request_handler.go:9`, `handler/lab_result_handler.go:10`, `handler/message_handler.go:11`, `handler/tracking_handler.go:10`, `handler/diet_plan_handler.go:11`
**Issue:** Each handler imports a domain entity package and defines `toXxxResponse()` / `XxxToMap()` functions that accept raw domain entity pointers and manually read their exported fields to build `gin.H` / `map[string]any` responses. Notable escalations: `client_handler.go` calls domain methods `u.FullName()` and `u.BMI()` inside the HTTP layer; `diet_plan_handler.go` traverses the full aggregate tree from the interfaces layer.
**DDD Rule:** Handlers MUST NOT import domain entity packages. Application services should return purpose-built response DTOs.
**Fix:** Define response structs in `internal/application/<domain>/` with `json:` tags. Application services return them. Handlers need no entity imports.

---

### [HIGH-02] `middleware/rate_limit.go` accepts concrete `*redis.Client` — infrastructure coupling in interfaces layer

**File:** `middleware/rate_limit.go:27`
**Issue:** `RateLimitByIP(rdb *redis.Client, max int64)` accepts a concrete infrastructure type. The interfaces layer must not depend on concrete infrastructure.
**DDD Rule:** `internal/interfaces/` MUST NOT import infrastructure directly.
**Fix:** Define a `RateLimiter` interface locally (mirroring `TokenRevocationChecker` pattern):
```go
type RateLimiter interface {
    Allow(ctx context.Context, key string, max int64, window time.Duration) (bool, error)
}
func RateLimitByIP(limiter RateLimiter, max int64) gin.HandlerFunc { ... }
```

---

### [HIGH-03] `router/router.go` accepts raw infrastructure types — interfaces layer imports infrastructure

**File:** `router/router.go:15`
**Issue:** `func New(db *pgxpool.Pool, rdb *redis.Client, cfg *configs.Config)` imports `pgxpool` and `go-redis` directly. Router serves as infrastructure passthrough.
**DDD Rule:** `internal/interfaces/` MUST NOT import infrastructure directly.
**Fix:** Move infra wiring upstream; have router accept `*bootstrap.Container` instead of raw infra types.

---

### [HIGH-04 through HIGH-09] Six individual handlers with direct domain entity imports

Covered under HIGH-01 for full details. Individual files: `handler/food_request_handler.go`, `handler/lab_result_handler.go`, `handler/message_handler.go`, `handler/tracking_handler.go` (six entity types), `handler/diet_plan_handler.go` (full aggregate traversal), `handler/medication_handler.go`.

---

### [HIGH-10] `client_handler.go` executes domain logic (`FullName()`, `BMI()`) in the HTTP layer

**File:** `handler/client_handler.go`
**Issue:** Calls `u.FullName()` and `u.BMI()` (domain entity methods) inside the HTTP handler's response mapping function. Domain logic executing at the HTTP boundary.
**Fix:** Application service computes these values and includes them in the response DTO.

---

### [HIGH-11] `diet_plan_handler.go` traverses full domain aggregate tree from interfaces layer

**File:** `handler/diet_plan_handler.go`
**Issue:** Handler traverses `DietPlan.Days[].Meals[].Options[].Items[]` from the HTTP layer. Complete knowledge of domain aggregate internals at the transport boundary.
**Fix:** Covered by HIGH-01 fix — application service returns flattened DTO.

---

## Medium Findings

### [MEDIUM-01] `avatar_handler.go` reads wrong gin context keys — **silent identity bypass (security bug)**

**File:** `handler/avatar_handler.go:35–36`
**Issue:** Reads `"userID"` and `"userRole"` but `RequireAuth` stores under `middleware.AuthUserIDKey = "auth_user_id"` and `middleware.AuthUserRoleKey = "auth_user_role"`. Both `c.Get()` calls return nil; callerID is `uuid.UUID{}` (all zeros) and callerRole is `""`, effectively bypassing caller-identity enforcement on every avatar upload.
**Fix (CRITICAL to apply immediately):**
```go
callerIDRaw, _ := c.Get(middleware.AuthUserIDKey)
callerRoleRaw, _ := c.Get(middleware.AuthUserRoleKey)
```

---

### [MEDIUM-02] RBAC role checks duplicated in handler bodies with hardcoded role strings

**Files:** `handler/food_request_handler.go:34,58,83,130`, `handler/message_handler.go:38,95,137,163`
**Issue:** Inline role guards use hardcoded literals (`"client"`, `"nutritionist"`) instead of `middleware.RequireRole()` applied at route registration.
**Fix:** Apply `middleware.RequireRole()` at route registration; remove inline checks from handler bodies.

---

### [MEDIUM-03] Request bodies defined as anonymous inline structs

**Files:** `handler/auth_handler.go`, `handler/client_handler.go`, `handler/nutritionist_handler.go`, `handler/food_request_handler.go`, `handler/message_handler.go`, `handler/push_handler.go`, and others.
**Issue:** Anonymous one-off request body structs cannot be reused, tested, or documented.
**Fix:** Define named request types in `dto/` (or domain-scoped sub-packages).

---

### [MEDIUM-04] Magic number `24*time.Hour` for token revocation TTL hardcoded in handler

**File:** `handler/auth_handler.go:146`
**Issue:** TTL is a business rule that should match configured token lifetime; hardcoded in the HTTP handler.
**Fix:** Pass JTI to `AuthService.LogoutWithJTI()`; service computes TTL from its own config.

---

## Low Findings

### [LOW-01] `toAppError()` defined only in `food_handler.go` — error conversion duplicated across handlers

**Fix:** Move `toAppError()` to `dto/errors.go` (exported) or keep as shared unexported helper in `handler` package.

### [LOW-02] `dto/pagination.go` — dual implementation for the same concern

**Issue:** `PaginationQuery` struct + `ParsePagination()` function do the same thing; used inconsistently.
**Fix:** Remove `ParsePagination()`; standardise on `c.ShouldBindQuery(&pg)` + `pg.Normalize()`.

### [LOW-03] `food_category_handler.go` — implicit domain entity field access without explicit import

**Issue:** Accesses exported entity fields via inferred return type. Resolved automatically when HIGH-01 is applied.

### [LOW-04] `diet_plan_handler.go` — missing `gofmt` indentation

**Fix:** Run `gofmt -w internal/interfaces/http/handler/diet_plan_handler.go` before structural changes.

---

## Compliant Patterns Found

- **`middleware/auth.go` — `TokenRevocationChecker` interface** correctly avoids direct infrastructure import. The pattern HIGH-02/03 should adopt. ✓
- **`middleware/RequireRole()` + group-level application in `router.go`** — correct DDD boundary for RBAC. ✓
- **`middleware/error_handler.go` — centralised error pipeline** — `c.Error()` / `dto.Abort()` pattern correct. ✓
- **Typed context key constants** — `AuthUserIDKey`, `AuthUserRoleKey`, `AuthTokenJTIKey` defined once (broken only in `avatar_handler.go`). ✓
- **`handler/admin_handler.go`, `notification_handler.go`, `push_handler.go`** — no domain entity imports. ✓
- **`handler/food_handler.go` — `callerContext()` helper** — good pattern for extracting auth context; should be promoted to shared utility. ✓

## Fix Priority Order

1. **SECURITY — fix immediately**: MEDIUM-01 — `avatar_handler.go` wrong context keys; callerID is always zero UUID
2. **[HIGH-01]** Introduce application-layer response DTOs; remove all domain/entity imports from handler package (9 files; start with `user`)
3. **[HIGH-02]** Replace `*redis.Client` in `RateLimitByIP` with a `RateLimiter` interface
4. **[HIGH-03]** Refactor `router.New()` to accept `*bootstrap.Container` instead of raw infra types
5. **[MEDIUM-02]** Move RBAC role checks to route-level `middleware.RequireRole()` in `router.go`
6. **[MEDIUM-03]** Extract anonymous request structs into named DTO types in `dto/`
7. **[MEDIUM-04]** Move token revocation TTL into `AuthService`
8. **[LOW-04]** Run `gofmt` on `diet_plan_handler.go`
9. **[LOW-01]** Move `toAppError()` to `dto/errors.go`
10. **[LOW-02]** Remove `ParsePagination()`; standardise on struct binding
