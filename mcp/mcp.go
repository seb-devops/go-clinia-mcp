package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"clinia-doc/mcp/langchain"
	"clinia-doc/mcp/openai"
	"clinia-doc/mcp/services"
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

	// Initialize providers
	openaiClient, err := openai.NewClient()
	if err != nil {
		return nil, fmt.Errorf("could not initialize openai client: %w", err)
	}

	supabaseClient, err := supabase.NewClient()
	if err != nil {
		return nil, fmt.Errorf("could not initialize supabase client: %w", err)
	}

	// Initialize service layer with dependency injection
	searchService := services.NewSearchService(openaiClient, supabaseClient)

	// LangChain decomposer
	decomposer, err := langchain.NewQueryDecomposer()
	if err != nil {
		return nil, fmt.Errorf("could not initialize langchain decomposer: %w", err)
	}

	// Register your get_embedding tool
	getEmbeddingTool := mcp.NewTool("get_embedding",
		mcp.WithDescription("Get document matches from Supabase using OpenAI embeddings"),
		mcp.WithString("query", mcp.Required(), mcp.Description("Text to embed and search")),
	)

	decomposeQueryTool := mcp.NewTool("decompose_query",
		mcp.WithDescription("Decompose a complex query into multiple subqueries using LangChain"),
		mcp.WithString("query", mcp.Required(), mcp.Description("The complex query to decompose")),
	)

	s.AddTool(decomposeQueryTool, decomposeQueryHandler(decomposer))
	s.AddTool(getEmbeddingTool, getEmbeddingHandler(searchService))

	// Start the SSE server on :8080 using HTTP handler
	sseServer := server.NewSSEServer(s, server.WithSSEEndpoint("/api/sse"))
	return sseServer, nil
}

func getEmbeddingHandler(searchService *services.SearchService) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("Received get_embedding request with params: %v", request.Params.Arguments)
		query, ok := request.Params.Arguments["query"].(string)
		if !ok || query == "" {
			return nil, errors.New("query must be a non-empty string")
		}

		results, err := searchService.SearchByText(ctx, query, 10)
		if err != nil {
			log.Printf("Search error: %v", err)
			return nil, fmt.Errorf("failed to search: %w", err)
		}

		jsonBytes, err := json.Marshal(results)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal results: %w", err)
		}

		return mcp.NewToolResultText(string(jsonBytes)), nil
	}
}

func decomposeQueryHandler(decomposer *langchain.QueryDecomposer) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("Received decompose_query request with params: %v", request.Params.Arguments)
		query, ok := request.Params.Arguments["query"].(string)
		if !ok || query == "" {
			return nil, errors.New("query must be a non-empty string")
		}

		results, err := decomposer.DecomposeQuery(ctx, query)
		if err != nil {
			log.Printf("Query decomposition error: %v", err)
			return nil, fmt.Errorf("failed to decompose query: %w", err)
		}

		jsonBytes, err := json.Marshal(results)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal results: %w", err)
		}

		return mcp.NewToolResultText(string(jsonBytes)), nil
	}
}
