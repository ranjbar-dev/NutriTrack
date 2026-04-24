# Phase 3: Client Offline Daily Loop - Research

**Researched:** 2026-04-23
**Domain:** Nuxt 4 client offline queue + tracking UX + sync replay transparency
**Confidence:** HIGH

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions
### Client daily loop surfaces
- **D-01:** Client home prioritizes Today actions: active plan snapshot, pending tracking actions, water target, and sync status indicators.
- **D-02:** Active plan UI must flatten nested API structure into readable mobile sections while preserving day/meal/option semantics.
- **D-03:** Tracking entry flows optimize for low-friction repeat actions over form density.

### Offline queue and sync behavior
- **D-04:** Offline write support is required for client tracking domains only: food, water, sleep, exercise, medication, and body measurements.
- **D-05:** Every queued write carries local_id and explicit sync state (`queued`, `syncing`, `synced`, `failed`) visible in UX.
- **D-06:** Sync replay triggers on reconnect and app-open, with manual retry available for failed entries.
- **D-07:** Conflict handling follows last-write-wins per PRD, but UI must still surface retry/error states clearly.

### Data boundaries and role isolation
- **D-08:** Nutritionist and admin routes remain online-first in this phase; no queueing layer is added for staff flows.
- **D-09:** Client offline cache scope is limited to active plan, recent tracking records, and required lookups for tracking UX.
- **D-10:** Cached authenticated payloads remain tightly bounded and user-scoped; logout must clear client offline state.

### UX and trust signals
- **D-11:** Sync and offline states are first-class UI elements, not hidden technical behavior.
- **D-12:** Persian-only copy explains state and recovery steps for queue failures in plain language.

### the agent's Discretion
- Exact Dexie schema partitioning and repository adapters.
- Sync batching granularity for replay throughput.
- Detailed component composition as long as UI-SPEC and Phase 1/2 contracts remain intact.

