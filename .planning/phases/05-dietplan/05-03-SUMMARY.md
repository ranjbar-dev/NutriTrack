---
phase: "05"
plan: "03"
subsystem: dietplan
tags: [exercise-recommendations, prescribed-medications, diet-plan-day, go, ddd, sqlc]
dependency_graph:
  requires: [05-01, 05-02]
  provides: [exercise-recommendations-api, prescribed-medications-api]
  affects: [diet-plan-day, get-full-plan]
tech_stack:
  added: []
  patterns: [DDD repository pattern, sqlc hand-written queries, pgx/v5 text[] scanning]
key_files:
  created:
    - migrations/000006_exercises_prescriptions.up.sql
    - migrations/000006_exercises_prescriptions.down.sql
    - db/queries/exercise_recommendations.sql
    - db/queries/day_prescribed_medications.sql
    - internal/infrastructure/persistence/sqlc/exercise_recommendations.sql.go
    - internal/infrastructure/persistence/sqlc/day_prescribed_medications.sql.go
  modified:
    - internal/infrastructure/persistence/sqlc/models.go
    - internal/domain/dietplan/entity/diet_plan.go
    - internal/domain/dietplan/repository/diet_plan_repository.go
    - internal/infrastructure/persistence/dietplan/mapper.go
    - internal/infrastructure/persistence/dietplan/pg_diet_plan_repository.go
    - internal/application/dietplan/diet_plan_service.go
    - internal/interfaces/http/handler/diet_plan_handler.go
    - internal/interfaces/http/router/router.go
decisions:
  - "Hand-wrote sqlc generated files (exercise_recommendations.sql.go, day_prescribed_medications.sql.go) instead of running sqlc generate due to Windows mmap lock issue"
  - "Used *time.Time for nullable date columns per sqlc.yaml override (date → time.Time) with emit_pointers_for_null_types: true"
  - "Initialized Exercises and Prescriptions as empty slices (not nil) in dietPlanDayToDomain to maintain consistent response shape"
metrics:
  duration: "~15 minutes"
  completed: "2026-04-21"
  tasks_completed: 12
  files_changed: 14
---

# Phase 05 Plan 03: Exercise Recommendations + Prescribed Medications Per Day Summary

**One-liner:** Exercise recommendations and prescribed medications attached to DietPlanDay, exposed via CRUD endpoints and included in GET /plans/:id full tree response.

## What Was Built

Two new domain concepts attached to `DietPlanDay`:

1. **Exercise Recommendations** — day-level exercise suggestions with name, duration, description, and calorie burn estimate
2. **Prescribed Medications** — day-level medication prescriptions referencing the medications table, with dosage, frequency, times array, instructions, and optional date range

Both are loaded as part of `GetFullPlan` (GET `/plans/:id`) and appear in each day's response. Four new CRUD endpoints were added for managing them.

## New API Endpoints

| Method | Route | Description |
|--------|-------|-------------|
| POST | `/plans/:id/days/:day_id/exercises` | Add exercise recommendation to a day |
| DELETE | `/plans/:id/days/:day_id/exercises/:exercise_id` | Remove exercise recommendation |
| POST | `/plans/:id/days/:day_id/prescriptions` | Add medication prescription to a day |
| DELETE | `/plans/:id/days/:day_id/prescriptions/:prescription_id` | Remove medication prescription |

## Architecture

Follows the established DDD layering:
- **Domain** (`entity/diet_plan.go`): `ExerciseRecommendation`, `PrescribedMedication`, `MedicationSnapshot` types; `DietPlanDay.Exercises` and `DietPlanDay.Prescriptions` fields
- **Repository interface** (`repository/diet_plan_repository.go`): 8 new methods
- **Infrastructure** (`mapper.go`, `pg_diet_plan_repository.go`): mappers + implementations
- **Application** (`diet_plan_service.go`): `AddExercise`, `RemoveExercise`, `AddPrescription`, `RemovePrescription`; `GetFullPlan` extended
- **HTTP** (`diet_plan_handler.go`, `router.go`): 4 new handlers, `planFullToMap` updated with `exercises` and `prescriptions` keys

## Deviations from Plan

### Auto-fixed Issues

None — plan executed exactly as written.

**Note:** `sqlc generate` was not run (Windows mmap lock); all generated files were hand-written following the exact patterns of existing `*.sql.go` files in the codebase.

## Known Stubs

None — exercises and prescriptions are fully wired through all layers. Empty slices (`[]`) are returned when no records exist, consistent with the existing meals/options/items pattern.

## Threat Flags

None — all new endpoints follow the existing ownership-check pattern (nutritionist or superadmin). Foreign key to `medications(id) ON DELETE RESTRICT` prevents orphaned prescriptions.

## Self-Check

**Result: PASSED**

All 7 checked files exist. Commit `da37feb` verified in git log. `go build ./...` passed with exit code 0.
