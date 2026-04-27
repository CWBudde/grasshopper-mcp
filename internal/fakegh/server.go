package fakegh

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"

	"github.com/cwbudde/grasshopper-mcp/internal/ghclient"
	"github.com/cwbudde/grasshopper-mcp/internal/version"
)

type Server struct {
	addressMu sync.RWMutex
	address   string

	mu          sync.Mutex
	components  map[string]ghclient.AddComponentResult
	inputs      map[ghclient.ParameterRef]any
	connections []ghclient.ConnectParams
	nextID      uint64
}

func New(address string) *Server {
	if address == "" {
		address = ghclient.DefaultAddress
	}
	return &Server{
		address:    address,
		components: make(map[string]ghclient.AddComponentResult),
		inputs:     make(map[ghclient.ParameterRef]any),
	}
}

func (s *Server) Address() string {
	s.addressMu.RLock()
	defer s.addressMu.RUnlock()
	return s.address
}

func (s *Server) ListenAndServe(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("listen fake grasshopper adapter: %w", err)
	}
	defer listener.Close()
	s.addressMu.Lock()
	s.address = listener.Addr().String()
	s.addressMu.Unlock()

	go func() {
		<-ctx.Done()
		_ = listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if ctx.Err() != nil || errors.Is(err, net.ErrClosed) {
				return nil
			}
			return fmt.Errorf("accept fake grasshopper connection: %w", err)
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	line, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		return
	}

	var request ghclient.Request
	if err := json.Unmarshal(line, &request); err != nil {
		_ = json.NewEncoder(conn).Encode(ghclient.Response{
			OK:    false,
			Error: &ghclient.ProtocolError{Code: "invalid_json", Message: err.Error()},
		})
		return
	}

	response := s.route(request)
	_ = json.NewEncoder(conn).Encode(response)
}

func (s *Server) route(request ghclient.Request) ghclient.Response {
	switch request.Method {
	case "health":
		return success(request.ID, ghclient.HealthResult{
			Version:           "fake-" + version.Version,
			ActiveDocument:    true,
			GrasshopperLoaded: true,
		})
	case "document_info":
		s.mu.Lock()
		count := len(s.components)
		s.mu.Unlock()
		return success(request.ID, ghclient.DocumentInfoResult{
			DocumentName:      "Fake Grasshopper Document",
			ObjectCount:       count,
			HasActiveDocument: true,
		})
	case "list_components":
		return success(request.ID, ghclient.ListComponentsResult{Components: fakeCatalog()})
	case "run_solver":
		return success(request.ID, ghclient.RunSolverResult{Completed: true, Message: "Fake solver completed."})
	case "add_component":
		return s.addComponent(request)
	case "set_input":
		return s.setInput(request)
	case "connect":
		return s.connect(request)
	case "get_output":
		return s.getOutput(request)
	default:
		return failure(request.ID, "unknown_method", fmt.Sprintf("Unknown method %q.", request.Method))
	}
}

func (s *Server) addComponent(request ghclient.Request) ghclient.Response {
	var params ghclient.AddComponentParams
	if err := decodeParams(request.Params, &params); err != nil {
		return failure(request.ID, "invalid_arguments", err.Error())
	}
	if params.Name == "" && params.Nickname == "" {
		return failure(request.ID, "invalid_arguments", "add_component requires name or nickname")
	}

	id := fmt.Sprintf("component-%d", atomic.AddUint64(&s.nextID, 1))
	name := params.Name
	if name == "" {
		name = params.Nickname
	}
	result := ghclient.AddComponentResult{ComponentID: id, Name: name, Nickname: params.Nickname}

	s.mu.Lock()
	s.components[id] = result
	s.mu.Unlock()

	return success(request.ID, result)
}

func (s *Server) setInput(request ghclient.Request) ghclient.Response {
	var params ghclient.SetInputParams
	if err := decodeParams(request.Params, &params); err != nil {
		return failure(request.ID, "invalid_arguments", err.Error())
	}
	if params.Target.ComponentID == "" || params.Target.Parameter == "" {
		return failure(request.ID, "invalid_arguments", "target componentId and parameter are required")
	}

	s.mu.Lock()
	s.inputs[params.Target] = params.Value
	s.mu.Unlock()

	return success(request.ID, ghclient.SetInputResult{Updated: true})
}

func (s *Server) connect(request ghclient.Request) ghclient.Response {
	var params ghclient.ConnectParams
	if err := decodeParams(request.Params, &params); err != nil {
		return failure(request.ID, "invalid_arguments", err.Error())
	}
	if params.Source.ComponentID == "" || params.Source.Parameter == "" ||
		params.Target.ComponentID == "" || params.Target.Parameter == "" {
		return failure(request.ID, "invalid_arguments", "source and target componentId and parameter are required")
	}

	s.mu.Lock()
	s.connections = append(s.connections, params)
	s.mu.Unlock()

	return success(request.ID, ghclient.ConnectResult{Connected: true})
}

func (s *Server) getOutput(request ghclient.Request) ghclient.Response {
	var params ghclient.GetOutputParams
	if err := decodeParams(request.Params, &params); err != nil {
		return failure(request.ID, "invalid_arguments", err.Error())
	}
	if params.Source.ComponentID == "" || params.Source.Parameter == "" {
		return failure(request.ID, "invalid_arguments", "source componentId and parameter are required")
	}

	return success(request.ID, ghclient.GetOutputResult{Value: 42.0, Type: "number"})
}

func fakeCatalog() []ghclient.ComponentSummary {
	return []ghclient.ComponentSummary{
		{ID: "grasshopper.maths.operators.addition", Name: "Addition", Nickname: "Add", Category: "Maths", Subcategory: "Operators", Description: "Add numbers."},
		{ID: "grasshopper.params.input.number", Name: "Number", Nickname: "Num", Category: "Params", Subcategory: "Input", Description: "Numeric parameter."},
		{ID: "grasshopper.params.input.panel", Name: "Panel", Nickname: "Panel", Category: "Params", Subcategory: "Input", Description: "Text panel."},
	}
}

func decodeParams(params any, target any) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func success(id string, result any) ghclient.Response {
	data, err := json.Marshal(result)
	if err != nil {
		return failure(id, "internal_error", err.Error())
	}
	return ghclient.Response{ID: id, OK: true, Result: data}
}

func failure(id, code, message string) ghclient.Response {
	return ghclient.Response{ID: id, OK: false, Error: &ghclient.ProtocolError{Code: code, Message: message}}
}
