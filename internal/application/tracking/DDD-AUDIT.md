# DDD Audit: internal/application/tracking
Layer: application
Audited: 2026-04-22
Files reviewed: 1 (tracking_service.go)

## Summary
- CRITICAL: 0
- HIGH: 2
- MEDIUM: 4
- LOW: 1
- PASS: 0

---

## Findings

### [HIGH] Direct entity struct literal construction bypasses domain factory pattern

**File:** `tracking_service.go:138, 161, 177, 195, 213, 231` and BulkSync at `:266, 290, 305, 323, 340, 356`
**Issue:** All six `Log*` methods and every `switch` branch of `BulkSync` construct domain entities using raw struct literals with exported fields. Invalid state (zero quantity, negative calories, empty names) is not rejected.
**DDD Rule:** Application services must create entities exclusively via `New*()` factory functions. Direct struct literal construction bypasses all domain invariants.
**Fix:** Add factory functions with unexported fields to `internal/domain/tracking/entity/tracking.go` for all six entity types (`NewFoodLog`, `NewWaterLog`, `NewSleepLog`, `NewExerciseLog`, `NewMedicationLog`, `NewBodyMeasurement`). Replace all struct literals with factory calls.

---

### [HIGH] No input validation in `Log*` service methods

**File:** `tracking_service.go:138, 161, 177, 195, 213, 231`
**Issue:** None of the six `Log*` methods validate request inputs before constructing and persisting entities. Missing guards: `Quantity > 0`, `Calories >= 0`, `FoodName != ""`, `AmountMl > 0`, `Quality` in valid range, `DurationMinutes > 0`, etc.
**DDD Rule:** Application services are the entry point for domain commands and must validate all inputs; preferred: enforce invariants inside domain factory `New*()` functions.
**Fix:** Resolve HIGH-1 first (domain factories with validation). As interim, add guard clauses at the start of each `Log*` method.

---

### [MEDIUM] `json:` struct tags on application-layer types

**File:** `tracking_service.go:89` (`SyncEntry`), `:96` (`BulkSyncResult`)
**Issue:** Both types carry `json:` serialization tags; `SyncEntry.Data` is `json.RawMessage`. These are HTTP/JSON wire-format DTOs that belong in the interfaces layer.
**DDD Rule:** Application layer MUST NOT define types with `json:` serialization tags.
**Fix:** Move `SyncEntry` and `BulkSyncResult` to `internal/interfaces/http/tracking/`. Replace with domain-oriented `BulkSyncCommand` and `BulkSyncSummary` (no json tags) in the application layer.

---

### [MEDIUM] `encoding/json` import and JSON parsing inside application service

**File:** `tracking_service.go:5, 262, 284, 299, 316, 333, 350`
**Issue:** The application service imports `"encoding/json"` and calls `json.Unmarshal` six times inside `BulkSync`. JSON deserialization is a transport concern belonging in the interfaces layer.
**DDD Rule:** Application layer MUST NOT perform JSON parsing.
**Fix:** Move all `json.Unmarshal` calls into the HTTP handler. Handler builds typed `BulkSyncCommand` and passes it to the service.

---

### [MEDIUM] `GetTracking` returns `any` — loss of compile-time type safety

**File:** `tracking_service.go:387`
**Issue:** `GetTracking(...) (any, error)` erases all type information. Every caller must perform unsafe type assertions.
**DDD Rule:** Application services must return domain-typed results.
**Fix:** Replace with six typed methods: `GetFoodLogs`, `GetWaterLogs`, `GetSleepLogs`, `GetExerciseLogs`, `GetMedicationLogs`, `GetBodyMeasurements`.

---

### [MEDIUM] `BulkSync` duplicates entity construction from `Log*` methods

**File:** `tracking_service.go:252–386`
**Issue:** `BulkSync` re-implements entity construction for all six entry types — near-exact copies of `Log*` methods. Any future change must be mirrored manually.
**DDD Rule:** Application services must not duplicate business logic.
**Fix:** Refactor each branch of `BulkSync` to delegate to the corresponding `Log*` method.

---

### [LOW] Silent error swallowing in `BulkSync` — no per-entry failure reporting

**File:** `tracking_service.go:372`
**Issue:** Repository errors are silently discarded (`continue` with no logging). Callers cannot distinguish a skipped duplicate from a failed write.
**Fix:** Add `Failures []SyncFailure` to the result type; track per-entry failures.

---

## Compliant Patterns Found

- **Repository interfaces injected** — accepts `trackRepo.TrackingRepository` and `userRepo.UserRepository` (domain interfaces). ✓
- **`NewTrackingService` factory present**. ✓
- **No forbidden layer imports** — No `internal/infrastructure` or `internal/interfaces` imports. ✓
- **Domain errors used correctly** — `shared.ErrForbidden`, `shared.ErrUserNotFound`, `shared.ErrValidation`. ✓
- **Application command objects are tag-free** — `LogFoodRequest`, etc. carry no `json:` tags. ✓
- **Access control delegates to domain** — `checkClientAccess` calls `client.BelongsTo(callerID)` (domain entity method). ✓

## Fix Priority Order
1. **[HIGH]** Add `New*()` factory functions with unexported fields to domain entity package for all six types
2. **[HIGH]** Add input validation (via domain factories preferred, or guard clauses as interim)
3. **[MEDIUM]** Move `SyncEntry`/`BulkSyncResult` to interfaces layer; remove `encoding/json` import
4. **[MEDIUM]** Refactor `BulkSync` to delegate to `Log*` methods
5. **[MEDIUM]** Replace `GetTracking(any)` with six typed getter methods
6. **[LOW]** Add per-entry failure tracking to bulk sync result
