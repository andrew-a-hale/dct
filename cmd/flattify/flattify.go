package flattify

import (
	"bytes"
	"dct/cmd/utils"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"maps"
	"os"
	"path"
	"slices"
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
	FlattifyCmd.Flags().BoolVarP(&sql, "sql", "s", false, "Generate DuckDB-compatible SQL statement select clause")
}

var FlattifyCmd = &cobra.Command{
	Use:   "flattify [FILE or JSON]",
	Short: "Convert nested JSON to flat structure",
	Long:  `Recursively unnest JSON structures to a single layer, with optional SQL output for database use`,
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		payload, ext, err := parseJsonArgs(args)
		if e, ok := err.(utils.UnsupportedFileTypeErr); ok {
			log.Fatalf("error: failed to parseJsonArgs: %v", e)
		}

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

		if sql {
			writeFromStatement(args[0], writer)
		}
	},
}

func parseJsonArgs(args []string) ([][]byte, string, error) {
	if len(args) != 1 {
		log.Fatalf("Error: expected 1 arg: %v\n", args)
	}

	filepath := args[0]
	file := path.Base(filepath)
	fileext := strings.ToLower(path.Ext(file))
	rawJson, err := os.ReadFile(filepath)
	if err != nil && !json.Valid([]byte(args[0])) {
		// not a file or json
		err := &utils.UnsupportedFileTypeErr{
			Msg:      "failed to read input json invalid file type",
			Filename: filepath,
			Ext:      fileext,
		}
		return [][]byte{}, "", err
	} else if err != nil {
		// valid json
		rawJson = []byte(args[0])
	}

	jsonType, err := detectJsonType(rawJson)
	if err != nil {
		log.Fatalf("%v", err)
	}

	var lines [][]byte
	switch jsonType {
	case utils.JSON:
		lines = append(lines, rawJson)
	case utils.NDJSON:
		jsonlines := bytes.Split(rawJson, []byte("\n"))
		for _, line := range jsonlines {
			lines = append(lines, line)
		}
	}

	return lines, jsonType, nil
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
					if sql {
						_flatten(v, fmt.Sprintf("json[%d]", i))
					} else {
						_flatten(v, fmt.Sprintf("$[%d]", i))
					}
				} else {
					_flatten(v, fmt.Sprintf("%s[%d]", path, i))
				}
			}
		case map[string]any:
			for k, v := range obj {
				if sql {
					if path == "" {
						_flatten(v, fmt.Sprintf(`json."%s"`, k))
					} else {
						_flatten(v, fmt.Sprintf(`%s."%s"`, path, k))
					}
				} else {
					k := fmt.Sprintf("['%s']", k)
					if path == "" {
						_flatten(v, "$"+k)
					} else {
						_flatten(v, fmt.Sprintf("%s%s", path, k))
					}
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
		if len(line) == 0 {
			break
		}

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
	fmt.Fprint(writer, "select\n")
	paths := make(map[string]string)
	for _, v := range line.([]Tuple) {
		path := v[0].(string)
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
	paths := make(map[string]string)
	for _, line := range lines.([][]Tuple) {
		for _, v := range line {
			path := v[0].(string)
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
		if !json.Valid(line) && len(line) > 0 {
			return "", fmt.Errorf("invalid json: %v", line)
		}
	}

	return utils.NDJSON, nil
}

func writeFromStatement(payload string, writer io.Writer) {
	if json.Valid([]byte(payload)) {
		fmt.Fprintf(writer, "from (select '%v'::json as json)\n", payload)
	} else if fs.ValidPath(payload) {
		fmt.Fprintf(writer, "from read_json_objects('%v', format='unstructured') as json\n", payload)
	} else {
		log.Fatal("unreachable")
	}
}
