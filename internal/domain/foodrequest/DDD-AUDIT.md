# DDD Audit: internal/domain/foodrequest
Layer: domain
Audited: 2026-04-22
Files reviewed: 2

## Summary
- CRITICAL: 1
- HIGH: 1
- MEDIUM: 1
- LOW: 1
- PASS: 1 (repository/food_request_repository.go — no violations)

---

## Findings

### [CRITICAL] FoodRequest aggregate exposes all fields as exported

**File:** `entity/food_request.go:17-26`
**Issue:** Every field on the `FoodRequest` struct is exported (`ID`, `ClientID`, `NutritionistID`, `FoodName`, `Status`, `RejectionReason`, `CreatedFoodID`, `CreatedAt`, `UpdatedAt`). Any consumer in any layer can read or mutate aggregate state directly, bypassing all domain invariants.
**DDD Rule:** Aggregates MUST have unexported fields. State changes MUST be mediated through exported methods that enforce business rules.
**Fix:** Make all fields lowercase and expose only what is needed via getter methods and domain-behaviour methods (`Approve()`, `Reject()`).

---

### [HIGH] Missing `NewFoodRequest()` factory function

**File:** `entity/food_request.go` (no factory present)
**Issue:** The `FoodRequest` aggregate has no constructor. Callers build the struct with raw literals, bypassing required-field validation.
**DDD Rule:** Aggregates MUST have a `New*()` factory function that validates required inputs and returns `(*T, error)`.
**Fix:** Add `NewFoodRequest(clientID, nutritionistID uuid.UUID, foodName string) (*FoodRequest, error)` with validation.

---

### [MEDIUM] No domain-behaviour mutation methods; state changes left to external callers

**File:** `entity/food_request.go:29-31`
**Issue:** Only read-only predicate methods exist (`IsPending()`, `IsApproved()`, `IsRejected()`). No `Approve()`, `Reject()` methods to enforce state transitions.
**DDD Rule:** Domain state transitions must be expressed as explicit aggregate methods.
**Fix:** Add `Approve(createdFoodID uuid.UUID) error` and `Reject(reason *string) error` with precondition checks.

---

### [LOW] No domain error variables defined

**File:** `entity/food_request.go` (none present)
**Issue:** No sentinel error values (`var Err* = errors.New(...)`).
**Fix:** Add `ErrFoodRequestEmptyName`, `ErrFoodRequestAlreadyReviewed`, `ErrFoodRequestNotFound`.

---

## Compliant Patterns Found

- **`repository/food_request_repository.go`** — `FoodRequestRepository` is correctly declared as a Go `interface`. ✓
- **No forbidden imports**. ✓
- **No struct tags** — `FoodRequest` has zero `json:`, `bson:`, or `db:` struct tags. ✓
- **`FoodRequestStatus` typed constant** — Named type with package-level constants. ✓

## Fix Priority Order
1. **[CRITICAL]** Make all `FoodRequest` fields unexported and add getter methods
2. **[HIGH]** Add `NewFoodRequest()` factory with validation
3. **[LOW]** Add `errors.go` with sentinel domain errors
4. **[MEDIUM]** Add `Approve()` and `Reject()` transition methods
