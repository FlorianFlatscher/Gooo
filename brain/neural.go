package brain

import (
	"Gooo/agent"
	"fmt"
	"gonum.org/v1/gonum/mat"
	"math"
	"math/rand"
)

const Layers = 3

// NeuralBrain is an implementation of the agent.DecisionEngine interface.
// It uses a neural network to decide based on certain observation.
type NeuralBrain struct {
	weights [Layers]*mat.Dense
}

func (n *NeuralBrain) DecideOnObservation(obs agent.Observation) agent.Action {
	input := mat.NewDense(3, 1, []float64{
		obs.DistanceForward,
		obs.HeightOfHole,
		obs.Position.Y,
	})

	curr := input
	for _, layer := range n.weights {
		// Prepare layer output
		layerRows, _ := layer.Dims()
		nextCurr := mat.NewDense(layerRows, 1, nil)
		// ReLu
		curr.Apply(func(i, j int, v float64) float64 {
			return math.Max(v, 0)
		}, curr)
		// Bias
		currRows, _ := curr.Dims()
		currWithBias := mat.NewDense(currRows+1, 1, nil)
		currWithBias.Stack(mat.NewDense(1, 1, []float64{1}), curr)
		// Forward
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

// CrossOver constructs a new NeuralBrain that uses k*100% weights of a, and (1-k)*100% weights of b, where k is between 0 and 1
func CrossOver(a *NeuralBrain, b *NeuralBrain, k float64) *NeuralBrain {
	if k < 0 || k > 1 {
		panic("invalid argument k")
	}

	// Parents
	weightsA := a.Weights()
	weightsB := a.Weights()

	// Child
	weightsC := [3]*mat.Dense{}
	for index := range weightsA {
		weightsC[index] = mat.DenseCopyOf(weightsA[index])
		// replace random weights by factor k
		weightsC[index].Apply(func(i, j int, v float64) float64 {
			if rand.Float64() > k {
				return weightsB[index].At(i, j)
			}
			return v
		}, weightsC[index])
	}

	return NewNeuralBrain(NeuralBrainOptions{
		InitialWeights: &weightsC,
	})
}

// Mutate tweaks n weights of the NeuralBrain by a random number between m and -m
func (a *NeuralBrain) Mutate(n int, m float64) {
	for nth := 0; nth < n; nth++ {
		layer := a.weights[rand.Intn(len(a.weights))]
		rows, cols := layer.Dims()
		r, c := rand.Intn(rows), rand.Intn(cols)
		oldV := layer.At(r, c)
		newV := oldV + m*(rand.Float64()*2-1)
		newV = math.Max(-1, math.Min(newV, 1))
		layer.Set(r, c, newV)
	}
}

func (n *NeuralBrain) Weights() [Layers]*mat.Dense {
	return n.weights
}

func (n *NeuralBrain) Print() {
	for i := range n.weights {
		fmt.Println(n.weights[i])
	}
}

type NeuralBrainOptions struct {
	InitialWeights *[Layers]*mat.Dense // if nil, the weights will be initialized with random values between 0 and 1
}

func NewNeuralBrain(options NeuralBrainOptions) *NeuralBrain {
	var weights [Layers]*mat.Dense

	if options.InitialWeights == nil {
		weights = [Layers]*mat.Dense{
			newDenseWithRandomValues(5, 4),
			newDenseWithRandomValues(4, 6),
			newDenseWithRandomValues(2, 5),
		}
	} else {
		weights = *options.InitialWeights
	}

	return &NeuralBrain{weights}
}

func newDenseWithRandomValues(r int, c int) *mat.Dense {
	data := make([]float64, r*c)
	for i := 0; i < r*c; i++ {
		data[i] = rand.Float64()*2 - 1
	}
	return mat.NewDense(r, c, data)
}
