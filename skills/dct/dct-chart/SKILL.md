---
name: dct-chart
description: Use this skill when the user wants to visualize data distributions, create ASCII histograms, generate simple charts from CSV/JSON data, plot column values, or see value frequencies in terminal-friendly format. Triggers include "chart this data", "visualize distribution", "histogram of values", "plot the data", "ascii chart", "terminal visualization", or when needing quick visual analysis without external plotting tools.
---

# DCT Chart - ASCII Data Visualization

Generate simple ASCII histograms from data file columns.

## When to Use

Use this skill when you need to:
- Quickly visualize value distributions
- Create terminal-friendly charts
- Analyze categorical data frequencies
- Get a quick histogram without external tools
- Share visualizations in text format

## Installation

```bash
which dct || go build -o dct && chmod +x ./dct
```

## Usage

```bash
dct chart [file] [column_index] [flags]
```

## Arguments

- `file`: Data file path (optional)
- `column_index`: 0-based column index to chart (default: 0)

## Flags

- `-w, --width <chars>`: Width of chart in characters (default: auto)

## Examples

### Basic Usage

Chart first column:
```bash
dct chart data.csv
```

Chart specific column (0-based index):
```bash
dct chart data.csv 2
```

### Custom Width

Wider chart:
```bash
dct chart sales.csv 1 -w 80
```

Compact chart:
```bash
dct chart data.csv 0 -w 40
```

### Different File Types

Chart from Parquet:
```bash
dct chart data.parquet 3
```

Chart from JSON:
```bash
dct chart data.json 0
```

## Output Format

Displays an ASCII histogram with:
- Column values on the left (Y-axis labels)
- Bar representation using box-drawing characters
- Count/frequency on the right
- Box border around the chart

Example output:
```
            ┌─ histogram of 'sales.csv' ──────────────────────────┐
    ProductA ┤ ╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍ 156 │
    ProductB ┤ ╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍ 98 │
    ProductC ┤ ╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍ 73 │
    ProductD ┤ ╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍ 45 │
    ProductE ┤ ╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍ 22 │
             └───────────────────────────────────────────────────┘
```

## How It Works

1. Reads the specified column from the data file
2. Groups values and counts frequencies
3. Sorts by frequency (descending)
4. Renders ASCII bar chart scaled to fit terminal

## Best Practices

- **Column indexing**: Remember column index is 0-based (first column = 0)
- **Use `-w` flag**: Adjust width to fit your terminal
- **Categorical data**: Works best with categorical or binned data
- **High cardinality**: Top values are shown; long tail may be truncated
- **Preview first**: Use `dct-peek` to see column names before charting

## Use Cases

### Sales Analysis

```bash
# Chart sales by product
dct chart sales.csv 0 -w 60

# Chart revenue distribution
dct chart orders.csv 2 -w 80
```

### Data Exploration

```bash
# Quick distribution check
dct peek data.csv -n 5
dct chart data.csv 1 -w 50
```

### Log Analysis

```bash
# Chart status codes from logs
dct chart logs.csv 3 -w 40
```

## Limitations

- Shows top N values (high-cardinality columns truncate)
- ASCII art has limited resolution
- Best for categorical or discrete numeric data
- Continuous numeric data should be pre-binned

## Column Selection Workflow

```bash
# 1. Preview file to see columns
dct peek data.csv

# 2. Note the column index (0-based)
#    Columns: id | name | category | amount
#    Index:    0    1        2         3

# 3. Chart desired column
dct chart data.csv 2 -w 60
```

## Integration Examples

### In Shell Scripts

```bash
#!/bin/bash
# Generate charts for all columns in a file
file="data.csv"
num_cols=$(head -1 "$file" | tr ',' '\n' | wc -l)
for i in $(seq 0 $((num_cols - 1))); do
    echo "Column $i:"
    dct chart "$file" $i -w 50
    echo
done
```

### With Other Tools

```bash
# Filter data then chart
grep "ERROR" logs.csv | dct chart - 1 -w 40

# Save chart to file
dct chart data.csv 0 -w 80 > chart.txt
```

## Related Skills

- `dct-peek`: Preview column names before charting
- `dct-profile`: Get detailed statistics about the column
- `dct-infer`: Understand column types

## Terminal Tips

- Charts look best in terminals with UTF-8 support
- Use wider terminals (`-w 100+`) for detailed views
- Pipe to `less -S` for wide charts: `dct chart data.csv 0 -w 120 | less -S`
