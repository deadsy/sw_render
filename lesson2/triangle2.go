package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/deadsy/sw_render/vec"
)

// return [u,v] such that p = a + u(b-a) + v(c-a)
func barycentric(a, b, c, p vec.V2) vec.V3 {
	ab := b.Sub(a)
	ac := c.Sub(a)
	pa := a.Sub(p)

	vx := vec.V3{ab[0], ac[0], pa[0]}
	vy := vec.V3{ab[1], ac[1], pa[1]}

	x := vx.Cross(vy)

	if x[2] == 0 {
		fmt.Printf("div by zero\n")
	}

	u := x[0] / x[2]
	v := x[1] / x[2]

	return vec.V3{1 - u - v, u, v}
}


func test_barycentric() {

	k := float32(1.4142135623730950488)

	x := barycentric(vec.V2{0, 0}, vec.V2{2, 0}, vec.V2{0, 2}, vec.V2{1, 1})
	fmt.Printf("%+v\n", x)
	x = barycentric(vec.V2{0, 0}, vec.V2{2, 0}, vec.V2{0, 2}, vec.V2{10, 10})
	fmt.Printf("%+v\n", x)
	x = barycentric(vec.V2{0, 0}, vec.V2{5, 0}, vec.V2{5, 5}, vec.V2{2, 2})
	fmt.Printf("%+v\n", x)

	x = barycentric(vec.V2{0, 0}, vec.V2{1, 1}, vec.V2{1, 1}, vec.V2{2, 2})
	fmt.Printf("%+v\n", x)

	x = barycentric(vec.V2{1, 1}, vec.V2{1, 1}, vec.V2{1, 1}, vec.V2{2, 2})
	fmt.Printf("%+v\n", x)

	x = barycentric(vec.V2{1, 1}, vec.V2{2, 2}, vec.V2{3, 3}, vec.V2{2, 2})
	fmt.Printf("%+v\n", x)

	x = barycentric(vec.V2{0, 0}, vec.V2{2, 0}, vec.V2{0, 2}, vec.V2{k, k})
	fmt.Printf("%+v\n", x)

	x = barycentric(vec.V2{1, 1}, vec.V2{3, 1}, vec.V2{1, 3}, vec.V2{1 + k, 1 + k})
	fmt.Printf("%+v\n", x)

}

func triangle2(v [3]*vec.V2, img *image.NRGBA, color color.NRGBA) {

}
