package brain

import (
	"Gooo/agent"
	"fmt"
	"gonum.org/v1/gonum/mat"
	"math/rand"
)

// NeuralBrain is an implementation of the agent.DecisionEngine interface.
// It uses a neural network to decide based on certain observation.
type NeuralBrain struct {
	weights [3]*mat.Dense
}

func (n *NeuralBrain) DecideOnObservation(obs agent.Observation) agent.Action {
	input := mat.NewDense(3, 1, []float64{
		obs.DistanceForward,
		obs.HeightOfHole,
		obs.Position.Y,
	})

	curr := input
	for _, layer := range n.weights {
		layerRows, _ := layer.Dims()
		nextCurr := mat.NewDense(layerRows, 1, nil)
		// Bias
		currRows, _ := curr.Dims()
		currWithBias := mat.NewDense(currRows+1, 1, nil)
		currWithBias.Stack(mat.NewDense(1, 1, []float64{1}), curr)
		nextCurr.Mul(layer, currWithBias)
		curr = nextCurr
	}

	resNothing := curr.At(0, 0)
	resJump := curr.At(1, 0)
	if resJump > resNothing {
		return agent.ActionJump
	}
	return agent.ActionNothing
}

func (n *NeuralBrain) DecideOnEvent(obs agent.Observation) agent.Action {
	return agent.ActionNothing
}

func (n *NeuralBrain) Print() {
	for i := range n.weights {
		fmt.Println(n.weights[i])
	}
}

func NewNeuralBrain() *NeuralBrain {

	weights := [3]*mat.Dense{
		newDenseWithRandomValues(5, 4),
		newDenseWithRandomValues(4, 6),
		newDenseWithRandomValues(2, 5),
	}

	return &NeuralBrain{weights}
}

func newDenseWithRandomValues(r int, c int) *mat.Dense {
	data := make([]float64, r*c)
	for i := 0; i < r*c; i++ {
		data[i] = rand.Float64()
	}
	return mat.NewDense(r, c, data)
}
