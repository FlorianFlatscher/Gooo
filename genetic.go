// # Genetic Algorithm
// This genetic algorithm trains FlappyBird AI agents to play FlappyBird. The learning phase is visualized using CUI, a packet for creating a CLI. The neural network of the agents is implemented using the gonum matrix library.

// ## Idea
// 1. The population, consisting of 30 birds, play FlappyBirds.
// 2. Each bird based on how long it was able to survive.
// 3. The best ones are elected to be the base of the new population.
// 4. The cycle repeats with the new population

//
package main

import (
	"GeneticAlgorithm/bird"     // covers the birds mechanics
	"GeneticAlgorithm/decision" // covers deciding on certain action based on observation
	"GeneticAlgorithm/game"     // covers the FlappyBird game, including CLI visualisation
	"math"
	"math/rand"
)

// Each game is played by 30 birds at the same time.
const populationSize = 30

// ## Algorithm

func main() {
	// The first generation's population consists of birds that are not trained yet.
	population := BuildPopulation()
	topScore := 0.0

	for generation := 0; true; generation++ {
		// FlappyBirdWithCui sets up a CLI FlappyBird game, with the bird agents as players. The game is over when al birds are dead. For details on the implementation see [game/flappybird.go](../game/docs/flappybird.html).
		game.FlappyBirdWithCui(game.Options{
			Birds:           population,
			Generation:      generation,
			TopScore:        int(topScore),
			FramesPerSecond: float64(generation/4)*20 + 30,
		})

		bestBird := CalculateBestBird(population)
		if bestBird.Score() > topScore {
			topScore = bestBird.Score()
		}

		// The mutation rate describes random modifications to the new generation. If it is higher, one could say the new generation is exploring new approaches. In this implementation, the mutation rate is decreasing, as the top score of the birds is increasing.
		mutation := math.Max(0.5-0.0001*float64(topScore), 0.001)

		// The new population is build through evolution. This marks the start of a new generation.
		population = Evolve(population, mutation)
	}
}

// ## Birds

// BuildPopulation builds a new population of birds with a neural network as decision engine. For more details on the implementation see [decision/neural.go](../decision/docs/neural.html).
func BuildPopulation() []bird.Bird {
	birds := make([]bird.Bird, populationSize)
	for i := 0; i < cap(birds); i++ {
		birds[i] = *bird.NewBird(bird.Options{
			DecisionEngine: decision.NewNeuralBrain(),
		})
	}
	return birds
}

// CalculateSummedScore calculates the score for the entire population together.
func CalculateSummedScore(birds []bird.Bird) float64 {
	summedScore := 0.0
	for i := 0; i < len(birds); i++ {
		summedScore += birds[i].Score()
	}
	return summedScore
}

// CalculateBestBird calculates the best scoring bird of the population.
func CalculateBestBird(birds []bird.Bird) bird.Bird {
	bestBird := birds[0]
	for i := 1; i < len(birds); i++ {
		if birds[i].Score() > bestBird.Score() {
			bestBird = birds[i]
		}
	}
	return bestBird
}

// ## Evolution

// Evolve builds the next generation's population, by reproducing and mutating the best scoring birds of the last generation.
func Evolve(birds []bird.Bird, mutation float64) []bird.Bird {
	summedScore := CalculateSummedScore(birds)
	bestBird := CalculateBestBird(birds)

	// A population with the same size is initialized.
	var children = make([]bird.Bird, populationSize)

	// The best scoring bird is revived and reused.
	children[0] = bestBird.Revive()

	// All other birds are "born by natural election".
	for i := 1; i < cap(children); i++ {
		// Two parents are elected by random. Better scoring birds are more likely to be elected.
		a, b := ElectBothParents(birds, summedScore)
		// Taking the parent's brains (the neural networks).
		brainA, brainB := a.DecisionEngine().(*decision.NeuralBrain), b.DecisionEngine().(*decision.NeuralBrain)
		// k is the cross-over coefficient that describes if the child should inherit more from a (k <= 0.5) or from b (k > 0.5).
		var k float64
		if a.Score() > b.Score() {
			k = 0.9
		}
		if b.Score() > a.Score() {
			k = 0.1
		}
		// The child's brain is generated.
		brainC := decision.CrossOver(brainA, brainB, k)
		// 5 weights of the neural network are modified by the mutation rate given as parameter.
		brainC.Mutate(5, mutation)
		children[i] = *bird.NewBird(bird.Options{
			DecisionEngine: brainC,
		})
	}

	return children
}

// ElectBothParents elects two parents
func ElectBothParents(birds []bird.Bird, totalScore float64) (bird.Bird, bird.Bird) {
	return ElectOneParent(birds, totalScore), ElectOneParent(birds, totalScore)
}

// ElectOneParent elects a parent. Higher scoring birds are more likely to be elected.
func ElectOneParent(birds []bird.Bird, totalScore float64) bird.Bird {
	pivot := rand.Float64() * totalScore
	for i := range birds {
		pivot -= birds[i].Score()
		if pivot <= 0 {
			return birds[i]
		}
	}
	return birds[len(birds)-1]
}
