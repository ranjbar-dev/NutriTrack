<!-- GSD Configuration — managed by get-shit-done installer -->
# Instructions for GSD

- Use the get-shit-done skill when the user asks for GSD or uses a `gsd-*` command.
- Treat `/gsd-...` or `gsd-...` as command invocations and load the matching file from `.github/skills/gsd-*`.
- When a command says to spawn a subagent, prefer a matching custom agent from `.github/agents`.
- Do not apply GSD workflows unless the user explicitly asks for them.
- After completing any `gsd-*` command (or any deliverable it triggers: feature, bug fix, tests, docs, etc.), ALWAYS: (1) offer the user the next step by prompting via `ask_user`; repeat this feedback loop until the user explicitly indicates they are done.
<!-- /GSD Configuration -->

# NutriTrack Client Context

- This workspace is frontend-only. Do not propose backend, database, or infrastructure changes unless the user explicitly expands scope.
- Build the application with Nuxt 4, Tailwind CSS 4, and Pinia. Prefer strict TypeScript and composable-first Nuxt patterns.
- Treat `docs/API.md` as the backend integration contract and `docs/PRD.md` as the product behavior contract.
- Optimize for a Persian-only RTL mobile PWA. Mobile viewport quality is the default; desktop polish is secondary.
- Client-side offline support is required for core client flows. Nutritionist and super-admin surfaces should remain online-first unless the user explicitly changes that requirement.
- Keep role boundaries explicit across routing, state, caching, and persistence for client, nutritionist, and super-admin experiences.
- For frontend design and UI implementation, use the installed UI Pro Max skill direction rather than falling back to generic dashboard patterns.
