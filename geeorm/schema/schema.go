package schema

import (
	"go/ast"
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
	schema := &Schema{
		Model:    dest,
		Name:     reflect.TypeOf(dest).Elem().Name(),
		FieldMap: make(map[string]*Field),
	}
	val := reflect.ValueOf(dest).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		if !field.Anonymous && ast.IsExported(field.Name) {
			schema.FieldNames = append(schema.FieldNames, field.Name)
			f := &Field{
				Name: field.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(field.Type))),
			}
			if v, ok := field.Tag.Lookup("geeorm"); ok {
				f.Tag = v
			}
			schema.Fields = append(schema.Fields, f)
			schema.FieldMap[field.Name] = f
		}
	}
	return schema
}
