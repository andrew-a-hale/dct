---
name: dct-peek
description: Use this skill when the user wants to preview or inspect the contents of a data file (CSV, JSON, NDJSON, Parquet). Triggers include "show me the data", "preview this file", "what's in this csv", "look at the first rows", "sample the data", or when needing to understand data structure before processing. This is often the first step before other data operations.
---

# DCT Peek - Preview File Contents

Preview the first N rows of any supported flat data file to understand its structure and content.

## When to Use

Use this skill as the first step when working with unfamiliar data files:
- Before running transformations or analysis
- To verify file format and structure
- To check data types and column names
- To sample data for quick inspection

## Installation

```bash
which dct || go build -o dct && chmod +x ./dct
```

## Usage

```bash
dct peek <file> [flags]
```

## Flags

- `-n, --lines <number>`: Number of lines to display (default: 10)
- `-o, --output <file>`: Output to file instead of stdout

## Examples

Preview default 10 lines:
```bash
dct peek data.csv
```

Preview specific number of lines:
```bash
dct peek large.parquet -n 5
```

Save preview to CSV:
```bash
dct peek data.json -o preview.csv
```

Quick check if file is readable:
```bash
dct peek data.csv -n 1
```

## Output Format

The output displays a formatted table showing:
- Column names (header row)
- Data types (second row)
- Data rows (up to N rows)

Example:
```
╭──────┬────────┬───────╮
│  id  │  name  │ price │
│BIGINT│VARCHAR │DOUBLE │
│──────│────────│───────│
│  1   │  Alice │ 10.99 │
│  2   │   Bob  │ 24.50 │
╰──────┴────────┴────────╯
```

## Best Practices

- **Always peek first** when working with new data files
- Use `-n 1` to quickly verify a file is readable before processing
- Use `-o` to save samples for documentation or testing
- Check the data types row to understand the schema
- Verify column names match expectations

## Related Skills

After peeking, you may want to use:
- `dct-profile`: For detailed data quality analysis
- `dct-infer`: To generate SQL schema from the data
- `dct-diff`: To compare with another file
