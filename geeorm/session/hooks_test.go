package session

import (
	"goTinyToys/geeorm/log"
	"testing"
)

type Account struct {
	ID       int64
	Password string
}

func (a *Account) BeforeInsert(s *Session) error {
	log.Info("before insert", a)
	a.ID += 1000
	return nil
}

func (a *Account) AfterQuery(s *Session) error {
	log.Info("after query", a)
	a.Password = "******"
	return nil
}

// test hooks
func TestSession_CallMethod(t *testing.T) {
	s := NewSession().Model(&Account{})
	s.DropTable()
	s.CreateTable()
	s.Insert(&Account{1, "123456"}, &Account{2, "123456"})

	u := &Account{}
	if err := s.First(u); err != nil || u.Password != "******" {
		t.Fatal("failed to call hooks")
	}
}
