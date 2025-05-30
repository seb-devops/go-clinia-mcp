package services

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

// Mock implementations for testing
type mockEmbeddingProvider struct {
	getEmbeddingFunc func(ctx context.Context, text string) ([]float32, error)
}

func (m *mockEmbeddingProvider) GetEmbedding(ctx context.Context, text string) ([]float32, error) {
	if m.getEmbeddingFunc != nil {
		return m.getEmbeddingFunc(ctx, text)
	}
	return []float32{1.0, 2.0, 3.0}, nil
}

type mockVectorSearchProvider struct {
	searchByVectorFunc func(ctx context.Context, vector []float32, limit int) ([]SearchResult, error)
}

func (m *mockVectorSearchProvider) SearchByVector(ctx context.Context, vector []float32, limit int) ([]SearchResult, error) {
	if m.searchByVectorFunc != nil {
		return m.searchByVectorFunc(ctx, vector, limit)
	}
	return []SearchResult{
		{
			ID:         1,
			URL:        "https://example.com/doc1",
			Title:      "Test Document",
			Content:    "Test content",
			Similarity: 0.95,
		},
	}, nil
}

func TestSearchService_SearchByText(t *testing.T) {
	tests := []struct {
		name                string
		embeddingFunc       func(ctx context.Context, text string) ([]float32, error)
		searchFunc          func(ctx context.Context, vector []float32, limit int) ([]SearchResult, error)
		query               string
		limit               int
		expectedResults     []SearchResult
		expectedError       bool
		expectedErrorString string
	}{
		{
			name: "successful search",
			embeddingFunc: func(ctx context.Context, text string) ([]float32, error) {
				return []float32{1.0, 2.0, 3.0}, nil
			},
			searchFunc: func(ctx context.Context, vector []float32, limit int) ([]SearchResult, error) {
				return []SearchResult{
					{ID: 1, Title: "Test Doc", Similarity: 0.9},
				}, nil
			},
			query: "test query",
			limit: 5,
			expectedResults: []SearchResult{
				{ID: 1, Title: "Test Doc", Similarity: 0.9},
			},
			expectedError: false,
		},
		{
			name: "embedding generation fails",
			embeddingFunc: func(ctx context.Context, text string) ([]float32, error) {
				return nil, errors.New("embedding failed")
			},
			query:               "test query",
			limit:               5,
			expectedError:       true,
			expectedErrorString: "failed to generate embedding",
		},
		{
			name: "vector search fails",
			embeddingFunc: func(ctx context.Context, text string) ([]float32, error) {
				return []float32{1.0, 2.0, 3.0}, nil
			},
			searchFunc: func(ctx context.Context, vector []float32, limit int) ([]SearchResult, error) {
				return nil, errors.New("search failed")
			},
			query:               "test query",
			limit:               5,
			expectedError:       true,
			expectedErrorString: "failed to perform vector search",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			embeddingProvider := &mockEmbeddingProvider{getEmbeddingFunc: tt.embeddingFunc}
			searchProvider := &mockVectorSearchProvider{searchByVectorFunc: tt.searchFunc}

			service := NewSearchService(embeddingProvider, searchProvider)

			results, err := service.SearchByText(context.Background(), tt.query, tt.limit)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.expectedErrorString != "" && !contains(err.Error(), tt.expectedErrorString) {
					t.Errorf("Expected error to contain '%s', got '%s'", tt.expectedErrorString, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if !reflect.DeepEqual(results, tt.expectedResults) {
					t.Errorf("Expected results %v, got %v", tt.expectedResults, results)
				}
			}
		})
	}
}

func TestSearchService_SearchByVector(t *testing.T) {
	tests := []struct {
		name            string
		searchFunc      func(ctx context.Context, vector []float32, limit int) ([]SearchResult, error)
		vector          []float32
		limit           int
		expectedResults []SearchResult
		expectedError   bool
	}{
		{
			name: "successful vector search",
			searchFunc: func(ctx context.Context, vector []float32, limit int) ([]SearchResult, error) {
				return []SearchResult{
					{ID: 1, Title: "Vector Doc", Similarity: 0.8},
				}, nil
			},
			vector: []float32{1.0, 2.0, 3.0},
			limit:  3,
			expectedResults: []SearchResult{
				{ID: 1, Title: "Vector Doc", Similarity: 0.8},
			},
			expectedError: false,
		},
		{
			name: "search provider fails",
			searchFunc: func(ctx context.Context, vector []float32, limit int) ([]SearchResult, error) {
				return nil, errors.New("provider error")
			},
			vector:        []float32{1.0, 2.0, 3.0},
			limit:         3,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			embeddingProvider := &mockEmbeddingProvider{}
			searchProvider := &mockVectorSearchProvider{searchByVectorFunc: tt.searchFunc}

			service := NewSearchService(embeddingProvider, searchProvider)

			results, err := service.SearchByVector(context.Background(), tt.vector, tt.limit)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if !reflect.DeepEqual(results, tt.expectedResults) {
					t.Errorf("Expected results %v, got %v", tt.expectedResults, results)
				}
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && func() bool {
		for i := 0; i <= len(s)-len(substr); i++ {
			if s[i:i+len(substr)] == substr {
				return true
			}
		}
		return false
	}()
}
