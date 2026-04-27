param(
    [string]$Configuration = "Release"
)

$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $PSScriptRoot
$Dist = Join-Path $Root "dist"

& (Join-Path $PSScriptRoot "build.ps1") -Configuration $Configuration

$PluginOut = Join-Path $Root "src/GrasshopperMcp.Plugin/bin/$Configuration"
$Gha = Join-Path $PluginOut "GrasshopperMcp.Plugin.gha"

if (-not (Test-Path $Gha)) {
    throw "Expected plugin artifact not found: $Gha"
}

Copy-Item $Gha (Join-Path $Dist "GrasshopperMcp.Plugin.gha") -Force
Write-Host "Package staged in $Dist"

