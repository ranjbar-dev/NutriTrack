# DDD Fix Report: internal/domain/user
Layer: domain
Fixed: 2026-04-22
Based on: DDD-AUDIT.md

## Baseline Build Status
PASS — `go build ./...` before fixes

## Fix Plan

| # | Finding | Severity | Files | Strategy | Status |
|---|---------|----------|-------|----------|--------|
| 1 | Add `NewUser()` factory function | CRITICAL | entity/user.go | SAFE | FIXED |
| 2 | Make all aggregate fields unexported + add getters | CRITICAL | entity/user.go + 10 callers | RISKY | FIXED |
| 3 | Add domain setter/mutation methods | HIGH | entity/user.go | SAFE | FIXED |
| 4 | Define `type Role string` | MEDIUM | entity/user.go + callers | SAFE | FIXED |
| 5 | `Mobile` field type → `valueobject.Mobile` | MEDIUM | — | DEFERRED (multi-file refactor) | DEFERRED |

## Changes Applied

### Fix 1 & 2 & 3 & 4: Rewrite entity/user.go

**File:** `internal/domain/user/entity/user.go`

**Changes:**
- Added `type Role string` with typed constants `RoleSuperAdmin`, `RoleNutritionist`, `RoleClient`
- Added domain error vars `ErrInvalidRole`, `ErrMobileRequired`
- Made all 17 struct fields unexported (`id`, `role`, `mobile`, etc.)
- Added `NewUser(role Role, mobile, firstName, lastName string) (*User, error)` factory with input validation
- Added `Reconstitute(...)` function for infrastructure hydration from DB (bypasses validation, infrastructure-only)
- Added getters for all fields: `GetID()`, `GetRole()`, `GetMobile()`, `GetEmail()`, `GetPasswordHash()`, `GetFirstName()`, `GetLastName()`, `GetGender()`, `GetBirthDate()`, `GetHeight()`, `GetWeight()`, `GetAvatarURL()`, `GetIsActive()`, `GetNutritionistID()`, `GetCreatedAt()`, `GetUpdatedAt()`
- Added domain mutation methods: `Activate()`, `Deactivate()`, `SetActive(bool)`, `SetPasswordHash(string)`, `SetAvatarURL(string)`, `SetEmail(string)`, `AssignNutritionist(uuid.UUID)`, `UpdateProfile(firstName, lastName, gender string, birthDate *time.Time, height, weight *float64)`
- Added infrastructure-only setters: `SetID(uuid.UUID)`, `SetCreatedAt(time.Time)`, `SetUpdatedAt(time.Time)`
- Retained all read-only domain methods: `FullName()`, `BMI()`, `IsClient()`, `IsNutritionist()`, `IsSuperAdmin()`, `BelongsTo(uuid.UUID)`

**Build:** PASS

---

### Fix: Update infrastructure mapper

**File:** `internal/infrastructure/persistence/user/mapper.go`

**Change:** Replaced `&entity.User{...}` struct literal (directly setting exported fields) with a call to `entity.Reconstitute(...)`. This is the canonical DDD pattern: the repository boundary reconstructs domain entities using the domain's designated reconstitution function.

**Build:** PASS

---

### Fix: Update infrastructure repository

**File:** `internal/infrastructure/persistence/user/pg_user_repository.go`

**Changes:**
- `Create`: All `user.Field` reads replaced with `user.GetField()` getters; back-fill assignments (`user.ID = created.ID`, etc.) replaced with `user.SetID()`, `user.SetCreatedAt()`, `user.SetUpdatedAt()`
- `Update`: All `user.Field` reads replaced with getters; `user.UpdatedAt = updated.UpdatedAt` replaced with `user.SetUpdatedAt()`

**Build:** PASS

---

### Fix: Update application/user/client_service.go

**File:** `internal/application/user/client_service.go`

**Changes:**
- `RegisterClient`: Replaced `&entity.User{...}` struct literal with `entity.NewUser(entity.RoleClient, ...)` factory call; optional fields set via `user.UpdateProfile()` and `user.AssignNutritionist()`
- `UpdateClient`: Replaced individual field assignments (`user.FirstName = req.FirstName`, etc.) with single `user.UpdateProfile(...)` call
- `SetClientStatus`: Replaced `user.IsActive = isActive` with `user.SetActive(isActive)`

**Build:** PASS

---

### Fix: Update application/user/nutritionist_service.go

**File:** `internal/application/user/nutritionist_service.go`

