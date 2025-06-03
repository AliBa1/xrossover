package game

import (
	"image/color"

	protocol "xrossover-client/flatbuffers/xrossover"

	rl "github.com/gen2brain/raylib-go/raylib"

	flatbuffers "github.com/google/flatbuffers/go"
)

type PlayerBox struct {
	id       string
	owner    string
	position rl.Vector3
	width    float32
	height   float32
	length   float32
	color    color.RGBA
}

func NewPlayerBox(id, owner string) *PlayerBox {
	return &PlayerBox{
		id:       id,
		owner:    owner,
		position: rl.Vector3{X: 0.0, Y: 1.0, Z: 0.0},
		width:    1.0,
		height:   1.0,
		length:   1.0,
		color:    rl.Red,
	}
}

func NewFBPlayerBox(id, owner string, pos protocol.Vector3) *PlayerBox {
	return &PlayerBox{
		id:       id,
		owner:    owner,
		position: rl.Vector3{X: pos.X(), Y: pos.Y(), Z: pos.Z()},
		width:    1.0,
		height:   1.0,
		length:   1.0,
		color:    rl.Red,
	}
}

func (p *PlayerBox) ID() string           { return p.id }
func (p *PlayerBox) Owner() string        { return p.owner }
func (p *PlayerBox) Position() rl.Vector3 { return p.position }
func (p *PlayerBox) Color() color.RGBA    { return p.color }

func (p *PlayerBox) Dimensions() Dimensions {
	return Dimensions{
		Width:  p.width,
		Height: p.height,
		Length: p.length,
	}
}

func (p *PlayerBox) Update(dt float32) {

}

func (p *PlayerBox) Move(x, y, z float32) {
	p.position.X += x
	p.position.Y += y
	p.position.Z += z
}

func (p *PlayerBox) UpdatePosition(x, y, z float32) {
	p.position.X = x
	p.position.Y = y
	p.position.Z = z
}

func (p *PlayerBox) Serialize() []byte {
	builder := flatbuffers.NewBuilder(1024)

	id := builder.CreateString(p.id)
	owner := builder.CreateString(p.owner)

	protocol.PlayerBoxStart(builder)
	protocol.PlayerBoxAddId(builder, id)
	protocol.PlayerBoxAddOwner(builder, owner)
	protocol.PlayerBoxAddPosition(builder, protocol.CreateVector3(builder, p.position.X, p.position.Y, p.position.Z))
	playerBox := protocol.PlayerBoxEnd(builder)

	protocol.NetworkMessageStart(builder)
	protocol.NetworkMessageAddPayloadType(builder, protocol.PayloadPlayerBox)
	protocol.NetworkMessageAddPayload(builder, playerBox)
	netMsg := protocol.NetworkMessageEnd(builder)

	builder.Finish(netMsg)

	return builder.FinishedBytes()
}

func (p *PlayerBox) SerializeMove(x float32, y float32, z float32) []byte {
	builder := flatbuffers.NewBuilder(1024)
	objID := builder.CreateString(p.id)

	protocol.MoveStart(builder)
	protocol.MoveAddDirection(builder, protocol.CreateVector3(builder, x, y, z))
	move := protocol.MoveEnd(builder)

	protocol.PlayerInputStart(builder)
	protocol.PlayerInputAddActionType(builder, protocol.ActionMove)
	protocol.PlayerInputAddAction(builder, move)
	protocol.PlayerInputAddObjectId(builder, objID)
	playerInput := protocol.PlayerInputEnd(builder)

	protocol.NetworkMessageStart(builder)
	protocol.NetworkMessageAddPayloadType(builder, protocol.PayloadPlayerInput)
	protocol.NetworkMessageAddPayload(builder, playerInput)
	netMsg := protocol.NetworkMessageEnd(builder)

	builder.Finish(netMsg)

	return builder.FinishedBytes()
}
