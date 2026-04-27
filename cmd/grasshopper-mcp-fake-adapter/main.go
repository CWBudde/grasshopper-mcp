package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/cwbudde/grasshopper-mcp/internal/fakegh"
	"github.com/cwbudde/grasshopper-mcp/internal/ghclient"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	address := os.Getenv("GRASSHOPPER_MCP_ADDR")
	if address == "" {
		address = ghclient.DefaultAddress
	}

	server := fakegh.New(address)
	fmt.Fprintf(os.Stderr, "fake Grasshopper adapter listening on %s\n", address)
	if err := server.ListenAndServe(ctx); err != nil && !errors.Is(err, net.ErrClosed) {
		fmt.Fprintf(os.Stderr, "fake adapter: %v\n", err)
		os.Exit(1)
	}
}
