package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

var (
	TCPPORT = "50000"
	UDPPORT = "50001"
)

type Client struct {
	tcpConn net.Conn
	udpConn net.UDPConn
}

func main() {
	// var clients []Client
	fmt.Println("Welcome to the Game Server!")
	go startTCP()
	go startUDP()
	select {}
}

func startTCP() {
	listener, err := net.Listen("tcp", ":"+TCPPORT)
	if err != nil {
		log.Fatalln("Error:", err)
	}
	defer listener.Close()

	fmt.Println("Listening TCP on port", TCPPORT)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("Error:", err)
		}
		go handleTCPConnection(conn)
	}
}

func handleTCPConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				log.Println("Client Disconnected")
				break
			}
			log.Println("Error:", err)
			break
		}
		log.Println("Recieved:", string(buffer[:n]))
	}
}

func startUDP() {
	addr, err := net.ResolveUDPAddr("udp", ":"+UDPPORT)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	fmt.Println("Listening UDP on port", UDPPORT)

	go handleUDPConnection(conn)
}

func handleUDPConnection(conn *net.UDPConn) {
	defer conn.Close()

	buffer := make([]byte, 1024)

	for {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Println("Error:", err)
		}

		log.Println("Recieved:", string(buffer[:n]))
	}
}
