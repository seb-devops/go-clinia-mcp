package langchain

import (
	"context"
	"testing"
)

func TestQueryDecomposerDecomposeQuery(t *testing.T) {
	decomposer := setupDecomposer(t)
	tests := getTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			result, err := decomposer.DecomposeQuery(ctx, tt.query)
			if err != nil {
				t.Errorf("DecomposeQuery() error = %v", err)
				return
			}

			validateDecompositionResult(t, result, tt.query, tt.want)
			validateSubQueries(t, result.SubQueries)
			logSubQueries(t, result.SubQueries)
		})
	}
}

func TestQueryDecomposerEmptyQuery(t *testing.T) {
	decomposer := setupDecomposer(t)

	ctx := context.Background()
	result, err := decomposer.DecomposeQuery(ctx, "")
	if err == nil {
		t.Error("DecomposeQuery() expected error for empty query, got nil")
	}
	if result != nil {
		t.Error("DecomposeQuery() expected nil result for empty query, got result")
	}
}

func TestParseDecompositionResponse(t *testing.T) {
	decomposer := &QueryDecomposer{}

	tests := []struct {
		name     string
		response string
		want     int
		wantErr  bool
	}{
		{
			name: "Valid numbered list",
			response: `1. Data Collection - Gather historical sales data
                       2. Analysis - Analyze trends and patterns  
                       3. Prediction - Create forecasting model`,
			want:    3,
			wantErr: false,
		},
		{
			name: "Valid list without descriptions",
			response: `1. What are the benefits of AI?
					   2. What are the challenges of AI?
					   3. What is the future of AI?`,
			want:    3,
			wantErr: false,
		},
		{
			name:     "Empty response",
			response: "",
			want:     0,
			wantErr:  true,
		},
		{
			name:     "No numbered items",
			response: "This is just some text without numbers",
			want:     0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decomposer.parseDecompositionResponse(tt.response)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDecompositionResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("parseDecompositionResponse() got %d subqueries, want %d", len(got), tt.want)
			}
		})
	}
}
