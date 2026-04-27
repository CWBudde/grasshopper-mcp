#!/usr/bin/env sh
set -eu

ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
DIST="$ROOT/dist/windows-amd64"

mkdir -p "$DIST"

: "${GOCACHE:=/tmp/grasshopper-mcp-gocache}"
: "${GOMODCACHE:=/tmp/grasshopper-mcp-gomodcache}"
export GOCACHE GOMODCACHE

if [ "${RUN_TESTS:-0}" = "1" ]; then
	go test ./...
fi

GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o "$DIST/grasshopper-mcp.exe" ./cmd/grasshopper-mcp
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o "$DIST/grasshopper-mcp-fake-adapter.exe" ./cmd/grasshopper-mcp-fake-adapter

printf 'Windows Go binaries written to %s\n' "$DIST"
