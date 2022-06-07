package bird

import (
	"GeneticAlgorithm/customMath"
	"GeneticAlgorithm/decision"
	"math"
)

type Bird struct {
	score          float64
	position       customMath.Point
	velocity       float64
	decisionEngine decision.DecisionEngine
	dead           bool
}

func (b *Bird) DoSomething(distanceForward float64, heightOfHole float64) {
	var observation = decision.Observation{
		DistanceForward: distanceForward,
		Position:        b.position,
		HeightOfHole:    heightOfHole,
	}
	var action = b.decisionEngine.DecideOnObservation(observation)

	switch action {
	case decision.ActionJump:
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
		Y: b.position.Y + 0.02 + b.velocity,
		X: b.position.X,
	}
}

type Options struct {
	DecisionEngine decision.DecisionEngine
}

// Constructor

func NewBird(o Options) *Bird {
	return &Bird{decisionEngine: o.DecisionEngine, position: customMath.Point{X: 0.25, Y: 0.5}}
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

func (b *Bird) Score() float64 {
	return b.score
}

func (b *Bird) SetScore(score float64) {
	b.score = score
}

func (b *Bird) IncrementScore(by float64) {
	b.score += by
}

func (b *Bird) DecisionEngine() decision.DecisionEngine {
	return b.decisionEngine
}

func (b *Bird) SetDecisionEngine(brain decision.DecisionEngine) {
	b.decisionEngine = brain
}

func (b *Bird) Dead() bool {
	return b.dead
}

func (b *Bird) SetDead(dead bool) {
	b.dead = dead
}

// Revive clones the bird
func (b *Bird) Revive() Bird {
	return *NewBird(Options{
		DecisionEngine: b.decisionEngine,
	})
}
