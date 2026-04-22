# DDD Fix Report: internal/infrastructure/redis
Layer: infrastructure
Fixed: 2026-04-22
Based on: DDD-AUDIT.md

## Baseline Build Status
PASS — `go build ./...` before fixes

## Fix Plan

| # | Finding | Severity | Files | Strategy | Status |
|---|---------|----------|-------|----------|--------|
| 1 | `TokenBlacklist` has no domain interface — app layer imports concrete type directly | HIGH | token_blacklist.go, domain/user/repository/token_blacklist_repository.go | SAFE | ALREADY FIXED |
| 2 | Business rule constants (`maxOTPAttempts`, `maxOTPRateLimit`) defined in infrastructure | MEDIUM | token_blacklist.go | DEFERRED | DEFERRED: constants not present in current file; no violation found |

## Changes Applied

### Fix 1: TokenBlacklist domain interface — ALREADY FIXED by prior agent

**Status:** ALREADY FIXED — no changes required in this pass.

**Verification:**
- `internal/domain/user/repository/token_blacklist_repository.go` defines `TokenBlacklistRepository` interface with `Revoke` and `IsRevoked` methods.
- `internal/infrastructure/redis/token_blacklist.go` `TokenBlacklist` struct implements both methods with matching signatures — satisfies the interface implicitly.
- `internal/application/auth/auth_service.go` field `blacklist` is typed as `userRepo.TokenBlacklistRepository` (the domain interface), not the concrete `*redis.TokenBlacklist`.
- The application layer has zero imports of `internal/infrastructure/redis`.

**Interface (domain/user/repository/token_blacklist_repository.go):**
```go
type TokenBlacklistRepository interface {
    Revoke(ctx context.Context, tokenID string, ttl time.Duration) error
    IsRevoked(ctx context.Context, tokenID string) (bool, error)
}
```

**Implementation (infrastructure/redis/token_blacklist.go):**
```go
func (b *TokenBlacklist) Revoke(ctx context.Context, tokenID string, ttl time.Duration) error { ... }
func (b *TokenBlacklist) IsRevoked(ctx context.Context, tokenID string) (bool, error) { ... }
```

## Deferred Items

- **[MEDIUM]** Business rule constants (`maxOTPAttempts`, `maxOTPRateLimit`): Not present in `token_blacklist.go`; the audit references a possible OTP helper — `otp_store.go` exists but the medium finding does not block any HIGH fix. Deferred to a future focused pass on `otp_store.go`.

## Final Build Status
PASS — `go build ./...` after all fixes
PASS — `go vet ./internal/...` after all fixes

## Remaining Violations
None at CRITICAL or HIGH severity.
