# NutriTrack

## What This Is

NutriTrack is a Persian-only, mobile-first Progressive Web Application (PWA) for managing the relationship between nutritionists (کارشناس تغذیه) and their clients. Nutritionists create personalized diet plans with meals, options, food items, exercise recommendations, and medication prescriptions. Clients view their plans, track daily intake (food, water, sleep, exercise, weight, body measurements, medications), and communicate with their nutritionist — all with full offline support for viewing and data entry.

## Core Value

A structured, digital workflow for nutritionists to manage clients and diet plans, and for clients to track daily health activities with full offline capability — in Persian, on mobile.

## Requirements

### Validated

(None yet — ship to validate)

### Active

- [ ] Super Admin can create, manage, and activate/deactivate nutritionist accounts
- [ ] Super Admin has full CRUD on shared food and medication databases
- [ ] Super Admin can view platform-wide statistics
- [ ] Nutritionists can manage own clients: register, activate, deactivate, view records
- [ ] Nutritionists can create and manage multi-day diet plans with meals, options, food items, exercise recommendations, and medication prescriptions
- [ ] Nutritionists can CRUD on shared food database (visible to all)
- [ ] Nutritionists can CRUD on shared medication database (visible to all)
- [ ] Nutritionists can view all client tracking data (food logs, weight, measurements, exercise, sleep, water, medication intake)
- [ ] Nutritionists can chat with own clients (text, image, file)
- [ ] Nutritionists can set daily water intake target per client
- [ ] Nutritionists can approve/reject food addition requests from clients
- [ ] Clients can view active diet plan with meal options and nutritional info
- [ ] Clients can log daily food intake (select meal option or skip)
- [ ] Clients can track daily water intake with timestamps
- [ ] Clients can track sleep times and quality
- [ ] Clients can track exercise with optional calorie burn
- [ ] Clients can log medication intake with timestamps
- [ ] Clients can record daily weight and body measurements
- [ ] Clients can upload lab results (PDF or link)
- [ ] Clients can chat with assigned nutritionist (text, image, file)
- [ ] Clients can submit food addition requests
- [ ] Full offline support for clients: view plans + record data, sync when online
- [ ] Push notifications via PWA Web Push for reminders and messages
- [ ] Platform is Persian-only (RTL) and mobile-first
- [ ] Authentication: Super Admin/Nutritionist via email+password, Client via mobile+OTP
- [ ] Only one active diet plan per client at a time
- [ ] Nutritional computation: real-time totals per option, meal, and day

### Out of Scope

- Desktop-optimized UI — mobile-only design, no desktop breakpoints
- Multi-language support — Persian only, no i18n infrastructure needed
- Real-time video/voice consultation — not needed for v1
- Payment processing or subscription billing — platform doesn't handle financial transactions
- External health device/wearable integration — manual entry only
- AI-powered diet recommendations — nutritionists make all decisions
- Calorie auto-detection from food photos — manual entry only
- Real-time WebSocket chat — polling is sufficient for this use case
- OAuth login — email/password and OTP are sufficient
- Native mobile app — PWA approach chosen for single codebase

## Context

- **Target market:** Iranian nutritionists and their clients
- **Language:** Persian-only with full RTL layout, Shamsi (Jalali) calendar for all dates
- **User scale:** Up to 50 nutritionists, 10,000 clients, ~500 concurrent users
- **Infrastructure:** Self-hosted on Hetzner with Docker/Docker Compose
- **Monitoring:** Grafana + Loki (existing stack)
- **SMS gateway:** Kavenegar (Iranian SMS provider) for OTP delivery
- **File storage:** Local filesystem on Hetzner (`/data/uploads/`)
- **Existing decisions:** 11 key architectural decisions documented in PRD (shared food DB, client-only offline, OTP for clients, polling for chat, etc.)
- **Food database:** Shared platform-wide with 8 categories, 12 measurement units, full nutritional data (calories, protein, carbs, fat, fiber, sugar, sodium)
- **Diet plan structure:** Deeply nested — Plan → Days → Meals → Options → Items, with prescribed medications and exercise recommendations

## Constraints

- **Tech stack**: Go (Gin framework) + PostgreSQL + Nuxt 4 (Vue 3) — chosen by product owner
- **Design**: Persian-only RTL, mobile-first viewport only, no desktop breakpoints
- **Hosting**: Hetzner dedicated/VPS with Docker + Docker Compose + Traefik
- **CI/CD**: GitLab CI/CD pipeline
- **Auth**: JWT (15-min access + 30-day refresh) with OTP for clients (6-digit, 2-min TTL)
- **Offline**: Client-side only (nutritionist/admin always requires connectivity)
- **Chat**: Polling every 10 seconds (not WebSocket)
- **Files**: Local filesystem storage (not S3/MinIO), max 5MB images, 10MB PDFs
- **Performance**: API < 200ms, diet plan load < 500ms, PWA initial load < 3s on 3G
- **Security**: bcrypt cost 12, parameterized queries, row-level authorization, CORS restricted

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Use Gin for HTTP framework | User preference — explicit request for Gin over Fiber/Echo | — Pending |
| Shared food/medication database across all nutritionists | Reduces duplication, builds richer database over time | — Pending |
| Client-only offline support | Keeps complexity manageable; nutritionists have stable internet | — Pending |
| OTP via SMS for client login | Most Iranian clients prefer SMS; simpler UX on mobile | — Pending |
| Email/password for nutritionist login | Nutritionists need stable sessions for long workflows | — Pending |
| Single active diet plan per client | Avoids confusion; one plan covers food, exercise, medication | — Pending |
| Polling for chat (not WebSocket) | Simpler to implement; acceptable latency; works better offline | — Pending |
| Files stored on Hetzner local disk | Simpler setup; sufficient for current scale | — Pending |
| PWA (not native app) | Single codebase, no app store dependency, easier updates | — Pending |
| Persian-only, no i18n | Target audience is Iranian; YAGNI | — Pending |
| RTL design with ui-ux-pro-max-skill | User specified for UI/UX design approach | — Pending |

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

---
*Last updated: 2026-04-19 after initialization*
