# Project Research Summary

**Project:** NutriTrack Client
**Domain:** Frontend-only Persian RTL nutrition management PWA
**Researched:** 2026-04-22
**Confidence:** HIGH

## Executive Summary

NutriTrack Client should be built as a single Nuxt 4 mobile-first PWA with strong role partitioning between client, nutritionist, and super admin surfaces, but with offline durability limited to the client experience. The product is not a generic admin panel. It is a daily-use Persian RTL care tool where the client must be able to open the app, understand today's diet actions immediately, log adherence quickly with one hand, exchange messages with a nutritionist, and trust what happens when connectivity drops.

The research is consistent about the recommended frontend posture: Nuxt 4 as the app shell, Tailwind 4 for a bespoke Persian mobile UI, Pinia for session and UI state, TanStack Query for server-state lifecycle, Dexie for client offline persistence, and a custom PWA setup with an explicit outbox and sync model. The architecture should isolate offline behavior as a bounded subsystem for client routes only. Nutritionist and admin flows should stay online-first. That keeps v1 aligned with the PRD instead of spreading offline complexity into operational surfaces that do not require it.

The main implementation risk is sequencing. If the team ships screens before it establishes role isolation, auth refresh discipline, cache boundaries, and a durable offline queue, later phases such as tracking, chat, and notifications will require rework. The roadmap should therefore front-load app shell, RTL design tokens, auth/session handling, and offline data contracts before feature-heavy client flows.

## Key Findings

### Recommended Stack Summary

The recommended stack is a frontend-only Nuxt 4.4.x application with Tailwind CSS 4.2.x and strict TypeScript, shaped around mobile PWA constraints rather than desktop-first layouts. Pinia should own session state, role state, UI state, install/update prompts, and sync metadata. TanStack Query should own API-backed server state. Dexie should own durable offline data, outbox items, and replay metadata for client-only flows.

This stack is a direct fit for the documented API contract: custom REST endpoints, JWT access and refresh behavior, OTP login for clients, multipart uploads, and push subscription endpoints. The research strongly rejects OAuth-centric auth frameworks, Pinia-only async caching, localStorage-based offline persistence, and generic generated service-worker setups that do not support selective runtime caching and explicit sync behavior.

**Core technologies:**
- Nuxt 4.4.x: app shell, routing, SSR-aware data flows, and PWA host framework.
- Tailwind CSS 4.2.x: custom mobile UI system with Persian RTL design tokens rather than a heavy component suite.
- Pinia: auth/session state, role state, app-level UI state, and sync/read-model state.
- TanStack Query: server-state fetching, cache freshness, invalidation, and deduplication.
- Dexie + IndexedDB: durable client offline cache, queued writes, and sync bookkeeping.
- @vite-pwa/nuxt with custom service worker behavior: installability, safe runtime caching, update prompts, and push integration.
- Zod + vee-validate: form validation and API boundary validation for OTP, tracking, messaging, settings, and plan authoring flows.
- Persian UX utilities: Vazirmatn, Jalali-aware date formatting, Persian numerals, and RTL-safe primitives.

### V1 Table-Stakes Summary

V1 is credible only if it behaves like a real daily-use care product on mobile. The client experience must revolve around today's plan, fast tracking, messaging, notifications, and transparent offline behavior. The nutritionist experience must support client lookup, client context, messaging, and plan management without forcing desktop assumptions into a narrow-screen UI.

**Must have for launch:**
- Role-based app shell with separate client, nutritionist, and admin route spaces.
- Client OTP login plus staff email/password login with stable session refresh and logout behavior.
- Plan-first client dashboard showing today's meals, water, medication, and pending actions.
- Full active-plan viewing with readable hierarchy for meals, options, exercise, prescriptions, and notes.
- Fast daily tracking for food, water, sleep, exercise, medication, and body measurements.
- Client offline read/write support with visible sync states for core tracking and message drafting.
- Messaging with unread state, attachments, and stable polling behavior.
- Push notifications and reminder preferences.
- Nutritionist client roster, client profile workspace, and basic diet plan authoring/editing.
- Food discovery/request flows and minimal lab-result upload/view behavior.

### Differentiators To Defer

The research is clear that v1 differentiation should come from product quality, not breadth. Persian-native information hierarchy, sync visibility, and one-thumb interaction patterns matter more than broad admin tooling or speculative intelligence features. That said, the strongest differentiators beyond table stakes should be staged after the foundation and core loop are stable.

**Defer until after the core v1 loop is proven:**
- Explicit sync center with detailed queue inspection and manual recovery controls beyond basic visible status.
- Rich “Today” timeline polish, plan-diff summaries, and deeper plan comprehension aids.
- Compact progress storytelling and adherence narratives beyond lightweight trends.
- Guided nutritionist authoring enhancements such as reusable defaults and advanced draft workflows.
- Advanced analytics dashboards, real-time chat transport, client-created custom foods, and automation builders.
- AI recommendations, food-photo logging, wearable integrations, payments, and broad admin reporting.

