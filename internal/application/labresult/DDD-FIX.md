# DDD Fix Report: internal/application/labresult
Layer: application
Fixed: 2026-04-22
Based on: DDD-AUDIT.md

## Baseline Build Status
PASS — `go build ./...` before fixes

## Fix Plan

| # | Finding | Severity | Files | Strategy | Status |
|---|---------|----------|-------|----------|--------|
| 1 | Application layer imports concrete infrastructure package | CRITICAL | lab_result_service.go | SAFE | FIXED |
| 2 | Service struct holds concrete `*storage.LocalStorage` instead of interface | HIGH | lab_result_service.go | SAFE | FIXED |
| 3 | Domain entity constructed via bare struct literal — no factory | MEDIUM | lab_result_service.go | DEFERRED | DEFERRED: LabResult entity has exported fields; factory refactor is cross-layer |
| 4 | Role-based access control uses raw string comparisons | MEDIUM | lab_result_service.go | DEFERRED | DEFERRED: cross-cutting shared.Role change |
| 5 | `UploadLabResult` backward-compat wrapper | LOW | lab_result_service.go | DEFERRED | DEFERRED: callers to audit separately |

## Changes Applied

### Fix 1 + 2: Define `LabResultStorage` interface in shared domain; replace concrete type in service
**File:** `internal/domain/shared/storage.go` (new — shared with message AttachmentStorage)
**Change:** Added `LabResultStorage` port interface:
```go
type LabResultStorage interface {
    SaveLabResult(src io.Reader, ext string) (string, error)
}
```
`*storage.LocalStorage` satisfies this interface implicitly.
**Build:** PASS

**File:** `internal/application/labresult/lab_result_service.go`
**Before:**
```go
import "github.com/ranjbar-dev/nutritrack/internal/infrastructure/storage"

type LabResultService struct {
    repo     labRepo.LabResultRepository
    userRepo userRepo.UserRepository
    storage  *storage.LocalStorage
}

func NewLabResultService(
    repo labRepo.LabResultRepository,
    userRepo userRepo.UserRepository,
    storage *storage.LocalStorage,
) *LabResultService { ... }
```
**After:**
```go
// infrastructure/storage import removed

type LabResultService struct {
    repo     labRepo.LabResultRepository
    userRepo userRepo.UserRepository
    storage  shared.LabResultStorage
}

func NewLabResultService(
    repo labRepo.LabResultRepository,
    userRepo userRepo.UserRepository,
    storage shared.LabResultStorage,
) *LabResultService { ... }
```
**Build:** PASS

No changes required to `bootstrap/wire.go` — `*storage.LocalStorage` satisfies `shared.LabResultStorage` implicitly via structural typing.

## Deferred Items
- **[MEDIUM]** `entity.LabResult` built as bare struct literal in `SubmitLabResult` (lines 142–151). Fixing requires adding a `NewLabResult()` factory and making fields unexported — a multi-file refactor touching infrastructure mappers and the interface layer. Deferred.
- **[MEDIUM]** Role string literals (`"superadmin"`, `"nutritionist"`, `"client"`) — define `shared.Role` type. Deferred as cross-cutting change.
- **[LOW]** `UploadLabResult` backward-compat wrapper — mark as `// Deprecated:` or remove after auditing callers. Deferred.

## Final Build Status
PASS — `go build ./...` after all fixes
PASS — `go vet ./internal/...` after all fixes

## Remaining Violations
- MEDIUM: LabResult struct literal construction — deferred (entity refactor required)
- MEDIUM: Untyped role string comparisons — deferred (cross-cutting)
- LOW: UploadLabResult compatibility wrapper — deferred
