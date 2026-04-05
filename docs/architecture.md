# Architecture

## Purpose

`devflow-config-service` is the metadata owner for:

- `Configuration`
- `ConfigurationRevision`

It provides configuration identity, immutable configuration revisions, environment-variable ownership, and revision lookup for release flows.
Configuration file content itself now comes from the centralized config repo; this service freezes repo snapshots into immutable revisions.

## Architecture Style

This repo uses a **layered metadata-service backend**:

```text
router -> api -> app -> infra/store
                \-> infra/config_repo
                \-> domain
```

## Request Flow

```text
Client
  -> router
  -> configuration handler
  -> configuration / revision app service
  -> config repo snapshot reader
  -> persistence store
  -> HTTP response
```

The target relational model is:

- `Configuration` = mutable identity + source path + latest revision pointer
- `ConfigurationRevision` = immutable repo-derived content snapshot

## Internal Package Layout

- `cmd/main.go`
  - process entrypoint only
- `pkg/infra/config`
  - config loading
  - runtime initialization
- `pkg/router`
  - route registration
  - middleware wiring
- `pkg/api`
  - configuration / revision handlers
- `pkg/app`
  - configuration identity behavior
  - explicit sync / revision freeze behavior
- `pkg/infra/store`
  - repo-owned configuration persistence
- `pkg/infra/config_repo`
  - centralized config repo snapshot loading
- `pkg/domain`
  - `Configuration`, `ConfigurationRevision`

## External Dependencies

- `Gin`
- PostgreSQL persistence
- centralized config repo filesystem layout
- `devflow-service-common`

## Non-Goals

- `Project`
- `Application`
- `Manifest`
- `Release`
- `Intent`
- verify ingress / writeback
- service-exposure ownership
