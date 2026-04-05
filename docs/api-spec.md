# API Spec

## Purpose

`devflow-config-service` defines the converged public API surface for:

- `Configuration`
- `ConfigurationRevision`

## Endpoint Groups

### `Configuration`
- `GET /api/v1/configurations`
- `POST /api/v1/configurations`
- `GET /api/v1/configurations/{id}`
- `POST /api/v1/configurations/{id}/sync`
- `PUT /api/v1/configurations/{id}`
- `DELETE /api/v1/configurations/{id}`

## Request Rules

- list endpoints use `page` and `page_size`
- `POST` and `PUT` use request DTOs, not raw domain models
- writable fields are `application_id`, `name`, `env`, and `source_path`
- `POST /api/v1/configurations/{id}/sync` freezes the current config-repo snapshot into an immutable revision

## Response Rules

- create returns `201` with `{ "data": ... }`
- get returns `200` with `{ "data": ... }`
- list returns `200` with `{ "data": [...], "pagination": { ... } }`
- update and delete return `204`

## Error Rules

- invalid ID or request body -> `400 invalid_argument`
- missing config repo wiring or missing source path -> `424 failed_precondition`
- resource not found -> `404 not_found`
- storage or uncategorized internal error -> `500 internal`

## Boundary Note

For repo scope and non-goals, see `docs/architecture.md`.

## Swagger Note

Generated Swagger artifacts must stay aligned with the current PostgreSQL-backed API contract. Regenerate them after route, request, or response changes.
