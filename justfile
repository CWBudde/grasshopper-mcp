# grasshopper-mcp development tasks

set shell := ["bash", "-uc"]

plugin_project := "src/GrasshopperMcp.Plugin/GrasshopperMcp.Plugin.csproj"
config := env("CONFIG", "Debug")
dist := "dist"

# List available recipes
default:
    @just --list

# Format sources
fmt:
    treefmt --no-cache --allow-missing-formatter
    @just --fmt --unstable --check >/dev/null 2>&1 || just --fmt --unstable

# Check formatting
check-formatted:
    treefmt --no-cache --allow-missing-formatter --fail-on-change
    just --fmt --unstable --check

# Format the C# plugin project with dotnet format
fmt-csharp:
    dotnet format {{ plugin_project }}

# Check C# plugin formatting without keeping changes
check-csharp-formatted:
    dotnet format {{ plugin_project }} --verify-no-changes

# Run Go tests
test:
    go test ./...

# Run Go tests with verbose output
test-verbose:
    go test -v ./...

# Run Go vet
vet:
    go vet ./...

# Run shellcheck on shell scripts
lint-shell:
    shellcheck build/*.sh

# Ensure go.mod/go.sum are tidy
check-tidy:
    go mod tidy
    @test -z "$$(git status --porcelain -- go.mod go.sum)" || (git status --short -- go.mod go.sum && exit 1)

# Build the Go MCP server
build-go:
    mkdir -p {{ dist }}
    go build -o {{ dist }}/grasshopper-mcp ./cmd/grasshopper-mcp

# Build the Grasshopper plugin; requires RHINO_SYSTEM_DIR
build-plugin:
    dotnet build {{ plugin_project }} -c {{ config }}

# Build Go artifacts and the Grasshopper plugin; requires RHINO_SYSTEM_DIR
build: build-go build-plugin

# Stage the Windows package via PowerShell
package config="Release":
    pwsh -NoProfile -File build/package.ps1 -Configuration "{{ config }}"

# Run Go-only checks that work without Rhino
check: check-formatted vet test check-tidy

# Run the MCP server over stdio
serve:
    go run ./cmd/grasshopper-mcp serve

# Run the CLI with arbitrary arguments
run *args:
    go run ./cmd/grasshopper-mcp {{ args }}

# Remove local build artifacts
clean:
    rm -rf {{ dist }}
    rm -f coverage.out coverage.html

# Apply automatic fixes
fix:
    just fmt