### Deferred Ideas (OUT OF SCOPE)
- Nutritionist/admin offline parity.
- Advanced analytics or adherence storytelling beyond lightweight trend/status views.
- Realtime messaging improvements (kept for Phase 4+).
</user_constraints>

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| CLNT-01 | Today view with active plan, pending actions, water target, sync status | Today aggregate model + SyncStatusStrip + water quick actions [VERIFIED: .planning/REQUIREMENTS.md] |
| CLNT-02 | Full active plan readability | Plan cache/read model flattened by day->meal->option view-mappers [VERIFIED: .planning/REQUIREMENTS.md] |
| CLNT-03 | Archived plan visibility with active context retained | Separate active_plan vs plan_history stores with explicit active badge [VERIFIED: .planning/REQUIREMENTS.md] |
| TRCK-01 | Food logging with mobile-friendly flows | Per-domain tracking composer + local-first enqueue using local_id [VERIFIED: docs/API.md] |
| TRCK-02 | Water/sleep/exercise/medication/body logging | Unified queue envelope + domain adapters for six tracking endpoints [VERIFIED: docs/API.md] |
| TRCK-03 | Recent tracking history + lightweight progress | Query local recent entries first, merge synced state metadata, derive simple summaries [VERIFIED: .planning/REQUIREMENTS.md] |
| OFFL-01 | Offline read of active plan + essential recent data | IndexedDB scoped cache for active plan + recent tracking [CITED: docs/PRD.md#6-offline--sync-strategy] |
| OFFL-02 | Local queue with durable IDs + visible states | Queue table keyed by local_id and explicit sync_state enum [VERIFIED: docs/API.md] |
| OFFL-03 | State visibility + auto/manual retry | Reconnect/app-open replay, backoff retry, failed retry UI controls [CITED: docs/PRD.md#6-offline--sync-strategy] |
</phase_requirements>

## Project Constraints (from copilot-instructions.md)

- Frontend-only scope; do not propose backend/database/infrastructure changes. [VERIFIED: .github/copilot-instructions.md]
- Stack must remain Nuxt 4 + Tailwind CSS 4 + Pinia with strict TypeScript and composable-first patterns. [VERIFIED: .github/copilot-instructions.md]
- docs/API.md is the backend contract; docs/PRD.md is the product behavior contract. [VERIFIED: .github/copilot-instructions.md]
- Persian-only RTL mobile PWA quality is primary; desktop is secondary. [VERIFIED: .github/copilot-instructions.md]
- Client offline support is required; nutritionist/admin stay online-first unless explicitly changed. [VERIFIED: .github/copilot-instructions.md]
- Keep role boundaries explicit across routing/state/cache/persistence. [VERIFIED: .github/copilot-instructions.md]

## Summary

Phase 3 should implement a local-first tracking write pipeline with explicit per-entry sync states, while keeping read cache scope intentionally narrow to active plan + recent tracking + minimal lookup data. This directly matches locked decisions D-04..D-10 and requirement OFFL-01..03. [VERIFIED: 03-CONTEXT.md] [VERIFIED: .planning/REQUIREMENTS.md]

The backend contract already supports this model: all tracking writes accept a client-generated local_id idempotency key, and bulk sync is available at POST /tracking/sync with type-tagged entries. This allows deterministic replay without duplicate server records if local_id is stable. [VERIFIED: docs/API.md]

UX trust should be treated as product logic, not decoration: global sync strip + row chips + retry affordances must remain visible until state transitions occur. That aligns with Phase 3 UI contract and PRD offline strategy (retry with capped backoff, manual recovery after max retries). [VERIFIED: 03-UI-SPEC.md] [CITED: docs/PRD.md#6-offline--sync-strategy]

**Primary recommendation:** Use Dexie-backed queue and cache repositories, replay via reconnect/app-open scheduler using /tracking/sync first, and surface queue transparency in both aggregate and per-entry UI from a single source of truth. [CITED: /websites/dexie] [VERIFIED: docs/API.md]

## Architectural Responsibility Map

| Capability | Primary Tier | Secondary Tier | Rationale |
|------------|-------------|----------------|-----------|
| Tracking entry capture and local enqueue | Browser / Client | API / Backend | Immediate offline durability and low-latency UX must happen client-side; backend receives replayed payloads. [VERIFIED: docs/PRD.md] |
| Queue replay orchestration | Browser / Client | API / Backend | Trigger conditions are reconnect/app-open/manual retry in client runtime; backend provides idempotent processing endpoint. [VERIFIED: 03-CONTEXT.md] [VERIFIED: docs/API.md] |
| Sync-state transparency (global + row) | Browser / Client | -- | State is local queue lifecycle metadata and must remain visible in client UI. [VERIFIED: 03-UI-SPEC.md] |
| Conflict resolution policy | API /Backend | Browser / Client | Policy is server-authoritative (last-write-wins), client communicates outcome and allows retry/review. [CITED: docs/PRD.md#6-offline--sync-strategy] |
| Active plan and recent tracking cache | Browser / Client | API / Backend | Client owns offline read availability; backend remains source of fresh truth on reconnect. [VERIFIED: 03-CONTEXT.md] |
| Role isolation (client-only offline queue) | Browser / Client (routing/state gating) | -- | Existing role boundaries in middleware/stores must prevent queue behavior from leaking to staff surfaces. [VERIFIED: app/middleware/auth-access.global.ts] [VERIFIED: app/middleware/role-shell.global.ts] |

## Standard Stack

### Core
| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| Nuxt | 4.0.3 (in repo) | App/runtime shell and composables | Existing project baseline; avoids introducing parallel architecture. [VERIFIED: package.json] |
| Pinia | 3.0.4 (latest) | App state and UI-derived queue stats | Existing state layer with plugin/subscription support for lifecycle hooks. [VERIFIED: npm registry] [CITED: /vuejs/pinia] |
| Dexie | 4.4.2 (latest) | IndexedDB schema, transactions, migrations for queue/cache | Proven IndexedDB ergonomics; explicit versioned schema and migration API for offline data evolution. [VERIFIED: npm registry] [CITED: /websites/dexie] |
| @vite-pwa/nuxt | 1.1.1 | SW registration/runtime cache behavior in Nuxt | Already installed; supports Workbox runtime strategies and update lifecycle. [VERIFIED: npm registry] [CITED: /vite-pwa/vite-plugin-pwa] |

### Supporting
| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| @vueuse/core | 14.2.1 | Reactive connectivity/composable helpers | Use for online/offline listeners and reconnect triggers if preferred over custom listeners. [VERIFIED: npm registry] |
| @pinia/nuxt | 0.11.3 (latest) | Nuxt integration for Pinia | Keep as integration layer; optional patch upgrade from 0.11.2 after phase validation. [VERIFIED: npm registry] |

### Alternatives Considered
| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| Dexie | idb (8.0.3) | idb is smaller but lower-level; more custom boilerplate for queue indexes/migrations. [VERIFIED: npm registry] [ASSUMED] |
| Custom replay per endpoint | POST /tracking/sync batch-first replay | Single endpoint reduces request overhead and keeps replay status deterministic per batch. [VERIFIED: docs/API.md] |

**Installation:**
```bash
npm install dexie @vueuse/core
```

**Version verification:**
- dexie 4.4.2, modified 2026-04-16T06:48:16.709Z [VERIFIED: npm registry]
- @vueuse/core 14.2.1, modified 2026-02-10T03:36:45.608Z [VERIFIED: npm registry]
- pinia 3.0.4, modified 2025-11-05T09:25:14.059Z [VERIFIED: npm registry]
- @pinia/nuxt 0.11.3, modified 2025-11-05T09:25:18.224Z [VERIFIED: npm registry]
- @vite-pwa/nuxt 1.1.1, modified 2026-02-06T10:27:12.488Z [VERIFIED: npm registry]

## Architecture Patterns

### System Architecture Diagram

```text
[Tracking UI Actions]
    |
    v
[Client Tracking Composer] -- validates/maps --> [Queue Repository (Dexie)]
    |                                               |
    | immediate optimistic state                    | persisted entry:
    |                                               | local_id, domain, payload,
    |                                               | sync_state, retry_count, error
    v                                               v
[Today/History Read Models] <--- joins --- [Plan Cache + Recent Tracking Cache]
    |
    v
[SyncStatusStrip + SyncStateChip]
    |
    +--> trigger: app-open / reconnect / manual retry
             |
             v
      [Replay Scheduler]
             |
             +--> batch by type/order --> POST /tracking/sync
             |
             +--> success: mark synced + synced_at
             |
             +--> partial error: mark failed with reason
             |
             +--> retry policy (exp backoff, max 3)
```

### Recommended Project Structure
```text
app/
  composables/
    offline/
      useOfflineQueue.ts          # enqueue + dequeue + state transitions
      useSyncReplay.ts            # reconnect/app-open/manual replay
      useSyncStatus.ts            # aggregate counters and labels
  lib/
    offline/
      db.ts                       # Dexie schema + migrations
      queue-repository.ts         # CRUD by local_id/state/domain
      replay-policy.ts            # backoff + max retry rules
      tracking-mapper.ts          # endpoint payload adapters
  stores/
    client-offline.ts             # queue summary + UI state transparency
  plugins/
    sync-replay.client.ts         # app-open bootstrap + online event hook
```

### Pattern 1: Transactional Queue Writes
**What:** Write queue entry + optimistic local projection in one Dexie transaction to prevent split-brain local state. [CITED: /websites/dexie]
**When to use:** Every tracking create/update path.
**Example:**
```typescript
// Source: Dexie docs (Version.stores + transactions)
await db.transaction('rw', db.queueEntries, db.trackingRecent, async () => {
  await db.queueEntries.put({
    localId,
    domain: 'water',
    syncState: 'queued',
    retryCount: 0,
    payload,
    createdAt: nowIso
  })

  await db.trackingRecent.put({
    localId,
    domain: 'water',
    display,
    syncState: 'queued'
  })
})
```

### Pattern 2: Batch Replay via /tracking/sync with Row-Level Reconciliation
**What:** Replay pending entries using /tracking/sync, then reconcile result back to each local_id row.
**When to use:** App open, reconnect, and manual retry.
**Example:**
```typescript
// Source: docs/API.md (15.7 Bulk Sync)
const entries = queued.map((q) => ({ type: q.domain, payload: q.payload }))
const result = await authFetch.request('/tracking/sync', { method: 'POST', body: { entries } })

for (const item of queued) {
  const hasError = result.data.errors.some((e) => e.local_id === item.localId)
  await queueRepo.updateState(item.localId, hasError ? 'failed' : 'synced')
}
```

### Pattern 3: Single Source for Sync Transparency
**What:** Derive both global strip metrics and row chips from queue table state (not from transient network callbacks).
**When to use:** Always; prevents hidden mismatch between list rows and header counters.
**Example:**
```typescript
const summary = computed(() => ({
  queued: entries.value.filter((x) => x.syncState === 'queued').length,
  syncing: entries.value.filter((x) => x.syncState === 'syncing').length,
  failed: entries.value.filter((x) => x.syncState === 'failed').length
}))
```

### Anti-Patterns to Avoid
- **Per-screen queue state copies:** Multiple local caches drift and break transparency; keep one queue source and derive views. [ASSUMED]
- **local_id regeneration on retry:** Destroys idempotency guarantees and can create duplicate server records. [VERIFIED: docs/API.md]
- **Queue execution on nutritionist/admin namespaces:** Violates locked D-08 scope and role-isolation requirements. [VERIFIED: 03-CONTEXT.md]
- **Toast-only failure communication:** Violates UI contract requiring persistent failed visibility and recovery actions. [VERIFIED: 03-UI-SPEC.md]

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| IndexedDB driver and migration layer | custom wrappers over native indexedDB API | Dexie versioned schema + upgrade hooks | Native API complexity and migration edge cases are costly and error-prone. [CITED: /websites/dexie] |
| SW registration/update lifecycle | manual service-worker registration orchestration | @vite-pwa/nuxt + existing pwa.client.ts | Already integrated in repo; keeps update/offline signaling consistent. [VERIFIED: app/plugins/pwa.client.ts] [CITED: /vite-pwa/vite-plugin-pwa] |
| State instrumentation around actions | custom event bus for store lifecycle | Pinia $subscribe/$onAction and plugins | Built-in hooks provide deterministic state/action tracing and cleanup behavior. [CITED: /vuejs/pinia] |
| Retry/backoff timing glue | ad hoc setTimeout chains in components | Central replay-policy module + queue metadata fields | Keeps retry behavior testable and consistent across all tracking domains. [ASSUMED] |

**Key insight:** The expensive failures in offline tracking systems are consistency failures, not UI failures; use battle-tested storage/state primitives and keep state transitions centralized. [ASSUMED]

## Common Pitfalls

### Pitfall 1: Queue/UI Divergence
**What goes wrong:** Header says synced while row chips remain queued/failed.
**Why it happens:** Aggregate computed from network events, rows from local DB snapshots.
**How to avoid:** Compute both from the same queue table projection.
**Warning signs:** Mismatched counters after reconnect replay.

### Pitfall 2: Losing Offline Data on Logout/Role Switch
**What goes wrong:** Next user sees stale client queue/cache or stale auth-scoped data persists.
**Why it happens:** Auth cleanup does not clear offline stores.
**How to avoid:** Add client store cleaner to logoutWithCleanup role cleaner map.
**Warning signs:** Entries remain after manual logout and relogin as different user.

### Pitfall 3: Infinite Replay Loops
**What goes wrong:** Failed entries loop forever and degrade battery/network.
**Why it happens:** Missing retry_count + terminal failed state.
**How to avoid:** Cap retries at 3 then switch to failed with explicit manual retry.
**Warning signs:** Repeated POST /tracking/sync with same payload after max retry.

### Pitfall 4: Over-caching Sensitive Payloads
**What goes wrong:** Offline DB contains broad authenticated data outside phase boundary.
**Why it happens:** Generic API response caching without scope filter.
**How to avoid:** Restrict offline stores to active plan + recent tracking + required lookups.
**Warning signs:** IndexedDB stores unrelated role/admin payloads.

## Code Examples

Verified patterns from official/local contracts:

### Dexie schema with state indexes
```typescript
// Source: Dexie docs + Phase 3 constraints
db.version(1).stores({
  queueEntries: 'localId, syncState, domain, createdAt, retryCount',
  trackingRecent: 'localId, domain, loggedDate, syncState',
  activePlanCache: 'planId, updatedAt'
})
```

### Reconnect + app-open replay trigger
```typescript
// Source: PRD offline strategy + D-06
onNuxtReady(() => {
  replayPending('app-open')
})

window.addEventListener('online', () => {
  replayPending('reconnect')
})
```

### UI state mapping contract
```typescript
// Source: 03-UI-SPEC.md Offline State UX Contract
const stateLabel: Record<'queued'|'syncing'|'synced'|'failed', string> = {
  queued: 'queued',
  syncing: 'syncing',
  synced: 'synced',
  failed: 'failed'
}
```

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| localStorage offline drafts + fire-and-forget replay | IndexedDB queue with per-entry lifecycle states and replay policy | Mature PWA offline guidance era (post-Workbox/IDB standardization) [ASSUMED] | Better durability, observability, and conflict handling |
| endpoint-by-endpoint immediate retries | batch replay with deterministic queue metadata | Supported by tracking sync endpoint in API contract [VERIFIED: docs/API.md] | Lower request overhead and cleaner failure accounting |

**Deprecated/outdated:**
- Hiding offline failures in temporary toast messages only; replaced by persistent state chips and recoverable queue views. [VERIFIED: 03-UI-SPEC.md]

## Assumptions Log

| # | Claim | Section | Risk if Wrong |
|---|-------|---------|---------------|
| A1 | idb would require materially more custom queue/migration boilerplate than Dexie in this codebase | Standard Stack / Alternatives | Medium: could overstate Dexie benefit and affect dependency choice |
| A2 | Central replay-policy module is the best fit over composable-local retry logic for this team | Don't Hand-Roll | Low-Medium: mainly maintainability tradeoff |
| A3 | Post-Workbox era is the key timeline shift for modern offline queue patterns | State of the Art | Low: historical framing only |

## Open Questions (RESOLVED)

1. Should Phase 3 include client message offline queue now or defer strictly to Phase 4 boundary?
- **Decision:** defer message offline queue implementation to Phase 4 and keep Phase 3 queue behavior tracking-domain only.
- **Reason:** preserves phase boundary while still allowing future extensibility through domain-tagged queue schema.

2. Should pinia/@pinia/nuxt patch upgrades be done inside Phase 3?
- **Decision:** defer package upgrades unless a required Phase 3 implementation API is blocked.
- **Reason:** avoids maintenance-risk scope creep inside a feature-critical offline phase.

## Environment Availability

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| node | Nuxt build/tests and tooling | yes | available in session (npm commands succeeded) | -- |
| npm | package/version verification and installs | yes | available in session | -- |
| npx | Context7 CLI fallback docs lookup | yes | available in session | Use official docs manually if unavailable |
| Browser Background Sync API | optional background replay behavior | unknown (runtime capability) | -- | reconnect/app-open/manual retry loop (required baseline) |

**Missing dependencies with no fallback:**
- None identified. [VERIFIED: command execution in session]

**Missing dependencies with fallback:**
- Background Sync API may be unsupported on some browsers; fallback is app-open + reconnect replay plus manual retry. [CITED: docs/PRD.md#6-offline--sync-strategy]

## Validation Architecture

### Test Framework
| Property | Value |
|----------|-------|
| Framework | Vitest 3.2.2 (unit/integration-style store/composable tests), Playwright 1.54.1 (E2E/mobile) [VERIFIED: package.json] |
| Config file | vitest.config.ts, playwright.config.ts [VERIFIED: repository] |
| Quick run command | `npm run test:unit -- tests/platform/cache-boundary.spec.ts` |
| Full suite command | `npm run test:unit && npm run test:e2e` |

### Phase Requirements -> Test Map
| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| CLNT-01 | Today aggregate cards + sync strip | unit/component | `npm run test:unit -- tests/client/today-sync-summary.spec.ts` | no (Wave 0) |
| CLNT-02 | Active plan flatten mapping correctness | unit | `npm run test:unit -- tests/client/plan-flatten-map.spec.ts` | no (Wave 0) |
| CLNT-03 | Active vs archived context separation | unit | `npm run test:unit -- tests/client/plan-history-context.spec.ts` | no (Wave 0) |
| TRCK-01 | Food log enqueue + replay mapping | unit | `npm run test:unit -- tests/offline/food-queue-replay.spec.ts` | no (Wave 0) |
| TRCK-02 | Other five tracking domains enqueue/replay | unit | `npm run test:unit -- tests/offline/tracking-domain-adapters.spec.ts` | no (Wave 0) |
| TRCK-03 | Recent history merge with sync states | unit | `npm run test:unit -- tests/client/tracking-history-sync-state.spec.ts` | no (Wave 0) |
| OFFL-01 | Offline read from active cache | unit | `npm run test:unit -- tests/offline/plan-cache-read.spec.ts` | no (Wave 0) |
| OFFL-02 | Durable local_id queue lifecycle | unit | `npm run test:unit -- tests/offline/queue-lifecycle.spec.ts` | no (Wave 0) |
| OFFL-03 | Reconnect retry + manual retry transitions | unit | `npm run test:unit -- tests/offline/replay-retry-policy.spec.ts` | no (Wave 0) |

### Sampling Rate
- **Per task commit:** `npm run test:unit -- tests/offline/*.spec.ts`
- **Per wave merge:** `npm run test:unit`
- **Phase gate:** `npm run test:unit && npm run test:e2e` green before verify-work

### Wave 0 Gaps
- [ ] tests/offline/queue-lifecycle.spec.ts - covers OFFL-02
- [ ] tests/offline/replay-retry-policy.spec.ts - covers OFFL-03
- [ ] tests/offline/plan-cache-read.spec.ts - covers OFFL-01
- [ ] tests/client/today-sync-summary.spec.ts - covers CLNT-01
- [ ] tests/client/plan-flatten-map.spec.ts - covers CLNT-02
- [ ] tests/client/plan-history-context.spec.ts - covers CLNT-03
- [ ] tests/client/tracking-history-sync-state.spec.ts - covers TRCK-03

## Security Domain

### Applicable ASVS Categories

| ASVS Category | Applies | Standard Control |
|---------------|---------|-----------------|
| V2 Authentication | yes | Reuse existing auth-session and auth-fetch token flow; no token storage in offline queue tables. [VERIFIED: app/stores/auth-session.ts] |
| V3 Session Management | yes | Logout cleanup must clear client offline stores and role cookies/state. [VERIFIED: app/stores/auth-session.ts] |
| V4 Access Control | yes | Keep offline queue and cache under client route/store scope only. [VERIFIED: 03-CONTEXT.md] |
| V5 Input Validation | yes | Validate tracking payload shape before enqueue and before replay mapping. [ASSUMED] |
| V6 Cryptography | no (phase scope) | No new crypto primitives; rely on HTTPS/JWT backend contract. [VERIFIED: docs/API.md] |

### Known Threat Patterns for Nuxt client offline queue

| Pattern | STRIDE | Standard Mitigation |
|---------|--------|---------------------|
| Stale data exposure across user switch | Information Disclosure | Clear client-scoped IndexedDB and in-memory summaries on logout/role change. |
| Offline payload tampering in local storage | Tampering | Re-validate payload shape and enforce per-endpoint required fields before replay. |
| Replay flood from failed loops | Denial of Service | Retry cap + exponential backoff + manual retry gate. [CITED: docs/PRD.md#6-offline--sync-strategy] |
| Unauthorized domain queueing (staff routes) | Elevation of Privilege | Route/store guard by role and namespace; no offline plugin registration outside client scope. |

## Sources

### Primary (HIGH confidence)
- docs/API.md - tracking endpoint contract, local_id idempotency, bulk sync endpoint.
- docs/PRD.md - offline strategy, client-only scope, retry policy.
- .planning/phases/03-client-offline-daily-loop/03-CONTEXT.md - locked decisions and scope boundaries.
- .planning/phases/03-client-offline-daily-loop/03-UI-SPEC.md - required sync-state UX visibility contract.
- /websites/dexie (Context7) - schema versioning, stores definition, upgrade patterns.
- /vite-pwa/vite-plugin-pwa (Context7) - runtime caching and registration/update behavior.
- /vuejs/pinia (Context7) - plugin/subscription/action instrumentation patterns.
- npm registry (`npm view`) - current package versions and modified timestamps.

### Secondary (MEDIUM confidence)
- package.json + current app store/plugin code for integration posture.

### Tertiary (LOW confidence)
- None.

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH - package versions verified and doc APIs fetched in session.
- Architecture: HIGH - directly constrained by phase context + API + UI spec.
- Pitfalls: MEDIUM-HIGH - mostly validated by phase contracts and established offline patterns.

**Research date:** 2026-04-23
**Valid until:** 2026-05-23 (30 days)
