# DDD Audit: internal/domain/shared
Layer: domain
Audited: 2026-04-22
Files reviewed: 4

## Summary
- CRITICAL: 2
- HIGH: 2
- MEDIUM: 0
- LOW: 1
- PASS: 0

---

## Findings

### [CRITICAL] `AppError` imports `net/http` and carries HTTP status codes

**File:** `errors/app_error.go`
**Issue:** The `AppError` domain type imports `net/http` and stores `HTTPStatus int`. HTTP status codes are a delivery/interface concern and MUST NOT appear in the domain layer. Any domain type that references HTTP semantics ties business logic to a specific transport protocol.
**DDD Rule:** Domain layer MUST NOT import any delivery/framework packages (`net/http`, gin, etc.).
**Fix:** Remove `net/http` import and `HTTPStatus` field. Callers (HTTP handlers) map domain error codes to HTTP status using a translation map in the interface layer.

---

### [CRITICAL] `json:` struct tags on `AppError` domain type

**File:** `errors/app_error.go`
**Issue:** `AppError` contains `json:` struct tags on its fields. Struct tags couple the domain type to a specific serialization format (JSON) and belong in the interface/infrastructure layer.
**DDD Rule:** Domain types MUST NOT carry serialization (json, xml, bson, db) struct tags.
**Fix:** Remove all `json:` struct tags from `AppError`. The interface layer should map to a response DTO with its own struct tags.

---

### [HIGH] `ToResponse()` HTTP serialization logic in the domain layer

**File:** `errors/app_error.go`
**Issue:** `AppError.ToResponse()` returns a serializable map / response struct. This is presentation logic and belongs exclusively in the interface/handler layer.
**DDD Rule:** Domain objects must not know about presentation format (JSON serialization, HTTP response bodies, etc.).
**Fix:** Delete `ToResponse()` from `AppError`. Create a response DTO in the interface layer that maps from `AppError`.

---

### [HIGH] `ImageInfo` value object has exported (mutable) fields

**File:** `valueobject/image_info.go`
**Issue:** `ImageInfo` struct has fully exported fields (`URL`, `AltText`, `Width`, `Height`). Value objects MUST be immutable — expose only getter methods, no direct field assignment.
**DDD Rule:** Value Objects — "immutable, no setters, only getters"
**Fix:** Lowercase all `ImageInfo` fields and add getter methods + a `NewImageInfo()` constructor.

---

### [LOW] `panic` in `init()` for timezone loading

**File:** (timezone utility file)
**Issue:** `init()` panics if the timezone file is missing. Panics in domain code are hard to test and bypass normal error handling.
**Fix:** Return or log a warning instead of panicking; supply a safe fallback (UTC).

---

## Compliant Patterns Found

- `errors/app_error.go` — Named error type with typed error codes is a valid DDD pattern. ✓
- `valueobject/` — ImageInfo correctly placed in valueobject package. ✓

## Fix Priority Order
1. **[CRITICAL]** Remove `net/http` import and `HTTPStatus` field from `AppError`
2. **[CRITICAL]** Remove all `json:` struct tags from `AppError`
3. **[HIGH]** Delete `ToResponse()` from `AppError`
4. **[HIGH]** Make `ImageInfo` fields unexported; add getters and `NewImageInfo()` constructor
5. **[LOW]** Replace `panic` in `init()` with graceful fallback to UTC
