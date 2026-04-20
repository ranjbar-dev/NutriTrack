---
phase: 06-offline-pwa
plan: 02
subsystem: ui
tags: [dexie, indexeddb, pinia, offline, vue, typescript, vitest]

requires:
  - phase: 06-01
    provides: dexie package installed, service worker skeleton

provides:
  - "NutriTrackDB Dexie singleton with 6 tables: activePlan, messages, syncQueue, syncMeta, notificationPreferences, uiState"
  - "clientPlan store cache-first: reads IndexedDB instantly, refreshes from API when online"
  - "iOS storage eviction detection: detects empty cache after recent fetch and sets eviction flag"
  - "OfflineBanner.vue component: offline indicator + eviction notice variants"
  - "useDb.test.ts: 4 unit tests for schema init and eviction detection"

affects:
  - 06-03 (syncQueue table used by useOfflineApi + useSyncQueueStore)
  - 06-04 (all tracking stores use db from this plan)
  - 06-05 (message store uses db.messages table)
  - 06-07 (notificationPreferences table used by notificationPrefs store)

tech-stack:
  added:
    - "fake-indexeddb@6.2.5 (dev, for vitest IDB mocking)"
  patterns:
    - "Dexie lazy singleton via Proxy — SSR-safe, never instantiated at module load time"
    - "Cache-first pattern: IDB read → show cached → refresh from API in background"
    - "iOS eviction detection via syncMeta plan_last_fetch timestamp comparison"

key-files:
  created:
    - frontend/app/db/index.ts
    - frontend/app/components/OfflineBanner.vue
    - frontend/tests/useDb.test.ts
  modified:
    - frontend/app/stores/clientPlan.ts (added Dexie cache-first read + iOS eviction check)
    - frontend/package.json (added fake-indexeddb devDep)

key-decisions:
  - "NutriTrackDB exported as a class + lazy Proxy singleton for SSR safety (Pitfall 3)"
  - "activePlan table uses id=1 singleton pattern for simplicity"
  - "iOS eviction check: if activePlan is empty AND plan_last_fetch is < 1 hour old, set eviction_detected flag"
  - "OfflineBanner uses window online/offline events for real-time offline detection"
  - "NotificationPrefRecord schema uses {id:1, data: object, updated_at} instead of per-field booleans (simplified)"

patterns-established:
  - "All client-only DB access via lazy db proxy from ~/db"
  - "Eviction detection uses uiState 'eviction_detected' key; set by store, cleared on successful fetch"

requirements-completed:
  - OFFL-02
  - OFFL-05
  - OFFL-10
  - OFFL-12

duration: ~20min
completed: 2026-04-20
---

# Phase 06-02: Dexie Schema + Plan Cache Summary

**NutriTrackDB with 6 Dexie tables, cache-first clientPlan store, iOS eviction detection, and OfflineBanner component**

## Performance

- **Duration:** ~20 min
- **Completed:** 2026-04-20
- **Tasks:** 2
- **Files modified:** 6

## Accomplishments
- Created NutriTrackDB Dexie singleton with 6 versioned tables (activePlan, messages, syncQueue, syncMeta, notificationPreferences, uiState)
- Extended clientPlan store with three-step cache-first pattern: IDB read → offline guard → API refresh with cache write
- Implemented iOS eviction detection using syncMeta timestamp comparison with 1-hour window
- Created OfflineBanner.vue showing offline indicator and dismissible eviction notice in Persian

## Task Commits

1. **Task 1: Create Dexie NutriTrackDB schema and vitest config** - `3b59bc3` (feat)
2. **Task 2: Extend clientPlan store + create OfflineBanner** - `3b59bc3` (feat)

## Files Created/Modified
- `frontend/app/db/index.ts` - NutriTrackDB Dexie v1 schema, lazy SSR-safe Proxy singleton
- `frontend/app/stores/clientPlan.ts` - Cache-first fetchActivePlan with iOS eviction detection
- `frontend/app/components/OfflineBanner.vue` - Offline banner + eviction notice with dismiss button
- `frontend/tests/useDb.test.ts` - 4 vitest unit tests: schema init, CRUD, iOS eviction scenario
- `frontend/package.json` - Added fake-indexeddb dev dependency

## Decisions Made
- NotificationPrefRecord uses simplified `{id, data: object, updated_at}` instead of per-field booleans — avoids schema migration if preference fields change
- Dexie exported as lazy Proxy (not direct `new NutriTrackDB()`) to prevent SSR-time instantiation errors

## Deviations from Plan
None - plan executed as specified.

## Issues Encountered
None

## User Setup Required
None - no external service configuration required.

## Next Phase Readiness
- syncQueue Dexie table ready for Plan 06-03 offline write pipeline
- messages table ready for Plan 06-05 message cache
- OfflineBanner component created but not yet wired into client.vue layout (done in Plan 06-04)

---
*Phase: 06-offline-pwa*
*Completed: 2026-04-20*
