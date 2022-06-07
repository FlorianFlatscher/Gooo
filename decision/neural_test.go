package decision

import (
	"GeneticAlgorithm/customMath"
	"fmt"
	"testing"
)

func TestNewNeuralBrain(t *testing.T) {
	t.Run("Test if neural decision gets constructed correctly", func(t *testing.T) {
		brain := NewNeuralBrain(NeuralBrainOptions{
			nil,
		})
		d := brain.DecideOnObservation(Observation{
			DistanceForward: 0,
			Position: customMath.Point{
				X: 0,
				Y: 0,
			},
			HeightOfHole: 0,
		})
		fmt.Printf("Decision: %d\n", d)
	})
}
