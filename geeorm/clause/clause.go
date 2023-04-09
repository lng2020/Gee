package clause

import "strings"

type Clause struct {
	sql     map[Type]string
	sqlVars map[Type][]interface{}
}

type Type int

const (
	INSERT Type = iota
	VALUES
	SELECT
	UPDATE
	SET
	DELETE
	LIMIT
	ORDERBY
	WHERE
)

func (c *Clause) Set(name Type, values ...interface{}) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlVars = make(map[Type][]interface{})
	}
	sql, vars := generators[name](values...)
	c.sql[name] = sql
	c.sqlVars[name] = vars
}

func (c *Clause) Build(names ...Type) (string, []interface{}) {
	var sqls []string
	var vars []interface{}
	for _, name := range names {
		sqls = append(sqls, c.sql[name])
		vars = append(vars, c.sqlVars[name]...)
	}
	return strings.Join(sqls, " "), vars
}
