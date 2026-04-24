# Phase 6: Admin Governance - Context

**Gathered:** 2026-04-24
**Status:** Ready for planning

<domain>
## Phase Boundary

Deliver the super-admin operational surface in the Nuxt client: platform stats, nutritionist account lifecycle management, read-only nutritionist client visibility, and elevated shared catalogue governance for foods, medications, and food categories.

Out of scope for this phase: backend API changes, offline support for admin surfaces, desktop-first data-grid patterns, and audit-log screens not backed by documented API endpoints.
</domain>

<decisions>
## Implementation Decisions

### Role, Routing, and Connectivity Boundaries
- **D-01:** Admin governance remains strictly online-first; Phase 3 offline queueing does not extend into `/admin/**` flows.
- **D-02:** All new super-admin pages stay under `/admin/**` and continue to rely on deny-by-default role guards from Phase 2.
- **D-03:** Admin catalogue operations must use elevated admin endpoints where available instead of reusing nutritionist-scoped ownership assumptions.

### Platform Oversight and Nutritionist Management
- **D-04:** Platform stats are presented as compact mobile KPI surfaces backed by `GET /admin/stats`, not derived from stitched client-side queries.
- **D-04a:** The admin stats UI must only show API-backed metrics; because `GET /admin/stats` exposes `total_clients` but not active/inactive client splits, Phase 6 should show total clients only unless the API contract changes.
- **D-05:** Nutritionist management covers create, list, detail, update, activate/deactivate, and read-only client-list visibility via the documented `/admin/nutritionists*` endpoints.
- **D-06:** Nutritionist status changes are explicit confirmation actions with clear Persian success/error states because they directly affect access control.

### Shared Catalogue Governance
- **D-07:** Super-admin food governance uses `GET /admin/foods` and `DELETE /admin/foods/:id` for elevated search and force-delete, while create/edit can continue through shared food CRUD where the API already grants `super_admin` authority.
- **D-08:** Super-admin medication governance uses `GET /admin/medications` and `DELETE /admin/medications/:id` for elevated search and force-delete, while create/edit continue through shared medication CRUD where permitted to `super_admin`.
- **D-09:** Food-category governance is admin-only for create/delete and should be treated as part of the catalogue control surface because category taxonomy materially affects shared food management.
- **D-10:** PRD mentions audit-log visibility for food and medication changes, but `docs/API.md` does not expose an audit-log contract; Phase 6 should not invent that surface and must defer it unless a documented endpoint appears.

### Agent's Discretion
- Mobile information architecture for balancing KPIs, list management, and destructive actions without collapsing into a desktop admin table.
- Search/filter defaults and detail-sheet density so long as Persian RTL mobile ergonomics remain primary.
</decisions>

<canonical_refs>
## Canonical References

### Product and Scope
- docs/PRD.md - Section 5.14 super-admin panel requirements and Section 6 offline boundary.
- .planning/REQUIREMENTS.md - ADMIN-01 and ADMIN-02 requirement definitions.
- .planning/ROADMAP.md - Phase 6 goal, dependency, and success criteria.

### API Contracts
- docs/API.md - Section 8 (admin stats), Section 9 (admin nutritionist management), Section 11.6-11.7 (admin food governance), Section 12.2-12.3 (admin food categories), and Section 13.6-13.7 (admin medication governance).

### Existing Platform Constraints
- .planning/phases/02-authentication-access-control/02-CONTEXT.md - Role guard and session-expiry constraints for admin-only routing.
- .planning/phases/05-nutritionist-workspace-plan-authoring/05-CONTEXT.md - Shared catalogue interaction patterns and the boundary between nutritionist authoring and admin governance.
</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- Existing protected shells, auth fetch utilities, and role middleware from Phases 1-2 should host admin pages without duplicating session logic.
- Shared list/search/form primitives from nutritionist workspace flows can be reused for admin roster and catalogue screens.
- Catalogue contracts and picker-related domain types established in Phase 5 should be extended rather than forked for admin governance.

### Established Patterns
- Contract-first `app/types/*` plus composable-driven API access.
- Persian-only RTL mobile-first layout with sheets, cards, confirmations, and safe-area-aware action bars.
- Explicit loading, empty, error, and destructive-action confirmation states instead of implicit optimistic admin mutations.

### Integration Points
- Admin screens should live under `/app/pages/admin/**` and `/app/components/admin/**` alongside existing role shell structure.
- Shared food, medication, and category types/composables likely need incremental extension to support elevated admin endpoints without breaking nutritionist flows.
- Nutritionist detail surfaces should link to existing read-only client-list and identity views rather than creating a second client management model.
</code_context>

<specifics>
## Specific Ideas

- Keep the admin home screen action-first on mobile: KPI summary first, then nutritionist roster and catalogue entry points.
- Prefer detail sheets/drawers for nutritionist profile and destructive catalogue actions instead of forcing multi-column layouts.
</specifics>

<deferred>
## Deferred Ideas

- Audit-log timeline UI until the API contract exists.
- Bulk admin actions, CSV export, and advanced analytics beyond the Phase 6 success criteria.
- Desktop-heavy moderation dashboards or multi-pane catalogue tooling.
</deferred>

---

*Phase: 06-admin-governance*
*Context gathered: 2026-04-24*