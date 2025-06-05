package main

import (
	"encoding/binary"
	"errors"
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
	g            = game.Game{}
	localUDPConn *net.UDPConn
)

type Client struct {
	Username string
	tcpConn  net.Conn
	udpAddr  *net.UDPAddr
}

func main() {
	g.Broadcast = broadcast
	go startTCP()
	go startUDP()
	go g.Run()

	select {}
}

func startTCP() {
	listener, err := net.Listen("tcp", ":"+TCPPort)
	if err != nil {
		log.Fatalln("Error:", err)
	}
	defer listener.Close()

	log.Println("Listening TCP on port", TCPPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error:", err)
			continue
		}
		go handleTCP(conn)
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
	localUDPConn = conn

	log.Println("Listening UDP on port", UDPPort)

	go handleUDP()
}

func handleTCP(conn net.Conn) {
	for {
		lengthPrefix := make([]byte, 4)
		_, err := conn.Read(lengthPrefix)
		if err != nil {
			log.Println("Failed to read message length:", err)
			break
		}

		dataLen := binary.BigEndian.Uint32(lengthPrefix)
		if dataLen > 10_000 {
			log.Println("Client: Message too large")
			break
		}

		data := make([]byte, dataLen)
		_, err = conn.Read(data)
		if err != nil {
			if err == io.EOF {
				log.Println("Client Disconnected")
			} else {
				log.Println("Error reading data on server:", err)
			}
			break
		}

		readData(conn, data, int(dataLen))

	}
}

func handleUDP() {
	for {
		data := make([]byte, 4096)
		n, err := localUDPConn.Read(data)
		if err != nil {
			if err == io.EOF {
				log.Println("Client Disconnected")
			} else {
				log.Println("Error reading data on client:", err)
			}
			break
		}

		readData(localUDPConn, data, n)
	}
}

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
				err := writeTCP(c.tcpConn, data)
				if err != nil {
					log.Println(err)
				}
			case "udp":
				err := writeUDP(c.udpAddr, data)
				if err != nil {
					log.Println(err)
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
		table := new(flatbuffers.Table)
		if msg.Payload(table) {
			connReq := new(protocol.ConnectionRequest)
			connReq.Init(table.Bytes, table.Pos)
			username := string(connReq.Username())
			udpStr := string(connReq.Udpaddr())
			addClient(username, conn, udpStr)
		}
	case protocol.PayloadPlayerBox:
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

			broadcast("tcp", "", g.ObjectRegistry.Serialize())
		}
	case protocol.PayloadBall:
		table := new(flatbuffers.Table)
		if msg.Payload(table) {
			fbPosition := new(protocol.Vector3)
			fbBall := new(protocol.Ball)
			fbBall.Init(table.Bytes, table.Pos)

			id := string(fbBall.Id())
			owner := string(fbBall.Owner())
			position := fbBall.Position(fbPosition)
			ball, err := g.ObjectRegistry.Get(id)
			// if ball not in registry
			if err != nil {
				ball = game.NewBall(id, owner, *position)
				g.ObjectRegistry.Add(ball)
				broadcast("tcp", "", g.ObjectRegistry.Serialize())
				return
			}

			ball.UpdatePosition(position.X(), position.Y(), position.Z())
			broadcast("udp", owner, ball.Serialize())
		}
	case protocol.PayloadPlayerInput:
		table := new(flatbuffers.Table)
		if msg.Payload(table) {
			fbInput := new(protocol.PlayerInput)
			fbInput.Init(table.Bytes, table.Pos)
			input := deserializeInput(fbInput)
			g.AddPlayerInput(input)
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
}

func writeTCP(conn net.Conn, data []byte) error {
	length := uint32(len(data))
	var lengthPrefix [4]byte
	binary.BigEndian.PutUint32(lengthPrefix[:], length)

	_, err := conn.Write(lengthPrefix[:])
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func writeUDP(remoteAddr *net.UDPAddr, data []byte) error {
	if localUDPConn == nil {
		return errors.New("udp connection not initialized on the server")
	}

	_, err := localUDPConn.WriteToUDP(data, remoteAddr)
	if err != nil {
		return err
	}

	return nil
}
