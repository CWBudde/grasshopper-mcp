package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/cwbudde/grasshopper-mcp/internal/ghclient"
)

type Server struct {
	client *ghclient.Client
	in     io.Reader
	out    io.Writer
}

type Option func(*Server)

func WithIO(in io.Reader, out io.Writer) Option {
	return func(s *Server) {
		s.in = in
		s.out = out
	}
}

func NewServer(client *ghclient.Client, options ...Option) *Server {
	server := &Server{
		client: client,
	}
	for _, option := range options {
		option(server)
	}
	return server
}

func (s *Server) Run(ctx context.Context) error {
	if s.client == nil {
		return errors.New("mcp server requires a grasshopper client")
	}
	if s.in == nil || s.out == nil {
		return errors.New("mcp server requires input and output streams")
	}

	scanner := bufio.NewScanner(s.in)
	writer := bufio.NewWriter(s.out)
	defer writer.Flush()

	for scanner.Scan() {
		if err := ctx.Err(); err != nil {
			return err
		}
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		response := s.handleLine(ctx, line)
		if response == nil {
			continue
		}
		if err := json.NewEncoder(writer).Encode(response); err != nil {
			return fmt.Errorf("write mcp response: %w", err)
		}
		if err := writer.Flush(); err != nil {
			return fmt.Errorf("flush mcp response: %w", err)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read mcp request: %w", err)
	}
	return nil
}

func (s *Server) handleLine(ctx context.Context, line []byte) *rpcResponse {
	var request rpcRequest
	if err := json.Unmarshal(line, &request); err != nil {
		return rpcError(nil, -32700, "Parse error")
	}
	if request.ID == nil {
		return nil
	}

	switch request.Method {
	case "initialize":
		return rpcResult(request.ID, map[string]any{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]any{
				"tools": map[string]any{},
			},
			"serverInfo": map[string]any{
				"name":    "grasshopper-mcp",
				"version": "0.1.0",
			},
		})
	case "tools/list":
		return rpcResult(request.ID, map[string]any{"tools": tools()})
	case "tools/call":
		var params toolCallParams
		if err := json.Unmarshal(request.Params, &params); err != nil {
			return rpcError(request.ID, -32602, "Invalid tool call parameters")
		}
		result := s.callTool(ctx, params)
		return rpcResult(request.ID, result)
	default:
		return rpcError(request.ID, -32601, "Method not found")
	}
}

func (s *Server) callTool(ctx context.Context, params toolCallParams) toolResult {
	switch params.Name {
	case "grasshopper_health":
		result, err := s.client.Health(ctx)
		return toolJSON(result, err)
	case "grasshopper_document_info":
		result, err := s.client.DocumentInfo(ctx)
		return toolJSON(result, err)
	case "grasshopper_list_components":
		result, err := s.client.ListComponents(ctx)
		return toolJSON(result, err)
	case "grasshopper_run_solver":
		result, err := s.client.RunSolver(ctx)
		return toolJSON(result, err)
	case "grasshopper_add_component":
		var args ghclient.AddComponentParams
		if err := decodeArguments(params.Arguments, &args); err != nil {
			return toolError(err.Error())
		}
		result, err := s.client.AddComponent(ctx, args)
		return toolJSON(result, err)
	case "grasshopper_set_input":
		var args ghclient.SetInputParams
		if err := decodeArguments(params.Arguments, &args); err != nil {
			return toolError(err.Error())
		}
		result, err := s.client.SetInput(ctx, args)
		return toolJSON(result, err)
	case "grasshopper_connect":
		var args ghclient.ConnectParams
		if err := decodeArguments(params.Arguments, &args); err != nil {
			return toolError(err.Error())
		}
		result, err := s.client.Connect(ctx, args)
		return toolJSON(result, err)
	case "grasshopper_get_output":
		var args ghclient.GetOutputParams
		if err := decodeArguments(params.Arguments, &args); err != nil {
			return toolError(err.Error())
		}
		result, err := s.client.GetOutput(ctx, args)
		return toolJSON(result, err)
	default:
		return toolError(fmt.Sprintf("Unknown tool %q.", params.Name))
	}
}

func decodeArguments(arguments map[string]any, target any) error {
	data, err := json.Marshal(arguments)
	if err != nil {
		return fmt.Errorf("encode tool arguments: %w", err)
	}
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("decode tool arguments: %w", err)
	}
	return nil
}

func toolJSON(value any, err error) toolResult {
	if err != nil {
		return toolError(err.Error())
	}
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return toolError(fmt.Sprintf("Encode tool result: %v", err))
	}
	return toolResult{
		Content: []contentItem{{Type: "text", Text: string(data)}},
	}
}

func toolError(message string) toolResult {
	return toolResult{
		IsError: true,
		Content: []contentItem{{Type: "text", Text: message}},
	}
}

func rpcResult(id json.RawMessage, result any) *rpcResponse {
	return &rpcResponse{JSONRPC: "2.0", ID: id, Result: result}
}

func rpcError(id json.RawMessage, code int, message string) *rpcResponse {
	return &rpcResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error:   &rpcErrorObject{Code: code, Message: message},
	}
}

type rpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type rpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Result  any             `json:"result,omitempty"`
	Error   *rpcErrorObject `json:"error,omitempty"`
}

type rpcErrorObject struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type toolCallParams struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments,omitempty"`
}

type toolResult struct {
	Content []contentItem `json:"content"`
	IsError bool          `json:"isError,omitempty"`
}

type contentItem struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
