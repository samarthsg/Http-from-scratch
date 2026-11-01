package server

import (
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
	return nil
}

func (s *Server) listen() {
	conn, err := s.Listener.Accept()
	if err != nil {
		conn.Close()
	}
}

func (s *Server) handle(conn net.Conn) {

}
