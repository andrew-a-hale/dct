# DCT (Data Check Tool)

A Swiss army knife for data engineers to quickly inspect, compare, and generate flat files.

## Overview

DCT provides a collection of command-line utilities for working with CSV, JSON, NDJSON, and Parquet files:

- **Peek**: Preview file contents
- **Diff**: Compare files with key matching and aggregates
- **Chart**: Generate simple visualisations from data files
- **Generator**: Generate synthetic data with customisable schemas
- **Flattify**: Convert nested JSON structures to flat formats or SQL
- **Prof**: Profile data files for values and characters

## Commands

Some examples are available in the `examples/` directory to use with `dct`.

### Peek

Preview file contents:

```bash
dct peek <file> [flags]
  -o, --output <file>    Output to file (default stdout)
  -n, --lines <number>   Number of lines to display

Examples
dct peek examples/left.parquet -n 5

╭──────┬──────┬───────╮
│  a   │  b   │   c   │
│BIGINT│BIGINT│VARCHAR│
│──────│──────│───────│
│  1   │  1   │  b%$  │
│  1   │  2   │  2%$  │
│  1   │  2   │  b%$  │
│  1   │  2   │  b%$  │
│  1   │  2   │  b%$  │
╰──────┴──────┴───────╯
```

### Diff

Compare two files with key matching and metrics:

```bash
dct diff <keys> <file1> <file2> [flags]
  -o, --output <file>    Output to file (default stdout)
  -m, --metrics <spec>   Metrics specification
  -a, --all              Show all metrics

Key spec format: left_key[=right_key]
Metrics spec:
  - JSON: {agg: {left: col, right: col}, ...}
  - File path: {file}.json
  - Aggregations: mean, median, min, max, count_distinct

Example
dct diff a examples/left.parquet examples/right.csv -m '{"metrics":[{"agg":"count_distinct","left":"c","right":"c"}]}'

╭──────┬──────┬──────┬───────┬──────────────────┬──────────────────┬───────────────────╮
│  a   │l_cnt │r_cnt │cnt_eq │l_c_count_distinct│r_c_count_distinct│c_count_distinct_eq│
│BIGINT│BIGINT│BIGINT│BOOLEAN│      BIGINT      │      BIGINT      │      BOOLEAN      │
│──────│──────│──────│───────│──────────────────│──────────────────│───────────────────│
│  1   │  6   │  7   │ false │        2         │        1         │       false       │
╰──────┴──────┴──────┴───────┴──────────────────┴──────────────────┴───────────────────╯
```

### Chart

Generate simple charts from data:

```bash
dct chart [file] [colIndex] [agg] [flags]

Flags:
  -w, --width int32   Width of the chart in characters

Examples
dct chart -w 50 examples/left.csv 1 count

  ┌─ histogram of 'left.csv' ────────────────────┐
1 ┤ ╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍ 6.00 │
2 ┤ ╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍ 5.00       │
  └──────────────────────────────────────────────┘
```

### Generator

Generate synthetic data:

```bash
dct gen <schema json file or json> [flags]

Flags:
  -n, --lines int        Number of data rows to generate (default 1)
  -o, --outfile string   Output file path (default stdout)
  -f, --format string    Output format: csv, ndjson (default "csv")

Examples
dct gen examples/faker-comp.json -n 5

first_name,last_name,company,email,phone_number
DAVID,MORALES,MACY'S,STEVEN.GUTIERREZ@FACEBOOK.com,49676067
JILLIAN,DAVIS,MURPHY USA,HUNTER.NELSON@LOEWS.com,43928996
DIEGO,PETERSON,SEMPRA ENERGY,RUBY.EVANS@IHEARTMEDIA.com,45321086
MAKENZIE,WILSON,EBAY,SERENITY.KING@STATEFARMINSURANCE.com,22255385
CALEB,MOORE,GROUP 1 AUTOMOTIVE,DIANA.TORRES@W.W.GRAINGER.com,58402680
```

#### DSL For Derived Fields

