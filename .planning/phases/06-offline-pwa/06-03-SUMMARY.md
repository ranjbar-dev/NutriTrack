---
phase: 06-offline-pwa
plan: 03
subsystem: ui
tags: [pinia, dexie, offline, background-sync, typescript, vitest]

requires:
  - phase: 06-02
    provides: db.syncQueue Dexie table, NutriTrackDB schema

provides:
  - "useOfflineApi composable: clientPost() queues to Dexie when offline or transport error"
  - "useSyncQueueStore Pinia store: enqueue, processQueue FIFO, exponential backoff (0/1/2/4s), retryFailed"
  - "useSyncManager composable: online events + Background Sync + 30s interval fallback"

affects:
  - 06-04 (all tracking stores use useOfflineApi.clientPost)
  - 06-05 (message store uses useOfflineApi.clientPost)

tech-stack:
  added: []
  patterns:
    - "clientPost wraps apiFetch: offline/transport error → queue, 4xx/5xx → rethrow"
    - "Single-flight processQueue guard via isProcessing ref"
    - "Exponential backoff: retry_count 0→0ms, 1→1s, 2→2s, ≥3→4s then failed"

key-files:
  created:
    - frontend/app/composables/useOfflineApi.ts
    - frontend/app/stores/syncQueue.ts
    - frontend/app/composables/useSyncManager.ts
    - frontend/tests/useSyncManager.test.ts
    - frontend/tests/setup.ts
    - frontend/tests/__mocks__/nuxt-app.ts
  modified:
    - frontend/vitest.config.ts

key-decisions:
  - "isTransportError detects TypeError (fetch fails) vs statusCode (server error)"
  - "Background Sync registration is best-effort (Chromium only) with 30s interval fallback"
  - "[Rule 3] Added #app alias + fake-indexeddb setupFiles to vitest.config.ts — #app not resolvable in Vitest without Nuxt runtime; fake-indexeddb/auto must polyfill before Dexie opens any connection"

requirements-completed:
  - OFFL-06
  - OFFL-07
  - OFFL-08
  - OFFL-09

duration: ~25min
completed: 2026-04-20
---

# Phase 06-03: Offline Write Pipeline Summary

**useOfflineApi + useSyncQueueStore + useSyncManager — POST requests queue to Dexie when offline with FIFO exponential backoff replay**

## What Was Built

### `useOfflineApi.ts`
Wraps `apiFetch` for all client POST endpoints. Decision tree:
1. Path not in `QUEUEABLE_PATH_PATTERNS` → pass through directly
2. `navigator.onLine === false` → enqueue immediately
3. Online → attempt `apiFetch`; if `TypeError` (network failure) → enqueue; if HTTP 4xx/5xx → rethrow

Queueable paths: `/client/food-logs`, `/client/water-logs`, `/client/sleep-logs`, `/client/exercise-logs`, `/client/medication-logs`, `/client/body-measurements`, `/client/lab-results`, `/messages`

### `useSyncQueueStore.ts`
Pinia store managing the Dexie `syncQueue` table:
- **`enqueue(path, body, options)`** — adds entry with `status: 'pending'`, assigns `local_id` (from body, options, or `crypto.randomUUID()`)
- **`processQueue()`** — single-flight guard via `isProcessing`, filters `pending` items with `next_attempt_at ≤ now`, processes in `created_at` FIFO order
- **`processSingleItem()`** — sets `syncing`, calls `apiFetch` (or `FormData` for blobs), fires `sync:itemSynced` Nuxt hook on success, applies exponential backoff on failure
- **`backoffMs(retryCount)`** — `0→0ms`, `1→1s`, `2→2s`, `≥3→4s`; item marked `failed` after 3 retries
- **`retryFailed()`** — resets failed items to pending, triggers processQueue
- **`refreshCounts()`** — updates `pendingCount` and `failedCount` reactive refs

### `useSyncManager.ts`
Lifecycle composable (call in root app layout):
- Registers 30s polling interval (processes queue when online)
- Listens for `window.online` events
- Registers Background Sync tag `nutritrack-sync` via SW (Chromium only, best-effort)
- Listens for `TRIGGER_SYNC` postMessage from SW (all browsers)
- `onMounted(start)` / `onUnmounted(stop)` for clean lifecycle management

## Verification

- sw.ts confirmed: has `TRIGGER_SYNC` background sync handler (from plan 06-01) ✓
- All 12 tests pass: 3 new (OFFL-06, OFFL-07, OFFL-08) + 9 existing ✓

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 3 - Blocking] Added `#app` alias to vitest.config.ts**
- **Found during:** Task 1 test run
- **Issue:** Vitest couldn't resolve `#app` (Nuxt internal alias) without Nuxt runtime; `vi.mock('#app', ...)` can't intercept what Vite can't resolve
- **Fix:** Added `'#app': resolve(__dirname, 'tests/__mocks__/nuxt-app.ts')` alias + minimal stub; vi.mock overrides the stub at test time
- **Files modified:** `frontend/vitest.config.ts`, `frontend/tests/__mocks__/nuxt-app.ts`

**2. [Rule 3 - Blocking] Added `setupFiles` with `fake-indexeddb/auto` to vitest.config.ts**
- **Found during:** Task 1 test run (2nd attempt)
- **Issue:** `fake-indexeddb/auto` inside async `vi.mock` factory didn't polyfill `indexedDB` before Dexie's first connection attempt; `DatabaseClosedError: MissingAPIError IndexedDB API missing`
- **Fix:** Added `frontend/tests/setup.ts` (imports `fake-indexeddb/auto`) as vitest `setupFiles` — ensures polyfill runs before any test module is evaluated
- **Files modified:** `frontend/vitest.config.ts`, `frontend/tests/setup.ts`

## Self-Check: PASSED

Files exist:
- `frontend/app/composables/useOfflineApi.ts` ✓
- `frontend/app/stores/syncQueue.ts` ✓
- `frontend/app/composables/useSyncManager.ts` ✓
- `frontend/tests/useSyncManager.test.ts` ✓
- `frontend/tests/setup.ts` ✓
- `frontend/tests/__mocks__/nuxt-app.ts` ✓

Commit `61292a4` exists ✓
