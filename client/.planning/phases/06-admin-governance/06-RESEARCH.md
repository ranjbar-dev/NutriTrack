# Phase 6: Admin Governance - Research

**Researched:** 2026-04-24
**Domain:** Nuxt 4 super-admin governance surfaces (stats, nutritionist lifecycle, elevated catalogue control)
**Confidence:** HIGH

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions
- **D-01:** Admin governance remains strictly online-first; Phase 3 offline queueing does not extend into `/admin/**` flows.
- **D-02:** All new super-admin pages stay under `/admin/**` and continue to rely on deny-by-default role guards from Phase 2.
- **D-03:** Admin catalogue operations must use elevated admin endpoints where available instead of reusing nutritionist-scoped ownership assumptions.
- **D-04:** Platform stats are presented as compact mobile KPI surfaces backed by `GET /admin/stats`, not derived from stitched client-side queries.
- **D-04a:** The admin stats UI must only show API-backed metrics; because `GET /admin/stats` exposes `total_clients` but not active/inactive client splits, Phase 6 should show total clients only unless the API contract changes.
- **D-05:** Nutritionist management covers create, list, detail, update, activate/deactivate, and read-only client-list visibility via the documented `/admin/nutritionists*` endpoints.
- **D-06:** Nutritionist status changes are explicit confirmation actions with clear Persian success/error states because they directly affect access control.
- **D-07:** Super-admin food governance uses `GET /admin/foods` and `DELETE /admin/foods/:id` for elevated search and force-delete, while create/edit can continue through shared food CRUD where the API already grants `super_admin` authority.
- **D-08:** Super-admin medication governance uses `GET /admin/medications` and `DELETE /admin/medications/:id` for elevated search and force-delete, while create/edit continue through shared medication CRUD where permitted to `super_admin`.
- **D-09:** Food-category governance is admin-only for create/delete and should be treated as part of the catalogue control surface because category taxonomy materially affects shared food management.
- **D-10:** PRD mentions audit-log visibility for food and medication changes, but `docs/API.md` does not expose an audit-log contract; Phase 6 should not invent that surface and must defer it unless a documented endpoint appears.

### the agent's Discretion
- Mobile information architecture for balancing KPIs, list management, and destructive actions without collapsing into a desktop admin table.
- Search/filter defaults and detail-sheet density so long as Persian RTL mobile ergonomics remain primary.

### Deferred Ideas (OUT OF SCOPE)
- Audit-log timeline UI until the API contract exists.
- Bulk admin actions, CSV export, and advanced analytics beyond the Phase 6 success criteria.
- Desktop-heavy moderation dashboards or multi-pane catalogue tooling.
</user_constraints>

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| ADMIN-01 | Super admin can view platform stats and manage nutritionist accounts from mobile-compatible admin screens. | Workstreams W1-W2 map `/admin/stats` and full `/admin/nutritionists*` lifecycle with mobile RTL card/list/sheet patterns and role guard reuse. [VERIFIED: .planning/REQUIREMENTS.md][VERIFIED: docs/API.md][VERIFIED: app/middleware/role-shell.global.ts] |
| ADMIN-02 | Super admin can manage shared food and medication catalogues with elevated admin endpoints. | Workstream W3 maps admin catalogue search + force-delete + category governance via `/admin/foods`, `/admin/medications`, `/admin/food-categories`, reusing catalogue typing/composable patterns. [VERIFIED: .planning/REQUIREMENTS.md][VERIFIED: docs/API.md][VERIFIED: app/composables/useCatalogueApi.ts] |
</phase_requirements>

## Project Constraints (from copilot-instructions context)

