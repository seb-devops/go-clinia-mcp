package openai

import (
	"context"
	"errors"
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

type OpenAIClientAdapter struct {
	client *openai.Client
}

func (a *OpenAIClientAdapter) Post(ctx context.Context, path string, body interface{}, out interface{}) error {
	return a.client.Post(ctx, path, body, out)
}

func NewOpenAIEmbedder() (*OpenAIEmbedder, error) {
	openaiKey := os.Getenv("OPENAI_API_KEY")
	if openaiKey == "" {
		return nil, errors.New("OPENAI_API_KEY environment variable not set")
	}
	client := openai.NewClient(option.WithAPIKey(openaiKey))
	adapter := &OpenAIClientAdapter{client: &client}
	return &OpenAIEmbedder{client: adapter}, nil
}

func (e *OpenAIEmbedder) GetEmbedding(query string) ([]float32, error) {
	ctx := context.Background()
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
		return nil, err
	}

	if len(embeddingResp.Data) == 0 {
		return nil, errors.New("no embedding returned")
	}

	return embeddingResp.Data[0].Embedding, nil
}
