package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/sourcegraph/jsonrpc2"
)

type MCPServer struct {
	executor *DCTExecutor
	conn     *jsonrpc2.Conn
}

func NewMCPServer(dctPath string) *MCPServer {
	return &MCPServer{
		executor: NewDCTExecutor(dctPath),
	}
}

func (s *MCPServer) Start(ctx context.Context) error {
	handler := jsonrpc2.HandlerWithError(s.handle)

	s.conn = jsonrpc2.NewConn(ctx, jsonrpc2.NewPlainObjectStream(stdinStdoutConn{}), handler)

	<-s.conn.DisconnectNotify()
	return nil
}

type stdinStdoutConn struct{}

func (stdinStdoutConn) Read(p []byte) (int, error) {
	return os.Stdin.Read(p)
}

func (stdinStdoutConn) Write(p []byte) (int, error) {
	return os.Stdout.Write(p)
}

func (stdinStdoutConn) Close() error {
	return nil
}

func (s *MCPServer) handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (any, error) {
	switch req.Method {
	case "initialize":
		return s.handleInitialize()
	case "tools/list":
		return s.handleToolsList()
	case "tools/call":
		return s.handleToolCall(req)
	default:
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeMethodNotFound, Message: fmt.Sprintf("method not found: %s", req.Method)}
	}
}

func (s *MCPServer) handleInitialize() (any, error) {
	return map[string]any{
		"protocolVersion": "2024-11-05",
		"capabilities": map[string]any{
			"tools": map[string]any{},
		},
		"serverInfo": map[string]any{
			"name":    "dct-mcp-server",
			"version": "1.0.0",
		},
	}, nil
}

func (s *MCPServer) handleToolsList() (any, error) {
	tools := []map[string]any{
		{
			"name":        "data_peek",
			"description": "Preview file contents - display the first few lines of a data file",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"file_path": map[string]any{
						"type":        "string",
						"description": "Path to the data file (CSV, JSON, NDJSON, Parquet)",
					},
					"lines": map[string]any{
						"type":        "integer",
						"description": "Number of lines to display (default: 10)",
						"minimum":     1,
					},
					"output_file": map[string]any{
						"type":        "string",
						"description": "Optional output file path",
					},
				},
				"required": []string{"file_path"},
			},
		},
		{
			"name":        "data_diff",
			"description": "Compare two data files with key matching and metrics",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"keys": map[string]any{
						"type":        "string",
						"description": "Key specification for matching (format: left_key[=right_key])",
					},
					"file1": map[string]any{
						"type":        "string",
						"description": "Path to the first data file",
					},
					"file2": map[string]any{
						"type":        "string",
						"description": "Path to the second data file",
					},
					"metrics": map[string]any{
						"type":        "string",
						"description": "JSON metrics specification (e.g., '[{\"agg\":\"count_distinct\",\"left\":\"col\",\"right\":\"col\"}]')",
					},
					"show_all": map[string]any{
						"type":        "boolean",
						"description": "Show all metrics",
					},
					"output_file": map[string]any{
						"type":        "string",
						"description": "Optional output file path",
					},
				},
				"required": []string{"keys", "file1", "file2"},
			},
		},
		{
			"name":        "data_chart",
			"description": "Generate simple charts from data files",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"file_path": map[string]any{
						"type":        "string",
						"description": "Path to the data file",
					},
					"column_index": map[string]any{
						"type":        "integer",
						"description": "Column index to chart",
						"minimum":     0,
					},
					"width": map[string]any{
						"type":        "integer",
						"description": "Width of the chart in characters",
						"minimum":     10,
					},
				},
				"required": []string{"file_path"},
			},
		},
		{
			"name":        "data_generate",
			"description": "Generate synthetic data with customizable schemas",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"schema": map[string]any{
						"type":        "string",
						"description": "JSON schema for data generation or path to schema file",
					},
					"lines": map[string]any{
						"type":        "integer",
						"description": "Number of data rows to generate (default: 1)",
						"minimum":     1,
					},
					"format": map[string]any{
						"type":        "string",
						"description": "Output format: csv or ndjson (default: csv)",
						"enum":        []string{"csv", "ndjson"},
					},
					"output_file": map[string]any{
						"type":        "string",
						"description": "Optional output file path",
					},
				},
				"required": []string{"schema"},
			},
		},
		{
			"name":        "data_flattify",
			"description": "Convert nested JSON structures to flat formats or SQL",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"input": map[string]any{
						"type":        "string",
						"description": "JSON content or path to JSON file",
					},
					"sql": map[string]any{
						"type":        "boolean",
						"description": "Create DuckDB-compliant SQL Select statement",
					},
					"output_file": map[string]any{
						"type":        "string",
						"description": "Optional output file path",
					},
				},
				"required": []string{"input"},
			},
		},
		{
			"name":        "data_js2sql",
			"description": "Convert JSON Schema to SQL CREATE TABLE statements",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"schema_file": map[string]any{
						"type":        "string",
						"description": "Path to JSON schema file",
					},
					"table_name": map[string]any{
						"type":        "string",
						"description": "Table name (default: test)",
					},
					"output_file": map[string]any{
						"type":        "string",
						"description": "Optional output file path",
					},
				},
				"required": []string{"schema_file"},
			},
		},
		{
			"name":        "data_profile",
			"description": "Provide detailed summaries and profiling for data files",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"file_path": map[string]any{
						"type":        "string",
						"description": "Path to the data file to profile",
					},
					"output_file": map[string]any{
						"type":        "string",
						"description": "Optional output file path",
					},
				},
				"required": []string{"file_path"},
			},
		},
	}

	return map[string]any{
		"tools": tools,
	}, nil
}

