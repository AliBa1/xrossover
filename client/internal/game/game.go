package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	WIDTH  = 800
	HEIGHT = 450
)

type Game struct {
	camera       rl.Camera3D
	cubePosition rl.Vector3
}

func (g *Game) Run() {
	rl.InitWindow(WIDTH, HEIGHT, "Game Window")
	defer rl.CloseWindow()

	g.initialize()
	g.loop()
}

func (g *Game) initialize() {
	g.camera = rl.Camera3D{
		Position:   rl.Vector3{X: 0.0, Y: 10.0, Z: 10.0},
		Target:     rl.Vector3{X: 0.0, Y: 0.0, Z: 0.0},
		Up:         rl.Vector3{X: 0.0, Y: 1.0, Z: 0.0},
		Fovy:       45.0,
		Projection: rl.CameraPerspective,
	}
	g.cubePosition = rl.Vector3{X: 0.0, Y: 1.0, Z: 0.0}

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

func (g *Game) updateDrawing() {
	rl.ClearBackground(rl.RayWhite)

	rl.BeginMode3D(g.camera)
	g.update3DOutput()
	rl.EndMode3D()

	// rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.LightGray)
	rl.DrawFPS(20, 20)
}

func (g *Game) update3DOutput() {
	rl.DrawCube(g.cubePosition, 1.0, 1.0, 1.0, rl.Red)
	rl.DrawCubeWires(g.cubePosition, 1.0, 1.0, 1.0, rl.Maroon)
	rl.DrawGrid(10, 1.0)
}

func (g *Game) processInput() {
	if rl.IsKeyDown(rl.KeyW) {
		g.cubePosition.Z -= 0.05
	} else if rl.IsKeyDown(rl.KeyS) {
		g.cubePosition.Z += 0.05
	} else if rl.IsKeyDown(rl.KeyA) {
		g.cubePosition.X -= 0.05
	} else if rl.IsKeyDown(rl.KeyD) {
		g.cubePosition.X += 0.05
	}
}
