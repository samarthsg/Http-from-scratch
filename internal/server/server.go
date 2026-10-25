package server

import (
	"net"
)

const port = "42069"

type serverState int

const (
	stateIntilaized serverState = iota
	stateDone
)

type Server struct {
	State serverState
}

func Serve(port int) (*Server, error) {
	server := &Server{}
	listner, err := net.Listen("tcp", "42069")
	if err != nil {
		return nil, err
	}
	defer listner.Close()
	for {
		conn, err := listner.Accept()
		if err != nil {
			return nil, err
		}
		defer conn.Close()
	}
	return server, nil
}

func (s *Server) Close() error {
	return nil
}

func (s *Server) listen() {

}

func (s *Server) handle(conn net.Conn) {

}
