# API Spec

## Purpose

`devflow-config-service` defines the converged public API surface for:

- `Configuration`
- `ConfigurationRevision`

## Endpoint Groups

### `Configuration`
- `GET /api/v1/configurations`
- `POST /api/v1/configurations`
- `GET /api/v1/configurations/:id`
- `PUT /api/v1/configurations/:id`
- `DELETE /api/v1/configurations/:id`

### `ConfigurationRevision`
- `GET /api/v1/configurations/:id/revisions`
- `POST /api/v1/configurations/:id/revisions`
- `GET /api/v1/configurations/:id/revisions/latest`
- `GET /api/v1/configuration-revisions/:id`

## Request Rules

- list endpoints support the common pagination parameters used in this repo
- `POST /api/v1/configurations` creates the logical configuration and its first revision together
- `Configuration` owns identity fields such as `application_id`, `name`, `env`, and revision pointers
- revisions are immutable once created
- environment variables belong to `ConfigurationRevision.env_vars`
- release flows should consume a revision, not mutable configuration state directly

## Response Rules

- create endpoints return the common create-response shape
- list endpoints return pagination headers
- success payloads must stay aligned with Swagger

## Error Rules

- invalid ID -> `400`
- resource not found -> `404`
- storage or uncategorized internal error -> `500`

## Boundary Note

For repo scope and non-goals, see `docs/architecture.md`.

## Swagger Note

Generated Swagger artifacts must stay aligned with the current PostgreSQL-backed API contract. Regenerate them after route, request, or response changes.
