---
phase: 06-offline-pwa
plan: 01
subsystem: ui
tags: [pwa, service-worker, workbox, vite-pwa, nuxt, dexie, offline]

requires:
  - phase: 01-foundation
    provides: nuxt.config.ts base configuration and frontend package.json

provides:
  - "@vite-pwa/nuxt module installed and configured with injectManifest strategy"
  - "Persian-only PWA manifest with shortcuts, RTL, standalone display"
  - "Custom service worker skeleton with precache, NetworkFirst routes, Background Sync stub, push stub"
  - "PWA icons (192, 384, 512px main + 3x 96px shortcut icons) in public/icons/"
  - "dexie@4.4.2 installed and available for Wave 2"

affects:
  - 06-02 (Dexie DB uses dexie package installed here)
  - 06-03 (sw.ts Background Sync handler updated here)
  - 06-07 (sw.ts push handler enhanced in final wave)

tech-stack:
  added:
    - "@vite-pwa/nuxt@1.1.1"
    - "dexie@4.4.2"
    - "workbox-precaching (transitive via vite-plugin-pwa)"
    - "workbox-routing"
    - "workbox-strategies"
    - "workbox-expiration"
  patterns:
    - "injectManifest strategy for custom SW (not generateSW preset)"
    - "srcDir: app/service-worker avoids Pitfall 2 (Rollup compilation path)"
    - "devOptions.enabled=true for SW dev mode testing"

key-files:
  created:
    - frontend/app/service-worker/sw.ts
    - frontend/public/icons/icon-192.png
    - frontend/public/icons/icon-384.png
    - frontend/public/icons/icon-512.png
    - frontend/public/icons/shortcut-plan.png
    - frontend/public/icons/shortcut-track.png
    - frontend/public/icons/shortcut-msg.png
  modified:
    - frontend/nuxt.config.ts (added pwa block + @vite-pwa/nuxt to modules)
    - frontend/package.json (added dexie, @vite-pwa/nuxt)

key-decisions:
  - "Used injectManifest strategy (D-01) for full SW control over caching, push, and Background Sync"
  - "Persian-only manifest: name/short_name='نوتری‌ترک', lang='fa', dir='rtl' (D-04)"
  - "Placeholder 1x1 PNG icons created via Node.js (ImageMagick not available); acceptable for dev"
  - "registerType='autoUpdate' for seamless SW lifecycle (UI-07)"
  - "Added /// <reference lib=\"webworker\" /> for TypeScript SW type support"

patterns-established:
  - "Service worker lives at app/service-worker/sw.ts (compiled by vite-plugin-pwa)"
  - "SW background sync uses TRIGGER_SYNC postMessage pattern to wake useSyncManager"

requirements-completed:
  - UI-06
  - UI-07
  - OFFL-01

duration: ~15min
completed: 2026-04-20
---

# Phase 06-01: PWA Shell Summary

**@vite-pwa/nuxt with injectManifest strategy, Persian manifest with shortcuts, and service worker skeleton precaching all static assets**

## Performance

- **Duration:** ~15 min
- **Completed:** 2026-04-20
- **Tasks:** 3
- **Files modified:** 10

## Accomplishments
- Installed @vite-pwa/nuxt@1.1.1 and dexie@4.4.2; nuxt.config.ts updated with full PWA block
- Created service worker with precacheAndRoute, NetworkFirst for plan + messages APIs, Background Sync stub, and push stub
- Generated 6 placeholder PNG icons (192/384/512px main + 96px shortcuts) for PWA manifest

## Task Commits

1. **Task 1: Install frontend PWA + offline dependencies** - `57bd745` (feat)
2. **Task 2: Configure @vite-pwa/nuxt with Persian manifest** - `57bd745` (feat)
3. **Task 3: Create service worker skeleton and PWA icons** - `57bd745` (feat)

## Files Created/Modified
- `frontend/app/service-worker/sw.ts` - Custom SW with precache, NetworkFirst routes, Background Sync + push stubs
- `frontend/nuxt.config.ts` - Added pwa block with injectManifest strategy and Persian manifest
- `frontend/package.json` - Added @vite-pwa/nuxt and dexie dependencies
- `frontend/public/icons/` - 6 placeholder PNG icons (192, 384, 512, 3x96px)

## Decisions Made
- injectManifest strategy chosen over generateSW for full caching control (D-01)
- Placeholder 1x1 PNGs used since ImageMagick unavailable; acceptable for development
- devOptions.enabled=true so SW can be tested during `npm run dev`

## Deviations from Plan
- Added `/// <reference lib="webworker" />` to sw.ts (not in plan) for TypeScript ServiceWorkerGlobalScope types — necessary to avoid TS errors

## Issues Encountered
None

## User Setup Required
None - no external service configuration required.

## Next Phase Readiness
- dexie package available for Plan 06-02 DB schema creation
- Service worker skeleton ready for Wave 2 runtime caching additions and Wave 3 sync integration

---
*Phase: 06-offline-pwa*
*Completed: 2026-04-20*
