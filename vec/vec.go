package vec

import (
	"math"
)

type V3 [3]float32

// Return the Euclidean length of a
func (a V3) Length() float32 {
	return float32(math.Sqrt(float64(a[0]*a[0] + a[1]*a[1] + a[2]*a[2])))
}

// Return a * k
func (a V3) Scale(k float32) V3 {
	return V3{
		a[0] * k,
		a[1] * k,
		a[2] * k,
	}
}

// Return a + b
func (a V3) Sum(b V3) V3 {
	return V3{
		a[0] + b[0],
		a[1] + b[1],
		a[2] + b[2],
	}
}

// Return a - b
func (a V3) Sub(b V3) V3 {
	return V3{
		a[0] - b[0],
		a[1] - b[1],
		a[2] - b[2],
	}
}

// Return a x b
func (a V3) Cross(b V3) V3 {
	return V3{
		a[1]*b[2] - a[2]*b[1],
		a[2]*b[0] - a[0]*b[2],
		a[0]*b[1] - a[1]*b[0],
	}
}

// Return a.b
func (a V3) Dot(b V3) float32 {
	return (a[0] * b[0]) +
		(a[1] * b[1]) +
		(a[2] * b[2])
}

// Normalize a
func (a V3) Normalize() V3 {
	l := a.Length()
	if l == 0 {
		return a
	} else {
		return V3{
			a[0] / l,
			a[1] / l,
			a[2] / l,
		}
	}
}
