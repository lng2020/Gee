package geeorm

import (
	"database/sql"
	"goTinyToys/geeorm/dialect"
	"goTinyToys/geeorm/log"
	"goTinyToys/geeorm/session"
)

type Engine struct {
	db *sql.DB
}

func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return
	}

	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}

	e = &Engine{db: db}
	log.Info("Connect database success")
	return
}

func (engine *Engine) Close() {
	if err := engine.db.Close(); err != nil {
		log.Error("Failed to close database")
	}
	log.Info("Close database success")
}

func (engine *Engine) NewSession() *session.Session {
	dialect, _ := dialect.GetDialect("sqlite3")
	return session.New(engine.db, dialect)
}

type TxFunc func(*session.Session) (interface{}, error)

func (engine *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := engine.NewSession()
	if err = s.Begin(); err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			_ = s.Rollback()
			panic(p)
		} else if err != nil {
			_ = s.Rollback()
		} else {
			err = s.Commit()
		}
	}()

	return f(s)
}

func difference(a []string, b []string) (diff []string) {
	mapB := make(map[string]bool)
	for _, item := range b {
		mapB[item] = true
	}
	for _, item := range a {
		if _, ok := mapB[item]; !ok {
			diff = append(diff, item)
		}
	}
	return
}

// Migrate migrate value's table
func (engine *Engine) Migrate(value interface{}) error {
	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		if !s.Model(value).HasTable() {
			log.Infof("table %s doesn't exist", s.RefTable().Name)
			return nil, s.CreateTable()
		}

		table := s.RefTable()
		columns, err := s.Raw("PRAGMA table_info(`" + table.Name + "`);").QueryRows()
		if err != nil {
			return nil, err
		}

		var columnNames []string
		for columns.Next() {
			var (
				name, tpe string
				notNull   bool
				dfltValue interface{}
				pk        int
			)
			if err = columns.Scan(&pk, &name, &tpe, &notNull, &dfltValue, &pk); err != nil {
				return nil, err
			}
			columnNames = append(columnNames, name)
		}

		addCols := difference(table.FieldNames, columnNames)
		delCols := difference(columnNames, table.FieldNames)
		log.Infof("added cols %v, deleted cols %v", addCols, delCols)
		for _, colName := range addCols {
			field := table.GetField(colName)
			if _, err = s.Raw("ALTER TABLE `" + table.Name + "` ADD COLUMN `" + colName + "` " + field.Type + ";").Exec(); err != nil {
				return
			}
		}

		if len(delCols) == 0 {
			return
		}

		return
	})
	return err
}
