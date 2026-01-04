package main

import (
	"fmt"
	"io"
	"log"
	"os"
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
				log.Println(err)
			}

			if isEOF == 0 {
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
	file, err := os.Open("message.txt")
	if err != nil {
		log.Println(err)
		return
	}

	for s := range getLinesChannel(file) {
		fmt.Println("read:", s)
	}
}
