package game

import (
	"fmt"
	"log"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	WIDTH   = 800
	HEIGHT  = 450
	GroundY = 0.0
	// Gravity = -9.8
	Gravity = -32.17
)

type Game struct {
	Username      string
	camera        rl.Camera3D
	box           *PlayerBox
	ball          *Ball
	hoop          *Hoop
	objRegistry   *ObjectRegistry
	network       *Network
	points        int
	scoreCooldown float32
}

func NewGame(username string, net *Network, objRegistry *ObjectRegistry) *Game {
	return &Game{Username: username, network: net, objRegistry: objRegistry}
}

func (g *Game) Run() {
	g.initialize()
	g.loop()
	defer g.shutdown()
}

func (g *Game) initialize() {
	rl.InitWindow(WIDTH, HEIGHT, g.Username+"'s Game Window")
	g.camera = rl.Camera3D{
		// Position:   rl.Vector3{X: 0, Y: 10, Z: 20},
		Position:   rl.Vector3{X: -5, Y: 15, Z: 20},
		Target:     rl.Vector3{X: 0, Y: 5, Z: 0},
		Up:         rl.Vector3{X: 0, Y: 1, Z: 0},
		Fovy:       45.0,
		Projection: rl.CameraPerspective,
	}

	g.box = NewPlayerBox(g.Username+"-Box", g.Username)
	g.ball = NewBall(g.Username+"-Ball", g.Username, g.box)
	g.hoop = NewHoop(0.0, -5.0)
	g.points = 0

	g.objRegistry.Add(g.box)
	g.objRegistry.Add(g.ball)

	rl.SetTargetFPS(60)
}

func (g *Game) loop() {
	lastTime := time.Now()
	for !rl.WindowShouldClose() {
		currentTime := time.Now()
		dt := float32(currentTime.Sub(lastTime).Seconds())
		lastTime = currentTime

		g.update(dt)

		g.processInput()
		rl.BeginDrawing()
		g.updateDrawing()

		rl.EndDrawing()
	}
}

func (g *Game) shutdown() {
	rl.CloseWindow()

	g.network.Disconnect()
}

func (g *Game) update(dt float32) {
	// rl.UpdateCamera(&g.camera, rl.CameraFree)
	g.ball.Update(dt)
	g.ball.DetectCollision(dt, *g.hoop)

	if g.network.IsConnected() {
		g.network.WriteUDP(g.ball.Serialize())
	}

	if g.scoreCooldown > 0 {
		g.scoreCooldown -= dt
	}

	g.detectScore()
}

func (g *Game) detectScore() {
	ballBelowRim := g.ball.position.Y+g.ball.radius < g.hoop.rim.position.Y
	ballInRimX := g.ball.position.X+g.ball.radius > g.hoop.rim.position.X-g.hoop.rim.radius && g.ball.position.X-g.ball.radius < g.hoop.rim.position.X+g.hoop.rim.radius
	ballInRimZ := g.ball.position.Z-g.ball.radius > g.hoop.rim.position.Z-g.hoop.rim.radius && g.ball.position.Z+g.ball.radius < g.hoop.rim.position.Z+g.hoop.rim.radius
	ballGoingDown := g.ball.velocity.Y < 0

	if g.scoreCooldown <= 0 && ballBelowRim && ballInRimX && ballInRimZ && ballGoingDown {
		fmt.Println("Score!!!!!")
		g.points += 1
		g.scoreCooldown = 1.0
	}
}

func (g *Game) updateDrawing() {
	rl.ClearBackground(rl.RayWhite)

	rl.BeginMode3D(g.camera)
	g.update3DOutput()
	rl.EndMode3D()

	for _, obj := range g.objRegistry.Objects {
		objScreenPosition := rl.GetWorldToScreen(rl.Vector3{X: obj.Position().X, Y: obj.Position().Y + 1.5, Z: obj.Position().Z}, g.camera)
		rl.DrawText(obj.ID(), int32(objScreenPosition.X)-rl.MeasureText(obj.ID(), 20)/2, int32(objScreenPosition.Y), 20, rl.Black)
	}

	textRight := make([]string, 0)

	if g.network.IsConnected() {
		textRight = append(textRight, "Connected to Server - Press [C] to disconnect")
		// text := "Connected to Server - Press [C] to disconnect"
		// rl.DrawText(text, WIDTH-rl.MeasureText(text, 20)-10, 20, 20, rl.Black)
	} else {
		textRight = append(textRight, "Press [C] to connect to server")
		// text := "Press [C] to connect to server"
		// rl.DrawText(text, WIDTH-rl.MeasureText(text, 20)-10, 20, 20, rl.Black)
	}

	textRight = append(textRight, "Press [B] to shoot")
	textRight = append(textRight, "Press [N] to retrieve ball")
	textRight = append(textRight, fmt.Sprintf("%s's Points: %d", g.Username, g.points))

	fontSize := int32(16)
	positionY := int32(20)
	gap := int32(20)
	for _, text := range textRight {
		rl.DrawText(text, WIDTH-rl.MeasureText(text, fontSize)-10, positionY, fontSize, rl.Black)
		positionY += gap
	}

	rl.DrawFPS(20, 20)
}

func (g *Game) update3DOutput() {
	g.objRegistry.Lock()
	defer g.objRegistry.Unlock()

	for _, obj := range g.objRegistry.Objects {
		obj.Draw()
	}
	// g.ball.Draw()
	g.hoop.Draw()
	rl.DrawGrid(10, 1.0)
}

func (g *Game) processInput() {
	// move cube
	if rl.IsKeyDown(rl.KeyW) {
		g.box.Move(0.0, 0.0, -0.05)
		g.sendMovement(0.0, 0.0, -0.05)
	} else if rl.IsKeyDown(rl.KeyS) {
		g.box.Move(0.0, 0.0, 0.05)
		g.sendMovement(0.0, 0.0, 0.05)
	} else if rl.IsKeyDown(rl.KeyA) {
		g.box.Move(-0.05, 0.0, 0.0)
		g.sendMovement(-0.05, 0.0, 0.0)
	} else if rl.IsKeyDown(rl.KeyD) {
		g.box.Move(0.05, 0.0, 0.0)
		g.sendMovement(0.05, 0.0, 0.0)
	}

	// connect to server
	if rl.IsKeyPressed(rl.KeyC) && !g.network.IsConnected() {
		playerObjects := []GameObject{g.box}
		err := g.network.ConnectTCP(g.Username, playerObjects)
		if err != nil {
			log.Println("Error connecting to server via TCP:", err)
		}

		err = g.network.ConnectUDP()
		if err != nil {
			log.Println("Error connecting to server via UDP:", err)
		}

	} else if rl.IsKeyPressed(rl.KeyC) && g.network.IsConnected() {
		g.network.Disconnect()
	}

	// shoot ball
	if rl.IsKeyPressed(rl.KeyB) {
		g.ball.Shoot(g.hoop.rim.position, g.hoop.rim.radius)
	}

	// retrieve ball
	if rl.IsKeyPressed(rl.KeyN) {
		g.ball.AssignTo(g.box)
	}
}

func (g *Game) sendMovement(x float32, y float32, z float32) {
	if g.network.IsConnected() {
		g.network.WriteUDP(g.box.SerializeMove(x, y, z))
	}
}
