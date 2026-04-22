# DDD Fix Report: internal/domain/notification
Layer: domain
Fixed: 2026-04-22
Based on: DDD-AUDIT.md

## Baseline Build Status
PASS — `go build ./...` before fixes

## Fix Plan

| # | Finding | Severity | Files | Strategy | Status |
|---|---------|----------|-------|----------|--------|
| 1 | NotificationPreference aggregate exposes all fields as exported | HIGH | entity/notification_preference.go + 4 callers | SAFE | FIXED |
| 2 | Missing factory function `NewNotificationPreference()` | HIGH | entity/notification_preference.go + application service | SAFE | FIXED |
| 3 | No domain error variables defined | LOW | entity/notification_preference.go | SAFE | FIXED (minimal errors required by factory) |

## Changes Applied

### Fix 1 & 2: Make fields unexported, add getters/setters, add factory

**File:** `internal/domain/notification/entity/notification_preference.go`
**Change:** Made all 6 struct fields unexported (`id`, `userID`, `mealReminders`, `waterReminders`, `messageAlerts`, `dietUpdates`). Added value-receiver getters for all fields, pointer-receiver setters for the 4 boolean fields. Added `NewNotificationPreference(id, userID uuid.UUID) (*NotificationPreference, error)` factory that validates `userID != uuid.Nil`. Added minimal domain error vars required by the factory.
**Before:**
```go
type NotificationPreference struct {
    ID             uuid.UUID
    UserID         uuid.UUID
    MealReminders  bool
    WaterReminders bool
    MessageAlerts  bool
    DietUpdates    bool
}
```
**After:**
```go
type NotificationPreference struct {
    id             uuid.UUID
    userID         uuid.UUID
    mealReminders  bool
    waterReminders bool
    messageAlerts  bool
    dietUpdates    bool
}

func NewNotificationPreference(id, userID uuid.UUID) (*NotificationPreference, error) { ... }
func (np NotificationPreference) GetID() uuid.UUID { ... }
// + 5 more getters, 4 setters
```
**Build:** PASS

### Fix 1 (caller): infrastructure/persistence/notification/mapper.go
**Change:** Replaced struct-literal construction with `NewNotificationPreference` factory + setters.
**Before:** `return entity.NotificationPreference{ID: row.ID, ...}`
**After:** `pref, _ := entity.NewNotificationPreference(row.ID, row.UserID); pref.SetMealReminders(...); ...; return *pref`
**Build:** PASS

### Fix 1 (caller): infrastructure/persistence/notification/pg_notification_preference_repository.go
**Change:** Replaced direct field access (`pref.UserID`, etc.) with getter calls (`pref.GetUserID()`, etc.) in `Upsert`.
**Build:** PASS

### Fix 2 (caller): application/notification/notification_service.go
**Change:** Replaced struct-literal construction in `UpdatePreferences` with factory + setters pattern.
**Before:** `entity.NotificationPreference{UserID: userID, MealReminders: ..., ...}`
**After:** `pref, err := entity.NewNotificationPreference(uuid.Nil, userID); ...; pref.SetMealReminders(...); ...; return s.prefRepo.Upsert(ctx, *pref)`
**Build:** PASS

### Fix 1 (caller): interfaces/http/handler/notification_handler.go
**Change:** Replaced direct field reads (`pref.ID`, `pref.UserID`, etc.) with getter calls in both `UpdatePreferences` and `GetPreferences` handlers.
**Build:** PASS

## Final Build Status
PASS — `go build ./...` after all fixes
PASS — `go vet ./internal/...` after all fixes

## Remaining Violations
None — all CRITICAL and HIGH findings resolved.
