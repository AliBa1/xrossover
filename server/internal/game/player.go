package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	protocol "xrossover-server/flatbuffers/xrossover"

	flatbuffers "github.com/google/flatbuffers/go"
)

type PlayerBox struct {
	id       string
	position rl.Vector3
	// width    float32
	// height   float32
	// length   float32
	// color    color.RGBA
}

func NewPlayerBox(id string, pos protocol.Vector3) *PlayerBox {
	return &PlayerBox{
		id:       id,
		position: rl.Vector3{X: pos.X(), Y: pos.Y(), Z: pos.Z()},
		// width:    1.0,
		// height:   1.0,
		// length:   1.0,
		// color:    rl.Red,
	}
}

func (p *PlayerBox) ID() string {
	return p.id
}

// func (p *PlayerBox) Position() rl.Vector3 {
// 	return p.position
// }

// func (p *PlayerBox) Update()

func (p *PlayerBox) Move(x float32, y float32, z float32) {
	p.position.X += x
	p.position.Y += y
	p.position.Z += z
}

func (p *PlayerBox) Serialize() ([]byte, error) {
	builder := flatbuffers.NewBuilder(1024)

	id := builder.CreateString(p.id)
	position := protocol.CreateVector3(builder, p.position.X, p.position.Y, p.position.Z)

	protocol.PlayerBoxStart(builder)
	protocol.PlayerBoxAddId(builder, id)
	protocol.PlayerBoxAddPosition(builder, position)
	playerBox := protocol.PlayerBoxEnd(builder)

	protocol.NetworkMessageStart(builder)
	protocol.NetworkMessageAddPayloadType(builder, protocol.PayloadPlayerBox)
	protocol.NetworkMessageAddPayload(builder, playerBox)
	netMsg := protocol.NetworkMessageEnd(builder)

	builder.Finish(netMsg)

	return builder.FinishedBytes(), nil
}
