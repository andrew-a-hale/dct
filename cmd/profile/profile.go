package profile

import (
	"dct/cmd/utils"
	"fmt"
	"io"
	"log"
	"maps"
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

		report, err := analyse(result)
		if err != nil {
			log.Fatalf("failed to analyse file: %v", err)
		}

		fmt.Fprintf(writer, "Profile\n%s", report)
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

func analyse(result utils.Result) (string, error) {
	var report string
	var err error
	for i := range len(result.Headers) {
		var col []string
		for _, row := range result.Rows {
			col = append(col, row[i])
		}
		report, err = analyseStringField(result.Headers[i].Name, col)
		if err != nil {
			return "", err
		}
	}

	return report, nil
}

func analyseStringField(header string, column []string) (string, error) {
	runeMap := make(map[rune]int)
	valueMap := make(map[string]int)
	for _, str := range column {
		for _, rune := range str {
			runeMap[rune]++
		}
		valueMap[str]++
	}

	valuesArr := utils.SortMap(valueMap, -1)
	values := "row: value -> count\n"
	for i, v := range valuesArr {
		values += fmt.Sprintf("%d: %v -> %d\n", i, v.X, v.Y)
	}
	var valueLengths []int
	for k := range maps.Keys(valueMap) {
		valueLengths = append(valueLengths, len(k))
	}

	runesArr := utils.SortMap(runeMap, -1)
	runes := "row: rune -> count\n"
	for i, v := range runesArr {
		runes += fmt.Sprintf("%d: %v -> %d\n", i, string(v.X), v.Y)
	}

	analysis := fmt.Sprintf("Field: %s\n", header)
	analysis += fmt.Sprintf("Value Occurrence\n%s\n", values)
	analysis += fmt.Sprint("Value Summary\n")
	analysis += fmt.Sprintf("%s\n\n", utils.Summarise(valueLengths))
	analysis += fmt.Sprintf("Char Occurrence\n%s", runes)

	return analysis, nil
}
