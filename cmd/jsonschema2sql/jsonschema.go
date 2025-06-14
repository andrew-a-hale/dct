package jsonschema2sql

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	// JSON
	JSON_ARRAY   string = "array"
	JSON_BOOLEAN string = "boolean"
	JSON_NULL    string = "null"
	JSON_NUMBER  string = "number"
	JSON_INTEGER string = "integer"
	JSON_OBJECT  string = "object"
	JSON_STRING  string = "string"
	JSON_NOTYPE  string = ""

	// SQL
	SQL_BOOLEAN string = "bool"
	SQL_NUMBER  string = "float"
	SQL_INTEGER string = "int"
	SQL_STRING  string = "varchar"
)

type Property struct {
	Type       any             `json:"type"`
	Reference  string          `json:"$ref"`
	Properties json.RawMessage `json:"properties"`
	Items      []string        `json:"items"`
}

type JsonSchema struct {
	Id         string              `json:"id"`
	Type       any                 `json:"type"`
	Reference  string              `json:"$ref"`
	Properties map[string]Property `json:"properties"`
}

type Field struct {
	Name string
	Type string
}

type (
	SqlSchema map[string]Table
	Table     map[string]Field
)

func (s SqlSchema) ToSql() (string, error) {
	return "", nil
}

type KeyedProperty struct {
	Key      string
	Property Property
}

type Stack []KeyedProperty

func (s *Stack) Len() int {
	return len(*s)
}

func (s *Stack) Pop() KeyedProperty {
	i := s.Len() - 1
	x := (*s)[i]
	*s = (*s)[i : i+1]
	return x
}

func (s *Stack) Push(kp KeyedProperty) {
	*s = append(*s, kp)
}

func process(data []byte) (SqlSchema, error) {
	var j JsonSchema
	err := json.Unmarshal(data, &j)
	if err != nil {
		return SqlSchema{}, err
	}

	var stack Stack
	for k, v := range j.Properties {
		stack.Push(KeyedProperty{k, v})
	}

	for len(stack) > 0 {
		kv := stack.Pop()
		k, v := kv.Key, kv.Property

		switch v.Type.(type) {
		case []string:
			field := Field{k, JSON_STRING}
			fmt.Println(field)
		case string:
			switch v.Type.(string) {
			case JSON_OBJECT:
			case JSON_ARRAY:
				array, _ := handleArray(v)
				fmt.Println(array)
			case JSON_NOTYPE:
				field, _ := followReference(j.Reference)
				fmt.Println(field)
			case JSON_NUMBER:
				fallthrough
			case JSON_STRING:
				fallthrough
			case JSON_INTEGER:
				fallthrough
			case JSON_BOOLEAN:
				field := Field{k, v.Type.(string)}
				fmt.Println(field)
			default:
				return SqlSchema{}, fmt.Errorf("failed to process jsonschema: type %v not implemented", v.Type)
			}
		default:
			return SqlSchema{}, fmt.Errorf("failed to process jsonschema: type %T not implemented", v.Type)
		}
	}

	return SqlSchema{}, errors.New("failed to find any properties")
}

func followReference(ref string) (Field, error) {
	fmt.Println(ref)
	return Field{}, nil
}

// can return any of array of types, or array of objects
func handleArray(prop Property) (any, error) {
	fmt.Println(prop)
	return "", nil
}
