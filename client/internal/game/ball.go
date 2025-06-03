package game

import (
	// protocol "xrossover-client/flatbuffers/xrossover"

	rl "github.com/gen2brain/raylib-go/raylib"

	flatbuffers "github.com/google/flatbuffers/go"
)

type Ball struct {
	id        string
	owner     string
	possessor GameObject
	position  rl.Vector3
	radius    float32
	velocity  rl.Vector3
}

func NewBall(id, owner string, possessor GameObject) *Ball {
	return &Ball{
		id:        id,
		owner:     owner,
		possessor: possessor,
		position:  rl.Vector3{X: 0.0, Y: 3.0, Z: 0.0},
		radius:    0.5,
		velocity:  rl.Vector3{X: 0.0, Y: 1.0, Z: 0.0},
	}
}

func (b *Ball) ID() string           { return b.id }
func (b *Ball) Owner() string        { return b.owner }
func (b *Ball) Position() rl.Vector3 { return b.position }

func (b *Ball) Update(dt float32) {
	if b.possessor != nil {
		b.position.X = b.possessor.Position().X + 0.5
		b.position.Z = b.possessor.Position().Z - 0.5
	} else {
		b.position.X += b.velocity.X * dt
		b.position.Z += b.velocity.Z * dt
	}

	b.position.Y += b.velocity.Y * dt

	if b.position.Y-b.radius < 0 {
		b.velocity.Y *= -1
		b.position.Y = b.radius
	} else if b.position.Y+b.radius > 5 {
		b.velocity.Y *= -1
		b.position.Y = 5 - b.radius
	}
}

func (b *Ball) UpdatePosition(x, y, z float32) {
	b.position.X = x
	b.position.Y = y
	b.position.Z = z
}

func (b *Ball) Serialize() []byte {
	builder := flatbuffers.NewBuilder(1024)

	// id := builder.CreateString(p.id)
	// owner := builder.CreateString(p.owner)
	//
	// protocol.PlayerBoxStart(builder)
	// protocol.PlayerBoxAddId(builder, id)
	// protocol.PlayerBoxAddOwner(builder, owner)
	// protocol.PlayerBoxAddPosition(builder, protocol.CreateVector3(builder, p.position.X, p.position.Y, p.position.Z))
	// playerBox := protocol.PlayerBoxEnd(builder)
	//
	// protocol.NetworkMessageStart(builder)
	// protocol.NetworkMessageAddPayloadType(builder, protocol.PayloadPlayerBox)
	// protocol.NetworkMessageAddPayload(builder, playerBox)
	// netMsg := protocol.NetworkMessageEnd(builder)
	//
	// builder.Finish(netMsg)

	return builder.FinishedBytes()
}
