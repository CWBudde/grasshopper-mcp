# Implementation Plan: Go-Based grasshopper-mcp

## Goal

Build a working `grasshopper-mcp` where Go owns the MCP server, command routing,
Grasshopper protocol client, validation, and most testable behavior. A thin C#
Grasshopper plugin remains necessary because Grasshopper loads .NET assemblies,
not native Go plugins.

The first usable target is:

- Rhino 8 / Grasshopper on Windows.
- A C# `.gha` plugin that runs inside Grasshopper and exposes a small local JSON
  control API.
- A Go MCP server that talks stdio MCP to clients and talks TCP/WebSocket JSON to
  the Grasshopper plugin.
- A minimal end-to-end tool set: health check, document status, component
  inventory, add component, connect parameters, set input values, run solver,
  and read results.

## Architecture Decision

Do not attempt a pure Go `.gha`. The Grasshopper side must be C#/.NET because
components and plugin lifecycle hooks are discovered through Grasshopper's .NET
APIs.

Use this split:

| Layer | Language | Responsibility |
|---|---|---|
| MCP server | Go | MCP stdio transport, tool schemas, request validation, logging, error mapping |
| Grasshopper client | Go | Typed client for the local Grasshopper control protocol |
| Domain model | Go | Component ids, parameter refs, graph operations, command/result DTOs |
| Grasshopper adapter | C# | Rhino/Grasshopper lifecycle, document access, component creation, solver calls |
| Packaging | PowerShell/Yak | Build `.gha`, build Go executable, produce installable layout |

This differs slightly from the pure in-process component MVP discussed in
`goal.md`: because this repository is named `grasshopper-mcp`, the first working
product should prove the MCP bridge path. A later phase can add custom
Grasshopper components backed by a Go native DLL if that becomes a product
requirement.

## Proposed Repository Layout