### Key Architecture Sequencing Implications

The architecture research strongly suggests sequencing by dependency, not by screen count. The app should be structured into a platform layer, role-based route areas, bounded domain modules, and reusable presentation primitives. Client offline capability must remain a bounded subsystem that only client routes can depend on. That means the roadmap should establish route partitions, middleware, repositories, local database schema, and sync policy before broad client tracking work begins.

**Major architectural implications:**
1. Foundation work must define role-specific route areas, layouts, middleware, and store reset rules before feature pages proliferate.
2. Repositories and an API wrapper should be introduced early so feature modules do not embed raw network logic and later fight offline requirements.
3. The offline subsystem needs explicit contracts for queued writes, durable local IDs, reconcile states, and feature-by-feature capability rules before tracking and messaging land.
4. Staff surfaces should share the design system and repository layer, but remain online-first and excluded from offline queue complexity.
5. Persian RTL behavior, Jalali rendering, typography, and mobile safe-area patterns belong in the design-system phase, not a final localization pass.

### Top Implementation Risks To Watch

1. **Treating offline as generic caching**: avoid this by defining a per-feature offline capability matrix and record-level sync states before building client flows.
2. **Relying on Background Sync as the main replay path**: use app-open, reconnect, and manual retry as the primary sync model; keep Background Sync as an enhancement only.
3. **Optimistic writes without durable reconciliation**: assign stable local IDs, persist outbox metadata, and test duplicate/conflict scenarios before broad tracking rollout.
4. **Auth refresh races under polling and role changes**: implement a single-flight refresh path, pause background fetchers during refresh/logout, and reset role-scoped stores on account switch.
5. **Role isolation implemented as hidden navigation only**: partition routing, layouts, persistence, and cache keys by role and user identity from the start.
6. **Overcaching authenticated data in the service worker**: cache static assets aggressively but whitelist sensitive authenticated resources carefully and clear user-scoped storage on logout.
7. **Treating Persian RTL as a CSS toggle**: lock in Persian numerals, Jalali display rules, bidi handling, and narrow-screen interaction conventions in the design system early.

## Implications for Roadmap

Based on the combined research, the safest roadmap is a dependency-first sequence that stabilizes platform behavior before feature-heavy surfaces.

### Phase 1: Foundation, RTL System, and PWA Boundaries

**Rationale:** Every later phase depends on a correct mobile shell, Persian RTL design tokens, route partitioning, and safe cache scope.

**Delivers:**
- Nuxt 4 app skeleton with Tailwind 4, Pinia, TypeScript strict mode, and role-aware layouts.
- Persian-only RTL design system baseline: typography, spacing, icons, form primitives, Jalali/persian formatting helpers, safe-area mobile shell.
- Initial PWA plumbing, manifest, install/update prompts, and conservative cache rules limited to app shell and static assets.

**Addresses:** mobile-first Persian UX, role-based shell expectations, PWA installability.

**Must avoid:** overcaching authenticated data, postponing RTL correctness, and collapsing all roles into one generic shell.

### Phase 2: Auth, Session, and Role Isolation

**Rationale:** Feature work will fail or be expensive to unwind if token refresh, role routing, and store partitioning are unstable.

**Delivers:**
- Client OTP login and staff email/password login flows.
- Shared API client wrapper, refresh logic, session bootstrap, logout, and route middleware.
- Role-scoped stores, cache namespaces, and account-switch reset policies.

**Addresses:** authentication table stakes and clean entry into client, nutritionist, and admin areas.

**Must avoid:** concurrent refresh storms, role leakage across sessions, and polling logic that ignores auth state.

### Phase 3: Client Offline Data Model and Sync Engine

**Rationale:** Client tracking and client messaging both depend on durable local persistence, queue state, and reconcile semantics.

**Delivers:**
- Dexie schema, outbox ledger, sync orchestrator, reconnect/manual retry flows, and record-level sync statuses.
- Offline capability matrix for readable client data, queueable writes, and intentionally online-only actions.
- Repository contracts that can switch between remote write-through and queued local persistence.

**Addresses:** mandatory client offline support and sync visibility.

**Must avoid:** generic offline banners with no durable queue, Background Sync dependence, and missing local identifiers.

### Phase 4: Client Core Loop

**Rationale:** Once the offline foundation is real, the product can safely add the main daily client workflow.

**Delivers:**
- Client home or “Today” surface.
- Active-plan viewing and clear plan hierarchy.
- Fast tracking flows for food, water, sleep, exercise, medication, and body measurements.
- Basic progress/history surfaces needed to confirm recent activity.

**Addresses:** the primary client value proposition and most critical v1 table stakes.

**Must avoid:** form-heavy tracking UX, hidden sync failure states, and plan screens that mirror backend nesting instead of human-readable structure.

