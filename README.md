# grasshopper-mcp

Go-first recreation of `alfredatnycu/grasshopper-mcp`.

The Go executable will provide the MCP server and most application behavior. A
small C# `.gha` adapter is still required because Grasshopper loads .NET plugin
assemblies and exposes its document/component APIs through the RhinoCommon and
Grasshopper SDKs.

## Current Status

Implemented so far:

- Go module and CLI.
- Newline-delimited JSON Grasshopper client in Go.
- Minimal dependency-free stdio MCP server in Go.
- MCP tools for health, document info, component listing, solver execution, and
  the graph mutation command surface.
- Fake Grasshopper adapter for local Linux/macOS development without Rhino.
- Protocol contract fixtures for representative adapter requests/responses.
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

Cross-compile the Go executables for Windows from Linux:

```bash
sh ./build/cross-build-windows.sh
```

To run Go tests as part of that script:

```bash
RUN_TESTS=1 sh ./build/cross-build-windows.sh
```

This writes:

```text
dist/windows-amd64/grasshopper-mcp.exe
dist/windows-amd64/grasshopper-mcp-fake-adapter.exe
```

The C# `.gha` can only be built when `dotnet` can resolve RhinoCommon and
Grasshopper assemblies through `RHINO_SYSTEM_DIR`. In practice, that means
building it on a Windows Rhino dev machine, or copying the Rhino reference DLLs
into an equivalent path and setting `RHINO_SYSTEM_DIR` before `dotnet build`.

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
go run ./cmd/grasshopper-mcp add-component '{"name":"Addition","x":10,"y":20}'
go run ./cmd/grasshopper-mcp set-input '{"target":{"componentId":"component-1","parameter":"A"},"value":5}'
go run ./cmd/grasshopper-mcp connect '{"source":{"componentId":"a","parameter":"R"},"target":{"componentId":"b","parameter":"A"}}'
go run ./cmd/grasshopper-mcp get-output '{"source":{"componentId":"component-1","parameter":"R"}}'
```

The adapter address defaults to `127.0.0.1:47820` and can be changed with
`GRASSHOPPER_MCP_ADDR`.

Run a fake adapter locally without Rhino:

```bash
go run ./cmd/grasshopper-mcp-fake-adapter
```

Then in another shell:

```bash
go run ./cmd/grasshopper-mcp health
go run ./cmd/grasshopper-mcp add-component '{"name":"Addition","x":10,"y":20}'
go run ./cmd/grasshopper-mcp get-output '{"source":{"componentId":"component-1","parameter":"R"}}'
```

The real C# adapter uses port `47820` by default. On Windows it can be changed
with `GRASSHOPPER_MCP_PORT`; the Go client side can be pointed at that address
with `GRASSHOPPER_MCP_ADDR`.

Protocol fixtures live in `internal/ghclient/testdata/protocol`.

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
