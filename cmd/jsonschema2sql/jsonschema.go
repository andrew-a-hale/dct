package jsonschema2sql

import "encoding/json"

type JsonValue interface{}

const (
	ARRAY   = "array"
	BOOLEAN = "boolean"
	NULL    = "null"
	NUMBER  = "number"
	INTEGER = "integer"
	OBJECT  = "object"
	STRING  = "string"
)

type Property struct {
	Type       JsonValue       `json:"type"`
	Reference  string          `json:"$ref"`
	Properties json.RawMessage `json:"properties"`
	Items      []string        `json:"items"`
}

type JsonSchema struct {
	Type       JsonValue           `json:"type"`
	Reference  string              `json:"$ref"`
	Properties map[string]Property `json:"properties"`
}
