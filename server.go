package main

import (
	"distributed-cache/cache"
	"log"
	"net"
)

type ServerOptions struct {
	ListenAddr string
	IsLeader   bool
}

type Server struct {
	ops   ServerOptions
	cache cache.Cache
}

func NewServer(ops ServerOptions, c cache.Cache) *Server {
	return &Server{
		ops:   ops,
		cache: c,
	}
}

func (s *Server) Serve() {
	listner, err := net.Listen("tcp", s.ops.ListenAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("listening on %s", s.ops.ListenAddr)

	for {
		conn, err := listner.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %v", err)
		}

		go s.connHandler(conn)
	}
}

func (s *Server) connHandler(conn net.Conn) {

	defer func() {
		conn.Close()
	}()

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("failed to read from connection: %v", err)
			break
		}

		msg := string(buf[:n])
		log.Printf("received message: %s", msg)
	}
}
