# Phase 3: Client Offline Daily Loop - Context

**Gathered:** 2026-04-23
**Status:** Ready for planning

<domain>
## Phase Boundary

Deliver the client daily loop for plan visibility and tracking with reliable offline queueing and sync-state transparency. This phase is client-surface only and does not expand offline support to nutritionist/admin surfaces.

</domain>

<decisions>
## Implementation Decisions

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

</decisions>

<specifics>
## Specific Ideas

- Prefer one-thumb entry controls for tracking actions.
- Keep Today screen action-oriented rather than dashboard-heavy.
- Ensure offline indicators are concise but always discoverable.

</specifics>

<canonical_refs>
## Canonical References

### Product and API contracts
- `docs/PRD.md` — Offline strategy, tracking expectations, and client/mobile product behavior.
- `docs/API.md` — Tracking endpoints, idempotent local_id behavior, sync endpoint, and response contracts.

### Planning and requirements
- `.planning/PROJECT.md` — frontend-only scope and offline boundary constraints.
- `.planning/REQUIREMENTS.md` — CLNT-01..03, TRCK-01..03, OFFL-01..03 requirements.
- `.planning/ROADMAP.md` — Phase 3 success criteria and sequencing.
- `.planning/STATE.md` — current completed baseline and continuity.

### Prior phase outputs
- `.planning/phases/01-platform-foundation/01-03-SUMMARY.md` — platform baseline references.
- `.planning/phases/02-authentication-access-control/02-04-SUMMARY.md` — auth/session/guard behavior to integrate with offline state lifecycle.
- `.planning/phases/02-authentication-access-control/02-UI-SPEC.md` — auth and shell UX continuity rules.

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- Role and auth middleware already enforce namespace access.
- Shared shell/layout primitives and notices from Phase 1 can host sync/offline state cues.
- Auth session store and plugins from Phase 2 provide logout hooks for cache/state cleanup.

### Established Patterns
- Strict role partitioning and Persian RTL visual baseline are already implemented.
- API error mapping and deterministic redirect behavior exist and should be reused.

### Integration Points
- New client offline repositories must integrate with current auth/session lifecycle.
- Tracking screens/components should attach to existing client layout and platform primitives.

</code_context>

<deferred>
## Deferred Ideas

- Nutritionist/admin offline parity.
- Advanced analytics or adherence storytelling beyond lightweight trend/status views.
- Realtime messaging improvements (kept for Phase 4+).

</deferred>

---

*Phase: 03-client-offline-daily-loop*
*Context gathered: 2026-04-23*