package mcp

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"net"
	"strings"
	"testing"

	"github.com/cwbudde/grasshopper-mcp/internal/ghclient"
)

func TestToolsList(t *testing.T) {
	input := strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"tools/list"}` + "\n")
	var output bytes.Buffer
	server := NewServer(ghclient.New("127.0.0.1:1"), WithIO(input, &output))

	if err := server.Run(context.Background()); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	var response rpcResponse
	if err := json.Unmarshal(output.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Error != nil {
		t.Fatalf("unexpected rpc error: %+v", response.Error)
	}
	result := response.Result.(map[string]any)
	toolList := result["tools"].([]any)
	if len(toolList) != 4 {
		t.Fatalf("tool count = %d, want 4", len(toolList))
	}
}

func TestHealthToolUsesGrasshopperClient(t *testing.T) {
	address, stop := startFakeGrasshopper(t)
	defer stop()

	input := strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"grasshopper_health","arguments":{}}}` + "\n")
	var output bytes.Buffer
	server := NewServer(ghclient.New(address), WithIO(input, &output))

	if err := server.Run(context.Background()); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if !strings.Contains(output.String(), `\"version\": \"test\"`) {
		t.Fatalf("output did not contain health result: %s", output.String())
	}
}

func TestHealthToolReturnsToolErrorWhenGrasshopperIsUnavailable(t *testing.T) {
	input := strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"grasshopper_health","arguments":{}}}` + "\n")
	var output bytes.Buffer
	server := NewServer(ghclient.New("127.0.0.1:1"), WithIO(input, &output))

	if err := server.Run(context.Background()); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if !strings.Contains(output.String(), `"isError":true`) {
		t.Fatalf("output did not contain tool error: %s", output.String())
	}
}

func startFakeGrasshopper(t *testing.T) (string, func()) {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	done := make(chan struct{})
	go func() {
		defer close(done)
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		line, err := bufio.NewReader(conn).ReadBytes('\n')
		if err != nil {
			t.Errorf("read request: %v", err)
			return
		}
		var request ghclient.Request
		if err := json.Unmarshal(line, &request); err != nil {
			t.Errorf("decode request: %v", err)
			return
		}
		response := map[string]any{
			"id": request.ID,
			"ok": true,
			"result": map[string]any{
				"version":           "test",
				"activeDocument":    true,
				"grasshopperLoaded": true,
			},
		}
		if err := json.NewEncoder(conn).Encode(response); err != nil {
			t.Errorf("write response: %v", err)
		}
	}()
	stop := func() {
		_ = listener.Close()
		<-done
	}
	return listener.Addr().String(), stop
}
