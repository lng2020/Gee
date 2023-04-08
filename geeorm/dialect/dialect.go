package dialect

import "reflect"

var dialectsMap = map[string]Dialect{}

type Dialect interface {
	// DataTypeOf returns the data type string of a Field
	DataTypeOf(typ reflect.Value) string
	// TableExistSQL returns the SQL statement to check if a table exists
	TableExistSQL(tableName string) (string, []interface{})
}

func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}
