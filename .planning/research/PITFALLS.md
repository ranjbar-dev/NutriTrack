# Domain Pitfalls

**Domain:** Nutrition Management PWA (Go + PostgreSQL + Nuxt 4, Persian RTL, Offline-first)
**Researched:** 2025-07-14
**Overall Confidence:** HIGH — based on PRD/phases analysis, official library documentation (Context7), and domain expertise

---

## Critical Pitfalls

Mistakes that cause rewrites, data loss, or security breaches.

---

### Pitfall 1: N+1 Query Explosion on Diet Plan Loading

**What goes wrong:** The diet plan is 5 levels deep: Plan → Days → Meals → Options → Items. A naïve Go implementation loads each level with separate queries per parent row. For a realistic 7-day plan with 5 meals/day, 3 options/meal, and 4 items/option, this produces: 1 (plan) + 7 (days) + 35 (meals) + 105 (options) + 420 (items) = **568 queries** for a single plan load. Response times blow past the 500ms target into multi-second territory.

**Why it happens:** Go lacks an ORM with automatic eager loading like Rails or Django. When using `pgx` directly or `sqlc`, developers write "get children by parent ID" queries in loops. It works in development with small plans, then collapses under realistic data.

**Consequences:** API exceeds 500ms target. Client plan view loads slowly. Mobile users on 3G abandon the app. The problem worsens with every plan created because query count is data-proportional.

**Prevention:**
- Load the plan tree in **2–3 batch queries**, not per-row queries:
  1. Query 1: `SELECT * FROM plan_days WHERE diet_plan_id = $1`
  2. Query 2: `SELECT m.*, mo.*, moi.* FROM meals m JOIN meal_options mo ON ... JOIN meal_option_items moi ON ... WHERE m.plan_day_id = ANY($1)` with the day IDs from query 1
  3. Assemble the tree in-memory in Go