func (s *MCPServer) handleToolCall(req *jsonrpc2.Request) (any, error) {
	var params struct {
		Name      string         `json:"name"`
		Arguments map[string]any `json:"arguments"`
	}

	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: "invalid parameters"}
	}

	var result *ExecutionResult
	var err error

	switch params.Name {
	case "data_peek":
		result, err = s.handleDataPeek(params.Arguments)
	case "data_diff":
		result, err = s.handleDataDiff(params.Arguments)
	case "data_chart":
		result, err = s.handleDataChart(params.Arguments)
	case "data_generate":
		result, err = s.handleDataGenerate(params.Arguments)
	case "data_flattify":
		result, err = s.handleDataFlattify(params.Arguments)
	case "data_js2sql":
		result, err = s.handleDataJs2Sql(params.Arguments)
	case "data_profile":
		result, err = s.handleDataProfile(params.Arguments)
	default:
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeMethodNotFound, Message: fmt.Sprintf("unknown tool: %s", params.Name)}
	}

	if err != nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInternalError, Message: err.Error()}
	}

	return map[string]any{
		"content": []map[string]any{
			{
				"type": "text",
				"text": formatToolResult(params.Name, result),
			},
		},
	}, nil
}

func (s *MCPServer) handleDataPeek(args map[string]any) (*ExecutionResult, error) {
	filePath, _ := args["file_path"].(string)
	lines, _ := args["lines"].(float64)
	outputFile, _ := args["output_file"].(string)

	if filePath == "" {
		return nil, fmt.Errorf("file_path is required")
	}

	return s.executor.ExecutePeek(filePath, int(lines), outputFile)
}

func (s *MCPServer) handleDataDiff(args map[string]any) (*ExecutionResult, error) {
	keys, _ := args["keys"].(string)
	file1, _ := args["file1"].(string)
	file2, _ := args["file2"].(string)
	metrics, _ := args["metrics"].(string)
	showAll, _ := args["show_all"].(bool)
	outputFile, _ := args["output_file"].(string)

	if keys == "" || file1 == "" || file2 == "" {
		return nil, fmt.Errorf("keys, file1, and file2 are required")
	}

	return s.executor.ExecuteDiff(keys, file1, file2, metrics, showAll, outputFile)
}

func (s *MCPServer) handleDataChart(args map[string]any) (*ExecutionResult, error) {
	filePath, _ := args["file_path"].(string)
	colIndex, _ := args["column_index"].(float64)
	width, _ := args["width"].(float64)

	if filePath == "" {
		return nil, fmt.Errorf("file_path is required")
	}

	return s.executor.ExecuteChart(filePath, int(colIndex), int32(width))
}

func (s *MCPServer) handleDataGenerate(args map[string]any) (*ExecutionResult, error) {
	schema, _ := args["schema"].(string)
	lines, _ := args["lines"].(float64)
	format, _ := args["format"].(string)
	outputFile, _ := args["output_file"].(string)

	if schema == "" {
		return nil, fmt.Errorf("schema is required")
	}

	return s.executor.ExecuteGenerate(schema, int(lines), format, outputFile)
}

func (s *MCPServer) handleDataFlattify(args map[string]any) (*ExecutionResult, error) {
	input, _ := args["input"].(string)
	sql, _ := args["sql"].(bool)
	outputFile, _ := args["output_file"].(string)

	if input == "" {
		return nil, fmt.Errorf("input is required")
	}

	return s.executor.ExecuteFlattify(input, sql, outputFile)
}

func (s *MCPServer) handleDataJs2Sql(args map[string]any) (*ExecutionResult, error) {
	schemaFile, _ := args["schema_file"].(string)
	tableName, _ := args["table_name"].(string)
	outputFile, _ := args["output_file"].(string)

	if schemaFile == "" {
		return nil, fmt.Errorf("schema_file is required")
	}

	return s.executor.ExecuteJs2Sql(schemaFile, tableName, outputFile)
}

func (s *MCPServer) handleDataProfile(args map[string]any) (*ExecutionResult, error) {
	filePath, _ := args["file_path"].(string)
	outputFile, _ := args["output_file"].(string)

	if filePath == "" {
		return nil, fmt.Errorf("file_path is required")
	}

	return s.executor.ExecuteProfile(filePath, outputFile)
}

func formatToolResult(toolName string, result *ExecutionResult) string {
	if !result.Success {
		return fmt.Sprintf("❌ %s failed (exit code: %d, duration: %s)\nError: %s\nOutput: %s",
			toolName, result.ExitCode, result.Duration, result.Error, result.Output)
	}

	return fmt.Sprintf("✅ %s completed successfully (duration: %s)\n\n%s",
		toolName, result.Duration, result.Output)
}
