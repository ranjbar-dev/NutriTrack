---
phase: 02
plan: 01
subsystem: domain/user + application/auth
tags: [domain, user, jwt, bcrypt, mobile, value-object, migration]
dependency_graph:
  requires: [01-04]
  provides: [UserRepository interface, JWTService, mobile value object, users migration]
  affects: [02-02, 02-03, 02-04]
tech_stack:
  added: [github.com/golang-jwt/jwt/v5 v5.3.1]
  patterns: [DDD aggregate, value object, repository interface, application service]
key_files:
  created:
    - internal/domain/user/entity/user.go
    - internal/domain/user/repository/user_repository.go
    - internal/domain/user/valueobject/mobile.go
    - internal/application/auth/jwt_service.go
    - internal/application/auth/password_service.go
    - migrations/000002_users.up.sql
    - migrations/000002_users.down.sql
    - db/queries/users.sql
  modified:
    - go.mod
    - go.sum
decisions:
  - "JWT uses separate access/refresh secrets with HS256 — simpler key rotation than RS256 for single-service deployment"
  - "bcrypt cost=12 — hardened against GPU attacks while remaining <300ms on modern hardware"
  - "Mobile value object normalises 09x → +98x at construction; domain always stores E.164"
metrics:
  duration: ~10 min
  completed: 2026-04-21
---

# Phase 2 Plan 1: User domain aggregate, JWT service, bcrypt password hashing — Summary

**One-liner:** Pure-domain User aggregate with Iranian mobile E.164 normalisation, JWT access(15 min)+refresh(7 day) token pair generation, and bcrypt cost-12 password hashing.

## What Was Built

### Domain Layer (`internal/domain/user/`)

| File | Purpose |
|------|---------|
| `entity/user.go` | User root aggregate — roles (superadmin/nutritionist/client), BMI(), BelongsTo() guard, helper predicates |
| `repository/user_repository.go` | `UserRepository` interface — 12 methods covering CRUD, pagination, exists-checks |
| `valueobject/mobile.go` | `Mobile` value object — Iranian mobile validation (09x / +989x), normalises to E.164 on construction |

### Application Layer (`internal/application/auth/`)

| File | Purpose |
|------|---------|
| `jwt_service.go` | `JWTService` — `GenerateTokenPair`, `ValidateAccessToken`, `ValidateRefreshToken`; separate secrets for access/refresh |
| `password_service.go` | `HashPassword` (bcrypt cost=12) and `CheckPassword` functions |

### Persistence Artefacts

| File | Purpose |
|------|---------|
| `migrations/000002_users.up.sql` | `users` table — role/gender CHECK constraints, 4 indexes, superadmin seed row |
| `migrations/000002_users.down.sql` | Rollback: `DROP TABLE IF EXISTS users` |
| `db/queries/users.sql` | 11 sqlc queries — GetByID/Mobile/Email, Create, Update, List/Count by nutritionist, List/Count nutritionists, Delete, ExistsBy* |

## Architecture Compliance

- ✅ `internal/domain/` — zero external dependencies (only `github.com/google/uuid` which is project-wide)
- ✅ `internal/application/` — calls `configs` for JWT config; no infrastructure imports
- ✅ Value objects are immutable structs with validation in constructors
- ✅ Repository interface in `internal/domain/user/repository/`; no implementation here

## Deviations from Plan

None — plan executed exactly as written. The `golang-jwt/jwt/v5` dependency was added via `go get` as expected (it was not yet in go.mod).

## Decisions Made

1. **JWT HS256 with separate secrets** — `JWTConfig.AccessSecret` and `JWTConfig.RefreshSecret` are distinct, allowing independent rotation without invalidating all tokens of the other type.
2. **bcrypt cost=12** — Hardened against GPU brute-force while staying under 300ms on modern hardware; configurable in future if benchmarks demand it.
3. **Mobile E.164 normalisation at construction** — `NewMobile("09123456789")` stores `+989123456789`; the domain always operates on E.164, making storage and comparison consistent.

## Known Stubs

None — all domain logic is fully wired. Infrastructure (PostgreSQL implementation of `UserRepository`) is deferred to Plan 02-02 as per the roadmap.

## Self-Check: PASSED

- `internal/domain/user/entity/user.go` — ✅ exists
- `internal/domain/user/repository/user_repository.go` — ✅ exists
- `internal/domain/user/valueobject/mobile.go` — ✅ exists
- `internal/application/auth/jwt_service.go` — ✅ exists
- `internal/application/auth/password_service.go` — ✅ exists
- `migrations/000002_users.up.sql` — ✅ exists
- `migrations/000002_users.down.sql` — ✅ exists
- `db/queries/users.sql` — ✅ exists
- Commit `f3a22b6` — ✅ exists (`feat(02-01): User domain aggregate...`)
- `go build ./...` — ✅ green (exit 0)
