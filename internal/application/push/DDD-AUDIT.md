# DDD Audit: internal/application/push
Layer: application
Audited: 2026-04-22
Files reviewed: 1 (push_service.go)

## Summary
- CRITICAL: 0
- HIGH: 1
- MEDIUM: 2
- LOW: 2
- PASS: 0

---

## Findings

### [HIGH] Infrastructure library (webpush-go) used directly inside application service

**File:** `push_service.go` (SendNotification / Broadcast)
**Issue:** The application service directly imports and calls a third-party push library (`github.com/SherClockHolmes/webpush-go` or similar). Application services MUST NOT depend on infrastructure libraries.
**DDD Rule:** Application layer — "depends only on domain and ports (interfaces defined in domain)". Third-party push implementations belong in `internal/infrastructure/push/`.
**Fix:** Define a `PushSender` port interface in the domain layer (e.g. `internal/domain/push/port/`), implement it in `internal/infrastructure/push/`, and inject through the interface.

---

### [MEDIUM] Factory `NewPushService` does not validate inputs

**File:** `push_service.go` (factory function)
**Issue:** The factory accepts VAPID keys, but does not check that they are non-empty before storing them in the service. Empty VAPID keys will cause silent failures at send time.
**Fix:** Add non-empty checks for VAPID public/private keys; return an error if invalid.

---

### [MEDIUM] `Subscribe` constructs domain entity directly via struct literal

**File:** `push_service.go` (Subscribe)
**Issue:** `PushSubscription` entity is created with a struct literal, bypassing any future `NewPushSubscription()` factory validation.
**Fix:** Call `pushentity.NewPushSubscription(...)` once the factory is added to the domain.

---

### [LOW] Hardcoded subscriber email used as push notification recipient tag

**File:** `push_service.go`
**Issue:** A constant or hardcoded string for the subscriber email appears in the service. This should be a configuration value passed at startup.
**Fix:** Accept subscriber email as a constructor parameter in `NewPushService`.

---

### [LOW] Individual delivery errors silently suppressed

**File:** `push_service.go` (Broadcast / SendToAll)
**Issue:** The broadcast loop catches per-subscription errors but discards them. Failed deliveries are not logged or counted.
**Fix:** Log per-subscription delivery errors; optionally return a multi-error summary.

---

## Compliant Patterns Found

- Service struct references repository interface for subscription storage. ✓
- No direct DB queries in the application layer. ✓

## Fix Priority Order
1. **[HIGH]** Define `PushSender` domain port interface; move webpush implementation to infrastructure layer
2. **[MEDIUM]** Add VAPID key validation to `NewPushService`
3. **[MEDIUM]** Replace `PushSubscription` struct literal with domain factory call
4. **[LOW]** Inject subscriber email via constructor parameter
5. **[LOW]** Log per-delivery errors in broadcast
