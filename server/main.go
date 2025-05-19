package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

var (
	TCPPort      = "50000"
	UDPPort      = "50001"
	clients      = map[string]*Client{}
	clientsMutex sync.Mutex
	// ClientID = "1"
)

type Client struct {
	// ID      string
	tcpConn net.Conn
	udpConn net.UDPConn
	// udpAddr *net.UDPAddr
}

func main() {
	fmt.Println("Welcome to the Game Server!")
	go startTCP()
	go startUDP()
	select {}
}

func startTCP() {
	listener, err := net.Listen("tcp", ":"+TCPPort)
	if err != nil {
		log.Fatalln("Error:", err)
	}
	defer listener.Close()

	fmt.Println("Listening TCP on port", TCPPort)

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

	clientsMutex.Lock()
	_, exists := clients[conn.RemoteAddr().String()]
	if !exists {
		clients[conn.RemoteAddr().String()] = &Client{tcpConn: conn}
	} else {
		clients[conn.RemoteAddr().String()].tcpConn = conn
	}
	// client := clients[conn.RemoteAddr().String()]
	clientsMutex.Unlock()

	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				log.Println("Client Disconnected")
			} else {
				log.Println("Error:", err)
			}

			clientsMutex.Lock()
			delete(clients, conn.RemoteAddr().String())
			clientsMutex.Unlock()
			break
		}
		log.Println("Recieved:", string(buffer[:n]))
	}
}

func startUDP() {
	addr, err := net.ResolveUDPAddr("udp", ":"+UDPPort)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	fmt.Println("Listening UDP on port", UDPPort)

	go handleUDPConnection(conn)
}

func handleUDPConnection(conn *net.UDPConn) {
	defer conn.Close()

	clientsMutex.Lock()
	_, exists := clients[conn.RemoteAddr().String()]
	if !exists {
		clients[conn.RemoteAddr().String()] = &Client{udpConn: *conn}
	} else {
		clients[conn.RemoteAddr().String()].udpConn = *conn
	}
	// client := clients[conn.RemoteAddr().String()]
	clientsMutex.Unlock()

	buffer := make([]byte, 1024)

	for {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Println("Error:", err)
		}

		log.Println("Recieved:", string(buffer[:n]))
	}
}
