package brain

import (
	"Gooo/agent"
	"math/rand"
)

type DumbBrain struct {
}

func (d DumbBrain) DecideOnObservation(obs agent.Observation) agent.Action {
	if rand.Float32() > 0.90 {
		return agent.JUMP
	}
	return agent.NOTHING
}

func (d DumbBrain) DecideOnEvent(obs agent.Observation) agent.Action {
	if rand.Float32() > 0.96 {
		return agent.JUMP
	}
	return agent.NOTHING
}

func NewDumbBrain() DumbBrain {
	return DumbBrain{}
}
