# Architecture Patterns

**Domain:** Frontend-only nutrition management PWA
**Project:** NutriTrack Client
**Researched:** 2026-04-22
**Overall confidence:** HIGH

## Recommended Architecture

Structure the frontend as a single Nuxt 4 application with four explicit layers:

1. **Platform layer** for app shell, auth/session, API client, storage, sync orchestration, notifications, RTL/Jalali utilities, and design tokens.
2. **Area layer** for role-specific route spaces: public/auth, client, nutritionist, and super admin.
3. **Domain feature layer** for bounded business modules such as plans, tracking, messages, foods, medications, clients, food requests, and notifications.
4. **Presentation layer** for reusable UI primitives plus feature-scoped components.

The key architectural decision is to make **client offline capability a bounded subsystem**, not a cross-app concern. Only the client area should depend on IndexedDB-backed caches and the offline sync queue. Nutritionist and super admin flows should stay online-first and simpler.

### Why this shape fits NutriTrack

- The PRD makes offline support mandatory only for clients.
- The API contract already separates role permissions cleanly, which maps well to route areas and domain repositories.
- Diet plans, tracking, messages, and notification preferences can be cached locally without forcing the entire app into an offline-first architecture.
- The frontend is greenfield, so separating platform concerns now avoids later rewrites when nutritionist/admin modules grow.

## System Shape

```text
Nuxt 4 App
|
+-- app/pages/auth/*                  Public auth flows
+-- app/pages/client/*                Offline-capable client workspace
+-- app/pages/nutritionist/*          Online-first nutritionist workspace
+-- app/pages/admin/*                 Online-first admin workspace
|
+-- app/layouts/*                     Area shells
+-- app/middleware/*                  Auth, role, online requirements
+-- app/plugins/*                     API, Pinia, PWA, push bootstrapping
+-- app/composables/*                 UI-facing feature facades
+-- app/stores/*                      Session, UI, queue, feature read models
+-- app/components/*                  Design system + feature components
|
+-- lib/api/*                         ofetch client, auth refresh, repositories
+-- lib/offline/*                     Dexie schema, queue engine, sync workers
+-- lib/domain/*                      DTO mapping, validators, normalizers
+-- lib/design/*                      tokens, theme contracts, icons, typography
+-- lib/i18n/*                        Persian text, numerals, Jalali helpers, RTL
+-- lib/utils/*                       Shared pure helpers
```

## Route and Area Boundaries

Use route prefixes as the main security and ownership boundary.

### Public and Auth Area

**Routes:** `app/pages/auth/*`

Responsibilities:

- Admin and nutritionist email/password login
- Client mobile number and OTP flow
- Refresh-session bootstrap and logout handoff
- App install education and notification consent education

Rules:

- No domain business logic here beyond auth.
- No direct access to client cached data.
- Successful login resolves the user role via `GET /auth/me`, then redirects to the correct area root.

### Client Area

**Routes:** `app/pages/client/*`

Responsibilities:

- Dashboard, active plan, daily tracking, measurements, water, sleep, exercise, medication, lab results, messages, food requests, settings
- Offline read access to cached plan, cached messages, cached notification preferences, and queued user actions
- Sync status visibility and retry UX

Rules:

- This is the only area allowed to depend on the offline queue and Dexie-backed repositories.
- Writes should go through repository commands that can choose remote write-through or queued offline persistence.
- File uploads remain online-required even in the client area.

### Nutritionist Area

**Routes:** `app/pages/nutritionist/*`

Responsibilities:

- Client list, client detail, plans, tracking review, messages, food requests, shared food/medication CRUD

Rules:

- Online-first only.
- Use the same repository contracts as client modules where possible, but remote-only implementations.
- No dependency on background sync or durable write queue.

### Super Admin Area

**Routes:** `app/pages/admin/*`

Responsibilities:

- Platform stats, nutritionist management, shared food/medication governance

Rules:

- Online-first only.
- Separate shell and navigation from nutritionist area, but reuse the same design system and shared table/filter components.

## Component Boundaries