- Workspace scope is frontend-only; Phase 6 must not include backend/database/infrastructure changes. [VERIFIED: .github/copilot-instructions.md]
- Stack baseline is Nuxt 4 + Tailwind 4 + Pinia with strict TypeScript and composable-first API access. [VERIFIED: .github/copilot-instructions.md]
- Product is Persian-only RTL mobile PWA; mobile viewport quality is primary and desktop is secondary. [VERIFIED: .github/copilot-instructions.md][VERIFIED: docs/PRD.md]
- Admin surfaces are online-first; offline queue/caching obligations apply to client flows, not admin. [VERIFIED: .github/copilot-instructions.md][VERIFIED: docs/PRD.md][VERIFIED: .planning/phases/06-admin-governance/06-CONTEXT.md]
- Role boundaries must stay explicit across routing/state/caching for client, nutritionist, and super-admin namespaces. [VERIFIED: .github/copilot-instructions.md][VERIFIED: app/middleware/role-shell.global.ts]

## Summary

Phase 6 is a bounded frontend integration phase with strong contract clarity for stats, nutritionist management, and elevated catalogue governance. The API already defines all mandatory ADMIN-01/ADMIN-02 endpoints, so planning risk is mostly UX orchestration and consistency with existing composable + typed-contract patterns rather than protocol uncertainty. [VERIFIED: docs/API.md][VERIFIED: .planning/REQUIREMENTS.md][VERIFIED: app/composables/useNutritionistClientApi.ts]

The most important planning discipline is separating "documented admin capabilities" from "PRD aspirations not yet contracted". Specifically, audit logs and active/inactive client split stats appear in PRD language but are not represented in current API responses, so they must be explicitly deferred or bounded as non-goals inside PLAN files to avoid hidden scope creep. [VERIFIED: docs/PRD.md][VERIFIED: docs/API.md]

Implementation should be decomposed into four workstreams: API/type/composable foundation, admin dashboard + nutritionist lifecycle UX, elevated catalogue governance UX, and targeted validation coverage (contracts + role routing + destructive-action confirmations). This decomposition aligns with existing phase patterns and test style already used in Phases 2 and 5. [VERIFIED: .planning/ROADMAP.md][VERIFIED: tests/auth/route-access-control.spec.ts][VERIFIED: tests/nutritionist/workspace-api-contracts.spec.ts]

**Primary recommendation:** Plan Phase 6 as API-contract-first Nuxt composables and typed models, then layer mobile RTL admin UI in narrow vertical slices (stats -> nutritionists -> catalogue) with explicit defer notes for non-contracted PRD items. [VERIFIED: docs/API.md][VERIFIED: app/types/catalogue.ts][VERIFIED: app/composables/useCatalogueApi.ts]

## Architectural Responsibility Map

| Capability | Primary Tier | Secondary Tier | Rationale |
|------------|-------------|----------------|-----------|
| Admin stats aggregation display | Browser/Client | API/Backend | Client renders KPI cards; source-of-truth aggregation is already exposed by `GET /admin/stats`. [VERIFIED: docs/API.md] |
| Nutritionist account lifecycle UX | Browser/Client | API/Backend | Client owns form states, confirmations, and role-safe navigation; API owns account persistence and status mutation. [VERIFIED: docs/API.md][VERIFIED: app/pages/nutritionist/clients/index.vue] |
| Read-only nutritionist client list visibility | Browser/Client | API/Backend | Client handles pagination/filter UX; API endpoint `/admin/nutritionists/:id/clients` is canonical data authority. [VERIFIED: docs/API.md] |
| Elevated food/medication governance | Browser/Client | API/Backend | Admin UI invokes elevated search/delete routes; backend enforces super_admin authorization boundary. [VERIFIED: docs/API.md] |
| Role namespace enforcement (`/admin/**`) | Frontend Server (Nuxt middleware/runtime) | Browser/Client | Global middleware maps session roles to route prefixes and prevents cross-namespace traversal. [VERIFIED: app/middleware/role-shell.global.ts] |

## Standard Stack

### Core
| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| nuxt | 4.4.2 (published 2026-03-12) | App/runtime shell, routing, middleware, composables (`useAsyncData`) | Already used across all completed phases and required by project constraints. [VERIFIED: package.json][VERIFIED: npm registry] |
| vue | 3.5.13 | Component/reactivity layer for mobile admin pages and sheets | Existing codebase already on Vue 3 with script setup and Composition API patterns. [VERIFIED: package.json][VERIFIED: app/pages/nutritionist/clients/index.vue] |
| pinia + @pinia/nuxt | 3.0.4 / 0.11.3 (published 2025-11-05) | Shared session/platform state integration | Existing role/session/pwa flows already depend on Pinia integration. [VERIFIED: package.json][VERIFIED: npm registry] |

