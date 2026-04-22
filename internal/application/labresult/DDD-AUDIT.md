# DDD Audit: internal/application/labresult
Layer: application
Audited: 2026-04-22
Files reviewed: 1 (lab_result_service.go)

## Summary
- CRITICAL: 1
- HIGH: 1
- MEDIUM: 2
- LOW: 1
- PASS: 0

---

## Findings

### [CRITICAL] Application layer imports concrete infrastructure package

**File:** `lab_result_service.go:15`
**Issue:** The application service imports `"github.com/ranjbar-dev/nutritrack/internal/infrastructure/storage"` — a concrete infrastructure package. The application layer must never reference infrastructure. Makes the service impossible to test without real storage.
**DDD Rule:** Application layer MUST NOT import `internal/infrastructure/*`. May only depend on domain interfaces.
**Fix:** Define a storage port interface in the domain or application layer. Accept it in `NewLabResultService`. Wire `*storage.LocalStorage` in `bootstrap/wire.go`.

---

### [HIGH] Service struct holds concrete infrastructure type instead of interface

**File:** `lab_result_service.go:22, 28`
**Issue:** `LabResultService.storage` is typed as `*storage.LocalStorage`. Tightly coupled to file-system implementation; prevents unit testing with a mock.
**DDD Rule:** Application services MUST accept port interfaces, not concrete implementations.
**Fix:** Replace field and constructor parameter with `LabResultStorage` interface:
```go
type LabResultStorage interface {
    SaveLabResult(src io.Reader, ext string) (string, error)
}
```

---

### [MEDIUM] Domain entity constructed via bare struct literal — no factory function

**File:** `lab_result_service.go:142–151`
**Issue:** `entity.LabResult` is assembled using a raw struct literal. No `NewLabResult()` factory exists. Invariant enforcement (e.g., "must have file or link", "ResultType must be valid") is scattered across the application layer.
**DDD Rule:** Aggregates MUST have a factory `New*()` that validates required fields and returns `(T, error)`.
**Fix:** Add `NewLabResult(clientID, nutritionistID uuid.UUID, title, resultType string, testDate *time.Time, notes string, link *string) (*LabResult, error)` to the domain entity package. Replace struct literal with factory call.

---

### [MEDIUM] Role-based access control uses raw string comparisons

**File:** `lab_result_service.go:117–133, 230–247`
**Issue:** `callerRole` is a plain `string`; role checks use hardcoded literals (`"superadmin"`, `"nutritionist"`, `"client"`). Typos produce silent failures; no compile-time safety.
**DDD Rule:** Domain concepts should be expressed as domain types or constants.
**Fix:** Define `type Role string` constants in `internal/domain/shared`. Update service signatures to accept `shared.Role`.

---

### [LOW] `UploadLabResult` backward-compat wrapper introduces dual API surface

**File:** `lab_result_service.go:170–189`
**Issue:** Documented as "kept for backward compatibility"; delegates to `SubmitLabResult` with hardcoded defaults. Two public methods for the same operation.
**Fix:** Audit callers; if none remain after migrating to `SubmitLabResult`, delete the wrapper. Otherwise add `// Deprecated:` godoc comment.

---

## Compliant Patterns Found

- **Repository interfaces injected correctly** — accepts `labRepo.LabResultRepository` and `userRepo.UserRepository`. ✓
- **`NewLabResultService` factory present**. ✓
- **No SQL, DB tags, or HTTP types**. ✓
- **Domain errors used correctly** — all use `shared.Err*` sentinel values. ✓
- **File validation in correct layer** — MIME detection and size limits in application service. ✓
- **`checkAccess` helper cleanly isolates access-control logic**. ✓

## Fix Priority Order
1. **[CRITICAL]** Define `LabResultStorage` interface; remove infrastructure import
2. **[HIGH]** Swap `*storage.LocalStorage` field/parameter for interface; wire in bootstrap
3. **[MEDIUM]** Add `entity.NewLabResult()` factory; replace struct literal in `SubmitLabResult`
4. **[MEDIUM]** Define `shared.Role` type + constants; update signatures
5. **[LOW]** Deprecate or remove `UploadLabResult` wrapper
