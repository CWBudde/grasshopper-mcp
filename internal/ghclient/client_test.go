package ghclient

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"net"
	"strings"
	"testing"
	"time"
)

func TestNewUsesDefaultAddress(t *testing.T) {
	client := New("")
	if client.Address() != DefaultAddress {
		t.Fatalf("address = %q, want %q", client.Address(), DefaultAddress)
	}
}

func TestNewUsesExplicitAddress(t *testing.T) {
	const address = "127.0.0.1:10000"
	client := New(address)
	if client.Address() != address {
		t.Fatalf("address = %q, want %q", client.Address(), address)
	}
}

func TestHealthSuccess(t *testing.T) {
	address, stop := startFakeServer(t, func(request Request) Response {
		if request.Method != "health" {
			t.Fatalf("method = %q, want health", request.Method)
		}
		return successResponse(t, request.ID, HealthResult{
			Version:           "0.1.0",
			ActiveDocument:    true,
			GrasshopperLoaded: true,
		})
	})
	defer stop()

	client := New(address)
	result, err := client.Health(context.Background())
	if err != nil {
		t.Fatalf("Health returned error: %v", err)
	}
	if result.Version != "0.1.0" || !result.ActiveDocument || !result.GrasshopperLoaded {
		t.Fatalf("unexpected health result: %+v", result)
	}
}

func TestDocumentInfoSuccess(t *testing.T) {
	address, stop := startFakeServer(t, func(request Request) Response {
		return successResponse(t, request.ID, DocumentInfoResult{
			DocumentName:      "Untitled",
			ObjectCount:       3,
			HasActiveDocument: true,
		})
	})
	defer stop()

	client := New(address)
	result, err := client.DocumentInfo(context.Background())
	if err != nil {
		t.Fatalf("DocumentInfo returned error: %v", err)
	}
	if result.DocumentName != "Untitled" || result.ObjectCount != 3 || !result.HasActiveDocument {
		t.Fatalf("unexpected document info: %+v", result)
	}
}

func TestProtocolError(t *testing.T) {
	address, stop := startFakeServer(t, func(request Request) Response {
		return Response{
			ID: request.ID,
			OK: false,
			Error: &ProtocolError{
				Code:    "no_document",
				Message: "No active Grasshopper document.",
			},
		}
	})
	defer stop()

	client := New(address)
	_, err := client.Health(context.Background())
	var protocolErr *ProtocolError
	if !errors.As(err, &protocolErr) {
		t.Fatalf("error = %T %v, want *ProtocolError", err, err)
	}
	if protocolErr.Code != "no_document" {
		t.Fatalf("code = %q, want no_document", protocolErr.Code)
	}
}

func TestMalformedResponse(t *testing.T) {
	address, stop := startRawServer(t, func(conn net.Conn) {
		_, _ = bufio.NewReader(conn).ReadBytes('\n')
		_, _ = conn.Write([]byte("{not-json}\n"))
	})
	defer stop()

	client := New(address)
	_, err := client.Health(context.Background())
	if err == nil || !strings.Contains(err.Error(), "decode grasshopper response") {
		t.Fatalf("error = %v, want decode error", err)
	}
}

func TestTimeout(t *testing.T) {
	address, stop := startRawServer(t, func(conn net.Conn) {
		_, _ = bufio.NewReader(conn).ReadBytes('\n')
		time.Sleep(200 * time.Millisecond)
	})
	defer stop()

	client := New(address, WithRequestTimeout(25*time.Millisecond))
	_, err := client.Health(context.Background())
	if err == nil || !strings.Contains(err.Error(), "read grasshopper response") {
		t.Fatalf("error = %v, want read timeout", err)
	}
}

func startFakeServer(t *testing.T, handler func(Request) Response) (string, func()) {
	t.Helper()
	return startRawServer(t, func(conn net.Conn) {
		reader := bufio.NewReader(conn)
		line, err := reader.ReadBytes('\n')
		if err != nil {
			t.Errorf("read request: %v", err)
			return
		}
		var request Request
		if err := json.Unmarshal(line, &request); err != nil {
			t.Errorf("decode request: %v", err)
			return
		}
		response := handler(request)
		if err := json.NewEncoder(conn).Encode(response); err != nil {
			t.Errorf("write response: %v", err)
		}
	})
}

func startRawServer(t *testing.T, handler func(net.Conn)) (string, func()) {
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
		handler(conn)
	}()
	stop := func() {
		_ = listener.Close()
		<-done
	}
	return listener.Addr().String(), stop
}

func successResponse(t *testing.T, id string, result any) Response {
	t.Helper()
	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("marshal result: %v", err)
	}
	return Response{ID: id, OK: true, Result: data}
}