### Supporting
| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| @vite-pwa/nuxt | 1.1.1 (published 2026-02-06) | Connectivity/update banners in shared shell | Keep existing shell behavior on admin layouts; no new offline queueing logic for admin. [VERIFIED: package.json][VERIFIED: app/layouts/admin.vue][VERIFIED: docs/PRD.md] |
| vitest | 4.1.5 latest (workspace currently 3.2.x) | Unit/contract-style regression tests | Reuse current tests style; no mandatory upgrade in this phase unless planned globally. [VERIFIED: npm registry][VERIFIED: package.json][VERIFIED: tests/nutritionist/workspace-api-contracts.spec.ts] |
| @playwright/test | 1.59.1 installed | Mobile-flow e2e coverage for admin critical paths | Add one lightweight happy-path + one role-guard e2e if phase validation requires browser-level confidence. [VERIFIED: package.json][VERIFIED: playwright.config.ts] |

### Alternatives Considered
| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| Existing composable + `useAsyncData` pattern | Introduce TanStack Query / another data-fetch layer | Adds new dependency and cognitive overhead late in roadmap with limited phase value. [ASSUMED] |
| Current card/list mobile layouts | Desktop admin data grid systems | Conflicts with Persian mobile-first requirement and increases responsive complexity. [VERIFIED: docs/PRD.md][VERIFIED: .github/copilot-instructions.md] |

**Installation:**
```bash
npm install
```

No additional package is required for Phase 6 implementation if current patterns are reused. [VERIFIED: package.json][VERIFIED: app/composables/useCatalogueApi.ts]

## Architecture Patterns

### System Architecture Diagram

```text
Super Admin User
  -> /admin/** route request
    -> role-shell middleware validates namespace and role
      -> Admin page composable calls (/api/v1/admin/*)
        -> API response envelope (data/meta or data/message)
          -> Mobile RTL UI states
             -> KPI cards (stats)
             -> nutritionist roster/detail/status sheets
             -> catalogue governance lists + destructive confirm
```

Data flow is request-driven and online-only for admin: there is no admin write queue or offline replay layer. [VERIFIED: .planning/phases/06-admin-governance/06-CONTEXT.md][VERIFIED: docs/PRD.md]

### Recommended Project Structure
```text
app/
  pages/
    admin/
      index.vue                 # KPI landing and action entry points
      nutritionists/index.vue   # list/create/search
      nutritionists/[id].vue    # detail/update/status + read-only clients
      catalogue/foods.vue       # elevated foods search/force-delete
      catalogue/medications.vue # elevated meds search/force-delete
      catalogue/categories.vue  # admin category create/delete
  components/
    admin/
      ...                       # cards, filters, rows, sheets, confirmations
  composables/
    useAdminStatsApi.ts
    useAdminNutritionistApi.ts
    useAdminCatalogueApi.ts
  types/
    admin.ts                    # stats, nutritionist DTOs, status payloads
```

This mirrors existing phase conventions: typed domain models in `app/types`, composable wrappers in `app/composables`, and page-first vertical slices under role namespace folders. [VERIFIED: app/types/catalogue.ts][VERIFIED: app/composables/useNutritionistClientApi.ts][VERIFIED: app/pages/nutritionist/clients/index.vue]

### Pattern 1: Composable-First API Boundary
**What:** Encapsulate each admin endpoint family in composables with stable keying, query builders, and typed responses. [VERIFIED: app/composables/useCatalogueApi.ts]
**When to use:** All list/detail/dashboard reads and mutations in admin pages. [VERIFIED: app/composables/useNutritionistClientApi.ts]
**Example:**
```ts
const baseUrl = '/api/v1'

export const useAdminStatsApi = () => {
  async function getStats() {
    return useAsyncData('admin-stats', () =>
      $fetch<{ data: AdminStats }>(`${baseUrl}/admin/stats`)
    )
  }

  return { getStats }
}
```
Source pattern: `app/composables/useCatalogueApi.ts`, `app/composables/useNutritionistClientApi.ts`. [VERIFIED: app/composables/useCatalogueApi.ts][VERIFIED: app/composables/useNutritionistClientApi.ts]

