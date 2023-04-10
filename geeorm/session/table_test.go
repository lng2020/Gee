package session

import (
	"testing"
)

type User struct {
	Name string `geeorm:"PRIMARY KEY"`
	Age  int
}

// test function for CreateTable()
func TestSession_CreateTable(t *testing.T) {
	s := NewSession().Model(&User{})
	_ = s.DropTable()
	if err := s.CreateTable(); err != nil {
		t.Fatal(err)
	}
	if !s.HasTable() {
		t.Fatalf("failed to create table for %s", s.RefTable().Name)
	}
}

// test function for DropTable()
func TestSession_DropTable(t *testing.T) {
	s := NewSession().Model(&User{})
	if err := s.DropTable(); err != nil {
		t.Fatal(err)
	}
	if s.HasTable() {
		t.Fatalf("failed to drop table %s", s.RefTable().Name)
	}
}
