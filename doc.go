// This file is for go doc, not for literate-programming

// Genetic Algorithm
//
// This genetic algorithm trains FlappyBird agents to play FlappyBird. The learning phase is visualized using CUI, a packet for creating a CLI. The neural network of the agents is implemented using the gonum matrix library.
//
// Idea
//
//		1. The population, consisting of 30 birds, play FlappyBirds.
//		2. Each bird based on how long it was able to survive.
//		3. The best ones are elected to be the base of the new population.
//		4. The cycle repeats with the new population
package main
