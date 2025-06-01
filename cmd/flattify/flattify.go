package flattify

import (
	"bytes"
	"dct/cmd/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"maps"
	"os"
	"path"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode"

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
	FlattifyCmd.Flags().BoolVarP(&sql, "sql", "s", false, "Generate DuckDB-compatible SQL statement select clause")
}

var FlattifyCmd = &cobra.Command{
	Use:   "flattify [FILE or JSON]",
	Short: "Convert nested JSON to flat structure",
	Long:  `Recursively unnest JSON structures to a single layer, with optional SQL output for database use`,
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		payload, ext := parseJsonArgs(args)

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
			flattify(payload, writer, writeSelectStatement)
		case ext == utils.NDJSON && sql:
			flattifyLines(payload, writer, writeMergedSelectStatement)
		case ext == utils.JSON && !sql:
			flattify(payload, writer, writeJson)
		case ext == utils.NDJSON && !sql:
			flattifyLines(payload, writer, writeJsonLines)
		}
	},
}

func parseJsonArgs(args []string) ([][]byte, string) {
	if len(args) != 1 {
		log.Fatalf("Error: expected 1 arg: %v\n", args)
	}

	rawJson := []byte(args[0])
	if fs.ValidPath(args[0]) && !json.Valid([]byte(args[0])) {
		filepath := args[0]
		file := path.Base(filepath)
		fileext := strings.ToLower(path.Ext(file))

		f, err := os.Open(filepath)
		if err != nil {
			log.Fatalf("failed to open file: %v", err)
		}
		defer f.Close()

		utils.Assert(
			slices.Contains(utils.FLATTIFY_SUPPORTED_FILETYPES, fileext),
			"unsupported file type: "+fileext,
		)

		rawJson, err = io.ReadAll(f)
		if err != nil {
			log.Fatalf("failed to read file: %v", err)
		}
	}

	var lines [][]byte
	var jsonType string
	jsonType, err := detectJsonType(rawJson)
	if err != nil {
		log.Fatalf("%v", err)
	}

	switch jsonType {
	case utils.JSON:
		lines = append(lines, rawJson)
	case utils.NDJSON:
		jsonlines := bytes.Split(rawJson, []byte("\n"))
		for _, line := range jsonlines {
			lines = append(lines, line)
		}
	}

	return lines, jsonType
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

func flattify(payload [][]byte, writer io.Writer, writerFunc func(any, io.Writer)) {
	var j any
	json.Unmarshal(payload[0], &j)
	obj := flattifyJson(j)
	writerFunc(obj, writer)
}

func flattifyLines(payload [][]byte, writer io.Writer, writerFunc func(any, io.Writer)) {
	var lines [][]Tuple
	for _, line := range payload {
		var j any
		json.Unmarshal(line, &j)
		lines = append(lines, flattifyJson(j))
	}
	writerFunc(lines, writer)
}

func writeJson(line any, writer io.Writer) {
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

func writeJsonLines(lines any, writer io.Writer) {
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

func writeSelectStatement(line any, writer io.Writer) {
	re := regexp.MustCompile(`\.?(\d+)`)

	fmt.Fprint(writer, "select\n")
	paths := make(map[string]string)
	for _, v := range line.([]Tuple) {
		key, ok := v[0].(string)
		if !ok {
			log.Fatalf("failed to create merged select statement: %v", v[0])
		}

		var path string
		if unicode.IsNumber(rune(key[0])) {
			path += "json"
		}

		path += string(re.ReplaceAll([]byte(key), []byte("[$1]")))
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
}

func writeMergedSelectStatement(lines any, writer io.Writer) {
	re := regexp.MustCompile(`\.?(\d+)`)

	paths := make(map[string]string)
	for _, line := range lines.([][]Tuple) {
		for _, v := range line {
			key, ok := v[0].(string)
			if !ok {
				log.Fatalf("failed to create merged select statement: %v", v[0])
			}

			var path string

			// edge case where key in top layer is an integer
			if unicode.IsNumber(rune(key[0])) {
				path += fmt.Sprintf(`json."%s"`, key[:strings.Index(key, ".")])
				key = key[strings.Index(key, "."):]
			}

			path += string(re.ReplaceAll([]byte(key), []byte("[$1]")))
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
}

func detectJsonType(content []byte) (string, error) {
	lines := bytes.Split(content, []byte("\n"))
	for _, line := range lines {
		_, ok := bytes.CutPrefix(line, []byte("{"))
		if !ok && json.Valid(content) {
			return utils.JSON, nil
		}
		_, ok = bytes.CutSuffix(line, []byte("}"))
		if !ok && json.Valid(content) {
			return utils.JSON, nil
		}
		_, ok = bytes.CutSuffix(line, []byte(","))
		if ok && json.Valid(content) {
			return utils.JSON, nil
		}
	}

	for _, line := range lines {
		if !json.Valid(line) {
			return "", errors.New("invalid json")
		}
	}

	return utils.NDJSON, nil
}
