package generator

import (
	"encoding/json"
	"io"
	"io/fs"
	"log"
	"math"
	"math/rand/v2"
	"os"
	"reflect"

	"github.com/spf13/cobra"
)

var (
	seed      int64
	rawSchema string
	lines     int
	out       string
)

func init() {
	GenCmd.Flags().StringVarP(&rawSchema, "schema", "s", "", "schema to generate")
	GenCmd.Flags().StringVarP(&out, "outfile", "o", "", "output file")
	GenCmd.Flags().Int64VarP(&seed, "seed", "S", 0, "fixed seed")
	GenCmd.Flags().IntVarP(&lines, "lines", "n", 0, "lines to generate")
}

var GenCmd = &cobra.Command{
	Use:   "gen -S -s [schema] -n [lines] -o [outfile]",
	Short: "generate dummy data",
	Long:  `generate dummy data`,
	Args:  nil,
	Run: func(cmd *cobra.Command, args []string) {
		schema := parseSchema(rawSchema)
		for _, f := range schema {
			if field, ok := f.(*RandomAsciiField); ok {
				field.Generate(lines)
			} else if field, ok := f.(*RandomUniformIntField); ok {
				field.Generate(lines)
			} else if field, ok := f.(*RandomNormalField); ok {
				field.Generate(lines)
			}
		}
	},
}

func parseSchema(rawSchema string) []Field {
	var schema []byte
	if fs.ValidPath(rawSchema) {
		f, err := os.Open(rawSchema)
		if err != nil {
			log.Fatalf("failed to open schema file: %v\n", err)
		}
		defer f.Close()

		schema, err = io.ReadAll(f)
		if err != nil {
			log.Fatalf("failed to open schema file: %v\n", err)
		}
	}

	var fields []interface{}
	err := json.Unmarshal(schema, &fields)
	if err != nil {
		log.Fatalf("failed to parse schema: %v\n", err)
	}

	var parsedFields []Field
	for _, field := range fields {
		reflectedField := reflect.ValueOf(field).MapIndex(reflect.ValueOf("field"))
		reflectedSource := reflect.ValueOf(field).MapIndex(reflect.ValueOf("source"))

		j, err := json.Marshal(field)
		if err != nil {
			log.Fatalf("failed to stringify field '%s'\n", reflectedField)
		}

		switch reflectedSource.Interface() {
		case "randomAscii":
			parsedFields = append(parsedFields, ParseRandomAsciiField(j))
		case "randomUniformInt":
			parsedFields = append(parsedFields, ParseRandomUniformIntField(j))
		case "randomNormal":
			parsedFields = append(parsedFields, ParseRandomNormalField(j))
		}
	}

	return parsedFields
}

type Field interface {
	Generate(int)
}

type RandomAsciiSource struct {
	Source string `json:"source"`
	Config struct {
		Length int `json:"length"`
	} `json:"config"`
}

type RandomAsciiField struct {
	Field    string `json:"field"`
	DataType string `json:"data_type"`
	Data     []string
	RandomAsciiSource
}

// randomly generated ascii string with chars from 33-126
func (s *RandomAsciiField) Generate(n int) {
	var res []string
	for i := 0; i < n; i++ {
		var sb string
		for i := 0; i < s.Config.Length; i++ {
			sb += string(uint8(rand.IntN(93) + 33))
		}
		res = append(res, sb)
	}

	s.Data = res
}

func ParseRandomAsciiField(raw []byte) *RandomAsciiField {
	var parsedField *RandomAsciiField
	err := json.Unmarshal(raw, &parsedField)
	if err != nil {
		log.Fatalf(
			"failed to parse schema field in '%s'\n",
			string(raw),
		)
	}

	return parsedField
}

type RandomUniformIntSource struct {
	Source string `json:"source"`
	Config struct {
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"config"`
}

type RandomUniformIntField struct {
	Field    string `json:"field"`
	DataType string `json:"data_type"`
	Data     []int
	RandomUniformIntSource
}

func (s *RandomUniformIntField) Generate(n int) {
	var res []int
	var x int
	for i := 0; i < n; i++ {
		x = rand.IntN(s.Config.Max-s.Config.Min) + s.Config.Min
		res = append(res, x)
	}
	s.Data = res
}

func ParseRandomUniformIntField(raw []byte) *RandomUniformIntField {
	var parsedField *RandomUniformIntField
	err := json.Unmarshal(raw, &parsedField)
	if err != nil {
		log.Fatalf(
			"failed to parse schema field in '%s'\n",
			string(raw),
		)
	}

	return parsedField
}

type RandomNormalSource struct {
	Source string `json:"source"`
	Config struct {
		Mean     float64 `json:"mean"`
		Std      float64 `json:"std"`
		Decimals int     `json:"decimals"`
	} `json:"config"`
}

type RandomNormalField struct {
	Field    string `json:"field"`
	DataType string `json:"data_type"`
	Data     []float64
	RandomNormalSource
}

func (s *RandomNormalField) Generate(n int) {
	var res []float64
	var x float64
	dec := math.Pow10(s.Config.Decimals)
	for i := 0; i < n; i++ {
		x = rand.NormFloat64()*s.Config.Std + s.Config.Mean
		x = float64(int(x*dec)) / dec // round
		res = append(res, x)
	}

	s.Data = res
}

func ParseRandomNormalField(raw []byte) *RandomNormalField {
	var parsedField *RandomNormalField
	err := json.Unmarshal(raw, &parsedField)
	if err != nil {
		log.Fatalf(
			"failed to parse schema field in '%s'\n",
			string(raw),
		)
	}

	return parsedField
}
