package main

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"github.com/deadsy/sw_render/vec"
	"github.com/deadsy/sw_render/wavefront"
	"github.com/disintegration/imaging"
)

const pixels_x = 1000
const pixels_ofs = 5

type V2 [2]int

func (a V2) Equal(b V2) bool {
	return (a[0] == b[0]) && (a[1] == b[1])
}

// Return a - b
func (a V2) Sub(b V2) V2 {
	return V2{
		a[0] - b[0],
		a[1] - b[1],
	}
}

func abs(a int) int {
	if a > 0 {
		return a
	}
	return -a
}

type plot_func func(int, int)

// bresenham's line algorithm
// dx > 0, dy >= 0, dx >= dy
func bresenham_line(dx, dy int, plot plot_func) {
	err_y := 2 * dy
	err_x := 2 * dx
	y := 0
	err := 0
	for x := 0; x <= dx; x++ {
		plot(x, y)
		err += err_y
		if err >= dx {
			y += 1
			err -= err_x
		}
	}
}

// major x-axis, quadrant 0
func plot_x0(ofs V2, img *image.NRGBA, color color.NRGBA) plot_func {
	return func(x, y int) {
		img.SetNRGBA(ofs[0]+x, ofs[1]+y, color)
	}
}

// major x-axis, quadrant 1
func plot_x1(ofs V2, img *image.NRGBA, color color.NRGBA) plot_func {
	return func(x, y int) {
		img.SetNRGBA(ofs[0]-x, ofs[1]+y, color)
	}
}

// major x-axis, quadrant 2
func plot_x2(ofs V2, img *image.NRGBA, color color.NRGBA) plot_func {
	return func(x, y int) {
		img.SetNRGBA(ofs[0]-x, ofs[1]-y, color)
	}
}

// major x-axis, quadrant 3
func plot_x3(ofs V2, img *image.NRGBA, color color.NRGBA) plot_func {
	return func(x, y int) {
		img.SetNRGBA(ofs[0]+x, ofs[1]-y, color)
	}
}

// major y-axis, quadrant 0
func plot_y0(ofs V2, img *image.NRGBA, color color.NRGBA) plot_func {
	return func(y, x int) {
		img.SetNRGBA(ofs[0]+x, ofs[1]+y, color)
	}
}

// major y-axis, quadrant 1
func plot_y1(ofs V2, img *image.NRGBA, color color.NRGBA) plot_func {
	return func(y, x int) {
		img.SetNRGBA(ofs[0]-x, ofs[1]+y, color)
	}
}

// major y-axis, quadrant 2
func plot_y2(ofs V2, img *image.NRGBA, color color.NRGBA) plot_func {
	return func(y, x int) {
		img.SetNRGBA(ofs[0]-x, ofs[1]-y, color)
	}
}

// major y-axis, quadrant 3
func plot_y3(ofs V2, img *image.NRGBA, color color.NRGBA) plot_func {
	return func(y, x int) {
		img.SetNRGBA(ofs[0]+x, ofs[1]-y, color)
	}
}

func line(a, b V2, img *image.NRGBA, color color.NRGBA) {
	if a.Equal(b) {
		return
	}
	x := b.Sub(a)
	if abs(x[0]) >= abs(x[1]) {
		// major x-axis
		if x[0] >= 0 {
			if x[1] >= 0 {
				bresenham_line(x[0], x[1], plot_x0(a, img, color))
			} else {
				bresenham_line(x[0], -x[1], plot_x3(a, img, color))
			}
		} else {
			if x[1] >= 0 {
				bresenham_line(-x[0], x[1], plot_x1(a, img, color))
			} else {
				bresenham_line(-x[0], -x[1], plot_x2(a, img, color))
			}
		}
	} else {
		// major y-axis
		if x[0] >= 0 {
			if x[1] >= 0 {
				bresenham_line(x[1], x[0], plot_y0(a, img, color))
			} else {
				bresenham_line(-x[1], x[0], plot_y3(a, img, color))
			}
		} else {
			if x[1] >= 0 {
				bresenham_line(x[1], -x[0], plot_y1(a, img, color))
			} else {
				bresenham_line(-x[1], -x[0], plot_y2(a, img, color))
			}
		}
	}
}

// object to image mapping
func Obj2Img(v, ofs vec.V3, scale float32) V2 {
	p := v.Sum(ofs).Scale(scale)
	return V2{int(p[0]), int(p[1])}
}

func main() {

	//objfile := "../obj/gopher.obj"
	objfile := "../obj/african_head.obj"
	//objfile := "../obj/test_triangle.obj"
	imgfile := "output.png"

	obj, err := wavefront.Read(objfile)
	if err != nil {
		fmt.Printf("%s: %s\n", objfile, err)
		os.Exit(1)
	}

	fmt.Printf("%s\n", obj)

	obj_ofs := obj.Offset()
	fmt.Printf("obj_ofs: %+v\n", obj_ofs)

	obj_range := obj.Range()
	fmt.Printf("obj_range: %+v\n", obj_range)

	// we want the image to be pixels_x wide on the x-axis
	// work out the scaling factor
	scale := (pixels_x - pixels_ofs) / obj_range[0]

	// work out the image size
	img_size := obj_range.Scale(scale)
	img_size = img_size.Sum(vec.V3{pixels_ofs, pixels_ofs, pixels_ofs})
	fmt.Printf("img_size: %+v\n", img_size)

	white := color.NRGBA{255, 255, 255, 255}
	black := color.NRGBA{0, 0, 0, 255}
	img := imaging.New(int(img_size[0]), int(img_size[1]), black)

	// iterate over the object faces
	for i := 0; i < obj.Len_F(); i++ {

		// get the vertices from the face
		v0 := obj.Get_V(i, 0).ToV3()
		v1 := obj.Get_V(i, 1).ToV3()
		v2 := obj.Get_V(i, 2).ToV3()

		p0 := Obj2Img(v0, obj_ofs, scale)
		p1 := Obj2Img(v1, obj_ofs, scale)
		p2 := Obj2Img(v2, obj_ofs, scale)

		// p0 to p1
		line(p0, p1, img, white)
		// p1 to p2
		line(p1, p2, img, white)
		// p2 to p0
		line(p2, p0, img, white)
	}

	img = imaging.FlipV(img)
	err = imaging.Save(img, imgfile)
	if err != nil {
		fmt.Printf("unable to save %s, %s\n", imgfile, err)
		os.Exit(1)
	}

	os.Exit(0)
}
