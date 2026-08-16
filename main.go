package main

import (
	"fmt"
	"io"
	"os"
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
	f, err := os.Open("messages.txt")
	if err != nil {
		fmt.Print("File dont exist")
	}
	lines := getLinesChannel(f)
	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}
}
