package generator

import (
	"context"
	"io"
	"log"
	"os"
	"reflect"

	"dct/cmd/utils"

	"github.com/spf13/cobra"
)

var (
	rawSchema string
	lines     int
	format    string
	outfile   string
	cache     utils.Cache = utils.NewCache()
)

func init() {
	GenCmd.Flags().StringVarP(&outfile, "outfile", "o", "", "Output file path (default: stdout)")
	GenCmd.Flags().StringVarP(&format, "format", "f", "csv", "Output format supports ndjson, csv")
	GenCmd.Flags().IntVarP(&lines, "lines", "n", 1, "Number of data rows to generate")
}

type (
	ctxKey   string
	Schema   []Field
	FieldMap map[string]int
)

const (
	FORMAT_KEY    ctxKey = "format"
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
		ctx = context.WithValue(ctx, FORMAT_KEY, "."+format)
		ctx = context.WithValue(ctx, SCHEMA_KEY, schema)
		ctx = context.WithValue(ctx, FIELD_MAP_KEY, fieldMap)
		Write(ctx, out, lines)
	},
}
