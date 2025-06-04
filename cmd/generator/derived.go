package generator

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strconv"

	"github.com/expr-lang/expr"
)

type DerivedField struct {
	Field    string `json:"field"`
	Source   string `json:"source"`
	DataType string `json:"data_type"`
	Value    string
	Config   struct {
		Expression string   `json:"expression"`
		Fields     []string `json:"fields"`
	} `json:"config"`
}

var env = map[string]any{
	"concat": func(a string, b string) string {
		return fmt.Sprintf("%s%s", a, b)
	},
}

func (s DerivedField) Generate(ctx context.Context) string {
	fieldMap := ctx.Value(FIELD_MAP_KEY).(FieldMap)
	schema := ctx.Value(SCHEMA_KEY).(Schema)
	fieldPtrs := make(map[string]reflect.Value)
	for _, f := range s.Config.Fields {
		idx := reflect.ValueOf(fieldMap).MapIndex(reflect.ValueOf(f)).Int()
		fieldPtrs[f] = reflect.ValueOf(schema).Index(int(idx))
	}

	for k, v := range fieldPtrs {
		field := v.Elem().Interface().(Field)
		cacheValue := CACHE[field.GetName()]
		var value any
		var err error
		switch field.GetType() {
		case BOOL:
			value, err = strconv.ParseBool(cacheValue)
		case INT:
			value, err = strconv.Atoi(cacheValue)
		case FLOAT:
			value, err = strconv.ParseFloat(cacheValue, 32)
		case STRING:
			value = CACHE[field.GetName()]
		default:
			log.Fatal("unimplemented type used in derived field")
		}

		if err != nil {
			log.Fatalf("failed to parse value `%v`: %v", cacheValue, err)
		}

		env[k] = value
	}

	program, err := expr.Compile(
		s.Config.Expression,
		expr.Env(env),
		expr.Operator("||", "concat"),
	)

	if err != nil {
		log.Fatalf(
			"failed to execute expression `%s` for field %s: %v",
			s.Config.Expression, s.Source, err,
		)
	}

	o, _ := expr.Run(program, env)
	return fmt.Sprintf("%v", o)
}

func (s DerivedField) GetType() string {
	return s.DataType
}

func (s DerivedField) GetName() string {
	return s.Field
}
