package generator

import (
	"fmt"
	"io"
	"log"
	"os"

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

		schema := parseSchema(rawSchema)
		for _, f := range schema {
			f.Generate(lines)
		}

		for i, r := range buildExport(schema, out).Rows {
			fmt.Printf("%d: %v\n", i, r)
		}
	},
}
