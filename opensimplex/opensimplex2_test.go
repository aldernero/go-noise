package opensimplex

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestNoise2D(t *testing.T) {
	seeds := []int64{1, 42, 0, -43}
	for _, seed := range seeds {
		noise := NewOpenSimplex2(seed)
		points := [][]float64{
			{0.0, 0.0},
			{1.0, 0.0},
			{0.0, 1.0},
			{1.0, 1.0},
			{-1000.3, 0.57},
		}
		for _, point := range points {
			val := noise.Noise2D(point[0], point[1])
			assert.GreaterOrEqual(t, val, -1.0)
			assert.LessOrEqual(t, val, 1.0)
		}
	}
}

func BenchmarkOpenSimplex2_Noise2D(b *testing.B) {
	randSeed := int64(42)
	noiseSeed := int64(1337)
	noise := NewOpenSimplex2(noiseSeed)
	rng := rand.New(rand.NewSource(randSeed))
	points := make([][]float64, 100)
	for i := 0; i < len(points); i++ {
		points[i] = []float64{1080 * rng.Float64(), 1080 * rng.Float64()}
	}
	for i := 0; i < b.N; i++ {
		for _, point := range points {
			noise.Noise2D(point[0], point[1])
		}
	}
}
