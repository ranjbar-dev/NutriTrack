---
phase: 01-foundation-infrastructure
plan: 02
subsystem: ui
tags: [nuxt4, tailwind-v4, rtl, persian, vazirmatn, jalaali-js, pinia, vue3]

# Dependency graph
requires:
  - phase: 01-01
    provides: Monorepo scaffold with backend/ and frontend/ directories
provides:
  - Nuxt 4 project with app/ directory structure and build pipeline
  - Tailwind CSS v4 with logical RTL properties
  - Vazirmatn variable font as default sans-serif
  - toPersianDigits() and toLatinDigits() utility functions
  - useShamsiDate() composable for Gregorian-to-Shamsi conversion
  - 4 role-based layouts (admin, nutritionist, client, auth)
  - BottomNav component with role-specific navigation items
  - AppButton, AppInput, LoadingSpinner base UI components
  - Placeholder pages for admin, nutritionist, and client dashboards
  - UserRole constants and ROLE_DEFAULT_ROUTES mapping
affects: [01-03-auth-frontend, 01-05-deployment, 02-all-plans, 03-all-plans]

# Tech tracking
tech-stack:
  added: [nuxt@4.4.2, tailwindcss@4.x, "@tailwindcss/postcss@4.x", pinia@3.x, "@pinia/nuxt@0.11.3", vazirmatn@33.0.3, jalaali-js@1.2.7, vitest@3.x, "@vue/test-utils@2.x"]
  patterns: [nuxt4-app-directory, tailwind-v4-logical-properties, composable-pattern, layout-based-role-routing]

key-files:
  created:
    - frontend/nuxt.config.ts
    - frontend/postcss.config.mjs
    - frontend/app/assets/css/main.css
    - frontend/app/app.vue
    - frontend/app/utils/persian-digits.ts
    - frontend/app/utils/constants.ts
    - frontend/app/composables/useShamsiDate.ts
    - frontend/app/components/ui/BottomNav.vue
    - frontend/app/components/ui/AppButton.vue
    - frontend/app/components/ui/AppInput.vue
    - frontend/app/components/ui/LoadingSpinner.vue
    - frontend/app/layouts/admin.vue
    - frontend/app/layouts/nutritionist.vue
    - frontend/app/layouts/client.vue
    - frontend/app/layouts/auth.vue
    - frontend/app/pages/admin/index.vue
    - frontend/app/pages/nutritionist/index.vue
    - frontend/app/pages/nutritionist/clients.vue
    - frontend/app/pages/client/index.vue
    - frontend/app/pages/client/plan.vue
  modified: []

key-decisions:
  - "Used Vazirmatn-Variable-font-face.css for optimal variable font loading"
  - "Used emoji icons for BottomNav (placeholder until proper icon library added)"
  - "Vazirmatn version 33.0.3 (latest on npm) instead of plan's assumed 35.0.1"
  - "jalaali-js imported via default CJS import since it has no ESM or TS types"

patterns-established:
  - "RTL-only logical properties: ms-/me-/ps-/pe-/text-start/text-end, never ml-/mr-/pl-/pr-/text-left/text-right"
  - "Layout-per-role: definePageMeta({ layout: 'rolename' }) on every page"
  - "Persian digit display: wrap all numeric output with toPersianDigits()"
  - "Shamsi date display: use useShamsiDate().formatShamsi() for all date rendering"
  - "Mobile-first: body max-width 430px, no desktop breakpoints"
  - "Auto-import: Nuxt 4 auto-imports utils from app/utils/ and composables from app/composables/"

requirements-completed: [UI-01, UI-02, UI-03, UI-04, UI-05]

# Metrics
duration: 13min
completed: 2026-04-19
---

# Phase 1 Plan 2: Nuxt 4 Frontend Foundation & Persian RTL Summary

**Nuxt 4 app with Tailwind v4 RTL logical properties, Vazirmatn font, Shamsi date composable, Persian digit utilities, and 4 role-based layouts with bottom navigation**

## Performance

- **Duration:** 13 min
- **Started:** 2026-04-19T15:32:52Z
- **Completed:** 2026-04-19T15:46:01Z
- **Tasks:** 3/3
- **Files modified:** 20

## Accomplishments
- Nuxt 4.4.2 project with app/ directory structure builds and runs successfully
- Tailwind CSS v4 with @tailwindcss/postcss processes styles with logical RTL properties
- Vazirmatn variable font loaded as default sans-serif across the entire app
- HTML root element configured with dir="rtl" and lang="fa" for Persian RTL
- toPersianDigits/toLatinDigits verified: "123" → "۱۲۳", "۱۲۳۴۵۶" → "123456"
- useShamsiDate composable converts Gregorian to Shamsi with Persian digit output
- 4 role-based layouts with bottom navigation (admin:3 tabs, nutritionist:5 tabs, client:4 tabs, auth:no nav)
- Base UI components: AppButton (3 variants + loading), AppInput (RTL + error), LoadingSpinner, BottomNav
- Mobile viewport constrained to 430px max width

## Task Commits

Each task was committed atomically:

1. **Task 1: Initialize Nuxt 4 project with all dependencies and configuration** - `2ce097c` (feat)
2. **Task 2: Persian utility functions — digits and Shamsi dates** - `ec7e705` (feat)
3. **Task 3: Role-based layouts, bottom navigation, and placeholder pages** - `44f7700` (feat)

