---
name: dct
description: Router skill for DCT (Data Check Tool). Use this skill whenever the user wants to work with flat data files (CSV, JSON, NDJSON, Parquet) for inspection, comparison, transformation, or generation. The main dct skill analyzes user intent and routes to appropriate sub-skills. Triggers include any mention of data files, previewing data, comparing datasets, generating test data, flattening JSON, creating SQL schemas, profiling data, or visualizing distributions.
tools:
  - dct-peek
  - dct-infer
  - dct-diff
  - dct-generate
  - dct-flattify
  - dct-profile
  - dct-js2sql
  - dct-chart
---

# DCT (Data Check Tool) - Skill Router

DCT is a Swiss army knife CLI tool for working with flat data files. This skill routes to appropriate sub-skills based on user intent.

## Quick Command Reference

| User Intent | Route To | Command Pattern |
|-------------|----------|-----------------|
| Preview/inspect data | `dct-peek` | `dct peek <file>` |
| Generate SQL schema | `dct-infer` | `dct infer <file>` |
| Compare two datasets | `dct-diff` | `dct diff <keys> <file1> <file2>` |
| Generate synthetic data | `dct-generate` | `dct gen <schema>` |
| Flatten nested JSON | `dct-flattify` | `dct flattify <json>` |
| Analyze data quality | `dct-profile` | `dct prof <file>` |
| JSON Schema to SQL | `dct-js2sql` | `dct js2sql <schema>` |
| Visualize data | `dct-chart` | `dct chart <file> <col>` |

## Routing Logic

Analyze the user's request and route to the appropriate sub-skill:

### Route to `dct-peek` when:
- User wants to preview data files
- Keywords: "peek", "preview", "show me", "look at", "first rows", "sample"
- Example: "Show me the first 10 rows of data.csv"

### Route to `dct-infer` when:
- User wants to generate SQL schemas
- Keywords: "infer", "schema", "create table", "sql from data", "ddl"
- Example: "Generate a CREATE TABLE statement from this CSV"

### Route to `dct-diff` when:
- User wants to compare two files
- Keywords: "diff", "compare", "differences", "match", "reconcile", "validate"
- Example: "Compare these two CSV files by the ID column"

### Route to `dct-generate` when:
- User wants to create synthetic test data
- Keywords: "generate", "synthetic", "mock", "fake data", "test data"
- Example: "Generate 1000 fake user records"

### Route to `dct-flattify` when:
- User wants to flatten nested JSON
- Keywords: "flatten", "unnest", "nested json", "make flat"
- Example: "Flatten this nested JSON from the API response"

### Route to `dct-profile` when:
- User wants to analyze data quality
- Keywords: "profile", "analyze", "data quality", "statistics", "distribution"
- Example: "Profile this data file for quality issues"

### Route to `dct-js2sql` when:
- User wants to convert JSON Schema to SQL
- Keywords: "json schema", "convert schema", "schema to sql"
- Example: "Convert this JSON Schema to a CREATE TABLE statement"

### Route to `dct-chart` when:
- User wants to visualize data
- Keywords: "chart", "visualize", "histogram", "plot", "graph"
- Example: "Create a chart of the sales column"

## Common Patterns

### Data Validation Workflow
1. `dct-peek`: Preview to understand structure
2. `dct-profile`: Check data quality
3. `dct-infer`: Generate schema for downstream use

### Data Comparison Workflow
1. `dct-peek`: Preview both files
2. `dct-diff`: Compare with appropriate keys

### Test Data Generation Workflow
1. `dct-generate`: Create synthetic data
2. `dct-peek`: Verify generated data
3. `dct-diff`: Compare with production sample

## Installation

All sub-skills require DCT to be installed:

```bash
which dct || go build -o dct && chmod +x ./dct
```

## Supported File Formats

All DCT sub-skills support:
- CSV (.csv)
- JSON (.json)
- NDJSON (.ndjson) - newline-delimited JSON
- Parquet (.parquet)

## Error Handling

If a sub-skill encounters errors:
- Verify the file exists and is readable
- Check file extension matches content format
- Ensure DCT binary is built and executable
