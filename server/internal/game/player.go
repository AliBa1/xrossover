package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	protocol "xrossover-server/flatbuffers/xrossover"

	flatbuffers "github.com/google/flatbuffers/go"
)

type PlayerBox struct {
	id       string
	owner    string
	position rl.Vector3
}

func NewPlayerBox(id, owner string, pos protocol.Vector3) *PlayerBox {
	return &PlayerBox{
		id:       id,
		owner:    owner,
		position: rl.Vector3{X: pos.X(), Y: pos.Y(), Z: pos.Z()},
	}
}

func (p *PlayerBox) ID() string {
	return p.id
}

func (p *PlayerBox) Owner() string {
	return p.owner
}

func (p *PlayerBox) Position() rl.Vector3 {
	return p.position
}

// func (p *PlayerBox) Update()

func (p *PlayerBox) UpdatePosition(x, y, z float32) {}

func (p *PlayerBox) Move(x, y, z float32) {
	p.position.X += x
	p.position.Y += y
	p.position.Z += z
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

func (p *PlayerBox) SerializeRegistry(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	id := builder.CreateString(p.id)
	owner := builder.CreateString(p.owner)

	protocol.PlayerBoxStart(builder)
	protocol.PlayerBoxAddId(builder, id)
	protocol.PlayerBoxAddOwner(builder, owner)
	protocol.PlayerBoxAddPosition(builder, protocol.CreateVector3(builder, p.position.X, p.position.Y, p.position.Z))
	playerBox := protocol.PlayerBoxEnd(builder)

	protocol.GameObjectWrapperStart(builder)
	protocol.GameObjectWrapperAddObjectType(builder, protocol.GameObjectUnionPlayerBox)
	protocol.GameObjectWrapperAddObject(builder, playerBox)
	gameObjectWrapper := protocol.GameObjectWrapperEnd(builder)
	return gameObjectWrapper
}
