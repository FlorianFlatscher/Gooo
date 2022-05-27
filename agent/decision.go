package agent

import "Gooo/math"

type Action int

const (
	JUMP    = Action(0)
	NOTHING = Action(1)
)

type Observation struct {
	DistanceForward float64
	Position        math.Point
	HeightOfHole    float64
}

type DecisionEngine interface {
	DecideOnObservation(obs Observation) Action
	DecideOnEvent(obs Observation) Action
}
