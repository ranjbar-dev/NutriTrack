# DDD Audit: internal/domain/medication
Layer: domain
Audited: 2026-04-22
Files reviewed: 2

## Summary
- CRITICAL: 0
- HIGH: 1
- MEDIUM: 1
- LOW: 1
- PASS: 1 (repository/medication_repository.go)

---

## Findings

### [HIGH] Aggregate exposes all fields as exported

**File:** `entity/medication.go:10`
**Issue:** All nine fields (`ID`, `Name`, `NameNormalized`, `Description`, `Unit`, `CreatedBy`, `IsActive`, `CreatedAt`, `UpdatedAt`) are exported. Any caller can read or mutate state directly, bypassing domain invariants.
**DDD Rule:** Aggregates MUST have unexported fields and expose state only through getter methods.
**Fix:** Make all fields unexported and add getter methods. Add a `Deactivate()` domain command.

---

### [MEDIUM] Missing `NewMedication()` factory function

**File:** `entity/medication.go` (entire file)
**Issue:** There is no `NewMedication(...)` factory function. Callers construct `Medication` via struct literals with no validation.
**DDD Rule:** Aggregates MUST expose a `New*()` factory function that validates required fields and returns `(*T, error)`.
**Fix:** Add factory with `name`, `description`, `unit`, `createdBy` params and validation.

---

### [LOW] No domain error variables defined

**File:** `entity/medication.go` (or a sibling `errors.go`)
**Issue:** No sentinel errors defined. Application and infrastructure layers cannot do typed error checking.
**Fix:** Add `ErrMedicationNotFound`, `ErrMedicationNameRequired`, `ErrMedicationUnitRequired`, `ErrMedicationInactive`.

---

## Compliant Patterns Found
- **`repository/medication_repository.go`** — `MedicationRepository` is correctly defined as a Go `interface`. ✓
- **No struct tags** — `entity/medication.go` has zero `json:`, `bson:`, or `db:` struct tags. ✓
- **No forbidden cross-layer imports** — No imports of `internal/infrastructure`, `internal/interfaces`, or `internal/application`. ✓
- **`valueobject/`** — Empty package, no violations. ✓

## Fix Priority Order
1. **[HIGH]** Make all `Medication` fields unexported and add getter methods
2. **[MEDIUM]** Add `NewMedication()` factory with field validation
3. **[LOW]** Add `errors.go` with sentinel domain errors
