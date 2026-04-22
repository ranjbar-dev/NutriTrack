# DDD Audit: internal/application/user
Layer: application
Audited: 2026-04-22
Files reviewed: 3 (avatar_service.go, client_service.go, nutritionist_service.go)

## Summary
- CRITICAL: 2
- HIGH: 3
- MEDIUM: 1
- LOW: 2
- PASS: 0

---

## Findings

### [CRITICAL] Application service imports concrete infrastructure package

**File:** `avatar_service.go:11`
**Issue:** `avatar_service.go` imports `github.com/ranjbar-dev/nutritrack/internal/infrastructure/storage` directly. The application layer must never depend on infrastructure — only on domain interfaces.
**DDD Rule:** `internal/application/` MUST NOT import `internal/infrastructure`.
**Fix:** Define a `FileStorage` port interface in the domain layer; inject through the interface. `*storage.LocalStorage` already satisfies this interface — only the import site and field type change.

---

### [CRITICAL] AvatarService field holds concrete infrastructure type

**File:** `avatar_service.go:17`
**Issue:** `storage *storage.LocalStorage` pins the service to a concrete struct from the infrastructure layer.
**DDD Rule:** Application services MUST depend on interfaces (ports), not concrete implementations.
**Fix:** Define `FileStorage` interface and change field type to `FileStorage`.

---

### [HIGH] `entity.User` aggregate constructed via raw struct literal (no factory function)

**File:** `client_service.go:65`, `nutritionist_service.go:70`
**Issue:** Both `RegisterClient` and `Create` build `entity.User` with a plain struct literal, bypassing validation.
**DDD Rule:** Factory functions `New*()` MUST validate required fields and return `(T, error)`. Application code must call the factory.
**Fix:** Add `NewClient(...)` and `NewNutritionist(...)` factory functions to the domain entity. Replace struct literals with factory calls.

---

### [HIGH] Direct mutation of exported aggregate fields (no domain setter methods)

**File:** `avatar_service.go:78`, `client_service.go:122–137`, `client_service.go:162`, `nutritionist_service.go:125–131`, `nutritionist_service.go:160`
**Issue:** Application code writes directly to exported `entity.User` fields (`FirstName`, `LastName`, `Gender`, `IsActive`, `AvatarURL`, `Height`, `Weight`, `BirthDate`), bypassing invariant enforcement.
**DDD Rule:** Aggregates MUST NOT expose raw entity fields — mutable state MUST go through exported methods.
**Fix:** Add domain mutation methods: `UpdateProfile(...)`, `SetAvatarURL(url string)`, `Activate()`, `Deactivate()`. Replace all direct field assignments with method calls.

---

### [MEDIUM] Cross-application-service import for password hashing

**File:** `nutritionist_service.go:8`
**Issue:** `NutritionistService` imports `internal/application/auth` solely to call `appAuth.HashPassword`. This creates horizontal coupling between application services.
**Fix:** Extract `HashPassword`/`CheckPassword` to `internal/security/password.go`. Both auth and user services import from `internal/security`.

---

### [LOW] Missing blank line before `ListClients` function declaration

**File:** `client_service.go:168`
**Fix:** Insert one blank line between `SetClientStatus` and `ListClients`.

---

### [LOW] Missing blank line before `SetStatus` function declaration

**File:** `nutritionist_service.go:152`
**Fix:** Insert one blank line between `GetClients` and `SetStatus`.

---

## Compliant Patterns Found

- **Repository interfaces injected at construction** — All three services accept `userRepo.UserRepository` (domain interface). ✓
- **Factory constructors exist** — `NewAvatarService`, `NewClientService`, `NewNutritionistService` all exist. ✓
- **No SQL, no DB connections** — Zero SQL or database/sql imports. ✓
- **No HTTP types** — No `http.Request`, `gin.Context` in any service. ✓
- **Value objects used for validation** — `valueobject.NewMobile(req.Mobile)` called before construction. ✓
- **Domain errors used consistently** — `shared.ErrInternal`, `shared.ErrUserNotFound`, etc. ✓

## Fix Priority Order
1. **[CRITICAL]** Define `FileStorage` interface in domain; remove infrastructure import from `avatar_service.go`
2. **[CRITICAL]** Change `storage` field type in `AvatarService` to `FileStorage` interface
3. **[HIGH]** Add `NewClient()`/`NewNutritionist()` factory functions to domain entity
4. **[HIGH]** Add domain mutation methods; replace direct field assignments with method calls
5. **[MEDIUM]** Extract password utilities to `internal/security/`
6. **[LOW]** Fix missing blank lines