## Files Created/Modified
- `frontend/package.json` - Nuxt 4 project with all dependencies
- `frontend/nuxt.config.ts` - RTL/fa config, Vazirmatn font, Pinia module, runtime config
- `frontend/postcss.config.mjs` - Tailwind CSS v4 PostCSS plugin
- `frontend/app/assets/css/main.css` - Tailwind import, Vazirmatn theme, mobile constraint
- `frontend/app/app.vue` - Root component with NuxtLayout + NuxtPage
- `frontend/app/utils/persian-digits.ts` - toPersianDigits() and toLatinDigits()
- `frontend/app/utils/constants.ts` - UserRole enum, ROLE_DEFAULT_ROUTES, ROLE_LAYOUT_MAP
- `frontend/app/composables/useShamsiDate.ts` - Shamsi date conversion wrapping jalaali-js
- `frontend/app/components/ui/BottomNav.vue` - Fixed bottom nav with active state
- `frontend/app/components/ui/AppButton.vue` - Button with primary/secondary/danger variants
- `frontend/app/components/ui/AppInput.vue` - Input with RTL support, label, error display
- `frontend/app/components/ui/LoadingSpinner.vue` - Animated spinner with size variants
- `frontend/app/layouts/admin.vue` - Admin layout with 3-tab bottom nav
- `frontend/app/layouts/nutritionist.vue` - Nutritionist layout with 5-tab bottom nav
- `frontend/app/layouts/client.vue` - Client layout with 4-tab bottom nav
- `frontend/app/layouts/auth.vue` - Auth layout with centered content, no nav
- `frontend/app/pages/admin/index.vue` - Admin dashboard placeholder
- `frontend/app/pages/nutritionist/index.vue` - Redirects to /nutritionist/clients
- `frontend/app/pages/nutritionist/clients.vue` - Clients list placeholder
- `frontend/app/pages/client/index.vue` - Redirects to /client/plan
- `frontend/app/pages/client/plan.vue` - Diet plan placeholder

## Decisions Made
- **Vazirmatn version:** Used 33.0.3 (latest available on npm) instead of planned 35.0.1 which doesn't exist
- **Variable font:** Used `Vazirmatn-Variable-font-face.css` for optimal variable font loading (verified in node_modules)
- **jalaali-js import:** Used `import jalaali from 'jalaali-js'` default import since the package is CJS-only without ESM/TS exports
- **Emoji icons:** Used emoji placeholders (👥🍽️📊💬👤📋✏️💊) for BottomNav icons - will be replaced with proper icon set in UI phase
- **Manual scaffold:** Created Nuxt 4 project structure manually instead of nuxi init (which hung on interactive prompts in this environment)

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 3 - Blocking] Vazirmatn npm version doesn't exist**
- **Found during:** Task 1 (npm install)
- **Issue:** Plan specified `vazirmatn@^35.0.1` but latest on npm is 33.0.3
- **Fix:** Changed to `vazirmatn@^33.0.3`
- **Files modified:** frontend/package.json
- **Verification:** npm install succeeds, font-face CSS loads, build passes
- **Committed in:** 2ce097c (Task 1 commit)

**2. [Rule 3 - Blocking] @pinia/nuxt incompatible with Nuxt 4.4.2**
- **Found during:** Task 1 (build verification)
- **Issue:** @pinia/nuxt@0.10.0 requires Nuxt ^3.15.0, incompatible with Nuxt 4.4.2
- **Fix:** Upgraded to @pinia/nuxt@0.11.3 which supports Nuxt 4
- **Files modified:** frontend/package.json, frontend/package-lock.json
- **Verification:** Build passes without Pinia compatibility warning
- **Committed in:** 2ce097c (Task 1 commit)

**3. [Rule 3 - Blocking] nuxi init interactive prompt hangs in non-TTY environment**
- **Found during:** Task 1 (project initialization)
- **Issue:** `npx nuxi@latest init frontend` requires interactive template selection, hangs in automated environment
- **Fix:** Manually created Nuxt 4 project structure with correct app/ directory layout and package.json
- **Files modified:** All frontend/ files
- **Verification:** npm install + nuxt prepare + npm run build all pass
- **Committed in:** 2ce097c (Task 1 commit)

---

**Total deviations:** 3 auto-fixed (3 blocking issues)
**Impact on plan:** All auto-fixes necessary to complete task execution. No scope creep. Identical end result to plan specification.

## Issues Encountered
- Tailwind v4 CSS syntax triggers esbuild minifier warnings (`Expected ";" but found "}"`) - cosmetic warning only, does not affect build output or functionality. Known Tailwind v4 + esbuild compatibility issue.
- Nuxt deprecation warnings about `useAppConfig` duplicated imports and `DEP0155` trailing slash patterns - framework-level warnings, no action needed.

## User Setup Required
None - no external service configuration required.

## Next Phase Readiness
- Frontend foundation complete: RTL rendering, Persian typography, layout system all functional
- Ready for Plan 03 (Auth frontend) to add login/OTP pages using the auth layout
- Ready for Plan 05 (Docker deployment) to containerize the frontend
- AppInput component ready for form implementations in auth and data entry pages
- BottomNav and layout system ready for all subsequent page additions

---
*Phase: 01-foundation-infrastructure*
*Completed: 2026-04-19*

## Self-Check: PASSED

All 20 created files verified present. All 3 task commits (2ce097c, ec7e705, 44f7700) verified in git log.
