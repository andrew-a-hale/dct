package chart

import (
	"dct/cmd/utils"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func init() {
	ChartCmd.Flags().Int32VarP(&width, "width", "w", 0, "Width of the chart in characters")
}

var (
	TEXTURE                []byte = []byte("╍")
	DIV_CHAR               []byte = []byte("┤")
	MIDDLE_HORIZONTAL_CHAR []byte = []byte("─")
	MIDDLE_VERTICAL_CHAR   []byte = []byte("│")
	DECIMAL_PLACES         int    = 2
	width                  int32
)

type Chart struct {
	Title        string
	XLabelWs     string
	TopBorder    string
	BottomBorder string
	Rows         []string
}

var CHART_TEMPLATE string = `
{{.XLabelWs}} ┌─{{.Title}}{{.TopBorder}}┐
{{range .Rows}} {{- println .}}{{end -}}
{{.XLabelWs}} └{{.BottomBorder}}┘
`

var ChartCmd = &cobra.Command{
	Use:   "chart [file] [col-name] [agg]",
	Short: "Generate visualisations from data",
	Long:  `Create a simple ASCII bar chart from data file using specified column and aggregation function`,
	Args:  cobra.MinimumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		checkArgs(args)
		xs, ys := processAgg(args)
		filename := filepath.Base(args[0])
		draw(filename, xs, ys)
	},
}

func draw(filename string, xs []string, ys []float64) {
	width := int(width)
	termWidth, _, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalf("failed to get terminal size: %v\n", err)
	}

	xMaxLength := maxStringWidth(xs)
	minWidth := xMaxLength + 20
	if termWidth < minWidth && width > 0 {
		log.Fatal("terminal width is too small to render chart\n")
	}
	if width < minWidth {
		fmt.Printf("provided width is too small, defaulting to %d\n", minWidth)
	}
	termWidth = max(width, minWidth)

	var xlabels []string
	for _, x := range xs {
		xLength := strings.Count(x, "") - 1
		leading := strings.Repeat(" ", xMaxLength-xLength)
		xlabels = append(xlabels, fmt.Sprintf("%s%s %s ", leading, x, DIV_CHAR))
	}

	yMaxLength := maxFloatWidth(ys, DECIMAL_PLACES)
	barWidth := (termWidth - xMaxLength - 3 - yMaxLength - 3)
	titleWidth := (termWidth - xMaxLength - 3 - 1)
	pixelValue := scale(ys, barWidth)
	title := makeTitle(filename, titleWidth)

	var rows []string
	for i, y := range ys {
		// label
		line := fmt.Sprint(xlabels[i])
		line += makeBar(y, pixelValue, TEXTURE)
		line += fmt.Sprintf(fmt.Sprintf(" %%.%df ", DECIMAL_PLACES), y)

		// fill
		lineLength := strings.Count(line, "") - 1
		line += strings.Repeat(" ", termWidth-lineLength-1)

		// border
		line += string(MIDDLE_VERTICAL_CHAR)
		rows = append(rows, line)
	}

	chart := Chart{
		Title:        title,
		XLabelWs:     strings.Repeat(" ", xMaxLength),
		TopBorder:    strings.Repeat(string(MIDDLE_HORIZONTAL_CHAR), termWidth-xMaxLength-1-1-strings.Count(title, "")-1),
		BottomBorder: strings.Repeat(string(MIDDLE_HORIZONTAL_CHAR), termWidth-xMaxLength-1-1-1),
		Rows:         rows,
	}

	tmpl, err := template.New("chart").Parse(CHART_TEMPLATE)
	if err != nil {
		log.Fatalf("invalid chart template: %v\n", err)
	}

	err = tmpl.Execute(os.Stdout, chart)
	if err != nil {
		log.Fatalf("failed to render template: %v\n", err)
	}
}

func makeBar(x float64, pixelValue float64, texture []byte) string {
	return strings.Repeat(string(texture), int(math.Ceil(x*pixelValue)))
}

func checkArgs(args []string) (filename, colName, agg string) {
	filename = args[0]
	_, err := os.Open(filename)
	if err != nil {
		log.Fatalf("file does not exist: %v\n", err)
	}

	colName = args[1]

	agg = args[2]
	utils.CheckAgg(agg)

	return filename, colName, agg
}

func processAgg(args []string) ([]string, []float64) {
	// format custom agg sql
	var agg string
	colName := args[1]

	switch args[2] {
	case utils.COUNT_DISTINCT:
		agg = fmt.Sprintf("count(distinct #%s)", colName)
	default:
		agg = fmt.Sprintf("%s(#%s)", args[2], colName)
	}

	// read file
	result, err := utils.Query(
		fmt.Sprintf(
			`select #%s, %s as agg from '%s' group by 1 order by agg desc`,
			colName,
			agg,
			args[0],
		),
	)
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

func makeTitle(filename string, width int) string {
	width -= 2 // padding
	formats := map[int]string{0: " histogram of '%s' ", 1: " hist of '%s' ", 2: " hist: '%s' ", 3: "'%s'"}

	var title string
	var length int
	for i := range len(formats) {
		title = fmt.Sprintf(formats[i], filename)
		length = strings.Count(title, "") - 1
		if length <= width {
			return title
		}
	}

	// skip title
	return string(MIDDLE_HORIZONTAL_CHAR)
}

func maxStringWidth(xs []string) int {
	m := 0
	for _, x := range xs {
		if strings.Count(x, "") > m {
			m = strings.Count(x, "")
		}
	}

	return m - 1
}

func maxFloatWidth(xs []float64, places int) int {
	largest := slices.Max(xs)
	digits := 1 + places
	for largest >= 1 {
		largest /= 10
		digits++
	}
	return digits
}

func scale(xs []float64, barWidth int) float64 {
	largest := slices.Max(xs)
	pixelValue := (float64(barWidth) / largest)
	return pixelValue
}
