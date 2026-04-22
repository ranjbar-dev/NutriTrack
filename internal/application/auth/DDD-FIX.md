# DDD Fix Report: internal/application/auth
Layer: application
Fixed: 2026-04-22
Based on: DDD-AUDIT.md

## Baseline Build Status
PASS — `go build ./...` before fixes

## Fix Plan

| # | Finding | Severity | Files | Strategy | Status |
|---|---------|----------|-------|----------|--------|
| 1 | Application layer imports `internal/infrastructure/redis` directly | CRITICAL | auth_service.go | SAFE | FIXED |
| 2 | AuthService holds concrete infrastructure type `*redis.TokenBlacklist` | CRITICAL | auth_service.go | SAFE | FIXED |
| 3 | Calls `redis.MaxOTPRateLimit()` / `redis.MaxOTPAttempts()` for business constants | HIGH | auth_service.go | SAFE | FIXED |
| 4 | AuthResponse carries `json:` struct tags in application layer | MEDIUM | dto.go | DEFERRED | DEFERRED: requires HTTP DTO in interfaces layer + mapping update |
| 5 | JWTService cryptographic signing in application layer | LOW | jwt_service.go | DEFERRED | DEFERRED: requires new infrastructure/jwt package + interface |

## Changes Applied

### Fix 1 & 2: TokenBlacklistRepository domain interface — remove redis import, use domain interface

**New file:** `internal/domain/user/repository/token_blacklist_repository.go`
```go
package repository

import (
    "context"
    "time"
)

// TokenBlacklistRepository defines the contract for revoking and checking JWT tokens.
// The Redis adapter implements this interface in internal/infrastructure/redis/.
type TokenBlacklistRepository interface {
    Revoke(ctx context.Context, tokenID string, ttl time.Duration) error
    IsRevoked(ctx context.Context, tokenID string) (bool, error)
}
```

**File:** `internal/application/auth/auth_service.go`

**Before (imports):**
```go
import (
    "context"
    "time"

    "github.com/google/uuid"
    "github.com/rs/zerolog/log"
    "github.com/ranjbar-dev/nutritrack/internal/domain/shared"
    userRepo "github.com/ranjbar-dev/nutritrack/internal/domain/user/repository"
    "github.com/ranjbar-dev/nutritrack/internal/domain/user/valueobject"
    "github.com/ranjbar-dev/nutritrack/internal/infrastructure/redis"  // ← VIOLATION
)

const otpLength = 6
```

**After (imports + constants):**
```go
import (
    "context"
    "time"

    "github.com/google/uuid"
    "github.com/rs/zerolog/log"
    "github.com/ranjbar-dev/nutritrack/internal/domain/shared"
    userRepo "github.com/ranjbar-dev/nutritrack/internal/domain/user/repository"
    "github.com/ranjbar-dev/nutritrack/internal/domain/user/valueobject"
)

const (
    otpLength       = 6
    maxOTPRateLimit = int64(3) // max OTP sends per rate-limit window
    maxOTPAttempts  = int64(3) // max failed OTP attempts before lock
)
```

**Before (struct + constructor):**
```go
type AuthService struct {
    userRepo    userRepo.UserRepository
    otpStore    userRepo.OTPRepository
    blacklist   *redis.TokenBlacklist          // ← VIOLATION
    jwtService  *JWTService
    smsProvider shared.SMSProvider
}

func NewAuthService(
    userRepo userRepo.UserRepository,
    otpStore userRepo.OTPRepository,
    blacklist *redis.TokenBlacklist,           // ← VIOLATION
    ...
```

**After (struct + constructor):**
```go
type AuthService struct {
    userRepo    userRepo.UserRepository
    otpStore    userRepo.OTPRepository
    blacklist   userRepo.TokenBlacklistRepository  // ← domain interface
    jwtService  *JWTService
    smsProvider shared.SMSProvider
}

func NewAuthService(
    userRepo userRepo.UserRepository,
    otpStore userRepo.OTPRepository,
    blacklist userRepo.TokenBlacklistRepository,   // ← domain interface
    ...
```

**Build:** PASS — `*redis.TokenBlacklist` satisfies `userRepo.TokenBlacklistRepository` implicitly; `bootstrap/wire.go` unchanged.

### Fix 3: Replace infrastructure constants with application-layer constants

**Before:**
```go
if count > redis.MaxOTPRateLimit() {   // line 97
    ...
}
if attempts >= redis.MaxOTPAttempts() { // line 131
    ...
}
```

**After:**
```go
if count > maxOTPRateLimit {
    ...
}
if attempts >= maxOTPAttempts {
    ...
}
```

Constants `maxOTPRateLimit = int64(3)` and `maxOTPAttempts = int64(3)` defined at package level in `auth_service.go`. The `redis.MaxOTPRateLimit()` and `redis.MaxOTPAttempts()` exported functions in the infrastructure package are no longer called from the application layer (they remain in the infrastructure package and can be removed or kept for infrastructure-internal use).

**Build:** PASS

## Deferred Items

- **[MEDIUM] `AuthResponse` json tags** (`dto.go:28`) — removing `json:` tags from `AuthResponse` requires adding a mapping HTTP DTO in `internal/interfaces/http/auth/` and updating all JSON marshaling call sites. Deferred to avoid breaking HTTP response serialization.
- **[LOW] JWTService in application layer** (`jwt_service.go`) — moving to `internal/infrastructure/jwt/` requires defining a `TokenService` interface and updating all callers. Deferred.

## Final Build Status
PASS — `go build ./...` after all fixes
PASS — `go vet ./internal/...` after all fixes

## Remaining Violations
None of the CRITICAL or HIGH findings remain unresolved.
