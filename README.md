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

Some examples are available in the `examples/` directory.

### Peek

Preview file contents:

```bash
dct peek <file> [options]
  -o, --output <file>    Output to file (default: stdout)
  -n, --lines <number>   Number of lines to display

Examples
dct peek examples/left.csv
dct peek examples/left.csv -n 10
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

Example
dct diff a examples/left.csv examples/right.csv -m examples/metrics.json
```

### Chart

Generate simple charts from data:

```bash
dct chart [file] [col-name] [agg] [flags]

Flags:
  -w, --width int32   Width of the chart in characters

Examples
dct chart -w 50 examples/left.csv 1 count
dct chart -w 23 examples/right.csv 1 sum
dct chart -w 10 examples/right.csv 1 max
dct chart -w 5 examples/chart.csv 1 count_distinct
dct chart examples/chart.csv 1 count
```

### Generator

Generate synthetic data:

```bash
dct gen -s [schema] -n [lines] -o [outfile] [flags]

Flags:
  -n, --lines int        Number of data rows to generate
  -o, --outfile string   Output file path (default: stdout)
  -s, --schema string    Schema definition file path

Examples
dct gen -n 200 -s examples/generator-schema.json
dct gen -n 20000 -s examples/faker-comp.json
```

#### DSL For Derived Fields

- Strings: ||
- Floats: +, *, /, ^, -
- Ints: +, *, /, ^, -, %

### Flattify

Convert nested JSON to flat formats:

```bash
dct flattify <file> [options]
  -s, --sql              Create DuckDB-compliant SQL statement
  -o, --output <file>    Output to file

Example
dct flattify -s examples/flattify.ndjson
dct flattify examples/flattify.ndjson
```
