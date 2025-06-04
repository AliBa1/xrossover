package game

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	protocol "xrossover-client/flatbuffers/xrossover"

	flatbuffers "github.com/google/flatbuffers/go"
)

const (
	ServerHost    = "localhost"
	ServerTCPPort = "50000"
	ServerTCPAddr = ServerHost + ":" + ServerTCPPort
	ServerUDPPort = "50001"
	ServerUDPAddr = ServerHost + ":" + ServerUDPPort
)

type Network struct {
	tcpConn      net.Conn
	udpConn      net.Conn
	localUDPAddr *net.UDPAddr
	bytesRead    int
	objRegistry  *ObjectRegistry
}

func NewNetwork(udpHost string, udpPort string, objRegistry *ObjectRegistry) *Network {
	localUDPAddr, err := net.ResolveUDPAddr("udp", udpHost+":"+udpPort)
	if err != nil {
		log.Println("Failed to resolve local UDP address:", err)
	}

	return &Network{localUDPAddr: localUDPAddr, objRegistry: objRegistry}
}

func (n *Network) IsConnected() bool {
	if n.tcpConn == nil || n.udpConn == nil {
		return false
	}

	return true
}

func (n *Network) ConnectTCP(username string, objects []GameObject) error {
	log.Println("Connecting to the server via TCP...")

	var d net.Dialer
	var err error
	n.tcpConn, err = d.Dial("tcp", ServerTCPAddr)
	if err != nil {
		return err
	}

	err = n.WriteTCP(n.serializeConnectionRequest(username))
	if err != nil {
		return err
	}

	for _, obj := range objects {
		err = n.WriteTCP(obj.Serialize())
		if err != nil {
			return err
		}
	}

	go n.handleTCP()
	return nil
}

func (n *Network) ConnectUDP() error {
	log.Println("Connecting to the server via UDP...")

	serverAddr, err := net.ResolveUDPAddr("udp", ServerUDPAddr)
	if err != nil {
		return err
	}

	n.udpConn, err = net.DialUDP("udp", n.localUDPAddr, serverAddr)
	if err != nil {
		return err
	}

	go n.handleUDP()
	return nil
}

func (n *Network) Disconnect() {
	if n.tcpConn != nil {
		n.tcpConn.Close()
		n.tcpConn = nil
	}

	if n.udpConn != nil {
		n.udpConn.Close()
		n.udpConn = nil
	}

	log.Println("Byte read by client in session:", n.bytesRead)
}

func (n *Network) handleTCP() {
	for n.tcpConn != nil {
		lengthPrefix := make([]byte, 4)
		bytes, err := n.tcpConn.Read(lengthPrefix)
		if err != nil {
			log.Println("Failed to read message length:", err)
			break
		}
		n.bytesRead += bytes

		dataLen := binary.BigEndian.Uint32(lengthPrefix)
		if dataLen > 10_000 {
			log.Println("Client: Message too large")
			break
		}

		data := make([]byte, dataLen)
		bytes, err = n.tcpConn.Read(data)
		if err != nil {
			switch err {
			case io.EOF:
				log.Println("Client Disconnected")
			default:
				log.Println("Error reading data on client:", err)
			}
			break
		}
		n.bytesRead += bytes

		n.readData(n.tcpConn, data, int(dataLen))
	}
}

func (n *Network) handleUDP() {
	for n.udpConn != nil {
		data := make([]byte, 4096)
		bytes, err := n.udpConn.Read(data)
		if err != nil {
			switch err {
			case io.EOF:
				log.Println("Client Disconnected")
			default:
				log.Println("Error reading data on client:", err)
			}
			break
		}
		n.bytesRead += bytes

		n.readData(n.udpConn, data, bytes)
	}
}

