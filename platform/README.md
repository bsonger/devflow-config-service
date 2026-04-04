# Platform Notes

This repository only owns the `devflow-config-service` boundary.

Runtime shape:

- `cmd/main.go` uses shared bootstrap from `../devflow-service-common`
- `pkg/router/` exposes only configuration routes
- `pkg/api/configuration.go` is the only HTTP handler surface
- `pkg/service/configuration.go` is the only service entrypoint

Shared infra:

- pagination and response helpers come from `devflow-service-common/httpx`
- middleware and telemetry helpers come from `devflow-service-common/routercore` and `devflow-service-common/observability`

Operational rules:

- outbound service or external calls must emit `metrics + trace + structured log`
- `Planner -> Generator -> Evaluator` is the default harness
- when delegation is supported, sub-agents must be spawned
