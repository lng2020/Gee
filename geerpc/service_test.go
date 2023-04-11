package geerpc

import (
	"reflect"
	"testing"
)

func TestNewService(t *testing.T) {
	arith := new(Arith)
	s := newService(arith)
	if s.name != "Arith" {
		t.Errorf("service name: expect %s, got %s", "Arith", s.name)
	}
	if len(s.methods) != 2 {
		t.Errorf("number of methods: expect %d, got %d", 2, len(s.methods))
	}
}

func TestCall(t *testing.T) {
	arith := new(Arith)
	s := newService(arith)
	mType := s.methods["Multiply"]
	if mType == nil {
		t.Fatal("method Multiply not found")
	}

	argv := reflect.New(mType.ArgType.Elem())
	replyv := reflect.New(mType.ReplyType.Elem())
	argv.Set(reflect.ValueOf(&Args{7, 8}))
	err := s.call("Multiply", argv, replyv)
	if err != nil {
		t.Fatal("call Multiply error:", err)
	}
	reply := replyv.Elem().Int()
	if reply != 56 {
		t.Fatalf("call Multiply: expect %d, got %d", 56, reply)
	}

	mType = s.methods["Add"]
	if mType == nil {
		t.Fatal("method Add not found")
	}
	argv = reflect.New(mType.ArgType.Elem())
	replyv = reflect.New(mType.ReplyType.Elem())
	argv.Set(reflect.ValueOf(&Args{7, 8}))
	err = s.call("Add", argv, replyv)
	if err != nil {
		t.Fatal("call Add error:", err)
	}
	reply = replyv.Elem().Int()
	if reply != 15 {
		t.Fatalf("call Add: expect %d, got %d", 15, reply)
	}
}
