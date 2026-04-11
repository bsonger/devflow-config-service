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

- app-config routes: `pkg/router/app_config.go`
- app-config handler: `pkg/api/app_config.go`
- app-config logic: `pkg/app/app_config.go`
- workload-config routes: `pkg/router/workload_config.go`
- workload-config handler: `pkg/api/workload_config.go`
- workload-config logic: `pkg/app/workload_config.go`
- runtime config: `pkg/infra/config/config.go`

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
