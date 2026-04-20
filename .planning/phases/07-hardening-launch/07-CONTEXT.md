# Phase 7: Hardening & Launch - Context

**Gathered:** 2026-04-20
**Status:** Ready for planning
**Mode:** --auto (advanced from `/gsd-next`)

<domain>
## Phase Boundary

Phase 7 turns the completed NutriTrack feature set into a production-ready release. The work is explicitly about hardening, verification, observability, backup/restore readiness, and launch polish — not about adding new product features.

This phase must validate the full stack: Gin API authorization and performance, Nuxt PWA load performance, Docker/Traefik production wiring, Grafana/Loki observability, and backup/restore procedures for PostgreSQL plus uploaded files.
</domain>

<decisions>
## Implementation Decisions

### Security and correctness

- **D-01:** Start with row-level authorization and cross-tenant access checks. Phase 7 should assume authorization bugs are the highest-risk launch blocker.
- **D-02:** Treat dependency and configuration audits as part of the same hardening pass: `govulncheck`, `npm audit`, CORS, HTTPS redirect/HSTS, refresh-token behavior, and file-upload revalidation all belong in this phase.
- **D-03:** Hardening work should prefer proving existing behavior with tests/audits over rewriting stable code paths unless a concrete issue is found.

### Performance and capacity

- **D-04:** Performance verification must use realistic data from the completed product, especially the diet-plan aggregate and offline/PWA startup path.
- **D-05:** The main targets are already locked by the project: API p95 under 200ms, diet plan under 500ms, and PWA initial load under 3 seconds on simulated 3G.
- **D-06:** Query/index fixes and bundle-size optimizations are in scope; new platform dependencies are not, unless an existing target cannot be met otherwise.

### Monitoring, backups, and launch readiness

- **D-07:** Monitoring is self-hosted and Docker-centric: Grafana dashboards, Loki log flow, uptime/health monitoring, and alert thresholds should all match the Hetzner + Docker Compose deployment model.
- **D-08:** Backup work is incomplete until restore is proven. PostgreSQL and upload backups both need an actual restore check, not just cron configuration.
- **D-09:** Final polish includes real Persian UX review, loading/empty/error state audits, and manual Android/iOS launch-path checks using the now-complete Phase 6 PWA.

### Agent's Discretion

- Exact sequencing between security, performance, monitoring, and backup work
- Whether to group audits by subsystem (backend/frontend/infra) or by verification goal
- Which existing commands and test harnesses to extend first, as long as no new product features are introduced
</decisions>

<specifics>
## Specific Ideas

- Reuse existing phase verification artifacts as the checklist baseline instead of rediscovering shipped behavior
- Prioritize endpoints that accept client identifiers or return file content when doing the row-level authorization audit
- Use the completed Phase 6 PWA as the benchmark target for 3G startup and mobile-device launch verification
</specifics>

<canonical_refs>
## Canonical References

- `.planning/ROADMAP.md` §Phase 7 — phase goal, requirements, and success criteria
- `.planning/REQUIREMENTS.md` — `INFRA-05`, `INFRA-07`, `INFRA-08`, `INFRA-09`, `INFRA-11`, `UI-08`
- `docs/phases.md` §Phase 7 — hardening scope, validation checklist, and launch guidance
- `.planning/phases/06-offline-pwa/VERIFICATION.md` — upstream verification state for the completed PWA/offline work
- `backend/cmd/api/main.go` — route surface for authorization/performance audit
- `frontend/nuxt.config.ts` and `frontend/app/service-worker/sw.ts` — PWA startup path and production bundle inputs
- `docker-compose.yml`, `docker-compose.dev.yml`, and Traefik-related config — deployment/observability entry points
</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- Backend already has structured JSON logging, health checks, JWT auth middleware, and sqlc-generated parameterized queries
- Frontend already has a production build, PWA manifest/service worker, offline cache, and notification settings surfaces from Phase 6
- Prior verification files document the main feature boundaries and should be reused when building launch checklists

### Established Patterns
- Handler → service → repository layering on the backend
- Persian-only mobile-first UI on the frontend
- Docker Compose + Traefik deployment assumptions
- GSD phase work is tracked by plan summaries plus phase verification artifacts

### Integration Points
- All completed client flows now cross multiple subsystems (auth, diet plan, tracking, messaging, offline sync, push), so hardening must validate end-to-end role journeys rather than isolated endpoints
- Monitoring and backup work must fit the same Hetzner-hosted Docker topology already chosen for the project
</code_context>

<deferred>
## Deferred Ideas

- New end-user features or product-scope changes
- Native/mobile-only capabilities outside the web/PWA stack
- Large architecture rewrites without evidence from the Phase 7 audits
</deferred>

---

*Phase: 07-hardening-launch*
*Context gathered: 2026-04-20*
