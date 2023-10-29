package value

import (
	"github.com/aldernero/gaul"
	"github.com/aldernero/interp"
	"math"
)

const (
	DefaultOctaves       = 4
	DefaultKnots         = 3
	DefaultMinX          = 0
	DefaultMaxX          = 1
	DefaultMinY          = -1
	DefaultMaxY          = 1
	SplineOverflowFactor = 0.9
)

type Noise1D struct {
	Seed       uint64
	Octaves    int
	Knots      int
	MinX, MaxX float64
	MinY, MaxY float64
	rng        gaul.LFSRLarge
	maxAmp     float64
	splines    []interp.CubicSpline
}

func NewNoise1D(seed uint64) *Noise1D {
	return &Noise1D{
		Seed:    seed,
		Octaves: DefaultOctaves,
		Knots:   DefaultKnots,
		MinX:    DefaultMinX,
		MaxX:    DefaultMaxX,
		MinY:    DefaultMinY,
		MaxY:    DefaultMaxY,
	}
}

func (n *Noise1D) init() {
	n.rng = gaul.NewLFSRLargeWithSeed(n.Seed)
	n.splines = []interp.CubicSpline{}
	amplitude := SplineOverflowFactor
	knots := n.Knots
	for i := 0; i < n.Octaves; i++ {
		p := math.Pow(2, float64(i))
		n.maxAmp += amplitude
		for j := 0; j < int(p); j++ {
			xs := gaul.Linspace(n.MinX, n.MaxX, knots+2, true)
			ys := make([]float64, knots+2)
			for k := 0; k < knots+2; k++ {
				ys[k] = amplitude * gaul.Map(0, 1, -1, 1, n.rng.Float64()) / p
			}
			spline, _ := interp.NewCubicSpline(xs, ys)
			n.splines = append(n.splines, spline)
		}
		knots *= 2
		amplitude /= 2
	}
}

func (n *Noise1D) Eval(x float64) float64 {
	y := 0.0
	for _, spline := range n.splines {
		y += spline.Eval(x)
	}
	return gaul.Map(-n.maxAmp, n.maxAmp, n.MinY, n.MaxY, y)
}
