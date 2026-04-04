# Devflow Config Service

`devflow-config-service` is the backend owner for `Configuration`.

## Backend Role

- own `Configuration`
- provide configuration metadata and content lookup
- act as a release-input source for other services and the future platform

## Backend Architecture

This repo uses a **layered metadata-service backend**:

```text
cmd
 -> config
 -> router
 -> api
 -> service
 -> store
 -> model
```

### Package responsibilities

- `cmd/`: service startup
- `pkg/config`: config loading and runtime init
- `pkg/router`: Gin router and middleware wiring
- `pkg/api`: HTTP handlers and status mapping
- `pkg/service`: configuration rules and resource behavior
- `pkg/store`: Mongo access
- `pkg/model`: `Configuration` model

## Non-Goals

- no `Project` ownership
- no `Application` ownership
- no `Manifest` ownership
- no `Release` ownership
- no `Intent` ownership
- no verify ingress

## Key Docs

- `docs/architecture.md`
- `docs/api-spec.md`
- `docs/constraints.md`
- `docs/resources/README.md`

## Local Run

- `go run ./cmd`
- `go build ./cmd/main.go`
- `go test ./...`
