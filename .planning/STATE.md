# Project State

## Project Reference

See: .planning/PROJECT.md (updated 2025-07-18)

**Core value:** Digitalize the nutritionist–client workflow in Iran — replacing WhatsApp + Excel + paper with a structured, offline-capable PWA
**Current focus:** Phase 2: Core Data Domain

## Current Position

Phase: 2 of 7 (Core Data Domain)
Plan: 3 of 4 in current phase
Status: In Progress
Last activity: 2026-04-19 — Completed 02-02 Food CRUD Backend

Progress: [███░░░░░░░] 23%

## Performance Metrics

**Velocity:**
- Total plans completed: 8
- Average duration: 8.75 min
- Total execution time: 1.17 hours

**By Phase:**

| Phase | Plans | Total | Avg/Plan |
|-------|-------|-------|----------|
| 1 | 6/6 | 45 min | 7.5 min |
| 2 | 2/4 | 25 min | 12.5 min |

**Recent Trend:**
- Last 5 plans: 10 min, 15 min, 13 min, 4 min, 4 min
- Trend: stable

*Updated after each plan completion*

## Accumulated Context

### Decisions

Decisions are logged in PROJECT.md Key Decisions table.
Recent decisions affecting current work:

- [Roadmap]: HTTP framework is Gin (not Fiber/Echo) per user directive
- [Roadmap]: RTL via Tailwind v4 native logical properties (no plugin needed)
- [Roadmap]: Phases 4 & 5 are parallel-safe but execute sequentially (solo dev)
- [Roadmap]: Phase 6 must follow both 4 and 5 (offline wraps all API endpoints)
- [Roadmap]: local_id dedup infrastructure built in Phase 4, consumed by Phase 6
- [01-01]: Used golang-migrate pgx/v5 driver (not postgres driver) to match pgxpool
- [01-01]: Go auto-upgraded 1.24.2 → 1.25.0 (Gin v1.12.0 requires Go 1.25+)
- [01-01]: Custom validator placed in internal/validator package for reusability
- [01-02]: Vazirmatn npm version 33.0.3 (plan assumed 35.0.1 which doesn't exist)
- [01-02]: Used Vazirmatn-Variable-font-face.css for optimal variable font loading
- [01-02]: jalaali-js CJS default import (no ESM/TS types available)
- [01-02]: @pinia/nuxt upgraded to 0.11.3 for Nuxt 4.4.2 compatibility
- [01-03]: JWT tokens signed with HMAC-SHA256 only; ParseToken validates signing method
- [01-03]: Rate limiter peeks JSON body for mobile field without consuming it
- [01-03]: Security headers middleware added (not in plan, added via Rule 2)
- [01-04]: Added GetRefreshTokenByHashAny sqlc query for replay detection (existing query blocked theft detection)
- [01-04]: Recovery middleware created as separate file (not in 01-03 middleware suite)
- [01-04]: AuthService.GetUserByID added to maintain layered arch (handlers never call repos directly)
- [01-05]: Updated apiBase to include /api prefix to match backend route structure
- [01-05]: Auth store normalizes both mobile and OTP code via toLatinDigits() for API calls
- [01-05]: Unauthorized page uses auth layout with middleware:[] to avoid redirect loop
- [01-06]: No version key in docker-compose (Docker Compose v2+ ignores it)
- [01-06]: Frontend priority=1 in Traefik labels ensures /api/* matches first
- [02-02]: Renamed food enums in migrations so sqlc could generate food models without colliding with the food_categories table name

### Pending Todos

None yet.

### Blockers/Concerns

- [Phase 3]: Highest technical risk — design spike needed on batch loading queries and plan builder UI state management
- [Phase 6]: iOS PWA storage eviction — test on real devices during Phase 6

## Deferred Items

| Category | Item | Status | Deferred At |
|----------|------|--------|-------------|
| *(none)* | | | |

## Session Continuity

Last session: 2026-04-19
Stopped at: Completed 02-02 Food CRUD Backend
Resume file: .planning/phases/02-core-data-domain/02-03-PLAN.md
