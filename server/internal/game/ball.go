package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	protocol "xrossover-server/flatbuffers/xrossover"

	flatbuffers "github.com/google/flatbuffers/go"
)

type Ball struct {
	id       string
	owner    string
	position rl.Vector3
}

func NewBall(id, owner string, pos protocol.Vector3) *Ball {
	return &Ball{
		id:       id,
		owner:    owner,
		position: rl.Vector3{X: pos.X(), Y: pos.Y(), Z: pos.Z()},
	}
}

func (b *Ball) ID() string {
	return b.id
}

func (b *Ball) Owner() string {
	return b.owner
}

func (b *Ball) Position() rl.Vector3 {
	return b.position
}

// func (p *Ball) Update()

func (b *Ball) Move(x float32, y float32, z float32) {
	b.position.X += x
	b.position.Y += y
	b.position.Z += z
}

func (b *Ball) Serialize() []byte {
	builder := flatbuffers.NewBuilder(1024)

	id := builder.CreateString(b.id)
	owner := builder.CreateString(b.owner)

	protocol.BallStart(builder)
	protocol.BallAddId(builder, id)
	protocol.BallAddOwner(builder, owner)
	protocol.BallAddPosition(builder, protocol.CreateVector3(builder, b.position.X, b.position.Y, b.position.Z))
	ball := protocol.BallEnd(builder)

	protocol.NetworkMessageStart(builder)
	protocol.NetworkMessageAddPayloadType(builder, protocol.PayloadBall)
	protocol.NetworkMessageAddPayload(builder, ball)
	netMsg := protocol.NetworkMessageEnd(builder)

	builder.Finish(netMsg)

	return builder.FinishedBytes()
}
