package generator

import "io"

type Export struct {
	Dest io.Writer
	Rows []any
	N    int
}

func buildExport(schema []Field, out io.Writer) Export {
	var rows []any

	for i := 0; i < lines; i++ {
		var row []any
		for _, f := range schema {
			row = append(row, f.GetValue(i))
		}
		rows = append(rows, row)
	}

	return Export{out, rows, lines}
}
