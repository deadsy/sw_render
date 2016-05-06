package vec

import (
//"math"
)

type V2 [2]float32

// Return a - b
func (a V2) Sub(b V2) V2 {
	return V2{
		a[0] - b[0],
		a[1] - b[1],
	}
}

func (a V2) Dot(b V2) float32 {
	return a[0]*b[0] + a[1]*b[1]
}
