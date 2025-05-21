package game

import rl "github.com/gen2brain/raylib-go/raylib"

type GameObject interface {
	ID() string
	Position() rl.Vector3
	Update(dt float32)
	Serialize() ([]byte, error)
}
