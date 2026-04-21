---
phase: "08"
plan: "01"
subsystem: admin
tags: [admin, stats, dashboard, food-management, medication-management]
dependency_graph:
  requires: [users table, foods table, diet_plans table, existing FoodService, MedicationService]
  provides: [GET /admin/stats, GET /admin/foods, DELETE /admin/foods/:id, GET /admin/medications, DELETE /admin/medications/:id]
  affects: [bootstrap/wire.go, router.go]
tech_stack:
  added: [internal/application/admin package]
  patterns: [sqlc hand-written query, DDD application service, route reuse via existing handlers]
key_files:
  created:
    - db/queries/admin_stats.sql
    - internal/infrastructure/persistence/sqlc/admin_stats.sql.go
    - internal/application/admin/admin_service.go
    - internal/interfaces/http/handler/admin_handler.go
  modified:
    - bootstrap/wire.go
    - internal/interfaces/http/router/router.go
decisions:
  - AdminService uses db.Queries directly (not a domain repository interface) — stats query is read-only aggregation with no domain logic
  - Single combined SQL query for stats (subselects for foods/diet_plans counts) — avoids N+1
  - Admin food/medication routes reuse existing FoodHandler.Search/Delete and MedicationHandler.Search/Delete — no duplicated handler code
metrics:
  duration: "10m"
  completed_date: "2026-04-21"
  tasks_completed: 8
  files_changed: 6
---

# Phase 8 Plan 01: Admin Stats & Food/Medication Management Summary

## One-liner

Super-admin dashboard stats endpoint (`GET /admin/stats`) plus admin-scoped food/medication list and force-delete routes, wired through a new `AdminService` backed by a single combined SQL query.

## What Was Built

### `GET /api/v1/admin/stats`

Returns aggregated platform counts in one SQL round-trip:
- `total_nutritionists`, `active_nutritionists`, `inactive_nutritionists`
- `total_clients`
- `total_foods`
- `active_diet_plans`

All counts come from a single `SELECT … FILTER` query on `users` with subselects for `foods` and `diet_plans`.

### Admin Food/Medication Routes

Added to the existing `adminGroup` (behind `RequireRole(super_admin)` middleware):

| Method | Path | Handler |
|--------|------|---------|
| `GET` | `/admin/stats` | `AdminHandler.GetStats` |
| `GET` | `/admin/foods` | `FoodHandler.Search` (reused) |
| `DELETE` | `/admin/foods/:id` | `FoodHandler.Delete` (reused) |
| `GET` | `/admin/medications` | `MedicationHandler.Search` (reused) |
| `DELETE` | `/admin/medications/:id` | `MedicationHandler.Delete` (reused) |

Food/medication service layer already encodes role-based ownership logic; `super_admin` callerRole bypasses ownership checks.

## Files Created

| File | Purpose |
|------|---------|
| `db/queries/admin_stats.sql` | SQL source for the combined stats query |
| `internal/infrastructure/persistence/sqlc/admin_stats.sql.go` | Hand-written sqlc file; `GetAdminStats()` on `*Queries` |
| `internal/application/admin/admin_service.go` | `AdminService` with `GetStats()` |
| `internal/interfaces/http/handler/admin_handler.go` | `AdminHandler.GetStats` HTTP handler |

## Files Modified

| File | Change |
|------|--------|
| `bootstrap/wire.go` | Added `dbsqlc` import alias; wired `AdminService`; added `AdminService` field to `Container` |
| `internal/interfaces/http/router/router.go` | Added `adminHandler`, `GET /admin/stats`; added admin food/med routes after respective handler declarations |

## Deviations from Plan

None — plan executed exactly as written. The `q.db` field is `DBTX` interface (not `*pgxpool.Pool`) so `QueryRow` is called via the interface, which is compatible with both `*pgxpool.Pool` and `pgx.Tx`. No adjustment needed.

## Self-Check: PASSED

- `db/queries/admin_stats.sql` ✅
- `internal/infrastructure/persistence/sqlc/admin_stats.sql.go` ✅
- `internal/application/admin/admin_service.go` ✅
- `internal/interfaces/http/handler/admin_handler.go` ✅
- `bootstrap/wire.go` ✅
- `internal/interfaces/http/router/router.go` ✅
- Commit `cb6e6c5` ✅
- `go build ./...` passes ✅