| Component | Responsibility | Communicates With |
|-----------|---------------|-------------------|
| `AppShell` | App bootstrap, safe-area/mobile shell, top-level suspense/error boundaries | Session store, PWA state, route middleware |
| `AreaLayout*` | Role-specific navigation and chrome | Session store, unread count store |
| `FeaturePage` components | Route-level orchestration and data loading | Feature composables |
| `Feature composables` | UI-facing commands and queries, no raw HTTP | Stores, repositories |
| `Pinia stores` | Reactive read models, UI state, session state, queue state | Repositories, sync engine |
| `Repositories` | Translate feature operations to API/cache/queue behavior | API client, IndexedDB, mappers |
| `API client` | Auth headers, refresh flow, error normalization, upload helpers | ofetch, session store |
| `Offline engine` | Queue persistence, retries, connectivity reactions, bulk sync | Dexie, repositories, network status |
| `Service worker` | Static/runtime caching and PWA lifecycle only | PWA plugin, browser cache |
| `Design system` | Tokens, primitives, form controls, feedback patterns | Tailwind theme, feature components |
| `Locale/date helpers` | Persian numerals, Jalali formatting, text normalization | Feature components, mappers |

## Data Ownership Model

Use this ownership rule to keep state clean:

- **Pinia owns reactive application state** used by the current session and UI.
- **IndexedDB owns durable offline data** such as active plan cache, cached messages, pending tracking entries, queued outgoing messages, and sync metadata.
- **Repositories own data access decisions** and are the only layer that can choose between remote, cache, and queue.
- **Components own only transient local form state** until submit or autosave boundaries.

Do not let route components call `useFetch` directly for feature logic once the app grows past bootstrap screens. Keep data access behind repositories so online and offline behavior stays swappable.

## Data Flow

### Online read flow

```text
Page -> Feature composable -> Repository -> API client -> Store
                                          -> optional cache persist
```

### Client write flow while online

```text
Form -> Feature composable -> Repository
                             -> optimistic local store update
                             -> API write
                             -> cache reconcile
```

### Client write flow while offline

```text
Form -> Feature composable -> Repository
                             -> validate payload locally
                             -> persist queue item in IndexedDB
                             -> optimistic local store update
                             -> mark entity as pending sync
```

### Reconnect sync flow

```text
Connectivity regained -> Sync orchestrator -> Queue batch builder
                                       -> POST /tracking/sync for tracking entries
                                       -> direct POST for queued messages
                                       -> reconcile success/failure per item
                                       -> update stores and IndexedDB metadata
```

## Module Ownership

Organize the codebase by bounded domains, not by file type alone.

### Platform modules

| Module | Owns |
|--------|------|
| `session` | tokens, current user, active role, logout, session bootstrap |
| `network` | connectivity, online-required decisions, sync trigger events |
| `pwa` | install prompt, update prompt, service worker state |
| `notifications` | web push subscription registration, permission state, notification preferences |
| `locale` | RTL flags, Persian numerals, Jalali formatting, text normalization |
| `ui` | global toasts, bottom sheets, modals, tabs, loading states |

### Client domain modules

| Module | Owns |
|--------|------|
| `client-plan` | active plan, plan days/meals/options, cached food subset |
| `client-tracking` | food, water, sleep, exercise, medication, body logs, queue metadata |
| `client-messages` | conversation cache, unread count, queued outbound messages |
| `client-labs` | upload form state, online-only upload workflow, lab list cache |
| `client-food-requests` | request submit/list flow |
| `client-settings` | avatar, notification preferences, profile read models |

### Nutritionist domain modules

| Module | Owns |
|--------|------|
| `nutritionist-clients` | client list/detail/read models |
| `nutritionist-plans` | plan builder/editor and active plan transitions |
| `nutritionist-tracking-review` | client tracking queries and charts |
| `nutritionist-messages` | client conversation list/detail |
| `nutritionist-foods` | shared food CRUD/search |
| `nutritionist-medications` | shared medication CRUD/search |
| `nutritionist-food-requests` | approve/reject queue |

### Admin domain modules

| Module | Owns |
|--------|------|
| `admin-stats` | platform dashboard KPIs |
| `admin-nutritionists` | nutritionist CRUD and activation status |
| `admin-foods` | global food governance |
| `admin-medications` | global medication governance |

## Suggested Directory Layout

