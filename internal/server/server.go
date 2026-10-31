package server

import (
	"net"
	"strconv"
)

type Server struct {
	Listener net.Listener
}

func Serve(port int) (*Server, error) {
	addr := ":" + strconv.Itoa(port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Server{Listener: listener}, nil
}

func (s *Server) Close() error {
	return nil
}

func (s *Server) listen() {

}

func (s *Server) handle(conn net.Conn) {

}
