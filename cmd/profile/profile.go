package profile

import (
	"dct/cmd/utils"
	"fmt"
	"io"
	"log"
	"maps"
	"math"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

var (
	defaultWriter = os.Stdout
	output        string
	writer        io.Writer
	topK          = 10000
)

func init() {
	ProfileCmd.Flags().StringVarP(&output, "output", "o", "", "Output file path (default: stdout)")
}

var ProfileCmd = &cobra.Command{
	Use:   "prof [FILE]",
	Short: "Analyse fields of data file.",
	Long:  `Analyse fields of data file to find edge cases`,
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		file := parseFileArg(args)

		var err error
		writer = defaultWriter
		if output != "" {
			writer, err = os.Create(output)
			if err != nil {
				log.Printf("Warning: failed to create out file defaulting to %v\n", defaultWriter)
			}
		}

		query := fmt.Sprintf("select * from '%s'", file)
		result, err := utils.Query(query)
		if err != nil {
			log.Fatalf("failed to read file: %v", err)
		}

		analyse(result, writer)
	},
}

func parseFileArg(args []string) string {
	if len(args) != 1 {
		log.Fatalf("Error: expected one file in args: %v\n", args)
	}

	filepath := args[0]
	file := path.Base(filepath)
	fileext := strings.ToLower(path.Ext(file))

	if slices.Contains(utils.PROFILE_SUPPORTED_FILETYPES, fileext) {
		return filepath
	}

	return ""
}

func analyse(result utils.Result, writer io.Writer) {
	fmt.Fprintln(writer, "-- PROFILE -- ")

	for i := range len(result.Headers) {
		var col []string
		for _, row := range result.Rows {
			col = append(col, row[i])
		}
		// writes directly to ouput
		analyseField(result.Headers[i].Name, col, writer)
	}

}

func analyseField(header string, column []string, writer io.Writer) {
	fmt.Fprintf(writer, "-- Field: `%s` -- \n", header)

	valueMap := make(map[string]int)
	for _, str := range column {
		valueMap[str]++
	}

	fmt.Fprintf(writer, "Count: %d\nUnique Count: %d\n\n", len(column), len(valueMap))
	fmt.Fprintln(writer, "Value Occurrence")

	i := 0
	// mostly unique values, just sample 10
	if len(valueMap) >= len(column)>>1 {
		fmt.Fprintln(writer, "MOSTLY UNIQUE VALUES JUST SHOWING SAMPLE OF 10")
		fmt.Fprintln(writer, "row: value -> count")
		for k, v := range valueMap {
			if i > 10 {
				break
			}
			fmt.Fprintf(writer, "%d: %v -> %d\n", i, k, v)
			i++
		}
	} else {
		fmt.Fprintln(writer, "row: value -> count")
		for k, v := range valueMap {
			fmt.Fprintf(writer, "%d: %v -> %d\n", i, k, v)
			i++
		}
	}
	fmt.Fprintln(writer)

	fmt.Fprint(writer, "Value Summary - String Lengths\n")
	fmt.Fprintf(writer, "%s\n\n", Summarise(valueMap))

	runeMap := make(map[rune]int)
	for k := range maps.Keys(valueMap) {
		for _, r := range k {
			runeMap[r]++
		}
	}

	runes := "row: rune -> count\n"
	leading := int(math.Ceil(math.Log10(float64(len(runeMap)))))

	i = 0
	for k, v := range runeMap {
		runes += fmt.Sprintf(
			"%0*d: %[3]q (hex: %[3]U) (dec: %[3]d) -> %[4]d\n",
			leading, i, k, v,
		)
		i++
	}

	fmt.Fprintf(writer, "Char Occurrence\n%s\n", runes)
	fmt.Fprintf(writer, "Char Analysis\n%s\n\n", AnalyseRunes(runeMap))
}
