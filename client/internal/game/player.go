package game

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type PlayerBox struct {
	id       string
	position rl.Vector3
	width    float32
	height   float32
	length   float32
	color    color.RGBA
}

func NewPlayerBox() *PlayerBox {
	return &PlayerBox{
		id:       "PBMJ",
		position: rl.Vector3{X: 0.0, Y: 1.0, Z: 0.0},
		width:    1.0,
		height:   1.0,
		length:   1.0,
		color:    rl.Red,
	}
}

func (p *PlayerBox) ID() string {
	return p.id
}

func (p *PlayerBox) Position() rl.Vector3 {
	return p.position
}

// func (p *PlayerBox) Update()

func (p *PlayerBox) Move(x float32, y float32, z float32) {
	p.position.X += x
	p.position.Y += y
	p.position.Z += z
}
