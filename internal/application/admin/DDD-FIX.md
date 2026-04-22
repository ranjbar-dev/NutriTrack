# DDD Fix Report: internal/application/admin
Layer: application
Fixed: 2026-04-22
Based on: DDD-AUDIT.md

## Baseline Build Status
PASS — `go build ./...` before fixes

## Fix Plan

| # | Finding | Severity | Files | Strategy | Status |
|---|---------|----------|-------|----------|--------|
| 1 | Direct import of `internal/infrastructure` in application layer | CRITICAL | admin_service.go | SAFE | FIXED |
| 2 | Service struct/factory depend on concrete `*db.Queries` | CRITICAL | admin_service.go | SAFE | FIXED |
| 3 | `GetStats` returns infrastructure-layer DTO (`db.AdminStats`) | HIGH | admin_service.go | SAFE | FIXED |
| 4 | No domain package exists for admin (`internal/domain/admin/` absent) | HIGH | (new package) | SAFE | FIXED |

## Changes Applied

### Fix 4 + 1: Create `internal/domain/admin/` package with `AdminStats` + `AdminRepository` interface
**File:** `internal/domain/admin/admin_repository.go` (new)
**Change:** Created domain package with:
- `AdminStats` struct — no `db:`/`json:` struct tags, pure domain DTO
- `AdminRepository` interface with `GetStats(ctx context.Context) (AdminStats, error)`
**Build:** PASS

### Fix 2 + 3: Create infrastructure concrete implementation
**File:** `internal/infrastructure/persistence/admin/admin_repository.go` (new)
**Change:** `PgAdminRepository` implements `domain/admin.AdminRepository`. Maps from `db.AdminStats` (SQLC infra DTO) to `domainadmin.AdminStats` (domain DTO) inside the repository, keeping infrastructure types contained within the infrastructure layer.

**Before:**
```go
// Infrastructure DTO leaked into application and handler layers
func (s *AdminService) GetStats(ctx context.Context) (db.AdminStats, error) {
    return s.queries.GetAdminStats(ctx)
}
```

**After:**
```go
// Infrastructure repo maps to domain DTO; application and handler see only domain type
func (r *PgAdminRepository) GetStats(ctx context.Context) (domainadmin.AdminStats, error) {
    row, err := r.queries.GetAdminStats(ctx)
    if err != nil {
        return domainadmin.AdminStats{}, err
    }
    return domainadmin.AdminStats{
        TotalNutritionists:    row.TotalNutritionists,
        // ...
    }, nil
}
```
**Build:** PASS

### Fix 1 + 2: Replace `*db.Queries` dependency with domain interface in AdminService
**File:** `internal/application/admin/admin_service.go`
**Before:**
```go
import db "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/sqlc"

type AdminService struct {
    queries *db.Queries
}

func NewAdminService(q *db.Queries) *AdminService {
    return &AdminService{queries: q}
}

func (s *AdminService) GetStats(ctx context.Context) (db.AdminStats, error) {
    return s.queries.GetAdminStats(ctx)
}
```
**After:**
```go
import domainadmin "github.com/ranjbar-dev/nutritrack/internal/domain/admin"

type AdminService struct {
    repo domainadmin.AdminRepository
}

func NewAdminService(repo domainadmin.AdminRepository) *AdminService {
    return &AdminService{repo: repo}
}

func (s *AdminService) GetStats(ctx context.Context) (domainadmin.AdminStats, error) {
    return s.repo.GetStats(ctx)
}
```
**Build:** PASS

### Wire update: `bootstrap/wire.go`
**Change:** Added `infraAdmin` import for new persistence/admin package. Changed wiring from direct `dbsqlc.New(db)` to `infraAdmin.NewPgAdminRepository(dbsqlc.New(db))`.
**Before:**
```go
adminSvc := appAdmin.NewAdminService(dbsqlc.New(db))
```
**After:**
```go
adminSvc := appAdmin.NewAdminService(infraAdmin.NewPgAdminRepository(dbsqlc.New(db)))
```
**Build:** PASS

## Final Build Status
PASS — `go build ./...` after all fixes
PASS — `go vet ./internal/...` after all fixes

## Remaining Violations
None — all CRITICAL and HIGH findings resolved.
