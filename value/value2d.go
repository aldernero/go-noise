package value

import (
	"github.com/aldernero/gaul"
	"github.com/aldernero/interp"
)

type Noise2D struct {
	Seed       uint64
	Octaves    int
	MinX, MaxX float64
	MinY, MaxY float64
	rng        gaul.LFSRLarge
	xKnots     []float64
	yKnots     []float64
	hNoise     []*Noise1D
	vNoise     []*Noise1D
}

func NewNoise2D(seed uint64, xKnots, yKnots, octaves int, rect gaul.Rect) *Noise2D {
	noise := Noise2D{
		Seed:    seed,
		Octaves: octaves,
		MinX:    rect.X,
		MaxX:    rect.X + rect.W,
		MinY:    rect.Y,
		MaxY:    rect.Y + rect.H,
	}
	noise.rng = gaul.NewLFSRLargeWithSeed(noise.Seed)
	noise.xKnots = gaul.Linspace(noise.MinX, noise.MaxX, xKnots, true)
	noise.yKnots = gaul.Linspace(noise.MinY, noise.MaxY, yKnots, true)
	noise.hNoise = make([]*Noise1D, xKnots)
	noise.vNoise = make([]*Noise1D, yKnots)
	for i := range noise.xKnots {
		noise.vNoise[i] = NewNoise1D(noise.Seed + uint64(i))
		noise.vNoise[i].Octaves = noise.Octaves
		noise.vNoise[i].MinX = noise.MinY
		noise.vNoise[i].MaxX = noise.MaxY
		noise.vNoise[i].Init()
	}
	for i := range noise.yKnots {
		noise.hNoise[i] = NewNoise1D(noise.Seed + 42*uint64(i))
		noise.hNoise[i].Octaves = noise.Octaves
		noise.hNoise[i].MinX = noise.MinX
		noise.hNoise[i].MaxX = noise.MaxX
		noise.hNoise[i].Init()
	}
	return &noise
}

func (n *Noise2D) Eval(x, y float64) float64 {
	xIndex := interp.Bisect(n.xKnots, x)
	yIndex := interp.Bisect(n.yKnots, y)
	xNoise := n.hNoise[xIndex].Eval(x)
	yNoise := n.vNoise[yIndex].Eval(y)
	return 0.5 * (xNoise + yNoise)
}
