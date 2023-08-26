package opensimplex

import (
	"math"
)

const (
	PRIME_X                  int64   = 0x5205402B9270C86F
	PRIME_Y                  int64   = 0x598CD327003817B5
	PRIME_Z                  int64   = 0x5BCC226E9FA0BACB
	PRIME_W                  int64   = 0x56CC5227E58F554B
	HASH_MULTIPLIER          int64   = 0x53A3F72DEEC546F5
	SEED_FLIP_3D             int64   = -0x52D547B2E96ED629
	SEED_OFFSET_4D           int64   = 0xE83DC3E0DA7164D
	ROOT2OVER2               float64 = 0.7071067811865476
	SKEW_2D                  float64 = 0.366025403784439
	UNSKEW_2D                float64 = -0.21132486540518713
	ROOT3OVER3               float64 = 0.5773502691896257
	FALLBACK_ROTATE_3D               = 2.0 / 3.0
	ROTATE_3D_ORTHOGONALIZER         = UNSKEW_2D
	SKEW_4D                  float64 = -0.138196601125011
	UNSKEW_4D                float64 = 0.309016994374947
	LATTICE_STEP_4D          float64 = 0.2
	N_GRADS_2D_EXPONENT      int64   = 7
	N_GRADS_3D_EXPONENT      int64   = 8
	N_GRADS_4D_EXPONENT      int64   = 9
	N_GRADS_2D                       = 1 << N_GRADS_2D_EXPONENT
	N_GRADS_3D                       = 1 << N_GRADS_3D_EXPONENT
	N_GRADS_4D                       = 1 << N_GRADS_4D_EXPONENT
	NORMALIZER_2D            float64 = 0.01001634121365712
	NORMALIZER_3D            float64 = 0.07969837668935331
	NORMALIZER_4D            float64 = 0.0220065933241897
	RSQUARED_2D              float64 = 0.5
	RSQUARED_3D              float64 = 0.6
	RSQUARED_4D              float64 = 0.6
)

var GRADIENTS_2D = []float64{
	0.38268343236509, 0.923879532511287,
	0.923879532511287, 0.38268343236509,
	0.923879532511287, -0.38268343236509,
	0.38268343236509, -0.923879532511287,
	-0.38268343236509, -0.923879532511287,
	-0.923879532511287, -0.38268343236509,
	-0.923879532511287, 0.38268343236509,
	-0.38268343236509, 0.923879532511287,
	0.130526192220052, 0.99144486137381,
	0.608761429008721, 0.793353340291235,
	0.793353340291235, 0.608761429008721,
	0.99144486137381, 0.130526192220051,
	0.99144486137381, -0.130526192220051,
	0.793353340291235, -0.60876142900872,
	0.608761429008721, -0.793353340291235,
	0.130526192220052, -0.99144486137381,
	-0.130526192220052, -0.99144486137381,
	-0.608761429008721, -0.793353340291235,
	-0.793353340291235, -0.608761429008721,
	-0.99144486137381, -0.130526192220052,
	-0.99144486137381, 0.130526192220051,
	-0.793353340291235, 0.608761429008721,
	-0.608761429008721, 0.793353340291235,
	-0.130526192220052, 0.99144486137381,
}

type OpenSimplex2 struct {
	seed int64
	dim  int64
}

func NewOpenSimplex2(seed int64) OpenSimplex2 {
	return OpenSimplex2{seed: seed, dim: 2}
}

func (o OpenSimplex2) Noise2D(x, y float64) float64 {
	xx := x * ROOT2OVER2
	yy := y * (ROOT2OVER2 * (1 + 2*SKEW_2D))
	return noise2UnskewedBase(o.seed, xx+yy, yy-xx)
}

func noise2UnskewedBase(seed int64, x, y float64) float64 {
	// Get base points and offsets.
	xInt, xFrac := math.Modf(x)
	yInt, yFrac := math.Modf(y)

	// Prime pre-multiplication for hash.
	xsbp := int64(xInt) * PRIME_X
	ysbp := int64(yInt) * PRIME_Y

	// Unskew
	t := (xFrac + yFrac) * UNSKEW_2D
	dx0 := xFrac + t
	dy0 := yFrac + t

	// First vertex
	var value float64
	a0 := RSQUARED_2D - dx0*dx0 - dy0*dy0
	if a0 > 0 {
		value = a0 * a0 * a0 * a0 * grad2D(seed, xsbp, ysbp, dx0, dy0)
	}

	// Second vertex
	a1 := (2*(1+2*UNSKEW_2D)*(1/UNSKEW_2D+2))*t + ((-2 * (1 + 2*UNSKEW_2D) * (1 + 2*UNSKEW_2D)) + a0)
	if a1 > 0 {
		dx := dx0 - (1 + 2*UNSKEW_2D)
		dy := dy0 - (1 + 2*UNSKEW_2D)
		value += a1 * a1 * a1 * a1 * grad2D(seed, xsbp+PRIME_X, ysbp+PRIME_Y, dx, dy)
	}

	// Third vertex
	if dy0 > dx0 {
		dx2 := dx0 - UNSKEW_2D
		dy2 := dy0 - UNSKEW_2D + 1
		a2 := RSQUARED_2D - dx2*dx2 - dy2*dy2
		if a2 > 0 {
			value += a2 * a2 * a2 * a2 * grad2D(seed, xsbp, ysbp+PRIME_Y, dx2, dy2)
		}
	} else {
		dx2 := dx0 - UNSKEW_2D + 1
		dy2 := dy0 - UNSKEW_2D
		a2 := RSQUARED_2D - dx2*dx2 - dy2*dy2
		if a2 > 0 {
			value += a2 * a2 * a2 * a2 * grad2D(seed, xsbp+PRIME_X, ysbp, dx2, dy2)
		}
	}
	return value
}

func grad2D(seed int64, xsvp int64, ysvp int64, dx0 float64, dy0 float64) float64 {
	hash := seed ^ xsvp ^ ysvp
	hash *= HASH_MULTIPLIER
	hash ^= hash >> (64 - N_GRADS_2D_EXPONENT + 1)
	gi := hash & ((N_GRADS_2D - 1) << 1) % 48
	return GRADIENTS_2D[gi|0]*dx0 + GRADIENTS_2D[gi|1]*dy0
}
