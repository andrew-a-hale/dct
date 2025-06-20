# DCT MCP Server

This is an MCP (Model Context Protocol) server that wraps the DCT (Data Check
Tool) CLI, making it available to AI agents and other MCP clients.

## Features

The MCP server exposes all DCT functionality as MCP tools:

- **data_peek** - Preview file contents
- **data_infer** - Infer a SQL Schema for a file
- **data_diff** - Compare files with key matching and metrics  
- **data_chart** - Generate simple visualizations from data files
- **data_generate** - Generate synthetic data with customizable schemas
- **data_flattify** - Convert nested JSON structures to flat formats or SQL
- **data_js2sql** - Convert JSON Schema to SQL CREATE TABLE statements
- **data_profile** - Profile data files for values and characters

## Installation

1. Build the MCP server:

```bash
cd mcp-server
go build -o dct-mcp-server
```

2. Run the MCP server (requires path to DCT binary):

```bash
./dct-mcp-server ../dct
```

## Usage with MCP Clients

The server communicates via JSON-RPC 2.0 over stdin/stdout. Here's an example of how to use it:

### Initialize

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "protocolVersion": "2024-11-05",
    "capabilities": {},
    "clientInfo": {"name": "test-client", "version": "1.0.0"}
  }
}
```

### List Tools

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/list",
  "params": {}
}
```

### Call a Tool

```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "data_peek",
    "arguments": {
      "file_path": "examples/left.parquet",
      "lines": 5
    }
  }
}
```

## Integration with Claude Desktop

Add this to your Claude Desktop MCP configuration:

```json
{
  "mcpServers": {
    "dct": {
      "command": "/path/to/dct-mcp-server",
      "args": ["/path/to/dct"]
    }
  }
}
```

## Development

The server is built with Go and uses the `github.com/sourcegraph/jsonrpc2`
library for MCP communication. It wraps the existing DCT CLI commands,
preserving all functionality while making it accessible to AI agents.

### Architecture

- `main.go` - Entry point and server initialization
- `server.go` - MCP protocol handlers and tool definitions
- `executor.go` - CLI command execution wrapper

The server creates temporary files when needed for JSON input and cleans them up automatically.
