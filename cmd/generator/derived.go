package generator

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/expr-lang/expr"
)

type DerivedField struct {
	Field    string `json:"field"`
	Source   string `json:"source"`
	DataType string `json:"data_type"`
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

func (s DerivedField) Generate(ctx context.Context) any {
	fieldMap := ctx.Value(FIELD_MAP_KEY).(FieldMap)
	schema := ctx.Value(SCHEMA_KEY).(Schema)
	fieldPtrs := make(map[string]reflect.Value)
	for _, f := range s.Config.Fields {
		idx := reflect.ValueOf(fieldMap).MapIndex(reflect.ValueOf(f)).Int()
		fieldPtrs[f] = reflect.ValueOf(schema).Index(int(idx))
	}

	for k, v := range fieldPtrs {
		field := v.Elem().Interface().(Field)
		cacheValue := cache.GetValue(field.GetName())
		var value any
		var ok bool
		switch cacheValue.(type) {
		case bool:
			value, ok = cacheValue.(bool)
		case int, int32, int64:
			value, ok = cacheValue.(int)
		case float32:
			value, ok = cacheValue.(float32)
		case float64:
			value, ok = cacheValue.(float64)
		case string:
			value, ok = cacheValue.(string)
		default:
			log.Fatal("unimplemented type used in derived field")
		}

		if !ok {
			log.Fatalf("failed to parse value `%v`", cacheValue)
		}

		env[k] = value
	}

	program, err := expr.Compile(
		s.Config.Expression,
		expr.Env(env),
	)

	if err != nil {
		log.Fatalf(
			"failed to execute expression `%s` for field %s: %v",
			s.Config.Expression, s.Source, err,
		)
	}

	o, _ := expr.Run(program, env)
	return o
}

func (s DerivedField) GetName() string {
	return s.Field
}
