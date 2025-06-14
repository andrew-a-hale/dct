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
}

// upgrade these to create sql definitions
// create table schema.release_schema ( id varchar primary key, name varchar, age_year int )
func TestNewJsonSchemaComplex(t *testing.T) {
	raw := []byte(`{
  "type": "object",
  "properties": {
    "name": { "type": "string" },
    "age": { "type": "object", "properties": { "year": { "type": int } } }
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
}

// upgrade these to create sql definitions
// create table schema.release_schema ( id varchar primary key, name varchar, age varchar )
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
}

// upgrade these to create sql definitions
// create table schema.release_schema ( id varchar primary key, name varchar, address_street varchar, address_city varchar )
func TestNewJsonSchemaReferenceComplexDenormalised(t *testing.T) {
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
}

// upgrade these to create sql definitions
// create table schema.release_schema ( id varchar primary key, name varchar, address_id string )
// create table schema.release_schema_addresses ( address_id string primary key, address_street string, address_city string )
func TestNewJsonSchemaReferenceComplexNormalised(t *testing.T) {
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
}

// upgrade these to create sql definitions
// create table schema.release_schema ( id varchar primary key, name varchar, address_id string )
// create table schema.release_schema_addresses ( address_id string primary key, address_street string, address_city string )
func TestNewJsonSchemaReferenceComplexNormalisedArray(t *testing.T) {
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
}

// upgrade these to create sql definitions
// create table schema.release_schema ( id varchar primary key, name varchar, address_id string )
// create table schema.release_schema_addresses ( address_id string primary key, address string )
func TestNewJsonSchemaReferenceArray(t *testing.T) {
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
}
