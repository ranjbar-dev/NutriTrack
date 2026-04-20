# Phase 5: Communication & Collaboration - Context

**Gathered:** 2026-04-20
**Status:** Ready for planning
**Mode:** --auto (sub-agent, no interactive session)

<domain>
## Phase Boundary

Phase 5 adds the collaboration layer on top of the completed diet-plan and tracking foundation: direct messaging between each client and their assigned nutritionist, client food-addition requests, and the nutritionist client-management workspace. It reuses Phase 4 history and lab-result flows inside the client profile, but does not add offline sync, push notifications, or any desktop-specific UI.

</domain>

<decisions>
## Implementation Decisions

### Messaging architecture

- **D-01:** Messaging is strictly one-to-one between a client and their assigned nutritionist; there are no group conversations, broadcast threads, or cross-nutritionist conversations.
- **D-02:** Message history is chronological within a conversation and immutable after send — no edit/delete support. The backend stores `sent_at`, optional `read_at`, and attachment metadata alongside the text body.
- **D-03:** Chat delivery uses 10-second polling only while the conversation screen is open. A lightweight unread-count endpoint powers badges in navigation and conversation lists.
- **D-04:** Message attachments follow the same security posture established in Phase 4 lab uploads: JPG/PNG up to 5 MB, PDF up to 10 MB, magic-byte validation, UUID filenames, and authenticated downloads with `Content-Disposition: attachment`.

### Food-request workflow

- **D-05:** Food requests are their own workflow, not a special message type. Each request belongs to one client, one assigned nutritionist, and a lifecycle of `pending`, `approved`, or `rejected`.
- **D-06:** Approving a request does not auto-create an incomplete food row. Approval routes the nutritionist into the existing food-creation flow with the request name/description prefilled, and the request is only marked approved once the food record is successfully saved.
- **D-07:** Client-facing notification for food-request outcomes is in-app: the request list shows live status/reason updates, and updated counts can be refreshed on the same polling cadence used for communication screens.

### Nutritionist client workspace

- **D-08:** The nutritionist client list remains the primary workspace and stays mobile-first card based. It must support search by name/mobile, active/inactive filtering, and sort by name or last activity.
- **D-09:** The client profile stays under the existing nutritionist client routes and reuses Phase 3 plan summaries plus Phase 4 tracking/history components. Quick actions stay pinned near the top: create plan, send message, activate/deactivate client, and edit profile fields.
- **D-10:** Height and date of birth are editable only by nutritionists. Activating/deactivating a client is a status toggle that preserves all plans, tracking history, messages, and uploads.

### File handling and security

- **D-11:** Phase 5 file storage remains local filesystem only, using a config-driven upload root with per-client subdirectories for message attachments. Direct filesystem paths are never exposed to the frontend.
- **D-12:** Per-client storage limits are enforced in the service layer before persisting a new attachment. Failed validations must return Persian errors and leave no partial files on disk.

### The agent's Discretion

- Exact chat bubble styling, spacing, and avatar treatment
- Whether the nutritionist conversation list shows the last message preview or only unread/status indicators
- Whether food-request status updates appear as inline badges, cards, or grouped timeline items
- Whether client-profile quick actions use sticky buttons or a compact action strip

</decisions>

<specifics>
## Specific Ideas

- Reuse Phase 4 file-validation and authenticated-download patterns instead of inventing a second upload/security model for messages.
- Keep the nutritionist client profile as the single place where archived plans, tracking history, lab results, and new communication actions converge.
- The food-request approval path should feel like a shortcut into the existing food-management UI, not a second parallel CRUD flow.

</specifics>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Product scope
- `.planning/ROADMAP.md` §Phase 5 — phase goal, dependency position, and success criteria
- `.planning/REQUIREMENTS.md` §Messaging, §Food Requests, §Client Management, §Security — MSG-01 through MSG-07, FREQ-01 through FREQ-04, CLNT-02 through CLNT-07, SEC-04, SEC-05, SEC-08
- `docs/phases.md` §Phase 5 — implementation guidance, API shape, and validation checklist

### Existing architecture and reusable phase contracts
- `.planning/phases/04-client-tracking-suite/04-CONTEXT.md` — file-validation, tracking-history, and client/nutritionist route conventions to reuse
- `.planning/phases/03-diet-plan-engine/03-CONTEXT.md` — active-plan summaries, archived-plan history, and client route patterns used in the profile workspace
- `.planning/phases/02-core-data-domain/02-CONTEXT.md` — food CRUD patterns and nutritionist/shared-data conventions for the food-request approval flow
- `.planning/research/STACK.md` — confirmed stack choices for polling chat, local filesystem storage, and Nuxt/Gin integration

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `frontend/app/composables/useApi.ts` — authenticated fetch wrapper for chat, food-request, and client-management API calls
- `frontend/app/components/tracking/WeightChart.vue` and Phase 4 tracking pages — reusable history tabs inside the nutritionist client profile
- `frontend/app/stores/*` Pinia pattern — established store structure for client/nutritionist data loading and optimistic updates
- `backend/internal/service/tracking_service.go` and `backend/internal/handler/tracking_handler.go` — existing file-validation, Persian error, and authenticated download patterns to mirror for message attachments

### Established Patterns
- Handler → service → repository layering with sqlc-backed queries
- Repository-level row authorization instead of trusting handler filtering
- Persian-only validation and error responses
- Mobile card/list UIs instead of desktop tables
- Local filesystem uploads served only through authenticated application endpoints

### Integration Points
- Client and nutritionist bottom navigation already reserve communication-oriented surface area for unread badges and messaging routes
- The existing nutritionist client route tree should host the new profile, history, and message entry points
- Food-request approval should link directly into the existing nutritionist food create/edit experience rather than duplicating food-entry forms
- Phase 4 tracking/history data and archived Phase 3 plans should surface inside the Phase 5 client profile without changing the underlying Phase 4 APIs

</code_context>

<deferred>
## Deferred Ideas

- Offline message caching, send queues, and reconnect sync — Phase 6
- Push notifications for unread messages and food-request outcomes — Phase 6
- Adaptive polling intervals or background polling when the chat is closed — future enhancement / backlog
- Voice notes, message search, and richer collaboration features — future backlog

</deferred>

---

*Phase: 05-communication-collaboration*
*Context gathered: 2026-04-20*
