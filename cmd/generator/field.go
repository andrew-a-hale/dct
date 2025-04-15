// add random enum from set
// add random email
package generator

import (
	"dct/cmd/generator/sources"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"math"
	"math/rand/v2"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type DataType = string

const (
	STRING = "string"
	INT    = "int"
	FLOAT  = "float"
	BOOL   = "bool"
)

func ParseField[T any](raw []byte) *T {
	var parsedField *T
	err := json.Unmarshal(raw, &parsedField)
	if err != nil {
		log.Fatalf(
			"failed to parse schema field in '%s'\n",
			string(raw),
		)
	}

	return parsedField
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

	var fields []any
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
			parsedFields = append(parsedFields, ParseField[RandomAsciiField](j))
		case "randomUniformInt":
			parsedFields = append(parsedFields, ParseField[RandomUniformIntField](j))
		case "randomNormal":
			parsedFields = append(parsedFields, ParseField[RandomNormalField](j))
		case "randomPoisson":
			parsedFields = append(parsedFields, ParseField[RandomPoissonField](j))
		case "firstNames":
			parsedFields = append(parsedFields, ParseField[FirstNameField](j))
		case "lastNames":
			parsedFields = append(parsedFields, ParseField[LastNameField](j))
		case "randomDatetime":
			parsedFields = append(parsedFields, ParseField[RandomDatetimeField](j))
		case "randomTime":
			parsedFields = append(parsedFields, ParseField[RandomTimeField](j))
		case "randomDate":
			parsedFields = append(parsedFields, ParseField[RandomDateField](j))
		case "uuid":
			parsedFields = append(parsedFields, ParseField[UuidField](j))
		case "derived":
			parsedFields = append(parsedFields, ParseField[DerivedField](j))
		}
	}

	return parsedFields
}

type Field interface {
	Generate(*[]Field, *map[string]int)
	GetValue() string
	GetType() string
	GetName() string
}

type RandomAsciiSource struct {
	Source   string `json:"source"`
	DataType string `json:"data_type"`
	Config   struct {
		Length int `json:"length"`
	} `json:"config"`
}

type RandomAsciiField struct {
	Field string `json:"field"`
	Data  string
	RandomAsciiSource
}

// randomly generated ascii string with chars from 33-126
func (s *RandomAsciiField) Generate(schema *[]Field, fieldMap *map[string]int) {
	var sb string
	for range s.Config.Length {
		sb += string(uint8(rand.IntN(93) + 33))
	}

	s.Data = sb
}

func (s *RandomAsciiField) GetValue() string {
	return s.Data
}

func (s *RandomAsciiField) GetType() string {
	return s.DataType
}

func (s *RandomAsciiField) GetName() string {
	return s.Field
}

