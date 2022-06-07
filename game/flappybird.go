package game

import (
	"Gooo/bird"
	"Gooo/customMath"
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
	"math/rand"
	"time"
)

type Options struct {
	Birds           []bird.Bird
	Generation      int
	TopScore        int
	FramesPerSecond float64
}

func FlappyBirdWithCui(opt Options) error {

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(generateGameFunction(opt))

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}
	// never
	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("hello", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		_, err := fmt.Fprintln(v, "Hello world!")
		if err != nil {
			panic(err)
		}
	}
	return nil
}

const PIPE_WIDTH = 0.05
const PIPE_HOLE_SIZE = 0.4
const PIPE_SPEED = 0.01
const PIPE_MARGIN = 0.2

func generateGameFunction(opt Options) func(g *gocui.Gui) error {
	return func(g *gocui.Gui) error {
		//var width, height = g.Size()

		debugV, err := g.SetView("debug", 0, 0, 100, 2)
		if err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
		}

		if v, err := g.SetView("flappy", 0, 3, 100, 20); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}

			v.Overwrite = true
			v.FgColor = gocui.ColorMagenta

			//var vWidth, vHeight = v.Size()

			var birds = opt.Birds
			var pipes = []customMath.Point{randomPipe(0.5), randomPipe(1)}
			fps := opt.FramesPerSecond
			ticker := timerFromFps(fps)
			go func() {
				totalScore := 0
				frame := 0
				for range ticker.C {

					allDead := true
					g.Update(func(g *gocui.Gui) error {
						// move pipes
						for i := range pipes {
							pipe := &pipes[i]
							pipe.X = pipe.X - PIPE_SPEED

							if pipe.X+PIPE_WIDTH <= 0 {
								pipes[i] = randomPipe(1)
							}
						}
						// move birds
						for i := range birds {
							bird := &birds[i]
							if bird.Dead() {
								continue
							}
							allDead = false
							var distanceForward = pipes[0].X
							var nextPipe = pipes[0]
							for _, pipe := range pipes {
								if pipe.X < distanceForward && pipe.X > bird.Position().X {
									distanceForward = pipe.X
									nextPipe = pipe
								}
							}
							bird.DoSomething(distanceForward, nextPipe.Y)
							bird.DoPhysics()
							var pos = bird.Position()

							if pos.Y < 0 || pos.Y > 1 {
								bird.SetDead(true)
								break
							}
							for _, pipe := range pipes {
								if pos.X >= pipe.X && pos.X <= pipe.X+PIPE_WIDTH {
									if pos.Y >= pipe.Y || pos.Y <= pipe.Y-PIPE_HOLE_SIZE {
										bird.SetDead(true)
										break
									}
								}
							}
							bird.IncrementScore(1)
						}

						v, err := g.View("flappy")
						if err != nil {
							panic(err)
						}
						v.Clear()
						var _, height = v.Size()

						// Drawing pipes
						for _, pipe := range pipes {
							err := drawBlock(pipe.X, pipe.X+PIPE_WIDTH, pipe.Y+1/float64(height), 1, v)
							err = drawBlock(pipe.X, pipe.X+PIPE_WIDTH, 0, pipe.Y-PIPE_HOLE_SIZE+1/float64(height), v)
							if err != nil {
								return err
							}
						}

						// drawing points
						var birdsPerY = make([]int, height)
						var xScreen = 0
						for _, bird := range birds {
							if bird.Dead() {
								continue
							}
							x, y := bird.Position().Spread()
							xNew, yScreen := Pos(x, y, v)
							xScreen = xNew
							if yScreen < len(birdsPerY) {
								birdsPerY[yScreen] += 1
							}
						}
						for yScreen, count := range birdsPerY {
							err := v.SetCursor(xScreen, yScreen)
							if err != nil {
								return err
							}
							switch count {
							case 0:
								continue
							case 1:
								v.EditWrite('░')
							case 2:
								v.EditWrite('▒')
							case 3:
								v.EditWrite('▓')
							default:
								v.EditWrite('█')
							}
						}

						// Debug
						debugV.Clear()
						birdCount := 0
						for i := range birds {
							if !birds[i].Dead() {
								birdCount++
							}
						}
						topScore := opt.TopScore
						if totalScore > topScore {
							topScore = totalScore
						}

						_, err = debugV.Write([]byte(fmt.Sprintf("Generation: %d | FPS: %d | Birds: %d | Top Score: %d | Score: %d", opt.Generation, int(fps), birdCount, topScore, totalScore)))
						if err != nil {
							return err
						}

						return nil
					})
					// wait for ui updating
					time.Sleep(100 * time.Microsecond)

					if allDead {
						g.Update(func(gui *gocui.Gui) error {
							return gocui.ErrQuit
						})
						break
					}
					totalScore += 1
					frame++
					if frame%50 == 0 {
						fps += 1
						ticker.Reset(fpsToDuration(fps))
					}
				}
			}()
		}
		return nil
	}
}

func timerFromFps(fps float64) *time.Ticker {
	return time.NewTicker(fpsToDuration(fps))
}

func fpsToDuration(fps float64) time.Duration {
	return time.Duration(float64(time.Second) / fps)
}

func randomPipe(x float64) customMath.Point {
	return customMath.Point{X: x, Y: rand.Float64()*(1-2*PIPE_MARGIN) + PIPE_MARGIN + PIPE_HOLE_SIZE/2}
}

// the function draws a block both horizontal and vertical
func drawBlock(x1, x2, y1, y2 float64, v *gocui.View) error {
	xx1, yy1 := Pos(x1, y1, v)
	xx2, yy2 := Pos(x2, y2, v)
	width, height := v.Size()
	for y := yy1; y <= yy2; y++ {
		for x := xx1; x <= xx2; x++ {
			if x < 0 || x >= width-1 || y < 0 || y >= height {
				continue
			}
			err := v.SetCursor(x, y)
			if err != nil {
				return err
			}
			v.EditWrite('#')
		}
	}
	return nil
}

func Pos(x float64, y float64, v *gocui.View) (int, int) {
	w, h := v.Size()
	return int(x * float64(w)), int(y * float64(h))
}