See [here](https://expr-lang.org/docs/language-definition) for DSL specification for schema.

### Flattify

Convert nested JSON to flat formats:

```bash
dct flattify <json file or json> [flags]
  -s, --sql              Create DuckDB-compliant SQL statement
  -o, --output <file>    Output to file (default stdout)

Examples
dct flattify examples/faker-comp.json

{
  "1.data_type": "string",
  "1.field": "first_name",
  "1.source": "firstNames",
  "2.data_type": "string",
  "2.field": "last_name",
  "2.source": "lastNames",
  "3.data_type": "string",
  "3.field": "company",
  "3.source": "companies",
  "4.source": "emails",
  "4.data_type": "string",
  "4.field": "email",
  "5.data_type": "int",
  "5.field": "phone_number",
  "5.source": "randomUniformInt",
  "5.config.min": 10000000,
  "5.config.max": 99999999
}
```

### Profile

Provide summaries for data files:

```bash
dct prof <file> [flags]
  -o, --output <file>    Output to file (default stdout)

Examples
dct prof examples/messy.csv

-- PROFILE -- 
-- Field: `Description` -- 
Count: 10
Unique Count: 10

Value Occurrence
MOSTLY UNIQUE VALUES JUST SHOWING SAMPLE OF 10
row: value -> count
0: NULL -> 1
1: Contains \0 null -> 1
2: Simple description -> 1
3: Special chars: \/, \\, | -> 1
4: SQL Injection attempt -> 1
5: Non-ASCII chars -> 1
6: Contains , commas -> 1
7: Contains "double quotes" -> 1
8: Contains
line breaks -> 1
9: Contains\ttabs -> 1

Value Summary - String Lengths
Min: 4
Mean: 17.300000
Max: 24

Char Occurrence
row: rune -> count
00: 'n' (unicode: U+006E) (UTF-8: 110) -> 16
01: 's' (unicode: U+0073) (UTF-8: 115) -> 12
02: 'L' (unicode: U+004C) (UTF-8: 76) -> 3
03: '|' (unicode: U+007C) (UTF-8: 124) -> 1
04: 't' (unicode: U+0074) (UTF-8: 116) -> 13
05: '"' (unicode: U+0022) (UTF-8: 34) -> 2
06: 'b' (unicode: U+0062) (UTF-8: 98) -> 3
07: 'r' (unicode: U+0072) (UTF-8: 114) -> 4
08: 'p' (unicode: U+0070) (UTF-8: 112) -> 4
09: 'c' (unicode: U+0063) (UTF-8: 99) -> 6
10: ':' (unicode: U+003A) (UTF-8: 58) -> 1
11: 'I' (unicode: U+0049) (UTF-8: 73) -> 3
12: 'q' (unicode: U+0071) (UTF-8: 113) -> 1
13: '\\' (unicode: U+005C) (UTF-8: 92) -> 5
14: 'S' (unicode: U+0053) (UTF-8: 83) -> 4
15: ',' (unicode: U+002C) (UTF-8: 44) -> 3
16: 'U' (unicode: U+0055) (UTF-8: 85) -> 1
17: 'a' (unicode: U+0061) (UTF-8: 97) -> 12
18: 'i' (unicode: U+0069) (UTF-8: 105) -> 11
19: 'u' (unicode: U+0075) (UTF-8: 117) -> 3
20: 'N' (unicode: U+004E) (UTF-8: 78) -> 2
21: 'Q' (unicode: U+0051) (UTF-8: 81) -> 1
22: 'A' (unicode: U+0041) (UTF-8: 65) -> 1
23: '\n' (unicode: U+000A) (UTF-8: 10) -> 1
24: 'k' (unicode: U+006B) (UTF-8: 107) -> 1
25: 'm' (unicode: U+006D) (UTF-8: 109) -> 4
26: 'C' (unicode: U+0043) (UTF-8: 67) -> 6
27: 'o' (unicode: U+006F) (UTF-8: 111) -> 11
28: ' ' (unicode: U+0020) (UTF-8: 32) -> 15
29: '-' (unicode: U+002D) (UTF-8: 45) -> 1
30: 'd' (unicode: U+0064) (UTF-8: 100) -> 2
31: 'l' (unicode: U+006C) (UTF-8: 108) -> 6
32: '/' (unicode: U+002F) (UTF-8: 47) -> 1
33: 'j' (unicode: U+006A) (UTF-8: 106) -> 1
34: '0' (unicode: U+0030) (UTF-8: 48) -> 1
35: 'h' (unicode: U+0068) (UTF-8: 104) -> 2
36: 'e' (unicode: U+0065) (UTF-8: 101) -> 9

Char Analysis
Control: 0
Comma: 3
Pipe: 1
Quotes: 2
Nonspace-Whitespace: 1
NonAscii: 0
Rest: 151
```
