package model

type ComponentID string

type ParameterRef struct {
	Component ComponentID `json:"component"`
	Parameter string      `json:"parameter"`
}
