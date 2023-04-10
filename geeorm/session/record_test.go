package session

import (
	"testing"
)

func testRecordInit(t *testing.T) *Session {
	t.Helper()
	s := NewSession().Model(&User{})
	s.DropTable()
	s.CreateTable()
	return s
}

// unit test for Insert function in record.go
func TestInsert(t *testing.T) {
	s := testRecordInit(t)
	result, err := s.Insert(&User{"Tom", 18}, &User{"Sam", 25})
	if err != nil || result == nil {
		t.Fatal("failed to insert data into database")
	}
}

// unit test for Find function in record.go
func TestFind(t *testing.T) {
	s := testRecordInit(t)
	_, _ = s.Insert(&User{"Tom", 18}, &User{"Sam", 25})
	var users []User
	if err := s.Find(&users); err != nil || len(users) != 2 {
		t.Fatal("failed to query data from database")
	}
	t.Log(users)
}

// unit test for Update function in record.go
func TestUpdate(t *testing.T) {
	s := testRecordInit(t)
	_, _ = s.Insert(&User{"Tom", 18}, &User{"Sam", 25})
	_, err := s.Where("Name = ?", "Tom").Update("Age", 20)
	if err != nil {
		t.Fatal("failed to update data in database")
	}
	var user User
	if err = s.First(&user); err != nil || user.Age != 20 {
		t.Fatal("failed to update data in database")
	}
	t.Log(user)
}

// unit test for Delete function in record.go
func TestDelete(t *testing.T) {
	s := testRecordInit(t)
	_, _ = s.Insert(&User{"Tom", 18}, &User{"Sam", 25})
	_, err := s.Where("Name = ?", "Tom").Delete()
	if err != nil {
		t.Fatal("failed to delete data from database")
	}
	var user User
	if err = s.First(&user); err != nil || user.Name != "Sam" {
		t.Fatal("failed to delete data from database")
	}
	t.Log(user)
}

// unit test for Count function in record.go
func TestCount(t *testing.T) {
	s := testRecordInit(t)
	_, _ = s.Insert(&User{"Tom", 18}, &User{"Sam", 25})
	n, err := s.Count()
	if n != 2 || err != nil {
		t.Fatal("failed to count data in database")
	}
}

// unit test for Limit function in record.go
func TestLimit(t *testing.T) {
	s := testRecordInit(t)
	_, _ = s.Insert(&User{"Tom", 18}, &User{"Sam", 25})
	var users []User
	if err := s.Limit(1).Find(&users); err != nil || len(users) != 1 {
		t.Fatal("failed to limit data in database")
	}
	t.Log(users)
}

// unit test for OrderBy function in record.go
func TestOrderBy(t *testing.T) {
	s := testRecordInit(t)
	_, _ = s.Insert(&User{"Tom", 18}, &User{"Sam", 25})
	var users []User
	if err := s.OrderBy("Age DESC").Find(&users); err != nil || users[0].Name != "Sam" {
		t.Fatal("failed to order data in database")
	}
	t.Log(users)
}
