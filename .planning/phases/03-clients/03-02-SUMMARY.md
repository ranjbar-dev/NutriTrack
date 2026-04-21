---
phase: "03"
plan: "02"
subsystem: "client-management"
tags: [client, nutritionist, crud, gin, ddd]
dependency_graph:
  requires: [02-auth, 03-01-nutritionist-management]
  provides: [client-registration, client-profile, client-list, client-update]
  affects: [diet-plan-phase]
tech_stack:
  added: []
  patterns: [DDD application-service, row-level-security, partial-update]
key_files:
  created:
    - internal/application/user/client_service.go
    - internal/interfaces/http/handler/client_handler.go
  modified:
    - bootstrap/wire.go
    - internal/interfaces/http/router/router.go
decisions:
  - "Ownership check (BelongsTo) performed in service layer, not handler, to enforce DDD boundary"
  - "birth_date parsed as 2006-01-02 string in handler; stored as *time.Time in domain"
  - "Context keys use middleware.AuthUserIDKey constant — not raw string — for type-safety"
metrics:
  duration: "~10 min"
  completed: "2026-04-21"
  tasks_completed: 1
  files_changed: 4
---

# Phase 03 Plan 02: Client Management Summary

**One-liner:** Nutritionist-scoped client CRUD — register/list/profile/update with row-level ownership enforcement via `BelongsTo()`.

## What Was Built

Four application-layer methods and four HTTP endpoints enabling a nutritionist to manage their clients:

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/clients` | Register a new client (mobile-only, no password) |
| GET | `/api/v1/clients` | List paginated clients for the caller nutritionist |
| GET | `/api/v1/clients/:id` | Get client profile (includes computed BMI) |
| PATCH | `/api/v1/clients/:id` | Partial update of client profile fields |

All routes are gated by `RequireAuth` + `RequireRole(nutritionist)`.

## Files Created

### `internal/application/user/client_service.go`
- `RegisterClient` — normalizes mobile via `valueobject.NewMobile`, checks `ExistsByMobile`, creates client with `NutritionistID` set, no email/password
- `GetClientProfile` — fetches user, asserts `IsClient()`, then `BelongsTo(nutritionistID)`
- `UpdateClient` — same ownership checks, applies partial field updates, calls `userRepo.Update`
- `ListClients` — delegates to `FindClientsByNutritionist` + `CountClientsByNutritionist`

### `internal/interfaces/http/handler/client_handler.go`
- Extracts `nutritionistID` from `middleware.AuthUserIDKey` context key
- Parses `birth_date` as `"2006-01-02"` string → `*time.Time`
- `toClientResponse()` helper maps domain `User` → `gin.H` with `bmi`, `full_name`, nullable fields

## Files Modified

### `bootstrap/wire.go`
- Added `ClientService *appUser.ClientService` field to `Container`
- Instantiated `clientSvc := appUser.NewClientService(userRepo)` and set in return struct

### `internal/interfaces/http/router/router.go`
- Added `/clients` group under `protected` with `RequireRole(RoleNutritionist)`
- Registered POST/GET/GET:id/PATCH:id routes

## Deviations from Plan

None — plan executed exactly as written.

## Known Stubs

None — all endpoints wire to real repository methods backed by the existing PgUserRepository implementation.

## Threat Flags

| Flag | File | Description |
|------|------|-------------|
| threat_flag: row-level-access | client_service.go | BelongsTo check enforced in service — handlers never return cross-nutritionist data |

## Self-Check: PASSED

- `internal/application/user/client_service.go` — FOUND
- `internal/interfaces/http/handler/client_handler.go` — FOUND
- `bootstrap/wire.go` — ClientService field present
- `internal/interfaces/http/router/router.go` — /clients routes registered
- Commit `88a9463` — FOUND (`feat(03-02): client management — register, list, profile, update`)
- `go build ./...` — PASSED (exit code 0)
