package generator

import (
	"context"
	"dct/cmd/generator/sources"
	"encoding/json"
	"io"
	"log"
	"math"
	"math/rand/v2"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func ParseField[T Field](raw []byte) *T {
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

func parseSchema(rawSchema string) Schema {
	var schema []byte
	if json.Valid([]byte(rawSchema)) {
		schema = []byte(rawSchema)
	} else {
		f, err := os.Open(rawSchema)
		if err != nil {
			log.Fatalf("failed to open schema file: %v\n", err)
		}
		defer f.Close()

		schema, err = io.ReadAll(f)
		if err != nil {
			log.Fatalf("failed to read schema file: %v\n", err)
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
		case "randomBool":
			parsedFields = append(parsedFields, ParseField[RandomBoolField](j))
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
		case "emails":
			parsedFields = append(parsedFields, ParseField[EmailField](j))
		case "companies":
			parsedFields = append(parsedFields, ParseField[CompanyField](j))
		case "derived":
			parsedFields = append(parsedFields, ParseField[DerivedField](j))
		}
	}

	return parsedFields
}

type Field interface {
	Generate(context.Context) any
	GetName() string
}

type RandomBoolField struct {
	Field    string `json:"field"`
	Source   string `json:"source"`
	DataType string `json:"data_type"`
}

// randomly generated ascii string with chars from 33-126
func (s RandomBoolField) Generate(ctx context.Context) any {
	var value bool
	if rand.Float32() > 0.5 {
		value = true
	} else {
		value = false
	}

	cache.PutValue(s.Field, value)
	return value
}

func (s RandomBoolField) GetName() string {
	return s.Field
}

type RandomAsciiField struct {
	Field    string `json:"field"`
	Source   string `json:"source"`
	DataType string `json:"data_type"`
	Config   struct {
		Length int `json:"length"`
	} `json:"config"`
}

// randomly generated ascii string with chars from 33-126
func (s RandomAsciiField) Generate(ctx context.Context) any {
	var value string
	for range s.Config.Length {
		value += string(uint8(rand.IntN(93) + 33))
	}

	cache.PutValue(s.Field, value)

	return value
}

func (s RandomAsciiField) GetName() string {
	return s.Field
}

type RandomUniformIntField struct {
	Field    string `json:"field"`
	Source   string `json:"source"`
	DataType string `json:"data_type"`
	Config   struct {
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"config"`
}

func (s RandomUniformIntField) Generate(ctx context.Context) any {
	value := rand.IntN(s.Config.Max-s.Config.Min) + s.Config.Min
	cache.PutValue(s.Field, value)
	return value
}

func (s RandomUniformIntField) GetName() string {
	return s.Field
}

type RandomNormalField struct {
	Field    string `json:"field"`
	Source   string `json:"source"`
	DataType string `json:"data_type"`
	Config   struct {
		Mean float64 `json:"mean"`
		Std  float64 `json:"std"`
	} `json:"config"`
}

func (s RandomNormalField) Generate(ctx context.Context) any {
	value := rand.NormFloat64()*s.Config.Std + s.Config.Mean
	cache.PutValue(s.Field, value)
	return value
}

func (s RandomNormalField) GetName() string {
	return s.Field
}

type RandomPoissonField struct {
	Field    string `json:"field"`
	Source   string `json:"source"`
	DataType string `json:"data_type"`
	Config   struct {
		Lambda int `json:"lambda"`
	} `json:"config"`
}

func (s RandomPoissonField) Generate(ctx context.Context) any {
	value := strconv.Itoa(generatePoisson(s.Config.Lambda))
	cache.PutValue(s.Field, value)
	return value
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

func (s RandomPoissonField) GetName() string {
	return s.Field
}

type LastNameField struct {
	Field    string `json:"field"`
	Source   string `json:"source"`
	DataType string `json:"data_type"`
}

func (s LastNameField) Generate(ctx context.Context) any {
	value := sources.LastNames[rand.IntN(len(sources.LastNames))]
	cache.PutValue(s.Field, value)
	return value
}

func (s LastNameField) GetName() string {
	return s.Field
}

type FirstNameField struct {
	Field    string `json:"field"`
	Source   string `json:"source"`
	DataType string `json:"data_type"`
}

func (s FirstNameField) Generate(ctx context.Context) any {
	value := sources.FirstNames[rand.IntN(len(sources.FirstNames))]
	cache.PutValue(s.Field, value)
	return value
}

func (s FirstNameField) GetName() string {
	return s.Field
}

type RandomDatetimeField struct {
	Field    string `json:"field"`
	Source   string `json:"source"`
	DataType string `json:"data_type"`
	Config   struct {
		Tz  string `json:"tz"`
		Min string `json:"min"`
		Max string `json:"max"`
	} `json:"config"`
}

func (s RandomDatetimeField) Generate(ctx context.Context) any {
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

	value := time.Unix(rand.Int64N(ub-lb)+lb, 0).In(loc).Format(time.RFC3339)
	cache.PutValue(s.Field, value)
	return value
}

func (s RandomDatetimeField) GetName() string {
	return s.Field
}

type RandomDateField struct {
	Field    string `json:"field"`
	Source   string `json:"source"`
	DataType string `json:"data_type"`
	Config   struct {
		Min string `json:"min"`
		Max string `json:"max"`
	} `json:"config"`
}

func (s RandomDateField) Generate(ctx context.Context) any {
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

	value := time.Unix(rand.Int64N(ub-lb)+lb, 0).Format(time.DateOnly)
	cache.PutValue(s.Field, value)
	return value
}

func (s RandomDateField) GetName() string {
	return s.Field
}

type RandomTimeField struct {
	Field    string `json:"field"`
	Source   string `json:"source"`
	DataType string `json:"data_type"`
	Config   struct {
		Min string `json:"min"`
		Max string `json:"max"`
	} `json:"config"`
}

func (s RandomTimeField) Generate(ctx context.Context) any {
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

	value := time.Unix(rand.Int64N(ub-lb)+lb, 0).In(time.UTC).Format(time.TimeOnly)
	cache.PutValue(s.Field, value)
	return value
}

func (s RandomTimeField) GetName() string {
	return s.Field
}

type UuidField struct {
	Field    string `json:"field"`
	Source   string `json:"source"`
	DataType string `json:"data_type"`
}

func (s UuidField) Generate(ctx context.Context) any {
	value := uuid.NewString()
	cache.PutValue(s.Field, value)
	return value
}

func (s UuidField) GetName() string {
	return s.Field
}

type EmailField struct {
	Field    string `json:"field"`
	Source   string `json:"source"`
	DataType string `json:"data_type"`
}

func (s EmailField) Generate(ctx context.Context) any {
	value := sources.Emails[rand.IntN(len(sources.Emails))]
	cache.PutValue(s.Field, value)
	return value
}

func (s EmailField) GetName() string {
	return s.Field
}

type CompanyField struct {
	Field    string `json:"field"`
	Source   string `json:"source"`
	DataType string `json:"data_type"`
	Value    string
}

func (s CompanyField) Generate(ctx context.Context) any {
	value := sources.Companies[rand.IntN(len(sources.Companies))]
	cache.PutValue(s.Field, value)
	return value
}

func (s CompanyField) GetName() string {
	return s.Field
}
