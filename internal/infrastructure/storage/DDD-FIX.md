# DDD Fix Report: internal/infrastructure/storage
Layer: infrastructure
Fixed: 2026-04-22
Based on: DDD-AUDIT.md

## Baseline Build Status
PASS — `go build ./...` before fixes

## Fix Plan

| # | Finding | Severity | Files | Strategy | Status |
|---|---------|----------|-------|----------|--------|
| 1 | No domain port interface for storage — app services import concrete type directly | HIGH | local_storage.go, domain/shared/ | SAFE | ALREADY FIXED |
| 2 | `NewLocalStorage` factory has no input validation | MEDIUM | local_storage.go | DEFERRED | DEFERRED: signature change requires callers update |
| 3 | `LocalStorage` has exported fields (`BasePath`, `BaseURL`) | LOW | local_storage.go | SAFE | FIXED |

## Changes Applied

### Fix 1: Storage domain port interfaces — ALREADY FIXED by prior agent

**Status:** ALREADY FIXED — no changes required in this pass.

**Verification:**
- `internal/domain/shared/file_storage.go` defines `FileStorage` interface with `SaveAvatar`, `SaveLabResult`, `SaveAttachment`.
- `internal/domain/shared/storage.go` defines `AttachmentStorage` (with `SaveAttachment`) and `LabResultStorage` (with `SaveLabResult`).
- `internal/application/user/avatar_service.go` field `storage` is typed as `shared.FileStorage`.
- `internal/application/message/message_service.go` field `storage` is typed as `shared.AttachmentStorage`.
- `internal/application/labresult/lab_result_service.go` field `storage` is typed as `shared.LabResultStorage`.
- `LocalStorage` implements all three interfaces via its existing methods.
- No application service imports `internal/infrastructure/storage` directly.

---

### Fix 3: Make `LocalStorage` fields unexported

**File:** `internal/infrastructure/storage/local_storage.go`
**Severity:** LOW
**Build:** PASS

**Before:**
```go
type LocalStorage struct {
    BasePath string // e.g. "uploads"
    BaseURL  string // e.g. "/uploads" or "https://..."
}

func NewLocalStorage(basePath, baseURL string) *LocalStorage {
    return &LocalStorage{BasePath: basePath, BaseURL: baseURL}
}

func (s *LocalStorage) SaveAvatar(src io.Reader, ext string) (string, error) {
    dir := filepath.Join(s.BasePath, "avatars")
    ...
    return fmt.Sprintf("%s/avatars/%s", s.BaseURL, filename), nil
}
// (similar for SaveLabResult and SaveAttachment)
```

**After:**
```go
type LocalStorage struct {
    basePath string // e.g. "uploads"
    baseURL  string // e.g. "/uploads" or "https://..."
}

func NewLocalStorage(basePath, baseURL string) *LocalStorage {
    return &LocalStorage{basePath: basePath, baseURL: baseURL}
}

func (s *LocalStorage) SaveAvatar(src io.Reader, ext string) (string, error) {
    dir := filepath.Join(s.basePath, "avatars")
    ...
    return fmt.Sprintf("%s/avatars/%s", s.baseURL, filename), nil
}
// (similar for SaveLabResult and SaveAttachment)
```

**Rationale:** No external file referenced `BasePath` or `BaseURL` directly. All access was within `local_storage.go` itself. Making them unexported prevents post-construction mutation and follows DDD infrastructure encapsulation rules.

## Deferred Items

- **[MEDIUM]** `NewLocalStorage` input validation: Adding a `(*LocalStorage, error)` return would break the caller in `bootstrap/wire.go`. Deferred — requires coordinated update across the bootstrap layer in a dedicated pass.

## Final Build Status
PASS — `go build ./...` after all fixes
PASS — `go vet ./internal/...` after all fixes

## Remaining Violations
None at CRITICAL or HIGH severity. One MEDIUM deferred (constructor validation).
