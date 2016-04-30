package vec

import (
	"math/rand"
)

type V2i [2]int

// return true if the vectors are equal
func (a V2i) Equal(b V2i) bool {
	return (a[0] == b[0]) && (a[1] == b[1])
}

// Return a - b
func (a V2i) Sub(b V2i) V2i {
	return V2i{
		a[0] - b[0],
		a[1] - b[1],
	}
}

// sort points by Y
func Sort_Y(p []*V2i) {
	if p[0][1] > p[1][1] {
		// swap p[0] with p[1]
		x := p[1]
		p[1] = p[0]
		p[0] = x
	}
}

// sort points by X
func Sort_X(p []*V2i) {
	if p[0][0] > p[1][0] {
		// swap p[0] with p[1]
		x := p[1]
		p[1] = p[0]
		p[0] = x
	}
}

// return a random V2i - limits set by passed vector
func (a V2i) Rand() V2i {
	return V2i{
		int(float32(a[0]) * rand.Float32()),
		int(float32(a[1]) * rand.Float32()),
	}
}

// return a random V2i - offset from a
func (a V2i) Rand_Delta(d int) V2i {
	return V2i{
		a[0] + int(float32(d)*(rand.Float32()-0.5)),
		a[1] + int(float32(d)*(rand.Float32()-0.5)),
	}
}
