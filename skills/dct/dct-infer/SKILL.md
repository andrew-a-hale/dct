---
name: dct-infer
description: Use this skill when the user wants to generate SQL CREATE TABLE statements from data files, infer schema from CSV/JSON/Parquet, create database schemas from existing data, or get column types from a file. Triggers include "generate schema", "create table from csv", "infer types", "what's the schema", "get column types", "sql ddl", or when preparing data for SQL databases like DuckDB, PostgreSQL, or similar.
---

# DCT Infer - Generate SQL Schema

Create DuckDB-compatible CREATE TABLE statements by analyzing data file contents.

## When to Use

Use this skill when you need to:
- Create database tables from existing data files
- Document the schema of a dataset
- Generate DDL for ETL pipelines
- Understand column types in a file
- Prepare data for SQL-based analysis

## Installation

```bash
which dct || go build -o dct && chmod +x ./dct
```

## Usage

```bash
dct infer <file> [flags]
```

## Flags

- `-t, --table <name>`: Table name (default: "default")
- `-n, --lines <number>`: Number of lines to analyze for type inference (useful for large files)
- `-o, --output <file>`: Output to file instead of stdout

## Examples

Basic schema inference:
```bash
dct infer data.csv
```

With custom table name:
```bash
dct infer data.parquet -t events
```

Save schema to file:
```bash
dct infer large.ndjson -n 1000 -t users -o schema.sql
```

Infer from specific number of rows:
```bash
dct infer bigfile.csv -n 500 -t transactions
```

## Output Format

DuckDB-compatible CREATE TABLE statement:

```sql
create table users (
    "id" bigint,
    "name" varchar,
    "email" varchar,
    "created_at" timestamp,
    "is_active" boolean
)
```

## Supported Data Types

The inferred schema uses DuckDB types:
- `bigint` - 64-bit integers
- `integer` - 32-bit integers
- `double` - Floating point numbers
- `varchar` - String/text data
- `timestamp` - Date and time
- `date` - Date only
- `time` - Time only
- `boolean` - True/false values
- `array(...)` - Array columns
- `row(...)` - Struct/nested columns

## Best Practices

- Use `-n` flag for large files to speed up inference
- Column names are quoted to handle special characters
- Output is compatible with DuckDB and similar SQL databases
- For Parquet files, types are read directly from metadata
- For CSV/JSON, types are inferred from sample data

## Integration Examples

### With DuckDB
```bash
# Create table directly
dct infer data.csv -t my_table | duckdb mydb.duckdb

# Or save and execute
dct infer data.csv -t my_table -o schema.sql
duckdb mydb.duckdb < schema.sql
```

### In Scripts
```bash
#!/bin/bash
for file in *.csv; do
    dct infer "$file" -t "$(basename "$file" .csv)" > "${file%.csv}.sql"
done
```

## Related Skills

- `dct-peek`: Preview data before inferring schema
- `dct-profile`: Check data quality before creating tables
