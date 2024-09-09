package main

import (
	"distributed-cache/cache"
	"flag"
	"log"
	"net"
	"time"
)

func main() {
	listenAdder := flag.String("listen", ":8080", "server listen address")
	leaderAdder := flag.String("leader", "", "leader address")
	ops := ServerOptions{
		ListenAddr:  *listenAdder,
		IsLeader:    true,
		leaderAdder: *leaderAdder,
	}

	// go run main.go -listen :8080 -leader :8081

	go func() {
		conn, err := net.Dial("tcp", ":8080")
		if err != nil {
			log.Fatalf("failed to dial: %v", err)
		}

		conn.Write([]byte("SET key1 value1 10"))

		time.Sleep(1 * time.Second)

		conn.Write([]byte("GET key1"))
		buf := make([]byte, 1024)
		n, _ := conn.Read(buf)
		log.Printf("Response: %s", string(buf[:n]))

	}()

	s := NewServer(ops, cache.NewCache())
	s.Serve()

}
