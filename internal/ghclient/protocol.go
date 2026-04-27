package ghclient

import "encoding/json"

type Request struct {
	ID     string `json:"id"`
	Method string `json:"method"`
	Params any    `json:"params,omitempty"`
}

type Response struct {
	ID     string          `json:"id"`
	OK     bool            `json:"ok"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  *ProtocolError  `json:"error,omitempty"`
}

type ProtocolError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type HealthResult struct {
	Version           string `json:"version"`
	ActiveDocument    bool   `json:"activeDocument"`
	GrasshopperLoaded bool   `json:"grasshopperLoaded"`
}

type DocumentInfoResult struct {
	DocumentName      string `json:"documentName"`
	ObjectCount       int    `json:"objectCount"`
	HasActiveDocument bool   `json:"hasActiveDocument"`
}

type ComponentSummary struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Nickname    string `json:"nickname"`
	Category    string `json:"category"`
	Subcategory string `json:"subcategory"`
	Description string `json:"description"`
}

type ListComponentsResult struct {
	Components []ComponentSummary `json:"components"`
}

type RunSolverResult struct {
	Completed bool   `json:"completed"`
	Message   string `json:"message,omitempty"`
}

type AddComponentParams struct {
	Name        string  `json:"name"`
	Nickname    string  `json:"nickname,omitempty"`
	Category    string  `json:"category,omitempty"`
	Subcategory string  `json:"subcategory,omitempty"`
	X           float64 `json:"x,omitempty"`
	Y           float64 `json:"y,omitempty"`
}

type AddComponentResult struct {
	ComponentID string `json:"componentId"`
	Name        string `json:"name"`
	Nickname    string `json:"nickname,omitempty"`
}

type ParameterRef struct {
	ComponentID string `json:"componentId"`
	Parameter   string `json:"parameter"`
}

type SetInputParams struct {
	Target ParameterRef `json:"target"`
	Value  any          `json:"value"`
}

type SetInputResult struct {
	Updated bool   `json:"updated"`
	Message string `json:"message,omitempty"`
}

type ConnectParams struct {
	Source ParameterRef `json:"source"`
	Target ParameterRef `json:"target"`
}

type ConnectResult struct {
	Connected bool   `json:"connected"`
	Message   string `json:"message,omitempty"`
}

type GetOutputParams struct {
	Source ParameterRef `json:"source"`
}

type GetOutputResult struct {
	Value any    `json:"value"`
	Type  string `json:"type,omitempty"`
}
