---
phase: "07"
plan: "04"
subsystem: notification-preferences
tags: [notifications, preferences, upsert, postgresql, ddd]
dependency_graph:
  requires: [07-03-push-subscriptions]
  provides: [notification-preferences-api]
  affects: [bootstrap/wire.go, router.go]
tech_stack:
  added: []
  patterns: [upsert-on-conflict, ddd-repository, gin-handler]
key_files:
  created:
    - migrations/000012_notification_preferences.up.sql
    - migrations/000012_notification_preferences.down.sql
    - db/queries/notification_preferences.sql
    - internal/infrastructure/persistence/sqlc/notification_preferences.sql.go
    - internal/domain/notification/entity/notification_preference.go
    - internal/domain/notification/repository/notification_preference_repository.go
    - internal/infrastructure/persistence/notification/mapper.go
    - internal/infrastructure/persistence/notification/pg_notification_preference_repository.go
    - internal/application/notification/notification_service.go
    - internal/interfaces/http/handler/notification_handler.go
  modified:
    - internal/infrastructure/persistence/sqlc/models.go
    - internal/domain/shared/apperror.go
    - bootstrap/wire.go
    - internal/interfaces/http/router/router.go
decisions:
  - "Used UPSERT ON CONFLICT(user_id) so first GET before upsert is never needed"
  - "ErrNotificationPreferenceNotFound uses string Code following AppError catalog pattern"
  - "GET /notifications/preferences returns ErrNotificationPreferenceNotFound when row absent"
metrics:
  duration: "~15 minutes"
  completed: "2026-04-21"
  tasks_completed: 11
  files_changed: 14
---

# Phase 7 Plan 04: Notification Preferences Summary

**One-liner:** PostgreSQL-backed notification preferences upsert + GET endpoints wired through full DDD stack (entity → repository interface → pg impl → application service → Gin handler).

## What Was Built

### Migration 000012
New `notification_preferences` table with columns `id`, `user_id` (FK → users), `meal_reminders`, `water_reminders`, `message_alerts`, `diet_updates`, `created_at`, `updated_at`. Unique constraint on `user_id` enables the upsert pattern.

### sqlc Hand-Written Layer
- `notification_preferences.sql.go` — `UpsertNotificationPreferences` and `GetNotificationPreferences` methods on `*Queries`
- `UpsertNotificationPreferencesParams` struct with `db:` tags
- `NotificationPreference` model appended to `models.go` with `db:` tags

### Domain Layer
- `entity.NotificationPreference` — pure Go struct, zero external deps
- `repository.NotificationPreferenceRepository` — interface with `Upsert` + `GetByUserID`

### Infrastructure Layer
- `PgNotificationPreferenceRepository` — implements repository interface using `db.New(pool)`
- `mapper.go` — `toDomain()` converts `db.NotificationPreference` → `entity.NotificationPreference`

### Application Layer
- `NotificationService` with `UpdatePreferences` (upsert) and `GetPreferences`

### HTTP API
| Method | Path | Description |
|--------|------|-------------|
| PATCH | /api/v1/notifications/preferences | Upsert all four boolean flags |
| GET | /api/v1/notifications/preferences | Read current preferences |

### Wiring
- `NotificationService` added to `Container` in `bootstrap/wire.go`
- Routes registered in `router.go` under protected group

### Persian Error
`ErrNotificationPreferenceNotFound` — code `NOTIFICATION_PREFERENCE_NOT_FOUND`, HTTP 404

## Deviations from Plan

**1. [Rule 1 - Bug] Corrected AppError field types**
- **Found during:** Task 10 (Persian errors)
- **Issue:** Plan spec had `Code: 404` (integer), but `AppError.Code` is a `string` and `HTTPStatus` is the int field
- **Fix:** Used `Code: "NOTIFICATION_PREFERENCE_NOT_FOUND"` + `HTTPStatus: http.StatusNotFound`
- **Files modified:** `internal/domain/shared/apperror.go`

**2. [Rule 1 - Bug] Removed spurious `time` import from notification_preferences.sql.go**
- **Found during:** Task 3
- **Issue:** Initial draft imported `"time"` but `time.Time` fields are on the struct in `models.go` (same package), so the import was unused
- **Fix:** Removed the `time` import from the `.sql.go` file

## Known Stubs

None — both endpoints wire real DB queries and return live data.

## Threat Flags

None — no new network endpoints beyond authenticated JWT-protected routes already established.

## Self-Check: PASSED

- [x] migrations/000012_notification_preferences.up.sql — exists
- [x] migrations/000012_notification_preferences.down.sql — exists
- [x] db/queries/notification_preferences.sql — exists
- [x] notification_preferences.sql.go — exists
- [x] models.go NotificationPreference struct — appended
- [x] domain entity + repository — exists
- [x] infrastructure mapper + repo — exists
- [x] application service — exists
- [x] handler — exists
- [x] apperror.go ErrNotificationPreferenceNotFound — added
- [x] wire.go NotificationService — wired
- [x] router.go routes — registered
- [x] `go build ./...` — exit code 0
- [x] git commit a1bc6be — verified
