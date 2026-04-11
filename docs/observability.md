# Observability

## Purpose

`devflow-config-service` emits the shared backend telemetry baseline plus config-sync and workload-template context.

## Logs

Required structured fields:
- `resource`
- `resource_id`
- `application_id`
- `app_config_id`
- `workload_config_id`
- `result`
- `error_code`

## Metrics

- use shared `devflow_http_*` ingress metrics
- if outbound repo sync or other dependencies are added, also emit `devflow_dependency_*` metrics
- forbid high-cardinality labels such as commit hashes, raw paths, or full repo URLs

## Tracing

- every business HTTP request should create a server span
- any future outbound Git or service calls must emit client spans with propagated trace context
- attach config resource identifiers as span attributes, not metric labels

## Health and readiness

- expose `/healthz`, `/readyz`, and `/metrics`
- keep diagnostics endpoints out of business counters/histograms

## Failure modes

Watch for:
- repo sync failures
- render/configmap snapshot failures
- workload template validation errors

## Dashboards and runbooks

Use the shared backend dashboard/runbook set until repo-specific views exist.
