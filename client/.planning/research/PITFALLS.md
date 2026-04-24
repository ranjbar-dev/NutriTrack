# Domain Pitfalls

**Domain:** Mobile-first Persian RTL nutrition management PWA
**Project:** NutriTrack client
**Researched:** 2026-04-22
**Overall confidence:** HIGH for product-specific pitfalls, MEDIUM for browser-behavior nuances

## Recommended Future Phases

Use these phase buckets when planning the frontend roadmap. The pitfall entries below refer to them directly.

| Phase | Focus |
|-------|-------|
| Phase 1 | App Foundation, PWA Shell, RTL Design System |
| Phase 2 | Auth, Session Management, Role Gating |
| Phase 3 | Offline Data Model, Cache Policy, Sync Engine |
| Phase 4 | Client Tracking Flows |
| Phase 5 | Chat, Polling, Uploads, and Media UX |
| Phase 6 | Push Notifications and Reminder Preferences |
| Phase 7 | Nutritionist and Admin Operational Surfaces |
| Phase 8 | Localization, Accessibility, Performance, and QA Hardening |

## Critical Pitfalls

### Pitfall 1: Treating offline mode as a cache feature instead of a product contract

**Address in:** Phase 1 and Phase 3

**What goes wrong:**
Teams cache a few API responses, add an offline banner, and assume the product is now offline-capable. For NutriTrack, that fails because offline behavior is not generic. Different entities have different guarantees: active diet plans and cached messages are readable offline, tracking writes and outgoing messages must queue for later sync, and lab uploads explicitly do not work offline.

**Why it happens:**
Frontend teams often start from service-worker setup instead of a feature-by-feature offline matrix. That produces a technical PWA shell without a product-accurate offline model.

**Consequences:**
- Users think entries were saved when they were only stored in memory.
- Nutrition workflows diverge from the PRD promises.
- Upload flows fail silently offline.
- Roadmap churn appears later because offline logic gets rebuilt per feature.

**Warning signs:**
- The app has one generic `isOffline` branch but no per-feature behavior matrix.
- Offline messaging and offline tracking use different ad hoc local formats.
- Lab result upload appears available even when there is no network.
- Acceptance criteria say "works offline" without specifying read, queue, retry, and reconcile semantics.

**Prevention strategy:**
- Define an explicit offline capability matrix before UI implementation and keep it aligned with the PRD.
- Model three states for every mutable feature: local draft, queued for sync, synced.
- Separate resource caching from user data persistence: Cache Storage for static/network resources, IndexedDB for structured user data and outbox state.
- Block or clearly degrade unsupported offline actions such as lab uploads.
- Make offline state visible in the UI at the record level, not only at the global shell level.

**Detection:**
- QA can create tracking entries offline, reload the app, and still see durable queued items with status.
- QA can attempt a lab upload offline and receive a clear, actionable message instead of a broken form.

### Pitfall 2: Assuming Background Sync will make offline delivery reliable on all target devices

**Address in:** Phase 3

**What goes wrong:**
Teams build the sync engine around Background Sync and discover late that it is not available on Safari or Firefox-class environments. The result is an outbox that only flushes reliably on some browsers.

**Why it happens:**
Background Sync is attractive in demos, but current platform support is limited. MDN still marks the API as limited availability, with no support in Safari and Firefox.

**Consequences:**
- Queued tracking entries appear stuck on part of the audience.
- Message outbox behavior differs by browser.
- Product owners think sync is automatic when in practice it depends on reopening the app.

**Warning signs:**
- The sync design depends on `registration.sync.register()` as the primary path.
- There is no manual retry or app-open reconciliation flow.
- Browser support is not part of acceptance criteria for offline sync.

**Prevention strategy:**
- Treat Background Sync as an enhancement only.
- Build the primary sync path around app open, explicit reconnect detection, and a visible retry action.
- Store retry metadata per queued item: attempts, last error, next retry time, and user-visible status.
- Make sync idempotent on the client side by reconciling local queue items against successful server responses.
- Test the outbox behavior on Chromium and Safari-class browsers before Phase 3 exits.

