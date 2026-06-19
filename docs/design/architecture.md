# git-forge — Initial Architecture

Tracks issue #1. This is the first-pass design that frames the work and seeds the implementation subissues. It is deliberately opinionated where a default is needed, and explicit about what is out of scope.

## Goals

A lightweight, self-hostable git forge — the slice of GitHub that matters for a small team or a solo developer working with agents in the loop. Specifically:

1. **Host repositories** over both `git` and Jujutsu (`jj`), with `jj` treated as a first-class citizen in the UI, not a translation layer over git.
2. **Provide a code review surface** — branches, diffs, comments, merge/PR — that maps cleanly to both `git` and `jj` workflows.
3. **Expose a stable HTTP + webhook surface that agents can drive end-to-end** — clone, read, comment, push, open and update PRs, react to events.
4. **Stay small.** Limited feature set, no enterprise sprawl, no in-product marketing surface, no required JavaScript for read-only browsing.

## Non-goals (for MVP)

- Organizations, teams, fine-grained RBAC — single-user-or-small-group model, owner-and-collaborator only.
- CI runners. We expose hooks; we do not host runners.
- Code search across all repos. Per-repo search via `git grep` / `jj` only.
- Packages, releases-as-artifacts, container registry, pages hosting.
- GitHub Actions equivalent. A webhook surface is the integration point.
- Wikis, discussions, projects.

These can be added later but they multiply scope and we should resist pulling them in early.

## Stack

- **Backend: Go.** Chosen for the spec — also: trivial deployment (single binary), strong stdlib for HTTP, good libraries for git plumbing (`go-git`), and acceptable for shelling out to `jj` and `git` CLIs where the libraries fall short.
- **Frontend: Angular.** Chosen for the spec. Standalone components, signals, the modern SSR-friendly setup. Build emits a static bundle the Go server can serve.
- **Storage:**
  - **Repositories:** bare git repos and jj-colocated repos on local disk under a configurable root (e.g. `/var/lib/git-forge/repos/<owner>/<name>.git`). No object store dependency for MVP.
  - **Metadata (users, issues, PRs, comments, sessions):** SQLite for MVP, with the schema written against `database/sql` so Postgres is a drop-in later. SQLite is single-binary-friendly and handles our expected scale; Postgres is the upgrade path.
- **Auth:** session cookies backed by the metadata DB. Personal access tokens for CLI / agent use. No OAuth in MVP — added later.
- **Transport:** HTTPS (operator-provided cert or behind a reverse proxy). Git over HTTPS via `git-http-backend` (or native Go implementation). SSH out of MVP scope — added once the HTTP path is proven.

### Why not gRPC, GraphQL, etc.

REST + JSON is enough. We have one client (Angular) and a small number of well-known endpoints. Agents prefer REST. We can layer a stricter typed surface later if it becomes painful.

## High-level architecture

```
+--------------------+        +---------------------+
|  Angular SPA       |  --->  |  Go HTTP server     |
|  (static bundle)   |  HTTP  |  - REST API (/api)  |
+--------------------+        |  - Git smart HTTP   |
                              |  - Webhook dispatch |
                              +----------+----------+
                                         |
                       +-----------------+------------------+
                       |                                    |
              +--------v--------+                  +--------v---------+
              |  SQLite (meta)  |                  |  Repo storage    |
              |  users, repos,  |                  |  /var/lib/.../   |
              |  issues, PRs,   |                  |  <owner>/<name>  |
              |  comments,      |                  |  (bare git +     |
              |  tokens, hooks  |                  |   jj colocated)  |
              +-----------------+                  +------------------+
```

Single Go process. Embedded static assets for the SPA. CLI shells out to `git` and `jj` binaries where library coverage is thin (notably `jj`).

## Data model (sketch)

Not final schema — just enough to anchor the subissues.

- `users(id, username, email, password_hash, created_at)`
- `tokens(id, user_id, name, hash, scopes, created_at, last_used_at)`
- `repos(id, owner_id, name, description, default_branch, visibility, vcs, created_at)` — `vcs ∈ {git, jj}`. A jj repo is colocated, so it also has a git view; `vcs` records the primary surface for the UI.
- `refs_cache(repo_id, ref, target, updated_at)` — denormalized for fast listing; rebuilt on push.
- `pulls(id, repo_id, number, title, body, author_id, head_ref, base_ref, state, created_at, merged_at, vcs_kind)` — `vcs_kind` lets jj-bookmark-based PRs coexist with git-branch-based PRs.
- `pull_reviews(id, pull_id, reviewer_id, state, submitted_at)`
- `pull_comments(id, pull_id, author_id, body, path, line, commit_oid, created_at)` — diff-anchored comments.
- `issues(id, repo_id, number, title, body, state, author_id, created_at, closed_at)`
- `issue_comments(id, issue_id, author_id, body, created_at)`
- `webhooks(id, repo_id, url, secret, events, active)`
- `hook_deliveries(id, webhook_id, event, payload, status_code, attempts, next_attempt_at)`

Numbers (`pulls.number`, `issues.number`) are per-repo, GitHub-style.

## Jujutsu integration

This is the part that warrants the most thought, because "git-with-jj-bolted-on" is a worse product than picking one model and making the other a guest.

**Approach:** treat the repo as a single store with two compatible surfaces.

