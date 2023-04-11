package main

import (
	"goTinyToys/geerpc"
	"log"
	"net"
	"sync"
	"time"
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
	log.SetFlags(0)
	addr := make(chan string)
	go startServer(addr)
	client, _ := geerpc.Dial("tcp", <-addr)
	defer func() { _ = client.Close() }()

	time.Sleep(time.Second)
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := &geerpc.Args{A: i, B: i * i}
			var reply int
			err := client.Call("Arith.Multiply", args, &reply)
			if err != nil {
				log.Println("call Arith.Multiply error:", err)
			} else {
				log.Printf("%d * %d = %d", args.A, args.B, reply)
			}
		}(i)
	}

	wg.Wait()
}
