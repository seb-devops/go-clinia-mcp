# MCP Clinia Doc

A Go 1.22+ project demonstrating how to build an MCP (Model Context Protocol) server with custom tools, including integration with Supabase and OpenAI for document embedding and retrieval.

## Features
- MCP server using [mcp-go](https://github.com/mark3labs/mcp-go)
- Custom tool: `hello_world` (greets a user)
- Custom tool: `get_embedding` (calls Supabase RPC with OpenAI embeddings)
- Uses Go's new `ServeMux` and idiomatic net/http patterns

## Requirements
- Go 1.22 or newer
- Supabase project (for document matching)
- OpenAI API key (for embeddings)

## Setup
1. **Clone the repository:**
   ```sh
   git clone <your-repo-url>
   cd <your-repo-directory>
   ```
2. **Set environment variables:**
   - `SUPABASE_URL`: Your Supabase project URL
   - `SUPABASE_KEY`: Your Supabase service role key
   - `OPENAI_API_KEY`: Your OpenAI API key

   Example (on Linux/macOS):
   ```sh
   export SUPABASE_URL="https://your-project.supabase.co"
   export SUPABASE_KEY="your-service-role-key"
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

## Adding a New Tool
- Define your tool using `mcp.NewTool`.
- Implement a handler: `func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error)`
- Register the tool and handler with `s.AddTool(tool, handler)` in `main()`.

## Example Usage
- `hello_world` tool: Greets a user by name.
- `get_embedding` tool: Returns document matches from Supabase using OpenAI embeddings.

## Project Structure
```
├── mcp/
│   ├── mcp.go           # Main server entrypoint
│   └── supabase/
│       └── client.go    # Supabase + OpenAI client logic
├── go.mod
├── go.sum
└── README.md
```

## Contributing
Pull requests and issues are welcome! Please open an issue to discuss your idea or bug before submitting a PR.

## License
MIT 