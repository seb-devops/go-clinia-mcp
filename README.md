# MCP Clinia Doc

A Go 1.22+ project demonstrating how to build an MCP (Model Context Protocol) server with custom tools, including integration with Supabase, OpenAI, and LangChain for document embedding, retrieval, and query decomposition.

## Features
- MCP server using [mcp-go](https://github.com/mark3labs/mcp-go)
- Custom tool: `get_embedding` (calls Supabase RPC with OpenAI embeddings)
- Custom tool: `decompose_query` (breaks down complex queries using LangChain agents)
- Uses Go's new `ServeMux` and idiomatic net/http patterns
- SSE (Server-Sent Events) transport for real-time communication

## Requirements
- Go 1.24 or newer
- Supabase project (for document matching)
- OpenAI API key (for embeddings and LangChain LLM)

## Setup
1. **Clone the repository:**
   ```sh
   git clone <your-repo-url>
   cd <your-repo-directory>
   ```
2. **Set environment variables:**
   - `SUPABASE_URL`: Your Supabase project URL
   - `SUPABASE_SERVICE_KEY`: Your Supabase service role key
   - `OPENAI_API_KEY`: Your OpenAI API key

   Example (on Linux/macOS):
   ```sh
   export SUPABASE_URL="https://your-project.supabase.co"
   export SUPABASE_SERVICE_KEY="your-service-role-key"
   export OPENAI_API_KEY="sk-..."
   ```

3. **Install dependencies:**
   ```sh
   go mod tidy
   ```

4. **Run the server:**
   ```sh
   go run ./mcp/mcp.go
   ```

## Available Tools

### `get_embedding`
Returns document matches from Supabase using OpenAI embeddings.

**Parameters:**
- `query` (string, required): Text to embed and search

**Example:**
```json
{
  "tool": "get_embedding",
  "arguments": {
    "query": "How to implement authentication in Go"
  }
}
```

### `decompose_query`
Breaks down complex queries into smaller, actionable subqueries using LangChain agents.

**Parameters:**
- `query` (string, required): Complex query to decompose

**Example:**
```json
{
  "tool": "decompose_query",
  "arguments": {
    "query": "I need to understand the impact of artificial intelligence on healthcare systems, including current applications, challenges, costs, and future trends in the next 5 years."
  }
}
```

**Response:**
```json
{
  "original_query": "I need to understand the impact of artificial intelligence...",
  "sub_queries": [
    {
      "id": 1,
      "query": "What are the current applications of AI in healthcare systems?",
      "description": "Current Applications",
      "priority": 1
    }
  ],
  "strategy": "agent_based_decomposition"
}
```

## Adding a New Tool
- Define your tool using `mcp.NewTool`.
- Implement a handler: `func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error)`
- Register the tool and handler with `s.AddTool(tool, handler)` in `setupServerAndTools()`.

## Project Structure
```
├── docs/
│   └── technical-decisions.md # Technical decisions and architecture choices
├── mcp/
│   ├── mcp.go              # Main server entrypoint
│   ├── supabase/
│   │   ├── client_impl.go  # Supabase client implementation
│   │   └── client_test.go  # Supabase client tests
│   ├── openai/
│   │   ├── openai_embed.go # OpenAI embedding client
│   │   └── openai_test.go  # OpenAI client tests
│   └── langchain/
│       ├── decomposer.go   # LangChain query decomposition
│       ├── decomposer_test.go # Query decomposition tests
│       └── README.md       # LangChain tool documentation
├── go.mod
├── go.sum
└── README.md
```

## Documentation

- **[Technical Decisions](docs/technical-decisions.md)**: Detailed documentation of architectural choices and implementation decisions for the LangChain integration
- **[LangChain Tool Guide](mcp/langchain/README.md)**: Comprehensive guide for the query decomposition tool

## Testing

Run all tests:
```sh
go test ./...
```

Run specific package tests:
```sh
# Test LangChain functionality (requires OPENAI_API_KEY)
go test ./mcp/langchain -v

# Test Supabase functionality
go test ./mcp/supabase -v

# Test OpenAI functionality
go test ./mcp/openai -v
```

## Dependencies

- **MCP**: `github.com/mark3labs/mcp-go` - Model Context Protocol implementation
- **Supabase**: `github.com/supabase-community/supabase-go` - Supabase client
- **OpenAI**: `github.com/openai/openai-go` - OpenAI API client
- **LangChain**: `github.com/tmc/langchaingo` - LangChain Go implementation
- **Logging**: `github.com/sirupsen/logrus` - Structured logging
- **Environment**: `github.com/joho/godotenv` - Environment variable loading

## Contributing
Pull requests and issues are welcome! Please open an issue to discuss your idea or bug before submitting a PR.

## License
MIT 