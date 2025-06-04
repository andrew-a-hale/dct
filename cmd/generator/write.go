package generator

import (
	"context"
	"dct/cmd/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
)

func Write(ctx context.Context, out io.Writer, lines int) {
	schema, ok := ctx.Value(SCHEMA_KEY).(Schema)
	if !ok {
		log.Fatalln("failed to read schema from context")
	}

	format, ok := ctx.Value(FORMAT_KEY).(string)
	if !ok {
		log.Fatalln("failed to read format from context")
	}

	switch format {
	case utils.NDJSON:
		writeJson(ctx, out, schema, lines)
	case utils.CSV:
		writeCsv(ctx, out, schema, lines)
	}
}

func writeCsv(ctx context.Context, out io.Writer, schema Schema, lines int) {
	fields := len(schema)

	for i := range lines {
		// headers
		if i == 0 {
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
}

func writeJson(ctx context.Context, out io.Writer, schema Schema, lines int) {
	for range lines {
		fmt.Fprint(out, "{")
		for i, f := range schema {
			value := f.Generate(ctx)

			switch f.GetType() {
			case FLOAT, INT, BOOL:
				fmt.Fprintf(out, `"%s":%s`, f.GetName(), value)
			case STRING:
				v, err := json.Marshal(value)
				if err != nil {
					log.Fatalf("failed to write `%s: %v` as json: %v", f.GetName(), value, err)
				}
				fmt.Fprintf(out, `"%s":%s`, f.GetName(), v)
			}

			if i != len(schema)-1 {
				fmt.Fprint(out, ",")
			}
		}

		fmt.Fprintln(out, "}")
	}
}
