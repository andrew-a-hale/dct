package peek

import (
	"dct/cmd/utils"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"slices"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	defaultWriter = os.Stdout
	defaultLines  = 10
	lines         string
	linesParsed   int
	output        string
	writer        io.Writer
)

func init() {
	PeekCmd.Flags().StringVarP(&output, "output", "o", "", "write to output file")
	PeekCmd.Flags().StringVarP(&lines, "lines", "n", "", "number of lines to output")
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

		linesParsed := defaultLines
		if lines != "" {
			x, err := strconv.ParseInt(lines, 10, 0)
			if err != nil {
				log.Printf("Warning: failed to parse -n arg defaulting to %v\n", defaultLines)
			}
			linesParsed = int(x)
		}

		peek(file, linesParsed, writer)
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

func peek(file string, lines int, writer io.Writer) {
	sql := fmt.Sprintf("select * from '%s' limit %d", file, lines)
	utils.Run(sql, writer)
}
