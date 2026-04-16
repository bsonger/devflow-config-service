# Architecture

## Purpose

`devflow-config-service` is the metadata owner for:

- `AppConfig`
- `AppConfigRevision`
- `WorkloadConfig`

It provides app-config identity, immutable app-config revisions, and workload template ownership for release flows.
App config file content is read from the fixed repo `git@github.com:bsonger/devflow-config-repo.git` on branch `main`, using a derived path based on `application_id + environment_id`.

## Architecture style

This repo uses a **layered metadata-service backend**:

```text
router -> api -> app -> infra/store
                \-> infra/config_repo
                \-> domain
```

## Request flow

```text
Client
  -> router
  -> app-config / workload-config handler
  -> app-config / revision / workload-config app service
  -> config repo snapshot reader
  -> persistence store
  -> HTTP response
```

The target relational model is:

- `AppConfig` = mutable identity + derived source path + latest revision pointer
- `AppConfigRevision` = immutable repo-derived file snapshot
- `WorkloadConfig` = runtime template plus strategy type

## Internal package layout

- `cmd/main.go`
  - process entrypoint only
- `pkg/infra/config`
  - config loading
  - runtime initialization
- `pkg/router`
  - route registration
  - middleware wiring
- `pkg/api`
  - app-config / workload-config handlers
- `pkg/app`
  - app-config identity behavior
  - explicit sync / revision freeze behavior
  - workload-config behavior
- `pkg/infra/store`
  - repo-owned app-config / workload-config persistence
- `pkg/infra/config_repo`
  - centralized config repo snapshot loading
- `pkg/domain`
  - `AppConfig`, `AppConfigRevision`, `WorkloadConfig`

## External dependencies

- `Gin`
- PostgreSQL persistence
- fixed config repo checkout from `git@github.com:bsonger/devflow-config-repo.git` (`main`)
- `devflow-service-common`

## Swagger generation

- `scripts/regen-swagger.sh` reruns `swag init -g cmd/main.go --parseDependency -o docs/generated/swagger`.
- `scripts/build.sh` calls regen then builds the binary to `bin/`.
- `Dockerfile` executes `swag init -g cmd/main.go --parseDependency -o docs/generated/swagger` during the build stage.
- Export scripts rely on `docs/generated/swagger` being populated at build time and `scripts/export_service_repo.sh` copies that folder into split repos.

## Non-goals

- `Project`
- `Application`
- `Image`
- `Release`
- `Intent`
- verify ingress / writeback
- service / route ownership
