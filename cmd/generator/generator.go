package generator

import (
	"context"
	"io"
	"log"
	"os"
	"reflect"

	"github.com/spf13/cobra"
)

var (
	seed      int
	rawSchema string
	lines     int
	outfile   string
)

func init() {
	GenCmd.Flags().StringVarP(&rawSchema, "schema", "s", "", "schema to generate")
	GenCmd.Flags().StringVarP(&outfile, "outfile", "o", "", "output file")
	GenCmd.Flags().IntVarP(&lines, "lines", "n", 0, "lines to generate")

	GenCmd.MarkFlagRequired("schema")
	GenCmd.MarkFlagRequired("lines")
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
	Use:   "gen -s [schema] -n [lines] -o [outfile]",
	Short: "generate dummy data",
	Long:  `generate dummy data`,
	Args:  nil,
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
