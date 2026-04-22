# DDD Audit: internal/application/food
Layer: application
Audited: 2026-04-22
Files reviewed: 2 (food_service.go, food_category_service.go)

## Summary
- CRITICAL: 2
- HIGH: 1
- MEDIUM: 0
- LOW: 1
- PASS: 0

---

## Findings

### [CRITICAL] Food aggregate constructed via raw struct literal — no factory called

**File:** `food_service.go:81`
**Issue:** `CreateFood` builds `entity.Food` through a direct struct literal, assigning all fields by hand without validation. No `entity.NewFood(...)` factory exists or is called.
**DDD Rule:** Aggregates MUST be created through a `New*()` factory function that validates required fields and returns `(*T, error)`.
**Fix:** Add `NewFood(name, unit string, calories, protein, carbohydrate, fat, fiber, sugar, sodium, amount float64, createdBy *uuid.UUID, categories []FoodCategory) (*Food, error)` factory to the domain. Replace struct literal in `CreateFood` with factory call.

---

### [CRITICAL] Application service directly mutates Food aggregate exported fields

**File:** `food_service.go:178`
**Issue:** `UpdateFood` writes directly to the fetched aggregate's exported fields (`food.Name`, `food.NameNormalized`, `food.Unit`, `food.Calories`, etc.). Any invariant check must be duplicated by every caller or skipped silently.
**DDD Rule:** Aggregates MUST encapsulate state behind unexported fields and expose behaviour through domain methods.
**Fix:** Make `Food` fields unexported; add `food.UpdateDetails(name, unit string, ...) error` domain method. Replace direct field assignments with method call.

---

### [HIGH] Authorization policy for `FoodCategoryService.Create` delegated to handler layer

**File:** `food_category_service.go:23`
**Issue:** Method comment acknowledges that "superadmin only" constraint is enforced in the handler, not the service. Any caller that bypasses the HTTP handler (test, CLI) can create categories without restriction.
**DDD Rule:** Application services are the gatekeepers of business policies. Authorization MUST be enforced inside the service.
**Fix:** Add `callerRole` parameter and enforce the constraint inside the service method. Update HTTP handler to pass `callerRole` from authenticated context.

---

### [LOW] Role identity expressed as untyped string literals

**File:** `food_service.go:64, 158, 163, 199, 204`
**Issue:** Role values `"nutritionist"` and `"superadmin"` are compared as raw strings. A single typo produces a silent authorization bypass with no compile-time signal.
**DDD Rule:** Domain concepts should be expressed as typed constants or value objects.
**Fix:** Define `type Role string` with `RoleNutritionist`, `RoleSuperAdmin` constants in `internal/domain/shared`. Replace all string literals with the constants.

---

## Compliant Patterns Found

- Both services depend on repository **interfaces** from the domain layer — no concrete infrastructure types. ✓
- `NewFoodService` and `NewFoodCategoryService` factory constructors accept domain repository interfaces. ✓
- No imports of `internal/infrastructure` or `internal/interfaces`. ✓
- Domain error variables (`shared.ErrForbidden`, `shared.ErrFoodNotFound`, etc.) used throughout. ✓
- Service struct fields are unexported. ✓

## Fix Priority Order
1. **[CRITICAL]** Add `entity.NewFood()` factory with validation; use it in `CreateFood`
2. **[CRITICAL]** Make `entity.Food` fields unexported; add `food.UpdateDetails()` domain method
3. **[HIGH]** Add `callerRole` parameter + guard to `FoodCategoryService.Create`
4. **[LOW]** Define `shared.Role` type + constants; replace all string literals
