package main

import (
	"fmt"
	"os"
	"xrossover-client/internal/game"
)

const (
	WIDTH  = 800
	HEIGHT = 450
)

func main() {
	username := os.Args[1]
	udpPort := os.Args[2]
	fmt.Println("Starting xrossover for", username)
	objRegistry := game.NewObjectRegistry()
	network := game.NewNetwork("localhost", udpPort, objRegistry)
	game := game.NewGame(username, network, objRegistry)

	game.Run()
}
