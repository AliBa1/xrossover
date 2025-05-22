package game

import (
	"log"
	"net"
	"xrossover-client/internal/buffer"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	WIDTH   = 800
	HEIGHT  = 450
	HOST    = "localhost"
	TCPPORT = "50000"
	UDPPORT = "50001"
)

type Game struct {
	Username string
	camera   rl.Camera3D
	box      *PlayerBox
	// cubePosition rl.Vector3
	tcpConn net.Conn
	udpConn net.Conn
}

func (g *Game) Run() {
	g.initialize()
	g.loop()
	g.shutdown()
}

func (g *Game) initialize() {
	rl.InitWindow(WIDTH, HEIGHT, "Game Window")
	g.camera = rl.Camera3D{
		Position:   rl.Vector3{X: 0.0, Y: 10.0, Z: 10.0},
		Target:     rl.Vector3{X: 0.0, Y: 0.0, Z: 0.0},
		Up:         rl.Vector3{X: 0.0, Y: 1.0, Z: 0.0},
		Fovy:       45.0,
		Projection: rl.CameraPerspective,
	}
	// g.cubePosition = rl.Vector3{X: 0.0, Y: 1.0, Z: 0.0}
	g.box = NewPlayerBox(g.Username)

	rl.SetTargetFPS(60)
}

func (g *Game) loop() {
	for !rl.WindowShouldClose() {
		g.processInput()
		rl.BeginDrawing()
		g.updateDrawing()
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

	if g.tcpConn != nil && g.udpConn != nil {
		rl.DrawText("Connected to Server", WIDTH-rl.MeasureText("Connected to Server", 20)-10, 20, 20, rl.Black)
	} else {
		rl.DrawText("Press [C] to connect to server", WIDTH-rl.MeasureText("Press [C] to connect to server", 20)-10, 20, 20, rl.Black)
	}

	rl.DrawFPS(20, 20)
}

func (g *Game) update3DOutput() {
	// rl.DrawCube(g.cubePosition, 1.0, 1.0, 1.0, rl.Red)
	rl.DrawCube(g.box.position, g.box.width, g.box.height, g.box.length, g.box.color)
	rl.DrawCubeWires(g.box.position, 1.0, 1.0, 1.0, rl.Maroon)
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
		data, err := g.box.SerializeMove(x, y, z)
		if err != nil {
			log.Println("Error serializing movement data:", err)
		}
		g.sendUDP(data)
	}
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

	data := buffer.SerializeConnectionRequest(g.Username, udpAddr)
	_, err = g.tcpConn.Write(data)
	if err != nil {
		log.Println("Error writing to server via TCP:", err)
		return
	}

	data, err = g.box.Serialize()
	_, err = g.tcpConn.Write(data)
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

	g.sendUDP([]byte("This is " + g.Username + " via UDP"))
	// _, err = g.udpConn.Write([]byte("This is " + g.Username + " via UDP"))
	// if err != nil {
	// 	log.Println("Error writing to server via UDP:", err)
	// 	return
	// }
}

func (g *Game) sendUDP(data []byte) {
	_, err := g.udpConn.Write(data)
	if err != nil {
		log.Println("Error writing to server via UDP:", err)
		return
	}
}