**Detection:**
- Disable network, create several entries, close and reopen the app, reconnect, and verify sync without relying on background events.

### Pitfall 3: Shipping optimistic writes without a durable reconciliation model

**Address in:** Phase 3 and Phase 4

**What goes wrong:**
The UI shows a tracking entry or message immediately, but after reconnect the same item duplicates, disappears, or comes back reordered because the client lacks a durable local identity and reconciliation ledger.

**Why it happens:**
Teams focus on optimistic UX but postpone sync bookkeeping. In this product, that is dangerous because the PRD already expects local identifiers and last-write-wins behavior, and the API exposes bulk sync for offline replay.

**Consequences:**
- Double food or water logs.
- Messages shown as sent locally but missing on the server.
- Manual support burden because users cannot tell whether a record truly synced.

**Warning signs:**
- Queue items do not have stable client-generated UUIDs.
- There is no mapping between local records and confirmed server records.
- Message lists sort only by client clock.
- The UI cannot distinguish `queued`, `syncing`, `failed`, and `synced`.

**Prevention strategy:**
- Give every offline-capable mutation a durable client-side ID at creation time.
- Persist an outbox ledger with status transitions and server correlation fields.
- Reconcile using server timestamps and explicit sync outcomes, not just array replacement.
- Normalize timeline rendering rules so local optimistic items and server-confirmed items can coexist safely.
- Write slice tests for duplicate prevention and conflict handling before broad tracking rollout.

**Detection:**
- Repeat the same reconnect scenario with intermittent failures and confirm no duplicate or vanished entries remain.

### Pitfall 4: Letting token refresh, polling, and role changes fight each other

**Address in:** Phase 2 and Phase 5

**What goes wrong:**
Multiple requests hit expired access tokens at once, each tries to refresh independently, and the session collapses into 401 loops, race conditions, or forced logout. This gets worse when chat polling is active.

**Why it happens:**
JWT refresh is treated as a simple interceptor concern instead of a shared session state machine. In NutriTrack, clients also have OTP login, nutritionists/admins use password login, and every role boundary changes which routes and stores should be live.

**Consequences:**
- Chat screen flickers between loaded and logged-out states.
- Background requests overwrite valid tokens with stale ones.
- Role-specific stores leak data across logout/login transitions on shared devices.

**Warning signs:**
- More than one refresh request can run concurrently.
- Polling continues during refresh or after logout.
- Role is inferred only from route naming instead of the auth payload.
- Pinia stores survive account switching without explicit reset.

**Prevention strategy:**
- Implement a single-flight refresh queue so only one refresh request runs at a time.
- Pause polling and other background fetchers during refresh, logout, and role transitions.
- Reset all role-scoped stores on logout and on successful login under a different role.
- Build auth around an explicit session machine: anonymous, authenticating, authenticated, refreshing, expired.
- Test concurrent expiry by forcing multiple protected requests at once.

**Detection:**
- Simulate token expiry while the unread count and chat conversation are polling and verify one refresh path, no duplicate refresh calls, and clean recovery.

### Pitfall 5: Treating role-based flows as hidden navigation instead of isolated product surfaces

**Address in:** Phase 2 and Phase 7

**What goes wrong:**
Teams build one generic shell, toggle menu items by role, and accidentally reuse stores, cache keys, and route assumptions across client, nutritionist, and admin experiences.

**Why it happens:**
It feels faster early on, but NutriTrack has materially different capabilities, route graphs, and data safety requirements per role. Client offline support is in scope; nutritionist/admin offline is not.

**Consequences:**
- Wrong-role data appears after account switching.
- A nutritionist-only endpoint gets prefetched for a client route.
- Offline caches contain data that should never have been stored for another role.

**Warning signs:**
- Shared query keys do not encode role.
- Layouts, stores, and local persistence are global rather than role-scoped.
- Route guards are added late, after screens already exist.

**Prevention strategy:**
- Partition the app by role at the routing, layout, store, and persistence layers.
- Namespace cache keys, IndexedDB stores, and persisted settings by role and user ID.
- Treat client offline storage as a distinct subsystem; do not extend it implicitly to operations roles.
- Make route protection and store reset behavior part of the initial architecture, not polish.