```text
app/
  app.vue
  assets/css/main.css
  components/
    ui/
    forms/
    feedback/
    charts/
    client/
    nutritionist/
    admin/
  composables/
    platform/
    session/
    client/
    nutritionist/
    admin/
  layouts/
    default.vue
    auth.vue
    client.vue
    nutritionist.vue
    admin.vue
  middleware/
    auth.global.ts
    role-client.ts
    role-nutritionist.ts
    role-admin.ts
    online-required.ts
  pages/
    index.vue
    auth/
    client/
    nutritionist/
    admin/
  plugins/
    api.ts
    pinia.ts
    pwa.client.ts
    push.client.ts
  stores/
    session.ts
    ui.ts
    sync.ts
    client-plan.ts
    client-tracking.ts
    client-messages.ts
    nutritionist-clients.ts
    nutritionist-plans.ts
    admin-stats.ts

lib/
  api/
    http.ts
    auth.ts
    errors.ts
    endpoints/
      auth.ts
      users.ts
      plans.ts
      tracking.ts
      messages.ts
      foods.ts
      medications.ts
      labs.ts
      notifications.ts
  domain/
    auth/
    plans/
    tracking/
    messages/
    foods/
    medications/
    users/
    notifications/
  offline/
    db.ts
    schema.ts
    queue.ts
    sync-orchestrator.ts
    sync-policies.ts
    message-queue.ts
  design/
    tokens.css
    semantic.ts
    motion.ts
  i18n/
    rtl.ts
    persian.ts
    jalali.ts
  utils/
    dates.ts
    files.ts
    ids.ts
    numbers.ts
```

## Routing and Middleware Strategy

Use middleware for access control and bootstrap only, not domain data fetching.

### Recommended middleware split

- `auth.global.ts`: restore session from persisted token state, fetch `GET /auth/me` when needed, and redirect anonymous users away from protected areas.
- `role-client.ts`, `role-nutritionist.ts`, `role-admin.ts`: enforce role-specific route entry.
- `online-required.ts`: block routes that cannot function offline, such as uploads, plan editing, admin management, and nutritionist CRUD flows.

### Redirect model

- Anonymous user -> `auth` routes only
- Authenticated client -> `client` area only
- Authenticated nutritionist -> `nutritionist` area only
- Authenticated super admin -> `admin` area only

Keep redirects centralized so role checks are not duplicated across pages.

## Auth and Session Handling

The API uses bearer access tokens plus refresh tokens. Architect the session layer around that contract.

### Recommended session behavior

- Keep the access token in memory and refresh proactively on API 401 or shortly before expiry.
- Persist the refresh token and minimal session metadata in browser storage through a session store adapter.
- On app boot, attempt silent refresh first, then call `GET /auth/me` to restore the role.
- Broadcast logout across tabs with `BroadcastChannel`.
- On explicit logout, call `POST /auth/logout`, clear IndexedDB queue state tied to that user if required by privacy policy, and clear all cached auth artifacts.

### Important constraint

Offline mode does **not** mean offline login. A user must have authenticated previously on the device. If the session cannot be refreshed because the device is offline, the client area may still expose already-cached plan/messages/queued drafts, but privileged server actions remain pending until connectivity returns.

## API Integration Pattern

Use an API layer based on `ofetch` with one shared client instance and feature-specific repositories.

### API client responsibilities

- Base URL and JSON defaults
- Bearer token injection
- Single-flight refresh handling so multiple 401s do not trigger refresh storms
- Error normalization from the API envelope into typed frontend errors
- Multipart helpers for avatar, messages, and lab uploads

### Repository contract pattern

Each domain should expose query and command methods, for example:

```text
plansRepository.getActiveClientPlan()
trackingRepository.logWater(entry)
trackingRepository.syncPending()
messagesRepository.getConversation(cursor)
messagesRepository.sendMessage(payload)
notificationsRepository.subscribePush(subscription)
```

This keeps route components and stores independent from raw endpoint details.

## Offline Architecture

The PRD already defines Service Worker + IndexedDB + Sync Manager. Keep those responsibilities separate.

### Service worker responsibilities

- App shell caching
- Static assets and font caching
- Safe caching of selected GET responses for previously visited client data
- Update lifecycle messaging (`offlineReady`, `needRefresh`)

### Service worker should not own

- Business conflict resolution
- Auth refresh logic
- Queue semantics for tracking or messages
- Multipart upload retry behavior

### IndexedDB stores to create

