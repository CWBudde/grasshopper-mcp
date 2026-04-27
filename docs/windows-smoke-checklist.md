# Windows Rhino Smoke Checklist

Use this checklist after copying the repository to a Windows machine with Rhino
8 installed.

## Build

If the Go executables were cross-compiled on Linux, copy these files to the
Windows machine:

```text
dist\windows-amd64\grasshopper-mcp.exe
dist\windows-amd64\grasshopper-mcp-fake-adapter.exe
```

The `.gha` still needs Rhino/Grasshopper references:

```powershell
$env:RHINO_SYSTEM_DIR = "C:\Program Files\Rhino 8\System"
go test ./...
.\build\build.ps1
```

Expected artifacts:

- `dist\grasshopper-mcp.exe`
- `src\GrasshopperMcp.Plugin\bin\Release\GrasshopperMcp.Plugin.gha`

## Install Plugin

Copy the `.gha` to:

```text
%APPDATA%\Grasshopper\Libraries\
```

Restart Rhino, run `Grasshopper`, and check the Rhino command history for:

```text
Grasshopper MCP adapter 0.1.0 listening on 127.0.0.1:47820.
```

## Adapter Smoke

With Rhino and Grasshopper open:

```powershell
go run .\cmd\grasshopper-mcp health
go run .\cmd\grasshopper-mcp document-info
go run .\cmd\grasshopper-mcp list-components
go run .\cmd\grasshopper-mcp run-solver
```

Expected at this stage:

- `health` returns version `0.1.0`.
- `document-info` returns active document state.
- `list-components` may still be placeholder data until the catalog is wired.
- graph mutation commands may return `graph_mutation_not_implemented` until
  `GraphMutationService` is connected to live Grasshopper APIs.

## MCP Smoke

Start the MCP server through the target MCP client with command:

```text
dist\grasshopper-mcp.exe serve
```

Expected:

- Tool list contains all eight MVP tools.
- `grasshopper_health` returns a response from the live adapter.
- If Rhino is closed, tool calls fail with a clear adapter-unavailable error.
