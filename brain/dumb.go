package brain

import (
	"Gooo/agent"
	"math/rand"
)

// DumbBrain is an implementation of the agent.DecisionEngine interface.
// It ignores all observations and events and chooses random action. This implementation is used for testing purposes.
type DumbBrain struct {
}

func NewDumbBrain() DumbBrain {
	return DumbBrain{}
}

func (d DumbBrain) DecideOnObservation(obs agent.Observation) agent.Action {
	if rand.Float32() > 0.90 {
		return agent.ActionJump
	}
	return agent.ActionNothing
}

func (d DumbBrain) DecideOnEvent(obs agent.Observation) agent.Action {
	if rand.Float32() > 0.96 {
		return agent.ActionJump
	}
	return agent.ActionNothing
}
