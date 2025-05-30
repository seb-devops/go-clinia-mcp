package langchain

import (
	"os"
	"testing"
)

// testCase represents a test case for query decomposition
type testCase struct {
	name  string
	query string
	want  int // expected minimum number of subqueries
}

// setupDecomposer creates a QueryDecomposer for testing, skipping if no API key
func setupDecomposer(t *testing.T) *QueryDecomposer {
	t.Helper()

	// Skip if no OpenAI API key is available
	if os.Getenv("OPENAI_API_KEY") == "" {
		t.Skip("OPENAI_API_KEY not set, skipping test")
	}

	decomposer, err := NewQueryDecomposer()
	if err != nil {
		t.Fatalf("Failed to create QueryDecomposer: %v", err)
	}

	return decomposer
}

// validateDecompositionResult checks the basic structure of a decomposition result
func validateDecompositionResult(t *testing.T, result *DecompositionResult, originalQuery string, minSubQueries int) {
	t.Helper()

	if result == nil {
		t.Error("DecomposeQuery() returned nil result")
		return
	}

	if result.OriginalQuery != originalQuery {
		t.Errorf("DecomposeQuery() original query = %v, want %v", result.OriginalQuery, originalQuery)
	}

	if len(result.SubQueries) < minSubQueries {
		t.Errorf("DecomposeQuery() got %d subqueries, want at least %d", len(result.SubQueries), minSubQueries)
	}
}

// validateSubQueries checks that all subqueries have required fields
func validateSubQueries(t *testing.T, subQueries []SubQuery) {
	t.Helper()

	for i, subquery := range subQueries {
		if subquery.ID == 0 {
			t.Errorf("Subquery %d has empty ID", i)
		}
		if subquery.Query == "" {
			t.Errorf("Subquery %d has empty Query", i)
		}
		if subquery.Description == "" {
			t.Errorf("Subquery %d has empty Description", i)
		}
	}
}

// logSubQueries logs the decomposed subqueries for debugging
func logSubQueries(t *testing.T, subQueries []SubQuery) {
	t.Helper()

	t.Logf("Successfully decomposed query into %d subqueries:", len(subQueries))
	for _, sq := range subQueries {
		t.Logf("  %d. %s - %s", sq.ID, sq.Description, sq.Query)
	}
}

// getTestCases returns the test cases for query decomposition
func getTestCases() []testCase {
	return []testCase{
		{
			name:  "Simple query",
			query: "What is the weather today?",
			want:  1,
		},
		{
			name:  "Complex research query",
			query: "I need to understand the impact of artificial intelligence on healthcare systems, including current applications, challenges, costs, and future trends in the next 5 years.",
			want:  3,
		},
		{
			name:  "Multi-faceted business query",
			query: "How can a small e-commerce business improve its online presence, increase customer retention, and optimize its supply chain management?",
			want:  3,
		},
	}
}
