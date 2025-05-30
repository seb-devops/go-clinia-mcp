package openai

import (
	"context"
	"fmt"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// HTTPClient interface for making HTTP requests to OpenAI API
type HTTPClient interface {
	Post(ctx context.Context, path string, body interface{}, out interface{}) error
}

// Client represents a generic OpenAI client
// Implements services.EmbeddingProvider interface
type Client struct {
	httpClient HTTPClient
}

// ClientAdapter adapts the official OpenAI client to our HTTPClient interface
type ClientAdapter struct {
	client *openai.Client
}

func (a *ClientAdapter) Post(ctx context.Context, path string, body interface{}, out interface{}) error {
	return a.client.Post(ctx, path, body, out)
}

// NewClient creates a new OpenAI client
func NewClient() (*Client, error) {
	openaiKey := os.Getenv("OPENAI_API_KEY")
	if openaiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	client := openai.NewClient(option.WithAPIKey(openaiKey))
	adapter := &ClientAdapter{client: &client}

	return &Client{httpClient: adapter}, nil
}

// GetEmbedding generates embeddings for the given text
func (c *Client) GetEmbedding(ctx context.Context, text string) ([]float32, error) {
	var embeddingResp struct {
		Data []struct {
			Embedding []float32 `json:"embedding"`
		} `json:"data"`
	}

	err := c.httpClient.Post(ctx, "/embeddings", map[string]interface{}{
		"model": "text-embedding-3-small",
		"input": []string{text},
	}, &embeddingResp)

	if err != nil {
		return nil, fmt.Errorf("failed to get embedding: %w", err)
	}

	if len(embeddingResp.Data) == 0 {
		return nil, fmt.Errorf("no embedding returned")
	}

	return embeddingResp.Data[0].Embedding, nil
}

// GenerateCompletion generates text completions (placeholder for future functionality)
func (c *Client) GenerateCompletion(ctx context.Context, prompt string) (string, error) {
	// This is a placeholder for future completion functionality
	// You can implement this method when needed
	return "", fmt.Errorf("completion functionality not yet implemented")
}

// GenerateImage generates images from text prompts (placeholder for future functionality)
func (c *Client) GenerateImage(ctx context.Context, prompt string) (string, error) {
	// This is a placeholder for future image generation functionality
	// You can implement this method when needed
	return "", fmt.Errorf("image generation functionality not yet implemented")
}
