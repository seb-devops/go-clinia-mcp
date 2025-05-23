package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"clinia-doc/mcp/openai"
	"clinia-doc/mcp/supabase"

	"github.com/joho/godotenv"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Load .env if present
	_ = godotenv.Load()

	// Optionally configure logrus (e.g., set format)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(log.InfoLevel)

	log.Println("Starting MCP server...")

	sseServer, err := setupServerAndTools()
	if err != nil {
		log.Fatalf("Server setup failed: %v", err)
	}

	log.Printf("SSE server listening on :8080")
	if err := sseServer.Start(":8080"); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// setupServerAndTools creates the MCP server, registers tools, and returns the SSE server.
func setupServerAndTools() (*server.SSEServer, error) {
	// Create a new MCP server with logging and recovery middleware
	s := server.NewMCPServer(
		"Clinia Doc Demo 🚀",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithRecovery(),
	)

	// Supabase client
	supabaseClient, err := supabase.NewClient()
	if err != nil {
		return nil, fmt.Errorf("could not initialize supabase client: %w", err)
	}

	openaiClient, err := openai.NewOpenAIEmbedder()
	if err != nil {
		return nil, fmt.Errorf("could not initialize openai client: %w", err)
	}

	// Register your get_embedding tool
	getEmbeddingTool := mcp.NewTool("get_embedding",
		mcp.WithDescription("Get document matches from Supabase using OpenAI embeddings"),
		mcp.WithString("query", mcp.Required(), mcp.Description("Text to embed and search")),
	)

	s.AddTool(getEmbeddingTool, getEmbeddingHandler(supabaseClient, openaiClient))

	// Start the SSE server on :8080 using HTTP handler
	sseServer := server.NewSSEServer(s, server.WithSSEEndpoint("/api/sse"))
	return sseServer, nil
}

func getEmbeddingHandler(client *supabase.Client, openaiClient *openai.OpenAIEmbedder) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("Received get_embedding request with params: %v", request.Params.Arguments)
		query, ok := request.Params.Arguments["query"].(string)
		if !ok || query == "" {
			return nil, errors.New("query must be a non-empty string")
		}

		embedding, err := openaiClient.GetEmbedding(query)
		if err != nil {
			log.Printf("OpenAI embedding error: %v", err)
			return nil, fmt.Errorf("failed to get embedding from OpenAI: %w", err)
		}

		results, err := client.GetEmbedding(embedding, 10)

		if err != nil {
			log.Printf("GetEmbedding error: %v", err)
			return nil, fmt.Errorf("failed to get embedding: %w", err)
		}

		jsonBytes, err := json.Marshal(results)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal results: %w", err)
		}

		return mcp.NewToolResultText(string(jsonBytes)), nil
	}
}