### Pattern 2: Mobile Card + Sheet Operation Model
**What:** List pages remain concise; updates/destructive actions run through modal/sheet confirmations with explicit busy/success/error state. [VERIFIED: app/pages/nutritionist/food-requests/index.vue]
**When to use:** Nutritionist activation/deactivation and force-delete food/medication flows. [VERIFIED: .planning/phases/06-admin-governance/06-CONTEXT.md]

### Pattern 3: Namespace Guard as Non-Bypassable Gate
**What:** Preserve global middleware role-prefix checks; admin pages must rely on existing deny-by-default behavior. [VERIFIED: app/middleware/role-shell.global.ts]
**When to use:** Every new `/admin/**` route introduction. [VERIFIED: app/middleware/role-shell.global.ts]

### Anti-Patterns to Avoid
- **Client-side derived stats beyond contract:** Do not infer active/inactive clients from nutritionist client listings to fill missing dashboard fields. [VERIFIED: docs/API.md][VERIFIED: .planning/phases/06-admin-governance/06-CONTEXT.md]
- **Reuse non-admin catalogue endpoint for admin moderation screens:** Admin governance must call elevated `/admin/*` endpoints where available. [VERIFIED: docs/API.md][VERIFIED: .planning/phases/06-admin-governance/06-CONTEXT.md]
- **Desktop table-first layout:** Avoid dense multi-column grids as the primary UX in this mobile-first phase. [VERIFIED: docs/PRD.md]

## Planning Workstream Decomposition

### W1 - Admin Contracts and Composables Foundation
- Define `app/types/admin.ts` for stats, nutritionist list/detail, status payload, and admin client-list DTOs. [VERIFIED: docs/API.md]
- Add composables for `/admin/stats` and `/admin/nutritionists*` with query helpers consistent with existing catalogue/client APIs. [VERIFIED: app/composables/useCatalogueApi.ts][VERIFIED: app/composables/useNutritionistClientApi.ts]
- Extend catalogue composables or add dedicated admin catalogue composable for `/admin/foods`, `/admin/medications`, `/admin/food-categories`. [VERIFIED: docs/API.md]

### W2 - Admin Dashboard + Nutritionist Governance UX
- Build `/admin/index` KPI cards from `GET /admin/stats` only. [VERIFIED: docs/API.md]
- Build nutritionist roster (`/admin/nutritionists`) with search/pagination and create action. [VERIFIED: docs/API.md]
- Build nutritionist detail/status (`/admin/nutritionists/[id]`) with explicit activate/deactivate confirmation and read-only clients panel. [VERIFIED: docs/API.md]

### W3 - Elevated Catalogue Governance UX
- Build admin food governance list using `GET /admin/foods` and `DELETE /admin/foods/:id`. [VERIFIED: docs/API.md]
- Build admin medication governance list using `GET /admin/medications` and `DELETE /admin/medications/:id`. [VERIFIED: docs/API.md]
- Build category governance page for `POST /admin/food-categories` and `DELETE /admin/food-categories/:id`. [VERIFIED: docs/API.md]

### W4 - Validation and Guard Rails
- Contract tests for new composables (endpoint paths, payload shapes). [VERIFIED: tests/nutritionist/workspace-api-contracts.spec.ts]
- Route access tests for `/admin/**` role isolation. [VERIFIED: tests/auth/route-access-control.spec.ts]
- UI tests for destructive confirmations and failure states (status toggle/force-delete). [ASSUMED]

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Data fetching cache/revalidation model | Custom global request cache layer | Existing `useAsyncData` + composable keys | Already established and consistent with current Nuxt codebase patterns. [VERIFIED: app/composables/useCatalogueApi.ts] |
| Role access matrix logic per page | Ad-hoc page-level auth conditionals | Existing global middleware role mapping | Centralized guard reduces divergence and regression risk. [VERIFIED: app/middleware/role-shell.global.ts][VERIFIED: tests/auth/route-access-control.spec.ts] |
| Complex admin data-grid subsystem | New desktop-style table framework | Existing card/list/filter/sheet mobile primitives | Phase requirement is mobile-compatible Persian UI, not desktop operations console. [VERIFIED: docs/PRD.md][VERIFIED: app/pages/nutritionist/clients/index.vue] |

