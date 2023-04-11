package main

import (
	"encoding/json"
	"goTinyToys/geerpc"
	"goTinyToys/geerpc/codec"
	"log"
	"net"
)

func startServer(addr chan string) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error:", err)
	}

	log.Println("start rpc server on", l.Addr())
	addr <- l.Addr().String()
	geerpc.Accept(l)
}

func main() {
	addr := make(chan string)
	go startServer(addr)

	conn, _ := net.Dial("tcp", <-addr)
	defer func() { _ = conn.Close() }()

	_ = json.NewEncoder(conn).Encode(geerpc.DefaultOption)
	cc := codec.NewGobCodec(conn)
	for i := 0; i < 5; i++ {
		// send option
		_ = json.NewEncoder(conn).Encode(geerpc.DefaultOption)
		// send request
		h := &codec.Header{
			ServiceMethod: "Foo.Sum",
			Seq:           uint64(i),
		}
		_ = cc.Write(h, geerpc.DefaultOption)
		// receive response
		_ = cc.ReadHeader(h)
		var reply int
		_ = cc.ReadBody(&reply)
	}
}
