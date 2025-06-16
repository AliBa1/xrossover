package game

import (
	"fmt"
	"math"
	"math/rand"
	protocol "xrossover-client/flatbuffers/xrossover"

	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"

	flatbuffers "github.com/google/flatbuffers/go"
)

type Ball struct {
	id           string
	owner        string
	possessor    GameObject
	position     rl.Vector3
	radius       float32
	velocity     rl.Vector3
	acceleration rl.Vector3
	color        color.RGBA
}

func (b *Ball) Shoot(target rl.Vector3, rimRadius float32) {
	// adjust for different arc
	time := rand.Float32() + 0.75

	// randomize shot target
	minOffsetX := target.X - rimRadius
	maxOffsetX := target.X + rimRadius
	offsetX := rand.Float32()*(maxOffsetX-minOffsetX) + minOffsetX
	target.X = offsetX
	minOffsetZ := target.Z - rimRadius
	maxOffsetZ := target.Z + rimRadius
	offsetZ := rand.Float32()*(maxOffsetZ-minOffsetZ) + minOffsetZ
	target.Z = offsetZ

	gravityVec := rl.Vector3{
		X: 0,
		Y: Gravity,
		Z: 0,
	}

	displacement := rl.Vector3Subtract(target, b.position)
	halfGT2 := rl.Vector3Scale(gravityVec, 0.5*time*time)
	numerator := rl.Vector3Subtract(displacement, halfGT2)
	newVelocity := rl.Vector3Scale(numerator, 1.0/time)

	b.velocity = newVelocity

	b.possessor = nil
}

func (b *Ball) AssignTo(player GameObject) {
	b.possessor = player
}

func NewBall(id, owner string, possessor GameObject) *Ball {
	return &Ball{
		id:           id,
		owner:        owner,
		possessor:    possessor,
		position:     rl.Vector3{X: 0.0, Y: 3.0, Z: 0.0},
		radius:       0.5,
		velocity:     rl.Vector3{X: 0.0, Y: 3.0, Z: 0.0},
		acceleration: rl.Vector3{X: 0.0, Y: Gravity, Z: 0.0},
		color:        rl.Orange,
	}
}

func NewFBBall(id, owner string, pos protocol.Vector3) *Ball {
	return &Ball{
		id:           id,
		owner:        owner,
		position:     rl.Vector3{X: pos.X(), Y: pos.Y(), Z: pos.Z()},
		radius:       0.5,
		velocity:     rl.Vector3{X: 0.0, Y: 3.0, Z: 0.0},
		acceleration: rl.Vector3{X: 0.0, Y: Gravity, Z: 0.0},
		color:        rl.Blue,
	}
}

func (b *Ball) ID() string           { return b.id }
func (b *Ball) Owner() string        { return b.owner }
func (b *Ball) Position() rl.Vector3 { return b.position }
func (b *Ball) Color() color.RGBA    { return b.color }

func (b *Ball) Update(dt float32) {
	var bounceHeight float32
	if b.possessor != nil {
		b.position.X = b.possessor.Position().X + 0.5
		b.position.Z = b.possessor.Position().Z - 0.5
		bounceHeight = 3.0
	} else {
		b.position.X += b.velocity.X * dt
		b.position.Z += b.velocity.Z * dt
		bounceHeight = 3000000
	}

	b.applyGravity(dt)

	if b.position.Y-b.radius < GroundY {
		b.velocity.Y *= -1
		// b.velocity.Y *= -0.7
		b.position.Y = b.radius
	} else if b.position.Y+b.radius > bounceHeight {
		b.velocity.Y *= -1
		b.position.Y = bounceHeight - b.radius
	}
}

func (b *Ball) DetectCollision(dt float32, hoop Hoop) {
	backboardLeftX := hoop.backboard.position.X - hoop.backboard.dimensions.Width
	backboardRightX := hoop.backboard.position.X + hoop.backboard.dimensions.Width
	ballLeftX := b.position.X - b.radius
	ballRightX := b.position.X + b.radius
	collidesBackboardX := ballRightX > backboardLeftX && ballLeftX < backboardRightX

	backboardDownY := hoop.backboard.position.Y - hoop.backboard.dimensions.Height
	backboardUpY := hoop.backboard.position.Y + hoop.backboard.dimensions.Height
	ballDownY := b.position.Y - b.radius
	ballUpY := b.position.Y + b.radius
	collidesBackboardY := ballUpY > backboardDownY && ballDownY < backboardUpY

	backboardBackZ := hoop.backboard.position.Z - hoop.backboard.dimensions.Length
	backboardForwardZ := hoop.backboard.position.Z + hoop.backboard.dimensions.Length
	ballBackZ := b.position.Z - b.radius
	ballForwardZ := b.position.Z + b.radius
	collidesBackboardZ := ballForwardZ > backboardBackZ && ballBackZ < backboardForwardZ

	if collidesBackboardX && collidesBackboardY && collidesBackboardZ {
		fmt.Println("ball and backboard collision")
		// b.velocity.X *= -1
		// b.velocity.Y *= -1
		b.velocity.Z *= -1
	}

	// chatgpt rim collision
	rimCenter := hoop.rim.position
	ballXZ := rl.Vector2{X: b.position.X, Y: b.position.Z}
	rimCenterXZ := rl.Vector2{X: rimCenter.X, Y: rimCenter.Z}

	// 2D distance from ball center to rim center (in XZ plane)
	distanceXZ := rl.Vector2Length(rl.Vector2Subtract(ballXZ, rimCenterXZ))

	// Rim collision condition: within rim radius +/- ball radius
	if distanceXZ <= hoop.rim.radius+b.radius && distanceXZ >= hoop.rim.radius-b.radius {
		// Optional: Add vertical position check too, if needed
		if math.Abs(float64(b.position.Y-rimCenter.Y)) <= float64(b.radius) {
			fmt.Println("ball and rim collision")
			// Respond to collision - e.g., reflect horizontal velocity
			// This is simplified and can be improved with proper normals
			b.velocity.X *= -0.7 // Invert and dampen horizontal velocity
			b.velocity.Z *= -0.7
		}
	}
}

func (b *Ball) applyGravity(dt float32) {
	b.velocity.Y += b.acceleration.Y * dt
	b.position.Y += b.velocity.Y * dt
}

func (b *Ball) Draw() {
	rl.DrawSphere(b.position, b.radius, b.color)
}

func (b *Ball) UpdatePosition(x, y, z float32) {
	b.position.X = x
	b.position.Y = y
	b.position.Z = z
}

func (b *Ball) Serialize() []byte {
	builder := flatbuffers.NewBuilder(1024)

	id := builder.CreateString(b.id)
	owner := builder.CreateString(b.owner)

	protocol.BallStart(builder)
	protocol.BallAddId(builder, id)
	protocol.BallAddOwner(builder, owner)
	protocol.BallAddPosition(builder, protocol.CreateVector3(builder, b.position.X, b.position.Y, b.position.Z))
	ball := protocol.BallEnd(builder)

	protocol.NetworkMessageStart(builder)
	protocol.NetworkMessageAddPayloadType(builder, protocol.PayloadBall)
	protocol.NetworkMessageAddPayload(builder, ball)
	netMsg := protocol.NetworkMessageEnd(builder)

	builder.Finish(netMsg)

	return builder.FinishedBytes()
}
