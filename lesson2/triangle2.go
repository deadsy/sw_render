package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/deadsy/sw_render/vec"
)

// return [u,v,w]  such that p = ua + vb + wc
func barycentric1(a, b, c, p vec.V2) vec.V3 {
	ab := b.Sub(a)
	ac := c.Sub(a)
	pa := a.Sub(p)

	vx := vec.V3{ab[0], ac[0], pa[0]}
	vy := vec.V3{ab[1], ac[1], pa[1]}

	x := vx.Cross(vy)

	v := x[0] / x[2]
	w := x[1] / x[2]

	return vec.V3{1 - v - w, v, w}
}

// return [u,v,w]  such that p = ua + vb + wc
func barycentric2(a, b, c, p vec.V2) vec.V3 {

	v0 := b.Sub(a)
	v1 := c.Sub(a)
	v2 := p.Sub(a)

	d00 := v0.Dot(v0)
	d01 := v0.Dot(v1)
	d11 := v1.Dot(v1)
	d20 := v2.Dot(v0)
	d21 := v2.Dot(v1)

	denom := d00*d11 - d01*d01

	v := (d11*d20 - d01*d21) / denom
	w := (d00*d21 - d01*d20) / denom

	return vec.V3{1 - v - w, v, w}
}

func test_bs(a, b, c, p vec.V2) {
	fmt.Printf("b1 %+v\n", barycentric1(a, b, c, p))
	fmt.Printf("b2 %+v\n", barycentric2(a, b, c, p))
}

func test_barycentric() {

	k := float32(1.4142135623730950488)

	test_bs(vec.V2{0, 0}, vec.V2{2, 0}, vec.V2{0, 2}, vec.V2{1, 1})
	test_bs(vec.V2{0, 0}, vec.V2{2, 0}, vec.V2{0, 2}, vec.V2{10, 10})
	test_bs(vec.V2{0, 0}, vec.V2{5, 0}, vec.V2{5, 5}, vec.V2{2, 2})
	test_bs(vec.V2{0, 0}, vec.V2{1, 1}, vec.V2{1, 1}, vec.V2{2, 2})
	test_bs(vec.V2{1, 1}, vec.V2{1, 1}, vec.V2{1, 1}, vec.V2{2, 2})
	test_bs(vec.V2{1, 1}, vec.V2{2, 2}, vec.V2{3, 3}, vec.V2{2, 2})
	test_bs(vec.V2{0, 0}, vec.V2{2, 0}, vec.V2{0, 2}, vec.V2{k, k})
	test_bs(vec.V2{1, 1}, vec.V2{3, 1}, vec.V2{1, 3}, vec.V2{1 + k, 1 + k})
}

func triangle2(v [3]*vec.V2, img *image.NRGBA, color color.NRGBA) {

}