- Use `pgx.Batch` (pgx's `SendBatch`) to send multiple queries in a single round-trip when JOINs become too wide
- Alternatively, use a single query with `json_build_object`/`json_agg` to return the entire tree as nested JSON from PostgreSQL directly
- Add a `GET /api/plans/:id/full` integration test that asserts query count ≤ 5 using a query counter middleware

**Detection:**
- Enable `pgx` query logging in development; grep logs for repeated query patterns
- Add request-scoped query counters in middleware — alert if any request exceeds 10 queries
- Monitor p95 response time on `/api/client/plan` from Phase 3 onward

**Phase to address:** Phase 3 (Diet Plan Engine) — this is the core query pattern; fix it before building anything on top of it.

**Confidence:** HIGH — pgx batch query pattern verified via Context7 docs.

---

### Pitfall 2: Persian Text Search Silently Fails or Returns No Results

**What goes wrong:** PostgreSQL's built-in full-text search (`to_tsvector`/`to_tsquery`) has **no Persian language dictionary**. If you use `to_tsvector('simple', name)`, Persian text tokenizes but gets no stemming or stop-word removal. Worse, if someone defaults to `'english'` configuration, Persian text may not tokenize properly at all. The PRD correctly notes using `pg_trgm` instead, but the implementation has its own trap: `pg_trgm` treats characters as word/non-word based on locale settings. With the wrong database locale (e.g., `C` locale), Persian characters may be treated as non-word characters and excluded from trigrams entirely, producing empty trigram sets and zero search results.

**Why it happens:** PostgreSQL ships with dictionaries for ~30 languages — Persian is not one of them. Developers set up the database without checking locale, use the default `C` locale (common in Docker images), and never test search with actual Persian data until late in development.

**Consequences:** Food search returns no results or irrelevant results. Nutritionists cannot find foods they just added. The entire diet plan builder workflow is blocked because food search is a prerequisite.

**Prevention:**
- Set PostgreSQL locale to `fa_IR.UTF-8` or at minimum `en_US.UTF-8` in the Docker image — **never** use `C` or `POSIX` locale. Use `POSTGRES_INITDB_ARGS: "--locale=en_US.UTF-8 --encoding=UTF8"` in Docker Compose
- Use `pg_trgm` with `gin_trgm_ops` GIN index on `foods.name` for fuzzy matching: `CREATE INDEX idx_foods_name_trgm ON foods USING GIN (name gin_trgm_ops);`
- Use `similarity()` or `%` operator, not `LIKE` or `to_tsvector`, for Persian text search
- **Normalize Persian text** before storing AND before searching (see Pitfall 3)
- Seed the database with 50+ real Persian food names during Phase 2 and write integration tests that verify search returns correct results for partial inputs like `بر` matching `برنج` (rice)
- Benchmark `pg_trgm` with 1,000+ food items; have Meilisearch as a documented fallback if performance degrades past 200ms

**Detection:**
- Search for a known Persian food name with 2–3 character substring — if zero results, locale is wrong
- Run `SELECT show_trgm('برنج');` — if the result is `{}` (empty array), the locale is misconfigured
- Test early in Phase 2 with actual Persian data, not Latin placeholders

**Phase to address:** Phase 2 (Core Data Domain) — the food/medication search is the first place this surfaces.

**Confidence:** HIGH — PostgreSQL `pg_trgm` docs verified via Context7 confirm trigram extraction depends on locale-dependent word character classification.

---

### Pitfall 3: Persian Character Normalization — ی vs ي and ک vs ك

**What goes wrong:** Persian and Arabic share a script but differ in key characters. Persian uses `ی` (U+06CC) and `ک` (U+06A9), while Arabic uses `ي` (U+064A) and `ك` (U+0643). Users type both interchangeably depending on their keyboard. A nutritionist stores `کباب` (with Persian ک), but a client searches `كباب` (with Arabic ك) — `pg_trgm` sees these as different characters and returns no results.

**Why it happens:** Mobile keyboards often mix Arabic and Persian character variants. Even official Persian keyboards sometimes produce Arabic variants. Developers who test with a single keyboard never encounter the mismatch.

**Consequences:** Search misses valid results. Duplicate food entries accumulate (same food stored with both character variants). Data integrity silently degrades.

**Prevention:**
- Create a `normalize_persian(text)` PostgreSQL function that replaces: `ي` → `ی`, `ك` → `ک`, Arabic numbers → Persian/Latin, `ة` → `ه`, zero-width non-joiner cleanup
- Apply this function via a trigger `BEFORE INSERT OR UPDATE` on `foods.name`, `medications.name`, and any other searchable Persian text column
- Apply the same normalization in the Go backend before any search query parameter hits the database
- Also apply in the Nuxt frontend for display consistency using a `normalizePersian()` utility
- Consider an `IMMUTABLE` SQL function wrapper so it can be used in GIN index expressions: `CREATE INDEX idx_foods_name_norm ON foods USING GIN (normalize_persian(name) gin_trgm_ops);`

**Detection:**
- Insert two food items: one with Arabic ك, one with Persian ک — if both are accepted without deduplication, normalization is missing
- Search for a word typed with the opposite keyboard variant — if zero results, normalization is missing

**Phase to address:** Phase 2 (Core Data Domain) — must be in place before any real data enters the system.

**Confidence:** HIGH — this is a well-known problem in every Persian/Arabic software project.

---

### Pitfall 4: iOS PWA Storage Eviction Destroys Offline Data

**What goes wrong:** iOS Safari evicts IndexedDB data for PWAs after approximately 7 days of non-use, or when the device is under storage pressure. A client opens the app on Monday, goes offline, returns on the next Monday — their cached diet plan and unsynced tracking logs are gone. **Data the user believes was saved is silently deleted by the OS.**

**Why it happens:** Apple treats PWA storage as "ephemeral" compared to native apps. There is no `navigator.storage.persist()` API support in iOS Safari (or it's non-functional). The Background Sync API is also not supported on iOS. This affects all PWA offline strategies, not just NutriTrack.

**Consequences:** Clients lose unsynced tracking data (food logs, water intake, measurements). Cached diet plan disappears, showing a blank screen when opened offline. Users lose trust in the app. This is the single most damaging pitfall for a client-facing offline-first PWA targeting iOS.

**Prevention:**
- **Sync aggressively when online:** Don't let unsynced data linger. Process the sync queue immediately on every app open and on every network connectivity change, not just when the queue is "full"
- **Show unsynced item count prominently:** A persistent badge like "3 items pending sync" makes users aware that data hasn't been saved to server yet, motivating them to connect
- **Minimize offline cache size:** Only cache the active diet plan and last 50 messages (as the PRD specifies). Don't cache historical data — reduce the risk of the OS reclaiming storage
- **Re-fetch on app resume:** When the app becomes visible (`visibilitychange` event), check if IndexedDB stores are intact. If the active plan is missing, show a "Reconnect to refresh your plan" message rather than a blank screen
- **Fallback polling:** Since iOS doesn't support Background Sync API, implement polling on app open as the PRD already specifies — don't rely on `BackgroundSyncPlugin` in Workbox
- **Test on actual iOS devices** with storage pressure simulation (fill device storage, wait 7+ days)

**Detection:**
- After Phase 6, test by: installing PWA on iOS, going offline, waiting 8 days, reopening — if diet plan is gone, eviction mitigation is insufficient
- Monitor for API requests that look like "first load" from clients who previously had cached data (server can detect if a client re-requests their full plan more than once per plan update)

**Phase to address:** Phase 6 (Offline & PWA) — but awareness needed from Phase 1 when designing the offline architecture.

**Confidence:** HIGH — iOS storage eviction is extensively documented and a known iOS PWA limitation.

---

### Pitfall 5: Offline Sync Queue Creates Duplicates or Loses Order

**What goes wrong:** The sync queue processes entries in FIFO order, but race conditions arise in multiple scenarios:
1. **Double-submit:** User taps "save" twice quickly → two queue entries with different `local_id` values for the same logical action
2. **Partial sync failure:** Queue has entries [A, B, C]. A syncs, B fails (server error), C syncs — now the client has a gap, and if B was a prerequisite for C (e.g., logging food for a meal that was updated), the data is inconsistent
3. **Concurrent tabs/windows:** Two browser tabs both process the same sync queue simultaneously → duplicate server-side records even with `local_id` deduplication, because both tabs read the queue before either marks items as syncing

**Why it happens:** IndexedDB doesn't have row-level locking across tabs. The `navigator.onLine` event fires in all tabs simultaneously. The sync manager in each tab starts processing the same queue.

**Consequences:** Duplicate tracking entries in the database. Nutritionist sees the client logged the same meal twice. Water intake shows double. Data integrity is silently compromised.

**Prevention:**
- **Use Web Locks API** (`navigator.locks.request('sync-queue', ...)`) to ensure only one tab processes the sync queue at a time
- **Debounce user actions:** Disable submit buttons after first tap; use a 300ms debounce on all tracking log submissions
- **Mark items as "syncing" in IndexedDB** before sending to server, and only mark "synced" on success. On failure, revert to "pending"
- **Process queue sequentially**, not in parallel — one item at a time, wait for server response before processing the next
- **Server-side `local_id` upsert** is the last line of defense: `INSERT ... ON CONFLICT (local_id) DO NOTHING` — already planned in Phase 4, but must be tested under concurrent conditions
- **Stop sync on first failure** for the same entity type (don't skip failed items and continue — this can cause ordering issues)

**Detection:**
- Query the database for duplicate `local_id` values — should be zero
- Query for suspiciously close timestamps on tracking entries (e.g., two water logs within 1 second)
- Integration test: open two browser tabs offline, create entries in both, go online, verify no duplicates

**Phase to address:** Phase 6 (Offline & PWA) — but `local_id` deduplication infrastructure must be in Phase 4.

**Confidence:** HIGH — Dexie.js docs (Context7) confirm IndexedDB has no cross-tab locking; Web Locks API is the standard solution.

---

### Pitfall 6: Row-Level Authorization Bypass on Indirect Access Paths

**What goes wrong:** Every endpoint that touches client data must verify `WHERE nutritionist_id = $current_user_id`. But authorization is forgotten on indirect access paths:
- `GET /api/plans/:plan_id` — checks plan exists but not that the current nutritionist owns the client
- `GET /api/nutritionist/clients/:client_id/food-logs` — checks client exists but not ownership
- `PATCH /api/messages/:id/read` — checks message exists but not that current user is a participant
- Nutritionist A guesses a valid `client_id` UUID belonging to Nutritionist B and accesses their data

**Why it happens:** Developers add authorization on the "main" CRUD endpoints but miss it on helper endpoints, report endpoints, or endpoints added later. Each new endpoint is a potential authorization gap. With 30+ endpoints touching client data, the probability of missing at least one is high.

**Consequences:** A nutritionist can view another nutritionist's client data — a privacy breach. Depending on Iranian data protection expectations, this could be a trust-destroying incident.

**Prevention:**
- **Authorization middleware at the repository layer**, not just the handler layer. Every repository method that accepts a `client_id` should JOIN to `users WHERE nutritionist_id = $current_user_id` and return 404 (not 403) if the client doesn't belong to the nutritionist. This makes it impossible to forget authorization in a handler.
- **Write a reusable `authorizeClientAccess(ctx, clientID)` function** that every handler calls before any data access. Make code review reject any handler that accesses client data without calling it.
- **Create an automated security test suite:** for every endpoint that accepts `client_id`, `plan_id`, `message_id`, etc., send a request with a valid JWT from Nutritionist B and verify 403/404 response. This suite runs in CI on every push.
- Return **404 (not 403)** for unauthorized access to prevent enumeration — the nutritionist shouldn't even know the resource exists.

**Detection:**
- Grep the codebase for endpoints that accept `:client_id` or `:plan_id` in the route — verify each one has authorization
- Run the security test suite from Phase 7 — any 200 response for cross-nutritionist access is a failure
- Code review checklist item: "Does this endpoint verify resource ownership?"

**Phase to address:** Phase 1 (Foundation — design the pattern), validated continuously through Phases 2–7, comprehensive audit in Phase 7.

**Confidence:** HIGH — this is the #1 security risk per the PRD's own non-functional requirements.

---

### Pitfall 7: JWT Refresh Token Race Condition

**What goes wrong:** With a 15-minute access token, multiple concurrent API requests from the same client can all discover the token is expired at the same time. Each one tries to refresh using the same refresh token. If the backend implements refresh token rotation (invalidate old token on use), only the first request succeeds — the rest get 401 errors because the refresh token has been consumed. The client gets logged out unexpectedly.

**Why it happens:** The Nuxt app makes parallel API calls on page load (e.g., load plan + load messages + load unread count). All use the same expired access token. All fail. All trigger the refresh flow simultaneously.

**Consequences:** Users get randomly logged out, especially after periods of inactivity (15+ minutes). On mobile with unstable connections, this happens frequently. The OTP re-login flow is disruptive (SMS delay, typing code).

**Prevention:**
- **Implement a refresh token queue** in the Pinia auth store: when multiple requests detect 401 simultaneously, only the FIRST one triggers the refresh. All others wait on a promise that resolves when the refresh completes, then retry with the new access token.
- In `ofetch`/`useFetch` interceptor: check if a refresh is already in progress (`isRefreshing` flag). If yes, push the failed request to a retry queue. When refresh completes, replay all queued requests.
- On the backend: implement a **grace period** for refresh token rotation — accept the old refresh token for 30 seconds after rotation, allowing concurrent requests to succeed. Use a `token_family` column to detect genuine token theft vs. race conditions.
- **Proactively refresh** the access token 1–2 minutes before expiry (check `exp` claim on each request), so it never actually expires during active use.

**Detection:**
- Users report "random logouts" — check server logs for clusters of 401s followed by multiple refresh token requests with the same old token
- Integration test: send 5 concurrent requests with an expired access token, verify all succeed after refresh

**Phase to address:** Phase 1 (Foundation) — the auth interceptor pattern must handle this from day one.

**Confidence:** HIGH — this is a well-documented JWT pattern problem. golang-jwt docs (Context7) confirm token validation but don't address rotation; this is application-level logic.

---

## Moderate Pitfalls

Issues that cause significant rework or user frustration but are recoverable.

---

### Pitfall 8: Shamsi (Jalali) Calendar Date Boundary Bugs

**What goes wrong:** Jalali and Gregorian calendars have different day boundaries, month lengths, and leap year rules. Common bugs:
1. **Diet plan date mapping:** A 7-day repeating plan uses `day_number` with modulo arithmetic. If the start date is Jalali but the modulo uses Gregorian date difference, days shift by 1 after the Jalali month boundary
2. **"Today" detection:** Server stores dates in Gregorian (PostgreSQL `date` type). Frontend converts to Jalali for display. If the timezone offset isn't accounted for, "today" in Tehran (UTC+3:30) may differ from "today" in UTC — a tracking log saved at 11 PM Tehran time appears under the wrong date
3. **Leap year:** Jalali years 1403, 1407, etc. have a leap day (30th of Esfand). If date validation rejects day 30 for Esfand, it breaks in leap years
4. **Date range queries:** Nutritionist queries client logs for "this month" — must convert Jalali month boundaries to Gregorian for the SQL WHERE clause

**Why it happens:** All dates are stored as Gregorian in PostgreSQL (correct approach), but conversion boundaries are tricky. `jalaali-js` handles conversion correctly, but developers apply it inconsistently.

**Prevention:**
- **All date storage is Gregorian** (never store Jalali dates in the database)
- **All date API parameters are Gregorian** (frontend converts before sending)
- **Create a `useShamsiDate` composable** (as docs/phases.md suggests) that wraps `jalaali-js` and handles: display formatting, "today" calculation with Tehran timezone offset, date range conversion
- **Set a fixed timezone** in the Go backend: `time.LoadLocation("Asia/Tehran")` — use this for all "today" calculations
- **Test Esfand 29/30 boundary** in both leap and non-leap years
- **Test midnight boundary:** log at 23:59 Tehran time → should be today's date, not tomorrow's

**Detection:**
- Create a tracking log at 11:30 PM Tehran time — if it shows under the wrong date, timezone handling is broken
- Check diet plan day mapping around Jalali month boundaries (especially months with 29 vs 30 vs 31 days)

**Phase to address:** Phase 1 (Foundation — establish the composable), Phase 3 (Diet Plan — date mapping), Phase 4 (Tracking — "today" detection).

**Confidence:** HIGH — Jalali calendar pitfalls are well-documented in Persian software development.

---

### Pitfall 9: Service Worker Caches Stale Diet Plan After Update

**What goes wrong:** The nutritionist updates the client's diet plan. The client opens the app — the service worker serves the old cached plan from the network-first strategy's stale fallback. If the client is on a slow connection where the network request times out, they see yesterday's plan indefinitely. Worse, if using `registerType: 'prompt'` and the user dismisses the update prompt, they stay on the old service worker and the cached API responses remain stale.

**Why it happens:** `vite-pwa/nuxt` with `NetworkFirst` strategy has a `networkTimeoutSeconds` option. If the network is slow but not down, the request times out and the cache serves stale data. The diet plan endpoint URL doesn't change (`GET /api/client/plan`), so the cache key is the same.

**Consequences:** Client follows the wrong diet plan. If the nutritionist changed meals for medical reasons, this is potentially harmful.

**Prevention:**
- **Don't use Workbox runtime caching for the diet plan endpoint.** Instead, use IndexedDB (Dexie) as the authoritative offline store. Fetch from API → update IndexedDB → render from IndexedDB. This gives you explicit control over freshness
- If using Workbox caching, add the plan's `updated_at` timestamp as a cache-busting query parameter: `GET /api/client/plan?v=${planUpdatedAt}`
- **Push notification on plan update:** When the nutritionist saves a new plan, send a web push notification to the client. The notification click handler forces a cache-busting fetch
- **Use ETag/If-None-Match** on the plan endpoint — if the plan hasn't changed, return 304 (fast). If it has, return the full plan. This is faster than always returning the full plan tree
- Set `registerType: 'autoUpdate'` (not `'prompt'`) to ensure service worker updates apply automatically without user interaction — critical for a mobile app where users don't understand update prompts

**Detection:**
- Update a plan as nutritionist, immediately check client view — if old plan shows, caching is too aggressive
- Test with Chrome DevTools network throttling set to "Slow 3G" — verify the plan updates within 10 seconds

**Phase to address:** Phase 6 (Offline & PWA) — service worker strategy configuration.

**Confidence:** HIGH — vite-pwa docs (Context7) confirm `registerType: 'prompt'` as default; must explicitly set `'autoUpdate'` for this use case.

---

### Pitfall 10: File Upload Security — Content Sniffing and Path Traversal

**What goes wrong:** Three attack vectors on file uploads:
1. **Content sniffing:** Attacker uploads a `.pdf` file that is actually an HTML file with malicious JavaScript. The browser sniffs the content type and executes the script when another user views/downloads it
2. **Path traversal:** Upload filename contains `../../etc/passwd` or `..\..\` — if the server uses the filename directly in the storage path, it writes outside the intended directory
3. **Storage exhaustion:** No per-user upload limit — one client uploads hundreds of large files, filling the 50GB annual budget in a week

**Why it happens:** Go's `multipart.FileHeader.Filename` contains the original filename as sent by the client. If used directly in `filepath.Join()`, path traversal is possible. Content-Type validation only checks the header, not actual content.

**Prevention:**
- **Never use the original filename for storage.** Generate a UUID filename: `{uuid}.{validated_extension}`
- **Validate actual content:** Read the first 512 bytes and use Go's `http.DetectContentType()` to verify. For PDFs, check for `%PDF` magic bytes. For JPEG, check `0xFF 0xD8`. Reject mismatches
- **Set `Content-Disposition: attachment`** on all file download responses to prevent browser content sniffing
- **Set `X-Content-Type-Options: nosniff`** header on all responses
- **Serve uploaded files through an authenticated proxy endpoint** (`GET /api/files/:id`) — never expose the filesystem path directly via static file serving
- **Enforce per-client storage limits:** max 100MB per client for lab results, max 50MB per conversation for message attachments. Track usage in the database
- **Store files outside the web root** (already in PRD: `/data/uploads/`)

**Detection:**
- Upload a `.pdf` file that is actually an HTML file — if the browser renders it as HTML, content sniffing protection is missing
- Upload a file with `../../test.txt` as filename — if a file appears outside the upload directory, path traversal is present
- Check server-side: `ls /data/uploads/` — filenames should all be UUIDs, never original names

**Phase to address:** Phase 5 (Communication Layer — messaging attachments, lab results upload).

**Confidence:** HIGH — Go's `multipart.FileHeader.Filename` path traversal is a documented security concern.

---

### Pitfall 11: Diet Plan Builder UI State Management Explosion

**What goes wrong:** The diet plan builder is a deeply nested reactive form: Plan → Days (tabs) → Meals (list) → Options (sub-list) → Items (sub-sub-list). Each item has a food picker modal, quantity input, and computed nutrition. Developers create a single massive Pinia store with deeply nested state, then watch for changes at every level. Vue's reactivity system triggers cascading recomputations — adding a single food item re-renders the entire plan.

**Why it happens:** The natural instinct is to model the store after the data structure (nested objects). Vue 3's reactivity deeply proxies nested objects by default, so every property change triggers watchers on the entire tree.

**Consequences:** Plan builder becomes slow with 5+ meals per day. Adding a food item causes a visible jank/delay. Persian mobile devices (often mid-range hardware) are hit hardest.

**Prevention:**
- **Flatten the Pinia store:** Use normalized state with IDs as keys, like a client-side database:
  ```
  planDays: { [dayId]: PlanDay }
  meals: { [mealId]: Meal }
  mealOptions: { [optionId]: MealOption }
  mealOptionItems: { [itemId]: MealOptionItem }
  ```
- **Use computed getters** for tree assembly: `getMealsForDay(dayId)` filters `meals` by `planDayId`
- **Compute nutrition per-option** only when that option's items change, not on every plan mutation. Use `computed` per option, not a single plan-wide watcher
- **Use `shallowRef`** for large lists that don't need deep reactivity
- **Debounce nutrition recalculation** by 100ms — user won't notice, but it prevents computing while they're still typing quantities

**Detection:**
- Build a plan with 7 days × 5 meals × 3 options × 4 items = 420 items. Add one more item — if there's visible jank, the state structure needs flattening
- Use Vue DevTools performance tab to check render count per component

**Phase to address:** Phase 3 (Diet Plan Engine) — frontend architecture decision.

**Confidence:** MEDIUM — specific to Vue 3 reactivity patterns; validated through general Vue performance documentation.

---

### Pitfall 12: Dexie.js IndexedDB Schema Versioning Breaks on Update

**What goes wrong:** Dexie requires keeping ALL previous version declarations when upgrading the schema. If you modify the schema in version 2 and remove the version 1 declaration, **existing users' databases fail to open** because Dexie can't determine how to upgrade from version 1. The app crashes silently (IndexedDB `onupgradeneeded` fails) and all offline data is lost.

**Why it happens:** Developers unfamiliar with Dexie treat version declarations like regular code — "we're on version 2 now, delete version 1." But Dexie uses version declarations as a migration chain, and users may be on any previous version.

**Consequences:** App crashes on launch for existing users after an update. All cached data lost. Users must clear browser data and re-login.

**Prevention:**
- **Never remove previous version declarations** from the Dexie database definition
- Keep a linear chain: `db.version(1).stores({...})`, `db.version(2).stores({...}).upgrade(...)`, etc.
- Use `.upgrade()` callbacks for data migration between versions
- Put the Dexie schema in a dedicated file (`~/db/schema.ts`) with clear comments: "DO NOT REMOVE OLD VERSIONS"
- **Test schema upgrades** by: opening app with version 1, storing data, deploying version 2, reopening — verify data is preserved and schema is updated
- Consider adding a version gate in the service worker update flow: force all queued data to sync before allowing the service worker update to activate (reducing the blast radius of schema migration bugs)

**Detection:**
- Open browser console — if IndexedDB errors appear after a deployment, schema versioning is broken
- Test with a user who hasn't opened the app since version N-1 — do they lose data?

**Phase to address:** Phase 6 (Offline & PWA) — but establish the schema versioning discipline from the first Dexie store creation.

**Confidence:** HIGH — Dexie.js official docs (Context7) explicitly document the version chain requirement.

---

### Pitfall 13: Push Notification Reminder Scheduler Sends Duplicates or Leaks Goroutines

**What goes wrong:** The reminder scheduler (goroutine that sends meal/medication reminders) has two common failure modes:
1. **Duplicate sends:** The goroutine checks "which reminders are due in the next minute?" every minute. If the query takes >1 minute, or if the server restarts mid-cycle, the same reminder is sent twice
2. **Goroutine leak:** Each reminder spawns a goroutine to call the Web Push API. If the push endpoint is slow or the subscription is stale (expired endpoint URL), the goroutine hangs forever. With 10,000 clients and 5 reminders/day, that's 50,000 goroutines — the server runs out of memory

**Why it happens:** Go makes goroutines cheap to spawn but doesn't enforce cleanup. Developers `go sendPush(subscription, payload)` without context cancellation or timeouts. The deduplication problem is classic distributed systems — without an idempotency key, retries produce duplicates.

**Prevention:**
- **Use a `processed_reminders` table** with a unique constraint on `(user_id, reminder_type, scheduled_time, date)`. Before sending, INSERT — if the constraint fires, skip. This survives server restarts
- **Set timeouts on all Web Push HTTP calls:** `http.Client{Timeout: 10 * time.Second}`. Never use the default client (no timeout)
- **Use a worker pool** instead of unbounded goroutines: buffer of 100 concurrent push notifications, queue the rest
- **Use context cancellation:** pass a context with timeout to each push goroutine. On server shutdown, cancel all pending pushes via context
- **Clean up stale subscriptions:** if a push endpoint returns 404 or 410 (Gone), delete the subscription from the database immediately
- **Compute daily reminders at midnight Tehran time** rather than per-minute checks — compute once, store in a sorted schedule, process with `time.AfterFunc`

**Detection:**
- Client receives the same meal reminder 2+ times — check `processed_reminders` table for duplicates
- Monitor goroutine count in Grafana — if it grows without bound, goroutines are leaking
- Check memory usage after running the scheduler for 24 hours — should be stable

**Phase to address:** Phase 6 (Offline & PWA) — push notification implementation.

**Confidence:** MEDIUM — Go goroutine leak pattern is well-documented; specific push notification scheduling needs validation at scale.

---

### Pitfall 14: Nuxt 4 Directory Structure Change Breaks Module Assumptions

**What goes wrong:** Nuxt 4 changes the default `srcDir` from project root (`./`) to `./app/`. The `~` alias now resolves to `app/` instead of the project root. If you scaffold with Nuxt 3 patterns (components at `./components/`, pages at `./pages/`) and then upgrade or use Nuxt 4 from the start without the new structure, import paths break, auto-imports fail, and the build silently produces wrong results or fails entirely.

**Why it happens:** Nuxt 4 was recently released. Many tutorials and community modules still assume Nuxt 3 directory structure. Copy-pasting config from Nuxt 3 projects or guides will produce subtle breaks.

**Consequences:** Build failures. Auto-imports stop working for components and composables. Middleware doesn't run. Pages don't render.

**Prevention:**
- **Use the Nuxt 4 directory structure from day one:** `app/` contains pages, components, composables, layouts, middleware, plugins, utils. `server/` is at root level. `shared/` for code shared between app and server
- Reference the exact Nuxt 4 directory layout from the official docs (confirmed via Context7):
  ```
  app/
    components/
    composables/
    layouts/
    middleware/
    pages/
    plugins/
    utils/
    app.vue
  server/
  shared/
  nuxt.config.ts
  ```
- **Don't set `srcDir: '.'`** to revert to Nuxt 3 structure — embrace the new structure for forward compatibility
- Verify all third-party Nuxt modules are compatible with Nuxt 4 before including them. Check their GitHub issues for "Nuxt 4" mentions
- Run the official migration codemod if starting from a Nuxt 3 template: `npx codemod@0.18.7 nuxt/4/migration-recipe`

**Detection:**
- `nuxi dev` shows "component not found" or "page not found" warnings — directory structure mismatch
- Auto-imports fail silently — composables return `undefined`

**Phase to address:** Phase 1 (Foundation) — project scaffolding must use correct structure.

**Confidence:** HIGH — Nuxt 4 directory structure changes confirmed via Context7 official docs.

---

## Minor Pitfalls

Issues that cause developer frustration or small UX problems, recoverable with low cost.

---

### Pitfall 15: RTL Layout Breaks with Mixed-Direction Content

**What goes wrong:** Nutritional values (numbers), measurement units, and food names mix RTL (Persian) and LTR (numbers) content on the same line. Without explicit `dir="auto"` or proper CSS logical properties, numbers appear on the wrong side, colons flip, and layouts look scrambled. Common issues:
- "۳۵۰ کالری" (350 calories) — the number should be left of the unit in RTL, but CSS may place it right
- Chart.js axis labels render LTR by default — Shamsi dates on the x-axis appear reversed
- Icons that have directional meaning (back arrows, progress indicators) point the wrong way

**Prevention:**
- Use Tailwind CSS RTL plugin (`tailwindcss-rtl`) and prefer **logical properties**: `ps-4` (padding-start) not `pl-4` (padding-left), `ms-2` not `ml-2`
- Set `dir="rtl"` on `<html>` element globally; use `dir="ltr"` explicitly on code blocks or numeric displays if needed
- Mirror directional icons in CSS: `.rtl .icon-back { transform: scaleX(-1); }`
- For Chart.js: set `options.scales.x.reverse = true` for RTL axis ordering; format date labels using the Shamsi composable
- Test every UI component with a full Persian data set, not lorem ipsum

**Phase to address:** Phase 1 (Foundation — Tailwind RTL setup), validated throughout all phases.

**Confidence:** HIGH — standard RTL development challenge.

---

### Pitfall 16: OTP Rate Limiting Bypassed by Phone Number Format Variations

**What goes wrong:** Iranian phone numbers can be formatted as `09123456789`, `9123456789`, `+989123456789`, or `00989123456789`. If the rate limiter keys on the raw input string, an attacker sends one OTP to each format — 4 OTPs instead of the intended max 3 per 10 minutes.

**Prevention:**
- **Normalize phone numbers** to a canonical format before any processing (rate limiting, OTP generation, user lookup). Strip all non-digits, remove leading `+98`, `0098`, or leading `0`. Store as `9123456789` (10 digits without country code)
- Validate against the pattern: `^9[0-9]{9}$` after normalization
- Apply rate limiting on the **normalized** number, not the raw input
- Also rate limit by **IP address** as a secondary control (max 10 OTP requests per IP per 10 minutes)

**Phase to address:** Phase 1 (Foundation — OTP flow).

**Confidence:** HIGH — standard Iranian telecom formatting issue.

---

### Pitfall 17: Connection Pool Exhaustion Under Load

**What goes wrong:** The diet plan full-tree query (Pitfall 1's solution using 2–3 queries with JOINs) holds a database connection for longer than simple CRUD queries. With 500 concurrent users and the default pgxpool size (likely 10–25 connections), connection starvation occurs during peak hours. Requests queue up, timeouts cascade, and the app appears "down" even though the server is running.

**Prevention:**
- Set `pgxpool` configuration explicitly: `MaxConns: 50`, `MinConns: 10`, `MaxConnLifetime: 1 hour`, `MaxConnIdleTime: 30 minutes`
- Add a **connection pool health metric** to the `/health` endpoint: report `pool.Stat().AcquiredConns()` vs `pool.Stat().TotalConns()`
- Set **query-level timeouts** using context: `ctx, cancel := context.WithTimeout(ctx, 5*time.Second)`
- Monitor pool utilization in Grafana — alert if acquired connections exceed 80% of max
- For the polling endpoint (`GET /api/messages/new`), ensure the query is fast (<10ms) — 500 users polling every 10 seconds = 50 req/sec of lightweight queries

**Phase to address:** Phase 1 (Foundation — database setup), validated in Phase 7 (load testing).

**Confidence:** HIGH — pgxpool configuration confirmed via Context7 docs.

---

### Pitfall 18: Message Polling Drains Mobile Battery

**What goes wrong:** The 10-second polling interval for the messaging system means the phone's radio never enters idle mode while the chat view is open. On mobile devices, this drains battery significantly. Users notice and blame the PWA.

**Prevention:**
- **Only poll when the chat view is active** (the PRD already specifies this — enforce it strictly)
- Use `visibilitychange` API: stop polling when the tab/app is backgrounded, resume when foregrounded
- Consider **adaptive polling**: 10 seconds for the first 2 minutes (active conversation), then back off to 30 seconds, then 60 seconds
- Use HTTP/2 on Traefik — connection reuse reduces the overhead of frequent small requests
- Include `last_message_id` in the poll request — if nothing new, return 204 No Content (minimal response body)

**Phase to address:** Phase 5 (Communication Layer).

**Confidence:** MEDIUM — polling frequency impact varies by device; adaptive polling is a common mitigation.

---

## Technical Debt Patterns

Patterns that feel fine initially but accumulate debt over 6–12 months.

| Pattern | Looks Like | Actually Is | Prevention |
|---------|-----------|-------------|------------|
| Inline authorization checks | `if plan.NutritionistID != currentUser.ID` in every handler | Inconsistent authorization that will be missed in new endpoints | Repository-level ownership verification |
| Persian strings hardcoded in Go handlers | Error messages like `"رژیم غذایی فعال یافت نشد"` scattered across handlers | Unmaintainable; changes require code deploys | Centralized error message map (`internal/errors/messages.go`) |
| Single Pinia store for all tracking data | `useTrackingStore` with food logs, water, sleep, exercise, meds, measurements | Massive store that re-renders everything when any tracking type updates | One store per tracking domain: `useWaterStore`, `useFoodLogStore`, etc. |
| Timestamp without timezone | `created_at timestamp` in PostgreSQL | Ambiguous times that break when server timezone changes | Always use `timestamptz` — `created_at timestamptz DEFAULT NOW()` |
| ENV vars without validation | `os.Getenv("SMS_API_KEY")` used directly | Silent failures when env vars are missing; empty API key means no SMS sent | Validate all required env vars on startup; fail fast with clear error |

---

## Integration Gotchas

### Kavenegar (Iranian SMS Gateway)

| Issue | Impact | Mitigation |
|-------|--------|------------|
| API rate limits during peak hours | OTP delivery delayed 10–30 seconds | Implement retry with exponential backoff; show "please wait" with timer on frontend |
| SMS delivery failures to certain carriers | Client never receives OTP, cannot login | Log delivery status callbacks; show "Didn't receive? Try again in 60 seconds" |
| Farsi text encoding in SMS body | OTP message shows garbled characters | Use UTF-8 explicitly in API call; test with actual Iranian SIM cards |
| API downtime (rare but happens) | All client logins blocked | Implement a fallback SMS provider adapter; alert on consecutive failures |

### Web Push (webpush-go)

| Issue | Impact | Mitigation |
|-------|--------|------------|
| iOS Safari 16.4+ required for push | Older iPhones cannot receive push notifications | Feature-detect `PushManager` before showing the permission prompt; degrade gracefully |
| Push subscription expires silently | Notifications stop without warning | Re-subscribe on every app open; handle 410 Gone responses by removing stale subscriptions |
| VAPID key rotation | All existing subscriptions become invalid | Plan VAPID key management; don't change keys after launch without a re-subscription migration |

### Chart.js in RTL Context

| Issue | Impact | Mitigation |
|-------|--------|------------|
| Axis labels default to LTR | Shamsi dates appear in wrong order on x-axis | Set `options.scales.x.reverse = true`; use custom tick callback for Shamsi formatting |
| Tooltip positioning breaks in RTL | Tooltips overflow off-screen on the left | Set `options.plugins.tooltip.rtl = true` (supported in Chart.js 4+) |
| Number formatting | Uses Latin numerals by default | Custom formatter using `new Intl.NumberFormat('fa-IR')` for Persian numerals |

---

## Performance Traps

| Trap | Where | Impact | Fix |
|------|-------|--------|-----|
| Full plan tree serialization | `GET /api/client/plan` | 100KB+ JSON response for complex plans | Compress with gzip (Traefik/Fiber handles this); consider pagination by day |
| Tracking log list without pagination | `GET /api/nutritionist/clients/:id/food-logs` | 10,000+ rows returned for long-term clients | Always paginate with cursor-based pagination; default limit 50 |
| Unrestricted food database search | `GET /api/foods?q=ب` | Single-character trigram search scans entire table | Minimum 2-character search query; abort previous search on new keystroke (debounce 300ms) |
| Large IndexedDB sync queue | Sync on reconnect | 100+ queued items processed sequentially = 100+ API calls = 1+ minute sync | Batch sync: send multiple items per request `POST /api/client/tracking/batch` |
| Message polling on every tab | All open chat tabs poll simultaneously | 5 tabs × 6 req/min = 30 req/min from one user | Use BroadcastChannel API to elect one tab as the poller |

---

## Security Mistakes

| Mistake | Risk | Phase | Prevention |
|---------|------|-------|------------|
| Storing OTP in plaintext | OTP readable if DB is compromised | Phase 1 | Hash the OTP with bcrypt (low cost factor 4) or SHA-256 before storing |
| Returning different error messages for "user not found" vs "wrong OTP" | Enumeration attack — reveals which phone numbers are registered | Phase 1 | Return the same error for both: `"کد تأیید نامعتبر است"` |
| JWT secret in code | Key compromise | Phase 1 | Load from environment variable; use minimum 256-bit random key |
| No request size limit on API | DoS via oversized payloads | Phase 1 | Set `Fiber`/`Echo` body size limit: 10MB for file uploads, 1MB for JSON |
| Serving uploaded files with original Content-Type | XSS via uploaded HTML/SVG files | Phase 5 | Force `Content-Disposition: attachment`; set `Content-Type: application/octet-stream` for downloads |
| Missing CSRF protection on auth endpoints | Session fixation on OTP verify | Phase 1 | Use `SameSite=Strict` on cookies; validate `Origin` header on mutation requests |

---

## UX Pitfalls

| Pitfall | Impact | Prevention |
|---------|--------|------------|
| No loading state during nutrition computation | Builder appears frozen while computing 100+ item totals | Show skeleton/spinner on nutrition summary cards; compute async |
| Optimistic UI without rollback | Water log shows "+200ml" immediately but sync fails silently | Show tentative indicator on optimistic entries; revert on sync failure with toast |
| Modal-in-modal for food picker in plan builder | Food picker opens a modal, food details open another — mobile stack overflow | Use a bottom sheet pattern; maximum 1 overlay at a time |
| No empty state for new nutritionists | First login shows empty client list, empty food database — user doesn't know what to do | Show onboarding cards: "Add your first client", "Browse food database" |
| Time picker doesn't default to common values | Every meal/medication time starts at 00:00 — tedious to set | Default to typical times: breakfast 07:30, lunch 12:30, dinner 19:30; medication 08:00/20:00 |
| Persian number keyboard not triggered | Quantity input shows Latin keyboard on some devices | Use `<input type="text" inputmode="decimal">` with pattern validation; don't rely on `type="number"` which behaves inconsistently with Persian numerals |

---

## "Looks Done But Isn't" Checklist

Items that pass happy-path testing but fail in real-world usage:

| Item | Happy Path | Real-World Failure | Test For |
|------|-----------|-------------------|----------|
| Diet plan archival | New plan archives the old one | Two nutritionists create plans for the same client simultaneously — race condition, both become active | Insert two plans concurrently; verify partial unique index constraint fires |
| OTP expiry | OTP expires after 2 minutes | User's phone clock is 3 minutes ahead — OTP appears expired immediately | Use server time exclusively; never trust client time for expiry |
| Offline food log | Works offline with active plan cached | Plan expires while offline (end_date passes) — user logs food against a plan that no longer applies | Allow logging against expired plans; resolve on sync by flagging as "logged during expired plan period" |
| File download | PDF downloads correctly in Chrome | iOS Safari opens PDF inline instead of downloading — user can't find the file later | Test on actual iOS Safari; provide explicit "Save to Files" instruction |
| Sleep log across midnight | Sleep 23:00–07:00 same date | Sleep 23:00 Tuesday → 03:00 Wednesday — which date does it belong to? | Assign to the date the user STARTED sleeping (the `date` field); document this convention |
| Soft-delete cascade | Deactivating a food item hides it from search | A diet plan references a deactivated food item — client sees "Unknown food" | Soft-deleted food items remain visible in existing plans; only hidden from new searches |
| Message read receipts | Marking as read works | User scrolls past a message without reading it — `read_at` is set on scroll, not on actual reading | Accept this limitation; mark as read when the conversation view is opened, not per-message scroll |

---

## Recovery Strategies

When a pitfall is discovered in production:

| Scenario | Immediate Action | Root Fix | Timeline |
|----------|-----------------|----------|----------|
| Duplicate tracking entries from sync | SQL: deduplicate by `local_id`, keep latest `created_at` | Fix sync queue with Web Locks API | 1–2 days |
| Persian search returns no results | Switch to `ILIKE '%term%'` as emergency fallback (slow but works) | Fix locale, rebuild trigram index | 1 day |
| JWT refresh causing mass logouts | Extend access token to 1 hour temporarily | Implement refresh queue in auth interceptor | 2–3 days |
| Diet plan cached stale | Add `?bust=${Date.now()}` to plan fetch URL as hotfix | Implement proper ETag + IndexedDB strategy | 1–2 days |
| Push notification duplicates | Disable reminder scheduler; send manually | Add `processed_reminders` deduplication table | 1 day |
| File upload vulnerability discovered | Disable file upload endpoints immediately | Fix content validation, path sanitization, redeploy | Same day |

---

## Pitfall-to-Phase Mapping

| Phase | Critical Pitfalls to Address | Moderate Pitfalls | Minor Pitfalls |
|-------|------------------------------|-------------------|----------------|
| **Phase 1: Foundation** | #6 (auth pattern), #7 (JWT refresh) | #8 (Shamsi composable), #14 (Nuxt 4 dirs) | #15 (RTL setup), #16 (OTP normalization), #17 (pool config) |
| **Phase 2: Core Data** | #2 (Persian search), #3 (char normalization) | | |
| **Phase 3: Diet Plan** | #1 (N+1 queries) | #11 (state management) | |
| **Phase 4: Tracking** | | #8 (date boundaries) | |
| **Phase 5: Communication** | #10 (file upload security) | #18 (polling battery) | |
| **Phase 6: Offline & PWA** | #4 (iOS eviction), #5 (sync duplicates) | #9 (stale cache), #12 (Dexie versioning), #13 (push duplicates) | |
| **Phase 7: Hardening** | #6 (auth audit — re-verify all endpoints) | #1 (query performance — load test) | All security items |

---

## Sources

| Source | Type | Used For |
|--------|------|----------|
| PostgreSQL pg_trgm docs (postgresql.org/docs/current/pgtrgm.html) | Official docs via Context7 | Persian text search, trigram behavior |
| Dexie.js docs (dexie.org) | Official docs via Context7 | IndexedDB schema versioning, offline sync patterns |
| vite-pwa/vite-plugin-pwa docs (GitHub) | Official docs via Context7 | Service worker strategies, caching, registerType |
| Nuxt 4 docs (nuxt.com/docs/4.x) | Official docs via Context7 | Directory structure changes, useFetch composable |
| pgx v5 docs (pkg.go.dev/github.com/jackc/pgx/v5) | Official docs via Context7 | Batch queries, connection pool configuration |
| golang-jwt/jwt docs (GitHub) | Official docs via Context7 | JWT signing, security notices |
| NutriTrack PRD (project file) | Project documentation | Domain requirements, data model, feature specs |
| NutriTrack docs/phases.md (project file) | Project documentation | Phase structure, risk heatmap, implementation guidance |
