package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	protocol "xrossover-server/flatbuffers/xrossover"
	"xrossover-server/internal/game"

	flatbuffers "github.com/google/flatbuffers/go"
)

const (
	TCPPort = "50000"
	UDPPort = "50001"
)

var (
	clients        = map[string]*Client{}
	clientsMutex   sync.Mutex
	objectRegistry = game.NewObjectRegistry()
)

type Client struct {
	Username string
	tcpConn  net.Conn
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

		readTCP(conn, buffer, n)
	}
}

func readTCP(conn net.Conn, buffer []byte, n int) {
	msg := protocol.GetRootAsNetworkMessage(buffer[:n], 0)
	switch msg.PayloadType() {
	case protocol.PayloadConnectionRequest:
		log.Println("adding to clients")
		table := new(flatbuffers.Table)
		if msg.Payload(table) {
			connReq := new(protocol.ConnectionRequest)
			connReq.Init(table.Bytes, table.Pos)
			username := string(connReq.Username())
			udpStr := string(connReq.Udpaddr())
			addClient(username, conn, udpStr)
		}
	case protocol.PayloadPlayerBox:
		log.Println("recieved a player box")
		table := new(flatbuffers.Table)
		if msg.Payload(table) {
			fbPosition := new(protocol.Vector3)
			fbBox := new(protocol.PlayerBox)
			fbBox.Init(table.Bytes, table.Pos)
			// add box to object registry
			id := string(fbBox.Id())
			position := fbBox.Position(fbPosition)
			playerBox := game.NewPlayerBox(id, *position)
			objectRegistry.Add(playerBox)
		}
	default:
		log.Println("Received without type:", msg.PayloadType())
	}
}

func addClient(username string, conn net.Conn, udpStr string) {
	udpAddr, err := net.ResolveUDPAddr("udp", udpStr)
	if err != nil {
		log.Println("Error getting client udp address:", err)
	}
	clientsMutex.Lock()
	client, exists := clients[username]
	if !exists {
		client = &Client{
			Username: username,
			tcpConn:  conn,
			udpAddr:  udpAddr,
		}
		clients[username] = client
	} else {
		clients[username].tcpConn = conn
		clients[username].udpAddr = udpAddr
	}
	clientsMutex.Unlock()
	log.Println(username, "added to clients")
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

	buffer := make([]byte, 1024)

	for {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Println("Error:", err)
		}

		log.Println("Recieved:", string(buffer[:n]))
		readUDP(conn, buffer, n)
		// conn.WriteToUDP([]byte("From UDP server!\n"), addr)
	}
}

func readUDP(conn net.Conn, buffer []byte, n int) {
	msg := protocol.GetRootAsNetworkMessage(buffer[:n], 0)
	switch msg.PayloadType() {
	case protocol.PayloadMovement:
		log.Println("recieved movement data")
		table := new(flatbuffers.Table)
		if msg.Payload(table) {
			fbDirection := new(protocol.Vector3)
			fbMovement := new(protocol.Movement)
			fbMovement.Init(table.Bytes, table.Pos)
			id := string(fbMovement.ObjectId())
			direction := fbMovement.Direction(fbDirection)
			obj, err := objectRegistry.Get(id)
			if err != nil {
				log.Println(err)
				return
			}
			obj.Move(direction.X(), direction.Y(), direction.Z())
		}
	default:
		log.Println("Received without type:", msg.PayloadType())
	}
}
