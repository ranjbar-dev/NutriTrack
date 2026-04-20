# Project State

## Project Reference

See: .planning/PROJECT.md (updated 2025-07-18)

**Core value:** Digitalize the nutritionist–client workflow in Iran — replacing WhatsApp + Excel + paper with a structured, offline-capable PWA
**Current focus:** Phase 6: Offline & PWA

## Current Position

Phase: 6 of 7 (Offline & PWA)
Plan: 6 of 7 (06-06 next)
Status: In Progress
Last activity: 2026-04-20 — Completed 06-03 offline write pipeline (useOfflineApi, useSyncQueueStore, useSyncManager)

Progress: [█████████░] 74%

## Performance Metrics

**Velocity:**
- Total plans completed: 29
- Average duration: retroactively reconciled
- Total execution time: retroactively reconciled

**By Phase:**

| Phase | Plans | Total | Avg/Plan |
|-------|-------|-------|----------|
| 1 | 6/6 | 45 min | 7.5 min |
| 2 | 4/4 | 37 min | 9.25 min |
| 3 | 10/10 | completed | n/a |
| 4 | 9/9 | completed | n/a |
| 5 | 7/7 | completed | n/a |

**Recent Trend:**
- Recent completion: Phase 5 completed — messaging, food requests, client management dashboard
- Trend: ready to begin Phase 6 planning

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
- [02-03]: labelMap[key] ?? key fallback used for strict TypeScript index access in validation loop
- [02-03]: draftReady ref gates localStorage watcher to prevent overwriting draft on initial form mount
- [02-04]: Reused optionalText/optionalBool/formatTimestamp helpers from food_service.go — same package, no duplication needed
- [02-04]: Duplicate check uses COALESCE pattern (same as food) to handle NULL exclude_id safely in SQL
- [06-03]: isTransportError detects TypeError (fetch fails) vs statusCode (server error) for queue routing
- [06-03]: Background Sync registration is best-effort (Chromium only) with 30s interval fallback
- [06-03]: #app alias + fake-indexeddb setupFiles added to vitest.config.ts for test env compatibility
- [06-05]: sendMessage return type changed to Promise<void>; local echo pushes directly to messages array
- [06-05]: hook callback cast `as never` matches syncQueue.ts callHook pattern for Nuxt strict hook typing

### Pending Todos

None yet.

### Blockers/Concerns

- [Phase 6]: iOS PWA storage eviction — test on real devices during Phase 6

## Deferred Items

| Category | Item | Status | Deferred At |
|----------|------|--------|-------------|
| *(none)* | | | |

## Session Continuity

Last session: 2026-04-20
Stopped at: Completed 06-05 — message store cache-first reads (D-07) + offline send queue (D-08) + offline tests
Resume file: .planning/phases/06-offline-pwa/06-06-PLAN.md
