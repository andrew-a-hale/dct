package generator

import (
	"fmt"
	"io"
)

func Export(schema []Field, out io.Writer) {
	fields := len(schema)
	for i := range lines {
		for j, f := range schema {
			fmt.Fprintf(out, "%v", f.GetValue(i))
			if j < fields-1 {
				fmt.Fprintf(out, ", ")
			} else {
				fmt.Fprintf(out, "\n")
			}
		}
	}
}
