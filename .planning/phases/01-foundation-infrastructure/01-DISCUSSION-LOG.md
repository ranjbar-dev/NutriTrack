# Phase 1: Foundation & Infrastructure - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2025-07-18
**Phase:** 01-foundation-infrastructure
**Areas discussed:** Authentication Architecture, Project Structure, Persian RTL Foundation, Role-based Layouts & Navigation, Deployment & Infrastructure, Security Foundation
**Mode:** --auto (all decisions auto-selected with recommended defaults)

---

## Authentication Architecture

| Option | Description | Selected |
|--------|-------------|----------|
| httpOnly cookies | JWT stored in secure httpOnly cookies — prevents XSS token theft | ✓ |
| localStorage | JWT in localStorage — simpler but vulnerable to XSS | |
| Memory + refresh cookie | Access token in memory, refresh in httpOnly cookie — hybrid approach | |

**User's choice:** httpOnly cookies (recommended default — strongest XSS protection)
**Notes:** Auto-selected. Refresh token queue pattern chosen for concurrent request handling (AUTH-09). Super Admin seeded via CLI command, not migration.

---

## Project Structure

| Option | Description | Selected |
|--------|-------------|----------|
| Monorepo | Single repo with backend/ and frontend/ directories | ✓ |
| Separate repos | Separate Git repos for Go API and Nuxt frontend | |

**User's choice:** Monorepo (recommended — simpler CI/CD, shared Docker Compose)
**Notes:** Auto-selected. Handler → Service → Repository layered architecture for Go. Nuxt 4 app/ directory structure.

---

## Persian RTL Foundation

| Option | Description | Selected |
|--------|-------------|----------|
| Tailwind v4 logical properties | Use ms-/me-/ps-/pe- throughout, dir="rtl" on root | ✓ |
| RTL plugin | Use a separate RTL plugin to auto-flip styles | |

**User's choice:** Tailwind v4 logical properties (recommended — native support, no plugin needed)
**Notes:** Auto-selected. Vazirmatn via npm, jalaali-js for Shamsi dates, custom toPersianDigits utility.

---

## Role-based Layouts & Navigation

| Option | Description | Selected |
|--------|-------------|----------|
| Bottom navigation bar | Mobile bottom nav with role-specific tabs | ✓ |
| Hamburger menu | Side drawer navigation | |
| Tab bar + header nav | Combined approach | |

**User's choice:** Bottom navigation bar (recommended — standard mobile-first pattern)
**Notes:** Auto-selected. Three separate Nuxt layouts. Route middleware for role-based access.

---

## Deployment & Infrastructure

| Option | Description | Selected |
|--------|-------------|----------|
| Docker Compose + Traefik | Full container orchestration with auto HTTPS | ✓ |
| Docker Compose + nginx | Traditional reverse proxy setup | |

**User's choice:** Docker Compose + Traefik (recommended — auto Let's Encrypt, Docker-native)
**Notes:** Auto-selected. 4-stage GitLab CI/CD pipeline. Multi-stage Go build for slim images.

---

## Security Foundation

| Option | Description | Selected |
|--------|-------------|----------|
| In-memory rate limiter | Simple rate limiting sufficient for v1 scale | ✓ |
| Redis-backed rate limiter | Distributed rate limiting for multi-instance | |

**User's choice:** In-memory rate limiter (recommended — sufficient for single-instance v1)
**Notes:** Auto-selected. CORS restricted to frontend domain. All queries parameterized via sqlc.

---

## Agent's Discretion

- Loading skeleton/spinner design
- Exact Tailwind color palette and design tokens
- Error page styling
- Exact CI/CD runner tags and Docker registry choice
- Go module path naming
- Health check response fields beyond status

## Deferred Ideas

- Persian pg_trgm search validation — Phase 2
- Plan builder UI state management — Phase 3
- iOS PWA storage eviction — Phase 6
