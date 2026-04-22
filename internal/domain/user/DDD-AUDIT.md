# DDD Audit: internal/domain/user
Layer: domain
Audited: 2026-04-22
Files reviewed: 4

## Summary
- CRITICAL: 2
- HIGH: 1
- MEDIUM: 2
- LOW: 1
- PASS: 3 files with no issues (`repository/user_repository.go`, `repository/otp_repository.go`, `valueobject/mobile.go`)

---

## Findings

### [CRITICAL] Aggregate fields are all exported (public)

**File:** `entity/user.go:23-39`
**Issue:** Every field on the `User` aggregate (`ID`, `Role`, `Mobile`, `Email`, `PasswordHash`, `FirstName`, `LastName`, `Gender`, `BirthDate`, `Height`, `Weight`, `AvatarURL`, `IsActive`, `NutritionistID`, `CreatedAt`, `UpdatedAt`) is exported. Any caller in any package can read or overwrite any field directly, completely bypassing domain invariants and encapsulation.
**DDD Rule:** Aggregates MUST have unexported fields; state is exposed only through getter methods and mutated only through domain-behavior methods.
**Fix:** Lowercase all fields and add getters/setters.

---

### [CRITICAL] No `NewUser()` factory function

**File:** `entity/user.go` (file-level — function absent)
**Issue:** There is no `NewUser(...)` constructor. Every consumer must manually assemble the struct with exported field assignments, meaning required fields can be omitted, invalid roles can be set, and invariants are never enforced at creation time.
**DDD Rule:** Aggregates MUST have a `New*()` factory function that validates inputs and returns `(*T, error)`.

---

### [HIGH] No domain setter/mutation methods — state changes are uncontrolled

**File:** `entity/user.go:41-76`
**Issue:** The only methods on `User` are read-only derivations. There are no methods to mutate state — callers must write directly to exported fields.
**DDD Rule:** Aggregates MUST expose getter/setter methods for controlled access; mutation should only happen through domain-behavior methods that enforce invariants.
**Fix:** Add `Activate()`, `Deactivate()`, `AssignNutritionist()`, `UpdateProfile()`, `SetPasswordHash()`, `SetAvatarURL()` methods.

---

### [MEDIUM] `Mobile` field is primitive `string` instead of `valueobject.Mobile`

**File:** `entity/user.go:26`
**Issue:** `Mobile string` accepts any string. The domain package already defines a validated `Mobile` value object but the aggregate does not use it.
**DDD Rule:** Value Objects MUST be used inside Aggregates to carry validated domain concepts.

---

### [MEDIUM] `Role` is an untyped `string` — no type-system enforcement

**File:** `entity/user.go:10-14`
**Issue:** Role constants are defined but the `Role` field is `string`, allowing any arbitrary string assignment.
**DDD Rule:** Domain concepts with a closed value set should be expressed as typed constants.
**Fix:** Define `type Role string` and use it for the field.

---

### [LOW] `PasswordHash` is an exported field — sensitive data with no access guard

**File:** `entity/user.go:28`
**Issue:** The `PasswordHash` field is exported, allowing any package to read or overwrite the raw bcrypt hash.
**Fix:** Make the field unexported and add controlled getter/setter.

---

## Compliant Patterns Found

1. **`valueobject/mobile.go`** — Textbook value object: single unexported `value` field, no setters, `NewMobile()` factory validates format. ✓
2. **`repository/user_repository.go`** — Correct interface-only definition in the domain layer. ✓
3. **`repository/otp_repository.go`** — Correct interface-only definition. ✓
4. **No struct tags on the aggregate** — `entity/user.go` has zero `json:`, `bson:`, or `db:` tags. ✓
5. **No forbidden imports** — None of the files import `internal/infrastructure`, `internal/interfaces`, or `internal/application`. ✓

## Fix Priority Order
1. **[CRITICAL]** Add `NewUser()` factory
2. **[CRITICAL]** Make all aggregate fields unexported + add getters
3. **[HIGH]** Add domain setter/mutation methods
4. **[MEDIUM]** Change `Mobile` field type to `valueobject.Mobile`
5. **[MEDIUM]** Define `type Role string`
6. **[LOW]** Ensure `passwordHash` stays unexported
