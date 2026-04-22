# DDD Audit: internal/application/dietplan
Layer: application
Audited: 2026-04-22
Files reviewed: 2 (diet_plan_service.go, nutrition.go)

## Summary
- CRITICAL: 0
- HIGH: 3
- MEDIUM: 1
- LOW: 1
- PASS: 0

---

## Findings

### [HIGH] Domain computation logic (nutritional aggregation) lives in the application layer

**File:** `nutrition.go:1–100`
**Issue:** `computeOptionTotals`, `computeMealRange`, and `computeDayRange` encode domain rules — how nutritional values are summed through the aggregate tree. These are domain invariants, not orchestration steps.
**DDD Rule:** Domain business logic MUST live in `internal/domain/`. The application layer coordinates; it must not implement domain calculations.
**Fix:** Move all three functions to the domain layer as methods on the respective aggregate/entity types: `MealOption.ComputeTotals()`, `DietMeal.ComputeTotalRange()`, `DietPlanDay.ComputeTotalRange()`. Service calls these methods.

---

### [HIGH] Domain entities constructed via struct literals — factory functions bypassed

**File:** `diet_plan_service.go:74, 169, 204, 249, 351, 459, 528`
**Issue:** Every creation path uses a raw struct literal (`&entity.DietPlan{...}`, `&entity.DietPlanDay{...}`, `&entity.DietMeal{...}`, etc.) rather than `New*()` factory functions. No domain-level validation is executed.
**DDD Rule:** Aggregates MUST have `New*()` factory functions. The application layer MUST call the factory, not the struct literal.
**Fix:** Add factory functions: `NewDietPlan`, `NewDietPlanDay`, `NewDietMeal`, `NewMealOption`, `NewMealOptionItem`, `NewExerciseRecommendation`, `NewPrescribedMedication`. Replace all struct literals with factory calls.

---

### [HIGH] Application layer directly mutates domain aggregate exported fields

**File:** `diet_plan_service.go:128–149` (GetFullPlan tree assembly) and `:580–581` (UpdatePlan)
**Issue:** `GetFullPlan` writes directly to exported fields to assemble the tree (`option.Items`, `meal.Options`, `day.Meals`, `plan.Days`, etc.). `UpdatePlan` directly writes `plan.Title`, `plan.Notes`, `plan.DailyWaterTargetML`.
**DDD Rule:** Aggregates MUST NOT expose raw writable fields — state changes go through setter/assembler methods.
**Fix:** Add setter methods to domain entities (`SetItems`, `SetOptions`, `SetMeals`, `SetDays`, `Update`); make fields unexported. Service calls these methods.

---

### [MEDIUM] Magic string role constants hardcoded at 16+ call sites

**File:** `diet_plan_service.go:103, 165, 200, 245, 264, 266, 293, 345, 403, 459, 503, 551, 605, 630, 659, 702`
**Issue:** `"superadmin"` and `"nutritionist"` appear as raw string literals at 16+ sites. A typo silently bypasses authorization.
**DDD Rule:** Domain concepts MUST be expressed as named constants or value objects in the domain layer.
**Fix:** Define `type Role string` constants in `internal/domain/shared`. Replace all inline strings.

---

### [LOW] `NewDietPlanService` does not use the functional options pattern

**File:** `diet_plan_service.go:57–59`
**Issue:** Adding optional dependencies requires a breaking signature change.
**Fix (optional):** Adopt functional options pattern for extensibility.

---

## Compliant Patterns Found

- **Repository interfaces accepted, not concrete implementations** — holds `dietRepo.DietPlanRepository` and `userRepo.UserRepository`. ✓
- **No forbidden imports** — No imports from `internal/infrastructure` or `internal/interfaces`. ✓
- **Domain errors used** — `shared.ErrForbidden`, `shared.ErrPlanNotFound`, etc. ✓
- **Factory function present** — `NewDietPlanService` exists. ✓
- **Authorization enforced in the application layer**. ✓
- **Service struct fields are unexported**. ✓

## Fix Priority Order
1. **[HIGH]** Add `New*()` factory functions to all domain entities; replace 7 struct literals in service
2. **[HIGH]** Move nutritional computation functions from `nutrition.go` to domain entity methods
3. **[HIGH]** Make aggregate fields unexported; add setter/assembler methods; replace 11 direct mutations
4. **[MEDIUM]** Extract role string constants to `internal/domain/shared`
5. **[LOW]** Adopt functional options in `NewDietPlanService`
