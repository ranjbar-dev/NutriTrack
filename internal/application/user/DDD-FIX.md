# DDD Fix Report: internal/application/user
Layer: application
Fixed: 2026-04-22
Based on: DDD-AUDIT.md

## Baseline Build Status
PASS — `go build ./...` before fixes

## Fix Plan

| # | Finding | Severity | Files | Strategy | Status |
|---|---------|----------|-------|----------|--------|
| 1 | Application service imports concrete infrastructure/storage package | CRITICAL | avatar_service.go | SAFE | FIXED |
| 2 | AvatarService field holds concrete infrastructure type | CRITICAL | avatar_service.go | SAFE | FIXED |
| 3 | `entity.User` constructed via raw struct literal (no factory) | HIGH | client_service.go, nutritionist_service.go | DEFERRED | DEFERRED: requires domain entity changes + all callers |
| 4 | Direct mutation of exported aggregate fields | HIGH | avatar_service.go, client_service.go, nutritionist_service.go | DEFERRED | DEFERRED: requires domain method additions + all callers |
| 5 | Cross-application-service import for password hashing | MEDIUM | nutritionist_service.go | DEFERRED | DEFERRED: requires new internal/security package |
| 6 | Missing blank line before ListClients | LOW | client_service.go | DEFERRED | DEFERRED: out of scope per task instructions |
| 7 | Missing blank line before SetStatus | LOW | nutritionist_service.go | DEFERRED | DEFERRED: out of scope per task instructions |

## Changes Applied

### Fix 1 & 2: FileStorage port interface — remove infrastructure import, use domain interface

**New file:** `internal/domain/shared/file_storage.go`
```go
package shared

import "io"

// FileStorage is the domain port for saving files.
type FileStorage interface {
    SaveAvatar(src io.Reader, ext string) (string, error)
    SaveLabResult(src io.Reader, ext string) (string, error)
    SaveAttachment(src io.Reader, ext string) (string, error)
}
```

**File:** `internal/application/user/avatar_service.go`

**Before (imports):**
```go
import (
    "context"
    "io"

    "github.com/google/uuid"
    "github.com/ranjbar-dev/nutritrack/internal/domain/shared"
    "github.com/ranjbar-dev/nutritrack/internal/domain/user/entity"
    userRepo "github.com/ranjbar-dev/nutritrack/internal/domain/user/repository"
    "github.com/ranjbar-dev/nutritrack/internal/infrastructure/storage"  // ← VIOLATION
)
```

**After (imports):**
```go
import (
    "context"
    "io"

    "github.com/google/uuid"
    "github.com/ranjbar-dev/nutritrack/internal/domain/shared"
    "github.com/ranjbar-dev/nutritrack/internal/domain/user/entity"
    userRepo "github.com/ranjbar-dev/nutritrack/internal/domain/user/repository"
)
```

**Before (struct + constructor):**
```go
type AvatarService struct {
    userRepo userRepo.UserRepository
    storage  *storage.LocalStorage           // ← VIOLATION
}

func NewAvatarService(repo userRepo.UserRepository, store *storage.LocalStorage) *AvatarService {
```

**After (struct + constructor):**
```go
type AvatarService struct {
    userRepo userRepo.UserRepository
    storage  shared.FileStorage              // ← domain port interface
}

func NewAvatarService(repo userRepo.UserRepository, store shared.FileStorage) *AvatarService {
```

**Build:** PASS — `*storage.LocalStorage` satisfies `shared.FileStorage` implicitly; `bootstrap/wire.go` unchanged.

## Deferred Items

- **[HIGH] Raw struct literal construction of `entity.User`** (`client_service.go:65`, `nutritionist_service.go:70`) — requires adding `NewClient()` / `NewNutritionist()` factory functions to the domain entity and updating all callers. Broader change than allowed in this single-pass fix.
- **[HIGH] Direct mutation of exported aggregate fields** — requires adding domain mutation methods (`UpdateProfile`, `SetAvatarURL`, etc.) to the entity and replacing all direct field assignments across 3 files. Deferred.
- **[MEDIUM] Cross-application-service password import** — requires extracting `HashPassword`/`CheckPassword` to a new `internal/security/` package. Deferred.
- **[LOW] Missing blank lines** — cosmetic; deferred.

## Final Build Status
PASS — `go build ./...` after all fixes
PASS — `go vet ./internal/...` after all fixes

## Remaining Violations
None of the CRITICAL findings remain. HIGH findings are documented as DEFERRED above.