**Detection:**
- Log in as client, then nutritionist, then client again on one device and verify no stale role data remains in UI or persistence.

### Pitfall 6: Adding service-worker caching that leaks or freezes authenticated data

**Address in:** Phase 1 and Phase 3

**What goes wrong:**
The service worker caches authenticated API responses too broadly. On a shared mobile device, a later user sees stale private diet plans, messages, or attachments; or the app keeps serving frozen responses after the backend changed.

**Why it happens:**
Teams apply generic stale-while-revalidate patterns to everything behind `/api` and `/uploads` without considering user identity, role, or sensitive payloads.

**Consequences:**
- Privacy leakage across sessions.
- Client sees an old plan while believing it is current.
- Debugging becomes difficult because network and cache disagree.

**Warning signs:**
- The cache strategy is defined by URL prefix only.
- There is no explicit list of cacheable authenticated resources.
- Logout does not clear user data caches or IndexedDB.

**Prevention strategy:**
- Cache static assets aggressively, but whitelist authenticated data explicitly.
- Store sensitive structured data in IndexedDB under user-scoped namespaces instead of opaque response caches when practical.
- Clear user-scoped caches and local databases on logout and account switch.
- Define freshness rules per resource: active plan, messages, notifications, and food catalogue subset all have different acceptable staleness.
- Avoid caching upload/download endpoints by default.

**Detection:**
- On a shared device test, logout and login as another role/user, then verify no prior authenticated data is shown when offline.

## Moderate Pitfalls

### Pitfall 7: Treating chat polling like a desktop convenience instead of a battery-sensitive mobile workflow

**Address in:** Phase 5

**What goes wrong:**
Polling runs continuously, unread counts and conversations fetch redundantly, and message ordering or read state drifts under flaky mobile connectivity.

**Prevention strategy:**
- Poll only while the relevant chat surface is visible.
- Use separate cadence rules for open conversation, unread badge, and background app states.
- Merge server pages with optimistic local messages through a stable ordering rule.
- Stop polling when offline, backgrounded, or refreshing auth.

**Warning signs:**
- Battery drain, repeated identical requests, or unread counts that oscillate.
- Messages jump position when the server response returns.

### Pitfall 8: Trying to make file uploads behave like offline-safe JSON mutations

**Address in:** Phase 5

**What goes wrong:**
Teams queue lab uploads or chat attachments offline, or they omit strong file constraints and progress/error handling because the API shape looks simple.

**Prevention strategy:**
- Keep uploads explicitly online-only unless the product contract changes.
- Validate file size and type before network submission and show progress, retry, and failure states.
- Generate attachment previews without assuming the final uploaded URL exists yet.
- Keep upload flows separate from the generic JSON outbox.

**Warning signs:**
- Attachments are pushed into the same offline queue as tracking entries.
- Large files fail without a clear client-side error.
- Users can leave the screen with no clue whether the upload finished.

### Pitfall 9: Treating Persian RTL as a stylesheet toggle instead of a full interaction requirement

**Address in:** Phase 1 and Phase 8

**What goes wrong:**
Layouts technically render right-to-left, but number formatting, Jalali dates, form fields, charts, swipe directions, and iconography still feel LTR or Western-first.

**Prevention strategy:**
- Define RTL and Persian formatting rules in the design system from day one.
- Standardize numeral rendering, Jalali date handling, and bidirectional text behavior for mixed Persian/Latin content.
- Test all core flows on actual narrow mobile widths, especially forms and charts.
- Audit component libraries before adoption; many break subtly in RTL.

**Warning signs:**
- Mixed Persian and numeric strings wrap badly.
- Time pickers, chart axes, or input placeholders read awkwardly.
- Teams mention "fix RTL later" during foundation work.

### Pitfall 10: Prompting for notifications too early and never recovering from denial

**Address in:** Phase 6

**What goes wrong:**
The app asks for notification permission immediately on first load or first login, users deny it, and the product loses its reminder and messaging reach permanently for many users.

