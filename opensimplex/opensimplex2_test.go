package opensimplex

import (
	"github.com/stretchr/testify/assert"
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
