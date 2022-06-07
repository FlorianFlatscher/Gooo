package decision

import (
	"math/rand"
)

// DumbBrain is an implementation of the bird.DecisionEngine interface.
// It ignores all observations and events and chooses random action. This implementation is used for testing purposes.
type DumbBrain struct {
}

func NewDumbBrain() DumbBrain {
	return DumbBrain{}
}

func (d *DumbBrain) DecideOnObservation(obs Observation) Action {
	if rand.Float32() > 0.90 {
		return ActionJump
	}
	return ActionNothing
}

func (d *DumbBrain) DecideOnEvent(obs Observation) Action {
	if rand.Float32() > 0.96 {
		return ActionJump
	}
	return ActionNothing
}
