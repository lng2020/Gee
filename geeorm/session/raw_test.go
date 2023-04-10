package session

import (
	"database/sql"
	"goTinyToys/geeorm/dialect"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

var TestDB *sql.DB
var TestDialect, _ = dialect.GetDialect("sqlite3")

func TestMain(m *testing.M) {
	TestDB, _ = sql.Open("sqlite3", "../geeorm.db")
	code := m.Run()
	_ = TestDB.Close()
	os.Exit(code)
}

func NewSession() *Session {
	return New(TestDB, TestDialect)
}

func TestSession_Exec(t *testing.T) {
	s := NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	result, _ := s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	if count, err := result.RowsAffected(); err != nil || count != 2 {
		t.Fatal("expect 2, but got", count)
	}
}

// TestSession_QueryRows tests the QueryRows function in raw.go
func TestSession_QueryRows(t *testing.T) {
	s := NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	_, _ = s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()

	// QueryRows returns a Rows object and an error
	rows, err := s.Raw("SELECT * FROM User").QueryRows()
	if err != nil {
		t.Fatal("failed to query rows:", err)
	}

	// Iterate through the rows and scan the values into variables
	var names []string
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			t.Fatal("failed to scan row:", err)
		}
		names = append(names, name)
	}

	// Check that the correct number of rows were returned
	if len(names) != 2 {
		t.Fatalf("expected 2 rows, but got %d", len(names))
	}

	// Check that the correct values were returned
	if names[0] != "Tom" || names[1] != "Sam" {
		t.Fatalf("expected names to be [Tom, Sam], but got %v", names)
	}
}
