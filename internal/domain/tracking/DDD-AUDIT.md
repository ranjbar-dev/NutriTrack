# DDD Audit: internal/domain/tracking
Layer: domain
Audited: 2026-04-22
Files reviewed: 2

## Summary
- CRITICAL: 0
- HIGH: 2
- MEDIUM: 1
- LOW: 2
- PASS: 1 (`repository/tracking_repository.go`)

---

## Findings

### [HIGH] All entity struct fields are exported — no encapsulation

**File:** `entity/tracking.go:10–97`
**Issue:** All 6 entity types (`FoodLog`, `WaterLog`, `SleepLog`, `ExerciseLog`, `MedicationLog`, `BodyMeasurement`) declare every field as exported. Any caller can read and mutate entity state directly, bypassing all domain invariants.
**DDD Rule:** Entities and Aggregates MUST have unexported fields. Access is controlled through exported getter and setter methods.
**Fix:** Lowercase every field, add exported getters, and add domain-validating setters.

---

### [HIGH] No factory `New*()` functions — entities can be created in invalid state

**File:** `entity/tracking.go` (entire file)
**Issue:** None of the 6 entity types have a `New*()` constructor function. Callers must populate structs directly.
**DDD Rule:** Every Aggregate/Entity MUST expose a `New*()` factory that validates required fields and returns `(*T, error)`.
**Fix:** Add one `New*()` per entity: `NewFoodLog`, `NewWaterLog`, `NewSleepLog`, `NewExerciseLog`, `NewMedicationLog`, `NewBodyMeasurement`.

---

### [MEDIUM] No domain sentinel errors defined

**File:** `entity/tracking.go`, `repository/tracking_repository.go`
**Issue:** No `var Err* = errors.New(...)` declarations anywhere under `internal/domain/tracking/`.
**DDD Rule:** Domain errors MUST be defined in the domain package.
**Fix:** Add `entity/errors.go` with `ErrTrackingNotFound`, `ErrTrackingUnauthorized`, `ErrInvalidQuantity`, `ErrInvalidSleepQuality`, `ErrInvalidMeasuredDate`.

---

### [LOW] `valueobject/` sub-package is empty

**File:** `valueobject/` (only `.gitkeep`)
**Issue:** Natural candidates for Value Objects: `SleepQuality` enum, `Quantity+Unit` pair.

---

### [LOW] Entity-specific logic (`Duration` computation) not encapsulated

**Issue:** Sleep duration should be computed as a domain method, not by callers.

---

## Compliant Patterns Found

- **`repository/tracking_repository.go`** — Correctly defined as a Go `interface`. ✓
- **No `json:`, `bson:`, `db:` struct tags** on any entity. ✓
- **No forbidden cross-layer imports**. ✓

## Fix Priority Order
1. **[HIGH]** Make all entity fields unexported; add getter methods
2. **[HIGH]** Add `New*()` factory functions with validation
3. **[MEDIUM]** Add `entity/errors.go` with domain error variables
