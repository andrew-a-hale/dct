package generator

import (
	"fmt"
	"io"
)

func Export(schema []Field, out io.Writer) {
	fields := len(schema)
	for i, f := range schema {
		value := f.GetValue()
		// handle quoting
		switch f.GetType() {
		case STRING:
			fmt.Fprintf(out, "\"%v\"", value)
		default:
			fmt.Fprintf(out, "%v", value)
		}
		if i < fields-1 {
			fmt.Fprintf(out, ",")
		} else {
			fmt.Fprintf(out, "\n")
		}
	}
}
