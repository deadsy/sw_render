package main

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"github.com/deadsy/sw_render/utils"
	"github.com/deadsy/sw_render/vec"
	"github.com/disintegration/imaging"
)

func collinear(a, b, c vec.V2i) bool {
	// if the slopes of the line segments are the same the points are collinear
	return (c[1]-b[1])*(b[0]-a[0]) == (c[0]-b[0])*(b[1]-a[1])
}

type line_func func(int, int, int)

// filled in triangle - horizontal raster between 2 bresenham lines
// dy > 0, dx0 <= dx1
func bresenham_triangle(dy, dx0, dx1 int, line line_func) {

	var x0, x1, err0, err1 int

	x0_inc := utils.Sgn(dx0)
	x1_inc := utils.Sgn(dx1)

	dx0 = utils.Abs(dx0)
	dx1 = utils.Abs(dx1)

	err_y := 2 * dy
	err_x0 := 2 * dx0
	err_x1 := 2 * dx1

	for y := 0; y < dy; y++ {
		line(x0, x1, y)
		// line 0
		err0 += err_x0
		for err0 >= dy {
			x0 += x0_inc
			err0 -= err_y
		}
		// line 1
		err1 += err_x1
		for err1 >= dy {
			x1 += x1_inc
			err1 -= err_y
		}
	}
}

// for flat top triangles
func line_ft(ofs vec.V2i, img *image.NRGBA, color color.NRGBA) line_func {
	return func(x0, x1, y int) {
		for x := x0; x <= x1; x++ {
			img.SetNRGBA(ofs[0]+x, ofs[1]+y, color)
		}
	}
}

// for flat bottom triangles (inverted from normal case)
func line_fb(ofs vec.V2i, img *image.NRGBA, color color.NRGBA) line_func {
	return func(x0, x1, y int) {
		for x := x0; x <= x1; x++ {
			img.SetNRGBA(ofs[0]+x, ofs[1]-y, color)
		}
	}
}

func triangle(a, b, c vec.V2i, img *image.NRGBA, color color.NRGBA) {

	if collinear(a, b, c) {
		// line or point - skip it
		fmt.Printf("collinear %+v %+v %+v\n", a, b, c)
		return
	}

	// bubble sort the vertices by y-order
	p := [3]*vec.V2i{&a, &b, &c}
	if p[0][1] > p[1][1] {
		// swap p[0] with p[1]
		x := p[1]
		p[1] = p[0]
		p[0] = x
	}
	if p[1][1] > p[2][1] {
		// swap p[1] with p[2]
		x := p[2]
		p[2] = p[1]
		p[1] = x
	}
	if p[0][1] > p[1][1] {
		// swap p[0] with p[1]
		x := p[1]
		p[1] = p[0]
		p[0] = x
	}

	if p[0][1] == p[1][1] {
		// flat bottom triangle
		fmt.Printf("flat bottom: %+v %+v %+v\n", *p[0], *p[1], *p[2])
		// TODO
		return
	}

	if p[1][1] == p[2][1] {
		// flat top triangle
		fmt.Printf("flat top: %+v %+v %+v\n", *p[0], *p[1], *p[2])
		// TODO
		return
	}

	fmt.Printf("general: %+v %+v %+v\n", *p[0], *p[1], *p[2])
	// TODO
}

func main() {

	imgfile := "output.png"

	white := color.NRGBA{255, 255, 255, 255}
	black := color.NRGBA{0, 0, 0, 255}
	img := imaging.New(800, 1000, black)

	// collinear
	triangle(vec.V2i{10, 20}, vec.V2i{10, 20}, vec.V2i{10, 20}, img, white)
	triangle(vec.V2i{10, -30}, vec.V2i{10, 20}, vec.V2i{10, -10}, img, white)
	triangle(vec.V2i{-10, -30}, vec.V2i{50, -30}, vec.V2i{-17, -30}, img, white)
	triangle(vec.V2i{0, 30}, vec.V2i{0, 30}, vec.V2i{10, 20}, img, white)
	triangle(vec.V2i{0, 30}, vec.V2i{0, 20}, vec.V2i{0, 10}, img, white)
	triangle(vec.V2i{10, 20}, vec.V2i{20, 40}, vec.V2i{30, 60}, img, white)
	triangle(vec.V2i{-10, -20}, vec.V2i{20, 40}, vec.V2i{13, 26}, img, white)

	// flat bottom
	triangle(vec.V2i{10, 30}, vec.V2i{20, 30}, vec.V2i{100, 200}, img, white)

	// flat top
	triangle(vec.V2i{10, 30}, vec.V2i{20, 30}, vec.V2i{0, 0}, img, white)

	// general
	triangle(vec.V2i{10, 20}, vec.V2i{20, 40}, vec.V2i{30, 61}, img, white)

	bresenham_triangle(95, -395, -200, line_ft(vec.V2i{400, 0}, img, white))
	bresenham_triangle(95, 200, 395, line_ft(vec.V2i{400, 100}, img, white))

	bresenham_triangle(95, -300, -80, line_ft(vec.V2i{400, 200}, img, white))
	bresenham_triangle(95, 80, 300, line_ft(vec.V2i{400, 300}, img, white))

	bresenham_triangle(95, -73, 0, line_ft(vec.V2i{400, 400}, img, white))
	bresenham_triangle(95, 0, 73, line_ft(vec.V2i{400, 500}, img, white))

	bresenham_triangle(95, -300, 200, line_ft(vec.V2i{400, 600}, img, white))
	bresenham_triangle(95, -200, 300, line_ft(vec.V2i{400, 700}, img, white))

	bresenham_triangle(95, -80, 71, line_ft(vec.V2i{400, 800}, img, white))
	bresenham_triangle(95, -71, 81, line_ft(vec.V2i{400, 900}, img, white))

	img = imaging.FlipV(img)
	err := imaging.Save(img, imgfile)
	if err != nil {
		fmt.Printf("unable to save %s, %s\n", imgfile, err)
		os.Exit(1)
	}

	os.Exit(0)
}
