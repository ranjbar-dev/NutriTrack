# DDD Audit: internal/domain/dietplan
Layer: domain
Audited: 2026-04-22
Files reviewed: 2

## Summary
- CRITICAL: 0
- HIGH: 3
- MEDIUM: 2
- LOW: 1
- PASS: 1 (repository/diet_plan_repository.go)

---

## Findings

### [HIGH-1] `json:` struct tags on domain-layer types

**File:** `entity/diet_plan.go:10`, `:18`, `:24`, `:118`
**Issue:** `NutritionalSummary`, `NutritionalRange`, `FoodSnapshot`, and `MedicationSnapshot` all carry `json:` struct tags inside the domain layer. Serialization is an interfaces concern.
**DDD Rule:** Aggregates MUST NOT have `json:` / `bson:` / `db:` struct tags.
**Fix:** Remove all `json:` tags from these types. Create mirrored DTO structs in `internal/interfaces/http/`.

---

### [HIGH-2] All aggregate and entity structs expose fully exported fields

**File:** `entity/diet_plan.go:44` (DietPlan), `:63` (DietPlanDay), `:73` (DietMeal), `:84` (MealOption), `:95` (MealOptionItem), `:107` (ExerciseRecommendation), `:124` (PrescribedMedication)
**Issue:** Every aggregate and entity has fully exported fields. Direct field mutation bypasses domain rules entirely.
**DDD Rule:** Aggregates MUST have unexported fields and expose state only through getter methods.
**Fix:** Lowercase all struct fields and add typed getters and validated setters.

---

### [HIGH-3] No factory `New*()` functions on any aggregate or entity

**File:** `entity/diet_plan.go` (entire file — absence of factory functions)
**Issue:** Not a single entity in this package has a `New*()` factory function. Callers use struct literals with no validation.
**DDD Rule:** Factory functions `New*()` MUST validate required fields and return `(*T, error)`.
**Fix:** Add factory functions for `DietPlan`, `DietPlanDay`, `DietMeal`, `MealOption`, `MealOptionItem`, `ExerciseRecommendation`, `PrescribedMedication`.

---

### [MEDIUM-1] No domain error variables defined

**File:** `entity/diet_plan.go` (absence)
**Issue:** The domain package defines no sentinel error variables.
**Fix:** Create `entity/errors.go` with `ErrInvalidClientID`, `ErrEmptyPlanTitle`, `ErrInvalidPlanDateRange`, `ErrPlanNotDraft`, `ErrPlanIncomplete`, etc.

---

### [MEDIUM-2] `NutritionalSummary` and `NutritionalRange` are Value Objects defined in entity/ with exported fields

**File:** `entity/diet_plan.go:10`, `:18`
**Issue:** These types have no UUID, represent immutable results, but are in `entity/` with exported fields. The `valueobject/` package exists but is unused.
**DDD Rule:** Value Objects MUST have all fields unexported.
**Fix:** Move to `internal/domain/dietplan/valueobject/nutritional.go`.

---

### [LOW-1] `service/` package is empty

**File:** `service/` directory (contains only `.gitkeep`)
**Issue:** Several operations warrant a domain service (nutritional calculation, plan completeness validation).

---

## Compliant Patterns Found

- **`repository/diet_plan_repository.go`** — Clean Go `interface` definition with no infrastructure details. ✓
- **`PlanStatus` type alias + const block** — Named type for status values. ✓
- **`IsActive()` method on `DietPlan`** — Domain behavior expressed as a method. ✓
- **No forbidden imports** — Neither file imports `internal/infrastructure`, `internal/interfaces`, or `internal/application`. ✓

## Fix Priority Order
1. **[HIGH-3]** Add `New*()` factory functions + domain errors
2. **[HIGH-2]** Make all fields unexported; add getters and validated domain methods
3. **[MEDIUM-1]** Add `entity/errors.go`
4. **[HIGH-1]** Remove `json:` struct tags; create DTOs in interfaces layer
5. **[MEDIUM-2]** Move `NutritionalSummary`/`NutritionalRange` to `valueobject/`
