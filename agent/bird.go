package agent

import (
	"Gooo/customMath"
	"math"
)

type Bird struct {
	score    string
	position customMath.Point
	velocity float64
	brain    DecisionEngine
	dead     bool
}

func (b *Bird) DoSomething(distanceForward float64, heightOfHole float64) {
	var observation = Observation{
		DistanceForward: distanceForward,
		Position:        b.position,
		HeightOfHole:    heightOfHole,
	}

	var action = b.brain.DecideOnObservation(observation)

	switch action {
	case ActionJump:
		b.velocity = -0.06
	}
}

func (b *Bird) DoPhysics() {
	if b.velocity < 0 {
		b.velocity += 0.01
	}
	if math.Abs(b.velocity) <= 0.01 {
		b.velocity = 0
	}
	b.position = customMath.Point{
		Y: b.position.Y + 0.01 + b.velocity,
		X: b.position.X,
	}
}

// Constructor

func NewBird(brain DecisionEngine) *Bird {
	return &Bird{brain: brain, position: customMath.Point{X: 0.25, Y: 0.5}}
}

// Setter & Getter

func (b *Bird) Velocity() float64 {
	return b.velocity
}

func (b *Bird) SetVelocity(velocity float64) {
	b.velocity = velocity
}

func (b *Bird) Position() customMath.Point {
	return b.position
}

func (b *Bird) SetPosition(position customMath.Point) {
	b.position = position
}

func (b *Bird) Score() string {
	return b.score
}

func (b *Bird) SetScore(score string) {
	b.score = score
}

func (b *Bird) Brain() DecisionEngine {
	return b.brain
}

func (b *Bird) SetBrain(brain DecisionEngine) {
	b.brain = brain
}

func (b *Bird) Dead() bool {
	return b.dead
}

func (b *Bird) SetDead(dead bool) {
	b.dead = dead
}
