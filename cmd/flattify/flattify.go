package flattify

import (
	"bufio"
	"dct/cmd/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"slices"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type Tuple [2]any

var (
	defaultWriter = os.Stdout
	output        string
	writer        io.Writer
)

func init() {
	FlattifyCmd.Flags().StringVarP(&output, "output", "o", "", "write to output file")
}

var FlattifyCmd = &cobra.Command{
	Use:   "flattify [FILE] [...[TYPE]]",
	Short: "flattify json",
	Long:  "recursively unnest json to a single layer",
	Args:  cobra.MatchAll(cobra.ExactArgs(2), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		file, ext := parseFileArg(args)
		log.Fatal("todo fix args!")

		var err error
		writer = defaultWriter
		if output != "" {
			writer, err = os.Create(output)
			if err != nil {
				log.Printf("Warning: failed to create out file defaulting to %v\n", defaultWriter)
			}
		}

		writerFunc := writeJsonLine

		switch ext {
		case utils.JSON:
			flattify(file, writer, writerFunc)
		case utils.NDJSON:
			flattifyLines(file, writer, writerFunc)
		}
	},
}

func parseFileArg(args []string) (string, string) {
	if len(args) > 2 {
		log.Fatalf("Error: expected 2 or less args: %v\n", args)
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
					_flatten(v, strconv.Itoa(i))
				} else {
					_flatten(v, fmt.Sprintf("%s.%d", path, i))
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

func flattify(file string, writer io.Writer, writerFunc func([]Tuple, io.Writer)) {
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
	writeJsonLine(obj, writer)
}

func flattifyLines(file string, writer io.Writer, writerFunc func([]Tuple, io.Writer)) {
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

	for _, line := range lines {
		writerFunc(line, writer)
	}
}

func writeJsonLine(line []Tuple, writer io.Writer) {
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
