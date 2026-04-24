# Phase 5 Validation: Nutritionist Workspace & Plan Authoring

**Generated:** 2026-04-23
**Phase:** 05-nutritionist-workspace-plan-authoring

---

## UAT Criteria

### SC-1 - Nutritionist roster and profile workspace (NUTR-01, NUTR-02)

- [ ] Nutritionist opens `/nutritionist/clients` and sees paginated client roster.
- [ ] Nutritionist can search roster by name/mobile and filter active/inactive state.
- [ ] Nutritionist can open `/nutritionist/clients/:id` profile and view identity summary.
- [ ] Profile shows current plan summary plus tabs for tracking history, messages, labs, and archived plans.

### SC-2 - Plan metadata authoring (NUTR-03)

- [ ] Nutritionist can create a client plan period with start/end dates and optional title/notes/water target.
- [ ] Nutritionist can edit existing plan metadata from plan edit surface.
- [ ] Validation prevents invalid date ranges and missing required fields.

### SC-3 - Hierarchical plan structure authoring (NUTR-04)

- [ ] Nutritionist can add/remove days, meals, options, and option items.
- [ ] Nutritionist can add/remove exercise recommendations per day.
- [ ] Nutritionist can add/remove prescriptions per day.
- [ ] Updated plan structure persists through API and reloads correctly.

### SC-4 - Catalogue lookup integration (CAT-01, CAT-02)

- [ ] Food picker supports search + category filter and selecting food items into meal option items.
- [ ] Medication picker supports search and selecting medications into day prescriptions.
- [ ] Authoring UI displays API-provided totals/ranges where available.

### SC-5 - Food request flow (CAT-03)

- [ ] Client can submit a food request from client surface.
- [ ] Nutritionist can list pending food requests.
- [ ] Nutritionist can approve a request with food payload or reject with reason.
- [ ] Request status updates and Persian feedback states are visible after moderation.

---

## Regression Guards

| Guard | Verification |
|---|---|
| Role isolation remains intact | Client cannot access `/nutritionist/**`; nutritionist cannot access client-only moderation actions |
| Auth/session behavior preserved | 401 responses from new endpoints route through existing auth failure handling |
| Offline scope unchanged | No nutritionist authoring endpoints are routed into client offline queue |
| Persian RTL fidelity | All new screens/components render in RTL with Persian copy and mobile spacing |

---

## Suggested Test Files

- `tests/nutritionist/client-roster.spec.ts`
- `tests/nutritionist/client-profile-shell.spec.ts`
- `tests/nutritionist/plan-authoring-metadata.spec.ts`
- `tests/nutritionist/plan-authoring-structure.spec.ts`
- `tests/nutritionist/catalog-food-search.spec.ts`
- `tests/nutritionist/catalog-medication-search.spec.ts`
- `tests/nutritionist/food-request-review.spec.ts`
- `tests/client/food-request-submit.spec.ts`

---

*Validation maps to ROADMAP.md Phase 5 success criteria SC-1 through SC-5*