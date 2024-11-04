package peek

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
	defaultWriter       = os.Stdout
	defaultLines  int32 = 10
	lines         int32
	output        string
	writer        io.Writer
)

func init() {
	PeekCmd.Flags().StringVarP(&output, "output", "o", "", "write to output file")
	PeekCmd.Flags().Int32VarP(&lines, "lines", "n", 0, "number of lines to output")
}

var PeekCmd = &cobra.Command{
	Use:   "peek [FILE]",
	Short: "peek into a file",
	Long:  "",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
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
			log.Printf("Error: expected -n to be at least 1 defaulting to %v\n", defaultLines)
			lines = defaultLines
		}

		peek(file, lines, writer)
	},
}

func parseFileArg(args []string) string {
	if len(args) != 1 {
		log.Fatalf("Error: expected one file in args: %v\n", args)
	}

	filepath := args[0]
	file := path.Base(filepath)
	fileext := strings.ToLower(path.Ext(file))

	if slices.Contains(utils.SUPPORTED_FILETYPES, fileext) {
		return filepath
	}

	return ""
}

func peek(file string, lines int32, writer io.Writer) {
	query := fmt.Sprintf("select * from '%s' limit %d", file, lines)
	result, err := utils.Query(query)
	if err != nil {
		log.Fatalf("failed to cmp files: %v", err)
	}

	if output == "" {
		result.Render(writer, int(lines))
	} else {
		result.ToCsv(writer)
	}
}