**Key insight:** Phase 6 should maximize reuse of established frontend primitives and contract wrappers; introducing new architectural abstractions this late creates disproportionate regression risk for minimal user value. [ASSUMED]

## Common Pitfalls

### Pitfall 1: PRD/API Contract Drift Hidden in UI Scope
**What goes wrong:** Planning includes audit-log timeline or active/inactive client split cards not returned by current endpoints. [VERIFIED: docs/PRD.md][VERIFIED: docs/API.md]
**Why it happens:** PRD capability language is broader than current contract in sections 5.14 vs endpoint list. [VERIFIED: docs/PRD.md][VERIFIED: docs/API.md]
**How to avoid:** Add explicit "deferred due to missing API contract" note in plan acceptance criteria. [VERIFIED: .planning/phases/06-admin-governance/06-CONTEXT.md]
**Warning signs:** Tasks contain verbs "derive", "infer", or "mock" for admin analytics/audit fields. [ASSUMED]

### Pitfall 2: Accidentally Applying Offline Client Patterns to Admin
**What goes wrong:** Admin actions get queued or treated as retry-safe offline workflows. [VERIFIED: docs/PRD.md]
**Why it happens:** Existing codebase has robust client offline infrastructure and developers overgeneralize it. [VERIFIED: .planning/STATE.md]
**How to avoid:** Keep admin pages online-only with immediate error/retry UX, no IndexedDB write queue integration. [VERIFIED: .planning/phases/06-admin-governance/06-CONTEXT.md]
**Warning signs:** New admin stores mention queue state, sync replay, or offline mutation flags. [ASSUMED]

### Pitfall 3: Endpoint Privilege Confusion in Catalogue Governance
**What goes wrong:** Admin moderation UI uses `/foods` and `/medications` list endpoints where elevated admin endpoints are required. [VERIFIED: docs/API.md]
**Why it happens:** Existing nutritionist catalogue search composable already targets non-admin routes and is easy to copy. [VERIFIED: app/composables/useCatalogueApi.ts]
**How to avoid:** Separate admin catalogue composable or explicit admin mode path switching with tests. [ASSUMED]
**Warning signs:** No test assertion that admin governance fetches `/api/v1/admin/foods` and `/api/v1/admin/medications`. [ASSUMED]

## Code Examples

Verified reuse patterns from repository:

### Query Builder + `useAsyncData` Keying
```ts
function buildQuery(params: Record<string, string | number | undefined>): string {
  const query = new URLSearchParams()
  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined && value !== '') {
      query.set(key, String(value))
    }
  })
  const serialized = query.toString()
  return serialized ? `?${serialized}` : ''
}
```
Source: `app/composables/useCatalogueApi.ts`. [VERIFIED: app/composables/useCatalogueApi.ts]

### Mobile List Page State Envelope
```ts
const listState = ref({ loading: true, error: '', clients: [] as NutritionistClient[] })

async function refreshRoster() {
  listState.value.loading = true
  listState.value.error = ''
  const { data, error } = await api.listClients(...)
  listState.value.clients = data.value?.data ?? []
  if (error.value) listState.value.error = 'خطا در دریافت لیست مراجعان'
  listState.value.loading = false
}
```
Source: `app/pages/nutritionist/clients/index.vue`. [VERIFIED: app/pages/nutritionist/clients/index.vue]

