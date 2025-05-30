package supabase

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"clinia-doc/mcp/services"

	"github.com/supabase-community/supabase-go"
)

type SupabaseClient interface {
	Rpc(fn, mode string, params map[string]interface{}) string
}

// Client implements services.VectorSearchProvider interface
type Client struct {
	supabase SupabaseClient
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

// SearchByVector implements services.VectorSearchProvider interface
func (c *Client) SearchByVector(ctx context.Context, vector []float32, limit int) ([]services.SearchResult, error) {
	var results []services.SearchResult

	// Supabase Rpc returns a string, not an error
	resStr := c.supabase.Rpc("match_site_pages", "exact", map[string]interface{}{
		"query_embedding": vector,
		"match_count":     limit,
	})

	if resStr == "" {
		return nil, fmt.Errorf("supabase rpc 'match_site_pages' returned empty response")
	}

	// Unmarshal the result string into results
	if err := json.Unmarshal([]byte(resStr), &results); err != nil {
		return nil, err
	}

	return results, nil
}

// Legacy method for backward compatibility - DEPRECATED: Use services.SearchService instead
func (c *Client) GetEmbedding(embedding []float32, matchCount int) ([]services.SearchResult, error) {
	return c.SearchByVector(context.Background(), embedding, matchCount)
}
