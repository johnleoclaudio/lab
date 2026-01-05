package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	strChan := make(chan string)

	var a = make([]byte, 8)
	var line string

	go func() {
		for {
			isEOF, err := f.Read(a)
			if err != nil {
				fmt.Println(err)
			}

			if isEOF == 0 {
				strChan <- "EOF"
				close(strChan)
				return
			}

			stringedA := string(a)

			if strings.Contains(stringedA, "\n") {
				str := strings.Split(stringedA, "\n")
				line = line + strings.Join(str[:1], "")
				strChan <- line

				line = ""
				line = line + strings.Join(str[1:], "")
				continue
			}

			line = line + string(a)
		}
	}()

	return strChan
}

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Println(err)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}

		fmt.Println("connection accepted!")

		for s := range getLinesChannel(conn) {
			if s == "EOF" {
				fmt.Println("connection closed!")
			}
			fmt.Println(s)
		}
	}
}
