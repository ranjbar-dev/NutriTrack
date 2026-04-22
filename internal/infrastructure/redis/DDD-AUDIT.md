# DDD Audit: internal/infrastructure/redis
Layer: infrastructure
Audited: 2026-04-22
Files reviewed: 1 (token_blacklist.go)

## Summary
- CRITICAL: 0
- HIGH: 1
- MEDIUM: 1
- LOW: 0
- PASS: 0

---

## Findings

### [HIGH] `TokenBlacklist` has no domain interface — application layer imports concrete type directly

**File:** `token_blacklist.go`
**Issue:** The `TokenBlacklist` struct is used directly in `internal/application/auth/` (or similar) without a corresponding domain port interface. The application layer imports the infrastructure package and holds a concrete `*redis.TokenBlacklist` pointer.
**DDD Rule:** Infrastructure implementations MUST be hidden behind domain port interfaces. Application services depend only on the interface.
**Fix:** Define a `TokenBlacklist` port interface in `internal/domain/shared/port/` (or `internal/domain/auth/port/`):
```go
type TokenBlacklist interface {
    Add(ctx context.Context, token string, ttl time.Duration) error
    IsBlacklisted(ctx context.Context, token string) (bool, error)
}
```
The Redis implementation implements this interface. The application service accepts the interface.

---

### [MEDIUM] Business rule constants (`maxOTPAttempts`, `maxOTPRateLimit`) defined in infrastructure

**File:** `token_blacklist.go` (or OTP helper in the same package)
**Issue:** Constants that represent domain/business rules (maximum OTP attempts, rate-limit window) are defined in the infrastructure package. Infrastructure should be configurable, not define policy.
**DDD Rule:** Business rules and thresholds belong in the domain or application layer.
**Fix:** Move these constants to `internal/domain/auth/` (domain rules) or pass them as configuration to the infrastructure adapter's constructor.

---

## Compliant Patterns Found

- Redis client is correctly scoped within the infrastructure package. ✓
- Token TTL is passed externally (not hardcoded per-operation). ✓

## Fix Priority Order
1. **[HIGH]** Define `TokenBlacklist` port interface in domain layer; update application service to use it
2. **[MEDIUM]** Move `maxOTPAttempts` / `maxOTPRateLimit` constants to domain layer or make them constructor parameters