func (n *Network) readData(conn net.Conn, data []byte, bytes int) {
	msg := protocol.GetRootAsNetworkMessage(data[:bytes], 0)
	switch msg.PayloadType() {
	case protocol.PayloadObjectRegistry:
		table := new(flatbuffers.Table)
		if msg.Payload(table) {
			objectRegistry := new(protocol.ObjectRegistry)
			objectRegistry.Init(table.Bytes, table.Pos)

			for i := 0; i < objectRegistry.ObjectsLength(); i++ {
				var objectWrapper protocol.GameObjectWrapper
				if objectRegistry.Objects(&objectWrapper, i) {
					objectUnionTable := new(flatbuffers.Table)
					if objectWrapper.Object(objectUnionTable) {
						switch objectWrapper.ObjectType() {
						case protocol.GameObjectUnionPlayerBox:
							fbPosition := new(protocol.Vector3)
							fbBox := new(protocol.PlayerBox)
							fbBox.Init(objectUnionTable.Bytes, objectUnionTable.Pos)

							id := string(fbBox.Id())
							owner := string(fbBox.Owner())
							position := fbBox.Position(fbPosition)
							playerBox := NewFBPlayerBox(id, owner, *position)
							n.objRegistry.Add(playerBox)
						case protocol.GameObjectUnionBall:
							fbPosition := new(protocol.Vector3)
							fbBall := new(protocol.Ball)
							fbBall.Init(objectUnionTable.Bytes, objectUnionTable.Pos)

							id := string(fbBall.Id())
							owner := string(fbBall.Owner())
							position := fbBall.Position(fbPosition)
							ball := NewFBBall(id, owner, *position)
							n.objRegistry.Add(ball)
						}
					}
				}
			}
		}
	case protocol.PayloadPlayerBox:
		table := new(flatbuffers.Table)
		if msg.Payload(table) {
			fbPosition := new(protocol.Vector3)
			fbBox := new(protocol.PlayerBox)
			fbBox.Init(table.Bytes, table.Pos)
			id := string(fbBox.Id())
			position := fbBox.Position(fbPosition)
			obj, err := n.objRegistry.Get(id)
			if err != nil {
				break
			}
			obj.UpdatePosition(position.X(), position.Y(), position.Z())
		}
	case protocol.PayloadBall:
		table := new(flatbuffers.Table)
		if msg.Payload(table) {
			fbPosition := new(protocol.Vector3)
			fbBall := new(protocol.Ball)
			fbBall.Init(table.Bytes, table.Pos)

			id := string(fbBall.Id())
			position := fbBall.Position(fbPosition)
			obj, err := n.objRegistry.Get(id)
			if err != nil {
				break
			}
			obj.UpdatePosition(position.X(), position.Y(), position.Z())
		}
	default:
		log.Println("Received without type:", msg.PayloadType())
	}
}

func (n *Network) WriteTCP(data []byte) error {
	length := uint32(len(data))
	var lengthPrefix [4]byte
	binary.BigEndian.PutUint32(lengthPrefix[:], length)

	_, err := n.tcpConn.Write(lengthPrefix[:])
	if err != nil {
		return err
	}

	_, err = n.tcpConn.Write(data)
	return err
}

func (n *Network) WriteUDP(data []byte) error {
	_, err := n.udpConn.Write(data)
	return err
}

func (n *Network) serializeConnectionRequest(username string) []byte {
	builder := flatbuffers.NewBuilder(1024)

	user := builder.CreateString(username)
	udp := builder.CreateString(n.localUDPAddr.String())

	protocol.ConnectionRequestStart(builder)
	protocol.ConnectionRequestAddUsername(builder, user)
	protocol.ConnectionRequestAddUdpaddr(builder, udp)
	connReq := protocol.ConnectionRequestEnd(builder)

	protocol.NetworkMessageStart(builder)
	protocol.NetworkMessageAddPayloadType(builder, protocol.PayloadConnectionRequest)
	protocol.NetworkMessageAddPayload(builder, connReq)
	netMsg := protocol.NetworkMessageEnd(builder)

	builder.Finish(netMsg)

	return builder.FinishedBytes()
}
