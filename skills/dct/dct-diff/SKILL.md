---
name: dct-diff
description: Use this skill when the user wants to compare two data files, find differences between datasets, validate data consistency, check if files have matching records, or reconcile data between sources. Triggers include "compare these files", "diff the datasets", "are these the same", "find differences", "validate data matches", "reconcile", "data comparison", or when doing data quality validation between two files.
---

# DCT Diff - Compare Datasets

Compare two data files with key matching and optional aggregation metrics.

## When to Use

Use this skill when you need to:
- Validate data consistency between two versions
- Compare production vs test data
- Reconcile data after ETL processes
- Check for data drift over time
- Validate data migrations

## Installation

```bash
which dct || go build -o dct && chmod +x ./dct
```

## Usage

```bash
dct diff <keys> <file1> <file2> [flags]
```

## Arguments

- `keys`: Key column(s) for matching records. Formats:
  - Single key: `id`
  - Composite keys: `key1,key2`
  - Different names: `left_col=right_col`
- `file1`: First data file (left side)
- `file2`: Second data file (right side)

## Flags

- `-m, --metrics <spec>`: Metrics specification (JSON string or file path)
- `-a, --all`: Show all metrics columns
- `-o, --output <file>`: Output to file instead of stdout

## Examples

### Basic Comparison

Compare by single key:
```bash
dct diff id left.csv right.csv
```

Compare by composite keys:
```bash
dct diff "first_name,last_name" file1.parquet file2.parquet
```

### Key Name Mapping

When key columns have different names:
```bash
dct diff user_id=customer_id old.csv new.csv
```

### With Metrics

Compare with count distinct metric:
```bash
dct diff id left.csv right.csv -m '[{"agg":"count_distinct","left":"email","right":"email"}]'
```

Multiple metrics:
```bash
dct diff id left.csv right.csv -m '[{"agg":"mean","left":"amount","right":"amount"},{"agg":"count_distinct","left":"category","right":"category"}]'
```

Load metrics from file:
```bash
dct diff id left.csv right.csv -m metrics.json -a
```

## Metrics Specification

JSON array of metric objects:
```json
[
  {
    "agg": "count_distinct",
    "left": "column_name",
    "right": "column_name"
  }
]
```

### Available Aggregations

- `mean` - Average value
- `median` - Median value
- `min` - Minimum value
- `max` - Maximum value
- `sum` - Sum of values
- `count` - Count of records
- `count_distinct` - Count of unique values

## Output Columns

Default output includes:
- Key column(s)
- `l_cnt` - Count from left file
- `r_cnt` - Count from right file
- `cnt_eq` - Whether counts match

With metrics and `-a` flag:
- `l_<col>_<agg>` - Left aggregation
- `r_<col>_<agg>` - Right aggregation
- `<col>_<agg>_eq` - Whether aggregations match

## Best Practices

- Use `-a` flag to see all comparison metrics
- Both files must contain the key columns
- Files must have at least one row of data
- Start with a small sample to verify keys work
- Use composite keys when single keys aren't unique

## Error Handling

Common issues:
- `attempted to diff when least one of the files have no data`: Check files aren't empty
- Key not found: Verify column names match exactly (case-sensitive)
- Format errors: Ensure metrics JSON is valid

## Example Workflow

```bash
# 1. Preview both files first
dct peek left.csv -n 3
dct peek right.csv -n 3

# 2. Compare by ID
dct diff id left.csv right.csv -a

# 3. Save results
dct diff id left.csv right.csv -m metrics.json -a -o comparison.csv
```

## Related Skills

- `dct-peek`: Preview files before comparing
- `dct-profile`: Check data quality of each file
