# DDD Audit: internal/application/foodrequest
Layer: application
Audited: 2026-04-22
Files reviewed: 1 (food_request_service.go)

## Summary
- CRITICAL: 0
- HIGH: 1
- MEDIUM: 2
- LOW: 2
- PASS: 0

---

## Findings

### [HIGH] Concrete sibling application service injected instead of interface

**File:** `food_request_service.go:30`
**Issue:** `FoodRequestService.foodSvc` is typed as `*appFood.FoodService` — a concrete struct from the sibling application package. Prevents independent testing and tightly couples bounded contexts.
**DDD Rule:** Application layer services MUST depend on abstractions (interfaces), not concrete types.
**Fix:** Define a narrow `FoodCreator` interface in this package:
```go
type FoodCreator interface {
    CreateFood(ctx context.Context, req appFood.CreateFoodRequest) (*foodEntity.Food, error)
}
```
Accept it in the constructor. `*appFood.FoodService` already satisfies this interface.

---

### [MEDIUM] Direct struct literal construction of domain entity bypasses factory and invariant protection

**File:** `food_request_service.go:52–57`
**Issue:** `Submit` constructs `frEntity.FoodRequest` as a raw struct literal and hard-codes the initial `Status: frEntity.FoodRequestStatusPending` invariant in the application service.
**DDD Rule:** Factory functions `New*()` MUST encode initial invariants. Application services must call the domain factory.
**Fix:** Add `NewFoodRequest(clientID, nutritionistID uuid.UUID, foodName string) (*FoodRequest, error)` to the domain entity package. Replace struct literal with factory call.

---

### [MEDIUM] Hard-coded role string literals bypass domain-defined constants and entity methods

**File:** `food_request_service.go:69, 108`
**Issue:** `user.Role != "nutritionist"` uses raw string literal. The `User` entity exports `IsNutritionist()` method and `RoleNutritionist` constant; these are not used.
**DDD Rule:** Application services must use domain-defined constants and entity behaviour methods.
**Fix:** Replace `user.Role != "nutritionist"` with `!user.IsNutritionist()`. Replace string literal at line 108 with `userEntity.RoleNutritionist`.

---

### [LOW] Raw exported field access on `User` entity bypasses domain encapsulation

**File:** `food_request_service.go:48`
**Issue:** `client.NutritionistID == nil` accesses exported field directly.
**DDD Rule:** Aggregates MUST NOT expose raw entity fields; use getter/behaviour methods.
**Fix:** Add `User.HasNutritionist() bool` method to domain entity; use it here.

---

### [LOW] `ApproveRequest` DTO duplicates `appFood.CreateFoodRequest` with no added semantics

**File:** `food_request_service.go:14–22`
**Issue:** `ApproveRequest` is a 7-field subset of `appFood.CreateFoodRequest` with no distinct semantics. Any food field addition requires changes in two places.
**Fix:** Either reuse `appFood.CreateFoodRequest` directly or document explicitly why the subset is intentional.

---

## Compliant Patterns Found

- **Repository interfaces used correctly** — accepts `frRepo.FoodRequestRepository` and `userRepo.UserRepository` as interfaces. ✓
- **No infrastructure imports**. ✓
- **Domain errors used** — `shared.ErrForbidden`, `shared.ErrFoodRequestNotFound`, etc. ✓
- **Domain entity behaviour methods used for status checks** — `foodReq.IsPending()` used correctly. ✓
- **Factory constructor present**. ✓

## Fix Priority Order
1. **[HIGH]** Define `FoodCreator` interface; decouple from concrete `*appFood.FoodService`
2. **[MEDIUM]** Add `NewFoodRequest()` domain factory; replace struct literal in `Submit`
3. **[MEDIUM]** Replace hard-coded role strings with domain methods/constants
4. **[LOW]** Add `HasNutritionist()` to `User` entity; remove raw field access
5. **[LOW]** Rationalise `ApproveRequest` DTO
