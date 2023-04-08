package geeorm

// Import the necessary packages
import (
	"testing"
)

func OpenDB(t *testing.T) *Engine {
	t.Helper()
	engine, err := NewEngine("sqlite3", "gee.db")
	if err != nil {
		t.Fatal("faliled to connect", err)
	}
	return engine
}

func TestNewEngine(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
}
