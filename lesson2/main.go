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

	for y := 0; y <= dy; y++ {
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

func line_simple(ofs vec.V2i, img *image.NRGBA, color color.NRGBA) line_func {
	return func(x0, x1, y int) {
		for x := x0; x <= x1; x++ {
			img.SetNRGBA(ofs[0]+x, ofs[1]+y, color)
		}
	}
}

func main() {

	imgfile := "output.png"

	white := color.NRGBA{255, 255, 255, 255}
	black := color.NRGBA{0, 0, 0, 255}
	img := imaging.New(800, 1000, black)

	bresenham_triangle(95, -395, -200, line_simple(vec.V2i{400, 0}, img, white))
	bresenham_triangle(95, 200, 395, line_simple(vec.V2i{400, 100}, img, white))

	bresenham_triangle(95, -300, -80, line_simple(vec.V2i{400, 200}, img, white))
	bresenham_triangle(95, 80, 300, line_simple(vec.V2i{400, 300}, img, white))

	bresenham_triangle(95, -73, 0, line_simple(vec.V2i{400, 400}, img, white))
	bresenham_triangle(95, 0, 73, line_simple(vec.V2i{400, 500}, img, white))

	bresenham_triangle(95, -300, 200, line_simple(vec.V2i{400, 600}, img, white))
	bresenham_triangle(95, -200, 300, line_simple(vec.V2i{400, 700}, img, white))

	bresenham_triangle(95, -80, 71, line_simple(vec.V2i{400, 800}, img, white))
	bresenham_triangle(95, -71, 81, line_simple(vec.V2i{400, 900}, img, white))

	img = imaging.FlipV(img)
	err := imaging.Save(img, imgfile)
	if err != nil {
		fmt.Printf("unable to save %s, %s\n", imgfile, err)
		os.Exit(1)
	}

	os.Exit(0)
}