### Phase 5: Messaging, Attachments, and Notifications

**Rationale:** Communication depends on stable session handling and benefits from the same queue/retry discipline established for client core flows.

**Delivers:**
- Conversation list and detail views, unread counts, disciplined polling strategy, and queued drafts where appropriate.
- Attachment handling with explicit online-only constraints and preview/progress states.
- Push subscription flows and notification preference management triggered from meaningful user actions.

**Addresses:** chat, unread state, care communication loop, and reminder expectations.

**Must avoid:** battery-heavy polling, early notification permission prompts, and treating uploads like offline-safe JSON mutations.

### Phase 6: Nutritionist Workspace and Operational Surfaces

**Rationale:** The staff experience can now reuse a stable design system, session model, and repository layer without inheriting client offline complexity.

**Delivers:**
- Nutritionist client list and client profile workspace.
- Basic plan authoring/editing, client tracking review, foods/medications search and management, food-request handling.
- Minimal super-admin operations for nutritionist management and core catalog governance.

**Addresses:** nutritionist v1 table stakes and essential admin operations.

**Must avoid:** desktop-style dense CRUD layouts, accidental offline support creep into staff flows, and overbuilding admin depth before core adoption is proven.

### Phase Ordering Rationale

- Foundation precedes auth because role layouts, RTL primitives, and cache boundaries influence how auth and route guards are wired.
- Auth precedes offline because repository contracts, user scoping, and privacy-safe persistence depend on a stable session model.
- Offline precedes client tracking and messaging because both need the same durable IDs, queue states, and reconcile logic.
- Client core loop precedes staff depth because the product promise centers on daily client adherence and offline reliability.
- Messaging and notifications come after the outbox model is proven so retry and permission behavior remain coherent.
- Staff surfaces come last because they are online-first and can be layered onto established platform and domain contracts with less rewrite risk.

### Research Flags

**Phases likely needing deeper research during planning:**
- **Phase 3:** browser-specific service worker, sync, and storage behavior need implementation-level validation on Safari-class mobile browsers.
- **Phase 5:** notification permission timing, push UX, and mobile polling cadence should be validated carefully against target-device behavior.
- **Phase 6:** nutritionist mobile plan-authoring interaction patterns need focused UX exploration to avoid cramped, hostile authoring flows.

**Phases with standard patterns and lighter research needs:**
- **Phase 1:** Nuxt 4, Tailwind 4, Pinia setup, and PWA baseline patterns are well documented.
- **Phase 2:** custom auth wrapper and role middleware are straightforward once the API contract is accepted as fixed.
- **Phase 4:** client tracking flows are domain-complex but architecturally standard once offline contracts are in place.

## Confidence Assessment

| Area | Confidence | Notes |
|------|------------|-------|
| Stack | HIGH | The Nuxt 4, Tailwind 4, Pinia, Dexie, PWA, and validation recommendations are strongly aligned with official ecosystem documentation and the project's API constraints. |
| Features | HIGH | Table stakes and deferrals are well supported by the PRD, PROJECT scope, and API contract. |
| Architecture | HIGH | The bounded offline subsystem, repository-first data access, and role-partitioned area model fit both the documented roles and product promises. |
| Pitfalls | HIGH | The major risks are product-specific and repeatedly reinforced by the PRD scope, API behavior, and current browser support constraints. |

**Overall confidence:** HIGH

### Gaps To Address During Planning

- Exact refresh-token persistence strategy depends on backend cookie/header behavior and should be validated before auth implementation is finalized.
- Safari-class behavior for PWA install prompts, push support expectations, and storage persistence needs explicit device testing during platform planning.
- The exact depth of nutritionist mobile plan authoring should be narrowed in planning so v1 delivers usable authoring without collapsing into a desktop CRUD port.
- The split between basic progress views and deferred analytics should be defined in requirements to prevent metric/dashboard scope creep.

## Sources

### Primary

- `.planning/PROJECT.md` — project scope, constraints, stack, and non-goals.
- `.planning/research/STACK.md` — technology recommendations and stack exclusions.
- `.planning/research/FEATURES.md` — table stakes, differentiators, and anti-features.
- `.planning/research/ARCHITECTURE.md` — recommended layering, ownership boundaries, and sequencing.
- `.planning/research/PITFALLS.md` — phase-level risks and mitigation guidance.

### Supporting references aggregated by research

- Official Nuxt 4 documentation — app structure, data fetching, and framework patterns.
- Official Tailwind CSS documentation — Nuxt integration with Tailwind 4.
- Official Pinia documentation — Nuxt integration and store patterns.
- Vite PWA documentation — Nuxt integration and service-worker options.
- MDN documentation — Push API, Notifications API, and Background Sync support constraints.

---
*Research completed: 2026-04-22*
*Ready for roadmap: yes*