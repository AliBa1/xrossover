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
	fmt.Println("Welcome to the xrossover client!")
	game := &game.Game{
		Username: username,
		UDPPort:  udpPort,
	}
	game.Run()
}
