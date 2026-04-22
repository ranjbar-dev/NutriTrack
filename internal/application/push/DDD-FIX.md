# DDD Fix Report: internal/application/push
Layer: application
Fixed: 2026-04-22
Based on: DDD-AUDIT.md

## Baseline Build Status
PASS — `go build ./...` before fixes

## Fix Plan

| # | Finding | Severity | Files | Strategy | Status |
|---|---------|----------|-------|----------|--------|
| 1 | webpush-go used directly in app service | HIGH | push_service.go | SAFE | FIXED |
| 2 | NewPushService does not validate VAPID keys | MEDIUM | push_service.go | N/A | ABSORBED by fix 1 |
| 3 | PushSubscription struct literal | MEDIUM | push_service.go | N/A | PRE-FIXED |
| 4 | Hardcoded subscriber email | LOW | push_service.go | FIXED (moved to infra constructor) |
| 5 | Individual delivery errors suppressed | LOW | push_service.go | DEFERRED | DEFERRED |

## Changes Applied

### Fix 1: Define PushSender domain port; move webpush to infrastructure layer

**New file:** `internal/domain/push/port/push_sender.go`
```go
package port

type PushSender interface {
    SendToSubscription(ctx context.Context, endpoint, p256dh, auth string, payload []byte) error
}
```

**New file:** `internal/infrastructure/push/webpush_sender.go`
- `WebpushSender` struct holds VAPID keys and subscriber email
- `NewWebpushSender(publicKey, privateKey, subscriberEmail string) *WebpushSender`
- Implements `pushPort.PushSender` via compile-time check: `var _ pushPort.PushSender = (*WebpushSender)(nil)`
- `SendToSubscription` wraps `webpush.SendNotification` from the third-party library

**Updated:** `internal/application/push/push_service.go`
- Removed `webpush-go` import
- Field changed from `vapidPublicKey/vapidPrivateKey string` → `sender pushPort.PushSender`
- `NewPushService(repo, sender pushPort.PushSender)` — cleaner constructor
- `Send()` delegates to `s.sender.SendToSubscription(...)`

**Updated:** `bootstrap/wire.go`
- Added import `infraPush "github.com/ranjbar-dev/nutritrack/internal/infrastructure/push"`
- Added `webpushSender := infraPush.NewWebpushSender(cfg.VAPID.PublicKey, cfg.VAPID.PrivateKey, "mailto:info@nutritrack.ir")`
- Changed `appPush.NewPushService(pgPushRepo, cfg.VAPID.PublicKey, cfg.VAPID.PrivateKey)` → `appPush.NewPushService(pgPushRepo, webpushSender)`

**Build:** PASS

### LOW Fix: Subscriber email moved to infra constructor
The hardcoded `"mailto:info@nutritrack.ir"` was in the app layer; it is now a constructor parameter of `WebpushSender` in `wire.go`.

## Pre-Fixed Findings (verified)

- **PushSubscription factory**: `Subscribe` already calls `entity.NewPushSubscription(userID, endpoint, p256dh, auth)`. ✓

## Deferred Items

- **[LOW]** Per-delivery error logging in `Send` — `_ = s.sender.SendToSubscription(...)` silences errors; logging deferred.

## Final Build Status
PASS — `go build ./...` after all fixes
PASS — `go vet ./internal/...` after all fixes

## Remaining Violations
None at HIGH severity.
