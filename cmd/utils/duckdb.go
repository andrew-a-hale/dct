package utils

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"

	_ "github.com/marcboeker/go-duckdb"
)

var style = lipgloss.NewStyle().Align(lipgloss.Center)

type Header struct {
	Name string
	Type string
}

type Result struct {
	Headers []Header
	Rows    [][]any
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
		return Result{}, fmt.Errorf("failed to query duckdb: %v", err)
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

	var out [][]any
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

		var tmp []any
		for _, v := range vals {
			deref := reflect.Indirect(reflect.ValueOf(v)).Interface()
			switch deref.(type) {
			case nil:
				tmp = append(tmp, nil)
			case string:
				tmp = append(tmp, deref.(string))
			case bool:
				tmp = append(tmp, deref.(bool))
			case int:
				tmp = append(tmp, deref.(int))
			case int32:
				tmp = append(tmp, int(deref.(int32)))
			case int64:
				tmp = append(tmp, int(deref.(int64))) // demote to architecture
			case float32:
				tmp = append(tmp, deref.(float32))
			case float64:
				tmp = append(tmp, deref.(float64))
			default:
				return Result{}, fmt.Errorf(
					"failed to serialise rows from duckdb, type `%T` not implemented yet",
					deref,
				)
			}
		}

		out = append(out, tmp)
		row++
	}

	return Result{headers, out}, nil
}

func (result *Result) ToCsv(writer io.Writer) error {
	var headers []string
	for _, header := range result.Headers {
		headers = append(headers, header.Name)
	}

	rs := strings.Join(headers, ",")
	if _, err := fmt.Fprintln(writer, rs); err != nil {
		return err
	}

	for _, row := range result.Rows {
		var tmp []string
		for _, v := range row {
			tmp = append(tmp, fmt.Sprintf("%v", v))
		}
		rs := strings.Join(tmp, ",")
		if _, err := fmt.Fprintln(writer, rs); err != nil {
			return err
		}
	}

	return nil
}

func (r *Result) RowsToString() [][]string {
	var rows [][]string

	for _, r := range r.Rows {
		var row []string
		for _, v := range r {
			row = append(row, fmt.Sprintf("%v", v))
		}
		rows = append(rows, row)
	}

	return rows
}

func (result *Result) Render(writer io.Writer, maxRows int) error {
	var headers []string
	var types []string
	for _, header := range result.Headers {
		headers = append(headers, header.Name)
		types = append(types, header.Type)
	}

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch row {
			case 3:
				// force border after type row in display
				return style.
					Border(lipgloss.NormalBorder(), true, false, false, false)
			default:
				return style
			}
		}).
		Rows(headers).
		Rows(types)

	rowsToDisplay := min(maxRows, len(result.Rows))

	t.Rows(result.RowsToString()[:rowsToDisplay]...)

	fmt.Println(t)
	return nil
}
