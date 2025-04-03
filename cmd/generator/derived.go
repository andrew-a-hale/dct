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

func (s *DerivedField) Generate(n int, schema *[]Field, fieldMap *map[string]int) {
	fieldPtrs := make(map[string]reflect.Value)
	for _, f := range s.Config.Fields {
		idx := reflect.ValueOf(*fieldMap).MapIndex(reflect.ValueOf(f)).Int()
		fieldPtrs[f] = reflect.ValueOf(*schema).Index(int(idx))
	}

	for i := range n {
		for k, v := range fieldPtrs {
			field := v.Elem().Interface().(Field)
			env[k] = field.GetValue(i)
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
		)
		if err != nil {
			log.Fatalf("failed to execute expression `%s` for field %s: %v", s.Config.Expression, s.Field, err)
		}
		o, _ := expr.Run(program, env)
		s.Data = append(s.Data, fmt.Sprintf("%v", o))
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
