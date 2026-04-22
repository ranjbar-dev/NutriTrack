# DDD Audit: internal/domain/labresult
Layer: domain
Audited: 2026-04-22
Files reviewed: 2

## Summary
- CRITICAL: 0
- HIGH: 2
- MEDIUM: 0
- LOW: 0
- PASS: 1 (repository/lab_result_repository.go)

---

## Findings

### [HIGH] LabResult aggregate exposes all fields as exported

**File:** `entity/lab_result.go:10`
**Issue:** Every field on `LabResult` is exported (uppercase). DDD requires aggregates to have unexported fields with controlled access through getter methods.
**DDD Rule:** Aggregates: unexported fields; getter/setter methods for controlled access.
**Fix:** Make all fields unexported and add public accessors:
- `ID()`, `ClientID()`, `NutritionistID()`, `Title()`, `ResultType()`, `TestDate()`, `FilePath()`, `OriginalName()`, `FileType()`, `FileSize()`, `Link()`, `Notes()`, `CreatedAt()`
- Setter only for mutable fields: `SetNotes(notes string)`

---

### [HIGH] Missing `NewLabResult()` factory function

**File:** `entity/lab_result.go` (entire file — no factory exists)
**Issue:** There is no `NewLabResult()` factory function. DDD mandates a factory that validates all required inputs and returns `(*LabResult, error)`.
**DDD Rule:** Factory `New*()` function MUST validate required fields and return `(T, error)` or `(*T, error)`.
**Fix:** Add factory that validates `clientID`, `title`, `resultType`, and that either `filePath` or `link` is provided. Add domain errors: `ErrLabResultMissingClientID`, `ErrLabResultMissingTitle`, `ErrLabResultMissingResultType`, `ErrLabResultNoSource`.

---

## Compliant Patterns Found

- **Repository as interface** (`repository/lab_result_repository.go`): `LabResultRepository` is correctly a pure Go `interface`. ✓
- **No struct tags on the aggregate** — zero `json:`, `bson:`, or `db:` struct tags. ✓
- **No forbidden cross-layer imports** — Only `time`, `github.com/google/uuid`, and sibling domain entity imports. ✓
- **Repository operations are domain-appropriate** — `Create`, `FindByID`, `ListByClientID`. ✓

## Fix Priority Order
1. **Add `NewLabResult()` factory with validation** — prevents invalid aggregates from being constructed.
2. **Make all `LabResult` fields unexported and add getters** — enforces encapsulation.