**Prevention strategy:**
- Ask from a clear user gesture after the value is obvious, even if subscription setup happens during first-login onboarding.
- Support the full permission state model: default, granted, denied.
- Use service-worker notifications on mobile, not the page-level constructor.
- Deduplicate notifications with tags and clear stale ones when the user is already in the app.

**Warning signs:**
- Permission prompt appears before the user reaches chat, reminders, or settings.
- There is no denied-state UX explaining how to re-enable notifications later.
- The app fires multiple notifications for one chat burst.

## Minor Pitfalls

### Pitfall 11: Ignoring storage pressure and local schema migration until production

**Address in:** Phase 3 and Phase 8

**What goes wrong:**
IndexedDB schemas change mid-project, old outbox data becomes unreadable, or cached data gets evicted under storage pressure and the app cannot explain what happened.

**Prevention strategy:**
- Version local schemas intentionally.
- Add migration tests for offline stores.
- Track storage usage and request persistent storage where supported.

**Warning signs:**
- Local schema changes are done ad hoc.
- Test plans do not include upgrade-from-old-build scenarios.

### Pitfall 12: Shipping a PWA shell that meets installability but misses mobile recovery UX

**Address in:** Phase 1 and Phase 8

**What goes wrong:**
The app is technically installable, but users do not understand connectivity state, queued work, update availability, or how to recover from a failed sync.

**Prevention strategy:**
- Design explicit mobile recovery patterns: offline banner, item-level sync states, retry actions, and update prompts.
- Treat installability and service-worker registration as baseline plumbing, not success criteria.

**Warning signs:**
- The app has a manifest and service worker but no user-facing sync or update UX.

## Phase-Specific Warnings

| Phase Topic | Likely Pitfall | Mitigation |
|-------------|---------------|------------|
| App shell and PWA setup | Overcaching authenticated data too early | Freeze cache scope to static assets first, then add user data caching through an explicit matrix |
| Auth and role architecture | Shared stores across roles and refresh races | Build a session state machine, single-flight refresh, and store reset policy before feature screens |
| Offline engine | Browser-specific sync assumptions | Make app-open reconciliation the primary path and treat Background Sync as optional |
| Tracking flows | Duplicate or ghost entries after reconnect | Use durable client IDs, outbox statuses, and reconciliation tests |
| Chat and media | Battery drain and ambiguous upload state | Visibility-aware polling, online-only uploads, clear attachment progress and failure UX |
| Notifications | Early permission denial and noisy delivery | Ask from a user gesture, support denied state, and dedupe via notification tags |
| Operations surfaces | Client-only offline patterns creeping into staff tools | Keep client offline storage separate from nutritionist/admin surfaces |
| Hardening | RTL, Jalali, and narrow-screen bugs discovered late | Run real-device Persian RTL QA before release candidates |

## Roadmap Guidance

The roadmap should treat offline correctness, auth/session behavior, and role isolation as foundation work, not feature polish. For this product class, the expensive rewrites usually come from discovering too late that the offline queue, service worker policy, or auth refresh model is incompatible with chat, uploads, and multi-role routing.

The safest sequence is:

1. Establish the app shell, RTL design tokens, route partitions, and cache boundaries.
2. Build auth/session and role isolation before broad feature implementation.
3. Implement the offline data model and sync engine before tracking and chat surfaces depend on it.
4. Add tracking flows on top of the proven outbox and reconciliation rules.
5. Add chat, uploads, and notifications only after auth, visibility, and retry behavior are stable.

## Sources

- Project scope: `.planning/PROJECT.md`
- Product contract: `docs/PRD.md`
- API contract: `docs/API.md`
- MDN Background Synchronization API: https://developer.mozilla.org/en-US/docs/Web/API/Background_Synchronization_API
- MDN Push API: https://developer.mozilla.org/en-US/docs/Web/API/Push_API
- MDN Using the Notifications API: https://developer.mozilla.org/en-US/docs/Web/API/Notifications_API/Using_the_Notifications_API
- web.dev Offline data: https://web.dev/learn/pwa/offline-data