### Role Namespace Mapping Pattern
```ts
if (role === 'super_admin') {
  return 'admin'
}
```
Source: `app/middleware/role-shell.global.ts`. [VERIFIED: app/middleware/role-shell.global.ts]

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| Monolithic role-agnostic dashboard | Namespace-specific shells and guarded routes (`/client`, `/nutritionist`, `/admin`) | Established by completed Phases 1-2 | Reduces cross-role leakage and simplifies per-role planning. [VERIFIED: .planning/ROADMAP.md][VERIFIED: app/middleware/role-shell.global.ts] |
| Backend-agnostic UI feature promises | Contract-first endpoint-driven implementation | Current roadmap execution style through Phase 5 | Keeps frontend-only scope realistic and testable. [VERIFIED: .planning/STATE.md][VERIFIED: tests/nutritionist/workspace-api-contracts.spec.ts] |

**Deprecated/outdated for this phase:**
- Assuming desktop-first admin table UX as baseline. Mobile-first Persian shell is explicit product direction. [VERIFIED: docs/PRD.md]

## Assumptions Log

| # | Claim | Section | Risk if Wrong |
|---|-------|---------|---------------|
| A1 | Contract tests should cover admin composable endpoint paths as primary phase safety net. | Planning Workstream Decomposition (W4) | Medium - planner may under-specify validation and miss privilege/path regressions. |
| A2 | New data-fetch library introduces unnecessary overhead versus current composable pattern. | Standard Stack / Alternatives | Low - could still succeed but with avoidable migration cost. |
| A3 | Warning-sign heuristics (`derive`/`infer`/`mock`) predict scope drift tasks. | Common Pitfalls | Low - heuristics may miss edge wording but intent remains valid. |

## Resolved Planning Decisions

1. **Nutritionist client-count display on roster**
  - Resolution: Treat client count as optional display metadata only. The Phase 6 roster plan must not require `client_count` for completion because `docs/API.md` does not guarantee that field on the list endpoint. If absent, the roster still ships with identity, status, and navigation actions intact.

2. **Admin catalogue create/edit surface strategy**
  - Resolution: Phase 6 plans dedicated admin governance pages for elevated search, force-delete, and category management, while leaving shared food/medication create-edit forms reusable only where the existing shared CRUD contracts already support `super_admin`. This keeps admin-specific privilege boundaries explicit without inventing new endpoints or a second form system.

## Environment Availability

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| Node.js | Nuxt/Vitest/Playwright toolchain | yes | v24.14.1 | none needed |
| npm | Package + script execution | yes | 11.11.0 | none needed |
| vitest CLI | Unit/contract test runs | yes | 3.2.4 (installed) | use `npm run test:unit` |
| playwright | E2E smoke checks | yes | 1.59.1 | unit-only validation if e2e not required |

Environment checks were executed from workspace shell. [VERIFIED: terminal command output]

**Missing dependencies with no fallback:**
- None. [VERIFIED: terminal command output]

**Missing dependencies with fallback:**
- None. [VERIFIED: terminal command output]

## Validation Architecture

Validation architecture is enabled because `.planning/config.json` sets `workflow.nyquist_validation` to `true`. [VERIFIED: .planning/config.json]

### Test Framework
| Property | Value |
|----------|-------|
| Framework | Vitest 3.2.x currently installed (`latest` 4.1.5 available) |
| Config file | `vitest.config.ts` |
| Quick run command | `npm run test:unit -- tests/auth/route-access-control.spec.ts tests/nutritionist/workspace-api-contracts.spec.ts` |
| Full suite command | `npm run test:unit` |

### Phase Requirements -> Test Map
| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| ADMIN-01 | Stats dashboard + nutritionist lifecycle endpoints are wired correctly | unit/contract | `npm run test:unit -- tests/admin/admin-api-contracts.spec.ts` | no - Wave 0 |
| ADMIN-01 | `/admin/**` route protection remains role-safe | unit | `npm run test:unit -- tests/auth/route-access-control.spec.ts` | yes |
| ADMIN-02 | Admin catalogue governance uses elevated endpoints and delete confirmations | unit/component | `npm run test:unit -- tests/admin/admin-catalogue-governance.spec.ts` | no - Wave 0 |

