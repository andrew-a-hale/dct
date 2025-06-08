package jsonschema2sql

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestNewJsonSchema(t *testing.T) {
	raw := []byte(`{
  "type": "object",
  "properties": {
    "name": { "type": "string" },
    "age": { "type": "integer" }
  }
}`)
	var j JsonSchema
	err := json.Unmarshal(raw, &j)
	if err != nil {
		t.Errorf("failed to unmarshal json: %v", err)
	}

	fmt.Println(j)
}
