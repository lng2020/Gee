package clause

import (
	"reflect"
	"testing"
)

// test genSelect function in clause.go
func TestGenSelect(t *testing.T) {
	var clause Clause
	clause.Set(SELECT, "name", "age")
	sql, vars := clause.Build(SELECT)
	t.Log(sql, vars)
	if sql != "SELECT name, age" {
		t.Fatal("failed to build SQL")
	}
	if len(vars) > 0 {
		t.Fatal("failed to build SQL vars")
	}
}

// test genInsert function in clause.go
func TestGenInsert(t *testing.T) {
	var clause Clause
	clause.Set(INSERT, "User", []string{"name", "age"})
	clause.Set(VALUES, []interface{}{"Tom", 18})
	sql, vars := clause.Build(INSERT, VALUES)
	t.Log(sql, vars)
	if sql != "INSERT INTO User (name, age) VALUES (?, ?)" {
		t.Fatal("failed to build SQL")
	}
	if !reflect.DeepEqual(vars, []interface{}{"Tom", 18}) {
		t.Fatal("failed to build SQL vars")
	}
}

// test genValues function in clause.go
func TestGenValues(t *testing.T) {
	var clause Clause
	clause.Set(INSERT, "User", []string{"name", "age"})
	clause.Set(VALUES, []interface{}{"Tom", 18}, []interface{}{"Sam", 20})
	sql, vars := clause.Build(VALUES)
	t.Log(sql, vars)
	if sql != "VALUES (?, ?), (?, ?)" {
		t.Fatal("failed to build SQL")
	}
	if !reflect.DeepEqual(vars, []interface{}{"Tom", 18, "Sam", 20}) {
		t.Fatal("failed to build SQL vars")
	}
}

// test genUpdate function in clause.go
func TestGenUpdate(t *testing.T) {
	var clause Clause
	clause.Set(UPDATE, "User")
	clause.Set(SET, []string{"name", "age"}, []interface{}{"Tom", 18})
	sql, vars := clause.Build(UPDATE, SET)
	t.Log(sql, vars)
	if sql != "UPDATE User SET name = ?, age = ?" {
		t.Fatal("failed to build SQL")
	}
	if !reflect.DeepEqual(vars, []interface{}{"Tom", 18}) {
		t.Fatal("failed to build SQL vars")
	}
}

// test genDelete function in clause.go
func TestGenDelete(t *testing.T) {
	var clause Clause
	clause.Set(DELETE, "User")
	clause.Set(WHERE, "name = ?", "Tom")
	sql, vars := clause.Build(DELETE, WHERE)
	t.Log(sql, vars)
	if sql != "DELETE FROM User WHERE name = ?" {
		t.Fatal("failed to build SQL")
	}
	if !reflect.DeepEqual(vars, []interface{}{"Tom"}) {
		t.Fatal("failed to build SQL vars")
	}
}

// test genLimit function in clause.go
func TestGenLimit(t *testing.T) {
	var clause Clause
	clause.Set(LIMIT, 10)
	sql, vars := clause.Build(LIMIT)
	t.Log(sql, vars)
	if sql != "LIMIT ?" {
		t.Fatal("failed to build SQL")
	}
	if !reflect.DeepEqual(vars, []interface{}{10}) {
		t.Fatal("failed to build SQL vars")
	}
}

// test genOrderBy function in clause.go
func TestGenOrderBy(t *testing.T) {
	var clause Clause
	clause.Set(ORDERBY, "age DESC")
	sql, vars := clause.Build(ORDERBY)
	t.Log(sql, vars)
	if sql != "ORDER BY age DESC" {
		t.Fatal("failed to build SQL")
	}
	if len(vars) > 0 {
		t.Fatal("failed to build SQL vars")
	}
}

// test genWhere function in clause.go
func TestGenWhere(t *testing.T) {
	var clause Clause
	clause.Set(WHERE, "name = ?", "Tom")
	clause.Set(WHERE, "age > ?", 18)
	sql, vars := clause.Build(WHERE)
	t.Log(sql, vars)
	if sql != "WHERE name = ? AND age > ?" {
		t.Fatal("failed to build SQL")
	}
	if !reflect.DeepEqual(vars, []interface{}{"Tom", 18}) {
		t.Fatal("failed to build SQL vars")
	}
}

// test Build function in clause.go
func TestBuild(t *testing.T) {
	var clause Clause
	clause.Set(SELECT, "name", "age")
	clause.Set(WHERE, "name = ?", "Tom")
	clause.Set(ORDERBY, "age DESC")
	clause.Set(LIMIT, 10)
	sql, vars := clause.Build(SELECT, WHERE, ORDERBY, LIMIT)
	t.Log(sql, vars)
	if sql != "SELECT name, age FROM User WHERE name = ? ORDER BY age DESC LIMIT ?" {
		t.Fatal("failed to build SQL")
	}
	if !reflect.DeepEqual(vars, []interface{}{"Tom", 10}) {
		t.Fatal("failed to build SQL vars")
	}
}

// test set function in clause.go
func TestSet(t *testing.T) {
	var clause Clause
	clause.Set(SELECT, "name", "age")
	clause.Set(WHERE, "name = ?", "Tom")
	clause.Set(ORDERBY, "age DESC")
	clause.Set(LIMIT, 10)
	sql, vars := clause.Build(SELECT, WHERE, ORDERBY, LIMIT)
	t.Log(sql, vars)
	if sql != "SELECT name, age FROM User WHERE name = ? ORDER BY age DESC LIMIT ?" {
		t.Fatal("failed to build SQL")
	}
	if !reflect.DeepEqual(vars, []interface{}{"Tom", 10}) {
		t.Fatal("failed to build SQL vars")
	}
}
