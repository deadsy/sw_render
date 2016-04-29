package vec

import (
	"math"
)

type V3f [3]float32

// Return the Euclidean length of a
func (a V3f) Length() float32 {
	return float32(math.Sqrt(float64(a[0]*a[0] + a[1]*a[1] + a[2]*a[2])))
}

// Return a * k
func (a V3f) Scale(k float32) V3f {
	return V3f{
		a[0] * k,
		a[1] * k,
		a[2] * k,
	}
}

// Return a + b
func (a V3f) Sum(b V3f) V3f {
	return V3f{
		a[0] + b[0],
		a[1] + b[1],
		a[2] + b[2],
	}
}

// Return a - b
func (a V3f) Sub(b V3f) V3f {
	return V3f{
		a[0] - b[0],
		a[1] - b[1],
		a[2] - b[2],
	}
}

// Return a x b
func (a V3f) Cross(b V3f) V3f {
	return V3f{
		a[1]*b[2] - a[2]*b[1],
		a[2]*b[0] - a[0]*b[2],
		a[0]*b[1] - a[1]*b[0],
	}
}

// Return a.b
func (a V3f) Dot(b V3f) float32 {
	return (a[0] * b[0]) +
		(a[1] * b[1]) +
		(a[2] * b[2])
}

// Normalize a
func (a V3f) Normalize() V3f {
	l := a.Length()
	if l == 0 {
		return a
	} else {
		return V3f{
			a[0] / l,
			a[1] / l,
			a[2] / l,
		}
	}
}
