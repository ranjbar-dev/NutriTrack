---
phase: 06
slug: admin-governance
status: completed
nyquist_compliant: true
wave_0_complete: true
created: 2026-04-24
---

# Phase 06 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | vitest |
| **Config file** | `vitest.config.ts` |
| **Quick run command** | `npm run test:unit -- tests/admin/admin-api-contracts.spec.ts` |
| **Full suite command** | `npm run test:unit` |
| **Estimated runtime** | ~20 seconds |

---

## Sampling Rate

- **After every task commit:** Run the narrowest relevant admin spec first, then add `tests/auth/route-access-control.spec.ts` only when route-guard behavior is touched
- **After every plan wave:** Run `npm run test:unit`
- **Before `/gsd-verify-work`:** Full suite must be green
- **Max feedback latency:** 30 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Threat Ref | Secure Behavior | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|------------|-----------------|-----------|-------------------|-------------|--------|
| 06-01-01 | 01 | 1 | ADMIN-01 | T-06-01 / V4 | Admin stats and nutritionist API wrappers use documented `/api/v1/admin/*` paths only. | unit | `npm run test:unit -- tests/admin/admin-api-contracts.spec.ts` | ✅ | ✅ green |
| 06-02-01 | 02 | 2 | ADMIN-01 | T-06-02 / V4 | `/admin/**` pages remain role-scoped and render only API-backed metrics. | unit/component | `npm run test:unit -- tests/admin/admin-dashboard-roster.spec.ts tests/auth/route-access-control.spec.ts` | ✅ | ✅ green |
| 06-03-01 | 03 | 3 | ADMIN-01 | T-06-03 / V5 | Nutritionist create/update/status flows validate inputs and require explicit confirmation for status changes. | unit/component | `npm run test:unit -- tests/admin/admin-nutritionist-detail.spec.ts tests/auth/route-access-control.spec.ts` | ✅ | ✅ green |
| 06-04-01 | 04 | 2 | ADMIN-02 | T-06-04 / V4 | Admin catalogue pages use elevated search/delete/category endpoints with destructive confirmations. | unit/component | `npm run test:unit -- tests/admin/admin-catalogue-governance.spec.ts tests/auth/route-access-control.spec.ts` | ✅ | ✅ green |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [x] `tests/admin/admin-api-contracts.spec.ts` — stats and nutritionist admin endpoint contract coverage
- [x] `tests/admin/admin-dashboard-roster.spec.ts` — admin dashboard plus roster/create coverage
- [x] `tests/admin/admin-nutritionist-detail.spec.ts` — nutritionist detail, status, and read-only client-list coverage
- [x] `tests/admin/admin-catalogue-governance.spec.ts` — elevated foods/medications/categories governance coverage

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| Mobile RTL admin information density remains usable across stats, roster, and destructive actions | ADMIN-01, ADMIN-02 | Visual ergonomics and sheet/card spacing are hard to assert fully in unit tests | Open admin dashboard and catalogue pages in a mobile viewport, verify KPI readability, confirmation copy, and safe-area behavior. |

---

## Validation Sign-Off

- [x] All tasks have `<automated>` verify or Wave 0 dependencies
- [x] Sampling continuity: no 3 consecutive tasks without automated verify
- [x] Wave 0 covers all MISSING references
- [x] No watch-mode flags
- [x] Feedback latency < 60s
- [x] `nyquist_compliant: true` set in frontmatter

**Approval:** approved 2026-04-24