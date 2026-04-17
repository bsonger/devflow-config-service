# Constraints

## Ownership

- `AppConfig` and `WorkloadConfig` belong exclusively to config-service

## Hard constraints

- do not add `Project`, `Application`, `Manifest`, `Release`, `Intent`, or `Verify` as public APIs in this repo
- do not move release or verify runtime logic back into config-service
- do not put high-cardinality business identifiers into metric labels

## Data rules

- deletion semantics follow the existing model definitions
- every write operation must update `updated_at`
- `AppConfig.source_path` is system-derived from `name`; clients must not treat it as user-editable data
- `sync-from-repo` must stay aligned with `devflow-config-repo` service layout under `applications/devflow-platform/services/<service>`
