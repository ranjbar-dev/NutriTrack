# DDD Fix Report: internal/domain/dietplan
Layer: domain
Fixed: 2025-07
Based on: DDD-AUDIT.md

## Baseline Build Status
PASS — `go build ./...` before fixes

## Fix Plan

| # | Finding | Severity | Files | Strategy | Status |
|---|---------|----------|-------|----------|--------|
| 1 | 7 aggregate/entity structs have exported fields | HIGH | entity/diet_plan.go | SAFE | FIXED |
| 2 | No `New*()` factory functions | HIGH | entity/diet_plan.go | SAFE | FIXED |
| 3 | Value objects have json struct tags | HIGH | entity/diet_plan.go | SAFE | FIXED |
| 4 | Callers use direct field access | HIGH | persistence/dietplan/mapper.go, pg_diet_plan_repository.go, cached_diet_plan_repository.go, application/dietplan/diet_plan_service.go, application/dietplan/nutrition.go, interfaces/http/handler/diet_plan_handler.go | SAFE | FIXED |

## Changes Applied

### Fix 1 + 2 + 3: Entity rewrite
**File:** `internal/domain/dietplan/entity/diet_plan.go`
**Change:**
- Value objects `NutritionalSummary`, `NutritionalRange`, `FoodSnapshot`, `MedicationSnapshot`: removed all `json:` struct tags; fields remain exported for cross-package struct literal construction in the infrastructure mapper.
- `DietPlan` (7 fields), `DietPlanDay`, `DietMeal`, `MealOption`, `MealOptionItem`, `ExerciseRecommendation`, `PrescribedMedication`: all fields made unexported. Added getters for every field. Added setters for mutable fields (SetID, SetStatus, SetCreatedAt, SetUpdatedAt, SetTitle, SetNotes, SetDailyWaterTargetML, SetDays, SetMeals, SetTotalRange, SetOptions, SetTotals, SetItems, SetExercises, SetPrescriptions, SetFood, SetComputed, SetMedication).
- `NewDietPlan(clientID, nutritionistID uuid.UUID, title string, startDate, endDate time.Time, notes string, dailyWaterTargetML int) *DietPlan`
- `NewDietPlanDay(planID uuid.UUID, dayNumber int) *DietPlanDay`
- `NewDietMeal(dayID uuid.UUID, title, scheduledTime string, displayOrder int) *DietMeal`
- `NewMealOption(mealID uuid.UUID, optionNumber int) *MealOption`
- `NewMealOptionItem(optionID, foodID uuid.UUID, quantity float64, unit, notes string) *MealOptionItem`
- `NewExerciseRecommendation(dayID uuid.UUID, exerciseName string, durationMinutes int, description string, caloriesBurnEstimate int) *ExerciseRecommendation`
- `NewPrescribedMedication(dayID, medicationID uuid.UUID, dosage, frequency string, times []string, instructions string, startDate, endDate time.Time) *PrescribedMedication`
- `Reconstitute*()` functions added for all 7 types.
**Build:** PASS

### Fix 4a: Mapper updated
**File:** `internal/infrastructure/persistence/dietplan/mapper.go`
**Change:** All 7 toDomain functions use `entity.Reconstitute*()`. FoodSnapshot/MedicationSnapshot struct literals preserved (valid since value object fields remain exported).
**Build:** PASS

### Fix 4b: Repository updated
**File:** `internal/infrastructure/persistence/dietplan/pg_diet_plan_repository.go`
**Change:** All 8 write methods (CreateWithArchive, Update, AddDay, AddMeal, AddOption, AddItem, AddExercise, AddPrescription) + FindPrescriptionByID use getters for DB params and setters for DB-generated fields.
**Build:** PASS

### Fix 4c: Cached repository updated
**File:** `internal/infrastructure/persistence/dietplan/cached_diet_plan_repository.go`
**Change:** `plan.ClientID` → `plan.ClientID()` in CreateWithArchive, Update, and Delete (3 occurrences).
**Build:** PASS

### Fix 4d: Service updated
**File:** `internal/application/dietplan/diet_plan_service.go`
**Change:** All 7 struct literals replaced with `entity.New*()` factory calls. All direct field reads replaced with getter calls. All direct field writes replaced with setter calls (SetTitle, SetNotes, SetDailyWaterTargetML, SetItems, SetTotals, SetOptions, SetTotalRange, SetMeals, SetExercises, SetPrescriptions, SetDays).
**Build:** PASS

### Fix 4e: Nutrition helpers updated
**File:** `internal/application/dietplan/nutrition.go`
**Change:** `it.Computed` → `it.Computed()`, `opt.Totals` → `opt.Totals()`, `m.TotalRange` → `m.TotalRange()`. NutritionalSummary/NutritionalRange field accesses (`.Calories`, `.Min.Calories`, etc.) remain unchanged since value object fields are exported.
**Build:** PASS

### Fix 4f: Handler updated
**File:** `internal/interfaces/http/handler/diet_plan_handler.go`
**Change:** Added `summaryToMap(*entity.NutritionalSummary) map[string]any` and `rangeToMap(*entity.NutritionalRange) map[string]any` nil-safe helpers. Updated `planToMap`, `planFullToMap`, `itemsToSlice`, `exercisesToSlice`, `prescriptionsToSlice` to use getter methods throughout. FoodSnapshot/MedicationSnapshot fields accessed directly since they remain exported.
**Build:** PASS

## Final Build Status
PASS — `go build ./...` after all fixes
PASS — `go vet ./internal/...` after all fixes

## Remaining Violations
None — all CRITICAL and HIGH findings resolved.

## Notes
- Value object fields (NutritionalSummary, NutritionalRange, FoodSnapshot, MedicationSnapshot) were left exported as a MEDIUM-level pragmatic decision: making them unexported would require adding getters to simple data-container value objects with no behavior, adding boilerplate without meaningful encapsulation benefit. The json tags (HIGH violation) were removed.
