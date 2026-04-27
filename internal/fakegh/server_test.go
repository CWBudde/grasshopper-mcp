package fakegh

import (
	"context"
	"testing"
	"time"

	"github.com/cwbudde/grasshopper-mcp/internal/ghclient"
)

func TestFakeAdapterEndToEnd(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server := New("127.0.0.1:0")
	errc := make(chan error, 1)
	go func() {
		errc <- server.ListenAndServe(ctx)
	}()
	t.Cleanup(func() {
		cancel()
		select {
		case err := <-errc:
			if err != nil {
				t.Fatalf("fake adapter returned error: %v", err)
			}
		case <-time.After(time.Second):
			t.Fatalf("fake adapter did not stop")
		}
	})

	client := ghclient.New(waitForAddress(t, server))
	health, err := client.Health(context.Background())
	if err != nil {
		t.Fatalf("health: %v", err)
	}
	if !health.ActiveDocument || !health.GrasshopperLoaded {
		t.Fatalf("unexpected health: %+v", health)
	}

	component, err := client.AddComponent(context.Background(), ghclient.AddComponentParams{Name: "Addition"})
	if err != nil {
		t.Fatalf("add component: %v", err)
	}

	if _, err := client.SetInput(context.Background(), ghclient.SetInputParams{
		Target: ghclient.ParameterRef{ComponentID: component.ComponentID, Parameter: "A"},
		Value:  2.0,
	}); err != nil {
		t.Fatalf("set input: %v", err)
	}

	output, err := client.GetOutput(context.Background(), ghclient.GetOutputParams{
		Source: ghclient.ParameterRef{ComponentID: component.ComponentID, Parameter: "R"},
	})
	if err != nil {
		t.Fatalf("get output: %v", err)
	}
	if output.Type != "number" {
		t.Fatalf("output = %+v, want number", output)
	}
}

func waitForAddress(t *testing.T, server *Server) string {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if server.Address() != "127.0.0.1:0" {
			return server.Address()
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("fake adapter did not start listening")
	return ""
}
