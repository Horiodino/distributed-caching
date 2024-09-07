package main

import (
	"distributed-cache/cache"
)

func main() {
	ops := ServerOptions{
		ListenAddr: ":8080",
		IsLeader:   true,
	}

	// go func() {
	// 	time.Sleep(1 * time.Second)
	// 	conn, err := net.Dial("tcp", ":8080")
	// 	if err != nil {
	// 		log.Fatalf("failed to dial: %v", err)
	// 	}

	// 	// conn.Write([]byte("hello"))
	// }()

	s := NewServer(ops, cache.NewCache())
	s.Serve()

}
