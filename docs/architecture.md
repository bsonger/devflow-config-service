# Architecture

## Purpose

`devflow-config-service` is the metadata owner for `Configuration`.
It provides configuration CRUD, configuration content storage, and configuration lookup for release flows.

## Architecture Style

This repo uses a **layered metadata-service backend**:

```text
router -> api -> service -> store
                    \-> model
```

## Request Flow

```text
Client
  -> router
  -> configuration handler
  -> configuration service
  -> Mongo store
  -> HTTP response
```

## Internal Package Layout

- `cmd/main.go`
  - process entrypoint only
- `pkg/config`
  - config loading
  - runtime initialization
- `pkg/router`
  - route registration
  - middleware wiring
- `pkg/api`
  - configuration handlers
- `pkg/service`
  - configuration behavior
- `pkg/store`
  - Mongo access
- `pkg/model`
  - `Configuration` model

## External Dependencies

- `Gin`
- `MongoDB`
- `devflow-service-common`

## Non-Goals

- `Project`
- `Application`
- `Manifest`
- `Release`
- `Intent`
- verify ingress / writeback
