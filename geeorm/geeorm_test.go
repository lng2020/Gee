package geeorm

// Import the necessary packages
import (
	"errors"
	"goTinyToys/geeorm/session"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Name string `geeorm:"PRIMARY KEY"`
	Age  int
}

func OpenDB(t *testing.T) *Engine {
	t.Helper()
	engine, err := NewEngine("sqlite3", "geeorm.db")
	if err != nil {
		t.Fatal("faliled to connect", err)
	}
	return engine
}

func TestNewEngine(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
}

func Test_Transaction(t *testing.T) {
	t.Run("Rollback", func(t *testing.T) {
		engine := OpenDB(t)
		defer engine.Close()

		s := engine.NewSession()
		_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()

		_, err := engine.Transaction(func(s *session.Session) (interface{}, error) {
			_ = s.Model(&User{}).CreateTable()
			_, _ = s.Insert(&User{"Tom", 18})
			return nil, errors.New("Error")
		})

		if err == nil || s.HasTable() {
			t.Fatal("Failed to rollback")
		}
	})
	t.Run("Commit", func(t *testing.T) {
		engine := OpenDB(t)
		defer engine.Close()

		s := engine.NewSession()
		_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()

		res, err := engine.Transaction(func(s *session.Session) (interface{}, error) {
			_ = s.Model(&User{}).CreateTable()
			_, _ = s.Insert(&User{"Tom", 18})
			return s.Count()
		})
		count, _ := res.(int)
		if err != nil || count != 1 {
			t.Fatal("Failed to commit", err, count)
		}
	})
}
