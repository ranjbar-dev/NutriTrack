# Research Summary — NutriTrack

**Project:** NutriTrack — Nutrition Management PWA  
**Domain:** B2B2C Nutritionist Practice Management (Persian Market)  
**Researched:** 2025-07-14 – 2025-07-18  
**Overall Confidence:** HIGH

---

## Executive Summary

NutriTrack is a Persian-only, mobile-first Progressive Web App that digitalizes the nutritionist–client workflow in Iran — replacing the current manual process of WhatsApp messaging, Excel diet plans, and paper tracking. The platform serves three roles (Super Admin, Nutritionist, Client) with a deeply nested diet plan engine at its core (Plan → Days → Meals → Options → Items), six daily tracking dimensions (food, water, sleep, exercise, medication, body measurements), and full offline support for clients. The competitive landscape is uniquely favorable: no credible Persian-language competitor exists, making NutriTrack's real competition the manual workflow itself, not other software.

The recommended stack — **Go/Gin + PostgreSQL 16 + Nuxt 4 + Tailwind CSS v4** — is well-validated across all research dimensions. Go's single-binary deployment and Gin's middleware architecture map cleanly to the 3-layer backend pattern (Handler → Service → Repository). PostgreSQL's `pg_trgm` extension solves Persian fuzzy text search without external services. Nuxt 4's `app/` directory structure with Pinia stores and Dexie.js IndexedDB provides the offline-first client architecture. Tailwind v4's native logical properties (`ms-`, `pe-`, `text-start`) eliminate the need for any RTL plugin — a key simplification. The target scale (50 nutritionists, 10K clients, ~500 concurrent users) is comfortably within single-server PostgreSQL + Go capacity, meaning no Redis, no microservices, and no message broker are needed.

The three highest-impact risks are: **(1)** the diet plan engine's N+1 query trap — a naïve implementation produces 500+ queries per plan load, which must be solved with batch loading from day one; **(2)** Persian text handling — `pg_trgm` silently fails with wrong database locale, and Arabic/Persian character variants (`ی`/`ي`, `ک`/`ك`) cause search misses unless normalized at every data boundary; and **(3)** iOS PWA storage eviction — Apple can silently delete IndexedDB data after 7 days of non-use, destroying unsynced client tracking logs. Each risk has a clear prevention strategy, but all three require architectural decisions made early (Phases 1–3) to avoid costly rework later.

---

## Cross-Cutting Themes

Patterns that emerged across multiple research dimensions:

