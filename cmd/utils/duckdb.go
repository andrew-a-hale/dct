package utils

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/marcboeker/go-duckdb"
)

type Header struct {
	Name string
	Type string
}

type Result struct {
	Headers []Header
	Rows    [][]string
}

func CheckFileHasRows(file string) (bool, error) {
	conn, err := sql.Open("duckdb", "")
	if err != nil {
		return false, err
	}

	row := conn.QueryRowContext(
		context.Background(),
		fmt.Sprintf("select count(*) from '%s'", file),
	)

	var cnt int
	row.Scan(&cnt)
	return cnt > 0, nil
}

func Execute(query string) error {
	conn, err := sql.Open("duckdb", "")
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(context.Background(), query)
	if err != nil {
		return err
	}

	return nil
}

func Query(query string) (Result, error) {
	conn, err := sql.Open("duckdb", "")
	if err != nil {
		return Result{}, err
	}
	defer conn.Close()

	rows, err := conn.QueryContext(context.Background(), query)
	if err != nil {
		return Result{}, err
	}

	cols, err := rows.ColumnTypes()
	if err != nil {
		return Result{}, err
	}

	var headers []Header
	for _, col := range cols {
		headers = append(headers, Header{col.Name(), col.DatabaseTypeName()})
	}

	var out [][]string
	vals := make([]any, len(cols))

	row := 0
	for rows.Next() {
		for i := range cols {
			vals[i] = new(any)
		}

		err = rows.Scan(vals...)
		if err != nil {
			return Result{}, err
		}

		var tmp []string
		for _, v := range vals {
			if s, ok := v.(*any); ok {
				if *s == nil {
					tmp = append(tmp, string("NULL"))
				} else if x, ok := (*s).(string); ok {
					tmp = append(tmp, x)
				} else if x, ok := (*s).(int32); ok {
					tmp = append(tmp, fmt.Sprintf("%d", x))
				} else if x, ok := (*s).(int64); ok {
					tmp = append(tmp, fmt.Sprintf("%d", x))
				} else if x, ok := (*s).(*big.Int); ok {
					tmp = append(tmp, fmt.Sprintf("%d", x))
				} else if x, ok := (*s).(float64); ok {
					tmp = append(tmp, fmt.Sprintf("%f", x))
				} else if x, ok := (*s).(bool); ok {
					tmp = append(tmp, strconv.FormatBool(x))
				} else if x, ok := (*s).(time.Time); ok {
					tmp = append(tmp, x.String())
				} else if x, ok := (*s).(duckdb.Decimal); ok {
					b := x.Value.String()
					sb := b[:x.Scale] + "." + b[x.Scale:]
					tmp = append(tmp, string(sb))
				} else if x, ok := (*s).([]any); ok { // json array
					j, err := json.Marshal(x)
					if err != nil {
						return Result{}, err
					}
					tmp = append(tmp, string(j))
				} else if x, ok := (*s).(map[string]any); ok { // json object
					j, err := json.Marshal(x)
					if err != nil {
						return Result{}, err
					}
					tmp = append(tmp, string(j))
				} else {
					err := fmt.Errorf("type `%v` not implemented yet", reflect.TypeOf(*s))
					return Result{}, err
				}
			}
		}

		out = append(out, tmp)
		row++
	}

	return Result{headers, out}, nil
}

func (result Result) ToCsv(writer io.Writer) error {
	var headers []string
	for _, header := range result.Headers {
		headers = append(headers, header.Name)
	}

	rs := strings.Join(headers, ",")
	if _, err := fmt.Fprintln(writer, rs); err != nil {
		return err
	}

	for _, row := range result.Rows {
		rs := strings.Join(row, ",")
		if _, err := fmt.Fprintln(writer, rs); err != nil {
			return err
		}
	}

	return nil
}
