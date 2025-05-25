package game

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	protocol "xrossover-client/flatbuffers/xrossover"

	rl "github.com/gen2brain/raylib-go/raylib"
	flatbuffers "github.com/google/flatbuffers/go"
)

const (
	WIDTH   = 800
	HEIGHT  = 450
	HOST    = "localhost"
	TCPPORT = "50000"
	UDPPORT = "50001"
)

type Game struct {
	Username       string
	camera         rl.Camera3D
	box            *PlayerBox
	tcpConn        net.Conn
	udpConn        net.Conn
	objectRegistry *ObjectRegistry
}

func (g *Game) Run() {
	g.initialize()
	g.loop()
	g.shutdown()
}

func (g *Game) initialize() {
	g.objectRegistry = NewObjectRegistry()
	rl.InitWindow(WIDTH, HEIGHT, "Game Window")
	g.camera = rl.Camera3D{
		Position:   rl.Vector3{X: 0.0, Y: 10.0, Z: 10.0},
		Target:     rl.Vector3{X: 0.0, Y: 0.0, Z: 0.0},
		Up:         rl.Vector3{X: 0.0, Y: 1.0, Z: 0.0},
		Fovy:       45.0,
		Projection: rl.CameraPerspective,
	}
	g.box = NewPlayerBox(g.Username)
	g.objectRegistry.Add(g.box)

	rl.SetTargetFPS(60)
}

func (g *Game) loop() {
	for !rl.WindowShouldClose() {
		g.processInput()
		rl.BeginDrawing()
		g.updateDrawing()
		if g.tcpConn != nil {
			go g.handleConnection(g.tcpConn)
		}
		if g.udpConn != nil {
			go g.handleConnection(g.udpConn)
		}
		rl.EndDrawing()
	}

}

func (g *Game) shutdown() {
	rl.CloseWindow()

	if g.tcpConn != nil {
		g.tcpConn.Close()
	}

	if g.udpConn != nil {
		g.udpConn.Close()
	}
}

func (g *Game) updateDrawing() {
	rl.ClearBackground(rl.RayWhite)

	rl.BeginMode3D(g.camera)
	g.update3DOutput()
	rl.EndMode3D()

	for _, obj := range g.objectRegistry.Objects {
		objScreenPosition := rl.GetWorldToScreen(rl.Vector3{X: obj.Position().X, Y: obj.Position().Y + 1.5, Z: obj.Position().Z}, g.camera)
		rl.DrawText(obj.ID(), int32(objScreenPosition.X)-rl.MeasureText(obj.ID(), 20)/2, int32(objScreenPosition.Y), 20, rl.Black)
	}

	if g.tcpConn != nil && g.udpConn != nil {
		rl.DrawText("Connected to Server", WIDTH-rl.MeasureText("Connected to Server", 20)-10, 20, 20, rl.Black)
	} else {
		rl.DrawText("Press [C] to connect to server", WIDTH-rl.MeasureText("Press [C] to connect to server", 20)-10, 20, 20, rl.Black)
	}

	rl.DrawFPS(20, 20)
}

func (g *Game) update3DOutput() {
	for _, obj := range g.objectRegistry.Objects {
		rl.DrawCube(obj.Position(), obj.Dimensions().Width, obj.Dimensions().Height, obj.Dimensions().Length, obj.Color())
		rl.DrawCubeWires(obj.Position(), 1.0, 1.0, 1.0, rl.Maroon)
	}
	rl.DrawGrid(10, 1.0)
}

func (g *Game) processInput() {
	// move cube
	if rl.IsKeyDown(rl.KeyW) {
		// g.cubePosition.Z -= 0.05
		g.box.Move(0.0, 0.0, -0.05)
		g.sendMovement(0.0, 0.0, -0.05)
	} else if rl.IsKeyDown(rl.KeyS) {
		// g.cubePosition.Z += 0.05
		g.box.Move(0.0, 0.0, 0.05)
		g.sendMovement(0.0, 0.0, 0.05)
	} else if rl.IsKeyDown(rl.KeyA) {
		// g.cubePosition.X -= 0.05
		g.box.Move(-0.05, 0.0, 0.0)
		g.sendMovement(-0.05, 0.0, 0.0)
	} else if rl.IsKeyDown(rl.KeyD) {
		// g.cubePosition.X += 0.05
		g.box.Move(0.05, 0.0, 0.0)
		g.sendMovement(0.05, 0.0, 0.0)
	}

	// connect to server
	if rl.IsKeyPressed(rl.KeyC) && g.tcpConn == nil {
		g.connectToTCP()
		g.connectToUDP()
	}
}

func (g *Game) sendMovement(x float32, y float32, z float32) {
	if g.udpConn != nil {
		g.sendMessage("udp", g.box.SerializeMove(x, y, z))
	}
}

func (g *Game) serializeConnectionRequest(udpAddr *net.UDPAddr) []byte {
	builder := flatbuffers.NewBuilder(1024)

	user := builder.CreateString(g.Username)
	udp := builder.CreateString(udpAddr.String())

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

func (g *Game) connectToTCP() {
	log.Println("Starting TCP connection to server...")
	var d net.Dialer
	var err error
	g.tcpConn, err = d.Dial("tcp", HOST+":"+TCPPORT)
	if err != nil {
		log.Println("Error connecting to server via TCP:", err)
		return
	}

	udpAddr, err := net.ResolveUDPAddr("udp", HOST+":"+UDPPORT)
	if err != nil {
		log.Println("Failed to get UDP address:", err)
	}

	// err = g.sendMessage("tcp", serialize.ConnectionRequest(g.Username, udpAddr))
	err = g.sendMessage("tcp", g.serializeConnectionRequest(udpAddr))
	if err != nil {
		log.Println("Error sending player box:", err)
		return
	}

	err = g.sendMessage("tcp", g.box.Serialize())
	if err != nil {
		log.Println("Error sending player box:", err)
		return
	}
}

func (g *Game) connectToUDP() {
	log.Println("Starting UDP connection to server...")
	var d net.Dialer
	var err error
	g.udpConn, err = d.Dial("udp", HOST+":"+UDPPORT)
	if err != nil {
		log.Println("Error connecting to server via UDP:", err)
		return
	}
}

func (g *Game) sendMessage(protocol string, data []byte) error {
	if protocol != "tcp" && protocol != "udp" {
		return errors.New("attempted to send message with an unsupported protocol")
	}

	length := uint32(len(data))
	var lengthPrefix [4]byte
	binary.BigEndian.PutUint32(lengthPrefix[:], length)

	var conn net.Conn
	if protocol == "tcp" {
		conn = g.tcpConn
	} else {
		conn = g.tcpConn
	}

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

func (g *Game) handleConnection(conn net.Conn) {
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
			log.Println("Message too large")
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

		g.readData(conn, data, int(dataLen))
	}
}

func (g *Game) readData(conn net.Conn, data []byte, n int) {
	msg := protocol.GetRootAsNetworkMessage(data[:n], 0)
	switch msg.PayloadType() {
	case protocol.PayloadObjectRegistry:
		log.Println("recieved object registry from server")
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
							position := fbBox.Position(fbPosition)
							playerBox := NewFBPlayerBox(id, *position)
							g.objectRegistry.Add(playerBox)
							log.Println(g.Username, "Added box with ID to client registry:", id)
						}
					}
				}
			}
		}
	default:
		log.Println("Received without type:", msg.PayloadType())
	}
}
