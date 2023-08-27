package opensimplex

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

var seeds = []int64{46, 38, 36, 40, 24, 52, 62, 49, 24, 55, math.MinInt64, math.MaxInt64}

var testPoints = [][]float64{
	{0.0, 0.0, 0.0, 0.0},
	{-122.23991945003587, -428.2528315542171, 489.37254518248164, 178.01174302724132},
	{184.41053746024738, 181.47471868357968, -400.2370541222765, -255.5953030302438},
	{122.37870493623137, -483.58106343488316, 148.4211926815907, -241.5396024919716},
	{-51.182468665792236, -484.65915775747794, -444.27837949390425, 327.671118599802},
	{-210.6314940469831, 134.8506313635961, 17.323852121434168, 285.8475716488992},
	{-264.78759963401956, 28.249052018116004, -223.6253979058257, 297.84780917467225},
	{33.55885748465681, -289.85726810789737, -482.58876139686237, 453.24757083916444},
	{286.3472640611923, -206.79124797231285, -70.33062967132109, -55.21955298620207},
	{350.5216010614357, 82.43443085499614, 142.68103352577754, 298.0398725703752},
	{104.14509776922986, 448.32829690915634, 371.21299001603927, 257.7996117494559},
	{-53.08418866599429, 398.8432023628333, -450.3075537164918, -47.08765074902854},
	{406.47815073104346, 350.06724506477184, -70.32882019517628, 20.61332892457479},
	{481.04685061554875, 136.162946309749, -114.80035764256336, -11.74252922621777},
	{273.79327682304364, -145.730339040963, -301.41129961450275, -247.2008611158578},
	{-67.2845876062883, 308.2655534340212, 364.0027263741823, 42.98928461571872},
	{-131.10872698248977, -250.01732391489872, 290.5198394095515, 308.9774649204279},
	{-273.8828247718006, -415.2073940611477, 421.3633385847502, 204.23686749655155},
	{-438.4126186155904, 298.6136665587418, -166.74214332356064, -97.57623425824814},
	{209.65011141985113, -286.0529615439887, 175.22853233110104, 331.9087430963689},
	{math.MaxFloat64, -math.MaxFloat64, -math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64},
}

func TestNoise2D(t *testing.T) {
	for _, seed := range seeds {
		noise := NewOpenSimplex2(seed)
		for _, point := range testPoints {
			x := point[0]
			y := point[1]
			val := noise.Noise2D(x, y)
			assert.GreaterOrEqual(t, val, 0.0)
			assert.LessOrEqual(t, val, 1.0)
			val = noise.Noise2DImproveX(x, y)
			assert.GreaterOrEqual(t, val, 0.0)
			assert.LessOrEqual(t, val, 1.0)
		}
	}
}

func TestNoise3D(t *testing.T) {
	for _, seed := range seeds {
		noise := NewOpenSimplex2(seed)
		for _, point := range testPoints {
			x := point[0]
			y := point[1]
			z := point[2]
			val := noise.Noise3D(x, y, z)
			assert.GreaterOrEqual(t, val, 0.0)
			assert.LessOrEqual(t, val, 1.0)
			val = noise.Noise3DImproveXY(x, y, z)
			assert.GreaterOrEqual(t, val, 0.0)
			assert.LessOrEqual(t, val, 0.0)
			val = noise.Noise3DImproveXZ(x, y, z)
			assert.GreaterOrEqual(t, val, 0.0)
			assert.LessOrEqual(t, val, 1.0)
			val = noise.Noise3DFallback(x, y, z)
			assert.GreaterOrEqual(t, val, 0.0)
			assert.LessOrEqual(t, val, 1.0)
		}
	}
}

func TestNoise4D(t *testing.T) {
	for _, seed := range seeds {
		noise := NewOpenSimplex2(seed)
		for _, point := range testPoints {
			x := point[0]
			y := point[1]
			z := point[2]
			w := point[3]
			val := noise.Noise4D(x, y, z, w)
			assert.GreaterOrEqual(t, val, 0.0)
			assert.LessOrEqual(t, val, 1.0)
			val = noise.Noise4DImproveXYZImproveXY(x, y, z, w)
			assert.GreaterOrEqual(t, val, 0.0)
			assert.LessOrEqual(t, val, 1.0)
			val = noise.Noise4DImproveXYZImproveXZ(x, y, z, w)
			assert.GreaterOrEqual(t, val, 0.0)
			assert.LessOrEqual(t, val, 1.0)
			val = noise.Noise4DImproveXYZ(x, y, z, w)
			assert.GreaterOrEqual(t, val, 0.0)
			assert.LessOrEqual(t, val, 1.0)
			val = noise.Noise4DImproveXZImproveZW(x, y, z, w)
			assert.GreaterOrEqual(t, val, 0.0)
			assert.LessOrEqual(t, val, 1.0)
			val = noise.Noise4DFallback(x, y, z, w)
			assert.GreaterOrEqual(t, val, 0.0)
			assert.LessOrEqual(t, val, 1.0)
		}
	}
}

func BenchmarkOpenSimplex2_Noise2D(b *testing.B) {
	noise := NewOpenSimplex2(42)
	for i := 0; i < b.N; i++ {
		for _, point := range testPoints {
			noise.Noise2D(point[0], point[1])
		}
	}
}

func BenchmarkOpenSimplex2_Noise3D(b *testing.B) {
	noise := NewOpenSimplex2(42)
	for i := 0; i < b.N; i++ {
		for _, point := range testPoints {
			noise.Noise3D(point[0], point[1], point[2])
		}
	}
}

func BenchmarkOpenSimplex2_Noise4D(b *testing.B) {
	noise := NewOpenSimplex2(42)
	for i := 0; i < b.N; i++ {
		for _, point := range testPoints {
			noise.Noise4D(point[0], point[1], point[2], point[3])
		}
	}
}
