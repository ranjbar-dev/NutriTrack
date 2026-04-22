# DDD Audit: internal/domain/food
Layer: domain
Audited: 2026-04-22
Files reviewed: 3

## Summary
- CRITICAL: 1
- HIGH: 4
- MEDIUM: 1
- LOW: 1
- PASS: 2 (repository/food_repository.go, repository/food_category_repository.go)

---

## Findings

### [CRITICAL] Food aggregate exposes all fields as exported

**File:** `entity/food.go:10-27`
**Issue:** Every field on the `Food` struct is exported (`ID`, `Name`, `NameNormalized`, `Unit`, `Calories`, `Protein`, `Carbohydrate`, `Fat`, `Fiber`, `Sugar`, `Sodium`, `Amount`, `CreatedBy`, `IsActive`, `Categories`, `CreatedAt`, `UpdatedAt`). External packages can read and write any field directly, bypassing all domain invariants.
**DDD Rule:** Aggregates MUST have unexported fields. State changes MUST go through exported methods.
**Fix:** Make all fields unexported and add getter/setter methods.

---

### [HIGH] No factory `NewFood()` function

**File:** `entity/food.go` — missing entirely
**Issue:** The `Food` aggregate has no factory constructor. Callers create `Food` via struct literal with no validation.
**DDD Rule:** Aggregates MUST have a `New*()` factory function that validates all required fields and returns `(*T, error)`.

---

### [HIGH] No domain errors defined in the package

**File:** `entity/food.go` — no `errors.go` or sentinel error variables exist
**Issue:** There are no sentinel domain errors. Application and interface layers have nothing domain-specific to match against.
**DDD Rule:** Domain errors (`var Err* = errors.New(...)`) MUST be defined in the domain package.
**Fix:** Create `entity/errors.go` with `ErrFoodNotFound`, `ErrFoodNameRequired`, `ErrFoodUnitRequired`, `ErrFoodInvalidNutrition`, `ErrCategoryNameRequired`, `ErrCategoryNotFound`.

---

### [HIGH] `FoodCategory` entity has all exported fields and no factory function

**File:** `entity/food.go:30-34`
**Issue:** `FoodCategory` has a UUID identity but all fields are exported and there is no `NewFoodCategory()` factory.
**DDD Rule:** Entities must have unexported fields and a validated factory function.

---

### [HIGH] `FoodCategory` defined in same file as `Food` aggregate

**File:** `entity/food.go:29-34`
**Issue:** `FoodCategory` is a distinct entity defined in `food.go` alongside the main aggregate.
**DDD Rule:** Each entity/aggregate should live in its own file.
**Fix:** Move `FoodCategory` to `entity/food_category.go`.

---

### [MEDIUM] `valueobject/` package is empty — nutritional data modelled as raw primitives

**File:** `entity/food.go:15-22`
**Issue:** Eight nutritional fields are plain `float64` primitives. They form a cohesive concept with shared validation rules.
**DDD Rule:** Groups of related primitives with shared validation should be modelled as Value Objects.
**Fix:** Create `valueobject/nutrition.go` with an immutable `Nutrition` value object.

---

### [LOW] `NameNormalized` is an exported field derived from `Name`

**File:** `entity/food.go:12`
**Issue:** `NameNormalized` should be computed internally when `Name` is set, not exposed as an independent exported field.

---

## Compliant Patterns Found

- **Repository interfaces correctly placed in the domain layer** — Both are pure Go `interface` types. ✓
- **No forbidden cross-layer imports** — No `internal/infrastructure`, `internal/interfaces`, or `internal/application` imports. ✓
- **No struct tags on domain types** — `Food` and `FoodCategory` carry zero `json:`, `bson:`, or `db:` struct tags. ✓

## Fix Priority Order
1. **[CRITICAL]** Unexport all `Food` struct fields and add getter/setter methods
2. **[HIGH]** Add `NewFood()` factory with input validation
3. **[HIGH]** Create `entity/errors.go` with domain sentinel error variables
4. **[HIGH]** Unexport `FoodCategory` fields and add `NewFoodCategory()` factory
5. **[HIGH]** Move `FoodCategory` to `entity/food_category.go`
6. **[MEDIUM]** Create `valueobject/nutrition.go`
