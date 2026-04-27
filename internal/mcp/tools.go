package mcp

func tools() []map[string]any {
	emptySchema := map[string]any{
		"type":       "object",
		"properties": map[string]any{},
	}
	parameterRefSchema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"componentId": map[string]any{"type": "string"},
			"parameter":   map[string]any{"type": "string"},
		},
		"required": []string{"componentId", "parameter"},
	}
	return []map[string]any{
		{
			"name":        "grasshopper_health",
			"description": "Check whether the Grasshopper adapter is reachable and report its version.",
			"inputSchema": emptySchema,
		},
		{
			"name":        "grasshopper_document_info",
			"description": "Return basic information about the active Grasshopper document.",
			"inputSchema": emptySchema,
		},
		{
			"name":        "grasshopper_list_components",
			"description": "List known Grasshopper components available through the adapter.",
			"inputSchema": emptySchema,
		},
		{
			"name":        "grasshopper_run_solver",
			"description": "Ask Grasshopper to recompute the active document.",
			"inputSchema": emptySchema,
		},
		{
			"name":        "grasshopper_add_component",
			"description": "Place a Grasshopper component in the active document.",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"name":        map[string]any{"type": "string"},
					"nickname":    map[string]any{"type": "string"},
					"category":    map[string]any{"type": "string"},
					"subcategory": map[string]any{"type": "string"},
					"x":           map[string]any{"type": "number"},
					"y":           map[string]any{"type": "number"},
				},
				"required": []string{"name"},
			},
		},
		{
			"name":        "grasshopper_set_input",
			"description": "Assign a scalar or list value to an input parameter.",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"target": parameterRefSchema,
					"value":  map[string]any{},
				},
				"required": []string{"target", "value"},
			},
		},
		{
			"name":        "grasshopper_connect",
			"description": "Connect an output parameter to an input parameter.",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"source": parameterRefSchema,
					"target": parameterRefSchema,
				},
				"required": []string{"source", "target"},
			},
		},
		{
			"name":        "grasshopper_get_output",
			"description": "Read the current value of an output parameter.",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"source": parameterRefSchema,
				},
				"required": []string{"source"},
			},
		},
	}
}
