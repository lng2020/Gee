package schema

import (
	"testing"

	"goTinyToys/geeorm/dialect"

	"github.com/stretchr/testify/assert"
)

type User struct {
	Name string `geeorm:"PRIMARY KEY"`
	Age  int
}

func TestParse(t *testing.T) {
	d, _ := dialect.GetDialect("sqlite3") // Replace with your dialect implementation
	schema := Parse(&User{}, d)

	assert.Equal(t, "User", schema.Name)
	assert.Equal(t, []string{"Name", "Age"}, schema.FieldNames)
	assert.Equal(t, "text", schema.FieldMap["Name"].Type)
	assert.Equal(t, "bigint", schema.FieldMap["Age"].Type)
	assert.Equal(t, "PRIMARY KEY", schema.FieldMap["Name"].Tag)
}