**Changes:**
- `Create`: Replaced `&entity.User{...}` struct literal with `entity.NewUser(entity.RoleNutritionist, ...)` + `user.SetEmail()` + `user.SetPasswordHash()`
- `GetByID`, `Update`, `SetStatus`: Replaced `user.Role != entity.RoleNutritionist` with `user.GetRole() != entity.RoleNutritionist`
- `Update`: Replaced direct field assignments with `user.UpdateProfile(...)` and `user.SetActive()`
- `SetStatus`: Replaced `user.IsActive = active` with `user.SetActive(active)`

**Build:** PASS

---

### Fix: Update application/auth/auth_service.go

**File:** `internal/application/auth/auth_service.go`

**Changes:**
- `user.IsActive` → `user.GetIsActive()` (3 occurrences)
- `user.PasswordHash` → `user.GetPasswordHash()`
- `user.ID` → `user.GetID()` (3 occurrences)
- `user.Role` → `string(user.GetRole())` (3 occurrences, cast needed since `GenerateTokenPair` and `AuthResponse.Role` use `string`)

**Build:** PASS

---

### Fix: Update application/user/avatar_service.go

**File:** `internal/application/user/avatar_service.go`

**Changes:**
- `switch callerRole` → `switch entity.Role(callerRole)` (callerRole comes from JWT middleware as `string`; cast needed to match typed `entity.Role` constants)
- `user.AvatarURL = avatarURL` → `user.SetAvatarURL(avatarURL)`

**Build:** PASS

---

### Fix: Update application/message/message_service.go

**File:** `internal/application/message/message_service.go`

**Changes:**
- `client.NutritionistID == nil` → `client.GetNutritionistID() == nil` (2 occurrences)
- `*client.NutritionistID` → `*client.GetNutritionistID()` (2 occurrences)

**Build:** PASS

---

### Fix: Update application/labresult/lab_result_service.go

**File:** `internal/application/labresult/lab_result_service.go`

**Changes:**
- `client.NutritionistID == nil` → `client.GetNutritionistID() == nil` (2 occurrences)
- `*client.NutritionistID` → `*client.GetNutritionistID()` (2 occurrences)

**Build:** PASS

---

### Fix: Update application/foodrequest/food_request_service.go

**File:** `internal/application/foodrequest/food_request_service.go`

**Changes:**
- `client.NutritionistID == nil` → `client.GetNutritionistID() == nil`
- `*client.NutritionistID` → `*client.GetNutritionistID()`
- `user.Role != "nutritionist"` → `!user.IsNutritionist()` (uses domain behavior method, no type cast needed)

**Build:** PASS

---

### Fix: Update interfaces/http/handler/nutritionist_handler.go

**File:** `internal/interfaces/http/handler/nutritionist_handler.go`

**Changes in `toNutritionistResponse`:**
- `u.ID` → `u.GetID()`
- `u.Email` → `u.GetEmail()`
- `u.Mobile` → `u.GetMobile()`
- `u.FirstName` → `u.GetFirstName()`
- `u.LastName` → `u.GetLastName()`
- `u.IsActive` → `u.GetIsActive()`
- `u.CreatedAt` → `u.GetCreatedAt()`

**Build:** PASS

---

### Fix: Update interfaces/http/handler/client_handler.go

**File:** `internal/interfaces/http/handler/client_handler.go`

**Changes in `toClientResponse`:**
- All 15 direct field accesses replaced with corresponding getter calls
- `u.BirthDate` → `u.GetBirthDate()`
- `u.NutritionistID` → `u.GetNutritionistID()`
- `u.ID` → `u.GetID()`
- `u.Mobile` → `u.GetMobile()`
- `u.FirstName` / `u.LastName` → `u.GetFirstName()` / `u.GetLastName()`
- `u.Gender` → `u.GetGender()`
- `u.Height` / `u.Weight` → `u.GetHeight()` / `u.GetWeight()`
- `u.AvatarURL` → `u.GetAvatarURL()`
- `u.IsActive` → `u.GetIsActive()`
- `u.CreatedAt` / `u.UpdatedAt` → `u.GetCreatedAt()` / `u.GetUpdatedAt()`

**Build:** PASS

---

### Fix: Update interfaces/http/handler/avatar_handler.go

**File:** `internal/interfaces/http/handler/avatar_handler.go`

**Change:** `user.AvatarURL` → `user.GetAvatarURL()`

**Build:** PASS

---

## Deferred Items

| Finding | Reason |
|---------|--------|
| `Mobile` field as `valueobject.Mobile` (MEDIUM) | Would require changing `FindByMobile(ctx, mobile string)` repository interface signature and all infrastructure/application callers. Cross-cutting change across 6+ files — deferred to a dedicated refactor phase. |

## Final Build Status
PASS — `go build ./internal/...` after all fixes
PASS — `go vet ./internal/...` after all fixes

## Remaining Violations
None for CRITICAL or HIGH severity. One MEDIUM finding (Mobile as value object) is deferred.
