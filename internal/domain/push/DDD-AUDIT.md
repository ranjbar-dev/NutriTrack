# DDD Audit: internal/domain/push
Layer: domain
Audited: 2026-04-22
Files reviewed: 2

## Summary
- CRITICAL: 0
- HIGH: 2
- MEDIUM: 1
- LOW: 1
- PASS: 1 file (`repository/push_subscription_repository.go`)

---

## Findings

### [HIGH] Exported fields on PushSubscription entity

**File:** `entity/push_subscription.go:11-17`
**Issue:** All six fields of `PushSubscription` are exported (`ID`, `UserID`, `Endpoint`, `P256dh`, `Auth`, `CreatedAt`). Any caller can freely read or overwrite any field.
**DDD Rule:** Entities — "unexported fields, mutable state via exported methods"
**Fix:** Make all fields unexported and add getter methods.

---

### [HIGH] Missing factory function with validation

**File:** `entity/push_subscription.go` (entire file — absent)
**Issue:** There is no `NewPushSubscription(...)` constructor. Without a factory, callers can create zero-value or invalid `PushSubscription` structs.
**DDD Rule:** Aggregates — "factory `New*()` function that validates inputs and returns an error"
**Fix:** Add factory validating `endpoint`, `p256dh`, `auth` fields. Add domain errors: `ErrEmptyEndpoint`, `ErrEmptyP256dh`, `ErrEmptyAuth`.

---

### [MEDIUM] No getter methods — no controlled access pattern

**File:** `entity/push_subscription.go`
**Issue:** Because fields are exported, there are no getter methods. Once fields are made unexported, callers will need getters.
**DDD Rule:** Aggregates — "getter/setter methods for controlled access"
**Fix:** Add read-only getters (addressed by the HIGH fix above).

---

### [LOW] No domain error variables defined

**File:** package-wide (absent)
**Issue:** No sentinel errors (`var Err* = errors.New(...)`).
**Fix:** Add `ErrEmptyEndpoint`, `ErrEmptyP256dh`, `ErrEmptyAuth`, `ErrSubscriptionNotFound`.

---

## Compliant Patterns Found

- **`repository/push_subscription_repository.go`** — Fully compliant, pure Go `interface`. ✓
- **No `json:`/`bson:`/`db:` struct tags** on `PushSubscription`. ✓
- **No imports of forbidden layers**. ✓
- **UUID identifier field** present. ✓

## Fix Priority Order
1. Make `PushSubscription` fields unexported
2. Add `NewPushSubscription()` factory with field validation
3. Add getter methods
4. Add `ErrSubscriptionNotFound` sentinel error
