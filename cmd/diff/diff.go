package diff

import (
	"dct/cmd/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	defaultWriter = os.Stdout
	output        string
	writer        io.Writer
	metrics       string
	all           bool
)

type key struct {
	left  string
	right string
	alias string
}

type keySpec struct {
	keys []key
}

type Metric struct {
	Agg   string `json:"agg"`
	Left  string `json:"left"`
	Right string `json:"right,omitempty"`
}

type metricSpec struct {
	Metrics []Metric `json:"metrics"`
}

func init() {
	DiffCmd.Flags().StringVarP(&output, "output", "o", "", "Output comparison to file")
	DiffCmd.Flags().StringVarP(&metrics, "metrics", "m", "",
		`Metrics specification for comparison, using JSON format:
  {{"agg": "mean", "left": "a", "right": "b"}, {"agg": "count_distinct", "left": "c"}}
  Supported aggregations: mean, median, min, max, count_distinct`)

	DiffCmd.Flags().BoolVarP(&all, "all", "a", false, "Show all rows, not just differences")
}

var DiffCmd = &cobra.Command{
	Use:   "diff [keys] [file1] [file2]",
	Short: "Compare files with key matching",
	Long: `Compare two files using key matching and metric calculations. 
	Specify keys in format: left_key[=right_key] (comma-separated for multiple keys)
	Use --metrics to define comparison metrics and --all to show all differences`,
	Args: cobra.MatchAll(cobra.ExactArgs(3), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		keys, left, right := parseArgs(args)

		var err error
		writer := defaultWriter
		if output != "" {
			writer, err = os.Create(output)
			if err != nil {
				log.Printf("Warning: failed to create out file defaulting to %v\n", defaultWriter)
			}
		}

		metricConf := metricSpec{}
		if metrics != "" {
			metricConf = parseMetrics(metrics)
		}

		diff(keys, left, right, metricConf, writer)
	},
}

func parseArgs(args []string) (keys keySpec, left, right string) {
	return parseKeys(args[0]), args[1], args[2]
}

func parseKeys(keyString string) keySpec {
	parts := strings.Split(keyString, ",")
	var keys []key

	for i, part := range parts {
		segments := strings.Split(part, "=")
		var k key
		switch s := len(segments); s {
		case 1:
			k = key{left: segments[0], right: segments[0], alias: segments[0]}
		case 2:
			k = key{left: segments[0], right: segments[1], alias: segments[0]}
		default:
			log.Fatalf("Error: malformed keys at %d: %s\n", i, keyString)
		}

		keys = append(keys, k)
	}

	return keySpec{keys: keys}
}

func parseMetrics(metricString string) metricSpec {
	metrics := []byte(metricString)

	file, err := os.Open(metricString)
	if err == nil {
		metrics, _ = io.ReadAll(file)
	}
	defer file.Close()

	var spec metricSpec
	err = json.Unmarshal(metrics, &spec)
	if err != nil {
		log.Fatalf("Error: failed to parse metric config: %v\n", err)
	}

	return spec
}

func generateKeySql(spec keySpec) (left, right string) {
	for i, key := range spec.keys {
		left += key.left
		right += fmt.Sprintf("%s as %s", key.right, key.left)
		if i < len(spec.keys)-1 {
			left += ", "
			right += ", "
		}
	}

	return left, right
}

func generateMetricSql(spec metricSpec) (left, right, main, check string) {
	if len(spec.Metrics) == 0 {
		return
	}

	for i, metric := range spec.Metrics {
		if metric.Agg == utils.COUNT_DISTINCT {
			left += fmt.Sprintf("count(distinct %s) as l_%s_count_distinct", metric.Left, metric.Left)

			if metric.Right == "" {
				metric.Right = metric.Left
			}
			right += fmt.Sprintf("count(distinct %s) as r_%s_count_distinct", metric.Right, metric.Left)
		} else {
			left += fmt.Sprintf("%s(%s) as l_%s_%s", metric.Agg, metric.Left, metric.Left, metric.Agg)

			if metric.Right == "" {
				metric.Right = metric.Left
			}
			right += fmt.Sprintf("%s(%s) as r_%s_%s", metric.Agg, metric.Right, metric.Left, metric.Agg)
		}

		main += fmt.Sprintf(
			"l_%s_%s, r_%s_%s, coalesce(l_%s_%s = r_%s_%s, false) as %s_%s_eq_flag",
			metric.Left,
			metric.Agg,
			metric.Left,
			metric.Agg,
			metric.Left,
			metric.Agg,
			metric.Left,
			metric.Agg,
			metric.Left,
			metric.Agg,
		)

		if !all {
			check += fmt.Sprintf(" or l_%s_%s != r_%s_%s", metric.Left, metric.Agg, metric.Left, metric.Agg)
		}

		if i < len(spec.Metrics)-1 {
			left += ", "
			right += ", "
			main += ", "
		}
	}

	return left, right, main, check
}

func generateSql(keys keySpec, left, right string, metrics metricSpec) string {
	leftKeys, rightKeys := generateKeySql(keys)
	leftMetrics, rightMetrics, mainMetrics, checkMetrics := generateMetricSql(metrics)
	leftSql := fmt.Sprintf("select %s, count(*) as l_cnt, %s from '%s' group by all", leftKeys, leftMetrics, left)
	rightSql := fmt.Sprintf("select %s, count(*) as r_cnt, %s from '%s' group by all", rightKeys, rightMetrics, right)

	sql := fmt.Sprintf(
		`with file1 as (
  %s
), file2 as (
  %s
)
select %s, l_cnt, r_cnt, coalesce(l_cnt = r_cnt, false) as cnt_eq_flag, %s
from file1
full join file2 using (%s)
where l_cnt <> r_cnt %s
order by cnt_eq_flag`,
		leftSql,
		rightSql,
		leftKeys,
		mainMetrics,
		leftKeys,
		checkMetrics,
	)

	return sql
}

func diff(keys keySpec, left string, right string, metrics metricSpec, writer io.Writer) {
	leftHasRows, err := utils.CheckFileHasRows(left)
	if err != nil {
		log.Fatalf("failed to check file: %v", err)
	}

	rightHasRows, err := utils.CheckFileHasRows(right)
	if err != nil {
		log.Fatalf("failed to check file: %v", err)
	}

	if !leftHasRows || !rightHasRows {
		log.Fatal("attempted to diff when least one of the files have no data")
	}

	query := generateSql(keys, left, right, metrics)
	result, err := utils.Query(query)
	if err != nil {
		log.Fatalf("failed to cmp files: %v", err)
	}

	if output == "" {
		maxRows := 5
		result.Render(writer, maxRows)
	} else {
		result.ToCsv(writer)
	}
}