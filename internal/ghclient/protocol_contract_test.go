package ghclient

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestProtocolContractFixtures(t *testing.T) {
	tests := []struct {
		name string
		file string
		want any
	}{
		{
			name: "health request",
			file: "health_request.json",
			want: Request{ID: "contract-1", Method: "health"},
		},
		{
			name: "health response",
			file: "health_response.json",
			want: Response{ID: "contract-1", OK: true},
		},
		{
			name: "add component request",
			file: "add_component_request.json",
			want: Request{ID: "contract-2", Method: "add_component"},
		},
		{
			name: "add component response",
			file: "add_component_response.json",
			want: Response{ID: "contract-2", OK: true},
		},
		{
			name: "connect request",
			file: "connect_request.json",
			want: Request{ID: "contract-3", Method: "connect"},
		},
		{
			name: "error response",
			file: "error_response.json",
			want: Response{ID: "contract-4", OK: false, Error: &ProtocolError{Code: "component_not_found"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := readFixture(t, tt.file)
			switch want := tt.want.(type) {
			case Request:
				var got Request
				if err := json.Unmarshal(data, &got); err != nil {
					t.Fatalf("decode request: %v", err)
				}
				if got.ID != want.ID || got.Method != want.Method {
					t.Fatalf("request = %+v, want id=%q method=%q", got, want.ID, want.Method)
				}
			case Response:
				var got Response
				if err := json.Unmarshal(data, &got); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if got.ID != want.ID || got.OK != want.OK {
					t.Fatalf("response = %+v, want id=%q ok=%v", got, want.ID, want.OK)
				}
				if want.Error != nil {
					if got.Error == nil || got.Error.Code != want.Error.Code {
						t.Fatalf("error = %+v, want code %q", got.Error, want.Error.Code)
					}
				}
			default:
				t.Fatalf("unsupported fixture type %T", tt.want)
			}
		})
	}
}

func TestAddComponentContractParams(t *testing.T) {
	var request Request
	if err := json.Unmarshal(readFixture(t, "add_component_request.json"), &request); err != nil {
		t.Fatalf("decode request: %v", err)
	}
	data, err := json.Marshal(request.Params)
	if err != nil {
		t.Fatalf("marshal params: %v", err)
	}
	var params AddComponentParams
	if err := json.Unmarshal(data, &params); err != nil {
		t.Fatalf("decode params: %v", err)
	}
	if params.Name != "Addition" || params.X != 10 || params.Y != 20 {
		t.Fatalf("params = %+v", params)
	}
}

func readFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", "protocol", name))
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	return data
}
