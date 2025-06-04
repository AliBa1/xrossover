package game

import (
	"image/color"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Hoop struct {
	backboard Backboard
	rim       Cylinder
	pole      Cylinder
	// net       Net
}

type Backboard struct {
	position   rl.Vector3
	dimensions Dimensions
	color      color.RGBA
}

type Cylinder struct {
	position rl.Vector3
	radius   float32
	height   float32
	color    color.RGBA
}

func NewHoop(x, z float32) *Hoop {
	const poleHeight float32 = 10.0
	return &Hoop{
		backboard: Backboard{
			position:   rl.Vector3{X: x, Y: poleHeight + 0.75, Z: z},
			dimensions: Dimensions{Width: 6.0, Height: 3.5, Length: 0.2},
			color:      rl.Beige,
		},
		rim: Cylinder{
			position: rl.Vector3{X: x, Y: poleHeight, Z: z - 0.5},
			radius:   0.75,
			height:   0.06,
			color:    rl.Red,
		},
		pole: Cylinder{
			position: rl.Vector3{X: x, Y: GroundY, Z: z - 0.33},
			radius:   0.25,
			height:   poleHeight,
			color:    rl.Gray,
		},
	}
}

// chatgpt wrote func
func drawRim(center rl.Vector3, radius float32, thickness float32, color rl.Color) {
	const segmentCount = 24
	angleStep := 2 * math.Pi / float64(segmentCount)

	for i := range segmentCount {
		angle := float64(i) * angleStep
		nextAngle := float64(i+1) * angleStep

		// Position of this segment
		x1 := center.X + radius*float32(math.Cos(angle))
		z1 := center.Z + radius*float32(math.Sin(angle))

		x2 := center.X + radius*float32(math.Cos(nextAngle))
		z2 := center.Z + radius*float32(math.Sin(nextAngle))

		start := rl.NewVector3(x1, center.Y, z1)
		end := rl.NewVector3(x2, center.Y, z2)

		// You can draw this segment as a line or tiny cylinder
		rl.DrawCylinderEx(start, end, thickness, thickness, 6, color)
	}
}

func (h *Hoop) Draw() {
	// backboard
	rl.DrawCube(h.backboard.position, h.backboard.dimensions.Width, h.backboard.dimensions.Height, h.backboard.dimensions.Length, h.backboard.color)

	// rim
	drawRim(h.rim.position, h.rim.radius, h.rim.height, h.rim.color)

	// pole
	rl.DrawCylinder(h.pole.position, h.pole.radius, h.pole.radius, h.pole.height, 1, h.pole.color)
}
