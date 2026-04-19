---
phase: 02-core-data-domain
plan: 04
subsystem: backend/medications
tags: [go, gin, sqlc, postgresql, crud, persian-search, row-level-auth]
dependency_graph:
  requires: [02-01]
  provides: [medication-crud-api]
  affects: [05-medication-frontend, 07-super-admin]
tech_stack:
  added: []
  patterns: [handler→service→repository, sqlc-typed-queries, persian-fuzzy-search, row-level-auth]
key_files:
  created:
    - backend/db/queries/medications.sql
    - backend/internal/repository/sqlc/medications.sql.go
    - backend/internal/model/dto/medication_dto.go
    - backend/internal/repository/medication_repo.go
    - backend/internal/service/medication_service.go
    - backend/internal/handler/medication_handler.go
  modified:
    - backend/cmd/api/main.go
decisions:
  - "[02-04]: Reused optionalText/optionalBool/formatTimestamp helpers from food_service.go — same package, no duplication needed"
  - "[02-04]: Duplicate check uses COALESCE pattern (same as food) to handle NULL exclude_id safely in SQL"
metrics:
  duration: 12min
  completed: 2026-04-19
  tasks_completed: 3
  files_changed: 7
---

# Phase 02 Plan 04: Medication CRUD Backend Summary

**One-liner:** Full medication CRUD API with Persian fuzzy search on name+generic_name, row-level auth (nutritionist=own, super_admin=any), and soft delete.

## What Was Built

Complete `/api/medications` REST endpoints following the identical layered architecture pattern established by the food backend (Plan 02-02).

### Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | /api/medications | nutritionist, super_admin | Create medication |
| GET | /api/medications | nutritionist, super_admin | List with Persian search, pagination |
| GET | /api/medications/:id | nutritionist, super_admin | Get single medication |
| PUT | /api/medications/:id | nutritionist, super_admin | Update with row-level auth |
| DELETE | /api/medications/:id | nutritionist, super_admin | Soft delete with row-level auth |

### SQL Queries (sqlc)

8 queries in `medications.sql`:
- `CreateMedication` — INSERT with auto-computed name_normalized and generic_name_normalized
- `GetMedicationByID` — JOIN with users for creator_name
- `ListMedications` — Persian ILIKE search on both name_normalized and generic_name_normalized, similarity ordering
- `CountMedications` — Same WHERE for pagination totals
- `UpdateMedication` — Full update with normalized recomputation
- `SoftDeleteMedication` — Super admin path (any record)
- `SoftDeleteMedicationByOwner` — Nutritionist path (own records only)
- `CheckDuplicateMedicationName` — COALESCE pattern for safe NULL exclusion

### Business Logic

- **Duplicate detection:** Persian-normalized name comparison, excludes self on update
- **Row-level auth (D-29):** Nutritionist → own items only; Super Admin → any item
- **Pagination:** Default page=1, limit=20, cap at 100
- **Persian search:** Searches both `name_normalized` AND `generic_name_normalized`
- **All errors in Persian** per project requirement

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Removed duplicate helper function declarations**
- **Found during:** Task 2 build
- **Issue:** `optionalText` and `optionalBool` are already declared in `food_service.go` (same package). Adding them again caused compilation error.
- **Fix:** Removed duplicate declarations from `medication_service.go`. Functions are shared within the package.
- **Files modified:** `backend/internal/service/medication_service.go`
- **Commit:** 0ae4c52

## Self-Check

### Created Files Exist
- ✅ backend/db/queries/medications.sql
- ✅ backend/internal/repository/sqlc/medications.sql.go
- ✅ backend/internal/model/dto/medication_dto.go
- ✅ backend/internal/repository/medication_repo.go
- ✅ backend/internal/service/medication_service.go
- ✅ backend/internal/handler/medication_handler.go

### Commits
- ada5e6e: feat(02-04): add medication SQL queries, sqlc generated code, and DTOs
- 0ae4c52: feat(02-04): add medication repository and service
- ae34226: feat(02-04): add medication handler and wire /api/medications routes

### Build Status
- `go build ./...` ✅
- `go vet ./...` ✅
- `sqlc generate` ✅

## Self-Check: PASSED
