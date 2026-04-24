# Phase 5: Nutritionist Workspace & Plan Authoring - Context

**Gathered:** 2026-04-23
**Status:** Ready for planning

<domain>
## Phase Boundary

Deliver the nutritionist operational workspace in the Nuxt client: client roster discovery, client profile workspace, mobile-first plan authoring, shared food/medication catalogue lookup for authoring, and food request submit/review actions.

Out of scope for this phase: backend API changes, admin-only catalogue moderation endpoints, realtime transport changes, desktop-first layout work, and offline queue support for nutritionist authoring surfaces.
</domain>

<decisions>
## Implementation Decisions

### Workspace and Role Boundaries
- **D-01:** Nutritionist workspace remains online-first and does not inherit client offline queue complexity from Phase 3.
- **D-02:** All new nutritionist pages remain strictly under `/nutritionist/**` and rely on existing role middleware for deny-by-default route access.
- **D-03:** Client food-request submission UI is exposed under `/client/**`, while review actions are exposed only under `/nutritionist/**`.

### Client Roster and Profile
- **D-04:** Roster UX must provide search by name/mobile plus active/inactive filtering and lightweight sorting using API-backed pagination (`GET /clients`).
- **D-05:** Client profile workspace is tabbed and must surface identity, current plan summary, tracking history, messages, labs, and archived plans in one coherent mobile flow.

### Plan Authoring
- **D-06:** Plan authoring follows API hierarchy exactly: plan -> days -> meals -> options -> items, with exercise and prescriptions managed per day.
- **D-07:** Plan metadata editing (period, title, notes, water target) and structural editing (days/meals/options/items) are separate UI concerns but share one authoring state model.
- **D-08:** Computed nutrition totals shown to nutritionists use API-provided totals/ranges from plan payloads; no hand-rolled nutrition math engine in the client.

### Catalogue and Requests
- **D-09:** Food and medication catalogue interactions are search/view first for Phase 5 authoring flows; elevated admin catalogue governance stays in Phase 6.
- **D-10:** Food request moderation actions (approve/reject) must show explicit irreversible-action confirmation and Persian feedback states.

### Agent's Discretion
- Precise card/sheet visual hierarchy for mobile authoring ergonomics.
- Polling cadence for profile-adjacent data (messages/tracking) within performance limits.
</decisions>

<canonical_refs>
## Canonical References

### Product and Scope
- docs/PRD.md - Sections 5.1, 5.2, 5.3, 5.12, 5.13 and non-goal/offline boundaries.
- .planning/REQUIREMENTS.md - NUTR-01..NUTR-04 and CAT-01..CAT-03 requirement definitions.
- .planning/ROADMAP.md - Phase 5 goal, dependencies, and success criteria.

### API Contracts
- docs/API.md - Sections 10 (Client Management), 11 (Foods), 12 (Food Categories), 13 (Medications), 14 (Diet Plans), and 18 (Food Requests).

### Existing Platform Constraints
- .planning/phases/02-authentication-access-control/02-CONTEXT.md - Role and auth guard constraints.
- .planning/phases/03-client-offline-daily-loop/03-CONTEXT.md - Client-only offline support boundary.
- .planning/phases/04-messaging-notifications-lab-exchange/04-CONTEXT.md - Messaging/lab integration patterns used in profile workspace tabs.
</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- Existing nutritionist shell and role middleware already establish protected namespace routing.
- Messaging/lab composables created in Phase 4 can be reused in profile-level tabs instead of duplicate transport logic.
- Shared mobile primitives (`EmptyState`, `InlineNotice`, sheets/cards) should be reused for roster and authoring screens.

### Established Patterns
- Typed `app/types/*` + `app/composables/*` contract-first pattern.
- Persian RTL mobile-first pages with safe-area-aware layouts.
- API-centric composables with `$fetch` for mutations and `useAsyncData` for list/read flows.

### Integration Points
- Nutritionist workspace screens live under `/app/pages/nutritionist/**` and `/app/components/nutritionist/**`.
- Client-side request submission under `/app/pages/client/**` must align with existing client shell navigation conventions.
</code_context>

<specifics>
## Specific Ideas

- Keep roster list rows compact but action-rich (status, active plan indicator, and quick open).
- Keep plan authoring progressive: metadata first, then structure editing by day/meal/option.
</specifics>

<deferred>
## Deferred Ideas

- Reusable plan templates and cloning workflows (v2 / EXP-03).
- Full nutritionist offline drafting mode and local conflict resolution.
- Advanced analytics overlays for authoring outcomes.
</deferred>

---

*Phase: 05-nutritionist-workspace-plan-authoring*
*Context gathered: 2026-04-23*