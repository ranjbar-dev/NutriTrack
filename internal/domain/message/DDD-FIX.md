# DDD Fix Report: internal/domain/message
Layer: domain
Fixed: 2025-07
Based on: DDD-AUDIT.md

## Baseline Build Status
PASS — `go build ./...` before fixes

## Fix Plan

| # | Finding | Severity | Files | Strategy | Status |
|---|---------|----------|-------|----------|--------|
| 1 | Aggregate has exported fields | HIGH | entity/message.go | SAFE | FIXED |
| 2 | No `NewMessage()` factory | HIGH | entity/message.go | SAFE | FIXED |
| 3 | Callers use direct field access | HIGH | persistence/message/mapper.go, pg_message_repository.go, application/message/message_service.go, interfaces/http/handler/message_handler.go | SAFE | FIXED |

## Changes Applied

### Fix 1 + 2: Unexported fields, getters, setters, factory functions
**File:** `internal/domain/message/entity/message.go`
**Change:** All 10 fields made unexported (id, senderID, receiverID, senderRole, content, isRead, isDeleted, threadID, replyToID, createdAt). Added getter for every field. Added `NewMessage(senderID, receiverID uuid.UUID, senderRole, content string, threadID, replyToID *uuid.UUID) *Message` factory. Added `ReconstituteMessage(...)` for DB loading.
**Build:** PASS

### Fix 3a: Mapper updated
**File:** `internal/infrastructure/persistence/message/mapper.go`
**Change:** Uses `entity.ReconstituteMessage(...)` instead of struct literal.
**Build:** PASS

### Fix 3b: Repository updated
**File:** `internal/infrastructure/persistence/message/pg_message_repository.go`
**Change:** All DB params use getter methods; SetID/SetCreatedAt used after insert.
**Build:** PASS

### Fix 3c: Service updated
**File:** `internal/application/message/message_service.go`
**Change:** `entity.Message{...}` struct literal replaced with `entity.NewMessage(...)`.
**Build:** PASS

### Fix 3d: Handler updated
**File:** `internal/interfaces/http/handler/message_handler.go`
**Change:** All `msg.X` and inline `result.ReceiverID` field accesses in both SendAsClient and SendAsNutritionist handlers replaced with getter calls (`msg.X()`, `result.ReceiverID()`).
**Build:** PASS

## Final Build Status
PASS — `go build ./...` after all fixes
PASS — `go vet ./internal/...` after all fixes

## Remaining Violations
None — all CRITICAL and HIGH findings resolved.
