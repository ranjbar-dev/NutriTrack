---
phase: 01-foundation-infrastructure
plan: 06
title: "Docker, Traefik & GitLab CI/CD"
subsystem: infrastructure
tags: [docker, traefik, gitlab-ci, compose, dockerfile, devops]
dependency_graph:
  requires: [01-01, 01-02]
  provides: [docker-infrastructure, traefik-routing, gitlab-ci-pipeline, env-documentation]
  affects: [all-subsequent-plans]
tech_stack:
  added: [traefik@v3.2, postgres@16-alpine, docker-compose, gitlab-ci]
  patterns: [multi-stage-docker-build, scratch-base-image, path-based-routing, acme-http-challenge, manual-deploy-gate]
key_files:
  created:
    - backend/Dockerfile
    - backend/.dockerignore
    - frontend/Dockerfile
    - frontend/.dockerignore
    - docker-compose.yml
    - docker-compose.dev.yml
    - .env.example
    - .gitlab-ci.yml
  modified: []
decisions:
  - "No version key in docker-compose (Docker Compose v2+ ignores it; avoids deprecation warning)"
  - "Frontend priority=1 in Traefik labels ensures /api/* matches api router first via longer prefix"
metrics:
  duration: "4 minutes"
  completed: "2026-04-19T15:59:30Z"
---

# Phase 01 Plan 06: Docker, Traefik & GitLab CI/CD Summary

Multi-stage Dockerfiles (Go scratch ~20MB, Nuxt Node 22 Alpine), production docker-compose with Traefik v3 path-based routing and Let's Encrypt TLS, dev compose with raw ports, and GitLab CI/CD pipeline with lint → test → build → deploy stages and manual deploy gate.

## What Was Done

### Task 1: Multi-stage Dockerfiles for Go and Nuxt
- Created `backend/Dockerfile` with multi-stage build: golang:1.25-alpine builder → scratch runtime
- `CGO_ENABLED=0 GOOS=linux GOARCH=amd64` for static binary, `-ldflags="-s -w"` strips symbols (~20MB final)
- Copies ca-certificates (outbound HTTPS for Kavenegar) and migrations directory
- Created `frontend/Dockerfile` with multi-stage build: node:22-alpine builder → node:22-alpine runtime
- `npm ci` for reproducible builds, copies only `.output/` to production stage
- Both `.dockerignore` files exclude `.git`, `node_modules`, `.env`, `*.md`, etc.
- **Commit:** `4d631e9`

### Task 2: Docker Compose (production + dev) with Traefik and .env.example
- Created `docker-compose.yml` with 4 services: traefik, api, frontend, postgres
- Traefik v3.2 configured entirely via Docker labels (no separate config file, per D-20)
- Path-based routing: `/api/*` → Go backend (port 8080), `/*` → Nuxt frontend (port 3000)
- Let's Encrypt ACME HTTP challenge for automatic TLS certificate renewal
- HTTP → HTTPS redirect enforced via entrypoint redirection
- Docker socket mounted read-only (`:ro`) per T-06-02
- PostgreSQL not exposed to host in production (internal Docker network only, per T-06-05)
- PostgreSQL health check with `pg_isready` and `service_healthy` condition
- Created `docker-compose.dev.yml` with 3 services (no Traefik): raw ports (8080, 3000, 5432), default env values, separate volume (`pgdata_dev`)
- Created root `.env.example` documenting all 14 environment variables with descriptions and generation hints
- **Commit:** `3f488ef`

### Task 3: GitLab CI/CD pipeline
- Created `.gitlab-ci.yml` with 4 stages: lint, test, build, deploy (per D-21, INFRA-03)
- **Lint stage:** `go vet` + optional `golangci-lint` for backend; `nuxi typecheck` for frontend
- **Test stage:** `go test -v -race -count=1` with PostgreSQL 16 service; `npm test` for frontend
- **Build stage:** Docker-in-Docker image builds pushed to GitLab Container Registry (main branch only)
- **Deploy stage:** SSH-based deployment to Hetzner, `docker compose pull + up` (manual gate for safety)
- Conditional execution: backend/frontend jobs only trigger when respective files change
- Frontend `node_modules` cached between pipeline runs
- **Commit:** `10e7a8e`

## Deviations from Plan

None - plan executed exactly as written.

## Verification Results

```
Test-Path backend/Dockerfile       → True
Test-Path frontend/Dockerfile      → True
Test-Path backend/.dockerignore    → True
Test-Path frontend/.dockerignore   → True
Test-Path docker-compose.yml       → True
Test-Path docker-compose.dev.yml   → True
Test-Path .env.example             → True
Test-Path .gitlab-ci.yml           → True

Production compose: 4 services (traefik, api, frontend, postgres) → PASS
Dev compose: 3 services (api, frontend, postgres, no Traefik)     → PASS
.env.example: All required vars documented                        → PASS
GitLab CI: 4 stages, 7 jobs, manual deploy gate                   → PASS
```

## Self-Check: PASSED

All 8 created files verified present. All 3 task commits (4d631e9, 3f488ef, 10e7a8e) verified in git log. No unexpected file deletions.
