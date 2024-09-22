# dct (data check tool)

## examples

## setup
- requires duckdb

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

## todo
- use viper
- add test cases
- add examples

