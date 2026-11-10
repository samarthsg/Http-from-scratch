package server

import (
	"fmt"
	"http-from-scratch/internal/request"
	"log"
	"net"
	"strconv"
	"sync/atomic"
)

type Server struct {
	Listener net.Listener
	Closed   atomic.Bool
}

func Serve(port int) (*Server, error) {
	addr := ":" + strconv.Itoa(port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	server := &Server{Listener: listener}
	go server.listen()
	return server, nil
}

func (s *Server) Close() error {
	s.Closed.Store(true)
	return s.Listener.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			if s.Closed.Load() {
				return
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Request line:")
	fmt.Println("- Method:", req.RequestLine.Method)
	fmt.Println("- Target:", req.RequestLine.RequestTarget)
	fmt.Println("- HTTP version:", req.RequestLine.HttpVersion)
	fmt.Println("Headers")
	req.Headers.ForEach(func(key, value string) {
		fmt.Printf("%s: %s\n", key, value)
	})
	fmt.Println("Body:")
	fmt.Println("BODY_STRING")
	fmt.Println(string(req.Body))
}
