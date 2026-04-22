# DDD Compliance Check & Fix — NutriTrack Backend

## Objective
Audit every sub-package in `internal/application`, `internal/domain`, `internal/infrastructure`, and `internal/interfaces` for Domain-Driven Design (DDD) compliance based on the patterns from https://programmingpercy.tech/blog/how-to-domain-driven-design-ddd-golang/ — then automatically plan and apply fixes for any violations found.

## DDD Rules Summary (reference for all agents)

| Layer | Expected Contents | Forbidden |
|-------|------------------|-----------|
| `internal/domain` | Entities (UUID+unexported fields), Value Objects (immutable+unexported), Aggregates (unexported fields+no DB tags+`New*()` factory), Repository interfaces, domain errors | Concrete repos, DB tags on aggregates, infra/app/http imports |
| `internal/application` | Service structs accepting domain repository interfaces, `New*Service()` factory, use-case methods | SQL/DB calls, infra imports, HTTP types |
| `internal/infrastructure` | Concrete repository implementations, internal DB DTOs with tags, `NewFromAggregate()`/`ToAggregate()` converters | Business logic, domain rule enforcement |
| `internal/interfaces` | HTTP handlers delegating to app services, request/response DTOs with json tags, input validation | Direct infra imports, domain logic |

## Execution Plan

### Wave 1 — Parallel Audit (spawn all at once)

Spawn the following agents **in parallel** using the `ddd-auditor` agent from `.github/agents/ddd-auditor.agent.md`. Each agent audits ONE folder:

**Domain layer:**
1. Audit `internal/domain/user`
2. Audit `internal/domain/food`
3. Audit `internal/domain/dietplan`
4. Audit `internal/domain/medication`
5. Audit `internal/domain/tracking`
6. Audit `internal/domain/labresult`
7. Audit `internal/domain/message`
8. Audit `internal/domain/messaging`
9. Audit `internal/domain/foodrequest`
10. Audit `internal/domain/notification`
11. Audit `internal/domain/push`
12. Audit `internal/domain/shared`

**Application layer:**
13. Audit `internal/application/user`
14. Audit `internal/application/auth`
15. Audit `internal/application/food`
16. Audit `internal/application/dietplan`
17. Audit `internal/application/medication`
18. Audit `internal/application/tracking`
19. Audit `internal/application/labresult`
20. Audit `internal/application/message`
21. Audit `internal/application/messaging`
22. Audit `internal/application/foodrequest`
23. Audit `internal/application/notification`
24. Audit `internal/application/push`
25. Audit `internal/application/admin`

**Infrastructure layer:**
26. Audit `internal/infrastructure/persistence`
27. Audit `internal/infrastructure/push`
28. Audit `internal/infrastructure/redis`
29. Audit `internal/infrastructure/sms`
30. Audit `internal/infrastructure/storage`

**Interfaces layer:**
31. Audit `internal/interfaces/http`

Each auditor agent MUST:
- Read all `.go` files in its assigned folder (recursively)
- Apply layer-specific DDD rules
- Write `DDD-AUDIT.md` into its target folder
- Return a one-line summary: `{folder}: PASS | {n} CRITICAL, {n} HIGH, {n} MEDIUM, {n} LOW`

---

### Wave 2 — Collect & Triage

After ALL Wave 1 agents complete:

1. Read every `DDD-AUDIT.md` produced
2. Build a consolidated triage table:

```
| Folder | CRITICAL | HIGH | MEDIUM | LOW | Action |
|--------|----------|------|--------|-----|--------|
| internal/domain/user | 0 | 1 | 2 | 0 | FIX |
| internal/domain/food | 0 | 0 | 0 | 0 | SKIP |
...
```

3. Any folder with CRITICAL > 0 OR HIGH > 0 → queue for Wave 3 (fix)
4. Folders with only MEDIUM/LOW → queue for Wave 4 (fix, lower priority)
5. PASS folders → no action needed

---

### Wave 3 — Fix CRITICAL & HIGH Violations (parallel, per folder)

For each folder queued from Wave 2 with CRITICAL or HIGH findings:

Spawn a `ddd-fixer` agent (from `.github/agents/ddd-fixer.agent.md`) per folder. Each fixer MUST:

1. Run `go build ./...` as baseline BEFORE making any changes
2. Read the folder's `DDD-AUDIT.md`
3. Read all referenced source files
4. Write `DDD-FIX.md` Plan section before touching any code
5. Apply fixes ONE at a time — CRITICAL first, then HIGH
6. After each edit: run `go build ./internal/...` to verify compilation
7. On build failure: `git checkout -- {file}`, mark as SKIPPED, continue
8. Never rename exported symbols without updating ALL callers in the same change
9. Never fix across more than 3 files in a single finding — DEFER if broader
10. Write final `DDD-FIX.md` with every change made or deferred

---

### Wave 4 — Fix MEDIUM & LOW Violations

After Wave 3 completes and all builds pass:

Spawn `ddd-fixer` agents for folders with only MEDIUM/LOW findings.
Same protocol as Wave 3.

---

### Wave 5 — Final Verification

After all fixers complete:

1. Run `go build ./...` — must PASS
2. Run `go vet ./internal/...` — must PASS  
3. Run `go test ./internal/...` — record results (do not fail on pre-existing test failures)
4. Print consolidated summary:
   - Total folders audited
   - Folders with no violations (PASS)
   - CRITICAL findings: found / fixed / deferred
   - HIGH findings: found / fixed / deferred
   - MEDIUM findings: found / fixed / deferred
   - LOW findings: found / fixed / deferred
   - Build status: PASS/FAIL
5. List any DEFERRED items that need human review

---

## Agent Instructions

### For each `ddd-auditor` agent invocation, use this prompt template:
```
You are the ddd-auditor agent. Load your instructions from `.github/agents/ddd-auditor.agent.md`.

Target folder: `{FOLDER_PATH}`
Project root: `c:\Users\root\Desktop\Projects\github\NutriTrack`
Go module: `github.com/ranjbar-dev/nutritrack`

Audit ALL `.go` files in the target folder for DDD compliance.
Write your findings to `{FOLDER_PATH}/DDD-AUDIT.md`.
Return a one-line summary when done.
```

### For each `ddd-fixer` agent invocation, use this prompt template:
```
You are the ddd-fixer agent. Load your instructions from `.github/agents/ddd-fixer.agent.md`.

Target folder: `{FOLDER_PATH}`
Project root: `c:\Users\root\Desktop\Projects\github\NutriTrack`  
Go module: `github.com/ranjbar-dev/nutritrack`
Audit report: `{FOLDER_PATH}/DDD-AUDIT.md`

Fix all CRITICAL and HIGH violations in the audit report without breaking compilation.
Write your changes to `{FOLDER_PATH}/DDD-FIX.md`.
Return a one-line summary when done.
```

---

## Constraints for All Agents

- **Never break compilation** — verify with `go build` after every change
- **Never rename exported types or functions** without updating all callers in the same atomic change
- **Never delete code** — use `// DDD-TODO:` comments for code to be removed later
- **Never add new external dependencies** — only use packages already in `go.mod`
- **Prefer minimal changes** — fix the specific violation, do not refactor beyond the audit finding
- **Defer risky changes** — if a fix touches more than 3 files, add to DEFERRED list in DDD-FIX.md
