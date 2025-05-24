package game

import (
	"errors"
	"image/color"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Dimensions struct {
	Width  float32
	Height float32
	Length float32
}

type GameObject interface {
	ID() string
	Position() rl.Vector3
	Dimensions() Dimensions
	Color() color.RGBA
	// Update(dt float32)
	Serialize() []byte
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

	_, ok := o.Objects[obj.ID()]
	if !ok {
		o.Objects[obj.ID()] = obj
	}
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
