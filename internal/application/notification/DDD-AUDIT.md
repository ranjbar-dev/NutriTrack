# DDD Audit: internal/application/notification
Layer: application
Audited: 2026-04-22
Files reviewed: 1 (notification_service.go)

## Summary
- CRITICAL: 0
- HIGH: 1
- MEDIUM: 1
- LOW: 0
- PASS: 0

---

## Findings

### [HIGH] Service constructs domain entity via raw struct literal (no factory)

**File:** `notification_service.go`
**Issue:** The application service instantiates the `NotificationPreference` entity with a raw struct literal (e.g. `notificationentity.NotificationPreference{ID: ..., UserID: ...}`). This bypasses all factory-level validation.
**DDD Rule:** Application services MUST call the domain's `New*()` factory to construct aggregates.
**Fix:** Call `notificationentity.NewNotificationPreference(id, userID)` once the factory is added.

---

### [MEDIUM] Factory `NewNotificationService` does not return an error

**File:** `notification_service.go` (factory function)
**Issue:** The service constructor `NewNotificationService(repo)` silently ignores nil dependency injection. If a nil repository is passed, the application will panic at runtime when the first method is invoked.
**DDD Rule:** Application service constructors SHOULD validate required dependencies and return an error for invalid configuration.
**Fix:** Change the signature to `NewNotificationService(repo notificationrepository.NotificationPreferenceRepository) (*NotificationService, error)` and return an error if `repo == nil`.

---

## Compliant Patterns Found

- Service constructor accepts a repository interface (not a concrete type). ✓
- Service file imports domain entity and repository packages only. ✓

## Fix Priority Order
1. **[HIGH]** Replace struct literal with `NewNotificationPreference()` factory call
2. **[MEDIUM]** Add nil-check and error return to `NewNotificationService` constructor
