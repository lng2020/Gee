package session

import (
	"database/sql"
	"errors"
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

// Find function is used to find records from database. e.g s.Find(&users)
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
	// write query result to values
	for rows.Next() {
		dest := reflect.New(destType).Elem()
		var values []interface{}
		for _, name := range table.FieldNames {
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		if err = rows.Scan(values...); err != nil {
			log.Error(err)
		}
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return
}

// Update function is used to update records in database
func (s *Session) Update(kv ...interface{}) (result sql.Result, err error) {
	m := make(map[string]interface{})
	for i := 0; i < len(kv); i += 2 {
		m[kv[i].(string)] = kv[i+1]
	}
	s.clause.Set(clause.UPDATE, s.RefTable().Name, m)
	sql, vars := s.clause.Build(clause.UPDATE, clause.WHERE)
	if result, err = s.Raw(sql, vars...).Exec(); err != nil {
		log.Error(err)
	}
	return
}

// Delete function is used to delete records from database
func (s *Session) Delete() (result sql.Result, err error) {
	s.clause.Set(clause.DELETE, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.DELETE, clause.WHERE)
	if result, err = s.Raw(sql, vars...).Exec(); err != nil {
		log.Error(err)
	}
	return
}

// Count function is used to count records in database
func (s *Session) Count() (count int, err error) {
	s.clause.Set(clause.COUNT, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.COUNT, clause.WHERE)
	row := s.Raw(sql, vars...).QueryRow()
	if err = row.Scan(&count); err != nil {
		log.Error(err)
	}
	return
}

// Limit function is used to limit the number of records
func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT, num)
	return s
}

// OrderBy function is used to sort the records
func (s *Session) OrderBy(desc string) *Session {
	s.clause.Set(clause.ORDERBY, desc)
	return s
}

// Where function is used to filter the records
func (s *Session) Where(desc string, args ...interface{}) *Session {
	var vars []interface{}
	s.clause.Set(clause.WHERE, append(append(vars, desc), args...)...)
	return s
}

// first function is used to get the first record
func (s *Session) First(value interface{}) (err error) {
	dest := reflect.Indirect(reflect.ValueOf(value))
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()
	if err = s.Limit(1).Find(destSlice.Addr().Interface()); err != nil {
		log.Error(err)
		return
	}
	if destSlice.Len() == 0 {
		err = errors.New("not found")
	}
	dest.Set(destSlice.Index(0))
	return
}
