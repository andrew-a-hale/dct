package chart

import (
	"dct/cmd/utils"
	"fmt"
	"log"
	"strconv"

	"github.com/spf13/cobra"
)

var ChartCmd = &cobra.Command{
	Use:   "chart [file] [col-index (1-based)]",
	Short: "create a simple histogram chart",
	Long:  "",
	Args:  cobra.MatchAll(cobra.ExactArgs(2), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		data := process(args[0], args[1])
		draw(data)
	},
}

func draw(data map[string]int64) {
	fmt.Println(data)
	// lipgloss
	// regions: y-axis, title, bars and x-axis
}

func process(file string, colindex string) map[string]int64 {
	result, err := utils.Query(fmt.Sprintf(`select #%s, count(*) as cnt from '%s' group by 1`, colindex, file))
	if err != nil {
		log.Fatalf("failed to process given file: %v", err)
	}

	data := make(map[string]int64)
	for _, row := range result.Rows {
		key := row[0]
		value := row[1]
		data[key], err = strconv.ParseInt(value, 10, 0)
		if err != nil {
			log.Fatalf("failed to parse aggregate column: %v", err)
		}
	}

	return data
}
