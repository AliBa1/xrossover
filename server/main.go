package main

import (
	"fmt"
	"log"
	"net"
)

var (
	PORT = "50000"
)

func main() {
	fmt.Println("Welcome to the Game Server!")
	// NOTE: switch to tcp6 if ipv6
	listener, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalln("Error:", err)
	}
	defer listener.Close()

	fmt.Println("Listening on port", PORT)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("Error:", err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
}
