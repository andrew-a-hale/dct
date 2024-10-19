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

var (
	TEXTURE                []byte = []byte("╍")
	DIV_CHAR               []byte = []byte("┤")
	MIDDLE_HORIZONTAL_CHAR []byte = []byte("─")
	MIDDLE_VERTICAL_CHAR   []byte = []byte("│")
	DECIMAL_PLACES         int    = 2
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
	Use:   "chart [file] [col-index (1-based)]",
	Short: "create a simple histogram chart",
	Long:  "",
	Args:  cobra.MatchAll(cobra.ExactArgs(2), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		xs, ys := process(args[0], args[1])
		filename := filepath.Base(args[0])
		draw(filename, xs, ys)
	},
}

func draw(filename string, xs []string, ys []float64) {
	termWidth, _, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal("failed to get terminal size")
	}

	width := int(math.Ceil(float64(termWidth) / 4.0))
	xMaxLength := maxStringWidth(xs)
	yMaxLength := maxFloatWidth(ys, DECIMAL_PLACES)
	barWidth := (width - xMaxLength - 3 - yMaxLength - 3)
	titleWidth := (width - xMaxLength - 3 - 1)
	pixelValue := scale(ys, barWidth)
	title := makeTitle(filename, titleWidth)

	var xlabels []string
	for _, x := range xs {
		xLength := strings.Count(x, "") - 1
		leading := strings.Repeat(" ", xMaxLength-xLength)
		xlabels = append(xlabels, fmt.Sprintf("%s%s %s ", leading, x, DIV_CHAR))
	}

	var rows []string
	for i, y := range ys {
		// label
		line := fmt.Sprint(xlabels[i])
		line += makeBar(y, pixelValue, TEXTURE)
		line += fmt.Sprintf(fmt.Sprintf(" %%.%df ", DECIMAL_PLACES), y)

		// fill
		lineLength := strings.Count(line, "") - 1
		line += strings.Repeat(" ", width-lineLength-1)

		// border
		line += string(MIDDLE_VERTICAL_CHAR)
		rows = append(rows, line)
	}

	chart := Chart{
		Title:        title,
		XLabelWs:     strings.Repeat(" ", xMaxLength),
		TopBorder:    strings.Repeat(string(MIDDLE_HORIZONTAL_CHAR), width-xMaxLength-1-1-strings.Count(title, "")-1),
		BottomBorder: strings.Repeat(string(MIDDLE_HORIZONTAL_CHAR), width-xMaxLength-1-1-1),
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

func process(file string, colindex string) ([]string, []float64) {
	result, err := utils.Query(
		fmt.Sprintf(
			`select #%s, count(*) as cnt from '%s' group by 1 order by cnt desc`,
			colindex,
			file,
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
	formats := map[int]string{0: " histogram of '%s' ", 1: " hist of '%s' ", 2: " hist: '%s' ", 3: "hist"}

	var title string
	var length int
	for i := 0; i < len(formats); i++ {
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
	for largest > 1 {
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
