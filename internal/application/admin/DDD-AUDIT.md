# DDD Audit: internal/application/admin
Layer: application
Audited: 2026-04-22
Files reviewed: 1 (admin_service.go)

## Summary
- CRITICAL: 2
- HIGH: 2
- MEDIUM: 0
- LOW: 0
- PASS: 0

---

## Findings

### [CRITICAL] Direct import of `internal/infrastructure` in application layer

**File:** `admin_service.go` (import block)
**Issue:** The application service file imports a package from `internal/infrastructure` directly. The application layer MUST NOT depend on infrastructure. All infrastructure access must be behind a domain port interface.
**DDD Rule:** "Application layer: depends only on domain packages and domain port interfaces."
**Fix:** Define an `AdminStatsReader` port interface in a domain package (or a new `internal/domain/admin/` package). Implement it in infrastructure. Inject through the interface.

---

### [CRITICAL] Service struct and factory depend on concrete infrastructure type (`*db.Queries`)

**File:** `admin_service.go`
**Issue:** The `AdminService` struct stores `*db.Queries` (a generated SQLC concrete type from the infrastructure layer) as a field. The constructor `NewAdminService(*db.Queries)` accepts the concrete type.
**DDD Rule:** Application services MUST depend on domain interfaces, never on infrastructure concrete types.
**Fix:** Define a `StatsRepository` or `AdminRepository` interface in the domain layer; inject that interface instead.

---

### [HIGH] `GetStats` returns infrastructure-layer DTO with `db:` struct tags

**File:** `admin_service.go` (GetStats return value)
**Issue:** The service returns a struct that originates from SQLC (`db.GetAdminStatsRow` or similar), complete with `db:` struct tags. Infrastructure DTOs must never cross the application boundary.
**DDD Rule:** Application services MUST return domain types or application-layer DTOs, never infrastructure types.
**Fix:** Define an `AdminStats` struct in the application or domain layer (no struct tags); map from the SQLC result inside the repository implementation.

---

### [HIGH] No domain package exists for admin (`internal/domain/admin/` is absent)

**File:** N/A (missing directory)
**Issue:** The admin bounded context has an application service but no corresponding domain layer. Business rules about which stats are computed or what constitutes valid admin activity have nowhere to live.
**DDD Rule:** Each bounded context should have a corresponding domain package.
**Fix (optional for now):** Create `internal/domain/admin/` with at minimum a `StatsRepository` interface.

---

## Compliant Patterns Found

- None — this package has the most severe architectural violations.

## Fix Priority Order
1. **[CRITICAL]** Define `AdminRepository` / `AdminStatsReader` interface in a domain or shared-domain package
2. **[CRITICAL]** Remove `*db.Queries` from service struct; inject through the domain interface
3. **[HIGH]** Define `AdminStats` application DTO; map from infrastructure in repository implementation
4. **[HIGH]** Create `internal/domain/admin/` package with `StatsRepository` interface
