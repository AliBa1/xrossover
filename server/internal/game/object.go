package game

import (
	"errors"
	"fmt"
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
	fmt.Println("Object to add:", obj)
	fmt.Println("Obj ID", obj.ID())
	o.Objects[obj.ID()] = obj
	fmt.Println("Objs after addition:", o.Objects)
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
