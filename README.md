# dct (data check tool)

## examples

```
art:
  go run main.go art ;

diff:
  go run main.go diff -m test/resources/metrics.json a test/resources/left.csv test/resources/right.csv ;

peek:
  go run main.go peek test/resources/left.csv ;

chart:
  go run main.go chart -w 50 test/resources/left.csv 1 count ;
  go run main.go chart -w 23 test/resources/right.csv 1 sum ;
  go run main.go chart -w 10 test/resources/right.csv 1 max ;
  go run main.go chart -w 5 test/resources/chart.csv 1 count_distinct ;
  go run main.go chart test/resources/chart.csv 1 count ;

gen:
  go run main.go gen -n 2 -s test/resources/generator-schema.json ;
```

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
  - error with testing it

## todo

- use viper
- add data generator
  - dct gen -o -n -s
    - output
    - n lines
    - schema
  - dct gen-sources

## DSL For Derived Fields

- Strings: ||
- Floats: +, *, /, ^, -
- Ints: +, *, /, ^, -, %
