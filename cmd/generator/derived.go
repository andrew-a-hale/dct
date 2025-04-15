package generator

import (
	"fmt"
	"log"
	"math"
	"reflect"
	"strconv"

	"github.com/expr-lang/expr"
)

type DerivedSource struct {
	Source   string `json:"source"`
	DataType string `json:"data_type"`
	Config   struct {
		Expression string `json:"expression"`
		Fields     []string
	} `json:"config"`
}

type DerivedField struct {
	Field string `json:"field"`
	DerivedSource
	Data string
}

var env = map[string]any{
	"power": func(a string, b int) string {
		x, _ := strconv.ParseFloat(a, 64)
		return fmt.Sprintf("%v", math.Pow(x, float64(b)))
	},
	"divide": func(a string, b string) string {
		x, _ := strconv.ParseFloat(a, 64)
		y, _ := strconv.ParseFloat(b, 64)
		return fmt.Sprintf("%v", x/y)
	},
	"plus": func(a string, b string) string {
		x, _ := strconv.ParseFloat(a, 64)
		y, _ := strconv.ParseFloat(b, 64)
		return fmt.Sprintf("%v", x+y)
	},
	"concat": func(a string, b string) string {
		return fmt.Sprintf("%s%s", a, b)
	},
	"minus": func(a string, b string) string {
		x, _ := strconv.ParseFloat(a, 64)
		y, _ := strconv.ParseFloat(b, 64)
		return fmt.Sprintf("%v", x-y)
	},
	"mult": func(a string, b string) string {
		x, _ := strconv.ParseFloat(a, 64)
		y, _ := strconv.ParseFloat(b, 64)
		return fmt.Sprintf("%v", x*y)
	},
	"mod": func(a string, b int) string {
		x, _ := strconv.Atoi(a)
		return fmt.Sprintf("%v", x%b)
	},
}

func (s *DerivedField) Generate(schema *[]Field, fieldMap *map[string]int) {
	fieldPtrs := make(map[string]reflect.Value)
	for _, f := range s.Config.Fields {
		idx := reflect.ValueOf(*fieldMap).MapIndex(reflect.ValueOf(f)).Int()
		fieldPtrs[f] = reflect.ValueOf(*schema).Index(int(idx))
	}

	for k, v := range fieldPtrs {
		field := v.Elem().Interface().(Field)
		env[k] = field.GetValue()
	}

	program, err := expr.Compile(
		s.Config.Expression,
		expr.Env(env),
		expr.Operator("+", "plus"),
		expr.Operator("-", "minus"),
		expr.Operator("*", "mult"),
		expr.Operator("%", "mod"),
		expr.Operator("/", "divide"),
		expr.Operator("^", "power"),
		expr.Operator("**", "power"),
		expr.Operator("||", "concat"),
	)
	if err != nil {
		log.Fatalf("failed to execute expression `%s` for field %s: %v", s.Config.Expression, s.Field, err)
	}
	o, _ := expr.Run(program, env)
	s.Data = fmt.Sprintf("%v", o)
}

func (s *DerivedField) GetValue() string {
	return s.Data
}

func (s *DerivedField) GetType() string {
	return s.DataType
}

func (s *DerivedField) GetName() string {
	return s.Field
}
