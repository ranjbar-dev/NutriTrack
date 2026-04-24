# Phase 1: Platform Foundation - Context

**Gathered:** 2026-04-22
**Status:** Ready for planning

<domain>
## Phase Boundary

Establish the Persian RTL mobile app shell, design primitives, and installable PWA baseline for client, nutritionist, and super admin surfaces. This phase defines the platform foundation only and does not implement feature-heavy business workflows.

</domain>

<decisions>
## Implementation Decisions

### RTL visual system
- **D-01:** Use a clinical-minimal visual tone for the platform shell: calm, readable, low-chroma UI with strong hierarchy.
- **D-02:** Use Vazirmatn as the base Persian-first typeface for the Phase 1 design system.
- **D-03:** Default to Persian digits and Jalali display in core UI surfaces.

### the agent's Discretion
- Mobile shell navigation pattern (bottom nav, tabs, or mixed) as long as role boundaries remain explicit.
- PWA runtime cache strategy details, as long as sensitive authenticated data is not broadly cached.
- Install and update prompt microcopy/timing details, as long as prompts are clear and non-disruptive.
- Exact role-layout composition per area (client, nutritionist, admin) while preserving strict route isolation.

</decisions>

<specifics>
## Specific Ideas

- Prioritize legibility and clarity over decorative visuals.
- Preserve a distinctly Persian-native feel through typography and numeral/date presentation.

</specifics>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Product and API contracts
- `docs/PRD.md` — Product behavior, role boundaries, Persian RTL mobile focus, and offline constraints.
- `docs/API.md` — Backend integration contract, auth/token behavior, and endpoint envelopes.

### Project scope and requirements
- `.planning/PROJECT.md` — Current scope, constraints, and core value.
- `.planning/REQUIREMENTS.md` — Phase-mapped requirements for platform foundation and downstream phases.
- `.planning/ROADMAP.md` — Phase sequencing and success criteria.

### Research findings for this phase
- `.planning/research/STACK.md` — Stack choices and exclusions for Nuxt 4, Tailwind 4, Pinia, and PWA support.
- `.planning/research/ARCHITECTURE.md` — Recommended boundaries for role areas and foundational layering.
- `.planning/research/PITFALLS.md` — Risks around PWA caching, role leakage, and RTL/Jalali quality.
- `.planning/research/SUMMARY.md` — Consolidated sequencing rationale used to structure this phase.

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- No reusable app code exists yet in this workspace; foundation must be created from scratch.

### Established Patterns
- Planning artifacts define dependency-first sequencing and frontend-only scope constraints.

### Integration Points
- Phase 1 output will be the integration base for Phase 2 auth/session and all role-specific route areas.

</code_context>

<deferred>
## Deferred Ideas

- Rich sync center UX and deeper offline diagnostics are deferred to later phases.
- Advanced visual storytelling, analytics depth, and realtime communication are outside Phase 1 scope.

</deferred>

---

*Phase: 01-platform-foundation*
*Context gathered: 2026-04-22*