### Sampling Rate
- **Per task commit:** run targeted admin + auth tests (`npm run test:unit -- <target files>`). [VERIFIED: package.json]
- **Per wave merge:** run full unit suite (`npm run test:unit`). [VERIFIED: package.json]
- **Phase gate:** unit suite green, then optional focused Playwright smoke for admin happy path. [ASSUMED]

### Wave 0 Gaps
- [ ] `tests/admin/admin-api-contracts.spec.ts` - endpoint path + payload contract checks for stats and nutritionist APIs.
- [ ] `tests/admin/admin-catalogue-governance.spec.ts` - verifies `/api/v1/admin/foods|medications|food-categories` usage and delete/status confirmation flows.
- [ ] `tests/admin/admin-dashboard-shell.spec.ts` - validates API-backed metric rendering only (no inferred fields).

## Security Domain

Security domain is required (security enforcement not disabled in config). [VERIFIED: .planning/config.json]

### Applicable ASVS Categories

| ASVS Category | Applies | Standard Control |
|---------------|---------|-----------------|
| V2 Authentication | yes | Reuse existing JWT session bootstrap and role cookie flow; no new auth mechanism in phase. [VERIFIED: app/middleware/role-shell.global.ts][VERIFIED: docs/API.md] |
| V3 Session Management | yes | Preserve current token/role handling and redirect-on-mismatch behavior. [VERIFIED: tests/auth/route-access-control.spec.ts] |
| V4 Access Control | yes | Enforce `/admin/**` namespace + endpoint-level super_admin contracts. [VERIFIED: app/middleware/role-shell.global.ts][VERIFIED: docs/API.md] |
| V5 Input Validation | yes | Client-side validation for admin forms + backend 422 handling from contract. [VERIFIED: docs/API.md] |
| V6 Cryptography | no (new crypto logic) | No phase-specific cryptographic implementation; rely on existing backend JWT/Web security. [VERIFIED: docs/API.md] |

### Known Threat Patterns for this stack

| Pattern | STRIDE | Standard Mitigation |
|---------|--------|---------------------|
| Privilege confusion between nutritionist/admin catalogue paths | Elevation of Privilege | Test and enforce admin endpoint usage in governance screens. [VERIFIED: docs/API.md][ASSUMED] |
| Destructive action misfire (force delete, deactivate) | Tampering | Two-step confirmation + busy-state lock + explicit success/error feedback. [VERIFIED: .planning/phases/06-admin-governance/06-CONTEXT.md][VERIFIED: app/pages/nutritionist/food-requests/index.vue] |
| Cross-namespace route access regression | Spoofing/Elevation | Keep global role middleware and regression tests for redirect behavior. [VERIFIED: app/middleware/role-shell.global.ts][VERIFIED: tests/auth/route-access-control.spec.ts] |

## Sources

### Primary (HIGH confidence)
- `docs/API.md` - Admin stats, nutritionist management, admin food/medication/category endpoint contracts.
- `.planning/phases/06-admin-governance/06-CONTEXT.md` - Locked phase decisions and explicit defer boundaries.
- `.planning/REQUIREMENTS.md` - ADMIN-01 / ADMIN-02 requirement definitions.
- `app/middleware/role-shell.global.ts` - route namespace and role mapping pattern.
- `app/composables/useCatalogueApi.ts`, `app/composables/useNutritionistClientApi.ts` - existing composable/query conventions.
- `tests/auth/route-access-control.spec.ts`, `tests/nutritionist/workspace-api-contracts.spec.ts` - validation style and guard coverage patterns.

### Secondary (MEDIUM confidence)
- `docs/PRD.md` - product-level admin expectations used to detect scope drift against API contracts.

### Tertiary (LOW confidence)
- npm latest-version comparison implications for upgrade timing (no official project policy mandates upgrade in this phase).

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH - validated against workspace `package.json` and npm registry version checks.
- Architecture: HIGH - grounded in existing middleware/composable/page patterns and locked phase context.
- Pitfalls: MEDIUM - two major drift risks are contract-verified, but warning-sign heuristics remain assumption-based.

**Research date:** 2026-04-24
**Valid until:** 2026-05-24