package generator

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/expr-lang/expr"
)

type DerivedField struct {
	Field  string `json:"field"`
	Source string `json:"source"`
	Config struct {
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
		switch v := cacheValue.(type) {
		case bool, int, int32, int64, float32, float64, string:
			env[k] = v
		default:
			log.Fatalf("unimplemented type used in derived field: %T", v)
		}
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