```text
.
|-- README.md
|-- PLAN.md
|-- go.mod
|-- cmd/
|   `-- grasshopper-mcp/
|       `-- main.go
|-- internal/
|   |-- ghclient/
|   |   |-- client.go
|   |   |-- protocol.go
|   |   `-- errors.go
|   |-- mcp/
|   |   |-- server.go
|   |   `-- tools.go
|   `-- model/
|       `-- graph.go
|-- src/
|   `-- GrasshopperMcp.Plugin/
|       |-- GrasshopperMcp.Plugin.csproj
|       |-- AssemblyInfo.cs
|       |-- PriorityLoader.cs
|       |-- Server/
|       |   |-- LocalServer.cs
|       |   |-- Protocol.cs
|       |   `-- CommandRouter.cs
|       `-- Grasshopper/
|           |-- DocumentService.cs
|           |-- ComponentCatalog.cs
|           `-- GraphMutationService.cs
|-- tests/
|   |-- ghclient/
|   |-- mcp/
|   `-- dotnet/
|-- build/
|   |-- build.ps1
|   |-- package.ps1
|   `-- smoke.ps1
`-- dist/
```

## Milestones

Current implementation status on Ubuntu:

- Milestone 1 is partially done. The portable scaffold is implemented; the C#
  plugin still requires a Windows Rhino machine for real compilation.
- Milestone 2 is partially done. Protocol types and loopback server source are
  implemented. The server is bound to `127.0.0.1` and supports `health`,
  `document_info`, `list_components`, and `run_solver` responses, but
  document/component data is still placeholder data until Rhino services are
  wired on Windows.
- Milestone 3 is done for the phase-2 commands. The Go client has timeout,
  malformed response, protocol error, and fake-server tests.
- Milestone 4 is partially done. The minimal stdio MCP server is implemented
  for the phase-2 tools and covered by Go tests. MCP inspector compatibility
  still needs a live external check.

### 1. Project Skeleton

Status: Partially done.

Deliverables:

- [x] Initialize Go module.
- [x] Add C# Grasshopper plugin project targeting Rhino 8.
- [x] Add build script placeholders for Go, .NET, and packaging.
- [x] Add README with Windows prerequisites and manual install path.
- [x] Add `.gitignore` and SDK selection files.

Acceptance criteria:

- [x] `go test ./...` runs on Ubuntu.
- [x] `go run ./cmd/grasshopper-mcp --help` prints available commands.
- [x] `dotnet build` fails clearly on Ubuntu when `RHINO_SYSTEM_DIR` is unset.
- [ ] `dotnet build` compiles the plugin project on a Windows Rhino dev
  machine.
- [ ] The generated `.gha` can be copied into Grasshopper's library folder.

### 2. Local Grasshopper Control Protocol

Status: Partially done.

Deliverables:

- [x] Define one JSON request/response envelope shared by Go and C#.
- [x] Implement C# local server source bound to loopback only.
- [x] Implement `health` command route.
- [x] Implement `document_info` command route.
- [x] Implement `list_components` command route.
- [x] Implement `run_solver` command route.
- [x] Add structured error responses with stable error codes.
- [ ] Replace placeholder C# document data with live Grasshopper document data.
- [ ] Replace placeholder C# component data with live component catalog data.

Acceptance criteria:

- [x] Protocol responses are newline-delimited JSON.
- [x] Unknown methods return a structured `unknown_method` error.
- [x] The server binds to `IPAddress.Loopback`, not a public interface.
- [ ] When Rhino/Grasshopper is open, a local client can call `health` and
  receive plugin version plus active-document state.
- [ ] Live Rhino smoke test confirms errors are JSON responses, not unhandled
  exceptions.

### 3. Go Grasshopper Client

Status: Done for phase-2 commands.

Deliverables:

- [x] Implement `internal/ghclient` with typed methods for all phase-2 commands.
- [x] Add connection timeout and request timeout.
- [x] Add clear "Grasshopper adapter is unavailable" errors.
- [x] Add unit tests using an in-process fake server.
- [x] Add CLI debug commands for the phase-2 commands.

Acceptance criteria:

- [x] Go tests cover success responses.
- [x] Go tests cover structured protocol errors.
- [x] Go tests cover malformed responses.
- [x] Go tests cover request timeout behavior.
- [x] CLI debug subcommands exist for `health`, `document-info`,
  `list-components`, and `run-solver`.
- [x] Adapter address can be configured with `GRASSHOPPER_MCP_ADDR`.

### 4. Minimal MCP Server

Status: Partially done.

Deliverables:

- [x] Implement Go stdio MCP server.
- [x] Expose `grasshopper_health`.
- [x] Expose `grasshopper_document_info`.
- [x] Expose `grasshopper_list_components`.
- [x] Expose `grasshopper_run_solver`.
- [x] Map Grasshopper protocol errors to MCP tool errors with useful messages.
- [ ] Run against MCP inspector or a real MCP client.

Acceptance criteria:

- [x] `tools/list` returns the initial four tools in Go tests.
- [x] `tools/call` delegates to the Go Grasshopper client in Go tests.
- [x] Tool calls fail cleanly when the adapter is unavailable.
- [ ] MCP inspector or a compatible MCP client can list tools.
- [ ] `grasshopper_health` works with Grasshopper open.

### 5. Graph Mutation MVP

Status: Not started.

Deliverables:

- [ ] Add C# services for safe graph mutation on the Grasshopper/Rhino UI
  thread.
- [ ] Implement `add_component` protocol command.
- [ ] Implement `set_input` protocol command.
- [ ] Implement `connect` protocol command.
- [ ] Implement `get_output` protocol command.
- [ ] Add matching Go client methods.
- [ ] Add matching MCP tools.
- [ ] Start with common built-in component lookup by name/category/nickname.

Acceptance criteria:

- [ ] From an MCP client, create a small graph in an empty Grasshopper document.
- [ ] Set numeric inputs on created components.
- [ ] Connect at least two component parameters.
- [ ] Run the solver after graph mutation.
- [ ] Read an expected output value through `get_output`.
- [ ] Invalid component names produce stable, actionable errors.
- [ ] Invalid parameter references produce stable, actionable errors.
- [ ] Graph mutation is dispatched safely on the Grasshopper/Rhino UI thread.

### 6. Packaging and Manual Distribution

Status: Not started beyond placeholders.

Deliverables:

- [ ] Build script produces `grasshopper-mcp.exe`.
- [ ] Build script produces `GrasshopperMcp.Plugin.gha`.
- [ ] Build script stages any required `.dll` dependencies.
- [ ] Build script stages an example MCP client config.
- [x] Manual install instructions for the `.gha` are documented.
- [ ] Optional Yak package layout once manual installation is stable.

Acceptance criteria:

- [ ] `build/build.ps1` completes on a Windows Rhino dev machine.
- [ ] `build/package.ps1` stages a usable `dist` directory.
- [ ] A clean Windows machine with Rhino 8, Go-built executable, and the plugin
  can run the full health-to-result smoke test.
- [ ] No absolute developer-machine paths are required at runtime.
- [ ] The staged MCP client config launches the Go server correctly.

### 7. Hardening

Status: Not started.

Deliverables:

- [ ] Add request correlation ids and structured logs on both sides.
- [ ] Add version negotiation between Go server and C# plugin.
- [ ] Add configurable port with safe default and discovery fallback.
- [ ] Add component catalog cache invalidation.
- [ ] Add Rhino compatibility smoke notes for Rhino 8 updates.
- [ ] Add long-running solver timeout behavior.

Acceptance criteria:

- [ ] Version mismatch reports a clear message.
- [ ] Logs are enough to diagnose "MCP client cannot reach Grasshopper" without
  a debugger.
- [ ] Long-running solver calls time out at the Go/MCP boundary without leaving
  the server wedged.
- [ ] Port conflicts report a clear startup error.
- [ ] Component catalog changes are reflected without restarting the Go server.

## Protocol Shape

Start with newline-delimited JSON over loopback TCP. It is easy to debug, easy to
fake in tests, and sufficient for local MCP calls.

Request:

```json
{
  "id": "01HT...",
  "method": "document_info",
  "params": {}
}
```

Success response:

```json
{
  "id": "01HT...",
  "ok": true,
  "result": {
    "documentName": "Untitled",
    "objectCount": 12
  }
}
```

Error response:

```json
{
  "id": "01HT...",
  "ok": false,
  "error": {
    "code": "component_not_found",
    "message": "No Grasshopper component matched 'Addition'"
  }
}
```

## MVP Tool List

| MCP tool | Grasshopper command | Purpose |
|---|---|---|
| `grasshopper_health` | `health` | Verify plugin connection and versions |
| `grasshopper_document_info` | `document_info` | Inspect active document state |
| `grasshopper_list_components` | `list_components` | Discover available components |
| `grasshopper_add_component` | `add_component` | Place a component on the canvas |
| `grasshopper_set_input` | `set_input` | Assign simple values to inputs |
| `grasshopper_connect` | `connect` | Connect output parameter to input parameter |
| `grasshopper_run_solver` | `run_solver` | Recompute the graph |
| `grasshopper_get_output` | `get_output` | Read computed output values |

Keep the first value types deliberately small: numbers, booleans, strings, and
flat lists. Geometry serialization can be added after the command loop is
stable.

## Testing Strategy

Go:

- Unit-test all DTO validation and MCP tool argument validation.
- Test `ghclient` against fake TCP servers.
- Add CLI smoke tests that do not require Rhino.

C#:

- Unit-test protocol parsing and command routing outside Rhino where practical.
- Keep Grasshopper API calls behind services so command routing stays testable.
- Use Rhino/Grasshopper manual smoke tests for document mutation until a reliable
  automated Rhino harness is available.

End-to-end:

- Start Rhino and Grasshopper with the plugin installed.
- Start `grasshopper-mcp.exe`.
- Call health.
- Create a numeric graph.
- Solve.
- Read expected output.

## Main Risks

| Risk | Mitigation |
|---|---|
| Grasshopper API calls from the wrong thread | Centralize UI-thread dispatch in the C# adapter |
| Protocol drift between Go and C# | Keep a versioned protocol document and contract tests |
| Component lookup ambiguity | Return candidates and require stable ids where possible |
| Rhino version differences | Start Rhino 8-only, add compatibility checks later |
| Complex geometry serialization too early | Defer geometry; prove scalar/list graph operations first |
| MCP debugging is opaque | Provide direct Go CLI debug commands beside MCP tools |

## First Implementation Sprint

1. Create `go.mod` and a minimal `cmd/grasshopper-mcp`.
2. Add the C# plugin project with `GH_AssemblyInfo` and `GH_AssemblyPriority`.
3. Implement C# `health` server on loopback.
4. Implement Go `ghclient.Health`.
5. Add MCP `grasshopper_health`.
6. Verify with a real Grasshopper session.

This sprint is successful when an MCP client can call one Go tool and receive a
live response from the Grasshopper plugin running inside Rhino.
