---
phase: 01-platform-foundation
plan: 03
subsystem: ui
tags: [pwa, banners, connectivity, cache-boundary]
requires:
  - phase: 01-01
    provides: pwa store baseline and conservative runtime config
  - phase: 01-02
    provides: shared shell primitives and persian ui baseline
  - phase: 01-04
    provides: role-isolated layouts
provides:
  - Reusable install, update, and connectivity banners
  - Role layout wiring for non-blocking PWA prompts
  - Cache-boundary regression test guards
affects: [authentication, client-offline, messaging]
tech-stack:
  added: []
  patterns: [store-driven prompt state, explicit cache-boundary marker tests]
key-files:
  created: [app/components/platform/InstallPromptBanner.vue, app/components/platform/UpdateAvailableBanner.vue, app/components/platform/ConnectivityBanner.vue, tests/platform/cache-boundary.spec.ts]
  modified: [app/stores/platform-pwa.ts, app/layouts/auth.vue, app/layouts/client.vue, app/layouts/nutritionist.vue, app/layouts/admin.vue, tests/platform/pwa-update-prompt.spec.ts, tests/platform/shell-role-isolation.spec.ts, nuxt.config.ts]
key-decisions:
  - "Install prompt remains client-only and appears after intentional role-shell-ready moment."
  - "Update and connectivity banners are shared across role layouts through common store state."
patterns-established:
  - "PWA banner visibility is derived from typed store helpers rather than per-layout local state."
  - "Cache policy regressions are protected by explicit test markers in nuxt config."
requirements-completed: [PLAT-02]
duration: 42min
completed: 2026-04-22
---

# Phase 1 Plan 03: PWA UX and Cache Guard Summary

**Store-driven install or update and connectivity banner orchestration with automated cache-boundary regression protection**

## Performance
- **Duration:** 42 min
- **Started:** 2026-04-22T02:30:00Z
- **Completed:** 2026-04-22T03:12:00Z
- **Tasks:** 3
- **Files modified:** 12

## Accomplishments
- Added reusable install/update/connectivity banners with Persian copy and non-blocking interaction controls.
- Wired banners across role layouts with install prompt restricted to client intentional timing.
- Added cache-boundary regression tests and explicit config marker to prevent authenticated API over-caching drift.

## Task Commits
1. **Task 1 RED** - `d5db892` (test)
2. **Task 1 GREEN** - `cd72aa7` (feat)
3. **Task 2 RED** - `b829ddb` (test)
4. **Task 2 GREEN** - `a072766` (feat)
5. **Task 3 RED** - `c9fb54c` (test)
6. **Task 3 GREEN** - `23cbdc4` (feat)

## Files Created or Modified
- `app/components/platform/InstallPromptBanner.vue` - install banner primitive.
- `app/components/platform/UpdateAvailableBanner.vue` - update banner with explicit refresh action.
- `app/components/platform/ConnectivityBanner.vue` - connectivity status banner.
- `app/stores/platform-pwa.ts` - banner state derivation and refresh action helper.
- `tests/platform/cache-boundary.spec.ts` - regression guards for safe runtime cache boundaries.
- `nuxt.config.ts` - explicit cache-boundary policy marker for test enforceability.

## Decisions Made
- Keep install prompt out of auth, nutritionist, and admin layouts to avoid intrusive first-contact prompts in operational surfaces.
- Require visible cache-boundary policy marker so test guards remain robust when runtime config evolves.

## Deviations from Plan
- Executed after Plan 01-04 instead of before it to satisfy declared dependency (`01-03 depends_on 01-04`).

## Verification Results
- `npm run lint` -> PASS
- `npm run typecheck` -> PASS
- `npm run test:unit -- tests/platform/pwa-update-prompt.spec.ts` -> PASS (5/5)
- `npm run test:unit -- tests/platform/cache-boundary.spec.ts` -> PASS (2/2)

## Self-Check: PASSED
- Summary file exists.
- Commits `d5db892`, `cd72aa7`, `b829ddb`, `a072766`, `c9fb54c`, and `23cbdc4` exist in git history.
