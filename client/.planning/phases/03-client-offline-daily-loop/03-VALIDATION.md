# Phase 03 Validation Strategy

**Phase:** 03 - Client Offline Daily Loop
**Date:** 2026-04-23
**Status:** Active

## Validation Objectives

- Validate client daily loop usability and correctness for plan visibility plus tracking entries.
- Validate offline queue durability, replay behavior, and sync-state transparency.
- Validate strict boundary that offline write support remains client-tracking only.

## Requirement Coverage Matrix

| Requirement | Validation Paths | Evidence Artifact |
|-------------|------------------|-------------------|
| CLNT-01 | `tests/client/today-sync-summary.spec.ts` | 03-02 summary + unit output |
| CLNT-02 | `tests/client/plan-flatten-map.spec.ts` | 03-03 summary + unit output |
| CLNT-03 | `tests/client/plan-history-context.spec.ts` | 03-03/03-05 summaries + unit output |
| TRCK-01 | `tests/offline/food-queue-replay.spec.ts` | 03-04 summary + unit output |
| TRCK-02 | `tests/offline/tracking-domain-adapters.spec.ts` | 03-04 summary + unit output |
| TRCK-03 | `tests/client/tracking-history-sync-state.spec.ts` | 03-05 summary + unit output |
| OFFL-01 | `tests/offline/plan-cache-read.spec.ts` | 03-01/03-03 summaries + unit output |
| OFFL-02 | `tests/offline/queue-lifecycle.spec.ts` | 03-01/03-04 summaries + unit output |
| OFFL-03 | `tests/offline/replay-retry-policy.spec.ts` | 03-05 summary + unit output |

## Verification Gates

1. `npm run lint`
2. `npm run typecheck`
3. `npm run test:unit -- tests/offline/plan-cache-read.spec.ts tests/offline/queue-lifecycle.spec.ts`
4. `npm run test:unit -- tests/client/today-sync-summary.spec.ts tests/client/plan-flatten-map.spec.ts tests/client/plan-history-context.spec.ts`
5. `npm run test:unit -- tests/offline/food-queue-replay.spec.ts tests/offline/tracking-domain-adapters.spec.ts`
6. `npm run test:unit -- tests/client/tracking-history-sync-state.spec.ts tests/offline/replay-retry-policy.spec.ts`

All gates must pass before Phase 3 verification is marked complete.

## Risk Checks

- Assert queue items retain durable local_id and never lose replay metadata across app restarts.
- Assert sync state transitions (`queued`, `syncing`, `synced`, `failed`) are visible and consistent in UI.
- Assert replay logic does not silently drop failed items and supports manual retry.
- Assert staff surfaces do not inherit client offline write queue behavior.
