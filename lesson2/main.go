package main

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"os"

	"github.com/deadsy/sw_render/utils"
	"github.com/deadsy/sw_render/vec"
	"github.com/deadsy/sw_render/wavefront"
	"github.com/disintegration/imaging"
)

func collinear(a, b, c vec.V2i) bool {
	// if the slopes of the line segments are the same the points are collinear
	return (c[1]-b[1])*(b[0]-a[0]) == (c[0]-b[0])*(b[1]-a[1])
}

type line_func func(int, int, int)

// filled in triangle - horizontal raster between 2 bresenham lines
// dy > 0
func bresenham_triangle(dy, dx0, dx1 int, line line_func) {

	var x0, x1, err0, err1 int

	if dx0 > dx1 {
		// swap them
		tmp := dx1
		dx1 = dx0
		dx0 = tmp
	}

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
	vec.Sort_Y(p[0:2])
	vec.Sort_Y(p[1:3])
	vec.Sort_Y(p[0:2])

	if p[0][1] == p[1][1] {
		// flat bottom triangle
		fmt.Printf("flat bottom: %+v %+v %+v\n", *p[0], *p[1], *p[2])
		dy := p[2][1] - p[0][1]
		dx0 := p[0][0] - p[2][0]
		dx1 := p[1][0] - p[2][0]
		bresenham_triangle(dy, dx0, dx1, line_fb(*p[2], img, color))
		return
	}

	if p[1][1] == p[2][1] {
		// flat top triangle
		fmt.Printf("flat top: %+v %+v %+v\n", *p[0], *p[1], *p[2])
		dy := p[1][1] - p[0][1]
		dx0 := p[1][0] - p[0][0]
		dx1 := p[2][0] - p[0][0]
		bresenham_triangle(dy, dx0, dx1, line_ft(*p[0], img, color))
		return
	}

	fmt.Printf("general: %+v %+v %+v\n", *p[0], *p[1], *p[2])
	// work out the x-coordinate of the unknown mid-point
	k := p[0][0] + utils.Round(float32((p[1][1]-p[0][1])*(p[2][0]-p[0][0]))/float32(p[2][1]-p[0][1]))

	fmt.Printf("K %d\n", k)

	// flat bottom triangle
	dy := p[2][1] - p[1][1]
	dx0 := p[1][0] - p[2][0]
	dx1 := k - p[2][0]
	bresenham_triangle(dy, dx0, dx1, line_fb(*p[2], img, color))

	// flat top triangle
	dy = p[1][1] - p[0][1]
	dx0 = p[1][0] - p[0][0]
	dx1 = k - p[0][0]
	bresenham_triangle(dy, dx0, dx1, line_ft(*p[0], img, color))
}

func Random_Color() color.NRGBA {
	return color.NRGBA{
		uint8(256 * rand.Float32()),
		uint8(256 * rand.Float32()),
		uint8(256 * rand.Float32()),
		255, //uint8(256 * rand.Float32()),
	}
}

func random_triangles(k vec.V2i, img *image.NRGBA) {
	for i := 0; i < 200; i++ {
		a := k.Rand()
		b := a.Rand_Delta(100)
		c := a.Rand_Delta(100)
		triangle(a, b, c, img, Random_Color())
	}
}

const pixels_x = 750
const pixels_ofs = 5

// object to image mapping
func Obj2Img(v, ofs vec.V3f, scale float32) vec.V2i {
	p := v.Sum(ofs).Scale(scale)
	return vec.V2i{int(p[0]), int(p[1])}
}

func main() {

	imgfile := "output.png"
	objfile := "../obj/african_head.obj"

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
	img_size = img_size.Sum(vec.V3f{pixels_ofs, pixels_ofs, pixels_ofs})
	fmt.Printf("img_size: %+v\n", img_size)

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

		triangle(p0, p1, p2, img, Random_Color())
	}

	//random_triangles(k, img)

	img = imaging.FlipV(img)
	err = imaging.Save(img, imgfile)
	if err != nil {
		fmt.Printf("unable to save %s, %s\n", imgfile, err)
		os.Exit(1)
	}

	os.Exit(0)
}
