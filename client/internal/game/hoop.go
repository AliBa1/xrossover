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
	net       Net
}

type Backboard struct {
	position   rl.Vector3
	dimensions Dimensions
	color      color.RGBA
}

type Cylinder struct {
	position rl.Vector3
	radius   float32
	length   float32
	color    color.RGBA
}

type Net struct {
	center rl.Vector3
	radius float32
	length float32
	color  color.RGBA
}

func NewHoop(x, z float32) *Hoop {
	const poleHeight float32 = 10.0
	const backboardThickness float32 = 0.2
	const rimRadius float32 = 0.75
	rimX := x
	rimY := poleHeight
	rimZ := z + rimRadius + backboardThickness
	return &Hoop{
		backboard: Backboard{
			position:   rl.Vector3{X: x, Y: poleHeight + 0.75, Z: z},
			dimensions: Dimensions{Width: 6.0, Height: 3.5, Length: backboardThickness},
			color:      rl.Beige,
		},
		rim: Cylinder{
			position: rl.Vector3{X: rimX, Y: rimY, Z: rimZ},
			radius:   rimRadius,
			length:   0.06,
			color:    rl.Red,
		},
		pole: Cylinder{
			position: rl.Vector3{X: x, Y: GroundY, Z: z - 0.33},
			radius:   0.25,
			length:   poleHeight,
			color:    rl.Gray,
		},
		net: Net{
			center: rl.Vector3{X: rimX, Y: rimY, Z: rimZ},
			radius: rimRadius,
			length: 1.5,
			color:  rl.White,
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

// chatgpt wrote func
func drawNet(center rl.Vector3, radius, height float32, segments int, color rl.Color) {
	angleStep := 2 * math.Pi / float64(segments)

	for i := range segments {
		angle := float64(i) * angleStep
		x := center.X + radius*float32(math.Cos(angle))
		z := center.Z + radius*float32(math.Sin(angle))

		top := rl.NewVector3(x, center.Y, z)
		bottom := rl.NewVector3(x, center.Y-height, z)

		rl.DrawLine3D(top, bottom, color)
	}
}

// chatgpt wrote func
func drawNetCrissCross(center rl.Vector3, radius, height float32, segments int, color rl.Color) {
	angleStep := 2 * math.Pi / float64(segments)

	topPoints := make([]rl.Vector3, segments)
	bottomPoints := make([]rl.Vector3, segments)

	for i := range segments {
		angle := float64(i) * angleStep
		x := center.X + radius*float32(math.Cos(angle))
		z := center.Z + radius*float32(math.Sin(angle))

		top := rl.NewVector3(x, center.Y, z)
		bottom := rl.NewVector3(x, center.Y-height, z)

		topPoints[i] = top
		bottomPoints[i] = bottom

		// Draw verticals
		rl.DrawLine3D(top, bottom, color)
	}

	// Draw criss-cross diagonals
	for i := range segments {
		next := (i + 1) % segments

		// Top[i] -> Bottom[next]
		rl.DrawLine3D(topPoints[i], bottomPoints[next], color)

		// Top[next] -> Bottom[i]
		rl.DrawLine3D(topPoints[next], bottomPoints[i], color)
	}
}

// chatgpt wrote func
func drawRealisticNet(center rl.Vector3, topRadius, height float32, segments, layers int, color rl.Color) {
	// Generate ring levels: top to bottom
	rings := make([][]rl.Vector3, layers+1)

	for layer := 0; layer <= layers; layer++ {
		t := float32(layer) / float32(layers)
		y := center.Y - t*height
		radius := topRadius * (1 - 0.5*t) // taper

		ring := make([]rl.Vector3, segments)
		for i := 0; i < segments; i++ {
			angle := float64(i) * 2 * math.Pi / float64(segments)
			x := center.X + radius*float32(math.Cos(angle))
			z := center.Z + radius*float32(math.Sin(angle))
			ring[i] = rl.NewVector3(x, y, z)
		}
		rings[layer] = ring
	}

	// Draw vertical and cross ropes using cylinders
	var thickness float32
	thickness = 0.02

	for i := 0; i < segments; i++ {
		next := (i + 1) % segments
		for layer := 0; layer < layers; layer++ {
			// Cross ropes (diagonals)
			drawRope(rings[layer][i], rings[layer+1][next], thickness*0.8, color)
			drawRope(rings[layer][next], rings[layer+1][i], thickness*0.8, color)
		}
	}
}

func drawRope(start, end rl.Vector3, thickness float32, color rl.Color) {
	rl.DrawCylinderEx(start, end, thickness, thickness, 6, color)
}

func (h *Hoop) Draw() {
	// backboard
	rl.DrawCube(h.backboard.position, h.backboard.dimensions.Width, h.backboard.dimensions.Height, h.backboard.dimensions.Length, h.backboard.color)

	// rim
	drawRim(h.rim.position, h.rim.radius, h.rim.length, h.rim.color)

	// pole
	rl.DrawCylinder(h.pole.position, h.pole.radius, h.pole.radius, h.pole.length, 1, h.pole.color)

	// net
	// drawNet(h.rim.position, h.rim.radius, 1.5, 12, rl.White)
	// drawNetCrissCross(h.rim.position, h.rim.radius, 1.5, 12, rl.White)
	// drawNetCrissCross(h.net.center, h.net.radius, h.net.length, 12, h.net.color)
	drawRealisticNet(h.net.center, h.net.radius, h.net.length, 16, 3, h.net.color)
}
