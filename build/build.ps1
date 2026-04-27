param(
    [string]$Configuration = "Release"
)

$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $PSScriptRoot
$Dist = Join-Path $Root "dist"

New-Item -ItemType Directory -Force -Path $Dist | Out-Null

go test ./...
go build -o (Join-Path $Dist "grasshopper-mcp.exe") ./cmd/grasshopper-mcp

if (-not $env:RHINO_SYSTEM_DIR) {
    throw "Set RHINO_SYSTEM_DIR to Rhino 8's System directory before building the Grasshopper plugin."
}

dotnet build (Join-Path $Root "src/GrasshopperMcp.Plugin/GrasshopperMcp.Plugin.csproj") -c $Configuration

