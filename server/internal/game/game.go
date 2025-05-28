package game

import (
	"log"
	"sync"
	"time"
)

const (
	TickRate     = 30
	TickInterval = time.Second / TickRate
)

type Game struct {
	inputs         map[string][]PlayerInput
	ObjectRegistry *ObjectRegistry
	sync.Mutex
}

type PlayerInput struct {
	Username string
	ObjectID string
	Input    string
}

func (g *Game) Run() {
	g.initialize()
	g.loop()
	// g.shutdown()
}

func (g *Game) initialize() {
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
				log.Println(err)
			} else {
				g.processInput(object, i.Input)
			}
		}
		g.inputs[username] = nil

	}
}

func (g *Game) processInput(obj GameObject, input string) {

}

func (g *Game) AddPlayerInput(username, objectID, input string) {
	g.Lock()
	defer g.Unlock()

	g.inputs[username] = append(g.inputs[username], PlayerInput{
		Username: username,
		ObjectID: objectID,
		Input:    input,
	})
}
