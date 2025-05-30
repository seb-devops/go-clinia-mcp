package openai

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
)

type mockHTTPClient struct {
	postFunc func(ctx context.Context, path string, body interface{}, out interface{}) error
}

func (m *mockHTTPClient) Post(ctx context.Context, path string, body interface{}, out interface{}) error {
	return m.postFunc(ctx, path, body, out)
}

type testCase struct {
	name       string
	postFunc   func(ctx context.Context, path string, body interface{}, out interface{}) error
	want       []float32
	wantErr    bool
	wantErrMsg string
}

func TestClient_GetEmbedding(t *testing.T) {
	tests := []testCase{
		{
			name: "success",
			postFunc: func(ctx context.Context, path string, body interface{}, out interface{}) error {
				// Simulate a successful embedding response
				resp := out.(*struct {
					Data []struct {
						Embedding []float32 `json:"embedding"`
					} `json:"data"`
				})
				resp.Data = []struct {
					Embedding []float32 `json:"embedding"`
				}{{Embedding: []float32{1.1, 2.2, 3.3}}}
				return nil
			},
			want:    []float32{1.1, 2.2, 3.3},
			wantErr: false,
		},
		{
			name: "api error",
			postFunc: func(ctx context.Context, path string, body interface{}, out interface{}) error {
				return errors.New("api error")
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "api error",
		},
		{
			name: "no embedding returned",
			postFunc: func(ctx context.Context, path string, body interface{}, out interface{}) error {
				resp := out.(*struct {
					Data []struct {
						Embedding []float32 `json:"embedding"`
					} `json:"data"`
				})
				resp.Data = nil
				return nil
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "no embedding returned",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{httpClient: &mockHTTPClient{postFunc: tt.postFunc}}
			got, err := c.GetEmbedding(context.Background(), "test")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEmbedding() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && tt.wantErrMsg != "" && !strings.Contains(err.Error(), tt.wantErrMsg) {
				t.Errorf("GetEmbedding() error = %v, wantErrMsg substring %v", err, tt.wantErrMsg)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEmbedding() = %v, want %v", got, tt.want)
			}
		})
	}
}
