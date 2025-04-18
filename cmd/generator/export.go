package generator

import (
	"fmt"
	"io"
	"strings"
)

func Export(schema []Field, out io.Writer, line int) {
	fields := len(schema)

	// headers
	if line == 0 {
		for i, f := range schema {
			name := f.GetName()
			if strings.Contains(name, `"`) {
				name = strings.ReplaceAll(name, `"`, `""`)
			}
			if strings.Contains(name, ",") {
				name = fmt.Sprintf(`"%s"`, name)
			}

			fmt.Fprintf(out, "%s", name)
			if i < fields-1 {
				fmt.Fprintf(out, ",")
			} else {
				fmt.Fprintf(out, "\n")
			}
		}
	}

	for i, f := range schema {
		value := f.GetValue()
		if strings.Contains(value, `"`) {
			value = strings.ReplaceAll(value, `"`, `""`)
		}
		if strings.Contains(value, ",") {
			value = fmt.Sprintf(`"%s"`, value)
		}

		fmt.Fprintf(out, "%s", value)
		if i < fields-1 {
			fmt.Fprintf(out, ",")
		} else {
			fmt.Fprintf(out, "\n")
		}
	}
}
