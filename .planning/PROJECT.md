# NutriTrack — Go Backend API

## What This Is

NutriTrack is a RESTful Go backend API for a Persian nutrition management platform that connects nutritionists with their clients. Nutritionists create personalized diet plans, prescribe medications, and monitor client progress. Clients track daily intake, exercise, sleep, water, and body measurements. The backend serves a Nuxt.js PWA frontend built exclusively for Iranian mobile users.

This project implements the **backend only** — all HTTP handlers, business logic, database interactions, and infrastructure configuration.

## Core Value

A nutritionist must be able to create a diet plan and assign it to a client — everything else serves this workflow.

## Requirements

### Validated

(None yet — ship to validate)

### Active

- [ ] JWT + OTP authentication for three user roles (super admin, nutritionist, client)
- [ ] Shared food and medication database with full CRUD and Persian full-text search
- [ ] Diet plan creation with nested structure (days → meals → options → food items) and computed nutritional totals
- [ ] Daily tracking APIs for food logs, water, sleep, exercise, medications, and body measurements
- [ ] Lab result upload and management (PDF/image, stored on local filesystem)
- [ ] Chat/messaging system between clients and nutritionists (polling-based)
- [ ] Food addition request workflow (client → nutritionist approval)
- [ ] Web Push notifications (VAPID) for reminders and messages
- [ ] Super admin panel APIs for nutritionist and database management
- [ ] Offline sync support — idempotent APIs with local_id deduplication
- [ ] Docker + Docker Compose deployment with PostgreSQL and Redis
- [ ] Persian error messages throughout all API responses
- [ ] Asia/Tehran timezone for all services and timestamp handling

### Out of Scope

- Frontend / PWA / Nuxt.js — backend only
- Real-time WebSocket chat — use polling as decided in PRD
- Payment processing or subscription billing
- AI-powered diet recommendations
- External health device integrations
- Multi-language support — Persian only, no i18n layer needed
- Desktop-optimized UI concerns

## Context

- **Target audience:** Iranian nutritionists and their clients; all UI copy, error messages, and user-facing strings in Persian (Farsi)
- **Timezone:** All timestamps stored in UTC, all time-based logic uses Asia/Tehran (UTC+3:30, DST-aware)
- **Backend language:** Go (Golang) with Gin HTTP framework
- **Database queries:** sqlc (type-safe generated queries from SQL) — no GORM or ORM
- **Databases:** PostgreSQL (primary data store) + Redis (OTP storage, rate limiting, caching, session tokens)
- **Design pattern:** Domain-Driven Design (DDD) — entities, value objects, aggregates, repositories, domain services, application services
- **Logging:** Structured JSON logging (zerolog or zap) to stdout, collected by Loki
- **Auth:** JWT access tokens (15 min) + refresh tokens (30 days) stored in Redis; OTP via SMS (Kavenegar adapter)
- **File storage:** Local filesystem on Hetzner (`/data/uploads/`) with path in DB
- **Push notifications:** Web Push via github.com/SherClockHolmes/webpush-go
- **Migrations:** golang-migrate with versioned SQL migration files
- **Infrastructure:** Docker + Docker Compose on Hetzner; Traefik as reverse proxy
- **Scale targets:** ~50 nutritionists, ~10,000 clients, ~500 concurrent users
- **SMS gateway:** Iranian SMS (Kavenegar/Melipayamak) configurable via adapter pattern

## Constraints

- **Tech Stack — Go + Gin:** Specified by user; no alternative web frameworks
- **Tech Stack — sqlc:** All DB queries via sqlc-generated code; no raw query strings or ORM
- **Tech Stack — Redis:** Required for OTP, rate limiting, token invalidation, and caching
- **Tech Stack — PostgreSQL:** Primary data store with pg_trgm for Persian full-text search
- **Tech Stack — Docker:** All services containerized; Docker Compose for local and production
- **Design — DDD:** Domain-Driven Design strictly applied; no anemic domain models
- **Language — Persian errors:** All `message` fields in API error responses must be in Persian
- **Timezone — Asia/Tehran:** TZ environment variable set to Asia/Tehran in all containers
- **Backend only:** No frontend code; pure JSON REST API
- **Security:** bcrypt (cost 12) for passwords, JWT short expiry, OTP rate limiting, row-level authorization

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Gin over Fiber/Echo | User requirement | — Pending |
| sqlc over GORM | Type-safe queries, no reflection overhead, explicit SQL | — Pending |
| DDD architecture | Clean separation of domain logic from infrastructure | — Pending |
| Redis for OTP + rate limiting | Fast TTL-based storage, atomic counters | — Pending |
| Polling for chat (not WebSocket) | Simpler, works offline, PRD decision #7 | — Pending |
| Persian-only error messages | Target audience is Iranian, PRD decision #11 | — Pending |
| Local filesystem for uploads | PRD decision #8 — simpler, Hetzner has sufficient disk | — Pending |
| JWT + OTP hybrid auth | Nutritionists use email/password; clients use SMS OTP | — Pending |
| golang-migrate for migrations | Version-controlled SQL files, works with sqlc | — Pending |

---
*Last updated: 2026-04-21 after initialization*

## Evolution

This document evolves at phase transitions and milestone boundaries.

**After each phase transition** (via `/gsd-transition`):
1. Requirements invalidated? → Move to Out of Scope with reason
2. Requirements validated? → Move to Validated with phase reference
3. New requirements emerged? → Add to Active
4. Decisions to log? → Add to Key Decisions
5. "What This Is" still accurate? → Update if drifted

**After each milestone** (via `/gsd-complete-milestone`):
1. Full review of all sections
2. Core Value check — still the right priority?
3. Audit Out of Scope — reasons still valid?
4. Update Context with current state
