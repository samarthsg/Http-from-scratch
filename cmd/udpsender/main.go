package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":42069")
	if err != nil {
		fmt.Print("Error:", err)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Fatal error: ", err)
	}
	defer conn.Close()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		_, err = conn.Write([]byte(line))
		if err != nil {
			log.Println(err)
		}
	}
}
