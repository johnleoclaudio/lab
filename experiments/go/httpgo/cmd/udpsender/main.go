package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Println(err)
		return
	}

	udpConn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Println(err)
		return
	}

	defer udpConn.Close()

	ioReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		str, err := ioReader.ReadString('\n')
		if err != nil {
			log.Println(err)
		}

		udpConn.Write([]byte(str))
	}

}
