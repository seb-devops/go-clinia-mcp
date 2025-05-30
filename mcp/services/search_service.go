package services

import (
	"context"
	"fmt"
)

// EmbeddingProvider defines the interface for generating embeddings
type EmbeddingProvider interface {
	GetEmbedding(ctx context.Context, text string) ([]float32, error)
}

// VectorSearchProvider defines the interface for vector-based search
type VectorSearchProvider interface {
	SearchByVector(ctx context.Context, vector []float32, limit int) ([]SearchResult, error)
}

// SearchResult represents a generic search result
type SearchResult struct {
	ID          int                    `json:"id"`
	URL         string                 `json:"url"`
	ChunkNumber int                    `json:"chunk_number"`
	Title       string                 `json:"title"`
	Summary     string                 `json:"summary"`
	Content     string                 `json:"content"`
	Metadata    map[string]interface{} `json:"metadata"`
	Similarity  float64                `json:"similarity"`
}

// SearchService coordinates embedding generation and vector search
type SearchService struct {
	embeddingProvider EmbeddingProvider
	searchProvider    VectorSearchProvider
}

// NewSearchService creates a new search service with injected dependencies
func NewSearchService(embeddingProvider EmbeddingProvider, searchProvider VectorSearchProvider) *SearchService {
	return &SearchService{
		embeddingProvider: embeddingProvider,
		searchProvider:    searchProvider,
	}
}

// SearchByText generates embeddings for text and performs vector search
func (s *SearchService) SearchByText(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	// Generate embedding for the query
	embedding, err := s.embeddingProvider.GetEmbedding(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	// Perform vector search
	results, err := s.searchProvider.SearchByVector(ctx, embedding, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to perform vector search: %w", err)
	}

	return results, nil
}

// SearchByVector performs vector search directly (useful when embedding is already available)
func (s *SearchService) SearchByVector(ctx context.Context, vector []float32, limit int) ([]SearchResult, error) {
	return s.searchProvider.SearchByVector(ctx, vector, limit)
}
