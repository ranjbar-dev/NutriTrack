# Phase 01: Platform Foundation - Research

**Researched:** 2026-04-22
**Domain:** Nuxt 4 platform foundation (frontend-only Persian RTL mobile PWA)
**Confidence:** HIGH

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions
### RTL visual system
- **D-01:** Use a clinical-minimal visual tone for the platform shell: calm, readable, low-chroma UI with strong hierarchy.
- **D-02:** Use Vazirmatn as the base Persian-first typeface for the Phase 1 design system.
- **D-03:** Default to Persian digits and Jalali display in core UI surfaces.

### the agent's Discretion
- Mobile shell navigation pattern (bottom nav, tabs, or mixed) as long as role boundaries remain explicit.
- PWA runtime cache strategy details, as long as sensitive authenticated data is not broadly cached.
- Install and update prompt microcopy/timing details, as long as prompts are clear and non-disruptive.
- Exact role-layout composition per area (client, nutritionist, admin) while preserving strict route isolation.

### Deferred Ideas (OUT OF SCOPE)
- Rich sync center UX and deeper offline diagnostics are deferred to later phases.
- Advanced visual storytelling, analytics depth, and realtime communication are outside Phase 1 scope.
</user_constraints>

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| PLAT-01 | Persian-only RTL mobile app shell with role-aware navigation for client/nutritionist/admin | Route-group isolation, dedicated layouts, and role-scoped navigation config baseline |
| PLAT-02 | PWA install + clear in-app update prompt | `@vite-pwa/nuxt` module with install/update plugin state (`showInstallPrompt`, `needRefresh`) |
| PLAT-03 | Persian typography, numerals, Jalali display, safe-area handling | Vazirmatn baseline, RTL root setup, Jalali formatter utilities, safe-area CSS tokens |
</phase_requirements>

## Project Constraints (from copilot-instructions.md)

- Frontend-only scope; no backend/database/infrastructure changes. [CITED: .github/copilot-instructions.md]
- Use Nuxt 4, Tailwind CSS 4, Pinia, and strict TypeScript/composable-first patterns. [CITED: .github/copilot-instructions.md]
- Treat `docs/API.md` as backend contract and `docs/PRD.md` as product behavior contract. [CITED: .github/copilot-instructions.md]
- Optimize for Persian-only RTL mobile PWA; mobile quality first. [CITED: .github/copilot-instructions.md]
- Client-side offline required for core client flows; nutritionist/admin remain online-first unless explicitly changed. [CITED: .github/copilot-instructions.md]
- Keep role boundaries explicit across routing, state, caching, and persistence. [CITED: .github/copilot-instructions.md]

## Summary

Phase 1 should deliver a strict platform baseline, not feature workflows: Nuxt 4 app shell, Tailwind 4 token system, Pinia platform stores, Persian RTL/Jalali foundation, role-isolated routing/layouts, and conservative PWA install/update behavior. [CITED: .planning/ROADMAP.md] [CITED: .planning/phases/01-platform-foundation/01-CONTEXT.md]

The most important implementation choice is boundary-first architecture: role segmentation and cache boundaries must be established now so later auth/offline features do not require rewrites. [CITED: .planning/research/ARCHITECTURE.md] [CITED: .planning/research/PITFALLS.md]

**Primary recommendation:** Build Phase 1 as a thin but strict platform layer with explicit route/layout partitioning and minimal-risk PWA caching (static assets and shell only), while deferring deeper offline sync mechanics to later phases. [CITED: .planning/phases/01-platform-foundation/01-CONTEXT.md] [CITED: .planning/research/SUMMARY.md]

## Architectural Responsibility Map

