$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $PSScriptRoot

go test ./...
go run (Join-Path $Root "cmd/grasshopper-mcp") --help

