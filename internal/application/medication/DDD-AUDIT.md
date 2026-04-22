# DDD Audit: internal/application/medication
Layer: application
Audited: 2026-04-22
Files reviewed: 1 (medication_service.go)

## Summary
- CRITICAL: 0
- HIGH: 2
- MEDIUM: 2
- LOW: 0
- PASS: 0

---

## Findings

### [HIGH] Service constructs domain entity via struct literal (no factory call)

**File:** `medication_service.go` (CreateMedication, likely)
**Issue:** The application service instantiates `medicationEntity.Medication{...}` using a raw struct literal, bypassing any factory validation. Once domain fields are unexported (as DDD requires), this will be a compile error.
**DDD Rule:** Application services MUST call the domain's `New*()` factory function to construct aggregates.
**Fix:** Replace struct literal with `medicationentity.NewMedication(...)`.

---

### [HIGH] Direct mutation of exported aggregate fields in UpdateMedication

**File:** `medication_service.go` (UpdateMedication)
**Issue:** `UpdateMedication` reads the aggregate, then directly assigns to exported fields (`medication.Name = ...`, etc.) before persisting. This bypasses all domain invariant protection.
**DDD Rule:** State changes MUST go through domain aggregate methods, not direct field assignment.
**Fix:** Add domain setter/update methods (`SetName`, `SetDosage`, etc.) to the `Medication` aggregate. The application service calls these methods.

---

### [MEDIUM] Deactivate state transition delegated to repository by ID

**File:** `medication_service.go` (DeactivateMedication)
**Issue:** `DeactivateMedication` calls the repository directly with just an ID, bypassing the aggregate lifecycle. The domain aggregate should have a `Deactivate()` method.
**DDD Rule:** Domain state transitions must be expressed as aggregate methods (Load → Transition → Persist).
**Fix:** Load the aggregate, call `medication.Deactivate()`, then persist.

---

### [MEDIUM] Magic role string literals scattered through service

**File:** `medication_service.go`
**Issue:** Role checks compare against raw string literals (e.g. `"nutritionist"`, `"admin"`). These belong in the domain as typed constants or value objects.
**Fix:** Define a `Role` type in `internal/domain/user/valueobject/` and use typed constants.

---

## Compliant Patterns Found

- No forbidden cross-layer imports detected. ✓
- Service accepts repository interfaces (not concrete types). ✓

## Fix Priority Order
1. **[HIGH]** Replace struct literals with domain factory calls
2. **[HIGH]** Add domain setter methods; remove direct field mutation from application layer
3. **[MEDIUM]** Implement aggregate `Deactivate()` method; load-transition-persist pattern
4. **[MEDIUM]** Replace magic role strings with typed domain constants