| Store | Purpose |
|-------|---------|
| `activePlan` | last active client diet plan snapshot |
| `trackingQueue` | pending food, water, sleep, exercise, medication, body entries |
| `messageCache` | last N client messages per conversation |
| `messageQueue` | unsent client messages created offline |
| `syncMeta` | retry counts, last sync time, error state |
| `notificationPrefsCache` | last-known preferences for settings UI |
| `profileCache` | minimal user profile and avatar URL |

### Sync rules

- Use `POST /tracking/sync` for queued tracking items.
- Send queued messages one by one with the message endpoint because chat attachments are multipart and message ordering matters.
- Do not queue lab-result file uploads offline. Store an interrupted draft locally and ask the user to retry online.
- Exponential backoff belongs to the sync orchestrator, not the page.
- Preserve server authority on timestamps and reconciliation.

### Conflict approach

The PRD specifies last-write-wins with server timestamp. Reflect that explicitly in repositories:

- queue item keeps `localId`, `createdAt`, `retryCount`, `lastError`
- sync response updates local entity state with server timestamps
- failed items remain visible in a sync center for manual retry

## Uploads and Attachments

Uploads need a separate path from normal JSON writes.

### Upload classes

- `avatar upload`: immediate multipart request, no offline queue
- `lab result upload`: immediate multipart request or link submission, no offline queue
- `chat attachment`: online multipart request when connected, offline queue allowed only for text-only drafts unless the product later accepts large local file staging complexity

Recommendation: phase 1 chat offline should queue **text messages only**. File attachments can stay online-required until the queueing and storage costs are validated.

## Push Subscription and Notification Boundaries

Keep push orchestration in a dedicated platform module.

### Notification flow

```text
Session restored -> browser permission check -> service worker ready
                 -> push subscription lookup/create
                 -> POST /push/subscribe
                 -> PATCH /notifications/preferences on settings updates
```

Rules:

- Subscription bootstrap should happen after login, not before auth.
- Permission prompting should be user-triggered from an in-app education surface, not automatic on first render.
- Persist only the subscription registration state needed for UI; the server remains source of truth.

## Design System Boundaries

Tailwind CSS 4 should define tokens and utilities, but feature teams should consume a small design system rather than raw utility soup everywhere.

### Recommended split

| Layer | Contents |
|-------|----------|
| `tokens` | colors, spacing, radius, shadows, z-index, motion, typography scales, semantic aliases |
| `primitives` | button, input, textarea, select, checkbox, radio, card, badge, tabs, sheet, dialog, toast, list row |
| `patterns` | mobile form section, metric card, meal timeline, message bubble, stat strip, empty state, sync banner |
| `feature components` | plan day viewer, water tracker, medication log form, client card, plan editor section |

Rules:

- Keep Tailwind tokens and semantic CSS variables in `lib/design`.
- Keep Jalali, Persian numerals, and RTL spacing behavior out of individual feature components; expose them through shared primitives and helpers.
- Charts, date pickers, and upload widgets should be wrapped so vendor libraries do not leak across the app.

## Persian, RTL, and Jalali Concerns

Treat locale handling as platform infrastructure, not as presentation detail.

### Centralize these concerns

- Jalali date parsing and formatting
- Persian numeral rendering and reverse conversion to ASCII digits for API payloads
- Text normalization for Persian search input and display consistency
- RTL-aware spacing, icon mirroring, swipe direction, and chart axis defaults
- Mobile keyboard-safe forms for Persian input

### Recommendation

Create a single locale helper package in `lib/i18n` and require all date/number transforms to go through it. Avoid mixing `Date`, ad hoc Jalali libraries, and manual string transforms in page code.

## Patterns to Follow

### Pattern 1: Repository + Store + Cache split

**What:** Store exposes reactive state, repository decides source of truth, cache persists offline state.

**When:** All client domains, and most nutritionist/admin read flows.

**Example:**

```typescript
// app/composables/client/useWaterTracking.ts
export function useWaterTracking() {
  const trackingStore = useClientTrackingStore()
  const trackingRepository = useTrackingRepository()

  const addWater = async (payload: AddWaterInput) => {
    await trackingRepository.logWater(payload)
    return trackingStore.waterEntries
  }

  return {
    entries: computed(() => trackingStore.waterEntries),
    addWater,
  }
}
```

### Pattern 2: Area shell isolation

**What:** Each role area gets its own layout, nav model, and middleware contract.

