package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	protocol "xrossover-server/flatbuffers/xrossover"
	"xrossover-server/internal/game"

	rl "github.com/gen2brain/raylib-go/raylib"
	flatbuffers "github.com/google/flatbuffers/go"
)

const (
	TCPPort = "50000"
	UDPPort = "50001"
)

var (
	clients      = map[string]*Client{}
	clientsMutex sync.Mutex
	// objectRegistry = game.NewObjectRegistry()
	g = game.Game{}
)

type Client struct {
	Username string
	tcpConn  net.Conn
	udpAddr  *net.UDPAddr
}

func main() {
	fmt.Println("Welcome to the Game Server!")
	g.Broadcast = broadcast
	go startTCP()
	go startUDP()
	go g.Run()

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
		go handleConnection(conn)
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

	go handleConnection(conn)
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		lengthPrefix := make([]byte, 4)
		_, err := conn.Read(lengthPrefix)
		if err != nil {
			log.Println("Failed to read message length:", err)
			break
		}
		dataLen := binary.BigEndian.Uint32(lengthPrefix)
		if dataLen > 10_000 {
			log.Println("Server: Message too large")
			break
		}

		data := make([]byte, dataLen)
		_, err = conn.Read(data)
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

		readData(conn, data, int(dataLen))
	}
}

// TODO: add a username as a param to send to all except a client since they already applied their changes through prediction
func broadcast(protocol, owner string, data []byte) {
	if protocol != "tcp" && protocol != "udp" {
		log.Println("Must broadcast message using TCP or UDP")
		return
	}
	clientsMutex.Lock()
	for _, c := range clients {
		if c.Username != owner {
			switch protocol {
			case "tcp":
				sendMessage(c.tcpConn, data)
			case "udp":
				conn, err := net.DialUDP("udp", nil, c.udpAddr)
				if err != nil {
					log.Println("Error broadcasting UDP message:", err)
				} else {
					sendMessage(conn, data)
				}

			}
		}
	}
	clientsMutex.Unlock()
}

func readData(conn net.Conn, data []byte, n int) {
	msg := protocol.GetRootAsNetworkMessage(data[:n], 0)
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
			owner := string(fbBox.Owner())
			position := fbBox.Position(fbPosition)
			playerBox := game.NewPlayerBox(id, owner, *position)
			g.ObjectRegistry.Add(playerBox)

			log.Println("Current object registry:", g.ObjectRegistry.Objects)
			broadcast("tcp", "", g.ObjectRegistry.Serialize())
		}
	// case protocol.PayloadMovement:
	// 	// log.Println("recieved movement data")
	// 	table := new(flatbuffers.Table)
	// 	if msg.Payload(table) {
	// 		fbDirection := new(protocol.Vector3)
	// 		fbMovement := new(protocol.Movement)
	// 		fbMovement.Init(table.Bytes, table.Pos)
	// 		id := string(fbMovement.ObjectId())
	// 		owner := string(fbMovement.ObjectOwner())
	// 		direction := fbMovement.Direction(fbDirection)
	// 		obj, err := g.ObjectRegistry.Get(id)
	// 		if err != nil {
	// 			log.Println(err)
	// 			return
	// 		}
	// 		obj.Move(direction.X(), direction.Y(), direction.Z())
	// 		// TODO: change to udp and fix for udp
	// 		broadcast("tcp", owner, obj.Serialize())
	// 	}
	case protocol.PayloadPlayerInput:
		log.Println("recieved movement data")
		table := new(flatbuffers.Table)
		if msg.Payload(table) {
			fbInput := new(protocol.PlayerInput)
			fbInput.Init(table.Bytes, table.Pos)
			input := deserializeInput(fbInput)
			g.AddPlayerInput(input)
			// broadcast("tcp", owner, obj.Serialize())
		}
	default:
		log.Println("Received without type:", msg.PayloadType())
	}

}

func deserializeInput(fbInput *protocol.PlayerInput) game.PlayerInput {
	i := game.PlayerInput{
		ObjectID: string(fbInput.ObjectId()),
	}

	switch fbInput.ActionType() {
	case protocol.ActionMove:
		table := new(flatbuffers.Table)
		if fbInput.Action(table) {
			fbMove := new(protocol.Move)
			fbMove.Init(table.Bytes, table.Pos)

			fbDir := new(protocol.Vector3)
			dir := fbMove.Direction(fbDir)

			i.Action = game.Move{
				Direction: rl.Vector3{
					X: dir.X(),
					Y: dir.Y(),
					Z: dir.Z(),
				},
			}
		}
	default:
		i.Action = nil
	}

	return i
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

func sendMessage(conn net.Conn, data []byte) error {
	length := uint32(len(data))
	var lengthPrefix [4]byte
	binary.BigEndian.PutUint32(lengthPrefix[:], length)

	_, err := conn.Write(lengthPrefix[:])
	if err != nil {
		return errors.New("error sending buffer length prefix to server")
	}

	_, err = conn.Write(data)
	if err != nil {
		return errors.New("error sending data to server")
	}

	return nil
}
