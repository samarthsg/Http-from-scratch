package main

import (
	"fmt"
	"http-from-scratch/internal/request"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Print("Error:", err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Print("Error:", err)
		}
		defer conn.Close()
		fmt.Println("A connection has been accepted")
		req, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Print("ERROR: ", err)
			break
		}
		fmt.Println("Request line:")
		fmt.Println("-> Method: ", req.RequestLine.Method)
		fmt.Println("-> Target: ", req.RequestLine.RequestTarget)
		fmt.Println("-> HTTP version: ", req.RequestLine.HttpVersion)
		fmt.Println()
	}
}
