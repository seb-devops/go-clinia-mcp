package supabase

import (
	"clinia-doc/mcp/services"
	"context"
	"encoding/json"
	"reflect"
	"testing"
)

type mockSupabaseClient struct {
	rpcFunc func(fn, mode string, params map[string]interface{}) string
}

func (m *mockSupabaseClient) Rpc(fn, mode string, params map[string]interface{}) string {
	return m.rpcFunc(fn, mode, params)
}

type testFields struct {
	supabase *mockSupabaseClient
}

type testArgs struct {
	vector []float32
	limit  int
}

func TestClient_SearchByVector(t *testing.T) {
	tests := []struct {
		testname string
		fields   testFields
		args     testArgs
		want     []services.SearchResult
		wantErr  bool
	}{
		{
			testname: "success",
			fields: testFields{
				supabase: &mockSupabaseClient{
					rpcFunc: func(fn, mode string, params map[string]interface{}) string {
						results := []services.SearchResult{{ID: 1, URL: "url", ChunkNumber: 1, Title: "title", Summary: "summary", Content: "content", Metadata: map[string]interface{}{"foo": "bar"}, Similarity: 0.99}}
						b, _ := json.Marshal(results)
						return string(b)
					},
				},
			},
			args:    testArgs{vector: []float32{1.0, 2.0}, limit: 1},
			want:    []services.SearchResult{{ID: 1, URL: "url", ChunkNumber: 1, Title: "title", Summary: "summary", Content: "content", Metadata: map[string]interface{}{"foo": "bar"}, Similarity: 0.99}},
			wantErr: false,
		},
		{
			testname: "supabase_returns_empty_string",
			fields: testFields{
				supabase: &mockSupabaseClient{
					rpcFunc: func(fn, mode string, params map[string]interface{}) string {
						return ""
					},
				},
			},
			args:    testArgs{vector: []float32{1.0, 2.0}, limit: 1},
			want:    nil,
			wantErr: true,
		},
		{
			testname: "supabase_returns_invalid_json",
			fields: testFields{
				supabase: &mockSupabaseClient{
					rpcFunc: func(fn, mode string, params map[string]interface{}) string {
						return "not json"
					},
				},
			},
			args:    testArgs{vector: []float32{1.0, 2.0}, limit: 1},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			c := &Client{
				supabase: tt.fields.supabase,
			}
			got, err := c.SearchByVector(context.Background(), tt.args.vector, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchByVector() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchByVector() = %v, want %v", got, tt.want)
			}
		})
	}
}
