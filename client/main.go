package main

import (
	"fmt"
	"xrossover-client/internal/game"
)

const (
	WIDTH  = 800
	HEIGHT = 450
)

func main() {
	fmt.Println("Welcome to the xrossover client!")
	game := &game.Game{}
	game.Run()
}
