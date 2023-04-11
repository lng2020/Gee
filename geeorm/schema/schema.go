package schema

import (
	"goTinyToys/geeorm/dialect"
	"reflect"
)

// The following code defines a struct named Field, which represents a column in a database table.
type Field struct {
	Name string // Name of the column
	Type string // Data type of the column
	Tag  string // Additional metadata for the column
}

// The Schema struct in this file maybe miss an attribute field map
type Schema struct {
	Model      interface{}       // Model represents the object that is being converted to a table
	Name       string            // Name of the table
	Fields     []*Field          // Fields represents the columns in the table
	FieldNames []string          // FieldNames represents the names of the columns in the table
	FieldMap   map[string]*Field // FieldMap represents a map of column names to their corresponding Field objects
}

func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:    dest,
		Name:     modelType.Name(),
		FieldMap: make(map[string]*Field),
	}
	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		field := &Field{
			Name: p.Name,
			Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			Tag:  p.Tag.Get("geeorm"),
		}
		schema.Fields = append(schema.Fields, field)
		schema.FieldNames = append(schema.FieldNames, p.Name)
		schema.FieldMap[p.Name] = field
	}
	return schema
}

// function turn record to fields
func (schema *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}
	for _, field := range schema.Fields {
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}

// function get field by name
func (schema *Schema) GetField(name string) *Field {
	return schema.FieldMap[name]
}
