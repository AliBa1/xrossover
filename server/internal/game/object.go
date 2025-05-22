package game

import (
	"errors"
	"sync"
)

type GameObject interface {
	ID() string
	Move(x float32, y float32, z float32)
	// Update(dt float32)
	Serialize() ([]byte, error)
}

type ObjectRegistry struct {
	sync.RWMutex
	objects map[string]GameObject
}

func NewObjectRegistry() *ObjectRegistry {
	return &ObjectRegistry{
		objects: make(map[string]GameObject),
	}
}

func (o *ObjectRegistry) Add(obj GameObject) {
	o.Lock()
	defer o.Unlock()
	o.objects[obj.ID()] = obj
}

func (o *ObjectRegistry) Remove(id string) {
	o.Lock()
	defer o.Unlock()
	delete(o.objects, id)
}

func (o *ObjectRegistry) Get(id string) (GameObject, error) {
	o.RLock()
	defer o.RUnlock()
	obj, ok := o.objects[id]
	if ok != false {
		return nil, errors.New("unable to find object (" + id + ")")
	}
	return obj, nil
}
