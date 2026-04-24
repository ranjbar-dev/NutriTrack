# Phase 5 Research: Nutritionist Workspace & Plan Authoring

**Generated:** 2026-04-23
**Phase:** 05-nutritionist-workspace-plan-authoring
**Sources:** docs/API.md sections 10, 11, 12, 13, 14, 18 · docs/PRD.md sections 5.1, 5.2, 5.3, 5.12, 5.13 · existing app patterns

---

## 1. Client Management Surface (NUTR-01, NUTR-02)

### Core endpoints
- `GET /clients?page=&page_size=` list roster
- `GET /clients/:id` profile summary
- `PATCH /clients/:id` profile update
- `PATCH /clients/:id/status` active/inactive toggle

### UX implications
- Search/filter/sort controls should drive API query state for pagination-safe behavior.
- Profile screen should aggregate linked data tabs using existing APIs already implemented in previous phases: tracking, messages, lab results, and plans.
- Keep profile shell mobile-first and tabbed instead of dense desktop table views.

## 2. Diet Plan Authoring Surface (NUTR-03, NUTR-04)

### Authoring hierarchy from API section 14
- Plan metadata: `POST /clients/:id/plans`, `PATCH /plans/:id`
- Structural nodes:
  - Day: add/delete
  - Meal: add/delete
  - Option: add/delete
  - Item: add/delete
  - Exercise: add/delete
  - Prescription: add/delete

### Modeling guidance
- Maintain explicit IDs at each node boundary to avoid accidental cross-node mutation.
- Authoring UI should represent each hierarchy level distinctly (day cards, meal panels, option rows).
- Prefer API-provided total ranges (`total_range`, `totals`) for display; no custom nutrient computation engine.

## 3. Shared Catalogue Lookup (CAT-01, CAT-02)

### Food contracts
- `GET /foods?q=&category_id=&page=&page_size=`
- `GET /food-categories`

### Medication contracts
- `GET /medications?q=&page=&page_size=`

### UX guidance
- Use picker sheets/dialogs for selecting foods and medications from authoring context.
- Keep lookup read-oriented for this phase; admin-level catalogue lifecycle is phase 6 scope.

## 4. Food Request Workflow (CAT-03)

### Client action
- `POST /food-requests` with requested food name.

### Nutritionist action
- `GET /food-requests` pending review list.
- `POST /food-requests/:id/approve` to create shared food.
- `POST /food-requests/:id/reject` with reason.

### UX guidance
- Request moderation should expose clear pending/approved/rejected status tags.
- Approve/reject actions should use explicit confirmation and failure handling copy in Persian.

## 5. Existing Codebase Patterns to Reuse

| Concern | Existing pattern | Reuse target |
|---|---|---|
| Typed API contracts | `app/types/*.ts` | New nutritionist/catalog/request types |
| Composable transport | `app/composables/useMessagingApi.ts`, `useLabApi.ts` | New workspace and authoring composables |
| Role shell | `app/layouts/nutritionist.vue` + middleware | New nutritionist routes |
| Mobile sheet/cards | Existing client/nutritionist components | Authoring editor and pickers |
| Test style | `tests/client/*.spec.ts`, `tests/platform/*.spec.ts` | New `tests/nutritionist/*.spec.ts` set |

## 6. Constraints and Pitfalls

- Keep frontend-only scope; no assumptions of backend behavior outside documented contracts.
- Nutritionist workspace is online-first by product rule; avoid introducing offline queue logic to authoring actions.
- Avoid route leakage across roles; verify all new paths stay under proper shells.
- API list endpoints are paginated; plan UI state around `meta` pagination envelope.

---

*Research generated for: 05-nutritionist-workspace-plan-authoring*