# API Spec

## Purpose

`devflow-config-service` only exposes public HTTP APIs for `Configuration`.

## Endpoint Groups

### `Configuration`
- `GET /api/v1/configurations`
- `POST /api/v1/configurations`
- `GET /api/v1/configurations/:id`
- `PUT /api/v1/configurations/:id`
- `DELETE /api/v1/configurations/:id`

## Request Rules

- list endpoints support the common pagination parameters used in this repo
- create/update handlers currently bind the repo-owned `Configuration` payload directly

## Response Rules

- create endpoints return the common create-response shape
- list endpoints return pagination headers
- success payloads must stay aligned with Swagger

## Error Rules

- invalid ObjectID -> `400`
- resource not found -> `404`
- storage or uncategorized internal error -> `500`

## Boundary Note

For repo scope and non-goals, see `docs/architecture.md`.
