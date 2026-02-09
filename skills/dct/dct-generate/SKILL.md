---
name: dct-generate
description: Use this skill when the user wants to create synthetic test data, generate fake datasets, create mock data for testing, produce realistic data with specific patterns, or need sample data with custom schemas. Triggers include "generate test data", "create fake data", "mock dataset", "synthetic data", "generate sample records", "create test data", "fake users", "mock data", or when needing test data with specific fields and relationships.
---

# DCT Generate - Create Synthetic Data

Generate realistic test data with customizable schemas and field types.

## When to Use

Use this skill when you need to:
- Create test datasets for development
- Generate mock data for demos
- Produce synthetic data for testing ETL pipelines
- Create data with specific distributions
- Generate data with referential integrity

## Installation

```bash
which dct || go build -o dct && chmod +x ./dct
```

## Usage

```bash
dct gen <schema> [flags]
```

## Arguments

- `schema`: JSON schema as a file path or inline JSON string

## Flags

- `-n, --lines <number>`: Number of rows to generate (default: 1)
- `-f, --format <format>`: Output format - csv, ndjson (default: csv)
- `-o, --outfile <file>`: Output file path (default: stdout)

## Examples

From schema file:
```bash
dct gen schema.json -n 1000 -o test_data.csv
```

Inline schema:
```bash
dct gen '[{"field":"name","source":"firstNames"}]' -n 100
```

NDJSON output:
```bash
dct gen schema.json -n 500 -f ndjson -o output.ndjson
```

Generate to stdout:
```bash
dct gen users-schema.json -n 10
```

## Schema Format

Array of field objects:
```json
[
  {
    "field": "column_name",
    "source": "source_type",
    "config": { ... }
  }
]
```

## Available Data Sources

### Random Generators

- `randomBool` - Boolean true/false
  ```json
  {"field": "active", "source": "randomBool"}
  ```

- `randomEnum` - Random value from list
  ```json
  {"field": "status", "source": "randomEnum", "config": {"values": ["pending", "active", "inactive"]}}
  ```

- `randomAscii` - Random ASCII string
  ```json
  {"field": "code", "source": "randomAscii", "config": {"length": 10}}
  ```

- `randomUniformInt` - Uniform integer distribution
  ```json
  {"field": "age", "source": "randomUniformInt", "config": {"min": 18, "max": 65}}
  ```

- `randomNormal` - Normal/Gaussian distribution
  ```json
  {"field": "score", "source": "randomNormal", "config": {"mean": 100, "std": 15}}
  ```

- `randomPoisson` - Poisson distribution
  ```json
  {"field": "events", "source": "randomPoisson", "config": {"lambda": 5}}
  ```

- `randomDatetime` - Random date/time
  ```json
  {"field": "created_at", "source": "randomDatetime", "config": {"min": "2024-01-01 00:00:00", "max": "2024-12-31 23:59:59", "tz": "UTC"}}
  ```

- `randomDate` - Random date
  ```json
  {"field": "birth_date", "source": "randomDate", "config": {"min": "1980-01-01", "max": "2005-12-31"}}
  ```

- `randomTime` - Random time
  ```json
  {"field": "meeting_time", "source": "randomTime", "config": {"min": "09:00:00", "max": "17:00:00"}}
  ```

### Data Generators

- `uuid` - UUID v4
  ```json
  {"field": "id", "source": "uuid"}
  ```

- `firstNames` - Random first names
  ```json
  {"field": "first_name", "source": "firstNames"}
  ```

- `lastNames` - Random last names
  ```json
  {"field": "last_name", "source": "lastNames"}
  ```

- `companies` - Company names
  ```json
  {"field": "company", "source": "companies"}
  ```

- `emails` - Email addresses
  ```json
  {"field": "email", "source": "emails"}
  ```

### Derived Fields

Create computed fields using the Expr language:

```json
{
  "field": "full_name",
  "source": "derived",
  "config": {
    "fields": ["first_name", "last_name"],
    "expression": "first_name + ' ' + last_name"
  }
}
```

Complex expressions:
```json
{
  "field": "display_name",
  "source": "derived",
  "config": {
    "fields": ["first_name", "last_name", "company"],
    "expression": "first_name + ' ' + last_name + ' (' + company + ')'"
  }
}
```

## Complete Schema Example

```json
[
  {"field": "id", "source": "uuid"},
  {"field": "first_name", "source": "firstNames"},
  {"field": "last_name", "source": "lastNames"},
  {"field": "email", "source": "emails"},
  {"field": "age", "source": "randomUniformInt", "config": {"min": 18, "max": 65}},
  {"field": "department", "source": "randomEnum", "config": {"values": ["Engineering", "Sales", "Marketing", "HR"]}},
  {"field": "salary", "source": "randomNormal", "config": {"mean": 75000, "std": 15000}},
  {"field": "is_active", "source": "randomBool"},
  {
    "field": "full_name",
    "source": "derived",
    "config": {
      "fields": ["first_name", "last_name"],
      "expression": "first_name + ' ' + last_name"
    }
  }
]
```

## Best Practices

- Generate small samples first (n=10) to verify schema
- Use derived fields to create realistic relationships
- Use NDJSON format for nested/complex data
- Save schemas to files for reuse
- Use appropriate distributions for realistic data

## Output Formats

**CSV** (default):
```csv
id,first_name,age
550e8400-e29b-41d4-a716-446655440000,John,34
```

**NDJSON**:
```json
{"id":"550e8400-e29b-41d4-a716-446655440000","first_name":"John","age":34}
{"id":"550e8400-e29b-41d4-a716-446655440001","first_name":"Jane","age":28}
```

## Related Skills

- `dct-peek`: Verify generated data looks correct
- `dct-infer`: Check schema of generated data
- `dct-diff`: Compare generated data with production samples
