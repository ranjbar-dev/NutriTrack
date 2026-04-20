# Phase 6: Offline & PWA - Context

**Gathered:** 2026-04-20
**Status:** Ready for planning
**Mode:** --auto (sub-agent, no interactive session)

<domain>
## Phase Boundary

Phase 6 wraps the completed client-facing plan, tracking, and communication flows in an offline-capable PWA shell. Clients must be able to install the app, view their active diet plan and recent messages offline, submit tracking entries and outgoing messages while disconnected, and have those writes synchronize automatically when connectivity returns. This phase also adds Web Push subscriptions, reminder preferences, and server-side push delivery for messages, plan changes, food-request outcomes, meal reminders, medication reminders, and water reminders.

Nutritionist and super-admin experiences remain online-only. This phase does not redesign existing core flows or add new business capabilities outside offline resilience, sync orchestration, and notification delivery.

</domain>

<decisions>
## Implementation Decisions

### PWA shell and service worker strategy

- **D-01:** Integrate `@vite-pwa/nuxt` into the Nuxt app using a custom service worker (`injectManifest` strategy), not the zero-config preset, because Phase 6 needs explicit runtime caching, notification click handling, and Background Sync registration hooks.
- **D-02:** The PWA applies to the whole frontend build, but runtime offline behavior is client-role scoped. Static assets and route shell files can be cached globally; diet-plan, tracking, messaging, and sync features only activate inside client routes and client stores.
- **D-03:** Static assets (JS, CSS, fonts, icons, images) use cache-first with versioned precache. Client API reads for active plan and message history use network-first with IndexedDB fallback so online freshness wins when available.
- **D-04:** The manifest is Persian-only with `name`/`short_name` based on `نوتری‌ترک`, `display: standalone`, app shortcuts for `برنامه`, `ثبت روزانه`, and `پیام‌ها`, and install/update prompts surfaced from a shared client-only PWA status component.

### Offline data model and cache boundaries

- **D-05:** Use Dexie.js as the single offline database. Create versioned tables for `activePlan`, `messages`, `syncQueue`, `syncMeta`, `notificationPreferences`, and `uiState`; tracking rows are not stored in per-domain tables because the authoritative offline write source is the normalized queue payload.
- **D-06:** Cache exactly one active plan snapshot per client user, including plan days, meals, options, exercise recommendations, medications, water target, and fetch metadata (`fetched_at`, `plan_id`, `updated_hint`). Refresh the snapshot on client app open when online.
- **D-07:** Cache the last 50 messages per conversation in IndexedDB. Opening a conversation reads cached messages immediately, then fetches fresh data online and merges by message ID.
- **D-08:** Outgoing offline messages may include attachment blobs. Store queued message payloads and attachment metadata/blob references in IndexedDB so text and file messages share one sync path.

### Sync queue and reconnect behavior

- **D-09:** Wrap existing `useApi()` calls with an offline-aware client API layer instead of duplicating fetch logic in each store. For eligible client POST requests, network failures caused by offline/unreachable transport become queued sync entries rather than visible hard failures.
- **D-10:** Queue all client write operations from Phases 4 and 5 that already carry or can carry `local_id`: food logs, water logs, sleep logs, exercise logs, medication logs, body measurements, lab-result metadata/file uploads, and outgoing messages.
- **D-11:** `syncQueue` entries store: `id`, `entity_type`, `request_path`, `method`, serialized payload, optional attachment blob, `local_id`, `created_at`, `status`, `retry_count`, `last_error`, and `next_attempt_at`. Queue processing is FIFO and single-flight to preserve causality and simplify UI state.
- **D-12:** Retry policy is fixed at exponential backoff `1s → 2s → 4s`, max 3 attempts. After the third failure the item becomes `failed` and remains visible in the pending-sync UI until manual retry.
- **D-13:** Reconnect detection uses both browser `online` events and an app-open sync sweep. Register Background Sync where supported, but the app must still sync correctly without it via foreground timers/polling.
- **D-14:** Sync success updates in-memory stores immediately: queued tracking items reconcile into existing store arrays, queued messages replace temporary local echoes, and plan/message caches update after successful pulls. Server-side `local_id` deduplication from Phase 4 remains the canonical duplicate guard.

### Push notification backend and reminder scheduling

- **D-15:** Add backend persistence for `push_subscriptions` and `notification_preferences`. Each client can have multiple browser/device subscriptions; preferences are per client and per notification category.
- **D-16:** Expose authenticated client endpoints to subscribe/unsubscribe Web Push and to read/update notification preferences. Subscription writes must be idempotent by endpoint.
- **D-17:** Use `webpush-go` with VAPID keys loaded from backend config/env. All push payloads share one envelope: Persian `title`, `body`, `action_url`, `icon`, `type`, and optional entity IDs for deep-link navigation.
- **D-18:** Event-driven pushes for new messages, new active plans, and food-request review outcomes fire directly from the existing Phase 3/5 service layer after successful writes. Scheduled reminders for meals, medications, and water run from a lightweight backend goroutine ticker every minute.
- **D-19:** Reminder scheduling is preference-aware and duplicate-safe. The scheduler only targets active client plans and records a dedup key per reminder window so the same reminder is not sent twice within the same minute bucket.

### Client UX for offline state and preferences

