# dct (data check tool)

## cli

- dct version
- dct peek -o {file=out.txt} -n number of lines
  - --output -o
  - --lines -n
- dct diff {keys} {file1} {file2} -o {file=out.txt} -m {spec} -a
  - key spec: `left_key[=right_key]`
  - --output -o
  - --metrics -m
  - --all -a
    - spec
      - {agg: {left: col, right: col}, ...}
      - {file}.json
      - aggs: mean, median, min, max, count_distinct
- dct art
- dct chart {file} {colIndex} {agg}
  - eg. `dct chart left.csv 1 sum`
- dct gen -s test/resources/faker-comp.json -n 10000
- dct flattify test/resources/flattify.ndjson -o out.sql

## DSL For Derived Fields

- Strings: ||
- Floats: +, *, /, ^, -
- Ints: +, *, /, ^, -, %