### 1. Persian/RTL Permeates Every Layer
This is not a "skin" concern. Persian text impacts database locale configuration (ARCHITECTURE, PITFALLS), search indexing strategy (STACK — `pg_trgm`), character normalization at storage and query boundaries (PITFALLS #3), date handling via Jalali calendar (PITFALLS #8), font loading (STACK — Vazirmatn), number formatting (FEATURES), CSS logical properties (STACK — Tailwind v4), and Chart.js axis rendering (PITFALLS #15). Every phase must validate RTL/Persian behavior, not just the "UI polish" phase.

### 2. Diet Plan Engine Is the Complexity Nexus
All four research files converge on the diet plan as the riskiest component. ARCHITECTURE details the 5-level nesting and aggregate root loading pattern. PITFALLS warns of N+1 queries (#1) and frontend state management explosion (#11). FEATURES identifies it as the P0 ship-blocking feature with the longest critical path. STACK recommends sqlc + pgx `SendBatch` specifically to handle this query pattern. This is where the project succeeds or fails technically.

### 3. Offline Support Is High-Value but Must Wait
FEATURES ranks offline support as the #1 differentiator (no competitor offers it), but all research agrees it must come after all online features are stable (Phase 6). ARCHITECTURE shows offline wraps every API endpoint's sync path. PITFALLS documents iOS storage eviction (#4), sync queue race conditions (#5), stale cache issues (#9), and Dexie.js versioning traps (#12). Building offline on top of unstable online features compounds every bug.

### 4. Row-Level Authorization Is a Security Pattern, Not Per-Endpoint Logic
PITFALLS (#6) warns that with 30+ endpoints touching client data, at least one will miss ownership verification. ARCHITECTURE recommends repository-level authorization (JOIN to ownership in every query) rather than handler-level checks. This cross-cutting pattern must be established in Phase 1 and validated continuously.

### 5. Keep It Simple — Scale Doesn't Demand Complexity
ARCHITECTURE explicitly confirms: 500 concurrent users is well within single-server PostgreSQL + Go capacity. STACK recommends against Redis, WebSockets, GraphQL, MinIO, and microservices. The PRD's own decisions (polling for chat, local filesystem for files, no real-time features) all reinforce simplicity. Every "should we add X?" question should default to "no" at this scale.

---

## Critical Decisions (Confirmed)

Decisions validated by cross-referencing multiple research dimensions:

| # | Decision | Source | Cross-Validation | Confidence |
|---|----------|--------|-------------------|------------|
| 1 | **sqlc over GORM** for SQL layer | STACK | ARCHITECTURE confirms aggregate loading needs raw SQL; PITFALLS #1 needs explicit query control | HIGH |
| 2 | **pg_trgm for Persian search** (not full-text search) | STACK, PITFALLS | PostgreSQL has no Persian FTS dictionary; pg_trgm + correct locale works; ARCHITECTURE indexes confirm | HIGH |
| 3 | **Tailwind v4 logical properties for RTL** (no plugin) | STACK | PITFALLS #15 confirms RTL issues; Tailwind v4's `ms-`/`pe-`/`text-start` solve them natively | HIGH |
| 4 | **Polling for chat** (not WebSocket) | PROJECT, ARCHITECTURE | PITFALLS #18 adds adaptive polling advice; ARCHITECTURE anti-pattern #5 rejects WebSocket explicitly | HIGH |
| 5 | **Dexie.js for IndexedDB** (standalone, not Dexie Cloud) | STACK | ARCHITECTURE Pattern 4 details sync queue; PITFALLS #5 and #12 identify Dexie-specific traps | HIGH |
| 6 | **pgx/v5 with SendBatch** for plan loading | STACK | ARCHITECTURE Pattern 2 implements 2–3 batch queries; PITFALLS #1 validates the N+1 prevention | HIGH |
| 7 | **No Redis at this scale** | STACK, ARCHITECTURE | 500 concurrent users; OTP/sessions in PostgreSQL; ARCHITECTURE scalability table confirms | HIGH |
| 8 | **Nuxt 4 `app/` directory structure** from day one | STACK | PITFALLS #14 warns Nuxt 3 structure breaks Nuxt 4; ARCHITECTURE project structure confirms | HIGH |
| 9 | **`registerType: 'autoUpdate'`** for PWA service worker | PITFALLS | PITFALLS #9 warns `prompt` default causes stale caches; critical for diet plan freshness | HIGH |
| 10 | **Persian character normalization at every boundary** | PITFALLS | #3 details `ی`/`ي` and `ک`/`ك` problem; must normalize in DB trigger, Go backend, and Nuxt frontend | HIGH |

---

## Risk Register

Top risks ranked by Impact × Likelihood, synthesized from all research:

| ID | Risk | Impact | Likelihood | Mitigation | Phase |
|----|------|--------|------------|------------|-------|
| R1 | **N+1 query explosion on diet plan load** — 500+ queries, multi-second response | Critical | High (without intervention) | Batch load in 2–3 queries with pgx SendBatch; test with query counter middleware; assert ≤5 queries per plan load | Phase 3 |
| R2 | **Persian search silently returns zero results** — wrong locale or missing normalization | Critical | High (default Docker locale is `C`) | Set locale `en_US.UTF-8` in Docker; create `normalize_persian()` DB trigger; seed 50+ real Persian foods; test in Phase 2 | Phase 2 |
| R3 | **iOS PWA storage eviction** — unsynced client data silently destroyed after 7 days | Critical | Medium (iOS-specific) | Aggressive sync on every app open and connectivity change; show pending count badge; minimize cache size; re-fetch on app resume | Phase 6 |
| R4 | **Row-level authorization bypass** — nutritionist accesses another's client data | Critical | Medium (30+ endpoints) | Repository-level ownership JOINs; `authorizeClientAccess()` reusable function; automated cross-nutritionist security tests in CI | Phase 1 (pattern), Phase 7 (audit) |
| R5 | **JWT refresh token race condition** — concurrent 401s cause mass logout | High | High (15-min token + parallel requests) | Auth interceptor with `isRefreshing` flag + retry queue; proactive refresh 1–2 min before expiry; server-side 30s grace period | Phase 1 |
| R6 | **Offline sync queue creates duplicates** — double-submit, concurrent tabs, partial failures | High | Medium | Web Locks API for cross-tab exclusion; debounce UI; sequential queue processing; server-side `ON CONFLICT (local_id) DO NOTHING` | Phase 4 (infra), Phase 6 (full) |
| R7 | **Shamsi date boundary bugs** — wrong "today" detection, timezone mismatch | Medium | High (UTC vs Tehran +3:30) | All storage Gregorian; fixed `Asia/Tehran` timezone in Go; `useShamsiDate` composable wrapping jalaali-js; test midnight boundaries | Phase 1 (composable), Phase 3–4 (validation) |
| R8 | **Diet plan builder UI state explosion** — deeply nested reactive state causes jank | Medium | Medium | Flatten Pinia store (normalized IDs, not nested objects); per-option computed nutrition; `shallowRef` for large lists; debounce recalc | Phase 3 |
| R9 | **Service worker caches stale diet plan** — client follows wrong/outdated plan | High | Medium | Use Dexie as authoritative offline store (not Workbox cache); push notification on plan update; ETag/If-None-Match; `autoUpdate` SW | Phase 6 |
| R10 | **File upload security** — content sniffing XSS, path traversal, storage exhaustion | High | Low (but exploitable) | UUID filenames only; validate magic bytes; `Content-Disposition: attachment`; per-client storage limits; authenticated download endpoints | Phase 5 |

---

## Phase Readiness Assessment

### Phase 1: Foundation
**Research readiness: HIGH — clear patterns, minimal unknowns**

| Clear | Needs Spike/Investigation |
|-------|--------------------------|
| Go project structure (3-layer: handler/service/repository) | Kavenegar SMS API integration — test with actual Iranian SIM |
| Gin route groups with middleware stacking (Pattern 3) | OTP hashing strategy (bcrypt cost 4 vs SHA-256) — minor decision |
| JWT access (15min) + refresh (30d) token flow | — |
| pgxpool configuration (MaxConns: 20–50) | — |
| Nuxt 4 `app/` directory structure | — |
| Tailwind v4 + Vazirmatn + RTL logical properties | — |
| Docker Compose (Go + PG + Traefik) | — |
| `useShamsiDate` composable with jalaali-js | — |
| Auth interceptor with refresh queue pattern (R5) | — |
| Row-level authorization pattern (R4) | — |

**Verdict:** Standard patterns — skip phase research. All technologies documented with Context7.

### Phase 2: Core Data Domain
**Research readiness: HIGH — but Persian search needs early validation**

| Clear | Needs Spike/Investigation |
|-------|--------------------------|
| Food/medication CRUD (3-layer pattern) | Persian `pg_trgm` search with correct locale — **validate immediately** |
| Food categories and measurement units (8 + 12 enums) | `normalize_persian()` function scope — test with real keyboard variants |
| Super Admin panel (Nuxt pages) | Minimum trigram query length (2 chars vs 3 chars for Persian) |
| NutritionLabel, FoodPicker, SearchInput components | — |

**Verdict:** Run a focused spike on `pg_trgm` with Persian data before full Phase 2 implementation. Seed 50+ real Persian food names and validate search queries.

### Phase 3: Diet Plan Engine
**Research readiness: MEDIUM — highest complexity, needs design spike**

| Clear | Needs Spike/Investigation |
|-------|--------------------------|
| Data model (5 tables: plans/days/meals/options/items) | Batch loading query design — prototype the 2–3 query aggregate load |
| Nutrition computation (sum item quantities × food macros) | Plan builder UI architecture — flatten Pinia store vs nested state (R8) |
| Plan lifecycle (create → active → archived) | Repeating day pattern (7-day cycle mapping) — modulo arithmetic with Jalali |
| Partial unique index for one-active-plan constraint | Plan creation payload validation — deeply nested JSON schema |
| `useMealBuilder` composable pattern | Transaction design — single PG transaction for full plan insert |

**Verdict:** **Needs phase research.** Run a design spike on: (a) aggregate loading query + in-memory tree assembly, (b) plan builder state management approach. This is the highest-risk phase.

### Phase 4: Client Tracking
**Research readiness: HIGH — six parallel CRUD domains with shared pattern**

| Clear | Needs Spike/Investigation |
|-------|--------------------------|
| 6 tracking tables with `(client_id, date)` indexes | `local_id` deduplication — `ON CONFLICT DO NOTHING` vs `DO UPDATE` |
| Tracking CRUD (same 3-layer pattern × 6 types) | Chart.js RTL configuration for weight/measurement charts |
| Daily dashboard with summary cards | "Today" detection with Tehran timezone (R7) — validate edge cases |
| `useShamsiDate` for date display | — |

**Verdict:** Standard patterns. The `local_id` dedup strategy is documented; just implement it consistently across all 6 tables.

### Phase 5: Communication Layer
**Research readiness: HIGH — well-documented patterns**

| Clear | Needs Spike/Investigation |
|-------|--------------------------|
| Message CRUD + polling (10s, `since` parameter) | File upload security implementation — content sniffing, path traversal (R10) |
| Chat UI component pattern | Adaptive polling strategy — 10s → 30s → 60s backoff (PITFALLS #18) |
| File upload handler (UUID names, magic byte validation) | BroadcastChannel API for single-tab polling (performance trap) |
| Food request workflow (submit → approve → create food) | — |
| Lab results (file upload + metadata) | — |

**Verdict:** Standard patterns. File upload security needs careful implementation but the patterns are well-documented.

### Phase 6: Offline & PWA
**Research readiness: MEDIUM — multiple interacting systems, iOS-specific risks**

| Clear | Needs Spike/Investigation |
|-------|--------------------------|
| @vite-pwa/nuxt configuration + `autoUpdate` SW type | iOS storage eviction mitigation strategy (R3) — test on real devices |
| Dexie.js schema for offline stores | Cross-tab sync coordination with Web Locks API (R6) |
| Sync queue pattern (ARCHITECTURE Pattern 4) | Service worker caching strategy — Dexie vs Workbox for API data (R9) |
| Push notification via webpush-go + VAPID | Dexie.js schema versioning discipline (PITFALLS #12) |
| Reminder worker (Go goroutine with minute-tick) | Reminder dedup with `processed_reminders` table (PITFALLS #13) |
| Notification preferences | Background Sync API absence on iOS — polling fallback |

**Verdict:** **Needs phase research.** Multiple interacting systems (SW + Dexie + sync queue + push + reminders) with iOS-specific constraints. Test on actual iOS devices before finalizing architecture.

### Phase 7: Hardening & Launch
**Research readiness: HIGH — standard security/performance checklist**

| Clear | Needs Spike/Investigation |
|-------|--------------------------|
| Row-level authorization audit (all 30+ endpoints) | Load testing tool selection (k6 vs wrk vs custom) |
| EXPLAIN ANALYZE on all critical queries | Grafana + Loki dashboard design for NutriTrack metrics |
| Database backup verification | SSL/TLS configuration on Hetzner with Traefik Let's Encrypt |
| GitLab CI/CD pipeline (lint → test → build → deploy) | — |
| Security test suite (cross-nutritionist access) | — |

**Verdict:** Standard patterns. No phase research needed — this is execution against a checklist.

---

## Open Questions

Consolidated from all research docs, deduplicated and prioritized:

| Priority | Question | Source | Resolution Path |
|----------|----------|--------|-----------------|
| **P0** | What is the Kavenegar API contract for OTP delivery, rate limits, and delivery callbacks? | PITFALLS (integration gotchas) | Spike during Phase 1 — test with real Iranian SIM card |
| **P0** | Does `pg_trgm` produce valid trigrams for Persian text with `en_US.UTF-8` locale? | PITFALLS #2, STACK | Spike during Phase 2 — run `SELECT show_trgm('برنج')` and validate |
| **P1** | How should the plan builder UI handle 420+ food items without jank on mid-range Persian mobile devices? | PITFALLS #11, ARCHITECTURE | Design spike during Phase 3 — prototype with normalized Pinia store |
| **P1** | What is the optimal Dexie.js schema for the offline store — mirror server tables or domain-optimized? | ARCHITECTURE (Dexie schema), PITFALLS #12 | Design during Phase 6 — start with minimal schema, expand per tracking type |
| **P2** | How does iOS 17/18 handle PWA storage persistence? Has Apple relaxed the 7-day eviction? | PITFALLS #4 | Research during Phase 6 — test on latest iOS devices |
| **P2** | Should the batch sync API (`POST /api/client/tracking/batch`) exist from Phase 4 or only Phase 6? | PITFALLS (performance traps) | Decide during Phase 4 planning — adding it early is low cost |
| **P2** | What is the initial food database seed set? How many items, which categories? | FEATURES (food database) | Content decision by product owner before Phase 2 |
| **P3** | What VAPID key management strategy prevents key rotation from invalidating all subscriptions? | PITFALLS (Web Push integration) | Design during Phase 6 — document key storage and rotation plan |

---

## Recommendations for Roadmap

### Phase Ordering Rationale

The 7-phase structure proposed by ARCHITECTURE is strongly validated by all research:

```
Phase 1 (Foundation) → Phase 2 (Core Data) → Phase 3 (Diet Plan Engine)
                                                       ↓
                                              ┌────────┴────────┐
                                              ↓                 ↓
                                    Phase 4 (Tracking)   Phase 5 (Communication)
                                              └────────┬────────┘
                                                       ↓
                                              Phase 6 (Offline & PWA)
                                                       ↓
                                              Phase 7 (Hardening)
```

**Why this order:**
1. **Phase 1 → 2 → 3** follows the strict data dependency chain: users → foods → diet plans. No shortcuts possible.
2. **Phases 4 and 5 can parallelize** after Phase 3 — FEATURES confirms they share no data dependencies (tracking tables don't reference messages; messages don't reference tracking).
3. **Phase 6 must follow 4+5** because offline wraps every API endpoint. Building sync before the endpoints exist is wasted work.
4. **Phase 7 is a hardening pass**, not new features. It validates everything built in Phases 1–6.

### Key Ordering Insights

1. **Front-load Persian text infrastructure (Phase 2).** The `normalize_persian()` function, `pg_trgm` index, and database locale must be proven before any real data enters the system. Every subsequent phase depends on text search working.

2. **Phase 3 needs the most calendar time.** FEATURES estimates 3–4 weeks. It contains the most complex backend (aggregate loading, transaction-based plan creation), the most complex frontend (plan builder UI), and the most complex computation (nutrition totals). Schedule a design spike before implementation.

3. **`local_id` dedup infrastructure belongs in Phase 4, not Phase 6.** PITFALLS #5 specifies this explicitly. The `UNIQUE` constraint on `local_id` columns and the `ON CONFLICT` pattern must exist in the tracking tables from their creation, not retrofitted when offline sync is added.

4. **Push notification subscription can start in Phase 5** (collect subscriptions when messaging is built), even though push sending logic is Phase 6. This decouples subscription management from notification delivery.

5. **Security patterns are Phase 1 decisions with Phase 7 validation.** The row-level authorization pattern, JWT refresh queue, and OTP normalization are all Phase 1 architectural choices. Phase 7 runs the comprehensive audit, but the patterns must be correct from day one.

### Research Flags

**Phases needing deeper research during planning:**
- **Phase 3 (Diet Plan Engine):** Design spike needed on aggregate loading query pattern and plan builder UI state management. Highest-complexity phase.
- **Phase 6 (Offline & PWA):** Multiple interacting systems (SW + Dexie + sync queue + push + reminders) with iOS-specific constraints. Test on real iOS devices.

**Phases with standard, well-documented patterns (skip phase research):**
- **Phase 1 (Foundation):** All technologies thoroughly documented via Context7. Gin, pgx, JWT, Nuxt 4 patterns are standard.
- **Phase 2 (Core Data):** CRUD + search. Only the `pg_trgm` Persian validation needs a focused spike (1–2 hours, not full research).
- **Phase 4 (Tracking):** Six parallel CRUD domains following the same 3-layer pattern.
- **Phase 5 (Communication):** Chat polling, file uploads, food requests — all well-established patterns.
- **Phase 7 (Hardening):** Checklist-based execution, no novel research needed.

### Conflicts Identified

One minor inconsistency was found across research dimensions:

| Conflict | STACK says | PITFALLS says | Resolution |
|----------|-----------|---------------|------------|
| Tailwind RTL approach | Use Tailwind v4 logical properties natively (no plugin needed) | Use `tailwindcss-rtl` plugin for RTL | **Use STACK recommendation.** Tailwind v4 has native logical property support (`ms-`, `pe-`, `text-start`). The PITFALLS reference to the RTL plugin appears to be based on Tailwind v3 assumptions. No plugin is needed. |
| pgxpool MaxConns | 20 connections sufficient for 500 concurrent users | 50 connections recommended | **Use 25 as default, scale to 50 if needed.** ARCHITECTURE suggests 25 for 500 concurrent users as the middle ground. Monitor pool utilization in Grafana. |

---

## Confidence Assessment

| Area | Confidence | Notes |
|------|------------|-------|
| Stack | **HIGH** | All core technologies verified via Context7 official docs. Version compatibility matrix validated. Go 1.25+ requirement for Gin confirmed. |
| Features | **MEDIUM** | Table stakes derived from Western competitor analysis + PRD, not primary Iranian market research. Anti-features list HIGH (sourced from PRD). Feature dependencies HIGH (from data model FKs). |
| Architecture | **HIGH** | 3-layer Go pattern is standard. Nuxt 4 directory structure confirmed. All data flows validated against PRD data model. pgx batch loading pattern verified. |
| Pitfalls | **HIGH** | 18 pitfalls documented, most verified via Context7 official docs. Persian-specific pitfalls (search, normalization, Jalali dates) are well-known in Persian software development. |

**Overall confidence: HIGH**

### Gaps to Address

| Gap | Impact | Resolution |
|-----|--------|------------|
| **No Persian market user research** — feature priorities based on Western competitor analysis | Medium | Validate table stakes with 2–3 real nutritionists before Phase 3 |
| **Kavenegar API specifics unknown** — rate limits, delivery guarantees, failover | Medium | Test during Phase 1 development with real SMS delivery |
| **iOS PWA storage eviction behavior on latest iOS** — may have changed in iOS 17/18 | Medium | Test on actual devices during Phase 6 |
| **Persian food database seed data** — no pre-existing Persian nutritional database identified | Low | Product owner must curate initial seed set before Phase 2 |
| **Actual mid-range phone performance** — plan builder UI jank untested | Low | Test Phase 3 output on representative devices (Samsung A-series, Xiaomi) |

---

## Sources

### Primary (HIGH confidence)
- **Gin Framework** — Context7 `/gin-gonic/gin` and `/websites/gin-gonic_en`: Go 1.25 requirement, middleware patterns, route groups
- **pgx/v5** — Context7 `/jackc/pgx` and `/websites/pkg_go_dev_github_com_jackc_pgx_v5`: SendBatch, pgxpool configuration
- **sqlc** — Context7 `/websites/sqlc_dev_en`: pgx/v5 integration, configuration format
- **golang-migrate** — Context7 `/golang-migrate/migrate` (v4.18.3): PostgreSQL setup, Go library usage
- **golang-jwt** — Context7 `/golang-jwt/jwt`: token creation/parsing, custom claims
- **Nuxt 4** — Context7 `/nuxt/nuxt` (v4.0.0, v4.1.3) and `/websites/nuxt_4_x`: `app/` directory, composables, middleware
- **Tailwind CSS v4** — Context7 `/websites/tailwindcss`: RTL logical properties, CSS-first config, PostCSS setup
- **Pinia** — Context7 `/vuejs/pinia`: SSR hydration, store plugins, Nuxt integration
- **@vite-pwa/nuxt** — Context7 `/websites/vite-pwa-org_netlify_app`: service worker lifecycle, registerType, `$pwa` composable
- **Dexie.js** — Context7 `/websites/dexie`: schema versioning, offline-first patterns (standalone)
- **Traefik** — Context7 `/websites/doc_traefik_io_traefik`: Docker Compose setup, ACME, HTTPS
- **PostgreSQL pg_trgm** — Context7: trigram extraction, locale dependency
- **go-playground/validator** — Context7 (v10.27.0): struct tag validation
- **zerolog** — Context7: zero-allocation JSON logging
- **Chart.js** — Context7 `/chartjs/chart.js`: RTL configuration, Vue integration

### Secondary (MEDIUM confidence)
- **Competitive platforms** (training data): Nutrium, Practice Better, Healthie, Foodzilla, That Clean Life, Nutritics — for feature landscape analysis
- **Consumer tracking apps** (training data): MyFitnessPal, Cronometer — for tracking feature benchmarking
- **webpush-go** — repository verified but not deeply documented via Context7

### Tertiary (LOW confidence)
- **Persian market context** — general knowledge of Iranian tech ecosystem; no primary user research conducted
- **iOS PWA storage eviction timing** — documented behavior but may vary by iOS version
- **Kavenegar API behavior under load** — known provider but rate limits not verified

---
*Research completed: 2025-07-18*  
*Ready for roadmap: yes*
