# grasshopper-mcp

Go-first recreation of `alfredatnycu/grasshopper-mcp`.

The Go executable will provide the MCP server and most application behavior. A
small C# `.gha` adapter is still required because Grasshopper loads .NET plugin
assemblies and exposes its document/component APIs through the RhinoCommon and
Grasshopper SDKs.

## Current Status

Implemented so far:

- Go module and placeholder CLI.
- Newline-delimited JSON Grasshopper client in Go.
- Minimal dependency-free stdio MCP server in Go.
- MCP tools for health, document info, component listing, and solver execution.
- Rhino 8 Grasshopper plugin project scaffold with a loopback TCP server source.
- Build/package script placeholders.

## Prerequisites

For Go-only development:

- Go 1.25 or newer.

For Grasshopper plugin builds on Windows:

- Rhino 8 with Grasshopper.
- .NET SDK compatible with Rhino 8 plugin development.
- `RHINO_SYSTEM_DIR` set to Rhino's `System` directory, for example:

```powershell
$env:RHINO_SYSTEM_DIR = "C:\Program Files\Rhino 8\System"
```

## Build

Go-only checks:

```bash
go test ./...
go run ./cmd/grasshopper-mcp
```

Run the MCP server over stdio:

```bash
go run ./cmd/grasshopper-mcp serve
```

Direct debug commands against a running Grasshopper adapter:

```bash
go run ./cmd/grasshopper-mcp health
go run ./cmd/grasshopper-mcp document-info
go run ./cmd/grasshopper-mcp list-components
go run ./cmd/grasshopper-mcp run-solver
```

The adapter address defaults to `127.0.0.1:47820` and can be changed with
`GRASSHOPPER_MCP_ADDR`.

Windows full build:

```powershell
.\build\build.ps1
```

Package staging:

```powershell
.\build\package.ps1
```

The plugin project copies its assembly to `GrasshopperMcp.Plugin.gha` after a
successful build.

## Manual Plugin Install

After building on Windows, copy `GrasshopperMcp.Plugin.gha` to:

```text
%APPDATA%\Grasshopper\Libraries\
```

Restart Rhino and Grasshopper after installing or replacing the `.gha`.
