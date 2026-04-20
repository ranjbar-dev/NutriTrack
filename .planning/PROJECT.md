# NutriTrack

## What This Is

NutriTrack is a Persian-only, mobile-first Progressive Web Application (PWA) for managing the relationship between nutritionists and their clients. It covers the full v1 workflow: authentication for all roles, shared food and medication data, diet plan authoring, client tracking, messaging, offline/PWA behavior, push reminders, and production hardening.

## Core Value

A structured, digital workflow for nutritionists to manage clients and diet plans, and for clients to track daily health activities with offline-capable mobile access in Persian.

## Current State

- **Version:** v1.0 Launch
- **Status:** Shipped on 2026-04-20
- **Codebase state:** All 7 roadmap phases are implemented, summarized, and archived for milestone close
- **Operational follow-up:** Real-device launch validation, live backup restore proof, staging load evidence, and live Grafana/Loki traffic proof still need to be executed outside this CLI environment

## Requirements

### Validated

- ✓ Multi-role authentication, authorization, and Persian RTL/mobile foundations — v1.0
- ✓ Shared food/medication management, admin controls, and nutritionist client management — v1.0
- ✓ Diet plan builder, active/archived plan delivery, and nutrition computation — v1.0
- ✓ Client tracking, lab uploads, messaging, and food request workflows — v1.0
- ✓ Offline caching, sync queue, installable PWA shell, and push notification flows — v1.0
- ✓ Launch hardening: security middleware/audits, backup/restore scripts, observability stack, and launch-critical UX polish — v1.0

### Active

None. Define a new milestone before adding more scope.

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
- **Language:** Persian-only with RTL layout, Shamsi/Jalali dates, Persian numerals
- **Stack:** Gin + PostgreSQL + Nuxt 4 + Tailwind v4 + Dexie + Docker Compose + Traefik
- **Hosting model:** Self-hosted on Hetzner with local filesystem storage
- **Observability:** Grafana, Loki, and Promtail provisioning shipped in-repo
- **Launch-readiness evidence still pending:** Physical-device UAT, staging restore proof, staging load validation, live observability capture

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Use Gin for the backend | Matches product-owner directive and middleware ecosystem needs | ✓ Good |
| Keep the app Persian-only and mobile-first | Aligns directly with the target audience and removes unnecessary i18n/desktop overhead | ✓ Good |
| Use sqlc + PostgreSQL instead of an ORM | Keeps query behavior explicit and safe for complex nested plan aggregates | ✓ Good |
| Restrict offline support to clients only | Maximizes product value without adding nutritionist/admin sync complexity | ✓ Good |
| Use polling instead of WebSockets for chat | Simpler infrastructure and compatible with the product's latency expectations | ✓ Good |
| Use a PWA instead of native mobile apps | Single codebase, easier updates, and enough capability for install/offline/push needs | ✓ Good |
| Store files on local disk | Fits current scale and hosting model without introducing object-storage complexity | ✓ Good |
| Add launch hardening in a dedicated final phase | Kept security, observability, and backup work visible instead of burying it in feature phases | ✓ Good |

## Next Milestone Goals

No next milestone is planned yet. If post-launch work is approved, begin with `/gsd-new-milestone` and treat the remaining operational sign-off items as launch evidence rather than feature scope.

---
*Last updated: 2026-04-20 after v1.0 Launch milestone close*
