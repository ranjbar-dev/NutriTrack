# DDD Fix Report: internal/infrastructure/persistence
Layer: infrastructure
Fixed: 2026-04-22
Based on: DDD-AUDIT.md

## Baseline Build Status
PASS — `go build ./...` before fixes

## Fix Plan

| # | Finding | Severity | Files | Strategy | Status |
|---|---------|----------|-------|----------|--------|
| 1 | Raw DB errors leak through domain boundary — foodrequest | HIGH | foodrequest/pg_food_request_repository.go | SAFE | FIXED |
| 2 | Raw DB errors leak through domain boundary — labresult | HIGH | labresult/pg_lab_result_repository.go | SAFE | FIXED |
| 3 | Raw DB errors leak through domain boundary — message | HIGH | message/pg_message_repository.go | SAFE | FIXED |
| 4 | Raw DB errors leak through domain boundary — notification | HIGH | notification/pg_notification_preference_repository.go | SAFE | FIXED |
| 5 | Raw DB errors leak through domain boundary — push | HIGH | push/pg_push_subscription_repository.go | SAFE | FIXED |
| 6 | `err == pgx.ErrNoRows` instead of `errors.Is` — foodrequest | LOW | foodrequest/mapper.go | SAFE | FIXED |
| 7 | `err == pgx.ErrNoRows` instead of `errors.Is` — labresult | LOW | labresult/mapper.go | SAFE | FIXED |
| 8 | `err == pgx.ErrNoRows` instead of `errors.Is` — message | LOW | message/mapper.go | SAFE | FIXED |

## Changes Applied

### Fix 1: foodrequest — add `shared` import; wrap all raw errors

**File:** `internal/infrastructure/persistence/foodrequest/pg_food_request_repository.go`

**Changes:**
- Added `"github.com/ranjbar-dev/nutritrack/internal/domain/shared"` to imports.
- `Create`: `return err` → `return shared.ErrInternal`
- `FindByID` (non-not-found path): `return nil, err` → `return nil, shared.ErrInternal`
- `ListPending`: `return nil, err` → `return nil, shared.ErrInternal`
- `CountPending`: split direct pass-through into `count, err := …; if err != nil { return 0, shared.ErrInternal }; return count, nil`
- `UpdateStatus`: `return nil, err` → `return nil, shared.ErrInternal`

**Build:** PASS

---

### Fix 2: labresult — wrap all raw errors

**File:** `internal/infrastructure/persistence/labresult/pg_lab_result_repository.go`

**Note:** `shared` import was already present.

**Changes:**
- `Create`: `return nil, err` → `return nil, shared.ErrInternal`
- `FindByID` (non-not-found path): `return nil, err` → `return nil, shared.ErrInternal`  
  _(The `isNotFound` path that returns `shared.ErrLabResultNotFound` was already correct — left unchanged.)_
- `ListByClientID` (ListLabResultsByClientID query error): `return nil, 0, err` → `return nil, 0, shared.ErrInternal`
- `ListByClientID` (CountLabResultsByClientID query error): `return nil, 0, err` → `return nil, 0, shared.ErrInternal`

**Build:** PASS

---

### Fix 3: message — add `shared` import; wrap all raw errors

**File:** `internal/infrastructure/persistence/message/pg_message_repository.go`

**Changes:**
- Added `"github.com/ranjbar-dev/nutritrack/internal/domain/shared"` to imports.
- `Create`: `return nil, err` → `return nil, shared.ErrInternal`
- `ListConversation` (rows query error): `return nil, 0, err` → `return nil, 0, shared.ErrInternal`
- `ListConversation` (count query error): `return nil, 0, err` → `return nil, 0, shared.ErrInternal`
- `MarkRead`: split `return r.queries.MarkConversationRead(…)` into `if err := …; err != nil { return shared.ErrInternal }; return nil`
- `CountUnread`: split `return r.queries.CountUnreadMessages(…)` into `count, err := …; if err != nil { return 0, shared.ErrInternal }; return count, nil`

**Build:** PASS

---

### Fix 4: notification — add `shared` import; wrap all raw errors

**File:** `internal/infrastructure/persistence/notification/pg_notification_preference_repository.go`

**Changes:**
- Added `"github.com/ranjbar-dev/nutritrack/internal/domain/shared"` to imports.
- `Upsert`: `return entity.NotificationPreference{}, err` → `return entity.NotificationPreference{}, shared.ErrInternal`
- `GetByUserID`: `return entity.NotificationPreference{}, err` → `return entity.NotificationPreference{}, shared.ErrInternal`

**Build:** PASS

---

### Fix 5: push — add `shared` import; wrap all raw errors

**File:** `internal/infrastructure/persistence/push/pg_push_subscription_repository.go`

**Changes:**
- Added `"github.com/ranjbar-dev/nutritrack/internal/domain/shared"` to imports.
- `Upsert`: `return nil, err` → `return nil, shared.ErrInternal`
- `Delete`: split `return r.q.DeletePushSubscription(…)` into `if err := …; err != nil { return shared.ErrInternal }; return nil`
- `ListByUser`: `return nil, err` → `return nil, shared.ErrInternal`

**Build:** PASS

---

### Fix 6–8: Replace `err == pgx.ErrNoRows` with `errors.Is(err, pgx.ErrNoRows)` in 3 mapper files

**Files:**
- `internal/infrastructure/persistence/foodrequest/mapper.go`
- `internal/infrastructure/persistence/labresult/mapper.go`
- `internal/infrastructure/persistence/message/mapper.go`

**Changes per file:**
- Added `"errors"` to import block.
- `isNotFound`: `return err == pgx.ErrNoRows` → `return errors.Is(err, pgx.ErrNoRows)`

**Build:** PASS

---

## Deferred Items

| Finding | Reason |
|---------|--------|
| [MEDIUM] Constructor return types should be domain interfaces | Requires updating all DI wire-up call sites across `bootstrap/wire.go` and related files — more than 3 files; deferred to avoid breaking changes. |
| [MEDIUM] Missing `NewFromAggregate()` reverse-mapper functions | Additive change; no DDD violation risk; deferred — out of scope for HIGH fix pass. |
| [MEDIUM] Infrastructure directly mutates exported aggregate fields post-insert | Requires domain aggregate refactoring (unexported fields + setters); cross-layer; deferred. |
| [LOW] No per-package internal DTO struct (sqlc types used directly) | Mitigation already in place per audit; deferred as optional. |

## Final Build Status
PASS — `go build ./internal/...` after all fixes  
PASS — `go vet ./internal/...` after all fixes

## Remaining Violations
None at CRITICAL or HIGH severity. All 5 HIGH sub-package violations resolved.
