# DDD Fix Report: internal/domain/push
Layer: domain
Fixed: 2026-04-22
Based on: DDD-AUDIT.md

## Baseline Build Status
PASS — `go build ./...` before fixes

## Fix Plan

| # | Finding | Severity | Files | Strategy | Status |
|---|---------|----------|-------|----------|--------|
| 1 | Exported fields on PushSubscription entity | HIGH | entity/push_subscription.go + 4 callers | SAFE | FIXED |
| 2 | Missing factory function with validation | HIGH | entity/push_subscription.go + application service | SAFE | FIXED |
| 3 | No getter methods — no controlled access pattern | MEDIUM | entity/push_subscription.go | SAFE | FIXED (addressed by Fix 1) |
| 4 | No domain error variables defined | LOW | entity/push_subscription.go | SAFE | FIXED (required by factory) |

## Changes Applied

### Fix 1 & 2 & 3 & 4: Make fields unexported, add factories, getters, and domain errors

**File:** `internal/domain/push/entity/push_subscription.go`
**Change:** Made all 6 struct fields unexported (`id`, `userID`, `endpoint`, `p256dh`, `auth`, `createdAt`). Added two constructors:
- `NewPushSubscription(userID uuid.UUID, endpoint, p256dh, auth string) (*PushSubscription, error)` — for creating new subscriptions with validation.
- `NewPushSubscriptionFromDB(id, userID uuid.UUID, endpoint, p256dh, auth string, createdAt time.Time) *PushSubscription` — for trusted infrastructure reconstruction from DB rows (no error return).

Added value-receiver getters for all 6 fields (`GetID`, `GetUserID`, `GetEndpoint`, `GetP256dh`, `GetAuth`, `GetCreatedAt`). Added domain error vars: `ErrEmptyEndpoint`, `ErrEmptyP256dh`, `ErrEmptyAuth`, `ErrSubscriptionNotFound`.

**Before:**
```go
type PushSubscription struct {
    ID        uuid.UUID
    UserID    uuid.UUID
    Endpoint  string
    P256dh    string
    Auth      string
    CreatedAt time.Time
}
```
**After:**
```go
type PushSubscription struct {
    id        uuid.UUID
    userID    uuid.UUID
    endpoint  string
    p256dh    string
    auth      string
    createdAt time.Time
}

func NewPushSubscription(userID uuid.UUID, endpoint, p256dh, auth string) (*PushSubscription, error) { ... }
func NewPushSubscriptionFromDB(id, userID uuid.UUID, ..., createdAt time.Time) *PushSubscription { ... }
func (ps PushSubscription) GetID() uuid.UUID { ... }
// + 5 more getters
```
**Build:** PASS

### Fix 1 (caller): infrastructure/persistence/push/mapper.go
**Change:** Replaced struct-literal construction with `NewPushSubscriptionFromDB` call.
**Before:** `return &entity.PushSubscription{ID: row.ID, ...}`
**After:** `return entity.NewPushSubscriptionFromDB(row.ID, row.UserID, row.Endpoint, row.P256dh, row.Auth, row.CreatedAt)`
**Build:** PASS

### Fix 1 (caller): infrastructure/persistence/push/pg_push_subscription_repository.go
**Change:** Replaced direct field reads (`sub.UserID`, `sub.Endpoint`, etc.) with getter calls (`sub.GetUserID()`, `sub.GetEndpoint()`, etc.) in `Upsert`.
**Build:** PASS

### Fix 2 (caller): application/push/push_service.go
**Change:** In `Subscribe`, replaced direct struct construction with `NewPushSubscription` factory call. In `Send`, replaced direct field reads (`sub.Endpoint`, `sub.Auth`, `sub.P256dh`) with getter calls.
**Before:** `sub := &entity.PushSubscription{UserID: userID, Endpoint: endpoint, ...}`
**After:** `sub, err := entity.NewPushSubscription(userID, endpoint, p256dh, auth); if err != nil { return nil, err }`
**Build:** PASS

### Fix 1 (caller): interfaces/http/handler/push_handler.go
**Change:** Replaced direct field reads (`sub.ID`, `sub.UserID`, `sub.Endpoint`, `sub.CreatedAt`) with getter calls in `Subscribe` handler.
**Build:** PASS

## Final Build Status
PASS — `go build ./...` after all fixes
PASS — `go vet ./internal/...` after all fixes

## Remaining Violations
None — all CRITICAL, HIGH, and MEDIUM findings resolved. LOW (domain errors) also resolved as a side-effect of factory validation requirements.
