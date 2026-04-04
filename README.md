# DevFlow Config Service

`devflow-config-service` is the backend owner for `Configuration`.

## Backend Role

- own `Configuration`
- provide configuration metadata and content lookup
- act as a release-input source for other services and the future platform

## Local Run

- `go run ./cmd`
- `go build ./cmd/main.go`
- `go test ./...`

## Key Docs

- `docs/architecture.md`
- `docs/api-spec.md`
- `docs/constraints.md`
- `docs/resources/README.md`
