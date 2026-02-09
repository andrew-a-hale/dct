---
name: dct-flattify
description: Use this skill when the user wants to flatten nested JSON structures, convert nested objects to flat format, generate SQL queries from nested JSON, unnest hierarchical data, or work with nested API responses that need to be tabular. Triggers include "flatten this json", "make json flat", "nested to flat", "unnest json", "json to sql", "flatten nested", or when dealing with deeply nested JSON from APIs or document stores.
---

# DCT Flattify - Flatten Nested JSON

Convert nested JSON structures to flat formats or SQL SELECT statements.

## When to Use

Use this skill when you need to:
- Convert nested API responses to flat format
- Transform hierarchical JSON for tabular analysis
- Generate SQL to query nested JSON files
- Unnest deeply nested structures
- Prepare JSON data for CSV export

## Installation

```bash
which dct || go build -o dct && chmod +x ./dct
```

## Usage

```bash
dct flattify <json> [flags]
```

## Arguments

- `json`: JSON file path or inline JSON string

## Flags

- `-s, --sql`: Generate DuckDB SQL SELECT statement
- `-o, --output <file>`: Output to file instead of stdout

## Examples

### Flatten JSON File

Basic flattening:
```bash
dct flattify nested.json -o flat.json
```

### Generate SQL Query

Create SQL to query the JSON:
```bash
dct flattify nested.json -s -o query.sql
```

### Inline JSON

Flatten inline JSON:
```bash
dct flattify '{"user":{"name":"John","age":30}}'
```

Flatten JSON array:
```bash
dct flattify '[{"a":1},{"b":2}]'
```

### Complex Nested Structure

```bash
dct flattify '{"data":{"users":[{"id":1,"profile":{"name":"John"}}]}}' -s
```

## Output Formats

### Without -s (Flat JSON)

Converts nested structure to flat key-value pairs using JSONPath-like notation:

**Input:**
```json
{
  "user": {
    "name": "John",
    "address": {
      "city": "NYC"
    }
  }
}
```

**Output:**
```json
{
  "$['user']['name']": "John",
  "$['user']['address']['city']": "NYC"
}
```

### With -s (SQL SELECT)

Generates DuckDB SQL to extract the flattened values:

**Input:**
```json
[{"a": 1, "b": {"c": 2}}]
```

**Output:**
```sql
select
    json[0]."a"::decimal,
    json[0]."b"."c"::decimal
from (select '[{"a": 1, "b": {"c": 2}}]'::json as json)
```

## Handling Arrays

Arrays are indexed in the output:

**Input:**
```json
[1, 2, 3]
```

**Output:**
```json
{
  "$[0]": 1,
  "$[1]": 2,
  "$[2]": 3
}
```

## Type Inference

The SQL mode infers types from sample values:
- Numbers → `decimal`
- Strings → `varchar`
- Booleans → `boolean`

## Best Practices

- Use `-s` flag when you need to query the data in SQL
- Flat JSON output is useful for ETL pipelines
- Works with NDJSON files (newline-delimited JSON)
- Handle large files by piping through the command

## Integration Examples

### With DuckDB

```bash
# Generate and execute SQL
dct flattify api_response.json -s | duckdb

# Or save and use
dct flattify data.json -s -o extract.sql
duckdb mydb.duckdb < extract.sql
```

### In Pipeline

```bash
# Flatten and convert to CSV
dct flattify nested.json | jq -r 'to_entries[] | [.key, .value] | @csv'
```

## Common Use Cases

### API Response Transformation

```bash
# Flatten a complex API response
curl -s https://api.example.com/users | dct flattify -s > users.sql

# Query with DuckDB
duckdb -c "$(cat users.sql)"
```

### Document Store to Relational

```bash
# Convert MongoDB export to flat format
dct flattify mongodb_export.json -o flat_export.json
```

## Related Skills

- `dct-peek`: Preview JSON structure before flattening
- `dct-infer`: Generate schema from flattened data
- `dct-js2sql`: Convert JSON Schema (not data) to SQL

## Limitations

- Very deeply nested structures (>100 levels) may hit limits
- Mixed types in arrays use the first type encountered
- Large files should be processed in chunks
