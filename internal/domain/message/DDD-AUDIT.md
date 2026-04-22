# DDD Audit: internal/domain/message
Layer: domain
Audited: 2026-04-22
Files reviewed: 2

## Summary
- CRITICAL: 0
- HIGH: 2
- MEDIUM: 1
- LOW: 1
- PASS: 1 file with no issues (repository/message_repository.go)

---

## Findings

### [HIGH] Message entity exposes all fields as exported

**File:** `entity/message.go:10-20`
**Issue:** Every field on the `Message` struct is exported. DDD entities must have unexported fields; external access must go through explicit getter/setter methods.
**DDD Rule:** Entities: "unexported fields, mutable state via exported methods"
**Fix:** Lowercase all fields and add getter methods. For controlled mutation expose a domain method: `MarkRead(at time.Time)`.

---

### [HIGH] No `NewMessage()` factory function

**File:** `entity/message.go` (absent)
**Issue:** There is no `NewMessage(...)` constructor. Without a factory, callers can create zero-value `Message{}` with no ID, empty sender/receiver, or empty content.
**DDD Rule:** Factory `New*()` function that validates inputs and returns an error.
**Fix:** Add factory validating `senderID`, `receiverID`, and `content`. Domain errors: `ErrEmptyContent`, `ErrMissingSender`, `ErrMissingReceiver`.

---

### [MEDIUM] No getter methods — direct field access bypasses domain encapsulation

**File:** `entity/message.go:10-20`
**Issue:** Because all fields are exported, consumers across every layer read and write fields directly. The domain loses the ability to enforce any future invariant.
**DDD Rule:** Aggregates MUST NOT expose raw entity fields — use getter methods.
**Fix:** Implement getter methods (addressed by the HIGH fix above).

---

### [LOW] No domain error variables defined

**File:** `entity/message.go` (absent)
**Issue:** No sentinel error values. Repository callers have no way to pattern-match on domain-specific errors.
**Fix:** Add `ErrMessageNotFound`, `ErrUnauthorizedRead`.

---

## Compliant Patterns Found
- **`repository/message_repository.go`** — Fully compliant; pure Go `interface`. ✓
- **No `json:`, `bson:`, `db:` struct tags** on `Message`. ✓
- **No forbidden cross-layer imports**. ✓
- **`HasAttachment()`** — Well-placed domain behaviour method. ✓

## Fix Priority Order
1. Add `NewMessage()` factory with input validation and domain error variables
2. Make all `Message` fields unexported and expose getter methods
3. Add domain-level `MarkRead(at time.Time)` method
4. Add sentinel error variables
