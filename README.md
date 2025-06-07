# DCT (Data Check Tool)

A Swiss army knife for data engineers to quickly inspect, compare, and generate
flat files.

## Overview

DCT provides a collection of command-line utilities for working with CSV, JSON,
NDJSON, and Parquet files:

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
  - JSON: [{agg: left: col, right: col}, ...]
  - Aggregations: mean, median, min, max, count_distinct

Example
dct diff a examples/left.parquet examples/right.csv -m '[{"agg":"count_distinct","left":"c","right":"c"}]'

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
dct chart [file] [colIndex] [flags]

Flags:
  -w, --width int32   Width of the chart in characters

Examples
dct chart -w 50 examples/chart.csv 1

            ┌─ histogram of 'chart.csv' ─────────┐
        xyz ┤ ╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍ 42.00 │
        bcd ┤ ╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍ 24.00             │
        abc ┤ ╍╍╍╍╍╍╍╍╍╍╍╍ 18.00                 │
        123 ┤ ╍╍╍╍ 6.00                          │
23467w81234 ┤ ╍ 1.00                             │
            └────────────────────────────────────┘
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
BRYAN,PETERSON,EMCOR GROUP,BRYAN.PETERSON@EMCORGROUP.COM,24262667
VALERIA,BROOKS,KRAFT HEINZ,VALERIA.BROOKS@KRAFTHEINZ.COM,52950975
KENDALL,CASTILLO,REINSURANCE GROUP OF AMERICA,KENDALL.CASTILLO@REINSURANCEGROUPOFAMERICA.COM,63120507
AIDEN,PRICE,FORD MOTOR,AIDEN.PRICE@FORDMOTOR.COM,74167250
DYLAN,ALLEN,JACOBS ENGINEERING GROUP,DYLAN.ALLEN@JACOBSENGINEERINGGROUP.COM,83166063
```

#### DSL For Derived Fields

Expressions for the Derived Fields use [Expr](https://expr-lang.org/docs/language-definition).

### Flattify

Convert nested JSON to flat formats:

```bash
dct flattify <json file or json> [flags]
  -s, --sql              Create DuckDB-compliant SQL Select statement
  -o, --output <file>    Output to file (default stdout)

Examples
dct flattify '[ 1, {"value": {"nested": [0, 1]}}]'

{
  "$[0]": 1,
  "$[1]['value']['nested'][0]": 0,
  "$[1]['value']['nested'][1]": 1
}

dct flattify '[ 1, {"value": {"nested": [0, 1]}}]' -s

select
        json[0]::decimal
        , json[1]."value"."nested"[0]::decimal
        , json[1]."value"."nested"[1]::decimal
from (select '[ 1, {"value": {"nested": [0, 1]}}]'::json as json)
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
MOSTLY UNIQUE VALUES SHOWING SAMPLE...
row: value -> count
0: Contains "double quotes" -> 1
1: Contains
line breaks -> 1
2: Contains\ttabs -> 1
3: Special chars: \/, \\, | -> 1
4: Contains \0 null -> 1
5: SQL Injection attempt -> 1
6: <nil> -> 1
7: Non-ASCII chars -> 1
8: Simple description -> 1
9: Contains , commas -> 1

Value Summary - String Lengths
Min: 5
Mean: 17.400000
Max: 24

Char Occurrence
row: rune -> count
00: 'l' (hex: U+006C) (dec: 108) -> 7
01: 'b' (hex: U+0062) (dec: 98) -> 3
02: '/' (hex: U+002F) (dec: 47) -> 1
03: 'Q' (hex: U+0051) (dec: 81) -> 1
04: 'L' (hex: U+004C) (dec: 76) -> 1
05: 'I' (hex: U+0049) (dec: 73) -> 3
06: 'A' (hex: U+0041) (dec: 65) -> 1
07: 'a' (hex: U+0061) (dec: 97) -> 12
08: 'u' (hex: U+0075) (dec: 117) -> 3
09: ':' (hex: U+003A) (dec: 58) -> 1
10: 'N' (hex: U+004E) (dec: 78) -> 1
11: '-' (hex: U+002D) (dec: 45) -> 1
12: 'S' (hex: U+0053) (dec: 83) -> 4
13: 'd' (hex: U+0064) (dec: 100) -> 2
14: 's' (hex: U+0073) (dec: 115) -> 12
15: '|' (hex: U+007C) (dec: 124) -> 1
16: 'm' (hex: U+006D) (dec: 109) -> 4
17: 'e' (hex: U+0065) (dec: 101) -> 9
18: ' ' (hex: U+0020) (dec: 32) -> 15
19: 't' (hex: U+0074) (dec: 116) -> 13
20: 'o' (hex: U+006F) (dec: 111) -> 11
21: ',' (hex: U+002C) (dec: 44) -> 3
22: '>' (hex: U+003E) (dec: 62) -> 1
23: 'q' (hex: U+0071) (dec: 113) -> 1
24: '\\' (hex: U+005C) (dec: 92) -> 5
25: 'n' (hex: U+006E) (dec: 110) -> 17
26: 'c' (hex: U+0063) (dec: 99) -> 6
27: 'r' (hex: U+0072) (dec: 114) -> 4
28: 'h' (hex: U+0068) (dec: 104) -> 2
29: '<' (hex: U+003C) (dec: 60) -> 1
30: 'i' (hex: U+0069) (dec: 105) -> 12
31: 'C' (hex: U+0043) (dec: 67) -> 6
32: '"' (hex: U+0022) (dec: 34) -> 2
33: '0' (hex: U+0030) (dec: 48) -> 1
34: 'j' (hex: U+006A) (dec: 106) -> 1
35: 'k' (hex: U+006B) (dec: 107) -> 1
36: 'p' (hex: U+0070) (dec: 112) -> 4
37: '\n' (hex: U+000A) (dec: 10) -> 1
```
