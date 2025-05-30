package langchain

import (
	"context"
	"fmt"
	"strings"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/tools"
)

// QueryDecomposer wraps LangchainGo agent functionality for breaking down queries
type QueryDecomposer struct {
	llm   *openai.LLM
	agent *agents.OneShotZeroAgent
}

// SubQuery represents a decomposed part of the original query
type SubQuery struct {
	ID          int    `json:"id"`
	Query       string `json:"query"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
}

// DecompositionResult represents the result of query decomposition
type DecompositionResult struct {
	OriginalQuery string     `json:"original_query"`
	SubQueries    []SubQuery `json:"sub_queries"`
	Strategy      string     `json:"strategy"`
}

// NewQueryDecomposer creates a new instance of QueryDecomposer
func NewQueryDecomposer() (*QueryDecomposer, error) {
	// Initialize OpenAI LLM
	llm, err := openai.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize OpenAI LLM: %w", err)
	}

	// Create a simple tool for query analysis
	analysisTools := []tools.Tool{
		&QueryAnalysisTool{},
	}

	// Create agent for query decomposition
	agent := agents.NewOneShotAgent(llm,
		analysisTools,
		agents.WithMaxIterations(3),
	)

	return &QueryDecomposer{
		llm:   llm,
		agent: agent,
	}, nil
}

// DecomposeQuery breaks down a complex query into smaller, actionable subqueries
func (qd *QueryDecomposer) DecomposeQuery(ctx context.Context, query string) (*DecompositionResult, error) {
	if strings.TrimSpace(query) == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}

	// Create a specialized prompt for query decomposition
	decompositionPrompt := fmt.Sprintf(`
		You are a query decomposition expert. Given the following complex query, break it down into 3-5 smaller, specific subqueries that can be answered independently.

		Original Query: "%s"

		Please decompose this into subqueries following this format:
		1. [Brief description] - [Specific query]
		2. [Brief description] - [Specific query]
		...

		Focus on:
		- Making each subquery specific and actionable
		- Ensuring subqueries can be answered independently
		- Ordering them by logical priority
		- Keeping the total number between 3-5 subqueries

		Return only the numbered list of subqueries with their descriptions.
	`, query)

	// Execute the agent with the decomposition prompt
	executor := agents.NewExecutor(qd.agent)
	response, err := chains.Run(ctx, executor, decompositionPrompt)
	if err != nil {
		return nil, fmt.Errorf("failed to decompose query: %w", err)
	}

	// Parse the response to extract subqueries
	subQueries, err := qd.parseDecompositionResponse(response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse decomposition response: %w", err)
	}

	return &DecompositionResult{
		OriginalQuery: query,
		SubQueries:    subQueries,
		Strategy:      "agent_based_decomposition",
	}, nil
}

// parseDecompositionResponse parses the agent's response into structured subqueries
func (qd *QueryDecomposer) parseDecompositionResponse(response string) ([]SubQuery, error) {
	lines := strings.Split(strings.TrimSpace(response), "\n")
	var subQueries []SubQuery

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Look for numbered list items (1., 2., etc.)
		if len(line) > 2 && line[1] == '.' {
			// Extract the content after the number
			content := strings.TrimSpace(line[2:])

			// Try to split on " - " to separate description and query
			parts := strings.Split(content, " - ")
			var description, query string

			if len(parts) >= 2 {
				description = strings.TrimSpace(parts[0])
				query = strings.TrimSpace(strings.Join(parts[1:], " - "))
			} else {
				// If no separator found, use the whole content as query
				description = fmt.Sprintf("Subquery %d", len(subQueries)+1)
				query = content
			}

			subQueries = append(subQueries, SubQuery{
				ID:          len(subQueries) + 1,
				Query:       query,
				Description: description,
				Priority:    i + 1,
			})
		}
	}

	if len(subQueries) == 0 {
		return nil, fmt.Errorf("no subqueries found in response: %s", response)
	}

	return subQueries, nil
}

// QueryAnalysisTool is a simple tool for query analysis
type QueryAnalysisTool struct{}

func (t *QueryAnalysisTool) Name() string {
	return "query_analysis"
}

func (t *QueryAnalysisTool) Description() string {
	return "Analyzes and decomposes complex queries into smaller parts"
}

func (t *QueryAnalysisTool) Call(ctx context.Context, input string) (string, error) {
	// This is a placeholder tool that helps the agent understand it can analyze queries
	return fmt.Sprintf("Query analyzed: %s", input), nil
}
