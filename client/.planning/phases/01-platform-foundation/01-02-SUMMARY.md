---
phase: 01-platform-foundation
plan: 02
subsystem: ui
tags: [rtl, persian, jalali, design-tokens, vue]
requires:
  - phase: 01-01
    provides: platform baseline and tests
provides:
  - Persian tokenized style system with safe-area primitives
  - Locale helpers for Persian digits and Jalali date display
  - Reusable platform shell components for role-aware mobile layouts
affects: [platform-foundation, auth, client-offline]
tech-stack:
  added: []
  patterns: [token-first styling, composable locale helpers, role-aware shell primitives]
key-files:
  created: [app/assets/css/main.css, app/lib/design/tokens.css, app/composables/useRtl.ts, app/composables/usePersianFormat.ts, app/lib/locale/numerals.ts, app/lib/locale/jalali.ts, app/components/platform/AppHeader.vue, app/components/platform/AppShell.vue, app/components/platform/BottomNavClient.vue, app/components/platform/InlineNotice.vue, app/components/platform/EmptyState.vue, app/components/platform/ErrorState.vue]
  modified: [tests/platform/persian-locale-baseline.spec.ts]
key-decisions:
  - "Default to Persian digit rendering while keeping OTP and identifiers Latin-safe by design."
  - "Use a restrained clinical-minimal token set to stabilize visual consistency in later feature phases."
patterns-established:
  - "Locale formatting is centralized in app/composables/usePersianFormat.ts and locale utilities."
  - "Shared platform UI state components are role-aware but do not contain role business logic."
requirements-completed: [PLAT-03]
duration: 48min
completed: 2026-04-22
---

# Phase 1 Plan 02: Persian RTL UI Foundation Summary

**Persian RTL token system and Jalali or numeral formatting helpers with reusable mobile shell primitives for role surfaces**

## Performance
- **Duration:** 48 min
- **Started:** 2026-04-22T00:56:00Z
- **Completed:** 2026-04-22T01:44:00Z
- **Tasks:** 3
- **Files modified:** 13

## Accomplishments
- Implemented tokenized color, spacing, typography, and safe-area foundations for the Persian mobile shell.
- Added locale utilities and composables for Persian digit defaults and Jalali date formatting.
- Built reusable shell/header/nav and empty/error/notice primitives for downstream role pages.

## Task Commits
1. **Task 1 RED** - `a00d21e` (test)
2. **Task 1 GREEN** - `30086a9` (feat)
3. **Task 2 RED** - `cd72783` (test)
4. **Task 2 GREEN** - `02c8de2` (feat)
5. **Task 3** - `69eae3f` (feat)

## Files Created or Modified
- `app/lib/design/tokens.css` - clinical-minimal palette, spacing, typography, and safe-area tokens.
- `app/assets/css/main.css` - global Persian RTL baseline styles and motion-safe defaults.
- `app/lib/locale/numerals.ts` - Persian and Latin numeral conversion helpers.
- `app/lib/locale/jalali.ts` - Jalali display formatting helper.
- `app/composables/usePersianFormat.ts` - centralized composable formatter API.
- `app/components/platform/AppShell.vue` - role-aware reusable shell wrapper.

## Decisions Made
- Keep role-agnostic shell primitives in shared components and keep role specifics in route or layout-level wiring.
- Keep safe-area and typography values in CSS tokens to avoid hardcoded per-component spacing.

## Deviations from Plan
None - plan executed as written.

## Verification Results
- `npm run lint` -> PASS
- `npm run typecheck` -> PASS
- `npm run test:unit -- tests/platform/persian-locale-baseline.spec.ts` -> PASS (5/5)
- `npm run test:unit -- tests/platform/shell-role-isolation.spec.ts` -> PASS (3/3)

## Known Stubs
- `app/components/platform/BottomNavClient.vue` uses baseline static links; role-specific destination rules remain enforced later by middleware and route guards.

## Self-Check: PASSED
- Summary file exists.
- Commits `a00d21e`, `30086a9`, `cd72783`, `02c8de2`, and `69eae3f` exist in git history.
