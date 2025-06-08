package jsonschema2sql

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

var (
	defaultWriter = os.Stdout
	output        string
	writer        io.Writer
)

func init() {
	JsonSchema2SqlCmd.Flags().StringVarP(&output, "output", "o", "", "Output file path (default: stdout)")
}

var JsonSchema2SqlCmd = &cobra.Command{
	Use:   "jsonschema2sql",
	Short: "Generate a SQL table from JSON Schema",
	Long:  `Generate a SQL table from JSON Schema`,
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// read jsonschema
// create type map for json to ansi sql
//   regular types
//   array types
//   row types
//   search for references
// generate create table statement
