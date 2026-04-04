# Config Service Platform Notes

## Purpose

This file is the repo-local runtime note for `devflow-config-service`.
For public API shape, ownership, and resource details, prefer:
- `../README.md`
- `../docs/`
- `../docs/resources/`

## Runtime entrypoints

- process entry: `cmd/main.go`
- shared bootstrap: `../devflow-service-common/bootstrap`
- router root: `pkg/router/router.go`

## Main local code paths

- configuration routes: `pkg/router/configuration.go`
- configuration handler: `pkg/api/configuration.go`
- configuration logic: `pkg/service/configuration.go`
- runtime config: `pkg/config/config.go`

## Platform dependencies

- shared response / pagination: `devflow-service-common/httpx`
- shared middleware: `devflow-service-common/routercore`
- shared observability: `devflow-service-common/observability`

## Service identity

- OTel `service.name`: `config-service`
- typical ports:
  - `CONFIG_SERVICE_PORT`
  - `CONFIG_SERVICE_METRICS_PORT`
  - `CONFIG_SERVICE_PPROF_PORT`
