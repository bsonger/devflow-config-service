#!/usr/bin/env bash
set -euo pipefail

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${DIR}"

export GOROOT="$(go env GOROOT)"
export GOPATH="$(go env GOPATH)"
export PATH="${GOROOT}/bin:${PATH}"

go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g cmd/main.go --parseDependency -o docs/generated/swagger
