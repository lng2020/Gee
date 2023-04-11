package session

import "goTinyToys/geeorm/log"

func (s *Session) Begin() (err error) {
	log.Info("transaction begin")
	s.tx, err = s.db.Begin()
	return
}

func (s *Session) Commit() (err error) {
	log.Info("transaction commit")
	err = s.tx.Commit()
	return
}

func (s *Session) Rollback() (err error) {
	log.Info("transaction rollback")
	err = s.tx.Rollback()
	return
}
