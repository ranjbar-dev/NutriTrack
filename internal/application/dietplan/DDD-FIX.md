# DDD Fix Report: internal/application/dietplan
Layer: application
Fixed: 2026-04-22
Based on: DDD-AUDIT.md

## Baseline Build Status
PASS — `go build ./...` before fixes

## Fix Plan

| # | Finding | Severity | Files | Strategy | Status |
|---|---------|----------|-------|----------|--------|
| 1 | Nutritional computation logic in app layer | HIGH | nutrition.go, diet_plan.go (domain) | SAFE | FIXED |
| 2 | Domain entities constructed via struct literals | HIGH | diet_plan_service.go | N/A | PRE-FIXED (factories already in use) |
| 3 | Direct mutation of aggregate exported fields | HIGH | diet_plan_service.go | N/A | PRE-FIXED (setters already in use) |
| 4 | Magic string role constants | MEDIUM | diet_plan_service.go | DEFERRED | DEFERRED |
| 5 | NewDietPlanService not using functional options | LOW | diet_plan_service.go | DEFERRED | DEFERRED |

## Changes Applied

### Fix 1: Move nutritional computation to domain entity methods

**Files changed:**
- `internal/domain/dietplan/entity/diet_plan.go` — added three domain methods
- `internal/application/dietplan/diet_plan_service.go` — updated to call domain methods
- `internal/application/dietplan/nutrition.go` — replaced with empty package stub

**Domain methods added to `diet_plan.go`:**
- `(*MealOption).ComputeTotals() *NutritionalSummary` — sums item computed values
- `(*DietMeal).ComputeTotalRange() *NutritionalRange` — min/max across options
- `(*DietPlanDay).ComputeTotalRange() *NutritionalRange` — sums meal ranges

**Service changes in `GetFullPlan` tree assembly:**
```go
// Before (app-layer functions)
option.SetTotals(computeOptionTotals(items))
meal.SetTotalRange(computeMealRange(options))
day.SetTotalRange(computeDayRange(meals))

// After (domain methods)
option.SetTotals(option.ComputeTotals())
meal.SetTotalRange(meal.ComputeTotalRange())
day.SetTotalRange(day.ComputeTotalRange())
```

**Build:** PASS

## Pre-Fixed Findings (verified)

- **Factory functions**: All 7 entity constructors (`NewDietPlan`, `NewDietPlanDay`, `NewDietMeal`, `NewMealOption`, `NewMealOptionItem`, `NewExerciseRecommendation`, `NewPrescribedMedication`) already existed in domain and were already called from the service. ✓
- **Setter methods**: `SetDays`, `SetMeals`, `SetOptions`, `SetItems`, `SetTotals`, `SetTotalRange`, `SetTitle`, `SetNotes`, `SetDailyWaterTargetML` all existed and were already used. ✓

## Deferred Items

- **[MEDIUM]** Magic role string constants — extracting `"superadmin"` / `"nutritionist"` to `internal/domain/shared` requires touching 16+ call sites across the service and has no impact on compilation safety; deferred to a dedicated refactoring phase.
- **[LOW]** Functional options pattern for `NewDietPlanService` — optional, low-risk.

## Final Build Status
PASS — `go build ./...` after all fixes
PASS — `go vet ./internal/...` after all fixes

## Remaining Violations
None at HIGH severity. MEDIUM and LOW items deferred as documented.
