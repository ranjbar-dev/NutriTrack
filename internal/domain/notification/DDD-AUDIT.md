# DDD Audit: internal/domain/notification
Layer: domain
Audited: 2026-04-22
Files reviewed: 2

## Summary
- CRITICAL: 0
- HIGH: 2
- MEDIUM: 0
- LOW: 1
- PASS: 1 file with no issues (`repository/notification_preference_repository.go`)

---

## Findings

### [HIGH] NotificationPreference aggregate exposes all fields as exported

**File:** `entity/notification_preference.go:6`
**Issue:** All six fields (`ID`, `UserID`, `MealReminders`, `WaterReminders`, `MessageAlerts`, `DietUpdates`) are exported. Direct field access bypasses any invariant enforcement.
**DDD Rule:** Aggregates: unexported fields, getter/setter methods for controlled access.
**Fix:** Make all fields unexported and add getter/setter methods.

---

### [HIGH] Missing factory function `NewNotificationPreference()`

**File:** `entity/notification_preference.go` (missing)
**Issue:** No `New*()` factory function exists. Callers construct the aggregate directly, bypassing validation.
**DDD Rule:** Factory functions `New*()` MUST validate required fields and return `(T, error)` or `(*T, error)`.
**Fix:** Add `NewNotificationPreference(id, userID uuid.UUID) (*NotificationPreference, error)` with UUID validation.

---

### [LOW] No domain error variables defined

**File:** `entity/notification_preference.go` (missing)
**Issue:** No sentinel errors defined.
**Fix:** Add `ErrNotificationPreferenceNotFound`, `ErrInvalidNotificationPreferenceID`, `ErrInvalidUserID`.

---

## Compliant Patterns Found
- `repository/notification_preference_repository.go` is a proper Go `interface`. ✓
- No `json:`, `bson:`, `db:` struct tags anywhere in domain structs. ✓
- No forbidden cross-layer imports. ✓
- `valueobject/` empty package, no violations. ✓

## Fix Priority Order
1. **[HIGH]** Make `NotificationPreference` fields unexported and add getter/setter methods
2. **[HIGH]** Add `NewNotificationPreference()` factory with UUID validation
3. **[LOW]** Define sentinel domain error variables