type RandomUniformIntSource struct {
	Source   string `json:"source"`
	DataType string `json:"data_type"`
	Config   struct {
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"config"`
}

type RandomUniformIntField struct {
	Field string `json:"field"`
	Data  int
	RandomUniformIntSource
}

func (s *RandomUniformIntField) Generate(schema *[]Field, fieldMap *map[string]int) {
	s.Data = rand.IntN(s.Config.Max-s.Config.Min) + s.Config.Min
}

func (s *RandomUniformIntField) GetValue() string {
	return strconv.Itoa(s.Data)
}

func (s *RandomUniformIntField) GetType() string {
	return s.DataType
}

func (s *RandomUniformIntField) GetName() string {
	return s.Field
}

type RandomNormalSource struct {
	Source   string `json:"source"`
	DataType string `json:"data_type"`
	Config   struct {
		Mean     float64 `json:"mean"`
		Std      float64 `json:"std"`
		Decimals int     `json:"decimals"`
	} `json:"config"`
}

type RandomNormalField struct {
	Field string `json:"field"`
	Data  float64
	RandomNormalSource
}

func (s *RandomNormalField) Generate(schema *[]Field, fieldMap *map[string]int) {
	s.Data = rand.NormFloat64()*s.Config.Std + s.Config.Mean
}

func formatFloat(x float64, places int) string {
	return fmt.Sprintf(fmt.Sprintf("%%0.%df", places), x)
}

func (s *RandomNormalField) GetValue() string {
	return formatFloat(s.Data, s.Config.Decimals)
}

func (s *RandomNormalField) GetType() string {
	return s.DataType
}

func (s *RandomNormalField) GetName() string {
	return s.Field
}

type RandomPoissonSource struct {
	Source   string `json:"source"`
	DataType string `json:"data_type"`
	Config   struct {
		Lambda int `json:"lambda"`
	} `json:"config"`
}

type RandomPoissonField struct {
	Field string `json:"field"`
	Data  int
	RandomPoissonSource
}

func (s *RandomPoissonField) Generate(schema *[]Field, fieldMap *map[string]int) {
	s.Data = generatePoisson(s.Config.Lambda)
}

func generatePoisson(lambda int) int {
	var n int

	for s := 0.0; s < 1; {
		u := rand.Float64()
		e := -math.Log(u) / float64(lambda)
		n += 1
		s += e
	}

	return n
}

func (s *RandomPoissonField) GetValue() string {
	return strconv.Itoa(s.Data)
}

func (s *RandomPoissonField) GetType() string {
	return s.DataType
}

func (s *RandomPoissonField) GetName() string {
	return s.Field
}

type LastNameSource struct {
	Source     string `json:"source"`
	DataType   string `json:"data_type"`
	SourceData []string
}

type LastNameField struct {
	Field string `json:"field"`
	LastNameSource
	Data string
}

func (s *LastNameField) Init() {
	s.SourceData = sources.LastNames
}

func (s *LastNameField) Generate(schema *[]Field, fieldMap *map[string]int) {
	if len(s.SourceData) == 0 {
		s.Init()
	}
	s.Data = s.SourceData[rand.IntN(len(s.SourceData))]
}

func (s *LastNameField) GetValue() string {
	return s.Data
}

func (s *LastNameField) GetType() string {
	return s.DataType
}

func (s *LastNameField) GetName() string {
	return s.Field
}

type FirstNameSource struct {
	Source     string `json:"source"`
	SourceData []string
	DataType   string `json:"data_type"`
}

type FirstNameField struct {
	Field string `json:"field"`
	FirstNameSource
	Data string
}

func (s *FirstNameField) Init() {
	for _, t := range sources.FirstNames {
		s.SourceData = append(s.SourceData, t.Name)
	}
}

func (s *FirstNameField) Generate(schema *[]Field, fieldMap *map[string]int) {
	if len(s.SourceData) == 0 {
		s.Init()
	}
	s.Data = s.SourceData[rand.IntN(len(s.SourceData))]
}

func (s *FirstNameField) GetValue() string {
	return s.Data
}

func (s *FirstNameField) GetType() string {
	return s.DataType
}

func (s *FirstNameField) GetName() string {
	return s.Field
}

type RandomDatetimeSource struct {
	Source   string `json:"source"`
	DataType string `json:"data_type"`
	Config   struct {
		Tz  string `json:"tz"`
		Min string `json:"min"`
		Max string `json:"max"`
	} `json:"config"`
}

type RandomDatetimeField struct {
	RandomDatetimeSource
	Field string `json:"field"`
	Data  time.Time
}

func (s *RandomDatetimeField) Generate(schema *[]Field, fieldMap *map[string]int) {
	MAX_TIME := time.Unix(1<<63-62135596801, 999999999)
	MIN_TIME := time.Unix(0, 0)

	loc, err := time.LoadLocation(s.Config.Tz)
	if err != nil {
		log.Fatalf("failed to parse tz: %v\n", err)
	}

	// handle min datetime
	var parsedDtMin time.Time
	if s.Config.Min != "" {
		parsedDtMin, err = time.ParseInLocation(time.DateTime, s.Config.Min, loc)
		if err != nil {
			log.Fatalf("failed to parse max datetime: %v\n", err)
		}
	} else {
		parsedDtMin = MIN_TIME
	}

	lb := MIN_TIME.Unix()
	if parsedDtMin.After(MIN_TIME) {
		lb = parsedDtMin.Unix()
	}

	// handle max datetime
	var parsedDtMax time.Time
	if s.Config.Max != "" {
		parsedDtMax, err = time.ParseInLocation(time.DateTime, s.Config.Max, loc)
		if err != nil {
			log.Fatalf("failed to parse max datetime: %v\n", err)
		}
	} else {
		parsedDtMax = MAX_TIME
	}

	ub := MAX_TIME.Unix()
	if parsedDtMax.Before(MAX_TIME) {
		ub = parsedDtMax.Unix()
	}

	s.Data = time.Unix(rand.Int64N(ub-lb)+lb, 0).In(loc)
}

func (s *RandomDatetimeField) GetValue() string {
	return s.Data.Format(time.RFC3339)
}

func (s *RandomDatetimeField) GetType() string {
	return s.DataType
}

func (s *RandomDatetimeField) GetName() string {
	return s.Field
}

type RandomDateSource struct {
	Source   string `json:"source"`
	DataType string `json:"data_type"`
	Config   struct {
		Min string `json:"min"`
		Max string `json:"max"`
	} `json:"config"`
}

type RandomDateField struct {
	RandomDateSource
	Field string `json:"field"`
	Data  time.Time
}

func (s *RandomDateField) Generate(schema *[]Field, fieldMap *map[string]int) {
	MAX_TIME := time.Unix(1<<63-62135596801, 999999999)
	MIN_TIME := time.Unix(0, 0)

	// handle min date
	var err error
	var parsedDtMin time.Time
	if s.Config.Min != "" {
		parsedDtMin, err = time.Parse(time.DateOnly, s.Config.Min)
		if err != nil {
			log.Fatalf("failed to parse max date: %v\n", err)
		}
	} else {
		parsedDtMin = MIN_TIME
	}

	lb := MIN_TIME.Unix()
	if parsedDtMin.After(MIN_TIME) {
		lb = parsedDtMin.Unix()
	}

	// handle max date
	var parsedDtMax time.Time
	if s.Config.Max != "" {
		parsedDtMax, err = time.Parse(time.DateOnly, s.Config.Max)
		if err != nil {
			log.Fatalf("failed to parse max date: %v\n", err)
		}
	} else {
		parsedDtMax = MAX_TIME
	}

	ub := MAX_TIME.Unix()
	if parsedDtMax.Before(MAX_TIME) {
		ub = parsedDtMax.Unix()
	}

	s.Data = time.Unix(rand.Int64N(ub-lb)+lb, 0)
}

func (s *RandomDateField) GetValue() string {
	return s.Data.Format(time.DateOnly)
}

func (s *RandomDateField) GetType() string {
	return s.DataType
}

func (s *RandomDateField) GetName() string {
	return s.Field
}

type RandomTimeSource struct {
	Source   string `json:"source"`
	DataType string `json:"data_type"`
	Config   struct {
		Min string `json:"min"`
		Max string `json:"max"`
	} `json:"config"`
}

type RandomTimeField struct {
	RandomDateSource
	Field string `json:"field"`
	Data  time.Time
}

func (s *RandomTimeField) Generate(schema *[]Field, fieldMap *map[string]int) {
	MAX_TIME, _ := time.ParseInLocation(time.TimeOnly, "23:59:59", time.UTC)
	MIN_TIME, _ := time.ParseInLocation(time.TimeOnly, "00:00:00", time.UTC)

	// handle min time
	var err error
	var parsedDtMin time.Time
	if s.Config.Min != "" {
		if len(s.Config.Min) < 8 {
			log.Fatalf("invalid format for min must be HH:MM:SS, not %v\n", s.Config.Min)
		}
		parsedDtMin, err = time.ParseInLocation(time.TimeOnly, s.Config.Min, time.UTC)
		if err != nil {
			log.Fatalf("failed to parse max time: %v\n", err)
		}
	} else {
		parsedDtMin = MIN_TIME
	}

	lb := MIN_TIME.Unix()
	if parsedDtMin.After(MIN_TIME) {
		lb = parsedDtMin.Unix()
	}

	// handle max time
	var parsedDtMax time.Time
	if s.Config.Max != "" {
		if len(s.Config.Max) < 8 {
			log.Fatalf("invalid format for max must be HH:MM:SS, not %v\n", s.Config.Max)
		}
		parsedDtMax, err = time.ParseInLocation(time.TimeOnly, s.Config.Max, time.UTC)
		if err != nil {
			log.Fatalf("failed to parse max time: %v\n", err)
		}
	} else {
		parsedDtMax = MAX_TIME
	}

	ub := MAX_TIME.Unix()
	if parsedDtMax.Before(MAX_TIME) {
		ub = parsedDtMax.Unix()
	}

	s.Data = time.Unix(rand.Int64N(ub-lb)+lb, 0).In(time.UTC)
}

func (s *RandomTimeField) GetValue() string {
	return s.Data.Format(time.TimeOnly)
}

func (s *RandomTimeField) GetType() string {
	return s.DataType
}

func (s *RandomTimeField) GetName() string {
	return s.Field
}

type UuidSource struct {
	Source   string `json:"source"`
	DataType string `json:"data_type"`
}

type UuidField struct {
	Field string `json:"field"`
	Data  string
	RandomAsciiSource
}

func (s *UuidField) Generate(schema *[]Field, fieldMap *map[string]int) {
	s.Data = uuid.NewString()
}

func (s *UuidField) GetValue() string {
	return s.Data
}

func (s *UuidField) GetType() string {
	return s.DataType
}

func (s *UuidField) GetName() string {
	return s.Field
}
