# DCT (Data Check Tool)

A Swiss army knife for data engineers to quickly inspect, compare, and generate flat files.

## Overview

DCT provides a collection of command-line utilities for working with CSV, JSON, NDJSON, and Parquet files:

- **Peek**: Preview file contents
- **Diff**: Compare files with key matching and aggregates
- **Chart**: Generate simple visualisations from data files
- **Generator**: Generate synthetic data with customisable schemas
- **Flattify**: Convert nested JSON structures to flat formats or SQL

## Commands

Some examples are available in the `Makefile`.

### Peek

Preview file contents:

```bash
dct peek <file> [options]
  -o, --output <file>    Output to file (default: stdout)
  -n, --lines <number>   Number of lines to display
```

### Diff

Compare two files with key matching and metrics:

```bash
dct diff <keys> <file1> <file2> [options]
  -o, --output <file>    Output to file (default: stdout)
  -m, --metrics <spec>   Metrics specification
  -a, --all              Show all metrics

Key spec format: left_key[=right_key]
Metrics spec:
  - JSON: {agg: {left: col, right: col}, ...}
  - File path: {file}.json
  - Aggregations: mean, median, min, max, count_distinct
```

### Chart

Generate simple charts from data:

```bash
dct chart <file> <colIndex> <aggregation>
Example: dct chart left.csv 1 sum
```

### Generator

Generate synthetic data:

```bash
dct gen -s <schema_file> -n <count>
Example: dct gen -s test/resources/faker-comp.json -n 10000
```

#### DSL For Derived Fields

- Strings: ||
- Floats: +, *, /, ^, -
- Ints: +, *, /, ^, -, %

### Flattify

Convert nested JSON to flat formats:

```bash
dct flattify <file> [options]
  -s, --sql                  Create DuckDB-compliant SQL statement
  -o, --output <file>    Output to file
```
