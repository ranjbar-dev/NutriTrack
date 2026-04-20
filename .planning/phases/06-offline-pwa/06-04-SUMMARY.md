---
phase: "06"
plan: "04"
subsystem: frontend/offline-pwa
tags: [offline, pwa, sync-queue, stores, composables]
dependency_graph:
  requires: [06-03]
  provides: [offline-store-writes, sync-status-ui, pwa-install-update]
  affects: [frontend/app/stores, frontend/app/layouts/client.vue]
tech_stack:
  added: []
  patterns: [offline-first clientPost pattern, queued-write guard]
key_files:
  created:
    - frontend/app/components/client/ClientSyncStatus.vue
    - frontend/app/composables/usePWA.client.ts
  modified:
    - frontend/app/stores/foodLog.ts
    - frontend/app/stores/waterLog.ts
    - frontend/app/stores/sleepLog.ts
    - frontend/app/stores/exerciseLog.ts
    - frontend/app/stores/medicationLog.ts
    - frontend/app/stores/bodyMeasurement.ts
    - frontend/app/stores/labResult.ts
    - frontend/app/layouts/client.vue
decisions:
  - "bodyMeasurement nutritionist path stays online-only (apiFetch); only client path uses clientPost"
  - "labResult file uploads remain online-only with early return if offline; metadata-only (link_url) path can be queued"
  - "usePWA.client.ts kept as explicit import in client.vue to avoid collision with @vite-pwa/nuxt auto-import"
metrics:
  duration: "~8 minutes"
  completed: "2025-07-14"
  tasks_completed: 2
  files_changed: 10
---

# Phase 06 Plan 04: Wire tracking stores + ClientSyncStatus + client.vue Summary

Wire all 7 client tracking stores to use `useOfflineApi().clientPost` for offline-capable POST writes, add `ClientSyncStatus.vue` for real-time sync queue feedback, create `usePWA.client.ts` for install/update prompts, and extend `client.vue` layout to surface all offline-PWA features.

## Tasks Completed

### Task 1: Wire 7 tracking stores with useOfflineApi.clientPost

All seven stores replaced their `apiFetch` POST calls with `clientPost` from `useOfflineApi`:

| Store | Entity Type | Guard Pattern |
|-------|------------|---------------|
| `foodLog.ts` | `food_log` | `if (!('queued' in result)) { upsertLocal(result) }` |
| `waterLog.ts` | `water_log` | `if (!('queued' in result)) { logs.value.push(result) }` |
| `sleepLog.ts` | `sleep_log` | `if (!('queued' in result)) { todayLog.value = result }` |
| `exerciseLog.ts` | `exercise_log` | `if (!('queued' in result)) { todayLogs.value.unshift(result) }` |
| `medicationLog.ts` | `medication_log` | `if (!('queued' in result)) { todayLogs.value.unshift(result) }` |
| `bodyMeasurement.ts` | `body_measurement` (client path only) | No guard needed (result unused, refetch after) |
| `labResult.ts` | `lab_result_meta` (metadata path) | `if (!('queued' in result)) { labResults.value.unshift(result) }` |

Special cases implemented:
- **bodyMeasurement**: nutritionist path (`/nutritionist/clients/:id/body-measurements`) stays online-only with `apiFetch`; client path uses `clientPost`
- **labResult**: file uploads (`formData.get('file')`) stay online-only with offline guard returning an error message; metadata-only (`link_url`) path is queued offline

### Task 2: Create ClientSyncStatus.vue, usePWA.client.ts, extend client.vue

- **ClientSyncStatus.vue**: Reads `pendingCount`, `failedCount`, `isProcessing` from `useSyncQueueStore` via `storeToRefs`; shows colored indicator dot + Persian status text; retry button for failed items
- **usePWA.client.ts**: Module-level `deferredPrompt` and `swRegistration` singletons; handles `beforeinstallprompt` event, SW `updatefound`/`statechange` cycle; `promptInstall()` and `applyUpdate()` exported
- **client.vue**: Extended with `<OfflineBanner />`, `<ClientSyncStatus />`, install prompt banner, update banner; `useSyncManager().processQueue` called on mount; RTL-safe logical properties (`start-4`, `end-4`)

## Deviations from Plan

### Pre-existing Issue (not auto-fixed — out of scope)

**[Pre-existing] PWA service worker build path error**
- **Found during:** Build verification
- **Issue:** `Could not resolve entry module "app/app/service-worker/sw.ts"` — double `app/` prefix in path configuration. This error existed before this plan.
- **Impact:** Final PWA SW bundling step fails, but TypeScript compilation for client and server both succeed (`✓ Client built in 3792ms`, `✓ Server built in 1971ms`)
- **Deferred to:** `deferred-items.md` — fix PWA nuxt.config.ts `srcDir` / SW path

### Auto-import collision (handled in code)

**`usePWA` name collision with `@vite-pwa/nuxt`**
- Nuxt auto-import warning: our `usePWA.client.ts` shadowed by `@vite-pwa/nuxt/dist/runtime/composables/index`
- **Resolution**: `client.vue` uses explicit path import `import { usePWA } from '~/composables/usePWA.client'` which bypasses auto-import resolution — our implementation is used correctly

## Known Stubs

None — all wired data flows through real store state.

## Self-Check: PASSED

- ✅ `frontend/app/components/client/ClientSyncStatus.vue` — exists
- ✅ `frontend/app/composables/usePWA.client.ts` — exists
- ✅ `frontend/app/layouts/client.vue` — modified with all required sections
- ✅ All 7 stores contain `useOfflineApi` import
- ✅ Commit `108550a` exists: `feat(06-04): wire 7 tracking stores with useOfflineApi + ClientSyncStatus + usePWA + client.vue wiring`
