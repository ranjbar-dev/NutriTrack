# Phase 4: Messaging, Notifications & Lab Exchange - Context

**Gathered:** 2026-04-23
**Status:** Ready for planning

<domain>
## Phase Boundary

Deliver client/nutritionist communication surfaces and notification controls for v1 using the documented APIs. Scope includes conversation history, unread state, text+attachment sending, lab result upload/view/download, and notification preferences.

Out of scope for this phase: realtime transport migration, offline durability for uploads, backend API changes, and expanding role permissions beyond documented contracts.
</domain>

<decisions>
## Implementation Decisions

### Messaging UX
- **D-01:** Messaging remains polling-based (not websocket/realtime) with stable refresh loops and deterministic unread badge updates (MSG-01).
- **D-02:** Message ordering in UI follows API contract chronology and preserves attachment metadata presentation without client-side re-sorting heuristics.
- **D-03:** Compose UI supports text-only, attachment-only, and mixed payloads while enforcing attachment size/type constraints from API/PRD.

### Attachment and Lab Handling
- **D-04:** Message attachments and lab results use explicit upload-state surfaces (idle/uploading/success/failure) with Persian recovery copy.
- **D-05:** Lab result upload is online-only; no offline queueing for file/link uploads (LAB-01 + PRD offline matrix).
- **D-06:** Lab result access flows split by resource type: file download action for file-backed results and external link action for link-backed results.

### Notifications
- **D-07:** Notification preferences UI is role-aware but contract-shared, backed by GET/PATCH preferences endpoints (NOTF-02).
- **D-08:** Push subscription controls are explicit opt-in/opt-out actions with clear permission-state feedback and failure reasons (NOTF-01).

### Role and Safety Boundaries
- **D-09:** Client message surface is restricted to assigned nutritionist conversation; nutritionist surface is restricted to own clients as enforced by existing role guards.
- **D-10:** Existing security posture from Phase 2 applies: auth interceptor and logout semantics must remain intact for all new endpoints.

### Agent's Discretion
- Polling interval tuning for conversation refresh (within performance and UX expectations).
- Visual treatment details for badges/chips/states as long as Persian RTL + token system contracts are respected.
</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Product and Scope
- docs/PRD.md - Messaging, lab upload, and notification behavior contracts including offline matrix and role boundaries.
- .planning/REQUIREMENTS.md - MSG-01, MSG-02, NOTF-01, NOTF-02, LAB-01 requirement definitions.
- .planning/ROADMAP.md - Phase 4 goal, dependencies, and success criteria.

### API Contracts
- docs/API.md - Messaging endpoints (`/messages`, `/clients/:id/messages`, unread count), lab endpoints, and notification preferences endpoints.

### Existing Platform Constraints
- .planning/phases/02-authentication-access-control/02-CONTEXT.md - Auth/session constraints and role access decisions.
- .planning/phases/03-client-offline-daily-loop/03-CONTEXT.md - Offline boundary decisions relevant to messaging/lab behavior.
</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- Existing role-isolated layouts and middleware from Phases 1-2 can host client/nutritionist communication screens.
- Shared error/empty/notice platform primitives can standardize message/lab/notification states.
- Existing auth-fetch and session-refresh plugins should be reused for new API surfaces.

### Established Patterns
- Persian-only RTL mobile-first composition with shared design tokens and safe-area-aware shells.
- Typed contract + composable pattern for endpoint access established in auth and tracking flows.
- Sync-state and retry UI patterns from Phase 3 can inform upload/send feedback patterns.

### Integration Points
- New messaging pages/components integrate under `/client/**` and `/nutritionist/**` namespaces.
- Notification preferences likely integrate into role settings/account surfaces.
- Lab workflows bridge client uploads with nutritionist visibility using shared endpoint contracts.
</code_context>

<specifics>
## Specific Ideas

- Keep conversation screen action-first on mobile: unread indicator, quick compose, attachment CTA, and resilient retry affordance.
- Use concise Persian labels for push permission states to reduce ambiguity during browser permission prompts.
</specifics>

<deferred>
## Deferred Ideas

- Realtime chat transport upgrade (websocket/SSE) - deferred to v2 scope (ADV-01).
- Advanced notification analytics or delivery diagnostics - outside v1 success criteria.
</deferred>

---

*Phase: 04-messaging-notifications-lab-exchange*
*Context gathered: 2026-04-23*