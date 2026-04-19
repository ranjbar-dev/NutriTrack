---
phase: 02-core-data-domain
plan: 03
subsystem: frontend
tags: [nuxt4, vue3, pinia, tailwind-v4, rtl, persian]

# Dependency graph
requires:
  - phase: 02-core-data-domain/02-02
    provides: Food CRUD API at /api/foods

provides:
  - Food management UI for nutritionists at /nutritionist/foods/*
  - useFoodStore Pinia store with full CRUD + pagination + search state
  - FoodCard and CategoryPills reusable components

affects:
  - phase-03-diet-plan-engine (food search UI reuse)

# Tech tracking
tech-stack:
  added:
    - vue-tsc@3.2.7 (devDependency for nuxi typecheck)
  patterns:
    - "Pinia store with Persian normalization for search (ی/ي, ک/ك)"
    - "300ms debounce via setTimeout/clearTimeout in watch() handler"
    - "localStorage draft auto-save with deep watcher and draftReady gate"
    - "Load more pagination (append mode vs reset mode in fetchFoods)"
    - "Inline Persian validation errors with labelMap ?? key fallback"

key-files:
  created:
    - frontend/app/stores/food.ts
    - frontend/app/components/food/FoodCard.vue
    - frontend/app/components/food/CategoryPills.vue
    - frontend/app/pages/nutritionist/foods/index.vue
    - frontend/app/pages/nutritionist/foods/new.vue
    - frontend/app/pages/nutritionist/foods/[id].vue
  modified: []

key-decisions:
  - "labelMap[key] ?? key fallback used for strict TypeScript index access (noUncheckedIndexedAccess)"
  - "normalizePersianText() applied to search in store (not just on page) to normalize ی/ي and ک/ك before API call"
  - "draftReady ref gates the localStorage watcher to prevent overwriting draft on initial form mount"

patterns-established:
  - "FoodCard emits edit/delete events with food.id for decoupled parent handling"
  - "CategoryPills emits select with string | null (null = all categories)"
  - "Edit form does not use localStorage draft — only new.vue does"

requirements-completed: [FOOD-01, FOOD-02, FOOD-07, FOOD-08]

# Metrics
duration: 15m
completed: 2026-04-19
---

# Phase 2 Plan 03: Food Management Frontend Summary

**Complete food management UI for nutritionists: card-based list with Persian fuzzy search (300ms debounce), category filter pills, active/inactive toggle, load more pagination (20/page), and multi-section add/edit forms with localStorage draft auto-save and inline Persian validation.**

## Performance

- **Duration:** ~15 min
- **Completed:** 2026-04-19
- **Tasks:** 3 implemented + 1 checkpoint (auto-approved)
- **Files created:** 6

## Accomplishments

- **Task 1:** Created `useFoodStore` Pinia store with full CRUD actions, append-mode pagination, Persian search normalization, and search/category/active filter state. Created `FoodCard` component with Persian category badges and `toPersianDigits()` calorie/protein summary. Created `CategoryPills` horizontal scrollable filter with 8 Persian category labels.
- **Task 2:** Built food list page at `/nutritionist/foods` with 300ms debounced search, category pills integration, active/inactive toggle, empty states (no foods / no results / no category), FoodCard list with edit/delete handlers, and "Load More" button.
- **Task 3:** Built `new.vue` (4-section add form with localStorage auto-save draft) and `[id].vue` (edit form pre-populated from store) with full Persian validation: required name, minimum 1 category, numeric range checks, measurement amount > 0. Both forms disable inputs and show loading state during submit.
- **Task 4 (Checkpoint):** ⚡ Auto-approved — running in automated non-interactive mode.

## Task Commits

| Task | Name | Commit | Files |
|------|------|--------|-------|
| 1 | Food store and shared components | `f4f7650` | food.ts, FoodCard.vue, CategoryPills.vue |
| 2 | Food list page | `e1b0ba5` | foods/index.vue |
| 3 | Food add and edit forms | `80a3f08` | foods/new.vue, foods/[id].vue |
| - | Install vue-tsc | `7b9c95e` | package.json, package-lock.json |

## Decisions Made

1. **labelMap fallback** — Used `labelMap[key] ?? key` instead of `labelMap[key]` to satisfy TypeScript strict index access (avoids `string | undefined` not assignable to `string` error).
2. **draftReady gate** — Added `draftReady` ref that is set to `true` only after onMounted restores draft from localStorage, preventing the watch from immediately overwriting with empty form state.
3. **Persian normalization in store** — `normalizePersianText()` applied inside `setSearch()` in the store so the API always receives normalized text regardless of what the component passes.

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Fixed TypeScript strict index access in validation loop**
- **Found during:** TypeScript typecheck run after Task 3
- **Issue:** `labelMap[key]` in `validateForm()` was typed as `string | undefined` causing TS2345 error in both new.vue and [id].vue
- **Fix:** Changed to `labelMap[key] ?? key` fallback — falls back to the raw field key name if mapping is missing
- **Files modified:** `frontend/app/pages/nutritionist/foods/new.vue`, `frontend/app/pages/nutritionist/foods/[id].vue`
- **Committed in:** `80a3f08`

### Pre-existing Errors (Out of Scope)

The following TypeScript errors existed before this plan and are not caused by plan changes:
- `app/composables/useShamsiDate.ts:1` — Missing `@types/jalaali-js` declaration (Phase 1 issue)
- `app/utils/persian-digits.ts:8` — Strict array index type on `persianDigits[parseInt(d)]` (Phase 1 issue)
- `nuxt.config.ts:22` — Missing `@types/node` for `process.env` (Phase 1 issue)

These are logged to deferred items and are out of scope for this plan.

## Checkpoints

**Task 4 — human-verify:**
⚡ Auto-approved: Complete food management UI for nutritionists — card-based list with Persian fuzzy search, category filter pills, load more pagination, add/edit forms with localStorage draft auto-save and validation.

## Known Stubs

None — all data flows through `useFoodStore` which calls real `/api/foods` endpoints.

## Threat Flags

None — no new trust boundaries introduced beyond what the plan's threat model covers.

---
*Phase: 02-core-data-domain*
*Completed: 2026-04-19*

## Self-Check: PASSED
- FOUND: frontend/app/stores/food.ts
- FOUND: frontend/app/components/food/FoodCard.vue
- FOUND: frontend/app/components/food/CategoryPills.vue
- FOUND: frontend/app/pages/nutritionist/foods/index.vue
- FOUND: frontend/app/pages/nutritionist/foods/new.vue
- FOUND: frontend/app/pages/nutritionist/foods/[id].vue
- FOUND: f4f7650
- FOUND: e1b0ba5
- FOUND: 80a3f08
- FOUND: 7b9c95e
