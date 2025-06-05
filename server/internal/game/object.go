package game

import (
	"errors"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
	flatbuffers "github.com/google/flatbuffers/go"
	protocol "xrossover-server/flatbuffers/xrossover"
)

type GameObject interface {
	ID() string
	Owner() string
	Position() rl.Vector3
	Move(x float32, y float32, z float32)
	// Update(dt float32)
	UpdatePosition(x, y, z float32)
	Serialize() []byte
	SerializeRegistry(builder *flatbuffers.Builder) flatbuffers.UOffsetT
}

type ObjectRegistry struct {
	sync.RWMutex
	Objects map[string]GameObject
}

func NewObjectRegistry() *ObjectRegistry {
	return &ObjectRegistry{
		Objects: make(map[string]GameObject),
	}
}

func (o *ObjectRegistry) Add(obj GameObject) {
	o.Lock()
	defer o.Unlock()

	o.Objects[obj.ID()] = obj
}

func (o *ObjectRegistry) Remove(id string) {
	o.Lock()
	defer o.Unlock()
	delete(o.Objects, id)
}

func (o *ObjectRegistry) Get(id string) (GameObject, error) {
	o.RLock()
	defer o.RUnlock()
	obj, ok := o.Objects[id]
	if ok != true {
		return nil, errors.New("unable to find object (" + id + ")")
	}
	return obj, nil
}

func (o *ObjectRegistry) Serialize() []byte {
	o.RLock()
	defer o.RUnlock()

	builder := flatbuffers.NewBuilder(1024)

	var objects []flatbuffers.UOffsetT

	for _, obj := range o.Objects {
		objects = append(objects, obj.SerializeRegistry(builder))
	}

	protocol.ObjectRegistryStartObjectsVector(builder, len(objects))
	for i := len(objects) - 1; i >= 0; i-- {
		builder.PrependUOffsetT(objects[i])
	}
	objectsVector := builder.EndVector(len(objects))

	protocol.ObjectRegistryStart(builder)
	protocol.ObjectRegistryAddObjects(builder, objectsVector)
	objectRegistry := protocol.ObjectRegistryEnd(builder)

	protocol.NetworkMessageStart(builder)
	protocol.NetworkMessageAddPayloadType(builder, protocol.PayloadObjectRegistry)
	protocol.NetworkMessageAddPayload(builder, objectRegistry)
	netMsg := protocol.NetworkMessageEnd(builder)

	builder.Finish(netMsg)

	return builder.FinishedBytes()
}
