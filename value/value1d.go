package value

import (
	"github.com/aldernero/gaul"
	"github.com/aldernero/interp"
)

const (
	DefaultOctaves       = 4
	DefaultKnots         = 3
	DefaultMinX          = 0
	DefaultMaxX          = 1
	SplineOverflowFactor = 0.9
)

var powerOfTwo func(int) int = func(i int) int {
	return 1 << i
}

type Noise1D struct {
	Seed            uint64
	Octaves         int
	Knots           int
	MinX, MaxX      float64
	SplineCountFunc func(int) float64
	KnotCountFunc   func(int) int
	AmpScaleFunc    func(int) float64
	rng             gaul.LFSRLarge
	maxAmp          float64
	splines         []interp.CubicSpline
}

func NewNoise1D(seed uint64) *Noise1D {
	noise := Noise1D{
		Seed:    seed,
		Octaves: DefaultOctaves,
		Knots:   DefaultKnots,
		MinX:    DefaultMinX,
		MaxX:    DefaultMaxX,
		SplineCountFunc: func(i int) float64 {
			return float64(powerOfTwo(i))
		},
		KnotCountFunc: func(i int) int {
			return powerOfTwo(i)
		},
		AmpScaleFunc: func(i int) float64 {
			return SplineOverflowFactor / float64(powerOfTwo(i))
		},
	}
	return &noise
}

func (n *Noise1D) Init() {
	n.rng = gaul.NewLFSRLargeWithSeed(n.Seed)
	n.splines = []interp.CubicSpline{}
	for i := 0; i < n.Octaves; i++ {
		p := n.SplineCountFunc(i)
		knots := n.KnotCountFunc(i)
		amplitude := n.AmpScaleFunc(i)
		n.maxAmp += amplitude * float64(p)
		for j := 0; j < int(p); j++ {
			xs := gaul.Linspace(n.MinX, n.MaxX, knots+2, true)
			ys := make([]float64, knots+2)
			for k := 0; k < knots+2; k++ {
				ys[k] = amplitude * gaul.Map(0, 1, -1, 1, n.rng.Float64())
			}
			spline, _ := interp.NewCubicSpline(xs, ys)
			n.splines = append(n.splines, spline)
		}
	}
}

func (n *Noise1D) Eval(x float64) float64 {
	y := 0.0
	for _, spline := range n.splines {
		y += spline.Eval(x)
	}
	return gaul.Map(-n.maxAmp, n.maxAmp, 0, 1, y)
}

func (n *Noise1D) EvalSigned(x float64) float64 {
	y := 0.0
	for _, spline := range n.splines {
		y += spline.Eval(x)
	}
	return gaul.Map(-n.maxAmp, n.maxAmp, -1, 1, y)
}