- **D-20:** Add a persistent client sync-status surface in the mobile UI showing `همگام‌سازی در حال انجام`, `همه‌چیز همگام است`, or `X مورد در انتظار`. Failed items expose a manual retry action; offline mode shows a distinct banner/toast when writes are saved locally.
- **D-21:** The existing client plan and client message pages must render cached data first when available. If there is no cache and the device is offline, show a Persian empty/offline state that clearly says fresh data will appear after reconnect.
- **D-22:** Notification preferences live under the client profile/settings area, not in a separate admin-like page. Toggles cover: new message, new plan, food-request result, meal reminders, medication reminders, and water reminders.
- **D-23:** iOS storage eviction is handled defensively: on app open, if required Dexie tables are missing or empty unexpectedly, clear stale UI flags, show a Persian “offline data was removed by the device” notice, and re-fetch plan/messages immediately when connectivity is available.

### Agent's Discretion

- Exact install-prompt copy and iconography, as long as it stays Persian and mobile-first
- Whether the sync-status indicator sits in the client layout header or as a compact floating pill, as long as it remains visible from plan/tracking/messages screens
- Whether failed sync items are listed inline on each screen or centralized in one retry sheet, as long as manual retry remains obvious
- The exact Workbox runtime cache names and Dexie schema version number

</decisions>

<specifics>
## Specific Ideas

- Reuse the existing `local_id` generation already present in tracking stores instead of inventing a second offline identifier scheme.
- Treat offline writes as “save now, sync later” with immediate optimistic UI updates and a Persian confirmation toast rather than a scary error state.
- Keep the first offline milestone focused on the client’s core daily loop: open plan, check meals/medications, log activity, send message, reconnect later.

</specifics>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Product scope
- `.planning/ROADMAP.md` §Phase 6 — goal, dependency chain, and success criteria for offline, sync, and push behavior
- `.planning/REQUIREMENTS.md` §Offline & Sync, §Push Notifications, §UI/UX & PWA — OFFL-01 through OFFL-12, NOTIF-01 through NOTIF-08, UI-06 and UI-07
- `docs/phases.md` §Phase 6 — implementation guidance, validation checklist, and delivery expectations for service worker, IndexedDB, sync manager, and Web Push

### Existing phase contracts
- `.planning/phases/04-client-tracking-suite/04-CONTEXT.md` — tracking data model, `local_id` contract, and client logging flows that must work offline
- `.planning/phases/05-communication-collaboration/05-CONTEXT.md` — messaging/file-request behavior, polling cadence, and message attachment constraints that Phase 6 wraps with caching and queueing
- `.planning/phases/03-diet-plan-engine/03-CONTEXT.md` — active-plan structure, meal times, medication times, exercise recommendations, and water-target fields reused for offline plan cache and reminders
- `.planning/STATE.md` — current project position and Phase 6 blocker note about iOS PWA storage eviction

### Stack and platform constraints
- `.planning/research/STACK.md` — approved stack choices for Nuxt 4, `@vite-pwa/nuxt`, Dexie.js, `webpush-go`, local filesystem storage, and Hetzner deployment
- `frontend/package.json` — current frontend dependency baseline to extend for PWA/offline support
- `backend/cmd/api/main.go` — existing server bootstrapping and service wiring where push services and reminder scheduling must integrate

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `frontend/app/composables/useApi.ts` — the existing authenticated fetch wrapper is the correct seam for an offline-aware client API layer
- `frontend/app/stores/clientPlan.ts` — current active-plan load/store pattern to extend with IndexedDB-backed cached reads
- `frontend/app/stores/message.ts` — existing message polling/send flows to wrap with cached history, optimistic offline echo, and queued send
- `frontend/app/stores/foodLog.ts`, `waterLog.ts`, `sleepLog.ts`, `exerciseLog.ts`, `medicationLog.ts`, `bodyMeasurement.ts`, `labResult.ts` — each already creates `local_id`, which Phase 6 should preserve when queueing offline writes
- `frontend/app/layouts/client.vue` — best place to surface global client sync/install status without affecting nutritionist/admin layouts
- `backend/internal/service/tracking_service.go` — existing tracking service boundary and filesystem upload behavior for lab results
- `backend/internal/service/communication_service.go` — existing messaging/file-attachment workflow to extend with push triggers and offline-compatible send semantics

### Established Patterns
- Pinia setup stores per domain and keep error/loading state local to each store
- `useApi()` handles auth refresh centrally; new offline behavior should layer on top of it rather than bypassing authentication
- Backend follows handler → service → repository with Persian error responses and config-driven filesystem paths
- Client UX is mobile-first, Persian-only, and card-based; new PWA/install/sync surfaces must follow the same constraint

### Integration Points
- Client plan fetch endpoint `/api/clients/me/active-plan` is the canonical source for plan cache refresh and reminder schedule data
- Client message endpoints under `/api/messages/*` supply both cached conversation reads and queued outgoing message sync
- Tracking endpoints from Phase 4 already accept `local_id` and are the primary consumers of the offline sync queue
- Client profile/settings routes are the natural home for notification permissions and preference toggles
- `backend/cmd/api/main.go` is where subscription services, notification services, and reminder scheduler startup should be wired into the running API process

</code_context>

<deferred>
## Deferred Ideas

- Background push delivery for nutritionist/admin roles — out of scope because offline/PWA support is client-only
- Rich message search, adaptive polling, or inbox summarization — future enhancement beyond the offline wrapper
- Native-device reminder scheduling outside Web Push/browser capabilities — future/native-app concern, not Phase 6

</deferred>

---

*Phase: 06-offline-pwa*
*Context gathered: 2026-04-20*
