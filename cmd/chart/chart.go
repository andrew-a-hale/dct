package chart

import (
	"dct/cmd/utils"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	TEXTURE                []byte = []byte("╍")
	DIV_CHAR               []byte = []byte("┤")
	BOTTOM_LEFT_CHAR       []byte = []byte("└")
	BOTTOM_RIGHT_CHAR      []byte = []byte("┘")
	TOP_LEFT_CHAR          []byte = []byte("┌")
	TOP_RIGHT_CHAR         []byte = []byte("┐")
	MIDDLE_HORIZONTAL_CHAR []byte = []byte("─")
	MIDDLE_VERTICAL_CHAR   []byte = []byte("│")
)

var ChartCmd = &cobra.Command{
	Use:   "chart [file] [col-index (1-based)]",
	Short: "create a simple histogram chart",
	Long:  "",
	Args:  cobra.MatchAll(cobra.ExactArgs(2), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		xs, ys := process(args[0], args[1])
		title := fmt.Sprintf("histogram of '%s'", filepath.Base(args[0]))
		fmt.Println()
		draw(title, xs, ys)
		fmt.Println()
	},
}

func draw(title string, xs []string, ys []float64) {
	termWidth, _, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal("failed to get terminal size")
	}

	width := int(math.Ceil(float64(termWidth) / 2.0))

	xlabels := make([]string, len(xs))
	xMaxLength := utils.MaxStringWidth(xs) + 1
	for i, x := range xs {
		leading := strings.Repeat(" ", xMaxLength-strings.Count(x, "")+1)
		xlabels[i] = fmt.Sprintf("%s%s %s ", leading, x, DIV_CHAR)
	}

	ticks := 5
	yticks, step := utils.TickScale(ys, ticks)
	ypadding := (width - xMaxLength - 1) / 4
	var ylabels string
	for i, y := range yticks {
		if i > 0 {
			ylabels += strings.Repeat(" ", ypadding-strings.Count(y, "")+1)
		}
		ylabels += y
	}

	fmt.Printf(
		"%s %s%s %s %s%s\n",
		strings.Repeat(" ", xMaxLength),
		TOP_LEFT_CHAR,
		strings.Repeat(string(MIDDLE_HORIZONTAL_CHAR), 2),
		title,
		strings.Repeat(
			string(MIDDLE_HORIZONTAL_CHAR),
			width-xMaxLength-1-2-strings.Count(title, "")-1,
		),
		TOP_RIGHT_CHAR,
	)
	for i, y := range ys {
		line := fmt.Sprint(xlabels[i])
		bar := makeBar(y, ypadding, step, TEXTURE)
		line += bar
		line += fmt.Sprintf(" %.2f ", y)
		lineLength := strings.Count(line, "") - 1
		fmt.Printf(
			"%s%s %s\n",
			line,
			strings.Repeat(" ", width-lineLength),
			MIDDLE_VERTICAL_CHAR,
		)
	}
	fmt.Printf("%s %s%s%s\n",
		strings.Repeat(" ", xMaxLength),
		BOTTOM_LEFT_CHAR,
		strings.Repeat(string(MIDDLE_HORIZONTAL_CHAR), width-xMaxLength-1),
		BOTTOM_RIGHT_CHAR,
	)
}

func makeBar(x float64, stepWidth int, step int, texture []byte) string {
	pixelValue := float64(step) / float64(stepWidth)
	return strings.Repeat(string(texture), int(math.Ceil(float64(x)/pixelValue)))
}

func process(file string, colindex string) ([]string, []float64) {
	result, err := utils.Query(fmt.Sprintf(`select #%s, count(*) as cnt from '%s' group by 1 order by cnt desc`, colindex, file))
	if err != nil {
		log.Fatalf("failed to process given file: %v", err)
	}

	xs := make([]string, len(result.Rows))
	ys := make([]float64, len(result.Rows))
	for i, row := range result.Rows {
		xs[i] = row[0]
		ys[i], err = strconv.ParseFloat(row[1], 64)
		if err != nil {
			log.Fatalf("failed to parse aggregate column: %v", err)
		}
	}

	return xs, ys
}
