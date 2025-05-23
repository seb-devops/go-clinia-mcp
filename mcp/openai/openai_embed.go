package openai

import (
	"context"
	"fmt"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type EmbedderClient interface {
	Post(ctx context.Context, path string, body interface{}, out interface{}) error
}

type OpenAIEmbedder struct {
	client EmbedderClient
}

// Used for testing purposes
type OpenAIClientAdapter struct {
	client *openai.Client
}

func (a *OpenAIClientAdapter) Post(ctx context.Context, path string, body interface{}, out interface{}) error {
	return a.client.Post(ctx, path, body, out)
}

func NewOpenAIEmbedder() (*OpenAIEmbedder, error) {
	openaiKey := os.Getenv("OPENAI_API_KEY")
	if openaiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}
	client := openai.NewClient(option.WithAPIKey(openaiKey))
	adapter := &OpenAIClientAdapter{client: &client}
	return &OpenAIEmbedder{client: adapter}, nil
}

func (e *OpenAIEmbedder) GetEmbedding(ctx context.Context, query string) ([]float32, error) {
	var embeddingResp struct {
		Data []struct {
			Embedding []float32 `json:"embedding"`
		} `json:"data"`
	}

	err := e.client.Post(ctx, "/embeddings", map[string]interface{}{
		"model": "text-embedding-3-small",
		"input": []string{query},
	}, &embeddingResp)

	if err != nil {
		return nil, fmt.Errorf("failed to get embedding: %w", err)
	}

	if len(embeddingResp.Data) == 0 {
		return nil, fmt.Errorf("no embedding returned")
	}

	return embeddingResp.Data[0].Embedding, nil
}
