package generator

import (
	"io"
	"log"
	"os"
	"reflect"

	"github.com/spf13/cobra"
)

var (
	seed      int64
	rawSchema string
	lines     int
	outfile   string
)

func init() {
	GenCmd.Flags().StringVarP(&rawSchema, "schema", "s", "", "schema to generate")
	GenCmd.Flags().StringVarP(&outfile, "outfile", "o", "", "output file")
	GenCmd.Flags().Int64VarP(&seed, "seed", "S", 0, "fixed seed")
	GenCmd.Flags().IntVarP(&lines, "lines", "n", 0, "lines to generate")
}

var GenCmd = &cobra.Command{
	Use:   "gen -S -s [schema] -n [lines] -o [outfile]",
	Short: "generate dummy data",
	Long:  `generate dummy data`,
	Args:  nil,
	Run: func(cmd *cobra.Command, args []string) {
		var out io.Writer
		var err error
		if outfile != "" {
			out, err = os.Open(outfile)
			if err != nil {
				log.Fatalf("failed to create out file: %v\n", err)
			}
		} else {
			out = os.Stdout
		}

		fieldMap := make(map[string]int)
		schema := parseSchema(rawSchema)
		for i, f := range schema {
			if f.GetSource() == "derived" {
				continue
			}
			fieldMap[reflect.ValueOf(f).Elem().FieldByName("Field").String()] = i
			f.Generate(lines, &schema, &fieldMap)
		}

		for _, f := range schema {
			if f.GetSource() == "derived" {
				f.Generate(lines, &schema, &fieldMap)
			}
		}

		Export(schema, out)
	},
}
