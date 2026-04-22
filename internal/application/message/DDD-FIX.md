# DDD Fix Report: internal/application/message
Layer: application
Fixed: 2026-04-22
Based on: DDD-AUDIT.md

## Baseline Build Status
PASS — `go build ./...` before fixes

## Fix Plan

| # | Finding | Severity | Files | Strategy | Status |
|---|---------|----------|-------|----------|--------|
| 1 | Application service imports concrete infrastructure (storage) package | CRITICAL | message_service.go | SAFE | FIXED |
| 2 | Domain entity constructed as raw struct literal; fields mutated directly | HIGH | message_service.go | DEFERRED | DEFERRED: Message entity still has exported fields; adding NewMessage() factory requires unexported-field refactor across all callers |
| 3 | Domain invariant (self-send check) enforced in application layer | MEDIUM | message_service.go | DEFERRED | DEFERRED: depends on Fix 2 (NewMessage factory) |

## Changes Applied

### Fix 1: Replace `*storage.LocalStorage` with `shared.AttachmentStorage` interface
**File:** `internal/domain/shared/storage.go` (new)
**Change:** Created `AttachmentStorage` port interface in the shared domain package:
```go
type AttachmentStorage interface {
    SaveAttachment(src io.Reader, ext string) (string, error)
}
```
`*storage.LocalStorage` satisfies this interface implicitly (structural typing).
**Build:** PASS

**File:** `internal/application/message/message_service.go`
**Before:**
```go
import "github.com/ranjbar-dev/nutritrack/internal/infrastructure/storage"

type MessageService struct {
    // ...
    storage  *storage.LocalStorage
}

func NewMessageService(
    msgRepo msgRepo.MessageRepository,
    userRepo userRepo.UserRepository,
    storage *storage.LocalStorage,
) *MessageService { ... }
```
**After:**
```go
// infrastructure/storage import removed

type MessageService struct {
    // ...
    storage  shared.AttachmentStorage
}

func NewMessageService(
    msgRepo msgRepo.MessageRepository,
    userRepo userRepo.UserRepository,
    storage shared.AttachmentStorage,
) *MessageService { ... }
```
**Build:** PASS

No changes required to `bootstrap/wire.go` — `*storage.LocalStorage` satisfies `shared.AttachmentStorage` implicitly.

## Deferred Items
- **[HIGH]** Domain entity `Message` constructed via raw struct literal in `SendAsClient` and `SendAsNutritionist`. Fixing this requires making `Message` fields unexported and adding a `NewMessage()` factory, which would require updating all callers across interfaces and infrastructure layers. Deferred to a dedicated domain entity refactor.
- **[MEDIUM]** Self-send invariant check lives in application layer; should move into `NewMessage()`. Deferred pending the entity refactor above.

## Final Build Status
PASS — `go build ./...` after all fixes
PASS — `go vet ./internal/...` after all fixes

## Remaining Violations
- HIGH: Message struct literal construction — deferred (multi-file entity refactor required)
- MEDIUM: Self-send invariant in application layer — deferred (depends on entity refactor)
