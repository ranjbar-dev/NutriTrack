# DDD Audit: internal/infrastructure/storage
Layer: infrastructure
Audited: 2026-04-22
Files reviewed: 1 (local_storage.go)

## Summary
- CRITICAL: 0
- HIGH: 1
- MEDIUM: 1
- LOW: 1
- PASS: 0

---

## Findings

### [HIGH] No domain port interface for storage — application services import concrete type directly

**File:** `local_storage.go`
**Issue:** The `LocalStorage` concrete struct is imported directly by one or more application services (`internal/application/message/`, `internal/application/user/`). There is no domain port interface for storage, so application services are tightly coupled to the concrete infrastructure type.
**DDD Rule:** Application layer MUST only import domain port interfaces. Infrastructure implementations must be behind interfaces.
**Fix:** Define a `FileStorage` port interface in the domain layer (e.g. `internal/domain/shared/port/`):
```go
type FileStorage interface {
    Save(ctx context.Context, filename string, content io.Reader) (url string, err error)
    Delete(ctx context.Context, url string) error
}
```
`LocalStorage` implements `FileStorage`. Application services accept `FileStorage`.

---

### [MEDIUM] `NewLocalStorage` factory has no input validation

**File:** `local_storage.go`
**Issue:** `NewLocalStorage(basePath, baseURL string)` creates the adapter without checking that `basePath` is a valid, writable directory or that `baseURL` is non-empty.
**Fix:** Validate both inputs; return `(*LocalStorage, error)`.

---

### [LOW] `LocalStorage` has exported fields (`BasePath`, `BaseURL` should be unexported)

**File:** `local_storage.go`
**Issue:** Infrastructure implementation fields are exported, allowing external modification of configuration after construction.
**Fix:** Rename to `basePath` and `baseURL` (unexported).

---

## Compliant Patterns Found

- `LocalStorage` is scoped entirely within the infrastructure package. ✓
- No imports of domain entity packages or application packages. ✓

## Fix Priority Order
1. **[HIGH]** Define `FileStorage` port interface in domain layer; update all application services to use it
2. **[MEDIUM]** Add input validation to `NewLocalStorage`; return error
3. **[LOW]** Make `basePath`/`baseURL` fields unexported
