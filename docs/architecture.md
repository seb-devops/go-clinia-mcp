# Architecture Documentation

## Overview

This document describes the maintainable architecture implemented in the MCP Clinia Doc project, focusing on separation of concerns, dependency injection, and testability.

## Architecture Principles

### 1. Separation of Concerns
- **Services Layer**: Business logic and coordination
- **Provider Layer**: Infrastructure-specific implementations
- **Interface Layer**: Clean contracts between components

### 2. Dependency Injection
- Components receive their dependencies through constructor injection
- Interfaces define contracts, not concrete implementations
- Easy to swap implementations for testing or different environments

### 3. Single Responsibility Principle
- Each component has one reason to change
- Clear boundaries between embedding, search, and coordination logic

## Layer Structure

```
┌─────────────────────────────────────────┐
│               MCP Server                │
│              (mcp/mcp.go)              │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│             Service Layer               │
│          (mcp/services/)                │
│                                         │
│  ┌─────────────────────────────────┐   │
│  │        SearchService            │   │
│  │                                 │   │
│  │  - SearchByText()              │   │
│  │  - SearchByVector()            │   │
│  └─────────────────────────────────┘   │
└─────────────┬───────────┬───────────────┘
              │           │
    ┌─────────▼──┐    ┌──▼──────────┐
    │ Embedding  │    │   Vector    │
    │ Provider   │    │   Search    │
    │ Interface  │    │  Provider   │
    │            │    │  Interface  │
    └─────────┬──┘    └──┬──────────┘
              │          │
    ┌─────────▼──┐    ┌──▼──────────┐
    │   OpenAI   │    │  Supabase   │
    │   Client   │    │   Client    │
    │(mcp/openai)│    │(mcp/supabase)│
    └────────────┘    └─────────────┘
```

## Component Details

### Service Layer (`mcp/services/`)

#### SearchService
The main service that coordinates embedding generation and vector search.

**Responsibilities:**
- Coordinate between embedding and search providers
- Implement business logic for text-to-vector search
- Provide clean API for consumers

**Interfaces:**
```go
type EmbeddingProvider interface {
    GetEmbedding(ctx context.Context, text string) ([]float32, error)
}

type VectorSearchProvider interface {
    SearchByVector(ctx context.Context, vector []float32, limit int) ([]SearchResult, error)
}
```

**Benefits:**
- ✅ Easy to test with mock implementations
- ✅ Can swap embedding providers (OpenAI → Anthropic, etc.)
- ✅ Can swap search backends (Supabase → Pinecone, etc.)
- ✅ Clear separation of concerns

### Provider Layer

#### OpenAI Client (`mcp/openai/`)
**Purpose:** Implements `EmbeddingProvider` interface
**Responsibilities:**
- Generate embeddings using OpenAI API
- Handle API authentication and error handling
- Provide foundation for other OpenAI services (completion, images)

#### Supabase Client (`mcp/supabase/`)
**Purpose:** Implements `VectorSearchProvider` interface
**Responsibilities:**
- Perform vector similarity search
- Handle Supabase-specific RPC calls
- Convert results to standard format

## Benefits of This Architecture

### 1. **Testability**
```go
// Easy to test with mocks
mockEmbedding := &mockEmbeddingProvider{}
mockSearch := &mockVectorSearchProvider{}
service := NewSearchService(mockEmbedding, mockSearch)
```

### 2. **Flexibility**
```go
// Easy to swap implementations
// Development: use mock providers
// Production: use real providers
// Testing: use different configurations
```

### 3. **Extensibility**
```go
// Add new providers without changing existing code
type AnthropicProvider struct { ... }
func (a *AnthropicProvider) GetEmbedding(...) { ... }

// Use with existing service
service := NewSearchService(anthropicProvider, supabaseClient)
```

### 4. **Maintainability**
- Clear boundaries between components
- Single responsibility for each layer
- Dependency injection makes relationships explicit
- Interface-based design enables easy refactoring

## Migration from Previous Architecture

### Before (Tightly Coupled)
```go
// Supabase client had direct OpenAI dependency
type SupabaseClient struct {
    supabase SupabaseClient
    embedder openai.OpenAIEmbedder  // ❌ Tight coupling
}

// Business logic mixed with infrastructure
func (c *Client) SearchWithQuery(query string) ([]Result, error) {
    embedding, err := c.embedder.GetEmbedding(query)  // ❌ Mixed concerns
    return c.SearchByVector(embedding)
}
```

### After (Loosely Coupled)
```go
// Clean separation with interfaces
type SearchService struct {
    embeddingProvider EmbeddingProvider    // ✅ Interface dependency
    searchProvider    VectorSearchProvider // ✅ Interface dependency
}

// Pure business logic
func (s *SearchService) SearchByText(query string) ([]Result, error) {
    embedding, err := s.embeddingProvider.GetEmbedding(query)  // ✅ Clear responsibility
    return s.searchProvider.SearchByVector(embedding)
}
```

## Future Enhancements

### 1. **Multiple Embedding Providers**
```go
type EmbeddingRouter struct {
    providers map[string]EmbeddingProvider
}

func (r *EmbeddingRouter) GetEmbedding(ctx context.Context, text string, provider string) ([]float32, error) {
    return r.providers[provider].GetEmbedding(ctx, text)
}
```

### 2. **Caching Layer**
```go
type CachedEmbeddingProvider struct {
    provider EmbeddingProvider
    cache    Cache
}
```

### 3. **Metrics and Observability**
```go
type InstrumentedSearchService struct {
    service SearchService
    metrics MetricsCollector
}
```

### 4. **Configuration-Driven Setup**
```yaml
providers:
  embedding:
    type: "openai"
    config:
      api_key: "${OPENAI_API_KEY}"
  search:
    type: "supabase"
    config:
      url: "${SUPABASE_URL}"
```

## Testing Strategy

### Unit Tests
- Test each provider independently with mocks
- Test service layer with mock providers
- Test error handling and edge cases

### Integration Tests
- Test with real providers in staging environment
- Test end-to-end workflows
- Performance and reliability testing

### Example Test Structure
```go
func TestSearchService_SearchByText(t *testing.T) {
    // Arrange
    mockEmbedding := &mockEmbeddingProvider{...}
    mockSearch := &mockVectorSearchProvider{...}
    service := NewSearchService(mockEmbedding, mockSearch)
    
    // Act
    results, err := service.SearchByText(ctx, "test query", 10)
    
    // Assert
    assert.NoError(t, err)
    assert.Len(t, results, expectedCount)
}
```

## Summary

This architecture provides:
- **Maintainability** through clear separation of concerns
- **Testability** through dependency injection and interfaces
- **Flexibility** to swap implementations
- **Extensibility** to add new features without breaking existing code

The service layer acts as a coordinator, while provider layers handle infrastructure concerns. This enables easy testing, development, and future enhancements while maintaining clean, understandable code. 