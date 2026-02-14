package chart

import (
	"fmt"
	"log"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
	"text/template"

	"dct/cmd/utils"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func init() {
	ChartCmd.Flags().IntVarP(&width, "width", "w", 0, "Width of the chart in characters")
}

const (
	DECIMAL_PLACES int = 2
	MINWIDTH       int = 20
)

var (
	texture              []byte = []byte("╍")
	divChar              []byte = []byte("┤")
	middleHorizontalChar []byte = []byte("─")
	middleVerticalChar   []byte = []byte("│")
	width                int
)

type Chart struct {
	Title        string
	XLabelWs     string
	TopBorder    string
	BottomBorder string
	Rows         []string
}

var ChartTemplate string = `
{{.XLabelWs}} ┌─{{.Title}}{{.TopBorder}}┐
{{range .Rows}} {{- println .}}{{end -}}
{{.XLabelWs}} └{{.BottomBorder}}┘
`

var ChartCmd = &cobra.Command{
	Use:   "chart [file] [colIndex]",
	Short: "Generate visualisations from data",
	Long:  `Create a simple ASCII bar chart from data file using specified column and aggregation function`,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		filename, colIndex := checkArgs(args)
		xs, ys := processAgg(filename, colIndex)
		draw(filename, xs, ys)
	},
}

func draw(filename string, xs []string, ys []int) {
	termWidth, _, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		termWidth = 50
	}

	xMaxLength := maxStringWidth(xs)
	minWidth := xMaxLength + MINWIDTH
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
		xlabels = append(xlabels, fmt.Sprintf("%s%s %s ", leading, x, divChar))
	}

	yMaxLength := int(math.Log10(float64(slices.Max(ys))))
	barWidth := (termWidth - xMaxLength - 3 - yMaxLength - 3 - 1)
	titleWidth := (termWidth - xMaxLength - 3 - 1)
	pixelValue := scale(ys, barWidth)
	title := makeTitle(filename, titleWidth)

	var rows []string
	for i, y := range ys {
		// label
		line := fmt.Sprint(xlabels[i])
		line += makeBar(y, pixelValue, texture)
		line += fmt.Sprintf(" %d ", y)

		// fill
		lineLength := strings.Count(line, "") - 1
		line += strings.Repeat(" ", termWidth-lineLength-1)

		// border
		line += string(middleVerticalChar)
		rows = append(rows, line)
	}

	chart := Chart{
		Title:        title,
		XLabelWs:     strings.Repeat(" ", xMaxLength),
		TopBorder:    strings.Repeat(string(middleHorizontalChar), termWidth-xMaxLength-1-1-strings.Count(title, "")-1),
		BottomBorder: strings.Repeat(string(middleHorizontalChar), termWidth-xMaxLength-1-1-1),
		Rows:         rows,
	}

	tmpl, err := template.New("chart").Parse(ChartTemplate)
	if err != nil {
		log.Fatalf("invalid chart template: %v\n", err)
	}

	err = tmpl.Execute(os.Stdout, chart)
	if err != nil {
		log.Fatalf("failed to render template: %v\n", err)
	}
}

func makeBar(x int, pixelValue float32, texture []byte) string {
	return strings.Repeat(string(texture), int(float32(x)*pixelValue))
}

func checkArgs(args []string) (filename string, colName int) {
	filename = args[0]
	_, err := os.Open(filename)
	if err != nil {
		log.Fatalf("file does not exist: %v\n", err)
	}

	colName, err = strconv.Atoi(args[1])
	if err != nil {
		log.Fatalf("failed to parse colIndex: %v\n", err)
	}

	return filename, colName
}

func processAgg(filename string, colIndex int) ([]string, []int) {
	// read file
	result, err := utils.Query(
		fmt.Sprintf(
			`select #%d, count(#%d) as agg from '%s' group by 1 order by agg desc`,
			colIndex,
			colIndex,
			filename,
		),
	)
	if err != nil {
		log.Fatalf("failed to process given file: %v", err)
	}

	var xs []string
	var ys []int
	for _, row := range result.Rows {
		xs = append(xs, fmt.Sprintf("%v", row[0]))
		y, ok := row[1].(int)
		if !ok {
			log.Fatalf("failed to parse aggregate column: %v", row[1])
		}
		ys = append(ys, y)
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
	return string(middleHorizontalChar)
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

func scale(xs []int, barWidth int) float32 {
	largest := slices.Max(xs)
	pixelValue := float32(barWidth) / float32(largest)
	return pixelValue
}
