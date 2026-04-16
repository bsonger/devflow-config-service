# DevFlow Config Service

`devflow-config-service` is the backend owner for `AppConfig` and `WorkloadConfig`.

## Role

- own `AppConfig` and `WorkloadConfig`
- provide app-level config identity plus explicit `sync-from-repo`
- act as a release-input source for other services and the platform
- read config files from the fixed repo `git@github.com:bsonger/devflow-config-repo.git` on branch `main`
- `sync-from-repo` pulls the latest `main` before freezing a revision

## Key Commands

- `go run ./cmd`
- `go build ./cmd/main.go`
- `go test ./...`
- Swagger UI: `/swagger/index.html`
- Staging Swagger UI: `/api/v1/config/swagger/index.html`

## Key Docs

- `docs/README.md`
- `scripts/README.md`
- `docs/architecture.md`
- `docs/constraints.md`
- `docs/observability.md`
- `docs/api-spec.md`
- `docs/resources/README.md`
- `docs/generated/swagger/swagger.yaml`
