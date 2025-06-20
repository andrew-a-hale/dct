package infer

import (
	"dct/cmd/utils"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

var (
	defaultWriter     = os.Stdout
	defaultLines  int = 10
	lines         int
	output        string
	writer        io.Writer
	table         string
)

func init() {
	InferCmd.Flags().StringVarP(&table, "table", "t", "default", "Table name used in create table statement (default default)")
	InferCmd.Flags().StringVarP(&output, "output", "o", "", "Output to file (default stdout)")
	InferCmd.Flags().IntVarP(&lines, "lines", "n", 0, "Number of lines to infer schema from")
}

var InferCmd = &cobra.Command{
	Use:   "infer <file>",
	Short: "Infer sql schema for file",
	Long:  `Infer sql schema for file`,
	Args:  cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		file := parseFileArg(args)
		log.Printf("peeking at %s...\n", file)

		var err error
		writer := defaultWriter
		if output != "" {
			writer, err = os.Create(output)
			if err != nil {
				log.Printf("Warning: failed to create out file defaulting to %v\n", defaultWriter)
			}
		}

		if lines < 1 {
			log.Printf("Warning: expected -n to be at least 1 defaulting to %v\n", defaultLines)
			lines = defaultLines
		}

		infer(file, lines, table, writer)
	},
}

func parseFileArg(args []string) string {
	if len(args) != 1 {
		log.Fatalf("Error: expected one file in args: %v\n", args)
	}

	filepath := args[0]
	file := path.Base(filepath)
	fileext := strings.ToLower(path.Ext(file))

	if slices.Contains(utils.PEEK_SUPPORTED_FILETYPES, fileext) {
		return filepath
	}

	return ""
}

func infer(file string, lines int, table string, writer io.Writer) {
	query := fmt.Sprintf("select * from '%s' limit %d", file, lines)
	result, err := utils.Query(query)
	if err != nil {
		log.Fatalf("failed to cmp files: %v", err)
	}

	fmt.Fprintln(writer, result.ToSql(table))
}
