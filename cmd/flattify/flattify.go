package flattify

import (
	"bufio"
	"dct/cmd/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"maps"
	"os"
	"path"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type Tuple [2]any

var (
	defaultWriter = os.Stdout
	output        string
	sql           bool
	writer        io.Writer
)

func init() {
	FlattifyCmd.Flags().StringVarP(&output, "output", "o", "", "Output file path (default: stdout)")
	FlattifyCmd.Flags().BoolVarP(&sql, "sql", "s", false, "Generate DuckDB-compatible SQL statement")
}

var FlattifyCmd = &cobra.Command{
	Use:   "flattify [FILE]",
	Short: "Convert nested JSON to flat structure",
	Long:  `Recursively unnest JSON structures to a single layer, with optional SQL output for database use`,
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		file, ext := parseFileArg(args)

		var err error
		writer = defaultWriter
		if output != "" {
			writer, err = os.Create(output)
			if err != nil {
				log.Printf("Warning: failed to create out file defaulting to %v\n", defaultWriter)
			}
		}

		switch {
		case ext == utils.JSON && sql:
			flattify(file, writer, writeSelectStatement)
		case ext == utils.NDJSON && sql:
			flattifyLines(file, writer, writeMergedSelectStatement)
		case ext == utils.JSON && !sql:
			flattify(file, writer, writeJsonLine)
		case ext == utils.NDJSON && !sql:
			flattifyLines(file, writer, writeJsonLines)
		}
	},
}

func parseFileArg(args []string) (string, string) {
	if len(args) != 1 {
		log.Fatalf("Error: expected 1 arg: %v\n", args)
	}

	filepath := args[0]
	file := path.Base(filepath)
	fileext := strings.ToLower(path.Ext(file))

	if slices.Contains(utils.FLATTIFY_SUPPORTED_FILETYPES, fileext) {
		return filepath, fileext
	}

	return "", ""
}

func flattifyJson(obj any) []Tuple {
	var res []Tuple

	var _flatten func(any, string)
	_flatten = func(obj any, path string) {
		switch obj := obj.(type) {
		case nil:
		case string, float64, bool:
			res = append(res, Tuple{path, obj})
		case []any:
			for i, v := range obj {
				if path == "" {
					_flatten(v, strconv.Itoa(i+1))
				} else {
					_flatten(v, fmt.Sprintf("%s.%d", path, i+1))
				}
			}
		case map[string]any:
			for k, v := range obj {
				if path == "" {
					_flatten(v, k)
				} else {
					_flatten(v, fmt.Sprintf("%s.%s", path, k))
				}
			}
		}
	}

	_flatten(obj, "")
	return res
}

func flattify(file string, writer io.Writer, writerFunc func(any, string, io.Writer)) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	var j any
	err = json.Unmarshal(content, &j)
	if err != nil {
		log.Fatalf("failed to marshal json: %v", err)
	}

	obj := flattifyJson(j)
	writerFunc(obj, file, writer)
}

func flattifyLines(file string, writer io.Writer, writerFunc func(any, string, io.Writer)) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer f.Close()

	var lines [][]Tuple
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var j any
		err = json.Unmarshal([]byte(scanner.Text()), &j)
		if err != nil {
			log.Fatalf("failed to marshal json: %v", err)
		}

		lines = append(lines, flattifyJson(j))
	}

	writerFunc(lines, file, writer)
}

func writeJsonLine(line any, _ string, writer io.Writer) {
	fmt.Fprint(writer, "{")
	for i, v := range line.([]Tuple) {
		k := v[0]

		var val string
		switch v := v[1].(type) {
		case nil:
		case string:
			val = fmt.Sprintf(`"%s"`, v)
		case float64, bool:
			val = fmt.Sprintf(`%v`, v)
		}
		if i == len(line.([]Tuple))-1 {
			fmt.Fprintf(writer, `"%v":%v`, k, val)
		} else {
			fmt.Fprintf(writer, `"%v":%v,`, k, val)
		}
	}
	fmt.Fprint(writer, "}\n")
}

func writeJsonLines(lines any, _ string, writer io.Writer) {
	for _, line := range lines.([][]Tuple) {
		fmt.Fprint(writer, "{")
		for i, v := range line {
			k := v[0]

			var val string
			switch v := v[1].(type) {
			case nil:
			case string:
				val = fmt.Sprintf(`"%s"`, v)
			case float64, bool:
				val = fmt.Sprintf(`%v`, v)
			}
			if i == len(line)-1 {
				fmt.Fprintf(writer, `"%v":%v`, k, val)
			} else {
				fmt.Fprintf(writer, `"%v":%v,`, k, val)
			}
		}
		fmt.Fprint(writer, "}\n")
	}
}

func writeSelectStatement(line any, file string, writer io.Writer) {
	fmt.Fprint(writer, "select\n")
	paths := make(map[string]string)
	for _, v := range line.([]Tuple) {
		switch v[1].(type) {
		case nil:
		case string:
			paths[v[0].(string)] = "varchar"
		case float64:
			paths[v[0].(string)] = "decimal"
		case bool:
			paths[v[0].(string)] = "bool"
		}
	}

	var keys []string
	for key := range maps.Keys(paths) {
		keys = append(keys, key)
	}
	slices.Sort(keys)

	for i, k := range keys {
		fmt.Fprint(writer, "\t")
		if i == 0 {
			fmt.Fprintf(writer, `%s::%s`, k, paths[k])
		} else {
			fmt.Fprintf(writer, `, %s::%s`, k, paths[k])
		}
		fmt.Fprint(writer, "\n")
	}

	fmt.Fprintf(writer, "from '%s';\n", file)
}

func writeMergedSelectStatement(lines any, file string, writer io.Writer) {
	paths := make(map[string]string)
	for _, line := range lines.([][]Tuple) {
		for _, v := range line {
			key, ok := v[0].(string)
			if !ok {
				log.Fatalf("failed to create merged select statement: %v", v[0])
			}

			re := regexp.MustCompile(`.(\d+)`)
			path := string(re.ReplaceAll([]byte(key), []byte("[$1]")))
			switch v[1].(type) {
			case nil:
			case string:
				paths[path] = "varchar"
			case float64:
				paths[path] = "decimal"
			case bool:
				paths[path] = "bool"
			}
		}
	}

	var keys []string
	for key := range maps.Keys(paths) {
		keys = append(keys, key)
	}
	slices.Sort(keys)

	fmt.Fprint(writer, "select\n")
	for i, k := range keys {
		fmt.Fprint(writer, "\t")
		if i == 0 {
			fmt.Fprintf(writer, `%s::%s`, k, paths[k])
		} else {
			fmt.Fprintf(writer, `, %s::%s`, k, paths[k])
		}
		fmt.Fprint(writer, "\n")
	}

	fmt.Fprintf(writer, "from '%s';\n", file)
}
