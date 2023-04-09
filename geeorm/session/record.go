package session

import (
	"database/sql"
	"goTinyToys/geeorm/clause"
	"goTinyToys/geeorm/log"
	"reflect"
)

// Insert function is used to insert a record into database
func (s *Session) Insert(values ...interface{}) (result sql.Result, err error) {
	recordValues := make([]interface{}, 0)
	for _, value := range values {
		table := s.Model(value).RefTable()
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
		recordValues = append(recordValues, table.RecordValues(value))
	}
	s.clause.Set(clause.VALUES, recordValues...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	if result, err = s.Raw(sql, vars...).Exec(); err != nil {
		log.Error(err)
	}
	return
}

// Find function is used to find records from database
func (s *Session) Find(values interface{}) (err error) {
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	destType := destSlice.Type().Elem()
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()

	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		log.Error(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		dest := reflect.New(destType).Elem()
		var values []interface{}
		for _, name := range table.FieldNames {
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		if err = rows.Scan(values...); err != nil {
			log.Error(err)
			return
		}
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return
}

// Where function is used to set where condition
func (s *Session) Where(values ...interface{}) *Session {
	s.clause.Set(clause.WHERE, values...)
	return s
}

// Update function is used to update records in database
func (s *Session) Update(values ...interface{}) (result sql.Result, err error) {
	table := s.Model(values).RefTable()
	for i := 0; i < len(values); i += 2 {
		s.clause.Set(clause.UPDATE, table.Name, values[i], values[i+1])
	}
	sql, vars := s.clause.Build(clause.UPDATE, clause.WHERE)
	if result, err = s.Raw(sql, vars...).Exec(); err != nil {
		log.Error(err)
	}

	return
}

// Delete function is used to delete records from database
func (s *Session) Delete() (result sql.Result, err error) {
	table := s.RefTable()
	s.clause.Set(clause.DELETE, table.Name)
	sql, vars := s.clause.Build(clause.DELETE, clause.WHERE)
	if result, err = s.Raw(sql, vars...).Exec(); err != nil {
		log.Error(err)
	}
	return
}

// Limit function is used to set limit condition
func (s *Session) Limit(limit int) *Session {
	s.clause.Set(clause.LIMIT, limit)
	return s
}

// OrderBy function is used to set order by condition
func (s *Session) OrderBy(orderBy string) *Session {
	s.clause.Set(clause.ORDERBY, orderBy)
	return s
}
