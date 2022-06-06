package customMath

type Point struct {
	X float64
	Y float64
}

func (p Point) Spread() (float64, float64) {
	return p.X, p.Y
}
