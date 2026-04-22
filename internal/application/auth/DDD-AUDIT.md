# DDD Audit: internal/application/auth
Layer: application
Audited: 2026-04-22
Files reviewed: 4 (auth_service.go, dto.go, jwt_service.go, password_service.go)

## Summary
- CRITICAL: 2
- HIGH: 1
- MEDIUM: 1
- LOW: 1
- PASS: 1 (password_service.go)

---

## Findings

### [CRITICAL] Application layer imports `internal/infrastructure/redis` directly

**File:** `auth_service.go:12`
**Issue:** `auth_service.go` imports `github.com/ranjbar-dev/nutritrack/internal/infrastructure/redis`. The application layer is strictly forbidden from importing any package under `internal/infrastructure`.
**DDD Rule:** `internal/application/` MUST NOT import `internal/infrastructure`. Dependency arrow must always point inward.
**Fix:** Define a `TokenBlacklistRepository` interface in `internal/domain/user/repository/`. Change `AuthService.blacklist` from `*redis.TokenBlacklist` to `userRepo.TokenBlacklistRepository`. Wire the concrete Redis impl at `bootstrap/wire.go`.

---

### [CRITICAL] AuthService holds a concrete infrastructure type as a field

**File:** `auth_service.go:22`
**Issue:** The `blacklist` field is declared as `*redis.TokenBlacklist` — a concrete struct from the infrastructure package. Makes the service impossible to unit-test without a real Redis connection.
**DDD Rule:** Application services MUST accept repository interfaces from `internal/domain`, not concrete implementations from `internal/infrastructure`.
**Fix:** Replace field type with the domain interface `userRepo.TokenBlacklistRepository`. Remove the `redis` import from `auth_service.go`.

---

### [HIGH] Application layer calls infrastructure-package functions for business constants

**File:** `auth_service.go:97, 131`
**Issue:** `auth_service.go` calls `redis.MaxOTPRateLimit()` and `redis.MaxOTPAttempts()` — functions exported by the infrastructure package — to obtain business rule thresholds used in conditionals.
**DDD Rule:** Business rules (thresholds, limits) belong in the domain or application layer, not sourced from infrastructure.
**Fix:** Define constants in the application layer or inject via `AuthServiceConfig` struct. Remove these exports from infrastructure.

---

### [MEDIUM] AuthResponse carries `json:` struct tags in the application layer

**File:** `dto.go:28`
**Issue:** `AuthResponse` is annotated with `json:` tags. JSON serialization is a presentation concern and belongs in the interfaces (HTTP) layer.
**DDD Rule:** DTOs with `json:` tags belong in `internal/interfaces/`. Application layer types should be plain Go structs.
**Fix:** Remove `json:` tags from `AuthResponse`. Create a mapping HTTP DTO in `internal/interfaces/http/auth/`.

---

### [LOW] JWTService implements cryptographic signing in the application layer

**File:** `jwt_service.go`
**Issue:** `JWTService` performs low-level JWT operations (HMAC-SHA256 signing, token parsing). These are infrastructure/adapter concerns, not application-layer orchestration.
**DDD Rule:** Application layer should orchestrate domain logic and delegate technical mechanisms to infrastructure.
**Fix (low priority):** Move `JWTService` to `internal/infrastructure/jwt/`; define a `TokenService` interface for `AuthService` to depend on.

---

## Compliant Patterns Found

- **`password_service.go`** — Clean utility functions with no cross-layer imports. Uses bcrypt with cost factor 12. ✓
- **`AuthService.userRepo` and `AuthService.otpStore`** — Correctly typed as domain repository interfaces. ✓
- **`NewAuthService` factory** — Accepts dependencies via constructor injection. ✓
- **Domain errors** — Consistently returns `shared.Err*` domain errors. ✓
- **`valueobject.NewMobile`** — Correctly delegates mobile validation to domain Value Object. ✓
- **Request types** — Clean plain structs with no json tags. ✓

## Fix Priority Order
1. **[CRITICAL]** Define `TokenBlacklistRepository` interface in domain; replace `*redis.TokenBlacklist` field in `AuthService`
2. **[CRITICAL]** Remove `redis` import from `auth_service.go`
3. **[HIGH]** Replace `redis.MaxOTPRateLimit()` / `redis.MaxOTPAttempts()` with application-layer constants
4. **[MEDIUM]** Strip `json:` tags from `AuthResponse`; add HTTP mapping DTO in interfaces layer
5. **[LOW]** Relocate `JWTService` to `internal/infrastructure/jwt/`; introduce `TokenService` interface
