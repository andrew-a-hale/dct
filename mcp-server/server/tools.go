package server

func getTools() []Tool {
	return []Tool{
		{
			Name:        "data_peek",
			Description: "Preview file contents - display the first few lines of a data file",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"file_path": {
						Type:        "string",
						Description: "Path to the data file (CSV, JSON, NDJSON, Parquet)",
					},
					"lines": {
						Type:        "integer",
						Description: "Number of lines to display (default: 10)",
						Minimum:     &[]int{1}[0],
					},
					"output_file": {
						Type:        "string",
						Description: "Optional output file path",
					},
				},
				Required: []string{"file_path"},
			},
		},
		{
			Name:        "data_infer",
			Description: "Generate a SQL Create Table statement from a file",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"file_path": {
						Type:        "string",
						Description: "Path to the data file (CSV, JSON, NDJSON, Parquet)",
					},
					"lines": {
						Type:        "integer",
						Description: "Number of lines to display (default: 10)",
						Minimum:     &[]int{1}[0],
					},
					"output_file": {
						Type:        "string",
						Description: "Optional output file path",
					},
					"table": {
						Type:        "string",
						Description: "Table name used in create table statement (default: default)",
					},
				},
				Required: []string{"file_path"},
			},
		},
		{
			Name:        "data_diff",
			Description: "Compare two data files with key matching and metrics",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"keys": {
						Type:        "string",
						Description: "Key specification for matching (format: left_key[=right_key])",
					},
					"file1": {
						Type:        "string",
						Description: "Path to the first data file",
					},
					"file2": {
						Type:        "string",
						Description: "Path to the second data file",
					},
					"metrics": {
						Type:        "string",
						Description: "JSON metrics specification (e.g., '[{\"agg\":\"count_distinct\",\"left\":\"col\",\"right\":\"col\"}]')",
					},
					"show_all": {
						Type:        "boolean",
						Description: "Show all metrics",
					},
					"output_file": {
						Type:        "string",
						Description: "Optional output file path",
					},
				},
				Required: []string{"keys", "file1", "file2"},
			},
		},
		{
			Name:        "data_chart",
			Description: "Generate simple charts from data files",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"file_path": {
						Type:        "string",
						Description: "Path to the data file",
					},
					"column_index": {
						Type:        "integer",
						Description: "Column index to chart",
						Minimum:     &[]int{0}[0],
					},
					"width": {
						Type:        "integer",
						Description: "Width of the chart in characters",
						Minimum:     &[]int{10}[0],
					},
				},
				Required: []string{"file_path"},
			},
		},
		{
			Name:        "data_generate",
			Description: "Generate synthetic data with customizable schemas",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"schema": {
						Type: "string",
						Description: `JSON schema for data generation or path to schema file.
The schema should be a JSON array of field objects, each containing:
{"field": "column_name", "source": "source_type", "config": {...}}.
Available sources:
  randomBool,
	randomEnum (config: {"values": array}),
  randomAscii (config: {"length": int}),
  randomUniformInt (config: {"min": int, "max": int}),
  randomNormal (config: {"mean": float, "std": float}),
  randomPoisson (config: {"lambda": int}),
  randomDatetime (config: {"tz": "timezone", "min": "YYYY-MM-DD HH:MM:SS", "max": "YYYY-MM-DD HH:MM:SS"}),
  randomDate (config: {"min": "YYYY-MM-DD", "max": "YYYY-MM-DD"}),
  randomTime (config: {"min": "HH:MM:SS", "max": "HH:MM:SS"}),
  uuid,
  firstNames,
  lastNames,
  companies,
  emails,
  derived (config: {"fields": ["field1", "field2"], "expression": "field1 + ' ' + field2"}).

Example: '[{"field":"name","source":"firstNames"},{"field":"age","source":"randomUniformInt","config":{"min":18,"max":65}}]'`,
					},
					"lines": {
						Type:        "integer",
						Description: "Number of data rows to generate (default: 1)",
						Minimum:     &[]int{1}[0],
					},
					"format": {
						Type:        "string",
						Description: "Output format: csv or ndjson (default: csv)",
						Enum:        []string{"csv", "ndjson"},
					},
					"output_file": {
						Type:        "string",
						Description: "Optional output file path",
					},
				},
				Required: []string{"schema"},
			},
		},
		{
			Name:        "data_flattify",
			Description: "Convert nested JSON structures to flat formats or SQL",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"input": {
						Type:        "string",
						Description: "JSON content or path to JSON file",
					},
					"sql": {
						Type:        "boolean",
						Description: "Create DuckDB-compliant SQL Select statement",
					},
					"output_file": {
						Type:        "string",
						Description: "Optional output file path",
					},
				},
				Required: []string{"input"},
			},
		},
		{
			Name:        "data_js2sql",
			Description: "Convert JSON Schema to SQL CREATE TABLE statements",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"schema_file": {
						Type:        "string",
						Description: "Path to JSON schema file",
					},
					"table_name": {
						Type:        "string",
						Description: "Table name (default: test)",
					},
					"output_file": {
						Type:        "string",
						Description: "Optional output file path",
					},
				},
				Required: []string{"schema_file"},
			},
		},
		{
			Name:        "data_profile",
			Description: "Provide detailed summaries and profiling for data files",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"file_path": {
						Type:        "string",
						Description: "Path to the data file to profile",
					},
					"output_file": {
						Type:        "string",
						Description: "Optional output file path",
					},
				},
				Required: []string{"file_path"},
			},
		},
	}
}
