# Repository Guidelines

## Boundary

Repo-local boundary summary:

- this repository is `devflow-config-service`
- public surface is `Configuration`

Authoritative boundary and resource semantics:

- `devflow-control/docs/system/boundaries.md`
- `devflow-control/docs/services/config-service.md`
- `devflow-control/docs/resources/configuration.md`

## Structure

- `cmd/main.go` uses shared bootstrap from `../devflow-service-common`.
- `pkg/api/` contains the configuration handler.
- `pkg/service/` contains configuration CRUD logic.
- `pkg/router/` contains config-only routes and middleware.
- `pkg/config/` initializes config, observability, PostgreSQL, and local store state.
- `docs/` contains the repository-level architecture, API, constraints, observability, and harness docs.

## Required Rules

- Any outbound service or external call must emit `metrics + trace + structured log`.
- Do not add high-cardinality business IDs to metrics labels.
- Default harness is `Planner -> Generator -> Evaluator`.
- When the runtime supports delegation, the harness must spawn those roles as separate sub-agents.
- Non-trivial work should use a run directory under `agents/runs/`.

## Doc And API Hygiene

- Regenerate Swagger after route or handler changes.
- Keep `README.md`, `AGENTS.md`, `agents/protocols/startup.md`, and `docs/*.md` aligned with the actual boundary.
- Do not reintroduce dead service/model/router/bootstrap files.
