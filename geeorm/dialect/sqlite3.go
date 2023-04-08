package dialect

import (
	"reflect"
	"time"
)

// sqlite3 struct that implements Dialect interface
type sqlite3 struct{}

var _ Dialect = (*sqlite3)(nil)

func init() {
	RegisterDialect("sqlite3", &sqlite3{})
}

// Return the data type of a field
func (s *sqlite3) DataTypeOf(typ reflect.Value) string {
	switch typ.Kind() {
	case reflect.Bool:
		return "bool"
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return "integer"
	case reflect.Int, reflect.Uint, reflect.Int64, reflect.Uint64:
		return "bigint"
	case reflect.Float32, reflect.Float64:
		return "real"
	case reflect.String:
		return "text"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "datetime"
		}
	}
	panic("invalid sql type")
}

// Generate TableExistSQL method to satisfy Dialect interface
func (s *sqlite3) TableExistSQL(tableName string) (string, []interface{}) {
	return "SELECT name FROM sqlite_master WHERE type='table' and name = ?", []interface{}{tableName}
}
