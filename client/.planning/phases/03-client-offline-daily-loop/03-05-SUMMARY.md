---
phase: 03-client-offline-daily-loop
plan: 05
type: execute
wave: 4
status: completed
completed_at: 2026-04-23T00:00:00Z
requirements:
  - CLNT-01
  - TRCK-03
  - OFFL-03
---

# Phase 3, Plan 03-05 - Tracking History and Manual Retry Summary

## Outcome
Implemented tracking history surface, lightweight progress summaries, and deterministic failed-sync manual retry actions to complete the client daily loop continuity.

## Implemented Work

### Tracking History and Progress Surface
- Added `app/pages/client/history/tracking.vue`:
  - shows recent tracking entries from offline queue/store
  - row-level sync chips per entry state
  - Persian-formatted timestamps
  - manual retry integration for failed entries
- Added `app/components/client/TrackingProgressSummary.vue`:
  - water completion summary block (based on available v1 data)
  - lightweight recent activity signal
- Added history utility functions in `app/lib/tracking/history.ts`:
  - timestamp formatting for Persian digits
  - domain label mapping
  - progress summary derivation from queue entries

### Failed Sync Inspection and Manual Retry
- Added `app/components/client/FailedSyncList.vue`:
  - visible failed entries list
  - single retry and retry-all actions
  - persistent Persian guidance copy for unresolved failures
- Added retry logic helpers in `app/lib/tracking/retry.ts`:
  - immutable failed-entry retry by `local_id`
  - bulk retry for all failed entries
  - failure guidance messaging for persistent failure states

### History Navigation Continuity
- Updated `app/components/platform/BottomNavClient.vue` history route to `/client/history/tracking`.
- Added `app/pages/client/history/index.vue` landing links for plans and tracking history surfaces.

## Threat Mitigations Applied
- T-03-41: Retry actions target immutable `local_id` and validate failed-state preconditions.
- T-03-42: Per-entry state transitions are explicit (`failed` -> `queued`) and reflected in aggregate metrics.
- T-03-43: Progress summary limited to v1 signals (recent records and water completion), avoiding unsupported analytics claims.

## Tests Added
- `tests/client/tracking-history-progress.spec.ts`
- `tests/client/manual-retry.spec.ts`

## Verification Results
- `npm run typecheck` -> PASS (exit 0)
- `npm run lint` -> PASS (exit 0)
- `npm run test:unit -- tests/client/` -> PASS (10 files, 68 tests)

## Files Changed (Plan Scope)
- `app/components/client/FailedSyncList.vue`
- `app/components/client/TrackingProgressSummary.vue`
- `app/pages/client/history/tracking.vue`
- `app/pages/client/history/index.vue`
- `app/components/platform/BottomNavClient.vue`
- `tests/client/tracking-history-progress.spec.ts`
- `tests/client/manual-retry.spec.ts`

## Additional Supporting Files
- `app/lib/tracking/history.ts`
- `app/lib/tracking/retry.ts`
- `app/components/client/SyncStateChip.vue`
