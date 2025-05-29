package game

import (
	"log"
	"sync"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	TickRate     = 30
	TickInterval = time.Second / TickRate
)

type Game struct {
	inputs         map[string][]PlayerInput
	ObjectRegistry *ObjectRegistry
	Broadcast      func(protocol, owner string, data []byte)
	sync.Mutex
}

// type PlayerInput struct {
// 	Username string
// 	ObjectID string
// 	Input    string
// }

type PlayerInput struct {
	ObjectID string
	Action   Action
}

type Action interface {
	Type() string
}

type Move struct {
	Direction rl.Vector3
}

func (m Move) Type() string { return "move" }

func (g *Game) Run() {
	g.initialize()
	g.loop()
	// g.shutdown()
}

func (g *Game) initialize() {
	g.inputs = make(map[string][]PlayerInput)
	g.ObjectRegistry = NewObjectRegistry()
}

func (g *Game) loop() {
	ticker := time.NewTicker(TickInterval)
	defer ticker.Stop()

	for tick := 0; ; tick++ {
		select {
		case <-ticker.C:
			g.update()
		}
	}
}

func (g *Game) update() {
	g.Lock()
	defer g.Unlock()

	for username, inputs := range g.inputs {
		for _, i := range inputs {
			object, err := g.ObjectRegistry.Get(i.ObjectID)
			if err != nil {
				log.Println("Error geting obj", err)
			} else {
				g.processInput(object, i.Action)
			}
		}
		g.inputs[username] = nil

	}
}

func (g *Game) processInput(obj GameObject, action Action) {
	switch v := action.(type) {
	case Move:
		obj.Move(v.Direction.X, v.Direction.Y, v.Direction.Z)
		log.Println("applied movement")
		// Broadcast()
		if g.Broadcast != nil {
			g.Broadcast("tcp", obj.Owner(), obj.Serialize())
		}
	}
}

func (g *Game) AddPlayerInput(input PlayerInput) {
	g.Lock()
	defer g.Unlock()

	obj, err := g.ObjectRegistry.Get(input.ObjectID)
	if err != nil {
		log.Println(err)
		return
	}

	g.inputs[obj.Owner()] = append(g.inputs[obj.Owner()], PlayerInput{
		ObjectID: input.ObjectID,
		Action:   input.Action,
	})
}
