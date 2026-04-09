# DevFlow Config Service

`devflow-config-service` is the backend owner for `Configuration`.

## Backend Role

- own `Configuration`
- provide configuration metadata and content lookup
- act as a release-input source for other services and the future platform
- read normalized `applications/devflow-platform/services/<service>/` config-repo layouts plus environment overlays

## Local Run

- `go run ./cmd`
- `go build ./cmd/main.go`
- `go test ./...`
- Swagger UI: `/swagger/index.html`

## Key Docs

- `docs/architecture.md`
- `docs/api-spec.md`
- `docs/constraints.md`
- `docs/resources/README.md`
- `docs/generated/swagger/swagger.yaml`
