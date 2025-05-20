package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	schema "xrossover-server/flatbuffers/xrossover"
)

const (
	TCPPort = "50000"
	UDPPort = "50001"
)

var (
	clients      = map[string]*Client{}
	clientsMutex sync.Mutex
	// ClientID = "1"
)

type Client struct {
	ID      string
	tcpConn net.Conn
	// udpConn net.UDPConn
	udpAddr *net.UDPAddr
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
		log.Println("Error:", err)
	}
	defer listener.Close()

	fmt.Println("Listening TCP on port", TCPPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error:", err)
			continue
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
			} else {
				log.Println("Error erere:", err)
			}

			// clientsMutex.Lock()
			// delete(clients, c.ID)
			// clientsMutex.Unlock()
			break
		}

		// msg := strings.TrimSpace(string(buffer[:n]))
		connReq := schema.GetRootAsConnectionRequest(buffer[:n], 0)
		username := string(connReq.Username())
		if len(username) > 0 {
			log.Println("adding to clients")
			addClient(username, conn)
		} else {
			// log.Println("Received:", msg)
			log.Println("Received without username")
		}

	}
}

func addClient(id string, conn net.Conn) {
	clientsMutex.Lock()
	client, exists := clients[id]
	if !exists {
		client = &Client{
			ID:      id,
			tcpConn: conn,
		}
		clients[id] = client
	} else {
		clients[id].tcpConn = conn
	}
	clientsMutex.Unlock()
	log.Println(id, "added to clients")
}

// func readTCP(c *Client) {
// 	defer c.tcpConn.Close()
//
// 	buffer := make([]byte, 1024)
//
// 	for {
// 		n, err := c.tcpConn.Read(buffer)
// 		if err != nil {
// 			if err == io.EOF {
// 				log.Println("Client Disconnected")
// 			} else {
// 				log.Println("Error erere:", err)
// 			}
//
// 			clientsMutex.Lock()
// 			delete(clients, c.ID)
// 			clientsMutex.Unlock()
// 			break
// 		}
// 		log.Println("Recieved:", string(buffer[:n]))
// 		if string(buffer[:n])[0] == 'A' {
// 			username := string(buffer[:n])
// 			_, err := c.tcpConn.Read(buffer)
// 			if err != nil {
// 				log.Println("Failed reading username:", err)
// 				return
// 			}
// 			addClient(username, c.tcpConn)
//
// 		}
// 	}
//
// }

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

	go handleUDPConnection(conn, addr)
}

func handleUDPConnection(conn *net.UDPConn, addr *net.UDPAddr) {
	defer conn.Close()

	// id := conn.RemoteAddr().String()
	username := make([]byte, 1024)
	n, err := conn.Read(username)
	if err != nil {
		log.Println("Failed reading username:", err)
		return
	}

	id := string(username[:n])

	clientsMutex.Lock()
	client, exists := clients[id]
	if !exists {
		client = &Client{
			ID:      id,
			udpAddr: addr,
		}
		clients[id] = client
	} else {
		clients[id].udpAddr = addr
	}
	clientsMutex.Unlock()

	buffer := make([]byte, 1024)

	for {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Println("Error:", err)
		}

		log.Println("Recieved:", string(buffer[:n]))
		// conn.WriteToUDP([]byte("From UDP server!\n"), addr)
	}
}
