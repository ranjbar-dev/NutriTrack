---
phase: "04"
plan: "03"
subsystem: medication-domain
tags: [medication, ddd, pg_trgm, soft-delete, crud, row-level-isolation]
dependency_graph:
  requires: [04-01, 04-02, 01-foundation]
  provides: [medication-aggregate, medication-crud-api]
  affects: [bootstrap, router]
tech_stack:
  added: []
  patterns: [ddd-aggregate, repository-interface, soft-delete, pg_trgm-search, row-level-write-isolation]
key_files:
  created:
    - migrations/000004_medications.up.sql
    - migrations/000004_medications.down.sql
    - db/queries/medications.sql
    - internal/infrastructure/persistence/sqlc/medications.sql.go
    - internal/domain/medication/entity/medication.go
    - internal/domain/medication/repository/medication_repository.go
    - internal/infrastructure/persistence/medication/mapper.go
    - internal/infrastructure/persistence/medication/pg_medication_repository.go
    - internal/application/medication/medication_service.go
    - internal/interfaces/http/handler/medication_handler.go
  modified:
    - internal/infrastructure/persistence/sqlc/models.go
    - bootstrap/wire.go
    - internal/interfaces/http/router/router.go
decisions:
  - "Soft delete (is_active=false) for nutritionist-owned medications; hard DELETE for superadmin only — mirrors food pattern"
  - "Row-level write isolation: nutritionist can only update/delete their own medications (CreatedBy check in service)"
  - "sqlc manually authored for medications.sql.go (VSCode/gopls had existing files memory-mapped, blocking sqlc overwrite)"
metrics:
  duration: "~20 min"
  completed: "2026-04-21"
  tasks_completed: 11
  files_changed: 13
---

# Phase 4 Plan 03: Medication Domain Aggregate — CRUD + pg_trgm Search + Soft Delete Summary

**One-liner:** Full DDD medication aggregate with pg_trgm similarity search, soft-delete soft isolation, and role-based CRUD over 5 REST endpoints.

## What Was Built

The Medication domain aggregate implemented end-to-end following the project's DDD layering:

- **Migration 000004** — `medications` table with `pg_trgm` GIN index on `name_normalized`, `created_by` FK, `is_active` soft-delete flag.
- **sqlc queries** — 7 queries: CreateMedication, GetMedicationByID, SearchMedications (pg_trgm similarity + ILIKE fallback), CountSearchMedications, UpdateMedication, DeactivateMedication, DeleteMedication.
- **Domain layer** — `entity.Medication` struct (zero external deps), `MedicationRepository` interface.
- **Infrastructure layer** — `PgMedicationRepository` implementing all 7 operations; `mapper.go` converts `db.Medication` (pgtype.UUID for `created_by`) to domain entity.
- **Application service** — `MedicationService` with 5 business operations enforcing row-level write isolation (nutritionist → own records only; superadmin → any).
- **HTTP handler** — 5 endpoints wired into the protected route group (any authenticated user; role logic in service).
- **Bootstrap + router** — `MedicationService` added to `Container`, wired in `NewContainer`, medication routes added to `router.go`.

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 3 - Blocking Issue] sqlc generate blocked by VSCode/gopls memory-mapped file locks**
- **Found during:** Step 3 (sqlc generate)
- **Issue:** `sqlc generate` rewrites all existing `.go` files in the sqlc output directory atomically. VSCode and gopls had existing files (`db.go`, `health.sql.go`, `models.go`, `users.sql.go`, `foods.sql.go`) memory-mapped, which prevented the write operations (`The requested operation cannot be performed on a file with a user-mapped section open`). Killing gopls did not release the locks.
- **Fix:** Manually authored `medications.sql.go` (new file — no locking conflict) following the exact sqlc code generation pattern derived from `foods.sql.go`. Updated `models.go` by writing the full file content via `[System.IO.File]::WriteAllText()` (bypasses mmap write restriction). The `Medication` struct was appended; all existing structs preserved verbatim.
- **Files modified:** `internal/infrastructure/persistence/sqlc/medications.sql.go` (created), `internal/infrastructure/persistence/sqlc/models.go` (updated)
- **Verification:** `go build ./...` succeeded with exit code 0.

## Known Stubs

None — all data sources are wired to the real database.

## Self-Check

**Created files exist:**
- `migrations/000004_medications.up.sql` ✅
- `internal/domain/medication/entity/medication.go` ✅
- `internal/application/medication/medication_service.go` ✅
- `internal/interfaces/http/handler/medication_handler.go` ✅

**Commit exists:**
- `c97fab1` feat(04-03): medication domain aggregate, pg_trgm search, CRUD handlers ✅

**Build:** `go build ./...` → exit 0 ✅

## Self-Check: PASSED
