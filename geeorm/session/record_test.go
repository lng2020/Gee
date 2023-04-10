package session

import (
	"testing"
)

// unit test for Insert function in record.go
func TestInsert(t *testing.T) {
	type User struct {
		Name string `geeorm:"PRIMARY KEY"`
		Age  int
	}
	s := NewSession().Model(&User{})
	s.DropTable()
	s.CreateTable()
	result, err := s.Insert(&User{"Tom", 18}, &User{"Sam", 25})
	if err != nil || result == nil {
		t.Fatal("failed to insert data into database")
	}
}

// unit test for Delete function in record.go
func TestDelete(t *testing.T) {
	type User struct {
		Name string `geeorm:"PRIMARY KEY"`
		Age  int
	}
	s := NewSession().Model(&User{})
	s.DropTable()
	s.CreateTable()
	_, _ = s.Insert(&User{"Tom", 18}, &User{"Sam", 25})
	result, err := s.Where("Name = ?", "Tom").Delete()
	if err != nil || result == nil {
		t.Fatal("failed to delete data from database")
	}
}

// unit test for Find function in record.go
func TestFind(t *testing.T) {
	type User struct {
		Name string `geeorm:"PRIMARY KEY"`
		Age  int
	}
	s := NewSession().Model(&User{})
	s.DropTable()
	s.CreateTable()
	_, _ = s.Insert(&User{"Tom", 18}, &User{"Sam", 25})
	var users []User
	if err := s.Find(&users); err != nil || len(users) != 2 {
		t.Fatal("failed to query data from database")
	}
	t.Log(users)
}
