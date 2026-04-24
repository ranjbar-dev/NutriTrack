# Project Retrospective

*A living document updated after each milestone. Lessons feed forward into future planning.*

## Milestone: v1.0 — NutriTrack Client MVP

**Shipped:** 2026-04-24  
**Phases:** 6 | **Plans:** 27 | **Timeline:** 5 days (2026-04-19 → 2026-04-24)

### What Was Built

- Persian RTL mobile PWA shell: installable, role-aware navigation, design tokens, Jalali date support
- Role-based auth: client OTP, nutritionist/admin email+password, token refresh, route guards
- Client offline daily loop: Today view, plan reading, multi-category tracking (food, water, sleep, exercise, meds, body), sync queue with durable local IDs and retry
- Messaging and lab exchange: conversation screens for both roles, file attachments, push notification controls, lab result upload/access
- Nutritionist workspace: client roster, profile, full plan authoring hierarchy, catalogue picker, food request moderation
- Super-admin governance: dashboard stats, nutritionist CRUD, catalogue approve/block/categorize, read-only client visibility

### What Worked

- **Dependency-first sequencing** — ordering phases platform → auth → offline → messaging → nutritionist → admin meant zero upstream rework. Each phase built cleanly on stable prior foundations.
- **Contract-first API typing** — establishing typed API contracts (and later admin contracts) as the first plan of each phase caught shape mismatches early, before UI code was written.
- **Bounded offline scope** — keeping offline durability strictly to client-side flows simplified Phase 4+ significantly. Admin and nutritionist surfaces are online-first by design and never needed offline complexity.
- **Plan-checker agent** — catching planning defects (unresolved research questions, broken validation sequences) before execution saved rework mid-wave.
- **Wave-based execution** — splitting execution into waves (foundation → UI → detail) within each phase made work parallelizable and reduced context sprawl.
- **Read-only constraint at type level** — enforcing `readOnly: true` in admin types caught accidental mutations at compile time rather than runtime.

### What Was Inefficient

- **REQUIREMENTS.md not updated during execution** — Phase 3 requirements were not checked off as plans completed. Discovered and corrected during milestone close. Future milestones should tick off requirements per plan as execution completes.
- **Path structure assumption** — commit subagent initially tried `src/` paths instead of the correct `app/` directory structure. A quick codebase scan before the first commit would have prevented this.
- **Textual bugs caught late** — "audit" and "edit" appearing in wrong UX contexts were caught during code review rather than during planning or spec review. Forbidden terminology should be embedded in phase CONTEXT.md files.

### Patterns Established

- **CONTEXT.md → RESEARCH.md → PLAN.md → execution** — the four-artifact phase lifecycle worked consistently across all 6 phases.
- **Deny-by-default role guards via global middleware** — fail-safe role isolation pattern that should persist in v1.1 and beyond.
- **Auth error forcing logout** — INVALID_TOKEN, TOKEN_REVOKED, and UNAUTHORIZED all trigger session termination with a session-expired marker. No partial session states.
- **Search context preservation on delete** — preserving and reapplying filters across destructive catalogue actions prevents disorientation. Apply to any future list+action surfaces.
- **PATCH over PUT for updates** — field-level PATCH merge to avoid blanking unmodified fields. Default pattern for all future update operations.

### Key Lessons

1. **Tick off requirements per plan during execution, not during milestone close.** Deferred ticking creates a reconciliation burden and risks requirements appearing incomplete when they are not.
2. **Scan the actual directory structure before the first commit in a new session.** Path assumptions (`src/` vs `app/`) cause easily avoidable failures.
3. **Embed forbidden/confusing terminology lists in CONTEXT.md.** Catching "audit" → "review" and "edit" → "view" at the planning stage costs less than finding them in code review.
4. **Offline scope boundary is a first-class architectural decision.** Establishing "client-only offline" at Phase 3 and holding it through Phase 6 prevented scope creep in every subsequent phase.
5. **Human gates (mobile RTL walkthrough) should be scheduled, not deferred indefinitely.** Documenting them as non-blocking is correct, but they need a concrete v1.1 schedule slot.

### Cost Observations

- Sessions: 1 (autonomous Phase 6 + milestone close workflow)
- Notable: Full 6-phase milestone delivered in a single 5-day autonomous run with zero upstream rework

---

## Cross-Milestone Trends

### Process Evolution

| Milestone | Timeline | Phases | Key Change |
|-----------|----------|--------|------------|
| v1.0 | 5 days | 6 | Baseline — dependency-first sequencing, contract-first typing, bounded offline scope |

### Cumulative Quality

| Milestone | Plans | Coverage | Deferred Items |
|-----------|-------|----------|----------------|
| v1.0 | 27 | 80%+ | 1 (mobile RTL UX walkthrough) |

### Top Lessons (Verified Across Milestones)

1. **Contract-first typing at phase start** — prevents UI/API shape mismatches from propagating. Verified v1.0.
2. **Bounded scope at architectural decision time** — offline scope decision at Phase 3 held cleanly through Phase 6. Verified v1.0.
