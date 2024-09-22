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
	DiffCmd.Flags().StringVarP(&output, "output", "o", "", "write to output file")
	DiffCmd.Flags().StringVarP(&metrics, "metrics", "m", "",
		`custom metrics to compare. metrics must be supplied in the following format
  {{"agg": "mean", "left": "a", "right": "b"}, {"agg": "count_distinct", "left": "c"}}
  supported agg are: mean, median, min, max, count_distinct`)

	DiffCmd.Flags().BoolVarP(&all, "all", "a", false, "show all rows")
}

var DiffCmd = &cobra.Command{
	Use:   "diff [keySpec] [file1] [file2]",
	Short: "diff two files",
	Long: `aggregate to files and compare.
keys must be supplied in the following format
  {{"left": "a", "right": "b", "alias": "a"}, {"left": "c"", "right": "c", "alias": "c"}}`,
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
	keys = parseKeys(args[0])

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
		left += fmt.Sprintf("%s", key.left)
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
			left += fmt.Sprintf("count(distinct %s) as l_%s", metric.Left, metric.Left)

			if metric.Right == "" {
				metric.Right = metric.Left
			}
			right += fmt.Sprintf("count(distinct %s) as r_%s", metric.Right, metric.Left)
		} else {
			left += fmt.Sprintf("%s(%s) as l_%s", metric.Agg, metric.Left, metric.Left)

			if metric.Right == "" {
				metric.Right = metric.Left
			}
			right += fmt.Sprintf("%s(%s) as r_%s", metric.Agg, metric.Right, metric.Left)
		}
		main += fmt.Sprintf(
			"l_%s, r_%s, l_%s = r_%s as %s_eq_flag",
			metric.Left,
			metric.Left,
			metric.Left,
			metric.Left,
			metric.Left,
		)

		if !all {
			check += fmt.Sprintf(" and l_%s != r_%s", metric.Left, metric.Left)
		}

		if i < len(spec.Metrics)-1 {
			left += ", "
			right += ", "
		}
	}

	return left, right, main, check
}

func generateSql(keys keySpec, left, right string, metrics metricSpec) string {
	leftKeys, rightKeys := generateKeySql(keys)
	leftMetrics, rightMetrics, mainMetrics, checkMetrics := generateMetricSql(metrics)
	leftSql := fmt.Sprintf("select %s, count(*) as a_cnt, %s from '%s' group by all", leftKeys, leftMetrics, left)
	rightSql := fmt.Sprintf("select %s, count(*) as b_cnt, %s from '%s' group by all", rightKeys, rightMetrics, right)

	sql := fmt.Sprintf(`with file1 as (
  %s
), file2 as (
  %s
)
select %s, a_cnt, b_cnt, a_cnt = b_cnt as cnt_eq_flag, %s
from file1
full join file2 using (%s)
where 1=1 %s
order by cnt_eq_flag`, leftSql, rightSql, leftKeys, mainMetrics, leftKeys, checkMetrics)

	return sql
}

func diff(keys keySpec, left string, right string, metrics metricSpec, writer io.Writer) {
	sql := generateSql(keys, left, right, metrics)
	utils.Run(sql, writer)
}
