//

//
package main

import (
	"Gooo/agent"
	"Gooo/brain"
	"Gooo/cui"
)

const birdsCount = 5

func main() {
	var birds = make([]agent.Bird, 5)
	for i := 0; i < birdsCount; i++ {
		birds[i] = *agent.NewBird(brain.NewDumbBrain())
		//birds = append(birds, *agent.NewBird(brain.NewDumbBrain()))
	}

	cui.LaunchGame(cui.GameOptions{
		Birds: birds,
	})
}
