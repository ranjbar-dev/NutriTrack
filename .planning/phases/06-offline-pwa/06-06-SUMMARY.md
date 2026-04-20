---
phase: "06"
plan: "06"
subsystem: "backend-push"
tags: [push-notifications, web-push, vapid, scheduler, reminder, websocket]
dependency_graph:
  requires: [06-05]
  provides: [push-subscription-api, notification-preferences-api, reminder-scheduler, notification-service]
  affects: [communication-service, diet-plan-service]
tech_stack:
  added: [github.com/SherClockHolmes/webpush-go@v1.4.0]
  patterns: [repository-pattern, fire-and-forget-goroutine, dedup-via-db, ticker-scheduler]
key_files:
  created:
    - backend/db/migrations/000010_create_push.up.sql
    - backend/db/migrations/000010_create_push.down.sql
    - backend/db/queries/push_subscriptions.sql
    - backend/db/queries/notification_preferences.sql
    - backend/internal/repository/sqlc/push_subscriptions.sql.go
    - backend/internal/repository/sqlc/notification_preferences.sql.go
    - backend/internal/repository/push_repo.go
    - backend/internal/model/dto/push_dto.go
    - backend/internal/service/push_service.go
    - backend/internal/service/push_service_test.go
    - backend/internal/service/reminder_scheduler.go
    - backend/internal/handler/push_handler.go
    - backend/tools/gen_vapid/main.go
  modified:
    - backend/internal/config/config.go
    - backend/internal/repository/diet_plan_repo.go
    - backend/internal/repository/communication_repo.go
    - backend/internal/service/communication_service.go
    - backend/internal/service/diet_plan_service.go
    - backend/cmd/api/main.go
    - backend/go.mod
    - backend/go.sum
    - backend/internal/repository/sqlc/models.go
    - backend/internal/repository/sqlc/querier.go
    - backend/internal/repository/sqlc/messages.sql.go
decisions:
  - "Push notification delivery is fire-and-forget (goroutines) — failures are logged but don't break HTTP responses"
  - "Reminder dedup uses DB (sent_reminders table) to prevent duplicate pushes across restarts"
  - "VAPID keys are optional config — app boots without them (pushes silently skip)"
  - "ListActivePlansWithSchedule uses raw pgx queries (no sqlc) to avoid complex JOINs in code-gen"
  - "Notification preferences default to enabled for all except water_reminder"
metrics:
  duration: "~30 minutes"
  completed: "2025-01-01"
  tasks_completed: 2
  files_changed: 26
---

# Phase 06 Plan 06: Backend Push Notification Infrastructure Summary

## One-liner
Web Push (VAPID) infrastructure with per-client preference gating, reminder scheduler, and fire-and-forget push on messages/plan activation.

## What Was Built

### Database Layer
- **Migration 000010**: Three new tables — `push_subscriptions` (one per device, upsert on re-register), `notification_preferences` (per-client boolean flags), `sent_reminders` (dedup log with 7-day retention)
- **sqlc queries**: Full CRUD for subscriptions and preferences; upsert-with-conflict patterns; EXISTS check for dedup; TTL-based purge for sent_reminders

### Repository Layer
- **PushRepository interface** with 8 methods covering subscription management, dedup tracking, and preferences
- **pgPushRepo** implementation using sqlc-generated queries with `pgtype.UUID` ↔ `string` conversion helpers
- **DietPlanRepository** extended with `ListActivePlansWithSchedule` — batched raw pgx query returning active plan schedules (meal times + JSON medication times) for the reminder scheduler

### Service Layer
- **NotificationService**: preference-gated push dispatch via `webpush-go`. Checks client's notification preferences before fetching subscriptions. Sends to all registered devices. VAPID keys from config.
- **StartReminderScheduler**: 60-second ticker checking 15-minute lookahead window for meal/medication reminders. Water reminders at fixed hours (8,10,12,14,16,18,20). All sends are fire-and-forget goroutines with DB dedup.

### Handler Layer
- **PushHandler**: 4 endpoints under `/api/client/push/`:
  - `POST /subscribe` — upsert subscription with User-Agent capture
  - `DELETE /subscribe` — remove subscription by endpoint
  - `GET /preferences` — return current preferences (defaults if none saved)
  - `PATCH /preferences` — upsert preferences

### Integration Points
- **CommunicationService.SendMessageTo**: fire-and-forget push to receiver on new message (D-18)
- **DietPlanService.ActivatePlan**: fire-and-forget push to client on plan activation
- **main.go**: scheduler started in background goroutine with context cancellation on shutdown

### Tools
- **gen_vapid**: standalone tool at `backend/tools/gen_vapid/main.go` for one-time VAPID key generation

## Tests
3 unit tests in `push_service_test.go`:
- `TestSendToClient_SkipsWhenPreferenceDisabled`: verifies no subscription fetch when preference is off
- `TestSendToClient_SkipsWhenNoSubscription`: verifies graceful no-op when no subscriptions registered
- `TestReminderAlreadySent`: verifies dedup mock logic

All 3 pass (`go test ./internal/service/... -run "TestSendToClient|TestReminderAlready"`)

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] sqlc regeneration changed ListMessages/ListMessagesSince param field names**
- **Found during:** build after sqlc generate
- **Issue:** sqlc v1.30.0 regenerated `messages.sql.go` with field names `SenderID`/`ReceiverID`/`SentAt` instead of the previous `UserA`/`UserB`/`Since`. The `communication_repo.go` still used the old names.
- **Fix:** Updated `communication_repo.go` `ListMessages` and `ListMessagesSince` calls to use new field names
- **Files modified:** `backend/internal/repository/communication_repo.go`
- **Commit:** 9140ab4

### Structural Notes

**ActivePlanSchedule types in repository package (not service)**
- The plan suggested defining `ActivePlanSchedule`, `MealScheduleItem`, `MedScheduleItem` as local types in the service package
- These are exported from the `repository` package instead, since `DietPlanRepository.ListActivePlansWithSchedule` is part of the repository interface (defining them in `service` would create a circular import)

## Known Stubs
None — all push logic is fully wired. VAPID keys not present in env will cause push sends to fail gracefully (logged, not panicked).

## Threat Flags

| Flag | File | Description |
|------|------|-------------|
| threat_flag: new-endpoint | backend/internal/handler/push_handler.go | 4 new client-authenticated endpoints for push subscription management |
| threat_flag: external-request | backend/internal/service/push_service.go | Outbound HTTPS requests to push endpoints (web push delivery) — uses client-provided endpoint URLs |

The push endpoint URLs are client-provided (from browser's PushSubscription object) and passed to the webpush library. This is standard Web Push Protocol behavior; the library handles TLS validation.

## Self-Check: PASSED

Files exist:
- backend/db/migrations/000010_create_push.up.sql ✓
- backend/db/migrations/000010_create_push.down.sql ✓
- backend/internal/repository/push_repo.go ✓
- backend/internal/service/push_service.go ✓
- backend/internal/service/reminder_scheduler.go ✓
- backend/internal/handler/push_handler.go ✓
- backend/tools/gen_vapid/main.go ✓

Commit 9140ab4 exists ✓

Build: `go build ./...` — PASS ✓
Tests: 3/3 PASS ✓