| Capability | Primary Tier | Secondary Tier | Rationale |
|------------|-------------|----------------|-----------|
| Nuxt app shell and role layout partition | Browser/Client (Vue UI) | Frontend Server (SSR rendering) | Layout and navigation chrome are rendered client-side but must remain SSR-compatible for hydration stability. [CITED: https://nuxt.com/docs/4.x/getting-started/introduction] |
| Tailwind token baseline | Browser/Client | CDN/Static | Design tokens compile to static CSS assets and are consumed in client rendering. [CITED: https://tailwindcss.com/docs/installation/framework-guides/nuxt] |
| Pinia platform stores baseline | Browser/Client | Frontend Server | Session/UI platform state is client-owned; Nuxt SSR integration is module-supported. [CITED: https://pinia.vuejs.org/ssr/nuxt.html] |
| Persian RTL + Jalali + numerals presentation | Browser/Client | Frontend Server | Locale rendering and bidi behavior are UI concerns with SSR-safe output requirements. [CITED: docs/PRD.md] [CITED: .planning/phases/01-platform-foundation/01-UI-SPEC.md] |
| Role route/layout isolation | Browser/Client | Frontend Server | Route groups and middleware are Nuxt app-level concerns and enforce role shell separation. [CITED: .planning/research/ARCHITECTURE.md] |
| PWA install/update baseline + safe cache boundaries | Browser/Client (SW + UI prompts) | CDN/Static | Service worker/install prompts run client-side; only static assets should be broadly cached in Phase 1. [CITED: https://vite-pwa-org.netlify.app/frameworks/nuxt.html] [CITED: .planning/research/PITFALLS.md] |

## Standard Stack

### Core

| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| nuxt | 4.4.2 | App framework, routing, SSR-compatible shell | Official Nuxt 4 framework baseline for file routing and SSR-friendly data patterns. [VERIFIED: npm registry] [CITED: https://nuxt.com/docs/4.x/getting-started/introduction] |
| tailwindcss | 4.2.4 | Utility CSS + design token implementation | Official Nuxt install guide uses Tailwind 4 + Vite plugin workflow. [VERIFIED: npm registry] [CITED: https://tailwindcss.com/docs/installation/framework-guides/nuxt] |
| @tailwindcss/vite | 4.2.4 | Tailwind integration in Nuxt Vite config | Recommended plugin path in Tailwind Nuxt guide. [VERIFIED: npm registry] [CITED: https://tailwindcss.com/docs/installation/framework-guides/nuxt] |
| pinia | 3.0.4 | Platform/client state stores | Pinia supports Nuxt 3/4 and auto-import flow through module integration. [VERIFIED: npm registry] [CITED: https://pinia.vuejs.org/ssr/nuxt.html] |
| @pinia/nuxt | 0.11.3 | Nuxt module for Pinia SSR integration | Official integration module for Pinia in Nuxt projects. [VERIFIED: npm registry] [CITED: https://pinia.vuejs.org/ssr/nuxt.html] |
| @vite-pwa/nuxt | 1.1.1 | PWA manifest/SW/install/update integration | Official Nuxt PWA module supporting install and refresh prompt states. [VERIFIED: npm registry] [CITED: https://vite-pwa-org.netlify.app/frameworks/nuxt.html] |

### Supporting

| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| vue | 3.5.33 | Runtime used by Nuxt 4 | Inherited by Nuxt; pin explicitly if dependency conflict occurs. [VERIFIED: npm registry] |
| dexie | 4.4.2 | IndexedDB abstraction for later offline phases | Add now only if Phase 1 scaffolds offline storage boundaries; otherwise defer to Phase 3. [VERIFIED: npm registry] [CITED: .planning/research/STACK.md] |
| @vite-pwa/assets-generator | 1.0.2 | PWA icon asset generation | Use if team wants deterministic icon pipeline from day one. [VERIFIED: npm registry] [CITED: https://vite-pwa-org.netlify.app/frameworks/nuxt.html] |

### Alternatives Considered

| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| @vite-pwa/nuxt | Manual SW + manifest wiring | More control but higher setup risk and slower delivery for Phase 1 baseline. [CITED: https://vite-pwa-org.netlify.app/frameworks/nuxt.html] |
| Pinia platform stores | Full server-state library in Phase 1 | Over-scoping Phase 1; server-state concern can be introduced in auth/offline phases. [CITED: .planning/research/SUMMARY.md] |
| Tailwind token baseline | Heavy prebuilt UI framework | Conflicts with locked UI direction (clinical-minimal, Persian-native shell). [CITED: .planning/phases/01-platform-foundation/01-UI-SPEC.md] |

**Installation:**
```bash
npm install nuxt tailwindcss @tailwindcss/vite pinia @pinia/nuxt @vite-pwa/nuxt
```

**Version verification (registry modified timestamps):**
- nuxt 4.4.2 (modified 2026-03-12T11:59:41.082Z) [VERIFIED: npm registry]
- tailwindcss 4.2.4 (modified 2026-04-22T15:33:29.666Z) [VERIFIED: npm registry]
- @tailwindcss/vite 4.2.4 (modified 2026-04-22T15:33:54.141Z) [VERIFIED: npm registry]
- pinia 3.0.4 (modified 2025-11-05T09:25:14.059Z) [VERIFIED: npm registry]
- @pinia/nuxt 0.11.3 (modified 2025-11-05T09:25:18.224Z) [VERIFIED: npm registry]
- @vite-pwa/nuxt 1.1.1 (modified 2026-02-06T10:27:12.488Z) [VERIFIED: npm registry]

## Architecture Patterns

### System Architecture Diagram

```text
User Opens App
   -> Nuxt App Entry (app.vue)
      -> Route Group Resolver (/client | /nutritionist | /admin | /auth)
         -> Role Layout (client.vue / nutritionist.vue / admin.vue / auth.vue)
            -> Shared Platform Layer
               -> Tailwind Token System + RTL Root + Vazirmatn Font
               -> Pinia Platform Stores (session/ui/pwa-shell)
               -> PWA Plugin State ($pwa install/update flags)
            -> Role-Specific Navigation Config
               -> Role Area Landing Screen

Service Worker Lifecycle
   -> register via @vite-pwa/nuxt
      -> static shell precache
      -> runtime cache allowlist (non-sensitive only)
      -> update detected (needRefresh)
      -> in-app update banner
```

### Recommended Project Structure

```text
app/
  app.vue
  assets/css/main.css
  layouts/
    auth.vue
    client.vue
    nutritionist.vue
    admin.vue
  pages/
    auth/
    client/
    nutritionist/
    admin/
  plugins/
    pwa.client.ts
    pinia.ts
  stores/
    platform-ui.ts
    platform-session.ts
    platform-pwa.ts
  composables/
    useRtl.ts
    usePersianFormat.ts

lib/
  design/
    tokens.css
  locale/
    jalali.ts
    numerals.ts
```

### Pattern 1: Role-Isolated Route Groups
**What:** Use top-level route groups with dedicated layouts per role area.
**When to use:** Immediately in Phase 1 before auth logic is complete.
**Example:**
```typescript
// Source: .planning/phases/01-platform-foundation/01-UI-SPEC.md
// Route groups to establish now:
// /client/** -> client layout
// /nutritionist/** -> nutritionist layout
// /admin/** -> admin layout
```

### Pattern 2: Token-First Tailwind Baseline
**What:** Define spacing/color/typography tokens once and consume through utility classes.
**When to use:** In initial `main.css` and first shell components.
**Example:**
```css
/* Source: .planning/phases/01-platform-foundation/01-UI-SPEC.md */
:root {
  --color-bg: #F4F7F6;
  --color-surface: #E6ECE9;
  --color-accent: #0F766E;
  --space-md: 16px;
  --space-lg: 24px;
}
```

### Pattern 3: PWA Prompt-Driven UX
**What:** Surface non-blocking install/update banners from PWA plugin state.
**When to use:** After module registration and shell mount.
**Example:**
```typescript
// Source: https://vite-pwa-org.netlify.app/frameworks/nuxt.html
const { $pwa } = useNuxtApp()
// use $pwa.showInstallPrompt and $pwa.needRefresh to render banners
```

### Anti-Patterns to Avoid
- **Single shared layout across all roles:** causes role leakage and future auth rewrite risk. [CITED: .planning/research/ARCHITECTURE.md]
- **Caching authenticated JSON responses broadly in SW:** risks stale/private data exposure on shared devices. [CITED: .planning/research/PITFALLS.md]
- **RTL as CSS-only toggle:** misses Persian numerals/Jalali/text-direction edge cases. [CITED: docs/PRD.md] [CITED: .planning/phases/01-platform-foundation/01-UI-SPEC.md]

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Nuxt-PWA integration | Custom SW registration/update/install plumbing | `@vite-pwa/nuxt` | Built-in install/update primitives reduce baseline risk. [CITED: https://vite-pwa-org.netlify.app/frameworks/nuxt.html] |
| Pinia SSR wiring | Manual Pinia bootstrapping and hydration glue | `@pinia/nuxt` module | Official Nuxt integration handles common SSR integration paths. [CITED: https://pinia.vuejs.org/ssr/nuxt.html] |
| Tailwind-Nuxt compile pipeline | Manual PostCSS plugin chain for Tailwind v4 | `@tailwindcss/vite` in Nuxt config | Official current path; less config drift. [CITED: https://tailwindcss.com/docs/installation/framework-guides/nuxt] |
| PWA install/update event state | Custom low-level browser event layer everywhere | `$pwa` reactive properties | Consistent app-level state for prompt UX. [CITED: https://vite-pwa-org.netlify.app/frameworks/nuxt.html] |

**Key insight:** In Phase 1, hand-rolled plumbing increases risk without adding product value; effort should go into Persian RTL quality, role boundaries, and safe cache decisions. [CITED: .planning/research/SUMMARY.md]

## Common Pitfalls

### Pitfall 1: Over-caching authenticated payloads in Phase 1 SW
**What goes wrong:** Role/user data can appear stale or leak on shared devices.
**Why it happens:** Generic URL-prefix cache rules are applied too early.
**How to avoid:** Precache app shell/static assets only in Phase 1; defer authenticated runtime caching policy hardening to later phases.
**Warning signs:** `/api` endpoints added to broad stale-while-revalidate policy.
[CITED: .planning/research/PITFALLS.md]

### Pitfall 2: Route shells not isolated from day one
**What goes wrong:** Client/nutritionist/admin nav and state become tangled.
**Why it happens:** Team starts with one shared dashboard shell for speed.
**How to avoid:** Create role-specific route roots and layouts before feature pages.
**Warning signs:** Shared layout renders cross-role placeholders.
[CITED: .planning/research/ARCHITECTURE.md] [CITED: .planning/phases/01-platform-foundation/01-UI-SPEC.md]

### Pitfall 3: RTL/Jalali fidelity deferred as "later polish"
**What goes wrong:** Deep UI corrections are needed across forms, lists, and date displays.
**Why it happens:** Initial work focuses only on `dir=rtl`.
**How to avoid:** Ship numeral/date/typography helpers in Phase 1 platform layer.
**Warning signs:** Mixed LTR/RTL text wrapping glitches and inconsistent digit rendering.
[CITED: docs/PRD.md] [CITED: .planning/phases/01-platform-foundation/01-UI-SPEC.md]

## Code Examples

Verified patterns from official sources:

### Nuxt + Tailwind 4 baseline
```typescript
// Source: https://tailwindcss.com/docs/installation/framework-guides/nuxt
import tailwindcss from '@tailwindcss/vite'

export default defineNuxtConfig({
  modules: ['@pinia/nuxt', '@vite-pwa/nuxt'],
  css: ['./app/assets/css/main.css'],
  vite: {
    plugins: [tailwindcss()]
  }
})
```

### Nuxt data-fetch baseline guidance
```typescript
// Source: https://nuxt.com/docs/4.x/api/utils/dollarfetch
// Prefer useFetch or useAsyncData + $fetch for SSR-friendly data transfer.
const { data } = await useFetch('/api/item')
```

### PWA install/update reactive state
```typescript
// Source: https://vite-pwa-org.netlify.app/frameworks/nuxt.html
const { $pwa } = useNuxtApp()

const showInstall = computed(() => $pwa?.showInstallPrompt)
const showUpdate = computed(() => $pwa?.needRefresh)
```

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| Nuxt 2-era HTTP modules for app data | Nuxt 4 `useFetch` / `useAsyncData` + `$fetch` | Nuxt 3/4 era | Better SSR transfer and reduced duplicate fetch behavior. [CITED: https://nuxt.com/docs/4.x/api/utils/dollarfetch] |
| Tailwind legacy Nuxt module patterns | Tailwind 4 with `@tailwindcss/vite` plugin | Tailwind v4 docs | Simpler, current integration path for Nuxt projects. [CITED: https://tailwindcss.com/docs/installation/framework-guides/nuxt] |
| Manual PWA event plumbing | `@vite-pwa/nuxt` plugin with `$pwa` state | Current module docs | Faster install/update UX baseline delivery. [CITED: https://vite-pwa-org.netlify.app/frameworks/nuxt.html] |

**Deprecated/outdated:**
- Broad reliance on Background Sync as primary replay mechanism is not cross-browser safe (notably unsupported in Safari/Firefox families). [CITED: https://developer.mozilla.org/en-US/docs/Web/API/Background_Synchronization_API]

## Assumptions Log

| # | Claim | Section | Risk if Wrong |
|---|-------|---------|---------------|
| A1 | `@vite-pwa/nuxt` Nuxt 3 docs behavior maps 1:1 to Nuxt 4 for required install/update primitives. [ASSUMED] | Standard Stack, Code Examples | Minor API mismatch could require small config adaptation |
| A2 | Phase 1 should include only platform-level Pinia stores (not domain stores). [ASSUMED] | Architecture Patterns | Slight plan refactor if team wants early domain placeholders |

## Open Questions

1. Should Phase 1 include only shell-level store stubs, or also minimal auth/session placeholder contracts for Phase 2 handoff?
   - What we know: phase goal is platform baseline, not full auth. [CITED: .planning/ROADMAP.md]
   - What's unclear: desired depth of future-proof scaffolding in this phase.
   - Recommendation: include typed store interfaces and empty actions only.

2. Should PWA icon generation be automated now or deferred?
   - What we know: module supports asset tooling and manifest components. [CITED: https://vite-pwa-org.netlify.app/frameworks/nuxt.html]
   - What's unclear: design team readiness of final icon set.
   - Recommendation: wire manifest and placeholders now; finalize brand assets in later design pass.

## Environment Availability

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| node | Nuxt/Tailwind build | Yes | v24.14.1 | — |
| npm | Package install/scripts | Yes | 11.11.0 | pnpm available |
| pnpm | Optional package manager | Yes | 10.33.0 | npm |
| git | source control workflow | Yes | 2.54.0.windows.1 | — |

**Missing dependencies with no fallback:**
- None. [VERIFIED: local environment commands]

**Missing dependencies with fallback:**
- None. [VERIFIED: local environment commands]

## Validation Architecture

### Test Framework

| Property | Value |
|----------|-------|
| Framework | Not initialized yet in repository (Phase 1 greenfield) |
| Config file | none - see Wave 0 |
| Quick run command | `npm run test:unit` (to add in Wave 0) |
| Full suite command | `npm run test && npm run test:e2e` (to add in Wave 0) |

### Phase Requirements -> Test Map

| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| PLAT-01 | Role-isolated Persian RTL shell routes/layouts render | unit/integration | `npm run test:unit -- shell-role-isolation` | No - Wave 0 |
| PLAT-02 | PWA install/update prompt surfaces correctly | integration/e2e | `npm run test:e2e -- pwa-prompt` | No - Wave 0 |
| PLAT-03 | Persian typography/numerals/Jalali/safe-area baseline | unit/visual-smoke | `npm run test:unit -- locale-baseline` | No - Wave 0 |

### Sampling Rate
- **Per task commit:** run targeted unit test for touched platform component/store.
- **Per wave merge:** run full unit suite + one mobile-shell e2e smoke.
- **Phase gate:** all PLAT-01/02/03 mapped tests green before `/gsd-verify-work`.

### Wave 0 Gaps
- [ ] `package.json` scripts for `test:unit`, `test:e2e`, and `test`
- [ ] `vitest.config.ts` for Nuxt component/store tests
- [ ] `playwright.config.ts` for mobile viewport PWA smoke
- [ ] `tests/platform/shell-role-isolation.spec.ts`
- [ ] `tests/platform/pwa-update-prompt.spec.ts`
- [ ] `tests/platform/persian-locale-baseline.spec.ts`

## Security Domain

### Applicable ASVS Categories

| ASVS Category | Applies | Standard Control |
|---------------|---------|-----------------|
| V2 Authentication | No (Phase 2 focus) | Route/layout prep only in Phase 1 |
| V3 Session Management | No (Phase 2 focus) | Session logic deferred |
| V4 Access Control | Yes | Route-group and layout isolation by role namespace |
| V5 Input Validation | Yes | Minimal schema checks for platform config/prompt state where input exists |
| V6 Cryptography | No | No crypto implementation in Phase 1 |

### Known Threat Patterns for this stack

| Pattern | STRIDE | Standard Mitigation |
|---------|--------|---------------------|
| Cross-role UI leakage via shared shell/state | Information Disclosure | Strict role route/layout partition and role-scoped store namespaces |
| Sensitive API response over-caching in SW | Information Disclosure | Conservative cache allowlist; avoid broad authenticated response caching in Phase 1 |
| Unsafe bidi/mixed-direction rendering for IDs/URLs | Tampering/Spoofing UX | Use explicit bidi-safe wrappers and RTL-aware text primitives |

## Sources

### Primary (HIGH confidence)
- npm registry package metadata via `npm view` for versions and modified dates. [VERIFIED: npm registry]
- https://nuxt.com/docs/4.x/getting-started/introduction - Nuxt 4 conventions and SSR architecture.
- https://nuxt.com/docs/4.x/api/composables/use-fetch - SSR-friendly fetch composable usage.
- https://nuxt.com/docs/4.x/api/utils/dollarfetch - official Nuxt fetch guidance.
- https://tailwindcss.com/docs/installation/framework-guides/nuxt - Tailwind 4 + Nuxt integration pattern.
- https://pinia.vuejs.org/ssr/nuxt.html - Pinia Nuxt 3/4 integration and auto-import behavior.
- https://vite-pwa-org.netlify.app/frameworks/nuxt.html - Nuxt PWA module usage and `$pwa` primitives.
- https://developer.mozilla.org/en-US/docs/Web/API/Background_Synchronization_API - compatibility and limited availability notes.
- .planning/phases/01-platform-foundation/01-CONTEXT.md
- .planning/phases/01-platform-foundation/01-UI-SPEC.md
- .planning/PROJECT.md
- .planning/REQUIREMENTS.md
- .planning/ROADMAP.md
- .planning/research/STACK.md
- .planning/research/ARCHITECTURE.md
- .planning/research/PITFALLS.md
- .planning/research/SUMMARY.md

### Secondary (MEDIUM confidence)
- None.

### Tertiary (LOW confidence)
- None.

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH - direct npm verification plus official docs.
- Architecture: HIGH - aligned between roadmap, context, and architecture research docs.
- Pitfalls: HIGH - repeated across project pitfall analysis and official browser compatibility docs.

**Research date:** 2026-04-22
**Valid until:** 2026-05-22 (30 days)
