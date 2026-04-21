---
phase: "06"
plan: "01"
subsystem: tracking
tags: [tracking, domain, infrastructure, sqlc, migrations, postgresql]
dependency_graph:
  requires: [05-dietplan, internal/domain/shared, internal/infrastructure/persistence/sqlc]
  provides: [internal/domain/tracking, internal/infrastructure/persistence/tracking]
  affects: []
tech_stack:
  added: []
  patterns: [xmax-insert-detection, ON-CONFLICT-DO-UPDATE, offline-sync-idempotency, DDD-repository-pattern]
key_files:
  created:
    - migrations/000007_tracking.up.sql
    - migrations/000007_tracking.down.sql
    - db/queries/food_logs.sql
    - db/queries/water_logs.sql
    - db/queries/sleep_logs.sql
    - db/queries/exercise_logs.sql
    - db/queries/medication_logs.sql
    - db/queries/body_measurements.sql
    - internal/infrastructure/persistence/sqlc/food_logs.sql.go
    - internal/infrastructure/persistence/sqlc/water_logs.sql.go
    - internal/infrastructure/persistence/sqlc/sleep_logs.sql.go
    - internal/infrastructure/persistence/sqlc/exercise_logs.sql.go
    - internal/infrastructure/persistence/sqlc/medication_logs.sql.go
    - internal/infrastructure/persistence/sqlc/body_measurements.sql.go
    - internal/domain/tracking/entity/tracking.go
    - internal/domain/tracking/repository/tracking_repository.go
    - internal/infrastructure/persistence/tracking/mapper.go
    - internal/infrastructure/persistence/tracking/pg_tracking_repository.go
  modified:
    - internal/infrastructure/persistence/sqlc/models.go
decisions:
  - "xmax trick used in all upsert queries to detect insert vs conflict without extra SELECT"
  - "ON CONFLICT DO UPDATE SET client_id = EXCLUDED.client_id ensures row always returned via RETURNING"
  - "UpsertXxx returns (inserted bool, err error) — entity populated in-place from DB row"
  - "Repository ListXxx methods always return empty slice (not nil) on no results"
metrics:
  duration: "~25 minutes"
  completed: "2026-04-21"
  tasks_completed: 8
  files_changed: 19
---

# Phase 06 Plan 01: Tracking Domain Aggregates, Migrations, and Infrastructure Summary

**One-liner:** Six tracking tables (food/water/sleep/exercise/medication/body) with offline-sync idempotency via UNIQUE(client_id, local_id) and xmax-based insert detection in hand-written sqlc files.

## What Was Built

### Migration 000007
- Created 6 tracking tables: `food_logs`, `water_logs`, `sleep_logs`, `exercise_logs`, `medication_logs`, `body_measurements`
- Each table has `UNIQUE(client_id, local_id)` for offline-sync idempotency
- Appropriate indexes on `(client_id, date_column)` for efficient date-range queries
- Nullable FK references: `food_id → foods(id)`, `medication_id → medications(id)` with `ON DELETE SET NULL`

### sqlc Query Files (db/queries/)
Six SQL files defining 4 operations each:
- `UpsertXxx` — `ON CONFLICT DO UPDATE` with `RETURNING *` plus xmax insert detection
- `ListXxxByClientAndDate` — filter by client + date, ordered by time
- `ListXxxByClient` — paginated list ordered by date DESC
- `CountXxxByClient` — total count for pagination metadata

### Hand-written sqlc Go Files
Six `.sql.go` files in `internal/infrastructure/persistence/sqlc/` following exact sqlc v1.30.0 patterns:
- All upsert functions return `(Model, bool, error)` — the bool indicates `true` if newly inserted
- xmax trick: `(xmax::text::bigint = 0) AS inserted` appended to RETURNING clause, scanned into `&inserted bool`

### models.go Update
Appended 6 new struct types: `FoodLog`, `WaterLog`, `SleepLog`, `ExerciseLog`, `MedicationLog`, `BodyMeasurement`

### Domain Layer
- `internal/domain/tracking/entity/tracking.go` — 6 entity types with Go-native types (float64, int, *float64)
- `internal/domain/tracking/repository/tracking_repository.go` — `TrackingRepository` interface with 18 methods (3 per tracking type)

### Infrastructure Layer
- `mapper.go` — bidirectional conversion helpers (pgtype.Numeric ↔ float64, *pgtype.Numeric ↔ *float64)
- `pg_tracking_repository.go` — `PgTrackingRepository` implementing all 18 repository methods

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] PowerShell backtick escaping in models.go append**
- **Found during:** Task 3
- **Issue:** Initial PowerShell script using `+"\`"+` pattern for backtick struct tags resulted in literal `+""+` characters written to models.go, causing syntax errors
- **Fix:** Used PowerShell here-string `@"..."@` syntax which treats double backticks (`\`\``) as literal single backtick characters
- **Files modified:** `internal/infrastructure/persistence/sqlc/models.go`
- **Commit:** Same task commit (fixed before committing)

## Build Verification

```
go build ./... → exit code 0
```

## Self-Check: PASSED

### Files verified:
- `migrations/000007_tracking.up.sql` ✅
- `migrations/000007_tracking.down.sql` ✅
- `internal/domain/tracking/entity/tracking.go` ✅
- `internal/domain/tracking/repository/tracking_repository.go` ✅
- `internal/infrastructure/persistence/tracking/mapper.go` ✅
- `internal/infrastructure/persistence/tracking/pg_tracking_repository.go` ✅

### Commits:
- `3fcff00` — feat(06-01): tracking domain aggregates, migrations, sqlc queries, infrastructure
