---
phase: "05"
plan: "02"
subsystem: dietplan
tags: [nutrition, items, totals, endpoints]
dependency_graph:
  requires: [05-01]
  provides: [meal-option-items-api, nutritional-totals]
  affects: [diet_plan_service, diet_plan_handler, diet_plan_repository]
tech_stack:
  added: []
  patterns: [join-query, bubble-up-aggregation, ownership-walk-up]
key_files:
  created:
    - internal/application/dietplan/nutrition.go
  modified:
    - internal/domain/dietplan/entity/diet_plan.go
    - internal/domain/dietplan/repository/diet_plan_repository.go
    - internal/infrastructure/persistence/dietplan/mapper.go
    - internal/infrastructure/persistence/dietplan/pg_diet_plan_repository.go
    - internal/infrastructure/persistence/sqlc/meal_option_items.sql.go
    - internal/application/dietplan/diet_plan_service.go
    - internal/interfaces/http/handler/diet_plan_handler.go
    - internal/interfaces/http/router/router.go
    - db/queries/meal_option_items.sql
decisions:
  - Hand-wrote ListMealOptionItemsWithFood sqlc function instead of running sqlc generate (Windows mmap lock avoidance)
  - Nutritional totals bubble-up: item→option (sum), option→meal (min/max range), meal→day (sum of ranges)
metrics:
  duration: "~15 minutes"
  completed: "2026-04-21"
  tasks_completed: 12
  files_changed: 10
---

# Phase 05 Plan 02: Meal Option Items + Nutritional Totals Summary

**One-liner:** AddItem/RemoveItem endpoints with food-JOIN query and Calorie/Protein/Carbs/Fat/Fiber bubble-up from item → option → meal → day.

## What Was Built

### Domain Layer
- Added `NutritionalSummary`, `NutritionalRange`, `FoodSnapshot` types to `entity/diet_plan.go`
- Extended `MealOptionItem` with `Food *FoodSnapshot` and `Computed *NutritionalSummary`
- Extended `MealOption` with `Totals *NutritionalSummary`
- Extended `DietMeal` and `DietPlanDay` with `TotalRange *NutritionalRange`

### Infrastructure Layer
- Added `ListMealOptionItemsWithFood` JOIN query to `db/queries/meal_option_items.sql`
- Hand-wrote `ListMealOptionItemsWithFoodRow` struct and `ListMealOptionItemsWithFood` function in the sqlc package
- Added `mealOptionItemWithFoodToDomain` mapper converting join rows to enriched domain entities
- Added `ListItemsWithFood` method to `PgDietPlanRepository`
- Added `ListItemsWithFood` to `DietPlanRepository` interface

### Application Layer
- Created `nutrition.go` with `computeOptionTotals`, `computeMealRange`, `computeDayRange` helpers
- Updated `GetFullPlan` to load items with food data and compute nutritional totals at all levels
- Added `AddItem` service method (ownership walk-up: option→meal→day→plan)
- Added `RemoveItem` service method (ownership walk-up: item→option→meal→day→plan)
- Added `AddItemRequest` type

### HTTP Layer
- Added `AddItem` handler: `POST /plans/:id/days/:day_id/meals/:meal_id/options/:option_id/items`
- Added `RemoveItem` handler: `DELETE /plans/:id/days/:day_id/meals/:meal_id/options/:option_id/items/:item_id`
- Updated `planFullToMap` to include `totals` on options, `total_range` on meals and days
- Added `itemsToSlice` helper for enriched item serialization (includes food snapshot + computed values)

## Decisions Made

1. **Hand-wrote sqlc JOIN function** — `sqlc generate` skipped (Windows mmap lock known issue); hand-wrote `ListMealOptionItemsWithFoodRow` and `ListMealOptionItemsWithFood` directly in the `.sql.go` file.

2. **Nutritional aggregation strategy** — Items sum directly into option totals; options produce min/max range per meal (since a user picks ONE option); day totals are sum of meal ranges (min=sum of all meal mins, max=sum of all meal maxes).

## Deviations from Plan

None — plan executed exactly as written.

## Known Stubs

None — all data is wired from DB via JOIN query.

## Threat Flags

None — new endpoints follow the same ownership walk-up auth pattern as existing plan endpoints.

## Self-Check: PASSED

- `internal/application/dietplan/nutrition.go` — FOUND
- `internal/domain/dietplan/entity/diet_plan.go` updated — FOUND
- `internal/infrastructure/persistence/sqlc/meal_option_items.sql.go` updated — FOUND
- Commit `0051859` — FOUND (`git log --oneline -1`)
- `go build ./...` — PASSED (exit code 0)
