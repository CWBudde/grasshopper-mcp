package ghclient

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"sync/atomic"
	"time"
)

const DefaultAddress = "127.0.0.1:47820"

var ErrGrasshopperUnavailable = errors.New("grasshopper adapter is unavailable")

var requestCounter uint64

type Client struct {
	address        string
	dialTimeout    time.Duration
	requestTimeout time.Duration
}

type Option func(*Client)

func WithDialTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.dialTimeout = timeout
	}
}

func WithRequestTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.requestTimeout = timeout
	}
}

func New(address string, options ...Option) *Client {
	if address == "" {
		address = DefaultAddress
	}
	client := &Client{
		address:        address,
		dialTimeout:    2 * time.Second,
		requestTimeout: 10 * time.Second,
	}
	for _, option := range options {
		option(client)
	}
	return client
}

func (c *Client) Address() string {
	return c.address
}

func (c *Client) Health(ctx context.Context) (HealthResult, error) {
	var result HealthResult
	if err := c.call(ctx, "health", nil, &result); err != nil {
		return HealthResult{}, err
	}
	return result, nil
}

func (c *Client) DocumentInfo(ctx context.Context) (DocumentInfoResult, error) {
	var result DocumentInfoResult
	if err := c.call(ctx, "document_info", nil, &result); err != nil {
		return DocumentInfoResult{}, err
	}
	return result, nil
}

func (c *Client) ListComponents(ctx context.Context) (ListComponentsResult, error) {
	var result ListComponentsResult
	if err := c.call(ctx, "list_components", nil, &result); err != nil {
		return ListComponentsResult{}, err
	}
	return result, nil
}

func (c *Client) RunSolver(ctx context.Context) (RunSolverResult, error) {
	var result RunSolverResult
	if err := c.call(ctx, "run_solver", nil, &result); err != nil {
		return RunSolverResult{}, err
	}
	return result, nil
}

func (c *Client) AddComponent(ctx context.Context, params AddComponentParams) (AddComponentResult, error) {
	var result AddComponentResult
	if err := c.call(ctx, "add_component", params, &result); err != nil {
		return AddComponentResult{}, err
	}
	return result, nil
}

func (c *Client) SetInput(ctx context.Context, params SetInputParams) (SetInputResult, error) {
	var result SetInputResult
	if err := c.call(ctx, "set_input", params, &result); err != nil {
		return SetInputResult{}, err
	}
	return result, nil
}

func (c *Client) Connect(ctx context.Context, params ConnectParams) (ConnectResult, error) {
	var result ConnectResult
	if err := c.call(ctx, "connect", params, &result); err != nil {
		return ConnectResult{}, err
	}
	return result, nil
}

func (c *Client) GetOutput(ctx context.Context, params GetOutputParams) (GetOutputResult, error) {
	var result GetOutputResult
	if err := c.call(ctx, "get_output", params, &result); err != nil {
		return GetOutputResult{}, err
	}
	return result, nil
}

func (c *Client) call(ctx context.Context, method string, params any, result any) error {
	ctx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()

	dialer := net.Dialer{Timeout: c.dialTimeout}
	conn, err := dialer.DialContext(ctx, "tcp", c.address)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrGrasshopperUnavailable, err)
	}
	defer conn.Close()

	if deadline, ok := ctx.Deadline(); ok {
		if err := conn.SetDeadline(deadline); err != nil {
			return err
		}
	}

	request := Request{
		ID:     nextRequestID(),
		Method: method,
		Params: params,
	}
	if err := json.NewEncoder(conn).Encode(request); err != nil {
		return fmt.Errorf("send grasshopper request: %w", err)
	}

	line, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		return fmt.Errorf("read grasshopper response: %w", err)
	}

	var response Response
	if err := json.Unmarshal(line, &response); err != nil {
		return fmt.Errorf("decode grasshopper response: %w", err)
	}
	if response.ID != request.ID {
		return fmt.Errorf("grasshopper response id %q does not match request id %q", response.ID, request.ID)
	}
	if !response.OK {
		if response.Error == nil {
			return &ProtocolError{Code: "protocol_error", Message: "Grasshopper returned an error without details"}
		}
		return response.Error
	}
	if len(response.Result) == 0 || result == nil {
		return nil
	}
	if err := json.Unmarshal(response.Result, result); err != nil {
		return fmt.Errorf("decode grasshopper result: %w", err)
	}
	return nil
}

func nextRequestID() string {
	id := atomic.AddUint64(&requestCounter, 1)
	return fmt.Sprintf("ghmcp-%d", id)
}

func (e *ProtocolError) Error() string {
	if e.Message == "" {
		return e.Code
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}
