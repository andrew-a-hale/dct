package jsonschema2sql

import (
	"encoding/json"
	"testing"
)

// upgrade these to create sql definitions
// create table schema.release_schema ( id varchar primary key, name varchar, age int )
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

	sql, err := process(raw)
	if err != nil {
		t.Errorf("failed to process json: %v", err)
	}

	expected := `\
create table test (
  id varchar primary key
  , name varchar
  , age int
)`
	if sql != expected {
		t.Errorf("invalid sql generated:\n`%v`\nexpected:\n`%v`", sql, expected)
	}
}

// upgrade these to create sql definitions
// create table schema.release_schema ( id varchar primary key, name varchar, age row(year int) )
func TestNewJsonSchemaComplex(t *testing.T) {
	raw := []byte(`{
  "type": "object",
  "properties": {
    "name": { "type": "string" },
    "age": { "type": "object", "properties": { "year": { "type": "integer" } } }
  }
}`)
	var j JsonSchema
	err := json.Unmarshal(raw, &j)
	if err != nil {
		t.Errorf("failed to unmarshal json: %v", err)
	}

	nested, err := j.Properties["age"].Properties.MarshalJSON()
	if err != nil {
		t.Errorf("failed to marshal nested json: %v", err)
	}

	var value Property
	err = json.Unmarshal(nested, &value)
	if err != nil {
		t.Errorf("failed to unmarshal json: %v", err)
	}

	sql, err := process(raw)
	if err != nil {
		t.Errorf("failed to process json: %v", err)
	}

	expected := `\
create table test (
  id varchar primary key
  , name varchar
  , age row(year int)
)`
	if sql != expected {
		t.Errorf("invalid sql generated:\n`%v`\nexpected:\n`%v`", sql, expected)
	}
}

func TestNewJsonSchemaReference(t *testing.T) {
	raw := []byte(`{
  "type": "object",
  "properties": {
    "name": { "type": "string" },
    "age": { "$ref": "#/properties/name" }
  }
}`)

	var j JsonSchema
	err := json.Unmarshal(raw, &j)
	if err != nil {
		t.Errorf("failed to unmarshal json: %v", err)
	}

	sql, err := process(raw)
	if err != nil {
		t.Errorf("failed to process json: %v", err)
	}

	expected := `\
create table test (
  id varchar primary key
  , name varchar
  , age varchar
)`
	if sql != expected {
		t.Errorf("invalid sql generated:\n`%v`\nexpected:\n`%v`", sql, expected)
	}
}

func TestNewJsonSchemaReferenceComplex(t *testing.T) {
	raw := []byte(`{
  "type": "object",
  "properties": {
    "name": { "type": "string" },
    "home": { "$ref": "#/definitions/address" }
  },
  "definitions": {
    "address": {
      "type": "object",
      "properties": {
        "street": { "type": "string" },
        "city": { "type": "string" }
      }
    }
  }
}`)

	var j JsonSchema
	err := json.Unmarshal(raw, &j)
	if err != nil {
		t.Errorf("failed to unmarshal json: %v", err)
	}

	sql, err := process(raw)
	if err != nil {
		t.Errorf("failed to process json: %v", err)
	}

	expected := `\
create table test (
  id varchar primary key
  , name varchar
  , home row(street varchar, city varchar)
)`
	if sql != expected {
		t.Errorf("invalid sql generated:\n`%v`\nexpected:\n`%v`", sql, expected)
	}
}

func TestNewJsonSchemaReferenceArray(t *testing.T) {
	raw := []byte(`{
   "type": "object",
   "properties": {
     "name": { "type": "string" },
     "houses": { "type": "array", "properties": { "address": { "$ref": "#/definitions/address" } } }
   },
   "definitions": {
     "address": {
       "type": "object",
       "properties": {
         "street": { "type": "string" },
         "city": { "type": "string" }
       }
     }
   }
 }`)

	var j JsonSchema
	err := json.Unmarshal(raw, &j)
	if err != nil {
		t.Errorf("failed to unmarshal json: %v", err)
	}

	sql, err := process(raw)
	if err != nil {
		t.Errorf("failed to process json: %v", err)
	}

	expected := `\
create table test (
  id varchar primary key
  , name varchar
  , houses array(row(street varchar, city varchar))
)`
	if sql != expected {
		t.Errorf("invalid sql generated:\n`%v`\nexpected:\n`%v`", sql, expected)
	}
}

func TestNewJsonSchemaNestedArray(t *testing.T) {
	raw := []byte(`{
   "type": "object",
   "properties": {
     "name": { "type": "string" },
     "houses": { "type": "array", "properties": { "address": { "type": "string" } } }
   }
 }`)

	var j JsonSchema
	err := json.Unmarshal(raw, &j)
	if err != nil {
		t.Errorf("failed to unmarshal json: %v", err)
	}

	sql, err := process(raw)
	if err != nil {
		t.Errorf("failed to process json: %v", err)
	}

	expected := `create table test (
        id varchar primary key
        , name varchar
        , houses array(row(address varchar))
)`
	if sql != expected {
		t.Errorf("invalid sql generated:\n`%v`\nexpected:\n`%v`", sql, expected)
	}
}
