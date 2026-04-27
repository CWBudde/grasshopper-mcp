package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/cwbudde/grasshopper-mcp/internal/ghclient"
	"github.com/cwbudde/grasshopper-mcp/internal/mcp"
)

func main() {
	if err := run(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "grasshopper-mcp: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		printUsage()
		return nil
	}

	client := ghclient.New(adapterAddress())

	switch args[0] {
	case "serve":
		server := mcp.NewServer(client, mcp.WithIO(os.Stdin, os.Stdout))
		return server.Run(ctx)
	case "health":
		health, err := client.Health(ctx)
		if err != nil {
			return err
		}
		return printJSON(health)
	case "document-info":
		result, err := client.DocumentInfo(ctx)
		if err != nil {
			return err
		}
		return printJSON(result)
	case "list-components":
		result, err := client.ListComponents(ctx)
		if err != nil {
			return err
		}
		return printJSON(result)
	case "run-solver":
		result, err := client.RunSolver(ctx)
		if err != nil {
			return err
		}
		return printJSON(result)
	case "add-component":
		var params ghclient.AddComponentParams
		if err := decodeJSONArg(args, &params); err != nil {
			return err
		}
		result, err := client.AddComponent(ctx, params)
		if err != nil {
			return err
		}
		return printJSON(result)
	case "set-input":
		var params ghclient.SetInputParams
		if err := decodeJSONArg(args, &params); err != nil {
			return err
		}
		result, err := client.SetInput(ctx, params)
		if err != nil {
			return err
		}
		return printJSON(result)
	case "connect":
		var params ghclient.ConnectParams
		if err := decodeJSONArg(args, &params); err != nil {
			return err
		}
		result, err := client.Connect(ctx, params)
		if err != nil {
			return err
		}
		return printJSON(result)
	case "get-output":
		var params ghclient.GetOutputParams
		if err := decodeJSONArg(args, &params); err != nil {
			return err
		}
		result, err := client.GetOutput(ctx, params)
		if err != nil {
			return err
		}
		return printJSON(result)
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func adapterAddress() string {
	if address := os.Getenv("GRASSHOPPER_MCP_ADDR"); address != "" {
		return address
	}
	return ghclient.DefaultAddress
}

func printJSON(value any) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}

func decodeJSONArg(args []string, target any) error {
	if len(args) < 2 {
		return fmt.Errorf("%s requires a JSON argument", args[0])
	}
	if err := json.Unmarshal([]byte(args[1]), target); err != nil {
		return fmt.Errorf("decode %s argument: %w", args[0], err)
	}
	return nil
}

func printUsage() {
	fmt.Println("grasshopper-mcp")
	fmt.Println("commands:")
	fmt.Println("  serve                         start stdio MCP server")
	fmt.Println("  health                        call Grasshopper health command")
	fmt.Println("  document-info                 call Grasshopper document_info command")
	fmt.Println("  list-components               call Grasshopper list_components command")
	fmt.Println("  run-solver                    call Grasshopper run_solver command")
	fmt.Println("  add-component JSON            call Grasshopper add_component command")
	fmt.Println("  set-input JSON                call Grasshopper set_input command")
	fmt.Println("  connect JSON                  call Grasshopper connect command")
	fmt.Println("  get-output JSON               call Grasshopper get_output command")
	fmt.Println("")
	fmt.Println("environment:")
	fmt.Printf("  GRASSHOPPER_MCP_ADDR  adapter address (default %s)\n", ghclient.DefaultAddress)
}
