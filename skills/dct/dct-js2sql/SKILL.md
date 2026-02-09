---
name: dct-js2sql
description: Use this skill when the user wants to convert JSON Schema to SQL CREATE TABLE statements, transform schema definitions to database DDL, create SQL tables from JSON Schema files, or generate database schemas from API specifications. Triggers include "json schema to sql", "convert schema to sql", "create table from json schema", "json schema ddl", "schema conversion", or when working with OpenAPI, JSON Schema, or API specifications that need database tables.
---

# DCT JS2SQL - JSON Schema to SQL

Convert JSON Schema files to DuckDB CREATE TABLE statements.

## When to Use

Use this skill when you need to:
- Convert API specifications to database schemas
- Generate SQL from JSON Schema definitions
- Create tables from OpenAPI schemas
- Transform schema files to DDL
- Migrate from document to relational models

## Installation

```bash
which dct || go build -o dct && chmod +x ./dct
```

## Usage

```bash
dct js2sql <schema_file> [flags]
```

## Arguments

- `schema_file`: Path to JSON Schema file

## Flags

- `-t, --table <name>`: Table name (default: "test")
- `-o, --output <file>`: Output to file instead of stdout

## Examples

Basic conversion:
```bash
dct js2sql schema.json
```

With custom table name:
```bash
dct js2sql api-schema.json -t api_events
```

Save to file:
```bash
dct js2sql user-schema.json -t users -o create_users.sql
```

## Supported JSON Schema Features

### Primitive Types

Maps JSON Schema types to DuckDB types:

| JSON Schema | DuckDB |
|-------------|---------|
| `string` | `varchar` |
| `integer` | `integer` |
| `number` | `double` |
| `boolean` | `boolean` |

### Arrays

Arrays with item types:
```json
{
  "type": "array",
  "items": {
    "type": "string"
  }
}
```

Output: `array(varchar)`

### Objects/Structs

Nested objects become row types:
```json
{
  "type": "object",
  "properties": {
    "name": {"type": "string"},
    "age": {"type": "integer"}
  }
}
```

Output: `row(name varchar, age integer)`

### References ($ref)

Supports local $ref references:
```json
{
  "$ref": "#/definitions/User"
}
```

## Example Schema Conversions

### Simple Schema

**Input:**
```json
{
  "type": "object",
  "properties": {
    "id": {"type": "integer"},
    "name": {"type": "string"},
    "email": {"type": "string"},
    "active": {"type": "boolean"}
  }
}
```

**Output:**
```sql
create table test (
    id integer,
    name varchar,
    email varchar,
    active boolean
);
```

### Complex Schema with Arrays

**Input:**
```json
{
  "type": "object",
  "properties": {
    "order_id": {"type": "integer"},
    "items": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "product_id": {"type": "integer"},
          "name": {"type": "string"},
          "price": {"type": "number"},
          "quantity": {"type": "integer"}
        }
      }
    }
  }
}
```

**Output:**
```sql
create table test (
    order_id integer,
    items array(row(product_id integer, name varchar, price double, quantity integer))
);
```

### Schema with References

**Input:**
```json
{
  "definitions": {
    "Address": {
      "type": "object",
      "properties": {
        "street": {"type": "string"},
        "city": {"type": "string"}
      }
    }
  },
  "type": "object",
  "properties": {
    "name": {"type": "string"},
    "address": {"$ref": "#/definitions/Address"}
  }
}
```

**Output:**
```sql
create table test (
    name varchar,
    address row(street varchar, city varchar)
);
```

## Best Practices

- Review generated SQL for type accuracy
- Add constraints (PRIMARY KEY, NOT NULL) manually after generation
- Consider normalizing nested arrays to separate tables
- Test the SQL in your target database
- Use meaningful table names with `-t` flag

## Integration Examples

### With DuckDB

```bash
# Convert and create table
dct js2sql schema.json -t events | duckdb mydb.duckdb

# Save and review first
dct js2sql schema.json -o schema.sql
# Review schema.sql
duckdb mydb.duckdb < schema.sql
```

### In Data Pipeline

```bash
#!/bin/bash
# Convert all JSON schemas in a directory
for schema in schemas/*.json; do
    table=$(basename "$schema" .json)
    dct js2sql "$schema" -t "$table" > "sql/${table}.sql"
done
```

## Limitations

- Supports JSON Schema Draft 7 features
- Complex validations (min/max, patterns) are not converted to constraints
- External $ref references (URLs) are not resolved
- Database-specific features (indexes, triggers) must be added manually

## Related Skills

- `dct-flattify`: Convert actual JSON data (not schemas) to SQL
- `dct-infer`: Generate SQL from data files rather than schemas
- `dct-peek`: Preview JSON Schema files

## Schema Sources

Works with schemas from:
- OpenAPI specifications
- JSON Schema Store
- API documentation
- Data validation frameworks
- TypeScript-to-JSON-Schema converters
