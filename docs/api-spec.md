# API Spec

## Purpose

`devflow-config-service` defines the converged public API surface for:

- `AppConfig`
- `WorkloadConfig`

## Endpoint Groups

### `AppConfig`
- `GET /api/v1/app-configs`
- `POST /api/v1/app-configs`
- `GET /api/v1/app-configs/{id}`
- `PUT /api/v1/app-configs/{id}`
- `DELETE /api/v1/app-configs/{id}`
- `POST /api/v1/app-configs/{id}/sync-from-repo`

### `WorkloadConfig`
- `GET /api/v1/workload-configs`
- `POST /api/v1/workload-configs`
- `GET /api/v1/workload-configs/{id}`
- `PUT /api/v1/workload-configs/{id}`
- `DELETE /api/v1/workload-configs/{id}`

## Swagger

- local UI: `/swagger/index.html`
- generated source: `docs/generated/swagger/swagger.yaml`

## Request Rules

- list endpoints use `page` and `page_size`
- `POST` and `PUT` use request DTOs, not raw domain models
- `AppConfig` writable fields are `application_id`, `environment_id`, and `name`
- `AppConfig` source repo is fixed to `git@github.com:bsonger/devflow-config-repo.git`
- `AppConfig` source branch is fixed to `main`
- `AppConfig` source path is derived from `name` as `applications/devflow-platform/services/<name>`
- `POST /api/v1/app-configs/{id}/sync-from-repo` reads `configuration.yaml`, `deployment.yaml`, `service.yaml`, plus optional `environments/{env}.yaml`
- `POST /api/v1/app-configs/{id}/sync-from-repo` pulls the latest `origin/main`, then freezes the current repo snapshot into an immutable revision
- `WorkloadConfig` writable fields are `application_id`, optional `environment_id`, `name`, `replicas`, `resources`, `probes`, `env`, `workload_type`, and `strategy`

## Response Rules

- create returns `201` with `{ "data": ... }`
- get returns `200` with `{ "data": ... }`
- `GET /api/v1/app-configs/{id}` includes latest synced revision payload when present:
  - `files`
  - `rendered_configmap`
  - `source_commit`
- list returns `200` with `{ "data": [...], "pagination": { ... } }`
- update and delete return `204`

## Error Rules

- invalid ID or request body -> `400 invalid_argument`
- missing config repo checkout, repo pull failure, or missing derived source path -> `424 failed_precondition`
- resource not found -> `404 not_found`
- storage or uncategorized internal error -> `500 internal`

## Boundary Note

For repo scope and non-goals, see `docs/architecture.md`.

## Swagger Note

Generated Swagger artifacts must stay aligned with the current PostgreSQL-backed API contract. Regenerate them after route, request, or response changes.
