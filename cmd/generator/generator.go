package generator

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"reflect"

	"github.com/spf13/cobra"
)

var (
	rawSchema string
	lines     int
	outfile   string
)

func init() {
	GenCmd.Flags().StringVarP(&outfile, "outfile", "o", "", "Output file path (default: stdout)")
	GenCmd.Flags().IntVarP(&lines, "lines", "n", 0, "Number of data rows to generate")
}

type (
	ctxKey   string
	Schema   []Field
	FieldMap map[string]int
)

const (
	SCHEMA_KEY    ctxKey = "schema"
	FIELD_MAP_KEY ctxKey = "fieldMap"
)

var GenCmd = &cobra.Command{
	Use:   "gen [schema]",
	Short: "Generate synthetic data",
	Long:  `Create realistic test data based on a schema definition with support for custom field types and derived fields`,
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		var out io.Writer
		var err error
		if outfile != "" {
			out, err = os.Create(outfile)
			if err != nil {
				log.Fatalf("failed to create out file: %v\n", err)
			}
		} else {
			out = os.Stdout
		}

		rawSchema = args[0]
		fieldMap := make(FieldMap)
		schema := parseSchema(rawSchema)
		for i, f := range schema {
			fieldMap[reflect.ValueOf(f).Elem().FieldByName("Field").String()] = i
		}

		ctx := context.Background()
		ctx = context.WithValue(ctx, SCHEMA_KEY, schema)
		ctx = context.WithValue(ctx, FIELD_MAP_KEY, fieldMap)
		for i := range lines {
			Write(ctx, out, i)
		}
	},
}

func parseInputSchema(schemaString string) any {
	sch := []byte(schemaString)

	file, err := os.Open(schemaString)
	if err == nil {
		sch, _ = io.ReadAll(file)
	}
	defer file.Close()

	var schema any
	err = json.Unmarshal(sch, &schema)
	if err != nil {
		log.Fatalf("Error: failed to parse metric config: %v\n", err)
	}

	return sch
}
