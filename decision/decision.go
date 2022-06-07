package decision

import "GeneticAlgorithm/customMath"

type Action int

const (
	ActionJump    = Action(0)
	ActionNothing = Action(1)
)

type Observation struct {
	DistanceForward float64
	Position        customMath.Point
	HeightOfHole    float64
}

type DecisionEngine interface {
	DecideOnObservation(obs Observation) Action
	DecideOnEvent(obs Observation) Action
}
