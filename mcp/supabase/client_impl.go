package supabase

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/supabase-community/supabase-go"

	"clinia-doc/mcp/openai"
)

type SupabaseClient interface {
	Rpc(fn, mode string, params map[string]interface{}) string
}

type Embedder interface {
	GetEmbedding(query string) ([]float32, error)
}

type Client struct {
	supabase SupabaseClient
	embedder openai.OpenAIEmbedder
}

type MatchResult struct {
	ID          int                    `json:"id"`
	URL         string                 `json:"url"`
	ChunkNumber int                    `json:"chunk_number"`
	Title       string                 `json:"title"`
	Summary     string                 `json:"summary"`
	Content     string                 `json:"content"`
	Metadata    map[string]interface{} `json:"metadata"`
	Similarity  float64                `json:"similarity"`
}

type SupabaseClientAdapter struct {
	client *supabase.Client
}

func (a *SupabaseClientAdapter) Rpc(fn, mode string, params map[string]interface{}) string {
	return a.client.Rpc(fn, mode, params)
}

func NewClient() (*Client, error) {
	supabaseUrl := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_SERVICE_KEY")
	supabaseClient, err := supabase.NewClient(supabaseUrl, supabaseKey, &supabase.ClientOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create supabase client: %w", err)
	}

	return &Client{
		supabase: &SupabaseClientAdapter{client: supabaseClient},
	}, nil
}

func (c *Client) GetEmbedding(embedding []float32, matchCount int) ([]MatchResult, error) {
	var results []MatchResult

	// Supabase Rpc returns a string, not an error
	resStr := c.supabase.Rpc("match_site_pages", "exact", map[string]interface{}{
		"query_embedding": embedding,
		"match_count":     matchCount,
	})

	if resStr == "" {
		return nil, errors.New("no results returned")
	}

	// Unmarshal the result string into results
	if err := json.Unmarshal([]byte(resStr), &results); err != nil {
		return nil, err
	}

	return results, nil
}
