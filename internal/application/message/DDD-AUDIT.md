# DDD Audit: internal/application/message
Layer: application
Audited: 2026-04-22
Files reviewed: 1 (message_service.go)

## Summary
- CRITICAL: 1
- HIGH: 1
- MEDIUM: 1
- LOW: 0
- PASS: 0

---

## Findings

### [CRITICAL] Application service imports concrete infrastructure (storage) package

**File:** `message_service.go` (import block)
**Issue:** The application service directly imports `internal/infrastructure/storage` (or a concrete storage adapter). Application services MUST depend on domain port interfaces only.
**DDD Rule:** Application layer → Domain ports only. NEVER import infrastructure packages directly.
**Fix:** Define a `FileStorage` port interface in the domain layer (e.g. `internal/domain/shared/port/`), have the storage adapter implement it, and inject through the interface.

---

### [HIGH] Domain entity constructed as raw struct literal; fields mutated directly

**File:** `message_service.go` (SendMessage / CreateMessage)
**Issue:** The `Message` entity is built as `messageentity.Message{SenderID: ..., ReceiverID: ..., ...}` via a struct literal, bypassing factory validation. Fields are also mutated post-construction.
**DDD Rule:** Aggregates MUST be constructed through `New*()` factories. Direct field assignment is forbidden once fields are unexported.
**Fix:** Call `messageentity.NewMessage(senderID, receiverID, content)` and use aggregate methods for state transitions.

---

### [MEDIUM] Domain invariant enforced in application layer

**File:** `message_service.go`
**Issue:** The check "cannot send message to yourself" is performed in the application service. This is a domain invariant and belongs inside the `NewMessage()` factory or a domain method.
**DDD Rule:** Domain invariants MUST be enforced within the domain aggregate/factory.
**Fix:** Move the self-send check inside `NewMessage()` — return `ErrSelfMessage` if senderID == receiverID.

---

## Compliant Patterns Found

- Service constructor accepts repository interface (not concrete type). ✓
- Pagination and listing logic is appropriately placed in the application service. ✓

## Fix Priority Order
1. **[CRITICAL]** Define storage domain port interface; remove direct infrastructure import
2. **[HIGH]** Replace struct literal with `NewMessage()` factory call
3. **[MEDIUM]** Move self-send invariant into `NewMessage()` domain factory
