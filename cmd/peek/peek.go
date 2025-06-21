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
	defaultWriter     = os.Stdout
	defaultLines  int = 10
	lines         int
	output        string
	writer        io.Writer
)

func init() {
	PeekCmd.Flags().StringVarP(&output, "output", "o", "", "Output to file instead of stdout")
	PeekCmd.Flags().IntVarP(&lines, "lines", "n", 0, "Number of lines to display")
}

var PeekCmd = &cobra.Command{
	Use:   "peek <file>",
	Short: "Preview file contents",
	Long:  `Display the first few lines of a data file to quickly inspect its structure and content`,
	Args:  cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		file := parseFileArg(args)

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

	if slices.Contains(utils.PEEK_SUPPORTED_FILETYPES, fileext) {
		return filepath
	}

	return ""
}

func peek(file string, lines int, writer io.Writer) {
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
