//

//
package main

import (
	"Gooo/agent"
	"Gooo/brain"
	"Gooo/cui"
	"math"
	"math/rand"
)

const birdsCount = 30

func main() {
	var birds = make([]agent.Bird, birdsCount)
	for i := 0; i < cap(birds); i++ {
		b := brain.NewNeuralBrain(brain.NeuralBrainOptions{
			InitialWeights: nil,
		})
		birds[i] = *agent.NewBird(b)
	}
	generation := 0
	topScore := 0
	for {
		err := cui.LaunchGame(cui.GameOptions{
			Birds:           birds,
			Generation:      generation,
			TopScore:        topScore,
			FramesPerSecond: float64(generation/4)*20 + 30,
		})
		if err != nil {
			panic(err)
		}

		bestBird := birds[0]
		totalScore := bestBird.Score()
		for i := 1; i < len(birds); i++ {
			totalScore += birds[i].Score()
			if birds[i].Score() > bestBird.Score() {
				bestBird = birds[i]
			}
		}

		if int(bestBird.Score()) > topScore {
			topScore = int(bestBird.Score())
		}

		var children = make([]agent.Bird, birdsCount)
		// Keep the bird with the best score as it is
		children[0] = bestBird.Revive()
		for i := 1; i < cap(children); i++ {
			a, b := electBothParents(birds, totalScore)
			brainA, brainB := a.Brain().(*brain.NeuralBrain), b.Brain().(*brain.NeuralBrain)
			var k float64
			if a.Score() > b.Score() {
				k = 0.6
			}
			if b.Score() > a.Score() {
				k = 0.4
			}
			brainC := brain.CrossOver(brainA, brainB, k)
			brainC.Mutate(5, math.Max(0.2-0.00001*float64(topScore), 0.001))
			children[i] = *agent.NewBird(brainC)
		}
		copy(birds, children)
		generation++
	}
}

func electBothParents(birds []agent.Bird, totalScore float64) (agent.Bird, agent.Bird) {
	return electOneParent(birds, totalScore), electOneParent(birds, totalScore)
}

func electOneParent(birds []agent.Bird, totalScore float64) agent.Bird {
	pivot := rand.Float64() * totalScore
	for i := range birds {
		pivot -= birds[i].Score()
		if pivot <= 0 {
			return birds[i]
		}
	}
	return birds[len(birds)-1]
}