- On disk: colocated jj+git working layout (`jj git init --colocate`) for jj repos; plain bare git for git-only repos.
- The UI's primary unit of work is an **anonymous changeset** for jj-mode repos and a **branch** for git-mode repos. PRs in jj-mode target a **bookmark**; PRs in git-mode target a **branch**. Internally these are the same `pulls` row with `vcs_kind` distinguishing the surface.
- Diffs and review comments are anchored to **commit OIDs**, which both git and jj agree on. This keeps the review system uniform regardless of the source VCS.
- Push protocol: git's smart HTTP for git pushes; jj-over-HTTPS via the standard jj-git remote path (jj pushes through the colocated git remote). No separate jj wire protocol for MVP.
- UI surfaces jj concepts directly when in jj mode: change IDs, working-copy commits, divergent changes, bookmarks (instead of "branches" terminology). The aim is for a jj user to feel that the forge speaks jj, not "git with jj support."

What we are explicitly *not* doing:
- Re-implementing jj's operation log or undo in the UI. Too much surface area for MVP.
- Exposing operations on `@-` working-copy commits server-side. Server is push-target only; jj-on-the-client owns the working copy.

## Agent integration

The forge should be usable by an agent without scraping HTML or exercising the SPA. That means:

1. **A complete REST surface** mirroring every UI action — list/read/comment on issues and PRs, push branches (already covered by smart HTTP), create reviews, request changes, approve, merge.
2. **Personal access tokens** with scoped permissions (`repo:read`, `repo:write`, `pull:comment`, `pull:merge`, `issue:write`). Tokens are the agent auth path.
3. **Webhooks** signed with HMAC-SHA256, with at-least-once delivery and a deliveries table for replay. Events: `push`, `pull.opened`, `pull.synchronized`, `pull.comment`, `pull.review`, `issue.opened`, `issue.comment`.
4. **An MCP server** as a thin layer over the REST API. This is a follow-on subissue, not part of the foundational slice — but the REST shape should be designed knowing an MCP server will sit on top.
5. **Stable event payloads** documented alongside the REST schema. Agents are downstream consumers; payload churn is their tax.

Open question: do we want a built-in "agent identity" type distinct from user accounts (so audit logs distinguish human vs. agent actions, and tokens can be revoked en-masse)? Leaning yes, but punting to a subissue.

## MVP slice

Smallest end-to-end vertical that demonstrates the architecture:

1. Single-user mode (no signup, operator-provisioned account).
2. Create a repo via API.
3. `git push` to the repo over HTTPS with token auth.
4. Browse the repo's default branch in the SPA (tree, file, commit list).
5. Open an issue, comment on it.
6. Open a PR, view the diff, leave a review comment.
7. Merge the PR (fast-forward or merge-commit; squash later).
8. One webhook fires on push, with HMAC signature, retried on failure.

A jj user can repeat steps 3–7 against a jj-mode repo. The UI surfaces jj terminology when `repos.vcs = jj`.

## Phasing into subissues

Each item below is intended to be a tractable PR (or a small number of PRs). Filed as subissues under #1.

1. **Backend: Go HTTP server skeleton.** Module layout, config loading, `/healthz`, embedded static frontend serving, logging, graceful shutdown. No business logic.
2. **Backend: SQLite metadata store and migrations.** `database/sql`-based migration runner, initial schema (users, tokens, repos), CRUD for users/repos used by later subissues.
3. **Backend: Auth — sessions and personal access tokens.** Login endpoint, cookie session middleware, token issuance and verification.
4. **Backend: Repository storage layer.** Disk layout, repo creation, ref enumeration, file/tree/blob reads, `go-git` for in-process operations.
5. **Backend: Git smart HTTP push/pull.** `git-http-backend` or native, token-authenticated.
6. **Backend: Jujutsu integration.** Colocated repo creation, `jj`-specific endpoints (list bookmarks, list changes), shell-out wrapper with timeouts.
7. **Backend: Issues and comments.** CRUD + list endpoints.
8. **Backend: Pull requests, reviews, diff-anchored comments.** Including merge endpoint (fast-forward + merge-commit).
9. **Backend: Webhooks and delivery worker.** Signed payloads, retries, deliveries table.
10. **Frontend: Angular app skeleton.** Routing, layout, auth context, REST client.
11. **Frontend: Repo browsing.** Tree, file view, commit list, branch/bookmark switcher (terminology adapts to `vcs`).
12. **Frontend: Issues UI.**
13. **Frontend: PR UI with diff view and inline comments.**
14. **Ops: Single-binary build and Dockerfile.** Embedded SPA, config via env + file, sample compose.
15. **Docs: REST API reference and webhook payload reference.**
16. **Follow-on: MCP server over the REST API.** Separate concern, tracked but not required for MVP.

Dependencies (rough):
- 1 → 2 → 3 → {4, 7, 10}
- 4 → {5, 6, 8}
- 8 → 9
- 10 → {11, 12, 13}

## Open questions to resolve before coding

These don't block opening the foundational subissues but should be settled before their PRs land.

- **Module path.** Suggesting `github.com/zredinger-ccc/git-forge`. Confirm before the skeleton PR.
- **Angular version + build target.** Latest stable, standalone components, signals. Output `dist/` embedded via `embed.FS` in Go.
- **SQLite driver.** `modernc.org/sqlite` (pure Go, no cgo) — keeps cross-compilation simple. Trade-off: marginally slower than `mattn/go-sqlite3`.
- **`jj` version pinning.** Decide a minimum supported `jj` version and document it; consider invoking with `jj --no-pager --color=never` everywhere for stable parsing.
- **Diff comment anchoring on jj rebases.** When the underlying commit OID rebases away, do comments follow the change ID (jj) or stay pinned to the original OID and be marked "outdated" (git-style)? Leaning git-style for uniformity, but worth deciding deliberately.

## Status

This document is the deliverable for issue #1. Implementation subissues are filed separately and reference this design. Substantive changes to the architecture should land as PRs editing this document.
