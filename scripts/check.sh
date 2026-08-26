#!/usr/bin/env bash
set -euo pipefail
gofmt -w $(find . -name '*.go')
go test ./...
