package clause

import (
	"strings"
)

type generator func(values ...interface{}) (string, []interface{})

var generators = map[Type]generator{}

func init() {
	generators[INSERT] = genInsert
	generators[VALUES] = genValues
	generators[SELECT] = genSelect
	generators[UPDATE] = genUpdate
	generators[DELETE] = genDelete
	generators[LIMIT] = genLimit
	generators[ORDERBY] = genOrderBy
	generators[WHERE] = genWhere
}

func genBindVar(num int) string {
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ", ")
}

func genInsert(values ...interface{}) (string, []interface{}) {
	tableName := values[0].(string)
	fields := strings.Join(values[1].([]string), ", ")
	return "INSERT INTO " + tableName + " (" + fields + ")", []interface{}{}
}

func genValues(values ...interface{}) (string, []interface{}) {
	var sql strings.Builder
	var vars []interface{}
	sql.WriteString("VALUES ")

	for i, value := range values {
		if i > 0 {
			sql.WriteString(", ")
		}
		sql.WriteString("(")
		sql.WriteString(genBindVar(len(value.([]interface{}))))
		sql.WriteString(")")
		vars = append(vars, value.([]interface{})...)
	}

	return sql.String(), vars
}

func genSelect(values ...interface{}) (string, []interface{}) {
	tableName := values[0].(string)
	fields := strings.Join(values[1].([]string), ", ")
	return "SELECT " + fields + " FROM " + tableName, []interface{}{}
}

func genUpdate(values ...interface{}) (string, []interface{}) {
	tableName := values[0].(string)
	var sql strings.Builder
	var vars []interface{}
	sql.WriteString("UPDATE " + tableName + " SET ")

	for i, field := range values[1].([]string) {
		if i > 0 {
			sql.WriteString(", ")
		}
		sql.WriteString(field + " = ?")
		vars = append(vars, values[i+2])
	}

	return sql.String(), vars
}

func genDelete(values ...interface{}) (string, []interface{}) {
	return "DELETE FROM " + values[0].(string), []interface{}{}
}

func genLimit(values ...interface{}) (string, []interface{}) {
	return "LIMIT ?", values
}

func genOrderBy(values ...interface{}) (string, []interface{}) {
	return "ORDER BY " + values[0].(string), []interface{}{}
}

func genWhere(values ...interface{}) (string, []interface{}) {
	var sql strings.Builder
	var vars []interface{}
	sql.WriteString("WHERE ")

	for i, expr := range values {
		if i%2 == 0 {
			if i > 0 {
				sql.WriteString(" AND ")
			}
			sql.WriteString(expr.(string))
		} else {
			vars = append(vars, expr)
		}
	}
	return sql.String(), vars
}
