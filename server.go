package main

import (
	"context"
	"distributed-cache/cache"
	"fmt"
	"log"
	"net"
)

type ServerOptions struct {
	ListenAddr  string
	IsLeader    bool
	leaderAdder string
}

type Server struct {
	ops       ServerOptions
	Followers map[net.Conn]struct{}
	cache     cache.Cache
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
	if !s.ops.IsLeader {
		conn, err := net.Dial("tcp", s.ops.leaderAdder)
		if err != nil {
			log.Fatalf("failed to dial: %v", err)
		}

		s.Followers[conn] = struct{}{}

	}
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

		go s.handleCommand(conn, buf[:n])
	}
}

func (s *Server) handleCommand(conn net.Conn, rawCmd []byte) {

	msg, err := parseMessage(rawCmd)
	if err != nil {
		conn.Write([]byte("invalid command\n"))
		return
	}

	switch msg.Cmd {
	case CMDset:
		if err := s.handleSetCommand(conn, msg); err != nil {
			conn.Write([]byte("failed to execute SET command\n"))
		}
	case CMDget:
		if err := s.handleGetCommand(conn, msg); err != nil {
			conn.Write([]byte("failed to execute GET command\n"))
		}
	}
}

func (s *Server) handleSetCommand(conn net.Conn, msg *Message) error {
	if err := s.cache.Set(msg.Key, msg.Value, msg.TTL); err != nil {
		return err
	}

	go s.sedntoFollowers(context.TODO(), msg)

	return nil
}

func (s *Server) handleGetCommand(conn net.Conn, msg *Message) error {
	value, err := s.cache.Get(msg.Key)
	if err != nil {
		fmt.Println("failed to get value from cache: ", err)
		return err
	}

	_, err = conn.Write(value)
	if err != nil {
		fmt.Println("failed to write to connection: ", err)
		return err
	}

	return nil
}

func (s *Server) sedntoFollowers(ctx context.Context, msg *Message) {
	for conn := range s.Followers {
		_, err := conn.Write(msg.Value)
		if err != nil {
			log.Printf("failed to write to connection: %v", err)
		}
	}
}

func (s *Server) ToBytes(msg *Message) []byte {
	cmd := fmt.Sprintf("%s %s %s %d", msg.Cmd, msg.Key, msg.Value, msg.TTL)
	return []byte(cmd)
}
