package generator

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
)

func Write(ctx context.Context, out io.Writer, line int) {
	schema, ok := ctx.Value(SCHEMA_KEY).(Schema)
	if !ok {
		log.Fatalln("failed to read schema from context")
	}

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
		value := f.Generate(ctx)
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
