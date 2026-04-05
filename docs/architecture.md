# Architecture

## Purpose

`devflow-config-service` is the metadata owner for:

- `Configuration`
- `ConfigurationRevision`

It provides configuration identity, immutable configuration revisions, environment-variable ownership, and revision lookup for release flows.

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
  -> configuration / revision service
  -> persistence store
  -> HTTP response
```

The target relational model is:

- `Configuration` = mutable identity + latest revision pointer
- `ConfigurationRevision` = immutable content snapshot

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
  - configuration / revision handlers
- `pkg/service`
  - configuration identity behavior
  - revision creation / lookup behavior
- `pkg/store`
  - repo-owned configuration persistence
- `pkg/model`
  - `Configuration`, `ConfigurationRevision`

## External Dependencies

- `Gin`
- PostgreSQL persistence
- `devflow-service-common`

## Non-Goals

- `Project`
- `Application`
- `Manifest`
- `Release`
- `Intent`
- verify ingress / writeback
- service-exposure ownership
