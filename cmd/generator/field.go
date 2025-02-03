package generator

import (
	"context"
	"dct/cmd/utils"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"math"
	"math/rand/v2"
	"os"
	"reflect"
	"slices"
	"strconv"
	"time"

	"github.com/expr-lang/expr"
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
		case "derived":
			parsedFields = append(parsedFields, ParseField[DerivedField](j))
		}
	}

	return parsedFields
}

type Field interface {
	Generate(int, context.Context)
	GetSource() string
	GetValues() []string
	GetValue(int) string
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
func (s *RandomAsciiField) Generate(n int, ctx context.Context) {
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

func (s *RandomAsciiField) GetValues() []string {
	return s.Data
}

func (s *RandomAsciiField) GetValue(i int) string {
	return s.Data[i]
}

func (s *RandomAsciiField) GetSource() string {
	return s.Source
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

func (s *RandomUniformIntField) Generate(n int, ctx context.Context) {
	var res []int
	var x int
	for i := 0; i < n; i++ {
		x = rand.IntN(s.Config.Max-s.Config.Min) + s.Config.Min
		res = append(res, x)
	}
	s.Data = res
}

func (s *RandomUniformIntField) GetValues() []string {
	var res []string
	for _, v := range s.Data {
		res = append(res, fmt.Sprintf("%d", v))
	}
	return res
}

func (s *RandomUniformIntField) GetValue(i int) string {
	return strconv.Itoa(s.Data[i])
}

func (s *RandomUniformIntField) GetSource() string {
	return s.Source
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

func (s *RandomNormalField) Generate(n int, ctx context.Context) {
	var res []float64
	var x float64
	for i := 0; i < n; i++ {
		x = rand.NormFloat64()*s.Config.Std + s.Config.Mean
		res = append(res, x)
	}

	s.Data = res
}

func formatFloat(x float64, places int) string {
	return fmt.Sprintf(fmt.Sprintf("%%0.%df", places), x)
}

func (s *RandomNormalField) GetValues() []string {
	var res []string
	for _, v := range s.Data {
		res = append(res, formatFloat(v, s.Config.Decimals))
	}
	return res
}

func (s *RandomNormalField) GetValue(i int) string {
	return formatFloat(s.Data[i], s.Config.Decimals)
}

func (s *RandomNormalField) GetSource() string {
	return s.Source
}

type RandomPoissonSource struct {
	Source string `json:"source"`
	Config struct {
		Lambda int `json:"lambda"`
	} `json:"config"`
}

type RandomPoissonField struct {
	Field    string `json:"field"`
	DataType string `json:"data_type"`
	Data     []int
	RandomPoissonSource
}

func (s *RandomPoissonField) Generate(n int, ctx context.Context) {
	var res []int
	var x int
	for i := 0; i < n; i++ {
		x = generatePoisson(s.Config.Lambda)
		res = append(res, x)
	}

	s.Data = res
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

func (s *RandomPoissonField) GetValues() []string {
	var res []string
	for _, v := range s.Data {
		res = append(res, strconv.Itoa(v))
	}
	return res
}

func (s *RandomPoissonField) GetValue(i int) string {
	return strconv.Itoa(s.Data[i])
}

func (s *RandomPoissonField) GetSource() string {
	return s.Source
}

type LastNameSource struct {
	Source string `json:"source"`
}

type LastNameField struct {
	Field    string `json:"field"`
	DataType string `json:"data_type"`
	LastNameSource
	Data []string
}

func (s *LastNameField) Generate(n int, ctx context.Context) {
	query := fmt.Sprintf(`
select name
from last_names
cross join generate_series(1, %d)
using sample reservoir(%d rows)`,
		max(n/100, 2),
		n,
	)
	result, err := utils.Query(query)
	if err != nil {
		log.Fatalf("failed to sample first_names from duckdb: %v\n", err)
	}

	s.Data = slices.Concat(result.Rows...)
}

func (s *LastNameField) GetValues() []string {
	return s.Data
}

func (s *LastNameField) GetValue(i int) string {
	return s.Data[i]
}

func (s *LastNameField) GetSource() string {
	return s.Source
}

type FirstNameSource struct {
	Source string `json:"source"`
}

type FirstNameField struct {
	Field    string `json:"field"`
	DataType string `json:"data_type"`
	FirstNameSource
	Data []string
}

func (s *FirstNameField) Generate(n int, ctx context.Context) {
	query := fmt.Sprintf(`
select name
from first_names
cross join generate_series(1, %d)
using sample reservoir(%d rows)`,
		max(n/100, 2),
		n,
	)
	result, err := utils.Query(query)
	if err != nil {
		log.Fatalf("failed to sample first_names from duckdb: %v\n", err)
	}

	s.Data = slices.Concat(result.Rows...)
}

func (s *FirstNameField) GetValues() []string {
	return s.Data
}

func (s *FirstNameField) GetValue(i int) string {
	return s.Data[i]
}

func (s *FirstNameField) GetSource() string {
	return s.Source
}

type RandomDatetimeSource struct {
	Source string `json:"source"`
	Config struct {
		Tz  string `json:"tz"`
		Min string `json:"min"`
		Max string `json:"max"`
	} `json:"config"`
}

type RandomDatetimeField struct {
	RandomDatetimeSource
	Field    string `json:"field"`
	DataType string `json:"data_type"`
	Data     []time.Time
}

func (s *RandomDatetimeField) Generate(n int, ctx context.Context) {
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

	var res []time.Time
	for i := 0; i < n; i++ {
		dt := time.Unix(rand.Int64N(ub-lb)+lb, 0).In(loc)
		res = append(res, dt)
	}
	s.Data = res
}

func (s *RandomDatetimeField) GetValues() []string {
	var res []string
	for _, dt := range s.Data {
		res = append(res, dt.Format(time.RFC3339))
	}
	return res
}

func (s *RandomDatetimeField) GetValue(i int) string {
	return s.Data[i].Format(time.RFC3339)
}

func (s *RandomDatetimeField) GetSource() string {
	return s.Source
}

type RandomDateSource struct {
	Source string `json:"source"`
	Config struct {
		Min string `json:"min"`
		Max string `json:"max"`
	} `json:"config"`
}

type RandomDateField struct {
	RandomDateSource
	Field    string `json:"field"`
	DataType string `json:"data_type"`
	Data     []time.Time
}

func (s *RandomDateField) Generate(n int, ctx context.Context) {
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

	var res []time.Time
	for i := 0; i < n; i++ {
		dt := time.Unix(rand.Int64N(ub-lb)+lb, 0)
		res = append(res, dt)
	}
	s.Data = res
}

func (s *RandomDateField) GetValues() []string {
	var res []string
	for _, dt := range s.Data {
		res = append(res, dt.Format(time.DateOnly))
	}
	return res
}

func (s *RandomDateField) GetValue(i int) string {
	return s.Data[i].Format(time.DateOnly)
}

func (s *RandomDateField) GetSource() string {
	return s.Source
}

type RandomTimeSource struct {
	Source string `json:"source"`
	Config struct {
		Min string `json:"min"`
		Max string `json:"max"`
	} `json:"config"`
}

type RandomTimeField struct {
	RandomDateSource
	Field    string `json:"field"`
	DataType string `json:"data_type"`
	Data     []time.Time
}

func (s *RandomTimeField) Generate(n int, ctx context.Context) {
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

	var res []time.Time
	for i := 0; i < n; i++ {
		dt := time.Unix(rand.Int64N(ub-lb)+lb, 0).In(time.UTC)
		res = append(res, dt)
	}
	s.Data = res
}

func (s *RandomTimeField) GetValues() []string {
	var res []string
	for _, dt := range s.Data {
		res = append(res, dt.Format(time.TimeOnly))
	}
	return res
}

func (s *RandomTimeField) GetValue(i int) string {
	return s.Data[i].Format(time.TimeOnly)
}

func (s *RandomTimeField) GetSource() string {
	return s.Source
}

type DerivedSource struct {
	Source string `json:"source"`
	Config struct {
		Expression string `json:"expression"`
		Fields     []string
	} `json:"config"`
}

type DerivedField struct {
	Field    string `json:"field"`
	DataType string `json:"data_type"`
	DerivedSource
	Data []string
}

func (s *DerivedField) Generate(n int, ctx context.Context) {
	fieldPtrs := make(map[string]reflect.Value)
	schema := ctx.Value("schema")
	fieldMap := ctx.Value("fieldMap")
	for _, f := range s.Config.Fields {
		idx := reflect.ValueOf(fieldMap).MapIndex(reflect.ValueOf(f)).Int()
		fieldPtrs[f] = reflect.ValueOf(schema).Index(int(idx))
	}

	for i := 0; i < n; i++ {
		env := make(map[string]interface{})
		for k, v := range fieldPtrs {
			field := v.Elem().Interface().(Field)
			fieldType := reflect.ValueOf(field).Elem().FieldByName("DataType")
			log.Fatal("todo", fieldType)
			env[k] = field.GetValue(i)
		}
		program, err := expr.Compile(s.Config.Expression, expr.Env(env))
		if err != nil {
			log.Fatalf("failed to execute expression `%s` for field %s: %v", s.Config.Expression, s.Field, err)
		}
		o, err := expr.Run(program, env)
		s.Data = append(s.Data, o.(string))
	}
}

func (s *DerivedField) GetValues() []string {
	return s.Data
}

func (s *DerivedField) GetValue(i int) string {
	return s.Data[i]
}

func (s *DerivedField) GetSource() string {
	return s.Source
}
