package js2sql

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

const TAB = "    " // 4 spaces

var (
	defaultWriter = os.Stdout
	output        string
	writer        io.Writer
	tableName     string
)

type SchemaType struct {
	Value string
}

func (st *SchemaType) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		st.Value = str
		return nil
	}

	var arr []string
	if err := json.Unmarshal(data, &arr); err == nil {
		for _, t := range arr {
			if t != "null" {
				st.Value = t
				return nil
			}
		}

		st.Value = "string"
		return nil
	}

	return fmt.Errorf("type must be string or array of strings")
}

type JSONSchema struct {
	Type        SchemaType            `json:"type"`
	Properties  map[string]JSONSchema `json:"properties"`
	Items       *JSONSchema           `json:"items"`
	Format      string                `json:"format"`
	Title       string                `json:"title"`
	Ref         string                `json:"$ref"`
	Definitions map[string]JSONSchema `json:"definitions,omitempty"`
	Defs        map[string]JSONSchema `json:"$defs,omitempty"`
}

func init() {
	Js2SqlCmd.Flags().StringVarP(&output, "output", "o", "", "Output file path (default: stdout)")
	Js2SqlCmd.Flags().StringVarP(&tableName, "table", "t", "test", "Table name for the generated SQL")
}

var Js2SqlCmd = &cobra.Command{
	Use:   "js2sql [jsonschema file]",
	Short: "Generate a SQL table from JSON Schema",
	Long:  `Generate a SQL table from JSON Schema. Provide a path to a JSON Schema file to generate a SQL CREATE TABLE statement.`,
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		data, err := os.ReadFile(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}

		sql, err := process(data, tableName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error processing schema: %v\n", err)
			os.Exit(1)
		}

		if output == "" {
			writer = defaultWriter
		} else {
			f, err := os.Create(output)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
				os.Exit(1)
			}
			defer f.Close()
			writer = f
		}

		fmt.Fprintln(writer, sql)
	},
}

func process(data []byte, tableName string) (string, error) {
	var schema JSONSchema
	err := json.Unmarshal(data, &schema)
	if err != nil {
		return "", fmt.Errorf("error parsing JSON Schema: %w", err)
	}

	var columns []string
	resolvedSchemas := make(map[string]JSONSchema)

	// get defs
	if schema.Definitions != nil {
		for name, def := range schema.Definitions {
			resolvedSchemas["#/definitions/"+name] = def
		}
	}

	if schema.Defs != nil {
		for name, def := range schema.Defs {
			resolvedSchemas["#/$defs/"+name] = def
		}
	}

	// early return if no properties
	if len(schema.Properties) == 0 {
		return "", nil
	}

	columnMap := make(map[string]string)
	for name, prop := range schema.Properties {
		if name == "id" {
			continue
		}

		if prop.Ref != "" {
			refSchema, err := resolveRef(prop.Ref, schema, resolvedSchemas)
			if err != nil {
				return "", err
			}
			prop = *refSchema
		}

		columnType, err := mapType(&prop, schema, resolvedSchemas)
		if err != nil {
			return "", err
		}

		column := fmt.Sprintf("%s%s %s", TAB, name, columnType)
		columnMap[name] = column
	}

	// sorted cols
	var sortedNames []string
	for name := range columnMap {
		sortedNames = append(sortedNames, name)
	}
	sort.Strings(sortedNames)

	idColumn := TAB + "id varchar primary key"
	columns = append(columns, idColumn)
	for _, name := range sortedNames {
		columns = append(columns, columnMap[name])
	}

	createTable := fmt.Sprintf("create table %s (\n%s\n);",
		tableName,
		strings.Join(columns, ",\n"))

	return createTable, nil
}

func resolveRef(ref string, rootSchema JSONSchema, resolvedSchemas map[string]JSONSchema) (*JSONSchema, error) {
	if schema, ok := resolvedSchemas[ref]; ok {
		return &schema, nil
	}

	if !strings.HasPrefix(ref, "#/") {
		return nil, fmt.Errorf("only local references are supported: %s", ref)
	}

	if strings.HasPrefix(ref, "#/$defs/") {
		defName := strings.TrimPrefix(ref, "#/$defs/")
		if rootSchema.Defs != nil {
			if def, ok := rootSchema.Defs[defName]; ok {
				return &def, nil
			}
		}
		return nil, fmt.Errorf("reference not found: %s", ref)
	}

	if strings.HasPrefix(ref, "#/definitions/") {
		defName := strings.TrimPrefix(ref, "#/definitions/")
		if rootSchema.Definitions != nil {
			if def, ok := rootSchema.Definitions[defName]; ok {
				return &def, nil
			}
		}
		return nil, fmt.Errorf("reference not found: %s", ref)
	}

	return nil, fmt.Errorf("unsupported reference format: %s", ref)
}

func mapType(schema *JSONSchema, rootSchema JSONSchema, resolvedSchemas map[string]JSONSchema) (string, error) {
	if schema == nil {
		return "VARCHAR", nil
	}

	if schema.Ref != "" {
		refSchema, err := resolveRef(schema.Ref, rootSchema, resolvedSchemas)
		if err != nil {
			return "", err
		}
		schema = refSchema
	}

	switch schema.Type.Value {
	case "string":
		switch schema.Format {
		case "date":
			return "date", nil
		case "date-time":
			return "timestamp", nil
		case "time":
			return "time", nil
		default:
			return "varchar", nil
		}
	case "integer":
		return "integer", nil
	case "number":
		return "double", nil
	case "boolean":
		return "boolean", nil
	case "null":
		return "null", nil
	case "array":
		if schema.Items != nil {
			if schema.Items.Ref != "" {
				refSchema, err := resolveRef(schema.Items.Ref, rootSchema, resolvedSchemas)
				if err != nil {
					return "", err
				}
				schema.Items = refSchema
			}

			itemType, err := mapType(schema.Items, rootSchema, resolvedSchemas)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("array(%s)", itemType), nil
		}
		return "array(varchar)", nil
	case "object":
		fieldMap := make(map[string]string)
		for name, prop := range schema.Properties {
			if prop.Ref != "" {
				refSchema, err := resolveRef(prop.Ref, rootSchema, resolvedSchemas)
				if err != nil {
					return "", err
				}
				prop = *refSchema
			}

			fieldType, err := mapType(&prop, rootSchema, resolvedSchemas)
			if err != nil {
				return "", err
			}
			fieldMap[name] = fmt.Sprintf("%s %s", name, fieldType)
		}

		if len(fieldMap) == 0 {
			return "row()", nil
		}

		var sortedNames []string
		for name := range fieldMap {
			sortedNames = append(sortedNames, name)
		}
		sort.Strings(sortedNames)

		var fields []string
		for _, name := range sortedNames {
			fields = append(fields, fieldMap[name])
		}

		return fmt.Sprintf("row(%s)", strings.Join(fields, ", ")), nil
	default:
		return "varchar", nil
	}
}
