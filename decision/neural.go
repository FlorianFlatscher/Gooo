package decision

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
	"math"
	"math/rand"
)

const Layers = 3

const Input = 3
const Layer1 = 3
const Layer2 = 3
const Layer3 = 2

// NeuralBrain is an implementation of the bird.DecisionEngine interface.
// It uses a neural network to decide based on certain observation.
type NeuralBrain struct {
	weights [Layers]*mat.Dense
}

func (nb *NeuralBrain) DecideOnObservation(obs Observation) Action {
	input := mat.NewDense(3, 1, []float64{
		obs.DistanceForward,
		obs.HeightOfHole,
		obs.Position.Y,
	})

	curr := input
	for _, layer := range nb.weights {
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
		return ActionJump
	}
	return ActionNothing
}

func (nb *NeuralBrain) DecideOnEvent(obs Observation) Action {
	return ActionNothing
}

// CrossOver constructs a new NeuralBrain that uses k*100% weights of a, and (1-k)*100% weights of b, where k is between 0 and 1
func CrossOver(a *NeuralBrain, b *NeuralBrain, k float64) *NeuralBrain {
	if k < 0 || k > 1 {
		panic("invalid argument k")
	}

	// Parents
	weightsA := a.Weights()
	weightsB := b.Weights()

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

	return NewNeuralBrainO(NeuralBrainOptions{
		InitialWeights: &weightsC,
	})
}

// Mutate tweaks n weights of the NeuralBrain by a random number between m and -m
func (nb *NeuralBrain) Mutate(n int, m float64) {
	for nth := 0; nth < n; nth++ {
		layer := nb.weights[rand.Intn(len(nb.weights))]
		rows, cols := layer.Dims()
		r, c := rand.Intn(rows), rand.Intn(cols)
		oldV := layer.At(r, c)
		newV := oldV + m*(rand.Float64()*2-1)
		newV = math.Max(-1, math.Min(newV, 1))
		layer.Set(r, c, newV)
	}
}

func (nb *NeuralBrain) Weights() [Layers]*mat.Dense {
	return nb.weights
}

func (nb *NeuralBrain) Print() {
	for i := range nb.weights {
		fmt.Println(nb.weights[i])
	}
}

type NeuralBrainOptions struct {
	InitialWeights *[Layers]*mat.Dense // if nil, the weights will be initialized with random values between 0 and 1
}

func NewNeuralBrainO(options NeuralBrainOptions) *NeuralBrain {
	var weights [Layers]*mat.Dense

	if options.InitialWeights == nil {
		// +1 for bias
		weights = [Layers]*mat.Dense{
			newDenseWithRandomValues(Layer1, Input+1),
			newDenseWithRandomValues(Layer2, Layer1+1),
			newDenseWithRandomValues(Layer3, Layer2+1),
		}
	} else {
		weights = *options.InitialWeights
	}

	return &NeuralBrain{weights}
}

func NewNeuralBrain() *NeuralBrain {
	return NewNeuralBrainO(NeuralBrainOptions{
		InitialWeights: nil,
	})
}

func newDenseWithRandomValues(r int, c int) *mat.Dense {
	data := make([]float64, r*c)
	for i := 0; i < r*c; i++ {
		data[i] = rand.Float64()*2 - 1
	}
	return mat.NewDense(r, c, data)
}
