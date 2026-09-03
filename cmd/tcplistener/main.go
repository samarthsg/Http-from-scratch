package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	currLine := ""
	go func() {
		defer f.Close()
		defer close(ch)
		for {
			data := make([]byte, 8)
			n, err := f.Read(data)
			if err != nil {
				break
			}
			parts := strings.Split(string(data[:n]), "\n")
			for i := 0; i < len(parts)-1; i++ {
				currLine += parts[i]
				ch <- currLine
				currLine = ""
			}
			currLine += parts[len(parts)-1]
		}
		if currLine != "" {
			ch <- currLine
		}
	}()
	return ch
}
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
		fmt.Println("A connection has been accepted")
		for line := range getLinesChannel(conn) {
			fmt.Println(line)
		}
		fmt.Println("Connection has been closed")
	}
}
