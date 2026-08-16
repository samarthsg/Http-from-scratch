package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println("Error opening the file")
	}
	currLine := ""
	for {
		data := make([]byte, 8)
		n, err := f.Read(data)
		if err != nil {
			break
		}
		parts := strings.Split(string(data[:n]), "\n")
		for i := 0; i < len(parts)-1; i++ {
			currLine += parts[i]
			fmt.Printf("read: %s\n", currLine)
			currLine = ""
		}
		currLine += parts[len(parts)-1]
	}
	if currLine != "" {
		fmt.Printf("read: %s\n", currLine)
	}
	f.Close()
}