**When:** From the first authenticated milestone.

**Why:** This prevents client offline concerns from polluting staff/admin UX and keeps role-based expansion phaseable.

### Pattern 3: Online-required capability flags

**What:** Feature routes declare whether connectivity is mandatory.

**When:** Uploads, plan editing, admin operations, and any screen without meaningful cached fallback.

**Why:** The app should fail explicitly instead of silently pretending an online-only feature works offline.

## Anti-Patterns to Avoid

### Anti-Pattern 1: Pinia as the offline database

**What:** Treating stores as the durable source of offline truth.

**Why bad:** Reloads, multi-tab behavior, large plan trees, and retry metadata become fragile.

**Instead:** Keep durable client data in IndexedDB and hydrate stores from repositories.

### Anti-Pattern 2: One giant `api.ts` with endpoint calls everywhere

**What:** Pages and components call raw endpoints directly.

**Why bad:** Endpoint details leak into UI code, role branching spreads, and offline support becomes impossible to retrofit.

**Instead:** Keep endpoint calls inside repositories organized by domain.

### Anti-Pattern 3: Shared route space with conditional UI by role

**What:** A single dashboard tree that hides or shows controls based on role.

**Why bad:** It creates authorization bugs, muddy navigation, and unclear phase boundaries.

**Instead:** Separate `client`, `nutritionist`, and `admin` route roots.

### Anti-Pattern 4: Service worker as business logic engine

**What:** Pushing queue reconciliation, auth, and domain rules into the service worker.

**Why bad:** Debugging becomes difficult and state coordination breaks.

**Instead:** Keep the worker focused on caching and lifecycle events; keep sync logic in app code.

## Build Order Implications

The architecture supports a clean phase order.

### Phase 1: Platform foundation

- Nuxt app skeleton
- Tailwind 4 token layer and core primitives
- Pinia and repository contracts
- API client, error model, auth bootstrap scaffolding
- RTL/Jalali/Persian utility layer
- PWA manifest and service worker baseline

### Phase 2: Auth and route-area shell

- Auth pages for admin/nutritionist login and client OTP
- Session restoration and refresh
- Role middleware and layouts
- Basic profile bootstrap

### Phase 3: Client read-first experience

- Active plan read flow
- Client home/dashboard shell
- Cache hydration for plan and profile
- Sync status primitives

### Phase 4: Client offline tracking core

- IndexedDB schema
- Tracking queue and bulk sync
- Offline banners, pending states, retry surface
- Water, food, sleep, exercise, medication, body logging

### Phase 5: Client communication and notifications

- Cached messages and unread count
- Text-message offline queue
- Push subscription and notification preferences
- Lab result online-only upload flow

### Phase 6: Nutritionist workspace

- Client list/detail
- Tracking review
- Messaging
- Plan management
- Shared food/medication CRUD

### Phase 7: Super admin workspace and hardening

- Nutritionist management
- Platform stats
- Shared database governance
- Permission edge cases, performance, and polish

This order builds platform invariants first, then the most product-critical client value, then staff operations.

## Scalability Considerations

| Concern | At 100 users | At 10K users | At 1M users |
|---------|--------------|--------------|-------------|
| Route complexity | Area prefixes are enough | Add route-level code splitting by area | Consider Nuxt layers or package extraction for staff/admin features |
| State size | Pinia + Dexie is simple | Normalize large plan/message caches per feature | Introduce stronger cache invalidation and storage pruning policies |
| Offline queue | Single queue table is fine | Split tracking vs message queues | Add richer observability and remote sync diagnostics |
| Design system | Shared primitives are enough | Add stricter semantic component contracts | Package the design system into its own workspace module |
| Locale logic | Helper module is enough | Enforce locale adapters on all inputs/charts | Treat locale/date layer as a standalone internal package |

## Sources

- Project scope: `.planning/PROJECT.md`
- Product constraints: `docs/PRD.md`
- API contract: `docs/API.md`
- Nuxt 4 docs via Context7 CLI: `/websites/nuxt_4_x`
- Pinia docs via Context7 CLI: `/vuejs/pinia`
- PWA docs via Context7 CLI: `/websites/vite-pwa-org_netlify_app` and `/vite-pwa/vite-plugin-pwa`
- Tailwind CSS 4 docs via Context7 CLI: `/tailwindlabs/tailwindcss.com`