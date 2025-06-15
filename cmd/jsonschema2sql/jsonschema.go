package jsonschema2sql

import (
	"dct/cmd/utils"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"reflect"
	"strings"
)

const (
	JSON_ARRAY   = "array"
	JSON_BOOLEAN = "boolean"
	JSON_NULL    = "null"
	JSON_NUMBER  = "number"
	JSON_INTEGER = "integer"
	JSON_OBJECT  = "object"
	JSON_STRING  = "string"
	SQL_BOOLEAN  = "bool"
	SQL_NUMBER   = "float"
	SQL_INTEGER  = "int"
	SQL_STRING   = "varchar"
)

var (
	jsonToSqlMap = make(map[string]string)
	logLevel     = slog.LevelDebug
)

// setup logger and jsonToSqlMap
func init() {
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	jsonToSqlMap[JSON_BOOLEAN] = SQL_BOOLEAN
	jsonToSqlMap[JSON_NUMBER] = SQL_NUMBER
	jsonToSqlMap[JSON_INTEGER] = SQL_INTEGER
	jsonToSqlMap[JSON_STRING] = SQL_STRING
}

type Property struct {
	Type       any             `json:"type"`
	Reference  string          `json:"$ref"`
	Properties json.RawMessage `json:"properties"`
}

type JsonSchema struct {
	Type        any                 `json:"type"`
	Reference   string              `json:"$ref"`
	Properties  map[string]Property `json:"properties"`
	Definitions map[string]Property `json:"definitions"`
}

type Field struct {
	Path  []string
	Type  string
	Array bool
}

func fieldsToDef(fields []Field) (string, error) {
	var sqlDef string

	tree := BuildTree(fields)

	var stack utils.Stack[*TreeNode]
	stack.Push(tree)
	slog.Debug("build sql", "tree", tree, "stack", stack)

	sqlDef += "create table %s (\n\tid varchar primary key"
	for stack.Len() > 0 {
		node := stack.Pop()
		if node.IsLeaf {
			sqlDef += fmt.Sprintf("\n\t, %s %s", node.Id, node.Type)
		} else if node.IsArray {
			sqlDef += fmt.Sprintf("\n\t, %s array(%s)", node.Id, node.Type)
		}

		for _, c := range node.Children {
			stack.Push(c)
		}
	}
	sqlDef += "\n)"
	sqlDef = fmt.Sprintf(sqlDef, "test")

	return sqlDef, nil
}

// Example
//
// create table schema.table (
//
//	id varchar primary key,
//	name varchar,
//	houses array(address row(street varchar, city varchar)
//
// )

type PropertyPath struct {
	Path     []string
	Property Property
	IsArray  bool
}

func process(data []byte) (string, error) {
	var j JsonSchema
	err := json.Unmarshal(data, &j)
	if err != nil {
		return "", err
	}

	var stack utils.Stack[PropertyPath]
	for k, v := range j.Properties {
		stack.Push(PropertyPath{[]string{k}, v, false})
	}

	var fields []Field
	for len(stack) > 0 {
		kv := stack.Pop()
		path, property, isArray := kv.Path, kv.Property, kv.IsArray

		if s, ok := property.Type.(string); ok && strings.HasPrefix(s, "[") {
			fields = append(fields, Field{path, JSON_STRING, isArray})
		}

		switch property.Type {
		case nil:
			prop, err := followReference(path, property.Reference, &j, isArray)
			if err != nil {
				log.Fatal(err)
			}
			stack.Push(prop)
		case JSON_ARRAY:
			props, _ := handleArray(path, property)
			kv.IsArray = true
			stack.Push(props...)
		case JSON_OBJECT:
			props := make(map[string]Property)
			err := json.Unmarshal(property.Properties, &props)
			if err != nil {
				return "", err
			}
			for k, v := range props {
				stack.Push(PropertyPath{append(path, k), v, isArray})
			}
		case JSON_NUMBER, JSON_STRING, JSON_INTEGER, JSON_BOOLEAN:
			fields = append(fields, Field{path, jsonToSqlMap[property.Type.(string)], isArray})
		default:
			return "", fmt.Errorf("failed to process jsonschema: type %v not implemented", property.Type)
		}
	}

	return fieldsToDef(fields)
}

func followReference(path []string, ref string, data *JsonSchema, fromArray bool) (PropertyPath, error) {
	var props PropertyPath

	refPath := strings.Split(ref, "/")
	utils.Assert(refPath[0] == "#", "only support internal references")
	refPath = refPath[1:]

	// follow path to reference
	schema := reflect.ValueOf(data)
	for len(refPath) > 0 {
		stub := refPath[0]
		switch schema.Kind() {
		case reflect.Ptr:
			schema = schema.Elem()
		case reflect.Struct:
			name := strings.ToUpper(stub[:1]) + stub[1:]
			schema = schema.FieldByName(name)
			refPath = refPath[1:]
		case reflect.Map:
			v := reflect.ValueOf(stub)
			schema = schema.MapIndex(v)
			refPath = refPath[1:]
		default:
			slog.Error("reflecting on ?", "refPath", refPath, "schema", schema.Kind())
			return props, fmt.Errorf("invalid reference path: %v", ref)
		}
	}

	// get prop from reference
	prop, ok := schema.Interface().(Property)
	if !ok {
		return props, errors.New("failed to follow reference")
	}

	return PropertyPath{append(path), prop, fromArray}, nil
}

func handleArray(path []string, property Property) ([]PropertyPath, error) {
	var props []PropertyPath
	obj := make(map[string]Property)
	err := json.Unmarshal(property.Properties, &obj)
	if err != nil {
		return props, err
	}

	for k, v := range obj {
		props = append(props, PropertyPath{append(path, k), v, true})
	}

	return props, nil
}
